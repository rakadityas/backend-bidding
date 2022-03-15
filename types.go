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
	ID       int64
	UserType int
	Username string
	Status   int
	Balance  int64
}

// Product table
type Product struct {
	ID          int64
	UserID      int64
	ProductName string
	ImageURL    string
	Status      int
}

// ProductDetail table
type ProductDetail struct {
	Product Product
	Auction Auction
}

// TimeWindow table
type TimeWindow struct {
	ID        int64
	StartTime time.Time
	EndTime   time.Time
	Status    int
}

// Auction table
type Auction struct {
	ID           int64
	ProductID    int64
	WinnerUserID int64
	Multiplier   int64
	Status       int
}

// BidCollection table
type BidCollection struct {
	ID         int64
	UserID     int64
	AuctionID  int64
	CurrentBid int64
	PaymentID  int64
}

// Payment table
type Payment struct {
	ID     int64
	UserID int64
	Amount int64
	Status int
}

// /get/auction/detail
type (
	GetAuctionDetailRequest struct {
		ProductID int64
		UserID    int64
	}

	GetAuctionDetailResponse struct {
		ProductDetail ProductDetail
	}
)

// /get/user/info
type (
	GetUserInfoRequest struct {
		UserID int64
	}

	GetUserInfoResposne struct {
		UserInfo User
	}
)

// /get/auction/list
type (
	GetAuctionListRequest struct {
		UserID  int64
		SortAsc bool
	}

	GetAuctionListResponse struct {
		ProductDetail []ProductDetail
	}
)

// /login
type (
	LoginRequest struct {
		Username string
		Password string
	}

	LoginResponse struct {
		Status int
	}
)

// /auction/create
type (
	CreateAuctionRequest struct {
		ProductName     string
		ProductImageURL string
		StartBid        int64
		Multiplier      int
		Date            time.Time
		ShopID          string
		UserID          string
	}

	CreateAuctionResponse struct {
		ResultStatus ResultStatus
	}
)

// /auction/bid
type (
	AuctionBidRequest struct {
		UserID    int64
		ProductID int64
		Amount    int64
	}

	AuctionBidResponse struct {
		ResultStatus ResultStatus
	}
)

type (
	ResultStatus struct {
		Message   string
		IsSuccess bool
	}
)
