package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
	"timeline/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	ErrConnection = errors.New("client is offline")
)

type Minio struct {
	cfg  config.S3
	conn *minio.Client
}

func New(cfg config.S3) *Minio {
	return &Minio{
		cfg: cfg,
	}
}

// Подключиться к MinIO, вернет ошибку если не удалось или сервис недоступен
func (m *Minio) Connect() error {
	minioClient, err := minio.New(m.cfg.Host+":"+m.cfg.DataPort, &minio.Options{
		Creds:  credentials.NewStaticV4(m.cfg.User, m.cfg.Password, ""),
		Secure: m.cfg.UseSSL,
	})
	if err != nil {
		return err
	}
	m.conn = minioClient
	StopHealthChecking, err := m.conn.HealthCheck(5 * time.Second)
	defer StopHealthChecking()
	if !m.conn.IsOnline() {
		return ErrConnection
	}
	return nil
}

func (m *Minio) Upload(ctx context.Context, URL string, fileName string, fileSize int64, reader io.Reader) error {
	if exists, errBucketExists := m.conn.BucketExists(ctx, m.cfg.DefaultBucket); errBucketExists != nil || !exists {
		err := m.conn.MakeBucket(ctx, m.cfg.DefaultBucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%s -> %w", errBucketExists, err)
		}
	}
	_, err := m.conn.PutObject(ctx, m.cfg.DefaultBucket, URL, reader, fileSize, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"Name": fileName,
		},
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Minio) Download(ctx context.Context, URL string) (*minio.Object, error) {
	obj, err := m.conn.GetObject(ctx, m.cfg.DefaultBucket, URL, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, err
}

func (m *Minio) Delete(ctx context.Context, URL string) error {
	if err := m.conn.RemoveObject(ctx, m.cfg.DefaultBucket, URL, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}