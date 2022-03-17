package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	redisKey = `bid-%d`

	AuctionStatusUnactive    = 0
	AuctionStatusActive      = 1
	AuctionStatusDone        = 2
	AuctionStatusDeactivated = 10
)

// core logic
func GetUserInfo(ctx context.Context, userID int64) (User, error) {
	userData, err := GetUserInfoDB(ctx, userID)
	if err != nil {
		return User{}, err
	}

	return userData, nil
}

func GetAuctionDetail(ctx context.Context, request GetAuctionDetailRequest) (GetAuctionDetailResponse, error) {

	var (
		highestBidder User
		highestBid    int64
		userID        int64
		err           error
	)

	auctionData, err := GetAuctionDB(ctx, request.ProductID)
	if err != nil {
		return GetAuctionDetailResponse{}, err
	}

	productData, err := GetProductDB(ctx, request.ProductID)
	if err != nil {
		return GetAuctionDetailResponse{}, err
	}

	// get timewindow
	timeWindow, err := GetTimeWindowDB(ctx, auctionData.ID)
	if err != nil {
		return GetAuctionDetailResponse{}, err
	}

	resp := GetHighestBid(ctx, auctionData.ID)
	if resp != nil {
		for _, val := range resp.Val() {
			highestBid = int64(val.Score)
			userID, _ = itfToInt64(val.Member)
		}
		if userID > 0 {
			highestBidder, err = GetUserInfoDB(ctx, userID)
			if err != nil {
				return GetAuctionDetailResponse{}, err
			}
		}
	}

	return GetAuctionDetailResponse{
		ProductDetail: ProductDetail{
			Product:       productData,
			Auction:       auctionData,
			HighestBidder: highestBidder,
			HighestBid:    highestBid,
			Countdown:     (timeWindow.EndTime.UnixMilli() - time.Now().UnixMilli()),
		},
	}, nil
}

func AuctionBidding(ctx context.Context, payload AuctionBidRequest) (response AuctionBidResponse, err error) {

	// Get auction
	auctionData, err := GetAuctionDB(ctx, payload.ProductID)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	// Check for auction status; Auction should be deactivated if TW expired w cron
	if auctionData.Status == AuctionStatusDeactivated {
		return AuctionBidResponse{
			ResultStatus: ResultStatus{
				IsSuccess: false,
				Message:   "Time's up guyz",
			},
		}, nil
	}

	// Multiplier validation
	if (payload.Amount % auctionData.Multiplier) != 0 {
		return AuctionBidResponse{
			ResultStatus: ResultStatus{
				IsSuccess: false,
				Message:   fmt.Sprintf("Hanya untuk kelipatan %d!", auctionData.Multiplier),
			},
		}, nil
	}

	// Get userinfo
	userData, err := GetUserInfo(ctx, payload.UserID)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	// Balance validation
	if userData.Balance <= 0 {
		return AuctionBidResponse{
			ResultStatus: ResultStatus{
				IsSuccess: false,
				Message:   "Top up dulu bang!",
			},
		}, nil
	}

	// Get total sum of prev bid(s)
	sumBid, err := GetSumBidCollection(ctx, payload.UserID, auctionData.ID)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	bidAmount := (payload.Amount - sumBid)

	// check balance query ke user get blanace (current bid - bidded balance)
	deductedBalance := userData.Balance - bidAmount
	if deductedBalance < 0 {
		return AuctionBidResponse{
			ResultStatus: ResultStatus{
				IsSuccess: false,
				Message:   "Top up dulu bang!",
			},
		}, nil
	}

	fmt.Println("SUMBID: ", sumBid)
	fmt.Println("DDDD: ", deductedBalance)
	fmt.Println("BID AMOUNT: ", bidAmount)

	// deduct balance
	err = UpdateBalance(ctx, userData.ID, deductedBalance)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	// NSQ NOT WORKING
	// err = DoPublishNSQ("Update_Scoreboard", UpdateScoreboardNSQ{
	// 	UserID:    payload.UserID,
	// 	BidAmount: bidAmount,
	// })
	// if err != nil {
	// 	return AuctionBidResponse{}, err
	// }

	// err = DoPublishNSQ("Insert_Collection_And_Payment", InsertPaymentAndBidCollectionNSQ{
	// 	Payment: Payment{
	// 		UserID: payload.UserID,
	// 		Amount: bidAmount,
	// 		Status: 1,
	// 	},
	// 	BidCollection: BidCollection{
	// 		UserID:     payload.UserID,
	// 		AuctionID:  auctionData.ID,
	// 		CurrentBid: bidAmount,
	// 		// PaymentID: ,
	// 	},
	// })
	// if err != nil {
	// 	return AuctionBidResponse{}, err
	// }

	// this is supposed to use NSQ
	err = InsertBidCollectionAndPayment(ctx, BidCollection{
		UserID:     payload.UserID,
		AuctionID:  auctionData.ID,
		CurrentBid: bidAmount,
	}, Payment{
		UserID: payload.UserID,
		Amount: bidAmount,
		Status: 1,
	})
	if err != nil {
		fmt.Println(err)
		return AuctionBidResponse{
			ResultStatus: ResultStatus{
				IsSuccess: true,
				Message:   "Maaf terjadi kendala",
			},
		}, nil
	}

	// this is supposed to use NSQ with max in flight 1 Update_Scoreboard
	CheckHighestBid(ctx, payload.Amount, payload.UserID, auctionData.ID)

	return AuctionBidResponse{
		ResultStatus: ResultStatus{
			IsSuccess: true,
			Message:   "Success bidding!",
		},
	}, nil
}

