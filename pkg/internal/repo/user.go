package repositories

import (
	models "models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db}
}

func (r UserRepository) Create(data *models.User) (*models.User, error) {
	result := r.db.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
