build:
	go build -o dist/trade_exchange cmd/exchange/main.go
	go build -o dist/trade_http_server cmd/exchangehttp/main.go

docker-build:
	docker build -f cmd/exchange/Dockerfile -t jesseobrien/exchange .
	docker build -f cmd/exchangehttp/Dockerfile -t jesseobrien/exchange-http .