package models

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	Id          uuid.UUID `json:"song_id"`
	GroupId     uuid.UUID `json:"group_id"`
	GroupName   string    `json:"group_name"`
	Title       string    `json:"title"`
	ReleaseDate string    `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
