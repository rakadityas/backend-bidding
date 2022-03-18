package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var clientFirestore *firestore.Client

func initFirebase() {
	// Use the application default credentials
	ctx := context.Background()
	opt := option.WithCredentialsFile("hackhaton-bgp-2022-firebase-adminsdk-ynchn-e22ec9f24b.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("app", app)

	clientFirestore, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("clientFirestore", clientFirestore)

}
