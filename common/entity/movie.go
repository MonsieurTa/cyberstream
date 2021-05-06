package entity

import "github.com/google/uuid"

type Movie struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;not null"`
	Name      string    `gorm:"column:name;unique;not null"`
	Path      string    `gorm:"column:path;unique;not null"`
	CreatedAt string    `gorm:"column:created_at"`
}

func NewMovie(name, path string) *Movie {
	id := uuid.New()
	return &Movie{
		ID:   id,
		Name: name,
		Path: path,
	}
}
