# base go image
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app
COPY templates /templates

WORKDIR /app

RUN go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailApp /app
COPY templates /templates

CMD ["/app/mailApp"]