package dto

import (
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
)

type StandartResponse struct {
	Song    models.Song `json:"song,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type SongsResponse struct {
	Songs   []models.Song `json:"songs,omitempty"`
	Message string        `json:"message,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type SongTextResponse struct {
	Text string `json:"text"`
}

type SongDetailResponse struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongTextPaginatedResponse struct {
	Verses []string `json:"verses"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
	Total  int      `json:"total"`
}

type PaginatedVersesResponse struct {
	Verses []string `json:"verses"`
	Page   int      `json:"page"`
	Limit  int      `json:"limit"`
	Total  int      `json:"total"`
}
