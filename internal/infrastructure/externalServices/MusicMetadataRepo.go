package externalServices

import (
	"context"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
	"net/http"
)

type MusicMetadataRepo struct {
	client *http.Client
	ApiUrl string
	Logger *logger.Logger
}

func NewExternalRepo(apiUrl string, logger *logger.Logger) *MusicMetadataRepo {
	return &MusicMetadataRepo{
		client: http.DefaultClient,
		ApiUrl: apiUrl,
		Logger: logger,
	}
}

func (r *MusicMetadataRepo) GetSongDetails(ctx context.Context, request dto.CreateSongRequest) (dto.SongDetailResponse, error) {

	resp := dto.SongDetailResponse{
		ReleaseDate: "11.11.2011",
		Text:        `Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight`,
		Link:        "somelink.com",
	}

	return resp, nil

	//apiUrl := fmt.Sprintf("%s/info?group=%s&song=%s", r.ApiUrl, url.QueryEscape(request.Group), url.QueryEscape(request.Title))
	//
	//r.Logger.Info.Info("Requesting song metadata from external API",
	//	"group", request.Group,
	//	"title", request.Title)
	//
	//req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	//if err != nil {
	//	r.Logger.Info.Error("Failed to create external API request",
	//		"error", err,
	//		"url", apiUrl)
	//	return dto.SongDetailResponse{}, fmt.Errorf("failed to create request: %w", err)
	//}
	//
	//resp, err := r.client.Do(req)
	//if err != nil {
	//	r.Logger.Info.Error("Failed to send request to external API",
	//		"error", err,
	//		"url", apiUrl)
	//	return dto.SongDetailResponse{}, fmt.Errorf("failed to send request: %w", err)
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	r.Logger.Info.Error("External API returned non-OK status",
	//		"status_code", resp.StatusCode,
	//		"url", apiUrl)
	//	return dto.SongDetailResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	//}
	//
	//var songDetails dto.SongDetailResponse
	//if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
	//	r.Logger.Info.Error("Failed to decode response from external API",
	//		"error", err,
	//		"url", apiUrl)
	//	return dto.SongDetailResponse{}, fmt.Errorf("failed to decode response: %w", err)
	//}
	//
	//r.Logger.Info.Info("Successfully retrieved song metadata from external API",
	//	"group", request.Group,
	//	"title", request.Title)
	//
	//return songDetails, nil
}
