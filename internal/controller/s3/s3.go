package s3

import (
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type S3Ctrl struct {
	Client *minio.Client
	Logger *zap.Logger
}

func New() {}
