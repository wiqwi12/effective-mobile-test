package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
	"github.com/wiqwi12/effective-mobile-test/internal/infrastructure/externalServices"
	"github.com/wiqwi12/effective-mobile-test/internal/infrastructure/postgres/repository"
	"github.com/wiqwi12/effective-mobile-test/pkg"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
	"log/slog"
	"strings"
	"time"
)

type SongSrvc struct {
	SongRepo          *repository.SongRepository
	GroupRepo         *repository.GroupRepository
	MusicMetadataRepo *externalServices.MusicMetadataRepo
	VerseRepo         *repository.VerseRepository
	Logger            *logger.Logger
}

func NewSongSrvc(songRepo *repository.SongRepository, groupRepo *repository.GroupRepository, metaDataRepo *externalServices.MusicMetadataRepo, verseRepo *repository.VerseRepository, logger *logger.Logger) *SongSrvc {
	return &SongSrvc{
		SongRepo:          songRepo,
		GroupRepo:         groupRepo,
		MusicMetadataRepo: metaDataRepo,
		VerseRepo:         verseRepo,
		Logger:            logger,
	}
}

func (s *SongSrvc) CreateSong(ctx context.Context, request dto.CreateSongRequest) (dto.StandartResponse, error) {
	var song models.Song

	details, err := s.MusicMetadataRepo.GetSongDetails(ctx, request)
	if err != nil {
		s.Logger.Info.Error("Failed to get song metadata",
			"error", err,
			"group", request.Group,
			"title", request.Title)
		return dto.StandartResponse{}, err
	}

	groupExsist, err := s.GroupRepo.GroupExsist(ctx, request.Group)
	if err != nil {
		s.Logger.Info.Error("Failed to get song group",
			"error", err)
		return dto.StandartResponse{}, err
	}
	fmt.Println(groupExsist)

	var group models.Group
	if !groupExsist {
		group.Name = request.Group
		group.Id = uuid.New()
		err := s.GroupRepo.CreateGroup(group)
		if err != nil {
			s.Logger.Info.Error("Failed to create group",
				"error", err)
			return dto.StandartResponse{}, err
		}
		fmt.Println("Created group", group)
	}

	group, err = s.GroupRepo.GetGroupByName(request.Group)

	song.GroupId = group.Id
	song.GroupName = group.Name
	song.CreatedAt = time.Now()
	song.UpdatedAt = time.Now()
	song.Id = uuid.New()
	song.Link = details.Link
	song.Text = details.Text
	song.Title = request.Title
	song.ReleaseDate = details.ReleaseDate

	slog.Info("text:", song.Text)

	exists, err := s.SongRepo.SongExistsByDetails(ctx, song)
	if err != nil {
		s.Logger.Info.Error("Error checking if song exists",
			"error", err)
		return dto.StandartResponse{
			Message: "something went wrong",
			Error:   err.Error(),
		}, err
	}
	fmt.Println("song exsists", exists)
	if exists {
		return dto.StandartResponse{
			Error: errors.New("song already exists").Error(),
		}, errors.New("song already exists")
	}

	err = s.SongRepo.CreateSong(ctx, song)
	if err != nil {
		s.Logger.Info.Error("Failed to create song in database",
			"error", err)
		return dto.StandartResponse{
			Message: "something went wrong",
			Error:   err.Error(),
		}, err
	}

	err = s.ProcessVerses(ctx, song)
	if err != nil {
		s.Logger.Info.Error("Failed to process verses", "error", err.Error())
		return dto.StandartResponse{
			Message: "something vent wrong",
			Error:   err.Error(),
		}, err
	}

	resp := dto.StandartResponse{
		Song:    song,
		Message: "Songs succsessfully created",
	}

	return resp, nil
}

func (s *SongSrvc) GetSongById(ctx context.Context, request dto.GetSongByIdRequest) (dto.StandartResponse, error) {
	var resp dto.StandartResponse

	song, err := s.SongRepo.GetSongById(ctx, request.Id)
	if err != nil {
		s.Logger.Info.Error("Failed to get song by ID",
			"error", err,
			"song_id", request.Id)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return resp, err
	}

	resp.Song = song

	return resp, nil
}

