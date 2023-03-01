package di

import (
	"log"
	"net/http"

	"github.com/moemoe89/file-storage/pkg/downloadfile"
)

// GetDownloadFile returns FileStorageUsecase instance.
func GetDownloadFile() downloadfile.DownloadFile {
	df, err := downloadfile.New(
		downloadfile.WithHTTPClient(http.DefaultClient),
	)
	if err != nil {
		log.Fatal("Minio Client", err)
	}

	return df
}
