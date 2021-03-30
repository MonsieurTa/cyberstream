package db

import (
	"models"
	repo "repo"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db   *gorm.DB
	User repo.UserRepository
}

func NewDatabase() Database {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic("failed to connect to database!")
	}
	return Database{
		db:   db,
		User: repo.NewUserRepository(db),
	}
}

func (d Database) AutoMigrate() {
	d.db.AutoMigrate(&models.User{})
	d.db.AutoMigrate(&models.Contact{})
	d.db.AutoMigrate(&models.Credential{})
}
