# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster as builder

WORKDIR /api

COPY go.* ./

RUN go mod download

ADD . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    make api

FROM ubuntu:20.04
COPY --from=builder /api/bin/api /bin

RUN apt update -y && apt install -y ca-certificates

EXPOSE 3000

CMD ["./bin/api"]
