# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster

WORKDIR /api

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY go.* ./

RUN go mod download

ADD . .

RUN go build -o hypertube-api cmd/api/main.go

EXPOSE 3000

CMD ["./hypertube-api"]
