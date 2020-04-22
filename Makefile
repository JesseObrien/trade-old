build:
	go build -o dist/exchange cmd/exchange/main.go
	go build -o dist/broker cmd/broker/main.go