package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;not null"`
	UserID       uuid.UUID `gorm:"column:user_id"`
	Username     string    `gorm:"column:username;unique"`
	PasswordHash string    `gorm:"column:password_hash"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func NewCredential(userID uuid.UUID, Username, password string) (*Credential, error) {
	v := Credential{
		ID:       uuid.New(),
		UserID:   userID,
		Username: Username,
	}
	err := v.SetPassword(password)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Credential) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	c.PasswordHash = string(passwordHash)
	return nil
}

func (c *Credential) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(c.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
