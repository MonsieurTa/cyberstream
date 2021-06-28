# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster as builder

WORKDIR /api

COPY go.* ./

RUN go mod download

ADD . .

RUN make api

FROM ubuntu:20.04
COPY --from=builder /api/bin/api /bin

EXPOSE 3000

CMD ["./bin/api"]
