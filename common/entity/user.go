package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;not null"`
	PublicInfo  PublicInfo
	Credentials Credentials
	CreatedAt   time.Time `gorm:"column:created_at"`
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
	credentials, err := NewCredentials(userID, c.Username, c.Password)
	if err != nil {
		return &User{ID: userID}, err
	}
	publicInfo := NewPublicInfo(userID, c.FirstName, c.LastName, c.Phone, c.Email)
	return &User{
		ID:          userID,
		Credentials: *credentials,
		PublicInfo:  *publicInfo,
	}, nil
}

func (u *User) Persisted() bool {
	return !u.CreatedAt.IsZero()
}
