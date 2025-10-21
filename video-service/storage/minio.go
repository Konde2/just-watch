package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client *minio.Client
	bucket string
}

func NewMinIOStorage(endpoint, accessKey, secretKey, bucket string) (*MinIOStorage, error) {
	useSSL := false
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		exists, err := minioClient.BucketExists(ctx, bucket)
		if err == nil && exists {
			log.Printf("Bucket %s already exists", bucket)
		} else {
			return nil, err
		}
	}

	return &MinIOStorage{client: minioClient, bucket: bucket}, nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	return err
}

func (s *MinIOStorage) GetObject(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *MinIOStorage) GetObjectSize(ctx context.Context, objectName string) (int64, error) {
	info, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}
	return info.Size, nil
}
