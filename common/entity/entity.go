package entity

import "github.com/google/uuid"

type ID uuid.UUID

func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

// TODO: uncouple gorm models (see repository pkg) and this package entities
// -> use runtime tags generation
