FROM golang:latest as build

WORKDIR /app
COPY . .
RUN GOOS=linux go build -o -ldflags="-s -w" -o http ./cmd/exchangehttp/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build /app/http .
CMD ["./exchange"]