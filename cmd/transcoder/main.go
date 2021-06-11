package main

import (
	"log"
	"net"
	"os"

	"github.com/MonsieurTa/hypertube/pkg/transcoder/handler"
	"github.com/joho/godotenv"
)

func initEnv() {
	env := os.Getenv("HYPERTUBE_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
}

func main() {
	initEnv()

	ln, err := net.Listen("tcp", "localhost:3010")
	if err != nil {
		panic("could not listen on tcp protocol at localhost:3010")
	}
	log.Println("Listening tcp connections on localhost:3010")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		handler.TCPHandler(conn)
	}
}
