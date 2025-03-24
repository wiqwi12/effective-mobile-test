package models

import "github.com/google/uuid"

type Group struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
