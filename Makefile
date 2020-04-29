build:
	go build -o dist/exchange cmd/exchange/main.go
	go build -o dist/broker cmd/broker/main.go

docker-build:
	docker build -f cmd/exchange/Dockerfile -t jesseobrien/exchange .
	docker build -f cmd/exchangehttp/Dockerfile -t jesseobrien/exchange-http .