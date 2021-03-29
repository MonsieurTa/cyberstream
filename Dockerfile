FROM golang:1.16.2-buster as build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY src .

RUN go build -o hypertube .

FROM debian:buster

WORKDIR /app

COPY --from=build /build/hypertube .

EXPOSE 8080

CMD ["/app/hypertube"]
