//nolint:lll
package usecases

//go:generate rm -f ./file_storage_uc_mock.go
//go:generate mockgen -destination file_storage_uc_mock.go -package usecases -mock_names FileStorageUsecase=GoMockFileStorageUsecase -source file_storage_uc.go

import (
	"context"
	"io"
	"time"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"
	"github.com/moemoe89/file-storage/pkg/logging"
	"github.com/moemoe89/file-storage/pkg/trace"
)

// FileStorageUsecase defines File Storage related domain functionality.
type FileStorageUsecase interface {
	// Upload uploads file to storage.
	Upload(ctx context.Context, file io.Reader, bucket, object string, expires time.Time) (*cloudstorage.CloudFile, error)
	// List lists the files by given bucket name.
	List(ctx context.Context, bucket string) ([]string, error)
	// Delete deletes the file by given bucket and object name.
	Delete(ctx context.Context, bucket, object string) error
}

// FileStorageUsecase time interface implementation check.
var _ FileStorageUsecase = (*fileStorageUsecase)(nil)

// NewFileStorageUsecase returns FileStorageUsecase.
func NewFileStorageUsecase(
	trace trace.Tracer,
	logger logging.Logger,
	minio cloudstorage.Client,
) FileStorageUsecase {
	return &fileStorageUsecase{
		trace:  trace,
		logger: logger,
		minio:  minio,
	}
}

// fileStorageUsecase is a struct for usecase.
type fileStorageUsecase struct {
	trace  trace.Tracer
	logger logging.Logger
	minio  cloudstorage.Client
}
