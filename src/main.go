package main

import (
	db "db"
	"fmt"
	"models"
)

func main() {
	db := db.NewDatabase()
	db.AutoMigrate()

	user, err := db.User.Create(&models.User{
		FirstName: "William",
		LastName:  "TA",
		Contact: models.Contact{
			Address: "10 impasse Jules Daunay",
			Phone:   "nil",
			Email:   "nil",
		},
		Credential: models.Credential{
			UserName: "wta",
			Password: []byte{},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", user)
}
