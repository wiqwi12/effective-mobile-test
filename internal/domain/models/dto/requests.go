package dto

import (
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
)

type CreateSongRequest struct {
	Group string `json:"group" validate:"required,min=1,max=100"`
	Title string `json:"title" validate:"required,min=1,max=100"`
}

type GetSongByIdRequest struct {
	Id uuid.UUID `json:"id" validate:"required,uuid"`
}

type UpdateSongRequest struct {
	GroupName   string `json:"group,omitempty"`
	Title       string `json:"title,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type DeleteSongByIdRequest struct {
	Id uuid.UUID `json:"id" validate:"required,uuid"`
}

type FilteredRequest struct {
	Title       string `json:"title,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type AddVersesRequest struct {
	Verses []string    `json:"verses"`
	Song   models.Song `json:"song"`
}

type PaginatedVersesRequest struct {
	SongId uuid.UUID `json:"-" validate:"required,uuid"`
	Page   int       `json:"page" validate:"required,min=1"`
	Limit  int       `json:"limit" validate:"required,min=1"` //куплетов на страницу
}
