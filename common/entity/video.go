package entity

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;not null" json:"id"`
	Name      string    `gorm:"column:name;unique;not null" json:"name"`
	Hash      string    `gorm:"column:hash;unique;not null" json:"hash"`
	Path      string    `gorm:"column:path;unique;not null" json:"path,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`

	Magnet string `gorm:"-" json:"magnet"`
}

func NewVideo(name, hash, path, magnet string) *Video {
	id := uuid.New()
	return &Video{id, name, hash, path, time.Time{}, magnet}
}

func (m *Video) Stored() bool {
	return !m.CreatedAt.IsZero()
}
