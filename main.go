package main

func main() {
	initFirebase()
	initDB()
	initRedis()
	initProducer()
	// initConsumer()
	initHTTP()
}
