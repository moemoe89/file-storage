package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
)

// List lists the files by given bucket name.
func (h *fileStorageHandler) List(ctx context.Context, req *rpc.ListRequest) (*rpc.ListResponse, error) {
	objects, err := h.uc.List(ctx, req.GetBucket())
	if err != nil {
		return nil, err
	}

	return &rpc.ListResponse{
		Files: objects,
	}, nil
}
