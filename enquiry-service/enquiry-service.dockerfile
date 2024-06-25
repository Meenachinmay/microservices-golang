# base go image
FROM golang:1.22-alpine as builder

# Set the timezone
ENV TZ="Asia/Tokyo"

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o enquiryApp ./cmd/api

RUN chmod +x /app/enquiryApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/enquiryApp /app

CMD ["/app/enquiryApp"]