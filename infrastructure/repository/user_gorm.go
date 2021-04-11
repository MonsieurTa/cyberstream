package repository

import (
	"github.com/MonsieurTa/hypertube/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserGORM struct {
	db *gorm.DB
}

func NewUserGORM(db *gorm.DB) *UserGORM {
	return &UserGORM{db}
}

func (u UserGORM) Create(user *entity.User) (*uuid.UUID, error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user.ID, nil
}

func (u UserGORM) CredentialsExist(username, password string) (*uuid.UUID, error) {
	var credentials entity.Credentials

	err := u.db.Where("username = ?", username).First(&credentials).Error
	if err != nil {
		return nil, err
	}

	err = credentials.CheckPassword(password)
	if err != nil {
		return nil, err
	}
	return &credentials.UserID, nil
}

func (u UserGORM) UpdateCredentials(userID *uuid.UUID, username, password string) error {
	credentials := entity.Credentials{}

	err := u.db.First(&credentials, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	err = credentials.Update(username, password)
	if err != nil {
		return err
	}
	return u.db.Model(&credentials).Updates(credentials).Error
}

func (u UserGORM) UpdatePublicInfo(userID *uuid.UUID, email, pictureURL string) error {
	publicInfo := entity.PublicInfo{}

	err := u.db.First(&publicInfo, "user_id = ?", userID).Error
	if err != nil {
		return err
	}
	publicInfo.Update(email, pictureURL)
	return u.db.Model(&publicInfo).Updates(publicInfo).Error
}
