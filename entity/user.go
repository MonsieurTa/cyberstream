package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	PublicInfo  PublicInfo
	Credentials Credentials
	CreatedAt   time.Time
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

func (u *User) FillWith(c CreateUserT) error {
	u.ID = uuid.New()

	err := u.Credentials.FillWith(u.ID, c.Username, c.Password)
	if err != nil {
		return err
	}
	u.PublicInfo.FillWith(u.ID, c.FirstName, c.LastName, c.Phone, c.Email)
	return nil
}
