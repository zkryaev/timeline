package infrastructure

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type S3 interface {
	Connect() error
	MediaStorage
}

type MediaStorage interface {
	Upload(ctx context.Context, URL string, fileName string, fileSize int64, reader io.Reader) error
	Download(ctx context.Context, URL string) (*minio.Object, error)
	Delete(ctx context.Context, URL string) error
}
