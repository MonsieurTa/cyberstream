package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Video struct {
	ID             uuid.UUID      `gorm:"column:id;type:uuid;not null" json:"id"`
	Name           string         `gorm:"column:name;unique;not null" json:"name"`
	Hash           string         `gorm:"column:hash;unique;not null" json:"hash"`
	FilePath       string         `gorm:"column:file_path;unique;not null" json:"file_path,omitempty"`
	SubtitlesPaths pq.StringArray `gorm:"type:text[];column:subtitles_paths;unique" json:"subtitle_path,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`

	Magnet string `gorm:"-" json:"magnet"`
}

func NewVideo(name, hash, filePath, magnet string, subtitlesPaths []string) *Video {
	id := uuid.New()
	return &Video{id, name, hash, filePath, subtitlesPaths, time.Time{}, magnet}
}

func (m *Video) Stored() bool {
	return !m.CreatedAt.IsZero()
}
