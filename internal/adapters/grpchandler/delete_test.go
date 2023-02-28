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

func TestFileStorageServer_Delete(t *testing.T) {
	type args struct {
		ctx    context.Context
		bucket string
		object string
	}

	type test struct {
		fields  fields
		args    args
		want    *rpc.Empty
		wantErr error
	}

	tests := map[string]func(t *testing.T, ctrl *gomock.Controller) test{
		"Given valid request of Delete file, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
				object: "object",
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Delete(args.ctx, args.bucket, args.object).Return(nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    new(rpc.Empty),
				wantErr: nil,
			}
		},
		"Given valid request of Delete file with empty bucket, When UC executed successfully, Return no error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				object: "object",
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Delete(args.ctx, "default", args.object).Return(nil)

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    new(rpc.Empty),
				wantErr: nil,
			}
		},
		"Given valid request of Delete file, When UC failed to executed, Return error": func(t *testing.T, ctrl *gomock.Controller) test {
			ctx := context.Background()

			args := args{
				ctx:    ctx,
				bucket: "bucket",
				object: "object",
			}

			ucMock := usecases.NewGoMockFileStorageUsecase(ctrl)
			ucMock.EXPECT().Delete(args.ctx, args.bucket, args.object).Return(errors.New("error"))

			return test{
				fields: fields{
					uc: ucMock,
				},
				args:    args,
				want:    new(rpc.Empty),
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

			got, err := sut.Delete(tt.args.ctx, &rpc.DeleteRequest{
				Object: tt.args.object,
				Bucket: tt.args.bucket,
			})
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
