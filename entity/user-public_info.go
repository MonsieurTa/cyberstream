package entity

import (
	"time"

	"github.com/google/uuid"
)

type PublicInfo struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;not null"`
	UserID     uuid.UUID `gorm:"column:user_id"`
	FirstName  string    `gorm:"column:first_name"`
	LastName   string    `gorm:"column:last_name"`
	Phone      string    `gorm:"column:phone"`
	Email      string    `gorm:"column:email;unique"`
	PictureURL string    `gorm:"column:picture_url"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
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
