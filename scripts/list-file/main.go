package main

import (
	"context"
	"log"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Dial gRPC server connection.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// Create the client connection.
	client := rpc.NewFileStorageServiceClient(conn)

	// Bucket name
	bucket := "default"

	// Call Upload RPC.
	files, err := client.List(context.Background(), &rpc.ListRequest{
		Bucket: bucket,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files.GetFiles() {
		log.Println(file)
	}

	log.Print("### Successfully lists files ###\n")
}
