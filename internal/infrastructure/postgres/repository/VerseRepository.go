package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
)

type VerseRepository struct {
	Db     *sql.DB
	Logger *logger.Logger
}

func NewVerseRepository(db *sql.DB, logger *logger.Logger) *VerseRepository {
	return &VerseRepository{
		Db:     db,
		Logger: logger,
	}
}

func (r *VerseRepository) AddVerses(ctx context.Context, req dto.AddVersesRequest) error {

	if req.Song.Id == uuid.Nil {
		r.Logger.Info.Error("Cannot add verses: song ID is nil")
		return errors.New("song ID cannot be nil")
	}

	if len(req.Verses) == 0 {
		r.Logger.Info.Error("Cannot add verses: no verses provided",
			"song_id", req.Song.Id)
		return errors.New("no verses provided")
	}

	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		r.Logger.Info.Error("Failed to begin transaction for adding verses",
			"error", err,
			"song_id", req.Song.Id)
		return err
	}
	defer tx.Rollback()

	deleteQuery, deleteArgs, err := squirrel.Delete("verses").
		Where(squirrel.Eq{"song_id": req.Song.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build delete query",
			"error", err,
			"song_id", req.Song.Id)
		return err
	}

	_, err = tx.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		r.Logger.Info.Error("Failed to delete existing verses",
			"error", err,
			"song_id", req.Song.Id)
		return err
	}

	for i, verse := range req.Verses {
		verseId := uuid.New()

		insertQuery, insertArgs, err := squirrel.Insert("verses").
			Columns("id", "song_id", "verse_number", "text").
			Values(verseId, req.Song.Id, i+1, verse).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			r.Logger.Info.Error("Failed to build insert query",
				"error", err,
				"song_id", req.Song.Id,
				"verse_number", i+1)
			return err
		}

		_, err = tx.ExecContext(ctx, insertQuery, insertArgs...)
		if err != nil {
			r.Logger.Info.Error("Failed to insert verse",
				"error", err,
				"song_id", req.Song.Id,
				"verse_number", i+1)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		r.Logger.Info.Error("Failed to commit transaction for adding verses",
			"error", err,
			"song_id", req.Song.Id)
		return err
	}

	return nil
}

func (r *VerseRepository) GetPaginatedVerses(ctx context.Context, request dto.PaginatedVersesRequest) (dto.PaginatedVersesResponse, error) {

	var resp dto.PaginatedVersesResponse
	firstVerse := (request.Page-1)*request.Limit + 1
	lastVerse := (request.Page) * request.Limit

	query, args, err := squirrel.Select("text").
		From("verses").
		Where(squirrel.Eq{
			"song_id": request.SongId,
		}).
		Where(squirrel.GtOrEq{
			"verse_number": firstVerse,
		}).Where(squirrel.LtOrEq{
		"verse_number": lastVerse,
	}).
		OrderBy("verse_number ASC").
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build query query", err)
		return dto.PaginatedVersesResponse{}, err
	}

	rows, err := r.Db.QueryContext(ctx, query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to query verses",
			"error", err,
			"song_id", request.SongId)
		return dto.PaginatedVersesResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var verse string
		err := rows.Scan(&verse)
		if err != nil {
			r.Logger.Info.Error("Failed to scan verse row",
				"error", err)
			return dto.PaginatedVersesResponse{}, err
		}

		resp.Verses = append(resp.Verses, verse)
	}
	err = rows.Err()
	if err != nil {
		r.Logger.Info.Error("Error iterating over rows",
			"error", err)
		return dto.PaginatedVersesResponse{}, err
	}

	if len(resp.Verses) == 0 && request.Page == 1 {
		return dto.PaginatedVersesResponse{
			Page:  request.Page,
			Limit: request.Limit,
			Total: 0,
		}, errors.New("no verses found")
	}

	resp.Limit = request.Limit
	resp.Page = request.Page
	resp.Total = len(resp.Verses)
	return resp, err
}
