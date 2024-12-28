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

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	SHOWCASE = "showcase"
	ORG      = "org"
	USER     = "user"
	WORKER   = "worker"
)

var (
	ErrGenUUID     = errors.New("failed to generate uuid")
	ErrSetUUID     = errors.New("failed to set uuid")
	ErrGetUUID     = errors.New("failed to get uuid")
	ErrSaveURL     = errors.New("failed to save url")
	ErrSaveFile    = errors.New("failed to save file")
	ErrURLEmpty    = errors.New("empty url")
	ErrInvalidUUID = errors.New("invalid url's uuid")
	ErrDownloading = errors.New("failed to download file from s3")
	ErrUploading   = errors.New("failed to upload file to s3")
	ErrDeleting    = errors.New("failed to delete file into s3")
)

type S3UseCase struct {
	user   infrastructure.UserRepository
	org    infrastructure.OrgRepository
	minio  infrastructure.MediaStorage
	Logger *zap.Logger
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, mediaRepo infrastructure.MediaStorage, logger *zap.Logger) *S3UseCase {
	return &S3UseCase{
		user:   userRepo,
		org:    orgRepo,
		minio:  mediaRepo,
		Logger: logger,
	}
}

// Генерирует UUID и сохраняет файл
func (s *S3UseCase) Upload(ctx context.Context, dto *s3dto.CreateFileDTO) error {
	NewUUID, err := uuid.NewV7()
	if err != nil {
		return ErrGenUUID
	}
	URL := NewUUID.String()
	// Сохраняем UUID/URL showcase
	switch dto.Entity {
	case ORG:
		if err := s.org.OrgSetUUID(ctx, dto.EntityID, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrSetUUID, err)
		}
	case USER:
		if err := s.user.UserSetUUID(ctx, dto.EntityID, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrSetUUID, err)
		}
	case WORKER:
		if err := s.org.WorkerSetUUID(ctx, dto.EntityID, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrSetUUID, err)
		}
	case SHOWCASE:
		idStr := strconv.Itoa(dto.EntityID)
		b := strings.Builder{}
		b.Grow(len("org") + 1 + len(idStr) + len(URL))
		b.WriteString("org")
		b.WriteString(idStr)
		b.WriteString("/")
		b.WriteString(URL)
		URL = b.String()
		if err := s.org.OrgSaveShowcaseImageURL(ctx, dto.EntityID, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrSaveURL, err)
		}
	default:
		return fmt.Errorf("entity doesn't exist: %s", dto.Entity)
	}
	if err := s.minio.Upload(ctx, URL, dto.Name, dto.Size, dto.Reader); err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}
	return nil
}

func (s *S3UseCase) Download(ctx context.Context, URL string) (*s3dto.File, error) {
	if err := validateURL(URL); err != nil {
		return nil, err
	}
	obj, err := s.minio.Download(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDownloading, err)
	}
	// Получение информации о файле
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDownloading, err)
	}
	buffer := make([]byte, objInfo.Size)
	// Считывание файла
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%s: %w", ErrDownloading, err)
	}
	defer obj.Close()
	return &s3dto.File{
		Name:  objInfo.Key,
		Size:  objInfo.Size,
		Bytes: buffer,
	}, nil
}

func (s *S3UseCase) Delete(ctx context.Context, entity string, URL string) error {
	if err := validateURL(URL); err != nil {
		return err
	}
	switch entity {
	case ORG:
		if err := s.org.OrgDeleteURL(ctx, URL, false); err != nil {
			return fmt.Errorf("%s: %w", ErrDeleting, err)
		}
	case SHOWCASE:
		if err := s.org.OrgDeleteURL(ctx, URL, true); err != nil {
			return fmt.Errorf("%s: %w", ErrDeleting, err)
		}
	case USER:
		if err := s.user.UserDeleteURL(ctx, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrDeleting, err)
		}
	case WORKER:
		if err := s.org.WorkerDeleteURL(ctx, URL); err != nil {
			return fmt.Errorf("%s: %w", ErrDeleting, err)
		}
	}
	if err := s.minio.Delete(ctx, URL); err != nil {
		return fmt.Errorf("%s: %w", ErrDeleting, err)
	}
	return nil
}

func validateURL(URL string) error {
	components := strings.Split(URL, "/")
	if len(components) == 0 {
		return ErrURLEmpty
	}
	for _, v := range components {
		if err := uuid.Validate(v); err != nil {
			return fmt.Errorf("%s: %w", ErrInvalidUUID, err)
		}
	}
	return nil
}
