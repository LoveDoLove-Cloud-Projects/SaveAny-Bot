package minio

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/krau/SaveAny-Bot/common"
	config "github.com/krau/SaveAny-Bot/config/storage"
	"github.com/krau/SaveAny-Bot/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	config config.MinioStorageConfig
	client *minio.Client
}

func (m *Minio) Init(cfg config.StorageConfig) error {
	minioConfig, ok := cfg.(*config.MinioStorageConfig)
	if !ok {
		return fmt.Errorf("failed to cast minio config")
	}
	if err := minioConfig.Validate(); err != nil {
		return err
	}
	m.config = *minioConfig

	client, err := minio.New(m.config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.config.AccessKeyID, m.config.SecretAccessKey, ""),
		Secure: m.config.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create minio client: %w", err)
	}

	exists, err := client.BucketExists(context.Background(), m.config.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("bucket %s does not exist", m.config.BucketName)
	}

	m.client = client
	return nil
}

func (m *Minio) Type() types.StorageType {
	return types.StorageTypeMinio
}

func (m *Minio) Name() string {
	return m.config.Name
}

func (m *Minio) JoinStoragePath(task types.Task) string {
	return path.Join(m.config.BasePath, task.StoragePath)
}

func (m *Minio) Save(ctx context.Context, r io.Reader, storagePath string) error {
	common.Log.Infof("Saving file from reader to %s", storagePath)

	_, err := m.client.PutObject(ctx, m.config.BucketName, storagePath, r, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file to minio: %w", err)
	}

	return nil
}

func (m *Minio) Exists(ctx context.Context, storagePath string) bool {
	common.Log.Debugf("Checking if file exists at %s", storagePath)
	// TODO: test it.
	_, err := m.client.StatObject(ctx, m.config.BucketName, storagePath, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false // File does not exist
		}
		return false
	}

	return true
}
