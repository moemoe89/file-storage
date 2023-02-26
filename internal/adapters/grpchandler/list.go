package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
)

// List lists the files by given bucket name.
func (h *fileStorageHandler) List(ctx context.Context, req *rpc.ListRequest) (*rpc.ListResponse, error) {
	bucket := req.GetBucket()
	if bucket == "" {
		bucket = defaultBucket
	}

	objects, err := h.uc.List(ctx, bucket)
	if err != nil {
		return nil, err
	}

	return &rpc.ListResponse{
		Files: objects,
	}, nil
}
