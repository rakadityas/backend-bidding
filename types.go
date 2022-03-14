package main

import "time"

// Test
const (
	UserTypeBidder = 0
	UserTypeSeller = 1

	UserStatusTest = 0
)

// User table
type User struct {
	ID       int64  `json:"id"`
	UserType int    `json:"user_type"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	Balance  int64  `json:"balance"`
}

// Product table
type Product struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	ProductName string `json:"product_name"`
	ImageURL    string `json:"image_url"`
	Status      int    `json:"status"`
}

// ProductDetail table
type ProductDetail struct {
}

// TimeWindow table
type TimeWindow struct {
	ID            int64     `json:"id"`
	ReferenceID   int64     `json:"reference_id"`
	ReferenceType int       `json:"reference_type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        int       `json:"status"`
}

// Auction table
type Auction struct {
	ID           int64 `json:"id"`
	ProductID    int64 `json:"product_id"`
	WinnerUserID int64 `json:"winner_user_id"`
	Multiplier   int64 `json:"multiplier"`
	Status       int   `json:"status"`
}

// BidCollection table
type BidCollection struct {
	ID         int64 `json:"id"`
	UserID     int64 `json:"user_id"`
	AuctionID  int64 `json:"auction_id"`
	CurrentBid int64 `json:"current_bid"`
	PaymentID  int64 `json:"payment_id"`
}

// Payment table
type Payment struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
	Status int   `json:"status"`
}

type UserInfo struct {
}

type GetAuctionDetailRequest struct {
	ProductID int64 `json:product_id`
	UserID    int64 `json:user_id`
}

type GetUserInfoRequest struct {
	UserID int64 `json:user_id`
}

type GetAuctionListRequest struct {
	UserID  int64 `json:user_id`
	SortAsc bool  `json:sort_asc`
}

type LoginRequest struct {
	Username string `json:username`
	Password string `json:password`
}

type CreateAuctionRequest struct {
	ProductName     string    `json:product_name`
	ProductImageUrl string    `json:product_image_url`
	StartBid        int64     `json:start_bid`
	Multiplier      int       `json:multiplier`
	Date            time.Time `json:date`
	ShopID          string    `json:shop_id`
	UserID          string    `json:user_id`
}

type GetAuctionDetailResponse struct {
	ProductDetail ProductDetail `json:product_detail`
	Auction       Auction       `json:auction`
}

type GetUserInfoResposne struct {
	UserInfo UserInfo `json:user_info`
}

type GetAuctionListResponse struct {
}

type LoginResponse struct {
}

type CreateAuctionResponse struct {
}
