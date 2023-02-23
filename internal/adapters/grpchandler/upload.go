package grpchandler

import (
	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/pkg/grpchealth"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// FileStorageServiceServer is File Storage Service server contract.
type FileStorageServiceServer interface {
	rpc.FileStorageServiceServer
	health.HealthServer
}

// NewFileStorageHandler returns a new gRPC handler that implements FileStorageServiceServer interface.
func NewFileStorageHandler() FileStorageServiceServer {
	return &fileStorageHandler{}
}

// fileStorageHandler is a struct for handler.
type fileStorageHandler struct {
	rpc.UnimplementedFileStorageServiceServer
	grpchealth.HealthChecker
}

// Upload uploads file to storage both request and response as stream.
func (h *fileStorageHandler) Upload(stream rpc.FileStorageService_UploadServer) error {
	return nil
}
