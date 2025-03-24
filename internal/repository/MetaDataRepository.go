package repository

import (
	"context"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
)

type MetaDataRepository interface {
	GetSongDetails(ctx context.Context, request dto.CreateSongRequest) (models.Song, error)
}
