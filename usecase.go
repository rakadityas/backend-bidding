package main

import (
	"context"

	"github.com/go-redis/redis"
)

const (
	RedisKey = "BID"
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

	auctionData, err := GetAuctionDB(ctx, request.ProductID)
	if err != nil {
		return GetAuctionDetailResponse{}, err
	}

	productData, err := GetProductDB(ctx, request.UserID)
	if err != nil {
		return GetAuctionDetailResponse{}, err
	}

	return GetAuctionDetailResponse{
		ProductDetail: ProductDetail{
			Product: productData,
			Auction: auctionData,
		},
	}, nil
}

func AuctionBidding(ctx context.Context, payload AuctionBidRequest) (response AuctionBidResponse, err error) {

	// get auction
	auctionData, err := GetAuctionDB(ctx, payload.ProductID)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	userData, err := GetUserInfo(ctx, payload.UserID)
	if err != nil {
		return AuctionBidResponse{}, err
	}

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

	// deduct balance
	err = UpdateBalance(ctx, userData.ID, bidAmount)
	if err != nil {
		return AuctionBidResponse{}, err
	}

	err = DoPublishNSQ("Update_Scoreboard", UpdateScoreboardNSQ{
		UserID:    payload.UserID,
		BidAmount: bidAmount,
	})
	if err != nil {
		return AuctionBidResponse{}, err
	}

	err = DoPublishNSQ("Insert_Collection_And_Payment", InsertPaymentAndBidCollectionNSQ{
		Payment: Payment{
			UserID: payload.UserID,
			Amount: bidAmount,
			Status: 1,
		},
		BidCollection: BidCollection{
			UserID:     payload.UserID,
			AuctionID:  auctionData.ID,
			CurrentBid: bidAmount,
			// PaymentID: ,
		},
	})
	if err != nil {
		return AuctionBidResponse{}, err
	}

	return AuctionBidResponse{
		ResultStatus: ResultStatus{
			IsSuccess: true,
			Message:   "Success bidding!",
		},
	}, nil
}

func CheckHighestBid(ctx context.Context, bid int64, userID int64) {

	var (
		highestBid int64
	)

	response := RedisClient.ZRevRangeWithScores(RedisKey, 0, 0)

	for _, val := range response.Val() {
		highestBid = int64(val.Score)
	}

	if bid > highestBid {
		RedisClient.ZAdd(RedisKey,
			redis.Z{
				Score:  float64(bid),
				Member: userID,
			})
	}

	return
}

func InsertBidCollectionAndPayment(ctx context.Context, bidCollection BidCollection, payment Payment) error {

	var (
		paymentID int64
		err       error
	)

	paymentID, err = InsertPayment(ctx, payment)
	if err != nil {
		return err
	}

	bidCollection.PaymentID = paymentID
	err = InsertBidCollection(ctx, bidCollection)
	if err != nil {
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
