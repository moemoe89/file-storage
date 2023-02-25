package usecases

import (
	"context"
	"io"
	"time"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"
)

// Upload uploads file to storage.
func (u *fileStorageUsecase) Upload(
	ctx context.Context, file io.Reader, bucket, object string, expires time.Time,
) (*cloudstorage.CloudFile, error) {
	ctx, span := u.trace.StartSpan(ctx, "UC.Upload", nil)
	defer span.End()

	return u.minio.Upload(ctx, file, bucket, object, expires)
}

// List lists the files by given bucket name.
func (u *fileStorageUsecase) List(ctx context.Context, bucket string) ([]string, error) {
	ctx, span := u.trace.StartSpan(ctx, "UC.List", nil)
	defer span.End()

	return u.minio.ListObjects(ctx, bucket)
}

// Delete deletes the file by given bucket and object name.
func (u *fileStorageUsecase) Delete(ctx context.Context, bucket, object string) error {
	ctx, span := u.trace.StartSpan(ctx, "UC.Delete", nil)
	defer span.End()

	return u.minio.Delete(ctx, bucket, object)
}
