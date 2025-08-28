package minio

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"

	"shop/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/slog"
)

type MinIO struct {
	Client *minio.Client
	Cnf    *config.Config
}

var bucketName = "shop"

func MinIOConnect(cnf *config.Config) (*MinIO, error) {
	// endpoint := cnf.MINIO_ENDPOINT
	accessKeyID := cnf.MINIO_ACCESS_KEY
	secretAccessKey := cnf.MINIO_SECRET_KEY

	minioClient, err := minio.New("31.187.74.228:9002", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		slog.Error("Failed to connect to MinIO: %v", err)
		return nil, err
	}

	// Create the bucket if it doesn't exist
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check if the bucket already exists
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			slog.Warn("Bucket already exists: %s\n", bucketName)
		} else {
			slog.Error("Error while making bucket %s: %v\n", bucketName, err)
		}
	} else {
		slog.Info("Successfully created bucket: %s\n", bucketName)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

	err = minioClient.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		slog.Error("Error while setting bucket policy: %v", err)
		return nil, err
	}

	return &MinIO{
		Client: minioClient,
		Cnf:    cnf,
	}, nil
}

func (m *MinIO) Upload(fileName, filePath string) (string, error) {
	ext := filepath.Ext(fileName)
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := m.Client.FPutObject(context.Background(), bucketName, fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		slog.Error("Error while uploading %s to bucket %s: %v\n", fileName, bucketName, err)
		return "", err
	}

	// serverHost := "minio"
	// domain := "31.187.74.228"
	// minioURL := fmt.Sprintf("http://%s:%d/%s/%s", "31.187.74.228", 9002, bucketName, fileName)

	minioURL := fmt.Sprintf("https://images.uz-dev-ai.uz/%s/%s", bucketName, fileName)


	return minioURL, nil
}
