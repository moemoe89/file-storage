package grpchandler_test

import (
	"context"
	"errors"
	"testing"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/internal/usecases"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFileStorageServer_List(t *testing.T) {
	type args struct {
		ctx    context.Context
		bucket string
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.ListResponse
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of List file, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
			}

			files := []string{"file1", "file2"}

			want := &rpc.ListResponse{
				Files: files,
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().List(args.ctx, args.bucket).Return(files, nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    want,
				wantErr: nil,
			}
		},
		"Given valid request of List file, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().List(args.ctx, args.bucket).Return(nil, errors.New("error"))

			return test{
				fields: fields{
					uc: ucMock,
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

			got, err := sut.List(tt.args.ctx, &rpc.ListRequest{
				Bucket: tt.args.bucket,
			})
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
