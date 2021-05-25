package repository

import (
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProviderGORM struct {
	db *gorm.DB
}

func NewProviderGORM(db *gorm.DB) *ProviderGORM {
	return &ProviderGORM{db}
}

func (m *ProviderGORM) RegisterProviders(providers []entity.Provider) error {
	return m.db.Clauses(clause.OnConflict{DoNothing: true}).Create(providers).Error
}

func (m *ProviderGORM) FindByName(name entity.ProviderName) (*entity.Provider, error) {
	broadcaster := entity.Provider{}

	err := m.db.First(&broadcaster, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &broadcaster, nil
}

func (m *ProviderGORM) StoreVideo(provider *entity.Provider, video *entity.Video) (uuid.UUID, error) {
	// https://gorm.io/docs/associations.html#Association-Mode
	err := m.db.Model(provider).Association("Videos").Append(video)
	if err != nil {
		return uuid.Nil, err
	}
	return video.ID, nil
}
