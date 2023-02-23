package minio

import "errors"

// Option configures Minio client
type Option func(m *minioClient) error

// defaultOptions is a default configuration for Minio.
var defaultOptions = []Option{
	WithEndpoint("localhost:9000"),
	WithSecure(false),
}

var (
	// errFailedSetEndpoint is an error message when failed to set endpoint.
	errFailedSetEndpoint = errors.New("failed to set minio.endpoint")
	// errFailedSetAccessKeyID is an error message when failed to set access key id.
	errFailedSetAccessKeyID = errors.New("failed to set minio.access_key_id")
	// errFailedSetSecretAccessKey is an error message when failed to set secret access key.
	errFailedSetSecretAccessKey = errors.New("failed to set minio.secret_access_key")
	// errFailedSetToken is an error message when failed to set token.
	errFailedSetToken = errors.New("failed to set minio.token")
)

// WithEndpoint returns an option that set the endpoint URL.
func WithEndpoint(value string) Option {
	return func(m *minioClient) error {
		if len(value) == 0 {
			return errFailedSetEndpoint
		}

		m.endpoint = value

		return nil
	}
}

// WithAccessKeyID returns an option that set the access key id.
func WithAccessKeyID(value string) Option {
	return func(m *minioClient) error {
		if len(value) == 0 {
			return errFailedSetAccessKeyID
		}

		m.accessKeyID = value

		return nil
	}
}

// WithSecretAccessKey returns an option that set the secret access key.
func WithSecretAccessKey(value string) Option {
	return func(m *minioClient) error {
		if len(value) == 0 {
			return errFailedSetSecretAccessKey
		}

		m.secretAccessKey = value

		return nil
	}
}

// WithToken returns an option that set the token.
func WithToken(value string) Option {
	return func(m *minioClient) error {
		if len(value) == 0 {
			return errFailedSetToken
		}

		m.token = value

		return nil
	}
}

// WithSecure returns an option that set the secure connection.
func WithSecure(value bool) Option {
	return func(m *minioClient) error {

		m.secure = value

		return nil
	}
}
