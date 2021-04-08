package repository

import (
	"github.com/MonsieurTa/hypertube/entity"
	"gorm.io/gorm"
)

type CredentialGORM struct {
	db *gorm.DB
}

func NewCredentialGORM(db *gorm.DB) *CredentialGORM {
	return &CredentialGORM{db}
}

func (u CredentialGORM) Validate(username, password string) error {
	var credential entity.Credential

	err := u.db.Where("username = ?", username).First(&credential).Error
	if err != nil {
		return err
	}

	err = credential.CheckPassword(password)
	if err != nil {
		return err
	}
	return nil
}
