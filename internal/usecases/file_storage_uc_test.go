package usecases_test

import (
	"github.com/moemoe89/file-storage/internal/di"
	"github.com/moemoe89/file-storage/internal/usecases"
	"github.com/moemoe89/file-storage/pkg/cloudstorage"
)

type fields struct {
	minio cloudstorage.Client
}

func sut(f fields) usecases.FileStorageUsecase {
	return usecases.NewFileStorageUsecase(
		di.GetTracer().Tracer(),
		di.GetLogger(),
		f.minio,
	)
}
