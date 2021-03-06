package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var db *sql.DB

func initDB() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

func GetUserInfoDB(ctx context.Context, userID int64) (userData User, err error) {
	err = db.QueryRowContext(ctx, queryGetUserInfoDB, userID).Scan(
		&userData.ID,
		&userData.UserType,
		&userData.Username,
		&userData.Status,
		&userData.Balance,
	)
	if err != nil {
		return User{}, err
	}

	return userData, nil
}

func GetAuctionDB(ctx context.Context, productID int64) (auction Auction, err error) {

	err = db.QueryRowContext(ctx, queryGetAuction, productID).Scan(
		&auction.ID,
		&auction.ProductID,
		&auction.WinnerUserID,
		&auction.InitialBid,
		&auction.HighestBid,
		&auction.Multiplier,
		&auction.Status,
	)
	if err != nil {
		return auction, err
	}

	return auction, nil
}

func GetProductDB(ctx context.Context, productID int64) (product Product, err error) {

	err = db.QueryRowContext(ctx, queryGetProduct, productID).Scan(
		&product.ID,
		&product.UserID,
		&product.ProductName,
		&product.ImageURL,
		&product.Status,
	)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func GetProductByUserIDDB(ctx context.Context, userID int64) (products []Product, err error) {

	rows, err := db.QueryContext(ctx, queryGetProductByUserID, userID)
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		rows.Scan(
			&product.ID,
			&product.UserID,
			&product.ProductName,
			&product.ImageURL,
			&product.Status,
		)

		products = append(products, product)
	}

	return products, nil
}

func GetSumBidCollection(ctx context.Context, userID, auctionID int64) (highestBid int64, err error) {
	err = db.QueryRowContext(ctx, queryGetSumBidCollection, userID, auctionID).Scan(
		&highestBid,
	)
	if err != nil && !IsMatchError(err, sql.ErrNoRows) {
		return 0, err
	}

	return highestBid, nil
}

func UpdateBalance(ctx context.Context, userID, amount int64) (err error) {
	_, err = db.ExecContext(ctx, queryUpdateBalance, amount, userID)
	if err != nil {
		return err
	}
	return nil
}

