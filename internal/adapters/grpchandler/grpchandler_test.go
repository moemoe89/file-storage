package grpchandler_test

import (
	"github.com/moemoe89/file-storage/internal/adapters/grpchandler"
	"github.com/moemoe89/file-storage/internal/usecases"
	"github.com/moemoe89/file-storage/pkg/downloadfile"
)

type fields struct {
	uc usecases.FileStorageUsecase
	df downloadfile.DownloadFile
}

func sut(f fields) grpchandler.FileStorageServiceServer {
	return grpchandler.NewFileStorageHandler(
		f.uc,
		f.df,
	)
}
