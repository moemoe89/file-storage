package grpchandler

import (
	"bytes"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/internal/usecases"
	"github.com/moemoe89/file-storage/pkg/downloadfile"
	"github.com/moemoe89/file-storage/pkg/grpchealth"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// FileStorageServiceServer is File Storage Service server contract.
type FileStorageServiceServer interface {
	rpc.FileStorageServiceServer
	health.HealthServer
}

// NewFileStorageHandler returns a new gRPC handler that implements FileStorageServiceServer interface.
func NewFileStorageHandler(uc usecases.FileStorageUsecase, df downloadfile.DownloadFile) FileStorageServiceServer {
	return &fileStorageHandler{uc: uc, df: df}
}

// fileStorageHandler is a struct for handler.
type fileStorageHandler struct {
	rpc.UnimplementedFileStorageServiceServer
	grpchealth.HealthChecker
	uc usecases.FileStorageUsecase
	df downloadfile.DownloadFile
}

type fileData struct {
	// id is a file ID will send to client.
	id string
	// bucket is a bucket name sent from client.
	bucket string
	// filename is a file name sent from client.
	filename string
	// objectName is a name for the uploaded file.
	objectName string
	// contentType is a Content-Type for the uploaded file.
	contentType string
	// targetSize is an expected max size from client.
	targetSize int64
	// offset is an accumulated chunk sent from client.
	offset int64
	// size is a size for uploaded file.
	size int64
	// firstChunk is a flag to know the stream is first chunk or not.
	firstChunk bool
	// buf is a variable-sized buffer for the uploaded file.
	buf *bytes.Buffer
}

const (
	// defaultBucket is a default bucket to use when it's not defined on the request.
	defaultBucket = "default"
)
