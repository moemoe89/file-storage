package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// gcsClient is a struct for gcs client.
type minioClient struct {
	*minio.Client

	endpoint        string
	accessKeyID     string
	secretAccessKey string
	token           string
	secure          bool

	bucketExists map[string]struct{}
}

// New returns Cloud Storage interface implementations.
func New(ctx context.Context, opts ...Option) (cloudstorage.Client, error) {
	m := new(minioClient)

	m.bucketExists = make(map[string]struct{}, 0)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(m); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	var err error

	m.Client, err = minio.New(m.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.accessKeyID, m.secretAccessKey, m.token),
		Secure: m.secure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed connect to minio client: %w", err)
	}

	return m, nil
}

// Upload uploads the file to the Cloud Storage given by bucket, object
// and return the CloudFile data structure.
func (m *minioClient) Upload(
	ctx context.Context, file io.Reader, bucket, object string, expires time.Time,
) (*cloudstorage.CloudFile, error) {
	err := m.checkBucket(ctx, bucket)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	fileSize := int64(buf.Len())

	info, err := m.Client.PutObject(
		ctx, bucket, object, bytes.NewReader(buf.Bytes()), fileSize, minio.PutObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	return &cloudstorage.CloudFile{
		ObjectName:      info.Key,
		Size:            info.Size,
		StorageLocation: info.Bucket,
	}, nil
}

func (m *minioClient) checkBucket(ctx context.Context, bucket string) error {
	if _, ok := m.bucketExists[bucket]; ok {
		return nil
	}

	exists, err := m.Client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existance: %w", err)
	}

	if exists {
		m.bucketExists[bucket] = struct{}{}

		return nil
	}

	err = m.Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %s: %w", bucket, err)
	}

	m.bucketExists[bucket] = struct{}{}

	return nil
}

// Delete deletes the given bucket and object from Cloud Storage.
func (m *minioClient) Delete(ctx context.Context, bucket, object string) error {
	_, err := m.Client.StatObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}

	err = m.Client.RemoveObject(ctx, bucket, object, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// ListObjects lists the objects by given bucket from Cloud Storage.
func (m *minioClient) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	objects := make([]string, 0)

	objectCh := m.Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to get list object: %w", object.Err)
		}

		objects = append(objects, object.Key)
	}

	return objects, nil
}
