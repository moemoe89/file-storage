package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

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

	// The target file path, call from root directory project.
	filepath := "scripts/upload-from-file/wp-logo.jpeg"
	// filename can be structured with directory e.g. images/wp-logo.jpeg
	// and in the cloud storage can create the directory directly.
	filename := "wp-logo.jpeg"

	// Open the target upload file.
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Close file after finished uploading.
	defer func() { _ = file.Close() }()

	// The offset for tracking the stream.
	var offset int64

	// Set the buffer for each stream.
	buf := make([]byte, 1024)

	go func() {
		for {
			// Read file as chunk.
			chunk, err := file.Read(buf)

			// Check whether the stream is finish or not.
			if errors.Is(err, io.EOF) {
				err = stream.CloseSend()
				if err != nil {
					log.Fatal(err)
				}

				// Stop the stream.
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			// Send stream to server.
			err = stream.Send(&rpc.UploadRequest{
				Type:     rpc.UploadType_UPLOAD_TYPE_FILE,
				Filename: filename,
				Detail: &rpc.UploadRequest_File{
					File: &rpc.FileUpload{
						Data:   buf[:chunk],
						Offset: offset,
					},
				},
				Validation: &rpc.Validation{
					// Validates what Content Type we expected.
					ContentTypes: []string{"image/jpeg", "image/png"},
					// Validates the max size we expected, 0 means no limit.
					MaxSize: 0,
				},
			})

			if err != nil {
				log.Fatal(err)
			}

			// Update the offset.
			offset += int64(chunk)
		}
	}()

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
