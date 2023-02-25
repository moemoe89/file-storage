package di

import "github.com/moemoe89/file-storage/internal/usecases"

// GetFileStorageUsecase returns FileStorageUsecase instance.
func GetFileStorageUsecase() usecases.FileStorageUsecase {
	return usecases.NewFileStorageUsecase(
		GetTracer().Tracer(),
		GetLogger(),
		GetMinio(),
	)
}
