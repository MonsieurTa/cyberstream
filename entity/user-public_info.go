package entity

import (
	"time"

	"github.com/google/uuid"
)

type PublicInfo struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	FirstName  string
	LastName   string
	Phone      string
	Email      string
	PictureURL string
	UpdatedAt  time.Time
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

func (p *PublicInfo) FillWith(userID uuid.UUID, firstname, lastname, phone, email string) {
	p.ID = uuid.New()
	p.UserID = userID
	p.FirstName = firstname
	p.LastName = lastname
	p.Phone = phone
	p.Email = email
}

func (p PublicInfo) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *PublicInfo) Update(email, pictureURL string) {
	if len(email) > 0 {
		p.Email = email
	}
	if len(pictureURL) > 0 {
		p.PictureURL = pictureURL
	}
}
