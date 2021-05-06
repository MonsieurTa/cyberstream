package repository

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieGORM struct {
	db *gorm.DB
}

func NewMovieGORM(db *gorm.DB) *MovieGORM {
	return &MovieGORM{db}
}

func (m MovieGORM) FindByID(movieID uuid.UUID) (*entity.Movie, error) {
	movie := entity.Movie{}

	err := m.db.First(&movie, "id = ?", movieID).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (m MovieGORM) SearchByName(pattern string) (*entity.Movie, error) {
	movie := entity.Movie{}

	param := `%` + pattern + `%`
	err := m.db.First(&movie, "name LIKE = ?", param).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
