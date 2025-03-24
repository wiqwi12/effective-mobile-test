package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
)

type SongRepository interface {
	CreateSong(ctx context.Context, song models.Song) error
	UpdateSong(ctx context.Context, song models.Song) error
	DeleteSong(ctx context.Context, songId uuid.UUID) error
	GetSongById(ctx context.Context, songId uuid.UUID) (models.Song, error)
	GetSongExsistsById(ctx context.Context, songId uuid.UUID) (bool, error) // можно было сделать проверку через sqlNoRows но я чет подзабил
	GetSongTextById(ctx context.Context, songId uuid.UUID) (string, error)
	GetSong(ctx context.Context, song models.Song) (models.Song, error)
	SongExistsByDetails(ctx context.Context, song models.Song) (bool, error)
	GetSongsWithFilter(ctx context.Context, request dto.FilteredRequest) ([]models.Song, error)
}
