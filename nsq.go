package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	producer *nsq.Producer
)

type mhUpdateScoreBoard struct{}
type mhInsertCollectionAndPayment struct{}

func DoPublishNSQ(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	//Publish the Message
	err = producer.Publish(topic, payload)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func initProducer() {
	var err error

	//The only valid way to instantiate the Config
	config := nsq.NewConfig()

	//Creating the Producer using NSQD Address
	producer, err = nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfigGeneral() *nsq.Config {
	config := nsq.NewConfig()
	config.MaxAttempts = 10
	config.MaxInFlight = 5
	config.MaxRequeueDelay = time.Second * 900
	config.DefaultRequeueDelay = time.Second * 0
	return config
}

func GetConfigTopicUpdateScoreBoard() *nsq.Config {
	config := nsq.NewConfig()
	config.MaxAttempts = 10
	config.MaxInFlight = 1
	config.MaxRequeueDelay = time.Second * 900
	config.DefaultRequeueDelay = time.Second * 0
	return config
}

func initConsumer() {

	//Creating the consumer
	cUpdateScoreBoard, err := nsq.NewConsumer("Update_Scoreboard", "update_scoreboard", GetConfigTopicUpdateScoreBoard())
	if err != nil {
		log.Fatal(err)
	}
	cUpdateScoreBoard.AddHandler(&mhUpdateScoreBoard{})

	cInsertCollectionAndPayment, err := nsq.NewConsumer("Insert_Collection_And_Payment", "insert_collection_and_payment", GetConfigGeneral())
	if err != nil {
		log.Fatal(err)
	}
	cInsertCollectionAndPayment.AddHandler(&mhInsertCollectionAndPayment{})

	err = cUpdateScoreBoard.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	err = cInsertCollectionAndPayment.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cUpdateScoreBoard.Stop()
	cInsertCollectionAndPayment.Stop()
}

// HandleMessage implements the Handler interface.
func (h *mhUpdateScoreBoard) HandleMessage(m *nsq.Message) error {

	var request UpdateScoreboardNSQ

	ctx := context.Background()

	if err := json.Unmarshal(m.Body, &request); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		return err
	}

	CheckHighestBid(ctx, request.BidAmount, request.UserID)

	m.Finish()
	return nil
}

func (h *mhInsertCollectionAndPayment) HandleMessage(m *nsq.Message) error {

	var (
		err     error
		request InsertPaymentAndBidCollectionNSQ
	)

	ctx := context.Background()

	if err := json.Unmarshal(m.Body, &request); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		return err
	}

	err = InsertBidCollectionAndPayment(ctx, request.BidCollection, request.Payment)
	if err != nil {
		log.Println("error: ", err)
		return err
	}

	m.Finish()
	return nil
}
