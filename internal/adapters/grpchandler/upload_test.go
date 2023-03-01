package grpchandler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/internal/usecases"
	"github.com/moemoe89/file-storage/pkg/cloudstorage"
	"github.com/moemoe89/file-storage/pkg/downloadfile"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockFileStorageService_UploadServer struct {
	grpc.ServerStream
	ctx               context.Context
	req               chan *rpc.UploadRequest
	resp              chan *rpc.UploadResponse
	errSend           error
	errRecv           error
	errSendFromClient error
	errRecvToClient   error
}

func (m *mockFileStorageService_UploadServer) Context() context.Context {
	return m.ctx
}

func (m *mockFileStorageService_UploadServer) Send(resp *rpc.UploadResponse) error {
	m.resp <- resp

	return m.errSend
}

func (m *mockFileStorageService_UploadServer) Recv() (*rpc.UploadRequest, error) {
	req, more := <-m.req
	if !more {
		return nil, errors.New("empty")
	}

	return req, m.errRecv
}

func (m *mockFileStorageService_UploadServer) SendFromClient(req *rpc.UploadRequest) error {
	m.req <- req

	return m.errSendFromClient
}

func (m *mockFileStorageService_UploadServer) RecvToClient() (*rpc.UploadResponse, error) {
	response, more := <-m.resp
	if !more {
		return nil, errors.New("empty")
	}

	return response, m.errRecvToClient
}