func InsertPayment(ctx context.Context, payment Payment) (id int64, err error) {

	err = db.QueryRowContext(ctx, queryInsertPayment,
		payment.UserID,
		payment.Amount,
		payment.Status,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertBidCollection(ctx context.Context, bidCollection BidCollection) (err error) {

	_, err = db.ExecContext(ctx, queryInsertBidCollection,
		bidCollection.UserID,
		bidCollection.AuctionID,
		bidCollection.CurrentBid,
		bidCollection.PaymentID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertAuctionDB(ctx context.Context, auction Auction) (id int64, err error) {

	err = db.QueryRowContext(ctx, queryInsertAuction,
		auction.ProductID,
		auction.WinnerUserID,
		auction.InitialBid,
		auction.Multiplier,
		auction.Status,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateAuctionBidDB(ctx context.Context, auction Auction) (err error) {

	_, err = db.ExecContext(ctx, queryUpdateAuctionBid,
		auction.WinnerUserID,
		auction.HighestBid,
		auction.ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertProductDB(ctx context.Context, product Product) (id int64, err error) {

	err = db.QueryRowContext(ctx, queryInsertProduct,
		product.UserID,
		product.ProductName,
		product.ImageURL,
		product.Status,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetAllAuction(ctx context.Context) (auctions []Auction, err error) {
	rows, err := db.QueryContext(ctx, queryGetAllAuction)
	if err != nil {
		return []Auction{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var auction Auction
		err = rows.Scan(&auction.ID,
			&auction.ProductID,
			&auction.WinnerUserID,
			&auction.InitialBid,
			&auction.HighestBid,
			&auction.Multiplier,
			&auction.Status)
		if err != nil {
			return []Auction{}, err
		}

		auctions = append(auctions, auction)
	}

	return auctions, nil
}
func GetAllProduct(ctx context.Context) (products []Product, err error) {
	rows, err := db.QueryContext(ctx, queryGetAllProduct)
	if err != nil {
		return []Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID,
			&product.UserID,
			&product.ProductName,
			&product.ImageURL,
			&product.Status,
		)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetUser(ctx context.Context, username, password string) (userData User, err error) {

	err = db.QueryRowContext(ctx, queryGetUserLogin, username, password).Scan(
		&userData.ID,
		&userData.UserType,
		&userData.Username,
		&userData.Status,
		&userData.Balance,
	)
	if err != nil {
		return User{}, err
	}

	return userData, nil
}

func GetTimeWindowDB(ctx context.Context, auctionID int64) (timeWindow TimeWindow, err error) {
	err = db.QueryRowContext(ctx, queryGeTimeWindow, auctionID).Scan(
		&timeWindow.ID,
		&timeWindow.AuctionID,
		&timeWindow.StartTime,
		&timeWindow.EndTime,
		&timeWindow.Status,
	)
	if err != nil {
		return TimeWindow{}, err
	}

	return timeWindow, nil
}

func InsertTWDB(ctx context.Context, tw TimeWindow) error {

	err := db.QueryRowContext(ctx, queryInsertTimeWindow,
		tw.AuctionID,
		tw.StartTime,
		tw.EndTime,
		tw.Status,
	)
	if err != nil {
		return err.Err()
	}

	return nil
}

const (
	queryGetUserInfoDB string = `
		SELECT 
		COALESCE(id, 0) as id, 
		COALESCE(user_type, 0) as user_type, 
		COALESCE(username, '') as username,
		COALESCE("status", 0) as "status",
		COALESCE(balance, 0) as balance
		FROM
			"user"
		WHERE
			id=$1
	`

	queryGetAuction string = `
		SELECT 
			COALESCE(id, 0) as id,
			COALESCE(product_id, 0) as product_id,
			COALESCE(winner_user_id, 0) as winner_user_id,
			COALESCE(initial_bid, 0) as initial_bid,
			COALESCE(highest_bid, 0) as highest_bid,
			COALESCE(multiplier, 0) as multiplier,
			COALESCE("status", 0) as "status"
		FROM 
			auction
		WHERE 
			product_id = $1
	`

	queryGetAllAuction string = `
	SELECT 
		COALESCE(id, 0) as id,
		COALESCE(product_id, 0) as product_id,
		COALESCE(winner_user_id, 0) as winner_user_id,
		COALESCE(initial_bid, 0) as initial_bid,
		COALESCE(highest_bid, 0) as highest_bid,
		COALESCE(multiplier, 0) as multiplier,
		COALESCE("status", 0) as "status"
	FROM 
		auction;
`

	queryGetProduct string = `
		SELECT 
			COALESCE(id, 0) as id,
			COALESCE(user_id, 0) as user_id,
			COALESCE(product_name, '') as product_name,
			COALESCE(image_url, '') as image_url,
			COALESCE("status", 0) as "status"
		FROM 
			product
		WHERE 
			id = $1;
	`

	queryGetProductByUserID string = `
	SELECT 
		COALESCE(id, 0) as id,
		COALESCE(user_id, 0) as user_id,
		COALESCE(product_name, '') as product_name,
		COALESCE(image_url, '') as image_url,
		COALESCE("status", 0) as "status"
	FROM 
		product
	WHERE 
		user_id = $1;
`

	queryGetAllProduct string = `
	SELECT 
		COALESCE(id, 0) as id,
		COALESCE(user_id, 0) as user_id,
		COALESCE(product_name, '') as product_name,
		COALESCE(image_url, '') as image_url,
		COALESCE("status", 0) as "status"
	FROM 
		product;
`

	queryGetSumBidCollection string = `
	SELECT 
		COALESCE(sum(current_bid), 0) as current_bid
	FROM 
		bid_collection
	WHERE 
		user_id = $1 and auction_id = $2;
    `

	queryUpdateBalance string = `
	UPDATE 
		"user"
	SET 
		balance = $1
	WHERE 
		id = $2;
    `

	queryInsertPayment string = `
	INSERT INTO 
		payment (
			user_id,
			amount,
			"status",
			create_time
		)
	VALUES ($1, $2, $3, $4)
	RETURNING
		id
	`

	queryInsertBidCollection string = `
	INSERT INTO 
		bid_collection (
			user_id,
			auction_id,
			current_bid,
			payment_id,
			create_time
		)
	VALUES ($1, $2, $3, $4, $5)
	`

	queryInsertProduct string = `
	INSERT INTO 
		product (
			user_id,
			product_name,
			image_url,
			status,
			create_time
		)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id
	`

	queryInsertTimeWindow string = `
	INSERT INTO 
		timewindow (
			auction_id,
			start_time,
			end_time,
			status,
			create_time,
			update_time
		)
	VALUES ($1, $2, $3, $4, 'NOW()', 'NOW()')
	RETURNING
		id
	`

	queryInsertAuction string = `
	INSERT INTO 
		auction (
			product_id,
			winner_user_id,
			initial_bid,
			multiplier,
			"status",
			create_time
		)
	VALUES ($1, $2, $3, $4, $5, $5)
	RETURNING id
	`

	queryUpdateAuctionBid string = `
	UPDATE auction
	SET 
		winner_user_id = $1,
		highest_bid = $2
	WHERE
		id = $3
	`

	queryGetUserLogin string = `
	SELECT 
	COALESCE(id, 0) as id, 
	COALESCE(user_type, 0) as user_type, 
	COALESCE(username, '') as username,
	COALESCE("status", 0) as "status",
	COALESCE(balance, 0) as balance
	FROM
		"user"
	WHERE 
		username = $1 AND	password = $2
`

	queryGeTimeWindow string = `
	SELECT 
	COALESCE(id, 0) as id, 
	COALESCE(auction_id, 0) as auction_id,
	COALESCE(start_time, '0001-01-01T00:00:00Z'::timestamp) as start_time, 
	COALESCE(end_time, '0001-01-01T00:00:00Z'::timestamp) as end_time,
	COALESCE("status", 0) as "status"
	FROM
		"timewindow"
	WHERE 
		auction_id = $1;
`
)