func (s *SongSrvc) UpdateSong(ctx context.Context, request dto.UpdateSongRequest, songId uuid.UUID) (dto.StandartResponse, error) {
	resp := dto.StandartResponse{}

	exists, err := s.SongRepo.SongExsistsById(ctx, songId)
	if err != nil {
		s.Logger.Info.Error("Error checking if song exists",
			"error", err,
			"song_id", songId)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return resp, err
	}
	if !exists {
		s.Logger.Info.Error("Song does not exist",
			"song_id", songId)
		resp.Message = "song does not exist"
		return resp, errors.New("song does not exist")
	}

	originalSong, err := s.SongRepo.GetSongById(ctx, songId)
	if err != nil {
		s.Logger.Info.Error("Failed to get song by ID", err)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return resp, err
	}

	if request.Link != "" {
		originalSong.Link = request.Link
	}

	if request.Text != "" {
		originalSong.Text = request.Text
	}

	if request.ReleaseDate != "" {
		originalSong.ReleaseDate = request.ReleaseDate
	}

	if request.Title != "" {
		originalSong.Title = request.Title
	}

	if request.GroupName != "" {
		groupExsists, err := s.GroupRepo.GroupExsist(context.Background(), request.GroupName)
		if err != nil {
			s.Logger.Info.Error("some error",
				"error", err)
			resp.Message = "some error occured"
			resp.Error = err.Error()
			return resp, err
		}
		if groupExsists {
			group, err := s.GroupRepo.GetGroupByName(request.GroupName)
			if err != nil {
				s.Logger.Info.Error("some error",
					"error", err)
				resp.Message = "some error occured"
				resp.Error = err.Error()
			}
			originalSong.GroupName = group.Name
			originalSong.GroupId = group.Id
		} else {
			group := models.Group{
				Id:   uuid.New(),
				Name: request.GroupName,
			}
			err := s.GroupRepo.CreateGroup(group)
			if err != nil {
				s.Logger.Info.Error("Failed to create group",
					"error", err)
				resp.Message = "some error occured"
				resp.Error = err.Error()
				return resp, err
			}
			originalSong.GroupId = group.Id
			originalSong.GroupName = group.Name
		}

	}

	originalSong.UpdatedAt = time.Now()
	originalSong.Id = songId

	err = s.SongRepo.UpdateSong(ctx, originalSong)
	if err != nil {
		s.Logger.Info.Error("Failed to update song",
			"error", err,
			"song_id", songId)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return resp, err
	}

	resp.Message = "Songs succsessfully updated"
	resp.Song = originalSong

	return resp, nil
}

func (s *SongSrvc) DeleteSong(ctx context.Context, req dto.DeleteSongByIdRequest) (dto.StandartResponse, error) {
	resp := dto.StandartResponse{}

	err := s.SongRepo.DeleteSong(ctx, req.Id)
	if err != nil {
		s.Logger.Info.Error("Failed to delete song",
			"error", err,
			"song_id", req.Id)
		resp.Message = "Some error occured"
		resp.Error = err.Error()
		return resp, err
	}

	resp.Message = "Songs succsessfully deleted"
	return resp, nil
}

func (s *SongSrvc) GetSongWithFilter(ctx context.Context, req dto.FilteredRequest) (dto.SongsResponse, error) {
	var resp dto.SongsResponse

	if pkg.IsEmpty(req) {
		s.Logger.Info.Error("Empty filters provided")
		resp.Message = "request must contain at least one filer"
		resp.Error = errors.New("empty filters").Error()
		return resp, errors.New("empty filters")
	}

	songs, err := s.SongRepo.GetSongsWithFilter(ctx, req)
	if err != nil {
		s.Logger.Info.Error("Failed to get songs with filter",
			"error", err)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return resp, err
	}

	resp.Songs = songs
	return resp, nil
}

func (s *SongSrvc) ProcessVerses(ctx context.Context, song models.Song) error {

	var req dto.AddVersesRequest

	req.Verses = strings.Split(song.Text, "\\n\\n")

	req.Song = song

	err := s.VerseRepo.AddVerses(ctx, req)
	if err != nil {
		s.Logger.Info.Error(err.Error())
		return err
	}
	return nil
}

// В тз к заданию ничего не было сказано, поэтому сделал page based
func (s *SongSrvc) GetPaginatedVerses(ctx context.Context, request dto.PaginatedVersesRequest) (dto.PaginatedVersesResponse, error) {

	resp, err := s.VerseRepo.GetPaginatedVerses(ctx, request)
	if err != nil {
		s.Logger.Info.Error(err.Error())
		return dto.PaginatedVersesResponse{}, err
	}

	return resp, nil

}
