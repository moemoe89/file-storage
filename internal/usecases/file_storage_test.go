package usecases_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/moemoe89/file-storage/pkg/cloudstorage"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFileStorageUC_Upload(t *testing.T) {
	type args struct {
		ctx     context.Context
		file    io.Reader
		bucket  string
		object  string
		expires time.Time
	}

	type test struct {
		fields  fields
		args    args
		want    *cloudstorage.CloudFile
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Upload file, When MinIO executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:     ctx,
				file:    strings.NewReader("my request"),
				bucket:  "bucket",
				object:  "object",
				expires: time.Time{},
			}

			want := &cloudstorage.CloudFile{
				ObjectName:      "object",
				Size:            1024,
				ContentType:     "text/plain",
				StorageLocation: "default/object",
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().Upload(args.ctx, args.file, args.bucket, args.object, args.expires).Return(want, nil)

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of Upload file, When MinIO failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:     ctx,
				file:    strings.NewReader("my request"),
				bucket:  "bucket",
				object:  "object",
				expires: time.Time{},
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().Upload(args.ctx, args.file, args.bucket, args.object, args.expires).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.Upload(tt.args.ctx, tt.args.file, tt.args.bucket, tt.args.object, tt.args.expires)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestFileStorageUC_List(t *testing.T) {
	type args struct {
		ctx    context.Context
		bucket string
	}

	type test struct {
		fields  fields
		args    args
		want    []string
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of List file, When MinIO executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
			}

			want := []string{
				"file-1",
				"file-2",
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().ListObjects(args.ctx, args.bucket).Return(want, nil)

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List file, When MinIO failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().ListObjects(args.ctx, args.bucket).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			got, err := sut.List(tt.args.ctx, tt.args.bucket)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestFileStorageUC_Delete(t *testing.T) {
	type args struct {
		ctx    context.Context
		bucket string
		object string
	}

	type test struct {
		fields  fields
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Delete file, When MinIO executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
				object: "object",
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().Delete(args.ctx, args.bucket, args.object).Return(nil)

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				wantErr: nil,
			}
		},
		"Given valid request of Delete file, When MinIO failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
				object: "object",
			}

			minioMock := cloudstorage.NewGoMockClient(ctrl)
			minioMock.EXPECT().Delete(args.ctx, args.bucket, args.object).Return(errors.New("error"))

			return test{
				fields: fields{
					minio: minioMock,
				},
				args:    args,
				wantErr: errors.New("error"),
			}
		},
	}

	for name, testFn := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt := testFn(t, ctrl)

			sut := sut(tt.fields)

			err := sut.Delete(tt.args.ctx, tt.args.bucket, tt.args.object)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
