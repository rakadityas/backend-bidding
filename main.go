package main

func main() {
	initDB()
	initRedis()
	initProducer()
	initConsumer()
	initHTTP()
}
