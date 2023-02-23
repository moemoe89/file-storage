package di

import (
	"log"
	"os"
	"sync"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/pkg/di"
	"github.com/moemoe89/file-storage/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcServerOnce    sync.Once
	fileStorageServer server.Server
)

// GetFileStorageGRPCServer returns gRPC server instance for File Storage Service server.
func GetFileStorageGRPCServer() server.Server {
	return getGRPCServer(fileStorageServer, func(server *grpc.Server) {
		h := GetFileStorageGRPCHandler()
		rpc.RegisterFileStorageServiceServer(server, h)
		grpc_health_v1.RegisterHealthServer(server, h)
	})
}

// getGRPCServer
func getGRPCServer(grpcServer server.Server, register server.HandlerRegister) server.Server {
	grpcServerOnce.Do(func() {
		opts := GetMiddleware()

		port := os.Getenv("SERVER_PORT")

		s, err := server.NewGRPCServer(port, register, opts...)
		if err != nil {
			log.Fatal("gRPC server", err)
		}

		di.RegisterCloser("gRPC server", di.NewCloser(s.GracefulStop))

		grpcServer = s
	})

	return grpcServer
}
