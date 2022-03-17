package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/cors"
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

func HandleGetAuctionList(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	userID, err := strconv.ParseInt(req.FormValue("user_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := GetAuctionList(ctx, userID, false)
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

func HandleAuctionCreate(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var payload CreateAuctionRequest

	err = jsoniter.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	resultStatus, err := InsertAuction(ctx, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := jsoniter.Marshal(resultStatus)
	if err != nil {
		fmt.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(resp))
}

func initHTTP() {
	mux := http.NewServeMux()
	// get
	mux.HandleFunc("/get/auction/detail", HandleGetAuctionDetail)
	mux.HandleFunc("/get/user/info", HandleGetUserInfo)
	mux.HandleFunc("/get/auction/list", HandleGetAuctionList)

	// post
	mux.HandleFunc("/auction/create", HandleAuctionCreate)
	mux.HandleFunc("/auction/bid", HandleAuctionBidding)
	mux.HandleFunc("/login", HandleLogin)

	handler := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
