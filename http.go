package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func HandleLogin(rw http.ResponseWriter, req *http.Request) {
	var request LoginRequest

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := jsoniter.Unmarshal(body, &request); err != nil {
		fmt.Println(err)
		return
	}

	user, err := Login(req.Context(), request.Username, request.Password)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jsoniter.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(resp))
	return
}

func HandleGetUserInfo(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	userID, err := strconv.ParseInt(req.FormValue("user_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	userData, err := GetUserInfo(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jsoniter.Marshal(userData)
	if err != nil {
		fmt.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(resp))
	return
}

func HandleGetAuctionDetail(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	productID, err := strconv.ParseInt(req.FormValue("product_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	userID, err := strconv.ParseInt(req.FormValue("user_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := &GetAuctionDetailRequest{
		ProductID: productID,
		UserID:    userID,
	}

	response, err := GetAuctionDetail(ctx, *request)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jsoniter.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(resp))
	return

}

func HandleAuctionBidding(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var payload AuctionBidRequest

	err = jsoniter.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := AuctionBidding(ctx, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jsoniter.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(resp))
}

func initHTTP() {
	// get
	http.HandleFunc("/get/auction/detail", HandleGetAuctionDetail)
	http.HandleFunc("/get/user/info", HandleGetUserInfo)
	http.HandleFunc("/get/auction/list", HandleLogin)

	// post
	http.HandleFunc("/auction/create", HandleLogin)
	http.HandleFunc("/auction/bid", HandleAuctionBidding)
	http.HandleFunc("/login", HandleLogin)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
