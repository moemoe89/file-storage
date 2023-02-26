package grpchandler_test

import (
	"github.com/moemoe89/file-storage/internal/adapters/grpchandler"
	"github.com/moemoe89/file-storage/internal/usecases"
)

type fields struct {
	uc usecases.FileStorageUsecase
}

func sut(f fields) grpchandler.FileStorageServiceServer {
	return grpchandler.NewFileStorageHandler(
		f.uc,
	)
}
