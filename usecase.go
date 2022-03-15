package main

import "context"

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
