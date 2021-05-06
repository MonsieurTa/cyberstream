# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster

WORKDIR /api

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY go.* ./

RUN go mod download

ADD . .

RUN mkdir /static

RUN go build -o hypertube-media cmd/media/main.go

EXPOSE 8080

CMD ["./hypertube-media"]
