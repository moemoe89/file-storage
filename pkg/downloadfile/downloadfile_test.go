package downloadfile

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

type mockedReadCloser struct {
	err error
}

func (m *mockedReadCloser) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func (m *mockedReadCloser) Close() error {
	return m.err
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: defaultOptions,
				},
				wantErr: nil,
			}
		},
		"Failed init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: []Option{
						WithHTTPClient(nil),
					},
				},
				wantErr: errInvalidArgument,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			_, err := New(tt.args.opts...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDownload(t *testing.T) {
	type args struct {
		ctx       context.Context
		targetURL string
	}

	type fields struct {
		mock *mockedHTTPClient
	}

	type test struct {
		args    args
		fields  fields
		want    []byte
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully download file": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: &http.Response{
							Body: io.NopCloser(strings.NewReader(`test download content`)),
						},
					},
				},
				want:    []byte(`test download content`),
				wantErr: nil,
			}
		},
		"Failed to create file": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
				},
				wantErr: errInternal,
			}
		},
		"Failed to create HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       nil,
					targetURL: "https://www.example.com",
				},
				wantErr: errInternal,
			}
		},
		"Failed do HTTP request": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						err: errInternal,
					},
				},
				wantErr: errInternal,
			}
		},
		"Failed to do Copy file": func(t *testing.T) test {
			t.Helper()

			mockedResp := &http.Response{
				Body: &mockedReadCloser{
					err: errInternal,
				},
			}

			return test{
				args: args{
					ctx:       context.Background(),
					targetURL: "https://www.example.com",
				},
				fields: fields{
					mock: &mockedHTTPClient{
						response: mockedResp,
						err:      nil,
					},
				},
				wantErr: errInternal,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			sut, err := New(WithHTTPClient(tt.fields.mock))
			assert.NoError(t, err)

			got, err := sut.DownloadByte(tt.args.ctx, tt.args.targetURL)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
