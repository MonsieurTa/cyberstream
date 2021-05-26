package repository

import (
	"errors"

	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoGORM struct {
	db *gorm.DB
}

func NewVideoGORM(db *gorm.DB) *VideoGORM {
	return &VideoGORM{db}
}

func (m VideoGORM) FindByID(videoID uuid.UUID) (*entity.Video, error) {
	video := entity.Video{}

	err := m.db.First(&video, "id = ?", videoID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &video, nil
}

func (m VideoGORM) FindByName(name string) (*entity.Video, error) {
	video := entity.Video{}

	err := m.db.First(&video, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &video, nil
}

func (m VideoGORM) SearchByName(pattern string) (*entity.Video, error) {
	video := entity.Video{}

	param := `%` + pattern + `%`
	err := m.db.First(&video, "name LIKE = ?", param).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (m VideoGORM) Create(video *entity.Video) (uuid.UUID, error) {
	err := m.db.Create(video).Error
	if err != nil {
		return uuid.Nil, err
	}
	return video.ID, nil
}
