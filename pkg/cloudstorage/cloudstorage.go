package cloudstorage

//go:generate rm -f ./cloudstorage_mock.go
//go:generate mockgen -destination cloudstorage_mock.go -package cloudstorage -mock_names Client=GoMockClient -source cloudstorage.go

import (
	"context"
	"io"
	"time"
)

// CloudFile is a data structure for file in the cloud.
type CloudFile struct {
	// ObjectName is the object name of the file.
	ObjectName string
	// Size is the size of the file.
	Size int64
	// ContentType is the content type of the file.
	ContentType string
	// StorageLocation is the storage location of the file.
	StorageLocation string
}

// Client is an interface for Cloud Storage.
type Client interface {
	// Upload uploads the file to the Cloud Storage given by bucket, object
	// and return the CloudFile data structure.
	Upload(ctx context.Context, file io.Reader, bucket, object string, expires time.Time) (*CloudFile, error)
	// ListObjects lists the objects by given bucket from Cloud Storage.
	ListObjects(ctx context.Context, bucket string) ([]string, error)
	// Delete deletes the given bucket and object from Cloud Storage.
	Delete(ctx context.Context, bucket, object string) error
}
