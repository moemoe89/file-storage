package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
)

// Delete deletes the file by given bucket and object name.
func (h *fileStorageHandler) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.Empty, error) {
	bucket := req.GetBucket()
	if bucket == "" {
		bucket = defaultBucket
	}

	return new(rpc.Empty), h.uc.Delete(ctx, bucket, req.GetObject())
}
