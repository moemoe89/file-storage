package main

import (
	"context"
	"errors"
	"io"
	"log"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
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

	// Disable retry since this call using stream.
	callOpts := []grpc.CallOption{
		grpc_retry.Disable(),
	}

	// Call Upload RPC.
	stream, err := client.Upload(context.Background(), callOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// filename can be structured with directory e.g. images/logo.jpeg
	// and in the cloud storage can create the directory directly.
	filename := "logo.jpeg"

	// Send stream to server.
	err = stream.Send(&rpc.UploadRequest{
		Type:     rpc.UploadType_UPLOAD_TYPE_URL,
		Filename: filename,
		Validation: &rpc.Validation{
			// Validates what Content-Type we expected.
			ContentTypes: []string{"image/jpeg", "image/png"},
			// Validates the max size we expected, 0 means no limit.
			MaxSize: 0,
		},
		Detail: &rpc.UploadRequest_Url{
			Url: "https://images.ctfassets.net/30xxrlh9suih/6jyUVNJs0bKkYwVGAxarUa/1fa3440d4c578d6ce05d38dbbf9fcb08/wovenplanet_60px_onwhite.jpeg",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	objectName := ""

	// Channel when done for both get the stream.
	done := make(chan bool)

	go func() {
		for {
			// Reads the chunk of stream.
			chunk, err := stream.Recv()

			// Check whether the stream is finish or not.
			if errors.Is(err, io.EOF) {
				// Close the channel to stop the service.
				close(done)

				// Stop the stream.
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			// Get the Object name.
			objectName = chunk.GetObjectName()
		}
	}()

	<-done

	log.Printf("### Uploading finished!! Uploaded to: %s ###\n", objectName)
}
