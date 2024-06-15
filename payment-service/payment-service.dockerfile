# base go image
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o paymentApp ./cmd/api

RUN chmod +x /app/paymentApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/paymentApp /app

CMD ["/app/paymentApp"]