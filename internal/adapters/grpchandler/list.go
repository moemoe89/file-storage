package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
)

// List lists the files by given bucket name.
func (h *fileStorageHandler) List(ctx context.Context, req *rpc.ListRequest) (*rpc.ListResponse, error) {
	ctx, span := h.trace.StartSpan(ctx, "Handler.List", nil)
	defer span.End()

	objects, err := h.minio.ListObjects(ctx, req.GetBucket())
	if err != nil {
		return nil, err
	}

	return &rpc.ListResponse{
		Files: objects,
	}, nil
}
