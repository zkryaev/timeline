package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"timeline/internal/entity/dto/s3dto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	GALLERY = "gallery"
	BANNER  = "banner"
	ORG     = "org"
	USER    = "user"
	WORKER  = "worker"
)

var (
	ErrGenUUID    = errors.New("failed to generate uuid")
	ErrSetUUID    = errors.New("failed to set uuid")
	ErrGetUUID    = errors.New("failed to get uuid")
	ErrSaveURL    = errors.New("failed to save url")
	ErrSaveFile   = errors.New("failed to save file")
	ErrURLEmpty   = errors.New("empty url")
	ErrInvalidURL = errors.New("invalid url")
	ErrDownload   = errors.New("failed to download file from s3")
	ErrUpload     = errors.New("failed to upload file to s3")
	ErrDelete     = errors.New("failed to delete file into s3")
)

type S3UseCase struct {
	user  infrastructure.UserRepository
	org   infrastructure.OrgRepository
	minio infrastructure.MediaStorage
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, mediaRepo infrastructure.MediaStorage) *S3UseCase {
	return &S3UseCase{
		user:  userRepo,
		org:   orgRepo,
		minio: mediaRepo,
	}
}

// Генерирует UUID и сохраняет файл
func (s *S3UseCase) Upload(ctx context.Context, logger *zap.Logger, dto *s3dto.CreateFileDTO) error {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return ErrGenUUID
	}
	url := newUUID.String()
	var prevURL string
	LogFetched := fmt.Sprintf("Fetched current %s uuid", dto.Entity)
	LogSaved := fmt.Sprintf("New %s uuid has been saved", dto.Entity)
	switch {
	case dto.Entity == ORG:
		prevURL, err = s.org.OrgUUID(ctx, dto.EntityID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetUUID, err)
		}
		logger.Info(LogFetched)
		if err = s.org.OrgSetUUID(ctx, dto.EntityID, url); err != nil {
			return fmt.Errorf("%w: %w", ErrSetUUID, err)
		}
		logger.Info(LogSaved)
	case dto.Entity == USER:
		prevURL, err = s.user.UserUUID(ctx, dto.EntityID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetUUID, err)
		}
		logger.Info(LogFetched)
		if err = s.user.UserSetUUID(ctx, dto.EntityID, url); err != nil {
			return fmt.Errorf("%w: %w", ErrSetUUID, err)
		}
		logger.Info(LogSaved)
	case dto.Entity == WORKER:
		prevURL, err = s.org.WorkerUUID(ctx, dto.EntityID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrGetUUID, err)
		}
		logger.Info(LogFetched)
		if err = s.org.WorkerSetUUID(ctx, dto.EntityID, url); err != nil {
			return fmt.Errorf("%w: %w", ErrSetUUID, err)
		}
		logger.Info(LogSaved)
	case (dto.Entity == GALLERY) || (dto.Entity == BANNER):
		idStr := strconv.Itoa(dto.EntityID)
		b := strings.Builder{}
		b.Grow(len("org") + 1 + len(idStr) + len(url))
		b.WriteString("org")
		b.WriteString(idStr)
		b.WriteString("/")
		b.WriteString(url)
		url = b.String()
		meta := &models.ImageMeta{
			URL:     url,
			DomenID: dto.EntityID,
			Type:    dto.Entity,
		}
		if err = s.org.OrgSaveShowcaseImageURL(ctx, meta); err != nil {
			return fmt.Errorf("%w: %w", ErrSaveURL, err)
		}
		logger.Info(LogSaved)
	default:
		return fmt.Errorf("entity \"%s\" doesn't exist", dto.Entity)
	}
	if err = s.minio.Upload(ctx, url, dto.Name, dto.Size, dto.Reader); err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}
	logger.Info("Image has been saved")
	if prevURL != "" {
		if err = s.minio.Delete(ctx, prevURL); err != nil {
			logger.Error("failed image delete ", zap.Error(err))
		}
	}
	return nil
}

func (s *S3UseCase) Download(ctx context.Context, logger *zap.Logger, url string) (*s3dto.File, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}
	obj, err := s.minio.Download(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDownload, err)
	}
	logger.Info("Open S3 stream")
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDownload, err)
	}
	logger.Info("Retrieved image metadata")
	buffer := make([]byte, objInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%w: %w", ErrDownload, err)
	}
	logger.Info("Image data has been saved to buffer")
	defer obj.Close()
	return &s3dto.File{
		Name:  objInfo.Key,
		Size:  objInfo.Size,
		Bytes: buffer,
	}, nil
}

func (s *S3UseCase) Delete(ctx context.Context, logger *zap.Logger, entity string, url string) error {
	if err := validateURL(url); err != nil {
		return err
	}
	logger.Info("URL is valid")
	meta := &models.ImageMeta{
		URL:  url,
		Type: entity,
	}
	LogDeleteURL := fmt.Sprintf("%s's URL has been deleted", entity)
	switch {
	case entity == ORG:
		if err := s.org.OrgDeleteURL(ctx, meta); err != nil {
			return fmt.Errorf("%w: %w", ErrDelete, err)
		}
	case (entity == GALLERY) || (entity == BANNER):
		if err := s.org.OrgDeleteURL(ctx, meta); err != nil {
			return fmt.Errorf("%w: %w", ErrDelete, err)
		}
	case entity == USER:
		if err := s.user.UserDeleteURL(ctx, url); err != nil {
			return fmt.Errorf("%w: %w", ErrDelete, err)
		}
	case entity == WORKER:
		if err := s.org.WorkerDeleteURL(ctx, url); err != nil {
			return fmt.Errorf("%w: %w", ErrDelete, err)
		}
	}
	logger.Info(LogDeleteURL)
	if err := s.minio.Delete(ctx, url); err != nil {
		return fmt.Errorf("%w: %w", ErrDelete, err)
	}
	logger.Info("Image has been deleted")
	return nil
}

func validateURL(url string) error {
	components := strings.Split(url, "/")
	switch n := len(components); n {
	case 0:
		return ErrURLEmpty
	case 1: // uuid
		if err := uuid.Validate(components[0]); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidURL, err)
		}
	case 2: // domain-name/uuid
		if err := uuid.Validate(components[1]); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidURL, err)
		}
	default:
		return ErrInvalidURL
	}
	return nil
}
