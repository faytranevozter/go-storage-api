run: 
	go run app/main.go
build: 
	go build -o storage-api app/main.go

.PHONY: build