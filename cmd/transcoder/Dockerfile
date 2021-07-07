# syntax = docker/dockerfile:1.2

FROM golang:1.16.2-buster as builder

WORKDIR /media

COPY go.* ./

RUN go mod download

ADD . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    make transcoder

FROM ubuntu:20.04
COPY --from=builder /media/bin/transcoder /bin

RUN apt update && apt install -y wget xz-utils ca-certificates
RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN tar -xvf ffmpeg-release-amd64-static.tar.xz \
    && cp ffmpeg-4.4-amd64-static/ffmpeg /bin \
    && cp ffmpeg-4.4-amd64-static/ffprobe /bin \
    && rm -rf ffmpeg-release-amd64-static.tar.xz ffmpeg-4.4-amd64-static

EXPOSE 3010

CMD ["/bin/transcoder"]
