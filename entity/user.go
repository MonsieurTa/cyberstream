package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;not null"`
	PublicInfo PublicInfo
	Credential Credential
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type PublicInfo struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;not null"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Phone     string    `gorm:"column:phone"`
	Email     string    `gorm:"column:email;unique"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Credential struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;not null"`
	UserID       uuid.UUID `gorm:"column:user_id"`
	Username     string    `gorm:"column:username;unique"`
	PasswordHash string    `gorm:"column:password_hash"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
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

func NewPublicInfo(userID uuid.UUID, firstname, lastname, phone, email string) *PublicInfo {
	return &PublicInfo{
		ID:        uuid.New(),
		UserID:    userID,
		FirstName: firstname,
		LastName:  lastname,
		Phone:     phone,
		Email:     email,
	}
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

func (p PublicInfo) FullName() string {
	return p.FirstName + " " + p.LastName
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
