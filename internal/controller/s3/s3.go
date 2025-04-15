package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/s3dto"

	"go.uber.org/zap"
)

type S3UseCase interface {
	Upload(ctx context.Context, logger *zap.Logger, dto *s3dto.CreateFileDTO) error
	Download(ctx context.Context, logger *zap.Logger, URL string) (*s3dto.File, error)
	Delete(ctx context.Context, logger *zap.Logger, req s3dto.DeleteReq) error
}

type S3Ctrl struct {
	usecase  S3UseCase
	logger   *zap.Logger
	settings *scope.Settings
}

func New(storage S3UseCase, logger *zap.Logger, settings *scope.Settings) *S3Ctrl {
	return &S3Ctrl{
		usecase:  storage,
		logger:   logger,
		settings: settings,
	}
}

// UploadFileHandler handles file uploads.
// @Summary Upload a file
// @Description Upload a single file with metadata. For orgs showcase: entity=showcase, but entity_id = org_id
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200
// @Failure 400
// @Failure 413
// @Failure 415
// @Failure 500
// @Router /media [post]
func (s3 *S3Ctrl) Upload(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(s3.settings, s3.logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(s3.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
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
	switch {
	case (entity == scope.BANNER || entity == scope.GALLERY || entity == scope.ORG || entity == scope.WORKER) && tdata.IsOrg:
		if !(entity == scope.WORKER) {
			entityID = tdata.ID
		}
	case entity == scope.USER && !tdata.IsOrg:
		entityID = tdata.ID
	default:
		var caller string
		if tdata.IsOrg {
			caller = scope.ORG
		} else {
			caller = scope.USER
		}
		logger.Info("forbid to upload media", zap.String("caller", caller), zap.Int("caller_id", tdata.ID), zap.String("entity", entity))
		http.Error(w, "", http.StatusUnauthorized)
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
		TData:    tdata,
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
// @Param url query string true "url for s3"
// @Success 200 {file} string "bytes"
// @Failure 400
// @Failure 500
// @Router /media [get]
func (s3 *S3Ctrl) Download(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(s3.settings, s3.logger, r.Context())
	url := query.NewParamString(scope.URL, true)
	params := query.NewParams(s3.settings, url)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	f, err := s3.usecase.Download(r.Context(), logger, url.Val)
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
// @Param url query string true "url for s3"
// @Param entity query string true "banner, gallery, org, user, worker"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /media [delete]
func (s3 *S3Ctrl) Delete(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(s3.settings, s3.logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(s3.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var (
		url    = query.NewParamString(scope.URL, true)
		entity = query.NewParamString(scope.ENTITY, true)
	)
	params := query.NewParams(s3.settings, url, entity)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	switch {
	case (entity.Val == scope.BANNER || entity.Val == scope.GALLERY || entity.Val == scope.ORG || entity.Val == scope.WORKER) && tdata.IsOrg:
	case entity.Val == scope.USER && !tdata.IsOrg:
	default:
		var caller string
		if tdata.IsOrg {
			caller = "org"
		} else {
			caller = "user"
		}
		logger.Info("forbid to delete media", zap.String("caller", caller), zap.Int("caller_id", tdata.ID), zap.String("entity", entity.Val), zap.String("url", url.Val))
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	req := s3dto.DeleteReq{Url: url.Val, Entity: entity.Val, TData: tdata}
	if err := s3.usecase.Delete(r.Context(), logger, req); err != nil {
		logger.Error("Delete", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
