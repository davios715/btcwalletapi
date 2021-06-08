FROM golang:alpine

WORKDIR /btcwalletapi-test

ADD . .

RUN go mod download

ENTRYPOINT go build  && ./btcwalletapi