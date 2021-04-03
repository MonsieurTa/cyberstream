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

func (u UserGORM) Create(user *entity.User) (uuid.UUID, error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return user.ID, result.Error
	}
	return user.ID, nil
}
