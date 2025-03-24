package repository

import (
	"context"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
)

type VersesRepository interface {
	AddVerses(ctx context.Context, req dto.AddVersesRequest) error
	GetPaginatedVerses(ctx context.Context, request dto.PaginatedVersesRequest) (dto.PaginatedVersesResponse, error)
}
