package di

import (
	"context"
	"log"
	"os"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"
	"github.com/moemoe89/file-storage/pkg/cloudstorage/minio"
)

// GetMinio returns cloud storage Client interface.
func GetMinio() cloudstorage.Client {
	client, err := minio.New(
		context.Background(),
		minio.WithEndpoint(os.Getenv("MINIO_HOST")),
		minio.WithAccessKeyID(os.Getenv("MINIO_ACCESS_KEY_ID")),
		minio.WithSecretAccessKey(os.Getenv("MINIO_SECRET_ACCESS_KEY")),
	)
	if err != nil {
		log.Fatal("Minio Client", err)
	}

	return client
}
