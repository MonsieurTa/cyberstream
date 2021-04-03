FROM golang:1.16.2-buster as build

WORKDIR /build

ADD . .

RUN go mod download

RUN go build -o hypertube api/*.go

FROM debian:buster

WORKDIR /app

COPY --from=build /build/hypertube .

EXPOSE 8080

CMD ["/app/hypertube"]
