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

	// Object name
	object := "logo.jpeg"

	// Call Upload RPC.
	_, err = client.Delete(context.Background(), &rpc.DeleteRequest{
		Object: object,
		Bucket: bucket,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("### Successfully deleted file: %s ###\n", object)
}
