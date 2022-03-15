package main

import (
	"context"
	"database/sql"
	"fmt"

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
		&auction.Multiplier,
		&auction.Status,
	)
	if err != nil {
		return auction, err
	}

	return auction, nil
}

func GetProductDB(ctx context.Context, productID int64) (product Product, err error) {

	err = db.QueryRowContext(ctx, queryGetAuction, productID).Scan(
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

const (
	queryGetUserInfoDB string = `
		SELECT 
		COALESCE(id, 0) as id, 
		COALESCE(user_type, 0) as user_type, 
		COALESCE(username, '') as username,
		COALESCE("status", 0) as "status"
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
			COALESCE(multiplier, 0) as multiplier,
			COALESCE("status", 0) as "status"
		FROM 
			auction
		WHERE 
			product_id = $1
	`

	queryGetProduct string = `
		SELECT 
			COALESCE(id, 0) as id,
			COALESCE(user_id, 0) as user_id,
			COALESCE(product_name, '') as product_name,
			COALESCE(image_url, '') as image_url,
			COALESCE("status", 0) as "status"
		FROM 
			auction
		WHERE 
			user_id = $1
	`
)
