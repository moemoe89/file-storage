package client

import (
	"context"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"path"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientConn() rpc.FileStorageServiceClient {
	// Dial gRPC server connection.
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	// Create the client connection.
	client := rpc.NewFileStorageServiceClient(conn)

	return client
}

func UploadFromFile(filepath, filename string, contentTypes []string, maxSize int64) { //nolint:funlen
	client := clientConn()

	// Disable retry since this call using stream.
	callOpts := []grpc.CallOption{
		grpc_retry.Disable(),
	}

	// Call Upload RPC.
	stream, err := client.Upload(context.Background(), callOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// Open the target upload file.
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// If file name not specified, get from file path.
	if filename == "" {
		filename = path.Base(filepath)
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
					ContentTypes: contentTypes,
					// Validates the max size we expected, 0 means no limit.
					MaxSize: maxSize,
				},
			})

			if err != nil {
				log.Fatal("asdas ", err)
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

func UploadFromURL(fileURL, filename string, contentTypes []string, maxSize int64) { //nolint:funlen
	client := clientConn()

	// Disable retry since this call using stream.
	callOpts := []grpc.CallOption{
		grpc_retry.Disable(),
	}

	// Call Upload RPC.
	stream, err := client.Upload(context.Background(), callOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// If file name not specified, get from URL.
	if filename == "" {
		u, err := url.Parse(fileURL)
		if err != nil {
			log.Fatal(err)
		}

		filename = path.Base(u.Path)
	}

	// Send stream to server.
	err = stream.Send(&rpc.UploadRequest{
		Type:     rpc.UploadType_UPLOAD_TYPE_URL,
		Filename: filename,
		Validation: &rpc.Validation{
			// Validates what Content-Type we expected.
			ContentTypes: contentTypes,
			// Validates the max size we expected, 0 means no limit.
			MaxSize: maxSize,
		},
		Detail: &rpc.UploadRequest_Url{
			Url: fileURL,
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

func ListFile(bucket string) {
	client := clientConn()

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

func DeleteFile(bucket, object string) {
	client := clientConn()

	// Call Upload RPC.
	_, err := client.Delete(context.Background(), &rpc.DeleteRequest{
		Object: object,
		Bucket: bucket,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("### Successfully deleted file: %s ###\n", object)
}