func TestFileStorageServer_Upload(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  []*rpc.UploadRequest
		mock *mockFileStorageService_UploadServer
	}

	type test struct {
		fields     fields
		args       args
		wantErr    error
		beforeFunc func(*testing.T, []*rpc.UploadRequest, *mockFileStorageService_UploadServer)
		afterFunc  func(*testing.T, *mockFileStorageService_UploadServer)
	}

	defaultAfterFunc := func(t *testing.T, m *mockFileStorageService_UploadServer) {
		t.Helper()

		close(m.req)
		close(m.resp)
	}

	defaultBeforeFunc := func(t *testing.T, req []*rpc.UploadRequest, m *mockFileStorageService_UploadServer) {
		t.Helper()

		for i := range req {
			err := m.SendFromClient(req[i])
			assert.NoError(t, err)
		}
	}

	defaultError := errors.New("error")

	ctx := context.Background()

	defaultMock := func(ctx context.Context, n int) *mockFileStorageService_UploadServer {
		return &mockFileStorageService_UploadServer{
			ctx:  ctx,
			req:  make(chan *rpc.UploadRequest, n),
			resp: make(chan *rpc.UploadResponse, n),
		}
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Upload file, When it executed successfully with Content-Type validation, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							ContentTypes: []string{"text/plain"},
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed with context canceled, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							ContentTypes: []string{"text/plain"},
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed with context error, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now())
			defer cancel()

			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							ContentTypes: []string{"text/plain"},
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    context.DeadlineExceeded,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it failed to send the stream response, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							ContentTypes: []string{"text/plain"},
						},
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:     ctx,
					req:     make(chan *rpc.UploadRequest, 1),
					resp:    make(chan *rpc.UploadResponse, 1),
					errSend: defaultError,
				},
			}

			return test{
				args:       args,
				wantErr:    defaultError,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed successfully with empty filename and bucket, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type: rpc.UploadType_UPLOAD_TYPE_FILE,
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed successfully without Content-Types, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed with invalid Content-Types, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type: rpc.UploadType_UPLOAD_TYPE_FILE,
						Validation: &rpc.Validation{
							ContentTypes: []string{"image/jpeg"},
						},
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    fmt.Errorf("file mime type: %s doesn't match with the target: %v", "text/plain", []string{"image/jpeg"}),
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed with mismatch offset, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Detail: &rpc.UploadRequest_File{
							File: &rpc.FileUpload{
								Data:   []byte(`my request`),
								Offset: 1,
							},
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			return test{
				args:       args,
				wantErr:    fmt.Errorf("unexpected offset, got %d, want %d", 1, 0),
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload file, When it executed and exceed max size, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							MaxSize: 1,
						},
						Detail: &rpc.UploadRequest_File{
							File: &rpc.FileUpload{
								Data:   []byte(`my`),
								Offset: 0,
							},
						},
					},
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
						Validation: &rpc.Validation{
							MaxSize: 1,
						},
						Detail: &rpc.UploadRequest_File{
							File: &rpc.FileUpload{
								Data:   []byte(`request`),
								Offset: 2,
							},
						},
					},
				},
				mock: defaultMock(ctx, 2),
			}

			return test{
				args:       args,
				wantErr:    fmt.Errorf("file size: %d exceed the target: %d", 2, 1),
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload URL, When it executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			targetURL := "http://test.com/test.txt"

			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_URL,
						Filename: "object",
						Bucket:   "bucket",
						Detail: &rpc.UploadRequest_Url{
							Url: targetURL,
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			byteData := []byte(`byte`)

			ucDownloadfile := downloadfile.NewGoMockDownloadFile(ctrl)
			ucDownloadfile.EXPECT().DownloadByte(args.ctx, targetURL).Return(byteData, nil).AnyTimes()

			want := &cloudstorage.CloudFile{
				ObjectName:      args.req[0].GetFilename(),
				Size:            1024,
				ContentType:     "text/plain",
				StorageLocation: args.req[0].GetBucket() + "/" + args.req[0].GetFilename(),
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, bytes.NewReader(byteData), args.req[0].GetBucket(), args.req[0].GetFilename(), time.Time{}).Return(want, nil).AnyTimes()

			return test{
				args: args,
				fields: fields{
					uc: ucMock,
					df: ucDownloadfile,
				},
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload URL, When failed to upload, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			targetURL := "http://test.com/test.txt"

			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_URL,
						Filename: "object",
						Bucket:   "bucket",
						Detail: &rpc.UploadRequest_Url{
							Url: targetURL,
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			byteData := []byte(`byte`)

			ucDownloadfile := downloadfile.NewGoMockDownloadFile(ctrl)
			ucDownloadfile.EXPECT().DownloadByte(args.ctx, targetURL).Return(byteData, nil).AnyTimes()

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, bytes.NewReader(byteData), args.req[0].GetBucket(), args.req[0].GetFilename(), time.Time{}).Return(nil, defaultError).AnyTimes()

			return test{
				args: args,
				fields: fields{
					uc: ucMock,
					df: ucDownloadfile,
				},
				wantErr:    defaultError,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload URL, When failed to execute download byte, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			targetURL := "http://test.com/test.txt"

			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_URL,
						Filename: "object",
						Bucket:   "bucket",
						Detail: &rpc.UploadRequest_Url{
							Url: targetURL,
						},
					},
				},
				mock: defaultMock(ctx, 1),
			}

			ucDownloadfile := downloadfile.NewGoMockDownloadFile(ctrl)
			ucDownloadfile.EXPECT().DownloadByte(args.ctx, targetURL).Return(nil, defaultError).AnyTimes()

			return test{
				args: args,
				fields: fields{
					df: ucDownloadfile,
				},
				wantErr:    defaultError,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload, When it executed with EOF error, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:     ctx,
					req:     make(chan *rpc.UploadRequest, 1),
					resp:    make(chan *rpc.UploadResponse, 1),
					errRecv: io.EOF,
				},
			}

			want := &cloudstorage.CloudFile{
				ObjectName:      args.req[0].GetFilename(),
				Size:            1024,
				ContentType:     "image/jpeg",
				StorageLocation: args.req[0].GetBucket() + "/" + args.req[0].GetFilename(),
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, &bytes.Buffer{}, args.req[0].GetBucket(), args.req[0].GetFilename(), time.Time{}).Return(want, nil).AnyTimes()

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:       args,
				wantErr:    nil,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given invalid request of Upload, When it executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_UNSPECIFIED,
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:  ctx,
					req:  make(chan *rpc.UploadRequest, 1),
					resp: make(chan *rpc.UploadResponse, 1),
				},
			}

			return test{
				args:       args,
				wantErr:    fmt.Errorf("upload type %d is not supported", rpc.UploadType_UPLOAD_TYPE_FILE),
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:     ctx,
					req:     make(chan *rpc.UploadRequest, 1),
					resp:    make(chan *rpc.UploadResponse, 1),
					errRecv: io.EOF,
				},
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Upload(args.ctx, &bytes.Buffer{}, args.req[0].GetBucket(), args.req[0].GetFilename(), time.Time{}).Return(nil, defaultError).AnyTimes()

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:       args,
				wantErr:    defaultError,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
			}
		},
		"Given valid request of Upload, When it executed with an error, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			args := args{
				ctx: ctx,
				req: []*rpc.UploadRequest{
					{
						Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
						Filename: "object",
						Bucket:   "bucket",
					},
				},
				mock: &mockFileStorageService_UploadServer{
					ctx:     ctx,
					req:     make(chan *rpc.UploadRequest, 1),
					resp:    make(chan *rpc.UploadResponse, 1),
					errRecv: defaultError,
				},
			}

			return test{
				args:       args,
				wantErr:    defaultError,
				beforeFunc: defaultBeforeFunc,
				afterFunc:  defaultAfterFunc,
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
