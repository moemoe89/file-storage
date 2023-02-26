package downloadfile

import (
	"net/http"
)

// Option configures downloadFile.
type Option func(d *downloadFile) error

// defaultOptions is a default configuration for downloadFile.
var defaultOptions = []Option{
	WithHTTPClient(http.DefaultClient),
}

// WithHTTPClient returns an option that set the http client.
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(d *downloadFile) error {
		if httpClient == nil {
			return errFailedSetHTTPClient
		}

		d.httpClient = httpClient

		return nil
	}
}
