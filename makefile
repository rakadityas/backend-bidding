BUILDHASH = $(shell git rev-parse --verify HEAD | cut -c 1-7)
VERSION = 1.0.0

# go build command
gorun:
	@go build && ./backend-bidding

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down