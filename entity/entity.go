package entity

import "github.com/google/uuid"

type ID uuid.UUID

func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}
