package grpchandler_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/internal/usecases"
	"github.com/moemoe89/file-storage/pkg/cloudstorage"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockFileStorageService_UploadServer struct {
	grpc.ServerStream
	ctx  context.Context
	req  chan *rpc.UploadRequest
	resp chan *rpc.UploadResponse
	err  error
}

func (m *mockFileStorageService_UploadServer) Context() context.Context {
	return m.ctx
}

func (m *mockFileStorageService_UploadServer) Send(resp *rpc.UploadResponse) error {
	m.resp <- resp

	return m.err
}

func (m *mockFileStorageService_UploadServer) Recv() (*rpc.UploadRequest, error) {
	req, more := <-m.req
	if !more {
		return nil, errors.New("empty")
	}

	return req, m.err
}

func (m *mockFileStorageService_UploadServer) SendFromClient(req *rpc.UploadRequest) error {
	m.req <- req

	return m.err
}

func (m *mockFileStorageService_UploadServer) RecvToClient() (*rpc.UploadResponse, error) {
	response, more := <-m.resp
	if !more {
		return nil, errors.New("empty")
	}

	return response, m.err
}

func TestFileStorageServer_Upload(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  *rpc.UploadRequest
		mock *mockFileStorageService_UploadServer
		file io.Reader
	}

	type test struct {
		fields  fields
		args    args
		wantErr error

		beforeFunc func(*testing.T, *rpc.UploadRequest, *mockFileStorageService_UploadServer)
		afterFunc  func(*testing.T, *mockFileStorageService_UploadServer)
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Upload file, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.UploadRequest{
					Type: rpc.UploadType_UPLOAD_TYPE_FILE,
					Validation: &rpc.Validation{
						ContentTypes: []string{"text/plain"},
						MaxSize:      0,
					},
					Filename: "object",
					Bucket:   "bucket",
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 10),
					resp: make(chan *rpc.UploadResponse, 10),
				},
				file: strings.NewReader("my request"),
			}

			want := &cloudstorage.CloudFile{
				ObjectName:      args.req.GetFilename(),
				Size:            1024,
				ContentType:     "image/jpeg",
				StorageLocation: args.req.GetBucket() + "/" + args.req.GetFilename(),
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, args.file, args.req.GetBucket(), args.req.GetFilename(), time.Time{}).Return(want, nil).AnyTimes()

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: nil,

				beforeFunc: func(t *testing.T, req *rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
					t.Helper()

					err := m.SendFromClient(req)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T, m *mockFileStorageService_UploadServer) {
					t.Helper()

					close(m.req)
					close(m.resp)
				},
			}
		},
		"Given valid request of Upload file with empty filename and bucket, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.UploadRequest{
					Type: rpc.UploadType_UPLOAD_TYPE_FILE,
					Validation: &rpc.Validation{
						ContentTypes: []string{"text/plain"},
						MaxSize:      0,
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 10),
					resp: make(chan *rpc.UploadResponse, 10),
				},
				file: strings.NewReader("my request"),
			}

			want := &cloudstorage.CloudFile{
				ObjectName:      args.req.GetFilename(),
				Size:            1024,
				ContentType:     "image/jpeg",
				StorageLocation: args.req.GetBucket() + "/" + args.req.GetFilename(),
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, args.file, args.req.GetBucket(), args.req.GetFilename(), time.Time{}).Return(want, nil).AnyTimes()

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: nil,

				beforeFunc: func(t *testing.T, req *rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
					t.Helper()

					err := m.SendFromClient(req)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T, m *mockFileStorageService_UploadServer) {
					t.Helper()

					close(m.req)
					close(m.resp)
				},
			}
		},
		"Given valid request of Upload file without Content-Types, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.UploadRequest{
					Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
					Filename: "object",
					Bucket:   "bucket",
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 10),
					resp: make(chan *rpc.UploadResponse, 10),
				},
				file: strings.NewReader("my request"),
			}

			want := &cloudstorage.CloudFile{
				ObjectName:      args.req.GetFilename(),
				Size:            1024,
				ContentType:     "image/jpeg",
				StorageLocation: args.req.GetBucket() + "/" + args.req.GetFilename(),
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, args.file, args.req.GetBucket(), args.req.GetFilename(), time.Time{}).Return(want, nil).AnyTimes()

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				wantErr: nil,
				beforeFunc: func(t *testing.T, req *rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
					t.Helper()

					err := m.SendFromClient(req)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T, m *mockFileStorageService_UploadServer) {
					t.Helper()

					close(m.req)
					close(m.resp)
				},
			}
		},
		"Given valid request of Upload file, When it executed with invalid Content-Types, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.UploadRequest{
					Type: rpc.UploadType_UPLOAD_TYPE_FILE,
					Validation: &rpc.Validation{
						ContentTypes: []string{"image/jpeg"},
						MaxSize:      0,
					},
					Filename: "object",
					Bucket:   "bucket",
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 10),
					resp: make(chan *rpc.UploadResponse, 10),
				},
				file: strings.NewReader("my request"),
			}

			return test{
				args:    args,
				wantErr: fmt.Errorf("file mime type: %s doesn't match with the target: %v", "text/plain", []string{"image/jpeg"}),
				beforeFunc: func(t *testing.T, req *rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
					t.Helper()

					err := m.SendFromClient(req)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T, m *mockFileStorageService_UploadServer) {
					t.Helper()

					close(m.req)
					close(m.resp)
				},
			}
		},
		"Given valid request of Upload file, When it executed with mismatch offset, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx: ctx,
				req: &rpc.UploadRequest{
					Type: rpc.UploadType_UPLOAD_TYPE_FILE,
					Validation: &rpc.Validation{
						MaxSize: 1,
					},
					Filename: "object",
					Bucket:   "bucket",
					Detail: &rpc.UploadRequest_File{
						File: &rpc.FileUpload{
							Data:   []byte(`my request`),
							Offset: 1,
						},
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 10),
					resp: make(chan *rpc.UploadResponse, 10),
				},
				file: strings.NewReader("my request"),
			}

			return test{
				args:    args,
				wantErr: fmt.Errorf("unexpected offset, got %d, want %d", 1, 0),
				beforeFunc: func(t *testing.T, req *rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
					t.Helper()

					err := m.SendFromClient(req)
					assert.NoError(t, err)
				},
				afterFunc: func(t *testing.T, m *mockFileStorageService_UploadServer) {
					t.Helper()

					close(m.req)
					close(m.resp)
				},
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, tt.args.req, tt.args.mock)
			}

			go func() {
				err := sut.Upload(tt.args.mock)
				assert.Equal(t, tt.wantErr, err)

				if tt.afterFunc != nil {
					tt.afterFunc(t, tt.args.mock)
				}
			}()
		})
	}
}
