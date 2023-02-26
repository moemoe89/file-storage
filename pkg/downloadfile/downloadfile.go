package downloadfile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

//go:generate rm -f ./downloadfile_mock.go
//go:generate mockgen -destination downloadfile_mock.go -package downloadfile -mock_names DownloadFile=GoMockDownloadFile -source downloadfile.go

// DownloadFile defines the methods needed to download a file.
type DownloadFile interface {
	// DownloadByte downloads file by given target URL and return as byte.
	DownloadByte(ctx context.Context, targetURL string) ([]byte, error)
}

var (
	// errFailedSetHTTPClient is an error message when failed to set http client.
	errFailedSetHTTPClient = errors.New("failed to set diskstorage.http_client")
	// errInternal is an error for indicates from internal.
	errInternal = errors.New("error internal")
	// errInvalidArgument is an error for indicates incorrect argument.
	errInvalidArgument = errors.New("error invalid argument")
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type downloadFile struct {
	httpClient HTTPClient
}

func wrapErr(err1 error, err2 error) error {
	return fmt.Errorf("%v: %w", err1, err2)
}

// New returns a new DiskStorage instance with an empty internal buffer.
func New(opts ...Option) (DownloadFile, error) {
	d := new(downloadFile)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(d); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", wrapErr(err, errInvalidArgument))
		}
	}

	return d, nil
}

// DownloadByte downloads file by given target URL and return as byte.
func (d *downloadFile) DownloadByte(ctx context.Context, targetURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed create http request: %w", wrapErr(err, errInternal))
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", wrapErr(err, errInternal))
	}

	return body, nil
}
