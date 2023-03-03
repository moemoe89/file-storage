package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	// publicHost is a public host for Google Cloud Storage e.g.
	// https://storage.googleapis.com/example-test/test.jpg
	publicHost = "https://storage.googleapis.com"
)

// gcsClient is a struct for gcs client.
type gcsClient struct {
	*storage.Client
}

// New returns Cloud Storage interface implementations.
func New(ctx context.Context, opts ...option.ClientOption) (cloudstorage.Client, error) {
	g := new(gcsClient)

	var err error

	g.Client, err = storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// Upload uploads the file to the Cloud Storage given by object
// and return the public url of the file.
func (g *gcsClient) Upload(
	ctx context.Context, file io.Reader, bucket, object string, expires time.Time,
) (*cloudstorage.CloudFile, error) {
	wc := g.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("failed to copy file %s to bucket %s: %w", object, bucket, err)
	}

	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer %s to bucket %s: %w", object, bucket, err)
	}

	// Get the uploaded file information
	attrs := wc.Attrs()

	cloudFile := &cloudstorage.CloudFile{
		ObjectName:      attrs.Name,
		Size:            attrs.Size,
		ContentType:     attrs.ContentType,
		StorageLocation: g.buildURL(bucket, object),
	}

	// immediately do return if expires time not configured.
	if expires.IsZero() {
		return cloudFile, nil
	}

	url, err := g.signedURL(bucket, object, expires)
	if err != nil {
		return nil, err
	}

	cloudFile.StorageLocation = url

	return cloudFile, nil
}

func (g *gcsClient) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	return nil, errors.New("unimplemented")
}

// Delete deletes the given object from Cloud Storage.
func (g *gcsClient) Delete(ctx context.Context, bucket, object string) error {
	err := g.Bucket(bucket).Object(object).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete file %s on bucket %s: %w", object, bucket, err)
	}

	return nil
}

// signedURL signed the object from cloud storage with expires time.
func (g *gcsClient) signedURL(bucket, object string, expires time.Time) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  http.MethodGet,
		Expires: expires,
	}

	url, err := g.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("failed to signed url for object %s in bucket %s: %w", object, bucket, err)
	}

	return url, nil
}

// buildURL builds the object URL from cloud storage.
func (g *gcsClient) buildURL(bucket, object string) string {
	return fmt.Sprintf("%s/%s/%s", publicHost, bucket, object)
}
