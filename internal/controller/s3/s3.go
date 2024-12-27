package s3

import (
	"context"
	"fmt"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/s3dto"
	"timeline/internal/libs/custom"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type MediaStorage interface {
	Upload(ctx context.Context, dto s3dto.CreateFileDTO) error
	Download(ctx context.Context, URL string) (s3dto.File, error)
	Delete(ctx context.Context, URL string) error
}

type S3Ctrl struct {
	usecase MediaStorage
	json    jsoniter.API
	logger  *zap.Logger
}

func New(storage MediaStorage, Logger *zap.Logger, jsoniter jsoniter.API) *S3Ctrl {
	return &S3Ctrl{
		usecase: storage,
		logger:  Logger,
		json:    jsoniter,
	}
}

func (s3 *S3Ctrl) Upload(w http.ResponseWriter, r *http.Request) {
	// Обработка query параметров
	const maxUploadSize = 5 << 20 // 5 MB
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	metainfo := s3dto.DomenInfo{}
	if err := s3.json.NewDecoder(r.Body).Decode(&metainfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) <= 0 {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := s3dto.CreateFileDTO{
		DomenInfo: metainfo,
		Name:      fileInfo.Filename,
		Size:      fileInfo.Size,
		Reader:    fileReader,
	}
	if err := s3.usecase.Upload(r.Context(), dto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s3 *S3Ctrl) Download(w http.ResponseWriter, r *http.Request) {
	params := map[string]bool{
		"url": true,
	}
	if !validation.IsQueryValid(r, params) {
		http.Error(w, "query invalid", http.StatusBadRequest)
		return
	}
	queryContract := map[string]string{
		"url": "string",
	}
	query, err := custom.QueryParamsConv(queryContract, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f, err := s3.usecase.Download(r.Context(), query["url"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", f.Name))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(f.Bytes)
}

func (s3 *S3Ctrl) Delete(w http.ResponseWriter, r *http.Request) {
	params := map[string]bool{
		"url": true,
	}
	if !validation.IsQueryValid(r, params) {
		http.Error(w, "query invalid", http.StatusBadRequest)
		return
	}
	queryContract := map[string]string{
		"url": "string",
	}
	query, err := custom.QueryParamsConv(queryContract, r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s3.usecase.Delete(r.Context(), query["url"].(string)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
