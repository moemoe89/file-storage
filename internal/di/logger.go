package di

import (
	"log"

	"github.com/moemoe89/file-storage/pkg/logging"
)

// GetLogger get the logger wrapper.
func GetLogger() logging.Logger {
	l, err := logging.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	return l
}