func CheckHighestBid(ctx context.Context, bid int64, userID int64, auctionID int64) {

	var (
		highestBid int64
	)

	response := GetHighestBid(ctx, auctionID)
	if response == nil {
		fmt.Println("error GetHighestBid: nil")
		return
	}

	for _, val := range response.Val() {
		highestBid = int64(val.Score)
	}

	if bid > highestBid {
		cmd := RedisClient.ZAdd(fmt.Sprintf(redisKey, auctionID),
			redis.Z{
				Score:  float64(bid),
				Member: userID,
			})

		if cmd.Err() != nil {
			fmt.Println(cmd.Err().Error())
			return
		}
	}

	return
}

func GetHighestBid(ctx context.Context, auctionID int64) *redis.ZSliceCmd {
	key := fmt.Sprintf(redisKey, auctionID)
	return RedisClient.ZRevRangeWithScores(key, 0, 0)
}

func InsertBidCollectionAndPayment(ctx context.Context, bidCollection BidCollection, payment Payment) error {

	var (
		paymentID int64
		err       error
	)

	paymentID, err = InsertPayment(ctx, payment)
	if err != nil {
		fmt.Println(err)
		return err
	}

	bidCollection.PaymentID = paymentID
	err = InsertBidCollection(ctx, bidCollection)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Login(ctx context.Context, username, password string) (User, error) {

	user, err := GetUser(ctx, username, password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func InsertAuction(ctx context.Context, auctionRequest CreateAuctionRequest) (ResultStatus, error) {

	product := Product{
		UserID:      auctionRequest.UserID,
		ProductName: auctionRequest.ProductName,
		ImageURL:    auctionRequest.ProductImageURL,
	}

	productID, err := InsertProductDB(ctx, product)
	if err != nil {
		fmt.Println("got error: ", err)
		return ResultStatus{
			Message:   "Terjadi kesalahan",
			IsSuccess: false,
		}, nil
	}

	auction := Auction{
		ProductID:  productID,
		Multiplier: auctionRequest.Multiplier,
		Status:     AuctionStatusActive, // harusnya ada cron buat activate; pas dibuat status unactivated
	}

	auctionID, err := InsertAuctionDB(ctx, auction)
	if err != nil {
		fmt.Println("got error: ", err)
		return ResultStatus{
			Message:   "Terjadi kesalahan",
			IsSuccess: false,
		}, nil
	}

	timeWindow := TimeWindow{
		AuctionID: auctionID,
		StartTime: auctionRequest.StartTime,
		EndTime:   auctionRequest.EndTime,
	}
	err = InsertTWDB(ctx, timeWindow)
	if err != nil {
		fmt.Println("got error: ", err)
		return ResultStatus{
			Message:   "Terjadi kesalahan",
			IsSuccess: false,
		}, nil
	}

	return ResultStatus{
		Message:   "Sukses",
		IsSuccess: true,
	}, nil
}

func GetAuctionList(ctx context.Context, userID int64, sortAsc bool) (response GetAuctionListResponse, err error) {

	switch userID > 0 {
	case true: // represent seller flow; will get based on user id
		response, err = GetAuctionListSeller(ctx, userID)
	default: // represent buyer flow; will get all
		response, err = GetAuctionListBuyer(ctx)
	}
	if err != nil {
		return GetAuctionListResponse{}, err
	}

	return response, nil
}

func GetAuctionListSeller(ctx context.Context, userID int64) (response GetAuctionListResponse, err error) {
	products, err := GetProductByUserIDDB(ctx, userID)
	if err != nil {
		return GetAuctionListResponse{}, err
	}

	for idx := range products {
		auction, err := GetAuctionDB(ctx, products[idx].ID)
		if err != nil {
			return GetAuctionListResponse{}, err
		}

		var (
			userID        int64
			highestBid    int64
			highestBidder User
		)

		resp := GetHighestBid(ctx, auction.ID)
		if resp != nil {
			for _, val := range resp.Val() {
				highestBid = int64(val.Score)
				userID, _ = itfToInt64(val.Member)
			}
			if userID > 0 {
				highestBidder, err = GetUserInfoDB(ctx, userID)
				if err != nil {
					return GetAuctionListResponse{}, err
				}
			}
		}

		response.ProductDetail = append(response.ProductDetail, ProductDetail{
			Product:       products[idx],
			Auction:       auction,
			HighestBidder: highestBidder,
			HighestBid:    highestBid,
		})
	}

	return response, nil
}

func itfToInt64(t interface{}) (int64, error) {
	switch t := t.(type) { // This is a type switch.
	case int64:
		return t, nil // All done if we got an int64.
	case int:
		return int64(t), nil // This uses a conversion from int to int64
	case string:
		return strconv.ParseInt(t, 10, 64)
	default:
		return 0, errors.New("data type invalid/unknown")
	}
}

func GetAuctionListBuyer(ctx context.Context) (response GetAuctionListResponse, err error) {
	// get all auction
	auctions, err := GetAllAuction(ctx)

	fmt.Println(err, " 1 //")
	if err != nil {
		return GetAuctionListResponse{}, err
	}

	// TODO: improve this logic
	for idx := range auctions {
		product, err := GetProductDB(ctx, auctions[idx].ProductID)
		fmt.Println(err, " 2")

		if err != nil {
			return GetAuctionListResponse{}, err
		}

		var (
			userID        int64
			highestBid    int64
			highestBidder User
		)

		resp := GetHighestBid(ctx, auctions[idx].ID)
		if resp != nil {
			for _, val := range resp.Val() {
				highestBid = int64(val.Score)
				userID, _ = itfToInt64(val.Member)
				fmt.Println("uid:" + fmt.Sprint(userID))
			}
			if userID > 0 {
				highestBidder, err = GetUserInfoDB(ctx, userID)
				if err != nil {
					return GetAuctionListResponse{}, err
				}
			}
		}

		response.ProductDetail = append(response.ProductDetail, ProductDetail{
			Product:       product,
			Auction:       auctions[idx],
			HighestBidder: highestBidder,
			HighestBid:    highestBid,
		})
	}

	return response, nil
}

func IsMatchError(err1 error, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}

	if err1 == nil {
		err1 = errors.New("nil")
	}

	if err2 == nil {
		err2 = errors.New("nil")
	}

	// for now comparing the message only, because if comparing errors will panic
	if err1.Error() == err2.Error() {
		return true
	}

	return false
}
