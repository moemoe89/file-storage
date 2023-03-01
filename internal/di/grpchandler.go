package di

import (
	"github.com/moemoe89/file-storage/internal/adapters/grpchandler"
)

// GetFileStorageGRPCHandler returns FileStorageServiceServer handler.
func GetFileStorageGRPCHandler() grpchandler.FileStorageServiceServer {
	return grpchandler.NewFileStorageHandler(GetFileStorageUsecase(), GetDownloadFile())
}
