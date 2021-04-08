package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;not null"`
	PublicInfo PublicInfo
	Credential Credential
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type CreateUserT struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Username  string
	Password  string
}

func NewUser(c CreateUserT) (*User, error) {
	userID := uuid.New()
	credential, err := NewCredential(userID, c.Username, c.Password)
	if err != nil {
		return &User{ID: userID}, err
	}
	publicInfo := NewPublicInfo(userID, c.FirstName, c.LastName, c.Phone, c.Email)
	return &User{
		ID:         userID,
		Credential: *credential,
		PublicInfo: *publicInfo,
	}, nil
}
