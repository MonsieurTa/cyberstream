package entity

import "github.com/google/uuid"

const (
	SUBSPLEASE ProviderName = "SUBSPLEASE"
)

type ProviderName string

type Provider struct {
	ID     uuid.UUID    `gorm:"column:id;type:uuid;not null"`
	Name   ProviderName `gorm:"column:provider;unique;not null"`
	Videos []Video      `gorm:"many2many:broadcaster_videos"`
}
