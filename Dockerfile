FROM golang:1.16.2-buster

WORKDIR /app

ADD . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 8080

ENTRYPOINT CompileDaemon --build="go build -o hypertube api/main.go api/app.go" --command=./hypertube
