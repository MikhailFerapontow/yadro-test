all:
	go build -o app main.go

test:
	go test ./...
