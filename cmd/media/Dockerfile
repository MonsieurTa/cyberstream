# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster

WORKDIR /media

RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ffmpeg

COPY go.* ./

RUN go mod download

ADD . .

RUN mkdir /static
RUN mkdir /download

RUN make media

EXPOSE 8080

CMD ["./media"]
