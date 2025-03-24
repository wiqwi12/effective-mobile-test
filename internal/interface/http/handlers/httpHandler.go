package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"
	"github.com/wiqwi12/effective-mobile-test/internal/service"
	"net/http"
)

type Handler struct {
	srvc      *service.SongSrvc
	validator *validator.Validate
}

func NewHandler(srvc *service.SongSrvc, validator *validator.Validate) *Handler {
	return &Handler{srvc: srvc, validator: validator}
}

// @Summary Создать новую песню
// @Description Создает новую песню в базе данных
// @Tags songs
// @Accept json
// @Produce json
// @Param request body dto.CreateSongRequest true "Данные для создания песни"
// @Success 200 {object} dto.StandartResponse "Песня успешно создана"
// @Failure 400 {object} dto.StandartResponse "Ошибка в запросе"
// @Router /api/song [post]
func (h *Handler) CreateSongHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var resp dto.StandartResponse

	var req dto.CreateSongRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "something wrong with request"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, err = h.srvc.CreateSong(context.Background(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

}

// @Summary Получить песню по ID
// @Description Возвращает данные песни по её ID
// @Tags songs
// @Produce json
// @Param id path string true "ID песни" format(uuid)
// @Success 200 {object} dto.StandartResponse "Данные песни"
// @Failure 400 {object} dto.StandartResponse "Ошибка в запросе"
// @Router /api/song/{id} [get]
func (h *Handler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp dto.StandartResponse

	songId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid song id"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	req := dto.GetSongByIdRequest{
		Id: songId,
	}
	err = h.validator.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid song id"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, err = h.srvc.GetSongById(context.Background(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "some error occured"
		resp.Error = err.Error()
		return
	}
	json.NewEncoder(w).Encode(resp)
	return
}

// @Summary Обновить песню
// @Description Обновляет данные существующей песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "ID песни" format(uuid)
// @Param request body dto.UpdateSongRequest true "Данные для обновления песни"
// @Success 200 {object} dto.StandartResponse "Песня успешно обновлена"
// @Failure 400 {object} dto.StandartResponse "Ошибка в запросе"
// @Router /api/song/{id} [put]
func (h *Handler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {

	var resp dto.StandartResponse
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid song id"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	req := dto.UpdateSongRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err = h.srvc.UpdateSong(context.Background(), req, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
	return
}

// @Summary Удалить песню
// @Description Удаляет песню из базы данных
// @Tags songs
// @Produce json
// @Param id path string true "ID песни" format(uuid)
// @Success 200 {object} dto.StandartResponse "Песня успешно удалена"
// @Failure 400 {object} dto.StandartResponse "Ошибка в запросе"
// @Router /api/song/{id} [delete]
func (h *Handler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {

	resp := dto.StandartResponse{}
	w.Header().Set("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid song id"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	request := dto.DeleteSongByIdRequest{
		Id: id,
	}
	err = h.validator.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid song id"
		resp.Error = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, err = h.srvc.DeleteSong(context.Background(), request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
	return
}

// @Summary Получить песни с фильтрацией
// @Description Возвращает список песен, соответствующих фильтрам
// @Tags songs
// @Accept json
// @Produce json
// @Param request body dto.FilteredRequest true "Параметры фильтрации"
// @Success 200 {object} dto.SongsResponse "Список найденных песен"
// @Failure 400 {object} dto.SongsResponse "Ошибка в запросе"
// @Router /api/song [get]
func (h *Handler) GetSongWithFilter(w http.ResponseWriter, r *http.Request) {
	resp := dto.SongsResponse{}
	w.Header().Set("Content-Type", "application/json")

	req := dto.FilteredRequest{}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.srvc.GetSongWithFilter(context.Background(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)

}

// @Summary Получить куплеты песни с пагинацией
// @Description Возвращает куплеты песни с указанной пагинацией
// @Tags verses
// @Accept json
// @Produce json
// @Param id path string true "ID песни" format(uuid)
// @Param request body dto.PaginatedVersesRequest true "Параметры пагинации"
// @Success 200 {object} dto.PaginatedVersesResponse "Список куплетов песни"
// @Failure 400 {object} dto.StandartResponse "Ошибка в запросе"
// @Router /api/verses/{id} [get]
func (h *Handler) GetPaginatedVerses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	songId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   err.Error(),
			Message: "Invalid song ID",
		})
		return
	}

	var req dto.PaginatedVersesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   err.Error(),
			Message: "Failed to decode request body",
		})
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   err.Error(),
			Message: "Invalid song ID",
		})
	}

	req.SongId = songId

	if req.Page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   "page must be at least 1",
			Message: "Validation failed",
		})
		return
	}

	if req.Limit < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   "limit must be at least 1",
			Message: "Validation failed",
		})
		return
	}
	
	resp, err := h.srvc.GetPaginatedVerses(context.Background(), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.StandartResponse{
			Error:   err.Error(),
			Message: "Failed to get paginated verses",
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}
