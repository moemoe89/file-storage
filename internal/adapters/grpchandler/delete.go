package grpchandler

import (
	"context"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
)

// Delete deletes the file by given bucket and object name.
func (h *fileStorageHandler) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.Empty, error) {
	ctx, span := h.trace.StartSpan(ctx, "Handler.Delete", nil)
	defer span.End()

	return new(rpc.Empty), h.minio.Delete(ctx, req.GetBucket(), req.GetObject())
}
