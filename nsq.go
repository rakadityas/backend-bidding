package main

import (
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

type messageHandler struct{}

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

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

func initConsumer() {
	config := nsq.NewConfig()
	config.MaxAttempts = 10
	config.MaxInFlight = 5
	config.MaxRequeueDelay = time.Second * 900
	config.DefaultRequeueDelay = time.Second * 0

	//Creating the consumer
	consumer, err := nsq.NewConsumer("Topic_Example", "Channel_Example", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(&messageHandler{})

	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	consumer.Stop()
}

// HandleMessage implements the Handler interface.
func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	//Process the Message
	var request Message
	if err := json.Unmarshal(m.Body, &request); err != nil {
		log.Println("Error when Unmarshaling the message body, Err : ", err)
		// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
		return err
	}
	//Print the Message
	log.Println("Message")
	log.Println("--------------------")
	log.Println("Name : ", request.Name)
	log.Println("Content : ", request.Content)
	log.Println("Timestamp : ", request.Timestamp)
	log.Println("--------------------")
	log.Println("")
	// Will automatically set the message as finish
	return nil
}
