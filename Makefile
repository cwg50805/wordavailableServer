install: 
ifeq (, $(wildcard $(shell which swag)))
	@go get -u github.com/swaggo/swag/cmd/swag
endif
	swag init
	@go build -v .
