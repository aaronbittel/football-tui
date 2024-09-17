all: algo

algo:
	@go run cmd/algo/main.go

grounds:
	@go run cmd/testing_grounds/main.go
