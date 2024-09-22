all: algo

algo:
	@go run cmd/algo/main.go

grounds:
	@go run cmd/testing_grounds/*

football:
	@go run cmd/football/main.go
