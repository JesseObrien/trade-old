FROM golang:latest as build

WORKDIR /app
COPY . .
RUN GOOS=linux go build -o -ldflags="-s -w" -o exchange ./cmd/exchange/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build /app/exchange .
CMD ["./exchange"]