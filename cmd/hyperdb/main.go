package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MonsieurTa/hypertube/common/db"
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/MonsieurTa/hypertube/common/infrastructure/repository"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	firstname string
	lastname  string
	phone     string
	email     string
	username  string
	password  string
)

func initEnv() {
	env := os.Getenv("HYPERTUBE_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
}

func init() {
	flag.Usage = func() {
		usage := "usage: hyperdb\n" +
			"\t-firstname=firstname\n" +
			"\t-lastname=lastname\n" +
			"\t-phone=phone\n" +
			"\t-email=email (mandatory)\n" +
			"\t-username=username (mandatory)\n" +
			"\t-password=password (mandatory)"
		fmt.Fprintln(os.Stderr, usage)
	}

	flag.StringVar(&firstname, "firstname", "", "user's firstname")
	flag.StringVar(&lastname, "lastname", "", "user's lastname")
	flag.StringVar(&phone, "phone", "", "user's phone")
	flag.StringVar(&email, "email", "", "user's email")
	flag.StringVar(&username, "username", "", "user's username")
	flag.StringVar(&password, "password", "", "user's password")
	flag.Parse()
}

func main() {
	initEnv()

	if username == "" || password == "" || email == "" {
		flag.Usage()
		return
	}

	if len(username) < 4 {
		fmt.Fprintln(os.Stderr, "username's length too short. 4 characters minimum")
		return
	}

	if len(password) < 8 {
		fmt.Fprintln(os.Stderr, "password's length too short. 8 characters minimum")
		return
	}

	db := db.InitDB(&db.PSQLConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Db:       os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}, &gorm.Config{})

	gormDB := db.DB().(*gorm.DB)
	repo := repository.NewUserGORM(gormDB)

	userParams := entity.CreateUserT{
		FirstName: firstname,
		LastName:  lastname,
		Phone:     phone,
		Email:     email,
		Username:  username,
		Password:  password,
	}

	user, err := entity.NewUser(userParams)
	if err != nil {
		panic(err)
	}

	ID, err := repo.Create(user)
	if err != nil {
		panic(err)
	}
	log.Println(ID.String())
}
