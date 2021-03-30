package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName  string
	LastName   string
	Contact    Contact
	Credential Credential
}

type Contact struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	Address   string
	Phone     string
	Email     string
	UpdatedAt time.Time
}

type Credential struct {
	ID       uint `gorm:"primarykey"`
	UserID   uint
	UserName string
	Password []byte
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
