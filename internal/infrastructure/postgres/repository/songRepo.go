package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
	"time"
)

type SongRepository struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewSongRepo(db *sql.DB, logger *logger.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		Logger: logger,
	}
}

func (r *SongRepository) CreateSong(ctx context.Context, song models.Song) error {
	exists, err := r.SongExistsByDetails(ctx, song)
	if err != nil {
		r.Logger.Info.Error("Failed to check if song exists",
			"error", err,
			"song_title", song.Title,
			"group_name", song.GroupName)
		return err
	}

	if exists {
		return errors.New("song already exists")
	}

	query, args, err := squirrel.Insert("Songs").Columns("id, group_id, group_name, title, release_date, text, link").
		Values(song.Id, song.GroupId, song.GroupName, song.Title, song.ReleaseDate, song.Text, song.Link).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for song creation",
			"error", err,
			"song_id", song.Id)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to execute song creation query",
			"error", err,
			"song_id", song.Id,
			"song_title", song.Title)
		return err
	}

	return nil
}

func (r *SongRepository) GetSongById(ctx context.Context, songId uuid.UUID) (models.Song, error) {
	query, args, err := squirrel.Select("id, group_id, group_name, title, release_date, text, link, created_at, updated_at").
		From("songs").
		Where(squirrel.Eq{
			"id": songId,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for song lookup by ID",
			"error", err,
			"song_id", songId)
		return models.Song{}, err
	}

	var song models.Song
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&song.Id,
		&song.GroupId,
		&song.GroupName,
		&song.Title,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			r.Logger.Info.Error("Error executing song lookup query",
				"error", err,
				"song_id", songId)
		}
		return models.Song{}, err
	}

	return song, nil
}

func (r *SongRepository) SongExsistsById(ctx context.Context, songId uuid.UUID) (bool, error) {
	query, args, err := squirrel.Select("1").
		From("songs").
		Where(squirrel.Eq{
			"id": songId,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query to check song existence by ID",
			"error", err,
			"song_id", songId)
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		r.Logger.Info.Error("Error checking if song exists by ID",
			"error", err,
			"song_id", songId)
		return false, err
	}

	return true, nil
}

func (r *SongRepository) SongExistsByDetails(ctx context.Context, song models.Song) (bool, error) {
	query, args, err := squirrel.Select("1").
		From("songs").
		Where(squirrel.Eq{
			"group_name":   song.GroupName,
			"title":        song.Title,
			"release_date": song.ReleaseDate,
			"text":         song.Text,
			"link":         song.Link,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query to check song existence by details",
			"error", err,
			"title", song.Title,
			"group_name", song.GroupName)
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		r.Logger.Info.Error("Error checking if song exists by details",
			"error", err,
			"title", song.Title,
			"group_name", song.GroupName)
		return false, err
	}

	return true, nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, song models.Song) error {
	exists, err := r.SongExsistsById(ctx, song.Id)
	if err != nil {
		r.Logger.Info.Error("Failed to check if song exists before update",
			"error", err,
			"song_id", song.Id)
		return err
	}

	if !exists {
		r.Logger.Info.Error("Song does not exist for update",
			"song_id", song.Id)
		return errors.New("song doesn't exists")
	}

	queryBuilder := squirrel.Update("songs").
		Where(squirrel.Eq{"id": song.Id}).
		Set("updated_at", time.Now()).
		PlaceholderFormat(squirrel.Dollar)

	if song.GroupId != uuid.Nil {
		queryBuilder = queryBuilder.Set("group_id", song.GroupId)
	}
	if song.Title != "" {
		queryBuilder = queryBuilder.Set("title", song.Title)
	}
	if song.ReleaseDate != "" {
		queryBuilder = queryBuilder.Set("release_date", song.ReleaseDate)
	}
	if song.Text != "" {
		queryBuilder = queryBuilder.Set("text", song.Text)
	}
	if song.Link != "" {
		queryBuilder = queryBuilder.Set("link", song.Link)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for song update",
			"error", err,
			"song_id", song.Id)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to execute song update query",
			"error", err,
			"song_id", song.Id)
		return err
	}

	return err
}

func (r *SongRepository) DeleteSong(ctx context.Context, songId uuid.UUID) error {
	exists, err := r.SongExsistsById(ctx, songId)
	if err != nil {
		r.Logger.Info.Error("Failed to check if song exists before deletion",
			"error", err,
			"song_id", songId)
		return err
	}

	if !exists {
		r.Logger.Info.Error("Song does not exist for deletion",
			"song_id", songId)
		return errors.New("song doesn't exist")
	}

	query, args, err := squirrel.Delete("songs").
		Where(squirrel.Eq{
			"id": songId,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for song deletion",
			"error", err,
			"song_id", songId)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to execute song deletion query",
			"error", err,
			"song_id", songId)
		return err
	}

	return nil
}

func (r *SongRepository) GetSongTextById(ctx context.Context, songId uuid.UUID) (string, error) {
	exists, err := r.SongExsistsById(ctx, songId)
	if err != nil {
		r.Logger.Info.Error("Failed to check if song exists before getting text",
			"error", err,
			"song_id", songId)
		return "", err
	}
	if !exists {
		r.Logger.Info.Error("Song does not exist for text retrieval",
			"song_id", songId)
		return "", errors.New("song doesn't exist")
	}

	query, args, err := squirrel.Select("text").
		From("songs").
		Where(squirrel.Eq{
			"id": songId,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for song text retrieval",
			"error", err,
			"song_id", songId)
		return "", err
	}

	var text string
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&text)
	if err != nil {
		r.Logger.Info.Error("Failed to retrieve song text",
			"error", err,
			"song_id", songId)
		return "", err
	}

	return text, nil
}

func (r *SongRepository) GetSongsWithFilter(ctx context.Context, request dto.FilteredRequest) ([]models.Song, error) {
	builder := squirrel.Select("id, group_id, group_name, title, release_date, text, link, created_at, updated_at").
		From("songs")

	// Apply filters
	if request.Title != "" {
		builder = builder.Where(squirrel.Eq{"title": request.Title})
	}
	if request.ReleaseDate != "" {
		builder = builder.Where(squirrel.Eq{"release_date": request.ReleaseDate})
	}
	if request.Text != "" {
		builder = builder.Where(squirrel.Eq{"text": request.Text})
	}
	if request.Link != "" {
		builder = builder.Where(squirrel.Eq{"link": request.Link})
	}
	if request.GroupName != "" {
		builder = builder.Where(squirrel.Eq{"group_name": request.GroupName})
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for filtered songs",
			"error", err)
		return nil, err
	}

	var songs []models.Song
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to execute filtered songs query",
			"error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song models.Song
		rows.Scan(
			&song.Id,
			&song.GroupId,
			&song.GroupName,
			&song.Title,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		songs = append(songs, song)
	}

	return songs, nil
}
