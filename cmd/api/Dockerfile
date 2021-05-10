# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster

WORKDIR /api

COPY go.* ./

RUN go mod download

ADD . .

RUN make api

EXPOSE 3000

CMD ["./api"]
