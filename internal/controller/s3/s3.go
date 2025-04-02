package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/s3dto"
	"timeline/internal/libs/custom"

	"go.uber.org/zap"
)

type S3UseCase interface {
	Upload(ctx context.Context, logger *zap.Logger, dto *s3dto.CreateFileDTO) error
	Download(ctx context.Context, logger *zap.Logger, URL string) (*s3dto.File, error)
	Delete(ctx context.Context, logger *zap.Logger, entity string, URL string) error
}

type S3Ctrl struct {
	usecase S3UseCase
	logger  *zap.Logger
}

func New(storage S3UseCase, logger *zap.Logger) *S3Ctrl {
	return &S3Ctrl{
		usecase: storage,
		logger:  logger,
	}
}

// UploadFileHandler handles file uploads.
// @Summary Upload a file
// @Description Upload a single file with metadata. For orgs showcase: entity=showcase, but entity_id = org_id
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param entity formData string true "Entity associated with the file"
// @Param entityID formData int true "Entity ID"
// @Param file formData file true "File to upload"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /media [post]
func (s3 *S3Ctrl) Upload(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := s3.logger.With(zap.String("uuid", uuid))
	const maxUploadSize = 2 << 20 // 2 MB
	if err := r.ParseMultipartForm(maxUploadSize); err != nil && !errors.Is(err, io.EOF) {
		logger.Error("ParseMultipartForm", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	entityID, err := strconv.Atoi(r.FormValue("entityID"))
	entity := r.FormValue("entity")
	switch {
	case err != nil || entityID <= 0:
		logger.Error("entity invalid", zap.Int("entity_id", entityID), zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	case entity == "":
		logger.Error("entity is empty")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	file, meta, err := r.FormFile("file")
	switch {
	case err != nil || meta.Size == 0:
		logger.Error("file not found", zap.Int64("file size", meta.Size), zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	case meta.Size > maxUploadSize:
		logger.Error("file size is too big", zap.Int64("file size", meta.Size))
		http.Error(w, "", http.StatusRequestEntityTooLarge)
		return
	}
	contentType := meta.Header.Get("Content-Type")
	if !validation.IsImage(contentType) {
		logger.Error("IsImage", zap.String("invalid image type", contentType))
		http.Error(w, "", http.StatusUnsupportedMediaType)
		return
	}
	domen := s3dto.DomenInfo{
		Entity:   entity,
		EntityID: entityID,
	}
	dto := &s3dto.CreateFileDTO{
		DomenInfo: domen,
		Name:      meta.Filename,
		Size:      meta.Size,
		Reader:    file,
	}
	if err = s3.usecase.Upload(r.Context(), logger, dto); err != nil {
		logger.Error("Upload", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DownloadFileHandler handles file downloads.
// @Summary Download a file
// @Description Download a file by its URL
// @Tags Media
// @Accept json
// @Produce application/octet-stream
// @Param url query string true "File URL"
// @Success 200 {file} string "File data"
// @Failure 400
// @Failure 500
// @Router /media [get]
func (s3 *S3Ctrl) Download(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := s3.logger.With(zap.String("uuid", uuid))
	params := map[string]bool{
		"url": true,
	}
	if err := validation.IsQueryValid(r, params); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	queryContract := map[string]string{
		"url": "string",
	}
	query, err := custom.QueryParamsConv(queryContract, r.URL.Query())
	if err != nil {
		logger.Error("QueryParamsConv", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	f, err := s3.usecase.Download(r.Context(), logger, query["url"].(string))
	if err != nil {
		logger.Error("Download", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", f.Name))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(f.Bytes)
}

// DeleteFileHandler handles file deletions.
// @Summary Delete a file
// @Description Delete a file by its URL and associated entity
// @Tags Media
// @Accept json
// @Produce json
// @Param url query string true "File URL"
// @Param entity query string true "Associated entity"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /media [delete]
func (s3 *S3Ctrl) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := s3.logger.With(zap.String("uuid", uuid))
	params := map[string]bool{
		"url":    true,
		"entity": true,
	}
	if err := validation.IsQueryValid(r, params); err != nil {
		logger.Error("IsQueryValid", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	queryContract := map[string]string{
		"url":    "string",
		"entity": "string",
	}
	query, err := custom.QueryParamsConv(queryContract, r.URL.Query())
	if err != nil {
		logger.Error("QueryParamsConv", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err = s3.usecase.Delete(r.Context(), logger, query["entity"].(string), query["url"].(string)); err != nil {
		logger.Error("Delete", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
