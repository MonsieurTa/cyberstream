package entity

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;not null" json:"id"`
	Name      string    `gorm:"column:name;unique;not null" json:"name"`
	Path      string    `gorm:"column:path;unique;not null" json:"path,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`

	Magnet string `gorm:"-" json:"magnet"`
}

func NewMovie(name, path, magnet string) *Movie {
	id := uuid.New()
	return &Movie{id, name, path, time.Time{}, magnet}
}

func (m *Movie) Stored() bool {
	return !m.CreatedAt.IsZero()
}
