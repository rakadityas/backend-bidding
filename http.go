package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func HandleLogin(rw http.ResponseWriter, req *http.Request) {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
}

func initHTTP() {
	// get
	http.HandleFunc("/get/auction/detail", HandleLogin)
	http.HandleFunc("/get/user/info", HandleLogin)
	http.HandleFunc("/get/auction/list", HandleLogin)

	// post
	http.HandleFunc("/auction/create", HandleLogin)
	http.HandleFunc("/auction/bid", HandleLogin)
	http.HandleFunc("/login", HandleLogin)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
