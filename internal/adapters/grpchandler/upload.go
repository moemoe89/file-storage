package grpchandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	rpc "github.com/moemoe89/file-storage/api/go/grpc"
	"github.com/moemoe89/file-storage/pkg/cloudstorage"
	"github.com/moemoe89/file-storage/pkg/diskstorage"
	"github.com/moemoe89/file-storage/pkg/grpchealth"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// FileStorageServiceServer is File Storage Service server contract.
type FileStorageServiceServer interface {
	rpc.FileStorageServiceServer
	health.HealthServer
}

// NewFileStorageHandler returns a new gRPC handler that implements FileStorageServiceServer interface.
func NewFileStorageHandler(minio cloudstorage.Client) FileStorageServiceServer {
	return &fileStorageHandler{minio: minio}
}

// fileStorageHandler is a struct for handler.
type fileStorageHandler struct {
	rpc.UnimplementedFileStorageServiceServer
	grpchealth.HealthChecker
	minio cloudstorage.Client
}

type fileData struct {
	// id is a file ID will send to client.
	id string
	// bucket is a bucket name sent from client.
	bucket string
	// filename is a file name sent from client.
	filename string
	// objectName is a name for the uploaded file.
	objectName string
	// contentType is a Content-Type for the uploaded file.
	contentType string
	// targetSize is an expected max size from client.
	targetSize int64
	// offset is an accumulated chunk sent from client.
	offset int64
	// size is a size for uploaded file.
	size int64
	// firstChunk is a flag to know the stream is first chunk or not.
	firstChunk bool
}

const (
	// defaultBucket is a default bucket to use when it's not defined on the request.
	defaultBucket = "default"
)

// Upload uploads file to storage both request and response as stream.
func (h *fileStorageHandler) Upload(stream rpc.FileStorageService_UploadServer) error { //nolint:funlen
	// Initialize diskstorage package.
	f, err := diskstorage.New()
	if err != nil {
		return err
	}

	fd := &fileData{id: uuid.New().String(), firstChunk: true}

	// Initialize context from stream Context.
	var ctx = stream.Context()

	for {
		// Handle if the Context is Done.
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}

			return ctx.Err()
		default:
		}

		// Reads the chunk stream data.
		chunk, err := stream.Recv()
		if fd.firstChunk {
			assignFileData(chunk, fd)

			switch chunk.GetType() {
			case rpc.UploadType_UPLOAD_TYPE_URL:
				return h.uploadFromURL(ctx, stream, chunk, f, fd)
			case rpc.UploadType_UPLOAD_TYPE_FILE:
				fd.contentType = mimetype.Detect(chunk.GetFile().GetData()).String()

				// Checks the expected Content-Type.
				err := validateContentType(chunk.GetValidation().GetContentTypes(), fd.contentType)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("upload type %d is not supported", rpc.UploadType_UPLOAD_TYPE_FILE)
			}
		}

		// Handle the end of file stream.
		if errors.Is(err, io.EOF) {
			return h.endStream(ctx, stream, f, fd)
		}

		if err != nil {
			return err
		}

		// Checks the expected max size.
		err = validateSize(chunk.GetFile().GetOffset(), fd.offset, fd.targetSize)
		if err != nil {
			return err
		}

		// Accumulated chunk.
		fd.offset += int64(len(chunk.GetFile().GetData()))

		// Writes the chunk data to buffer.
		err = f.Write(chunk.GetFile().GetData())
		if err != nil {
			return err
		}

		// File not uploaded yet, so leave the URL empty.
		err = sendSteam(stream, fd)
		if err != nil {
			return err
		}
	}
}

// assignFiledata assign furst chunk, file name, bucket and target size from request to FileData object.
func assignFileData(chunk *rpc.UploadRequest, fd *fileData) {
	fd.firstChunk = false

	fd.filename = chunk.GetFilename()

	if fd.filename == "" {
		fd.filename = uuid.New().String()
	}

	fd.bucket = chunk.GetBucket()

	if fd.bucket == "" {
		fd.bucket = defaultBucket
	}

	fd.targetSize = chunk.GetValidation().GetMaxSize()
}

func (h *fileStorageHandler) uploadFromURL(
	ctx context.Context,
	stream rpc.FileStorageService_UploadServer,
	chunk *rpc.UploadRequest,
	f diskstorage.DiskStorage,
	fd *fileData,
) error {
	fileUpload, err := f.DownloadByte(ctx, chunk.GetUrl())
	if err != nil {
		return err
	}

	// assign size of the file.
	fd.size = int64(len(fileUpload))

	// Uploads to Google Cloud Storage.
	cloudFile, err := h.minio.Upload(ctx, bytes.NewReader(fileUpload), fd.bucket, fd.filename, time.Time{})
	if err != nil {
		return err
	}

	fd.objectName = cloudFile.ObjectName

	// Send stream to the client after getting uploaded URL.
	return sendSteam(stream, fd)
}

// endStream ends the stream to client,
// including do compression (if requested by client)
// and uploading to cloud storage.
func (h *fileStorageHandler) endStream(
	ctx context.Context,
	stream rpc.FileStorageService_UploadServer,
	f diskstorage.DiskStorage,
	fd *fileData,
) error {
	// Reset the buffer.
	defer f.ResetBuffer()

	fileUpload := f.GetBuffer()

	fd.size = fd.offset

	// Uploads to Google Cloud Storage.
	cloudFile, err := h.minio.Upload(ctx, fileUpload, fd.bucket, fd.filename, time.Time{})
	if err != nil {
		return err
	}

	fd.objectName = cloudFile.ObjectName

	// NOTE: Writes file to disk.
	// Example if we want to write the file o disk.
	// return f.WriteFile(fmt.Sprintf(tmpUploadPath, filename), os.FileMode(0644))

	// Send stream to the client after getting uploaded URL.
	return sendSteam(stream, fd)
}

// sendSteam sends stream to client.
func sendSteam(
	stream rpc.FileStorageService_UploadServer, fd *fileData,
) error {
	return stream.Send(&rpc.UploadResponse{
		Id:              fd.id,
		ObjectName:      fd.objectName,
		StorageLocation: "",
		Offset:          fd.offset,
		Size:            fd.size,
		ContentType:     fd.contentType,
	})
}
