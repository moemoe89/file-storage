package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	iDI "github.com/moemoe89/file-storage/internal/di"
	"github.com/moemoe89/file-storage/pkg/di"

	"go.uber.org/zap"
)

func main() {
	logger := iDI.GetLogger()

	server := iDI.GetFileStorageGRPCServer()
	gateway := iDI.GetFileStorageGatewayServer()

	logger.Info("File Storage service is ready")

	go func() {
		// Run() keeps its process until receiving any error
		if err := server.Run(); err != nil {
			logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	go func() {
		// Run() keeps its process until receiving any error
		if err := gateway.Run(); err != nil {
			logger.Fatal("failed to serve Gateway", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	logger.Info(fmt.Sprintf("SIGNAL %d received, shutting down gracefully...", <-quit))
	di.CloseAll()

	logger.Info("finished graceful shut down")
}
