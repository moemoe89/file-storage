// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileStorageServiceClient is the client API for FileStorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileStorageServiceClient interface {
	// Upload uploads file to storage both request and response as stream.
	Upload(ctx context.Context, opts ...grpc.CallOption) (FileStorageService_UploadClient, error)
	// List lists the files by given bucket name.
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	// Delete deletes the file by given bucket and object name.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Empty, error)
}

type fileStorageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileStorageServiceClient(cc grpc.ClientConnInterface) FileStorageServiceClient {
	return &fileStorageServiceClient{cc}
}

func (c *fileStorageServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FileStorageService_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileStorageService_ServiceDesc.Streams[0], "/FileStorageService/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileStorageServiceUploadClient{stream}
	return x, nil
}

type FileStorageService_UploadClient interface {
	Send(*UploadRequest) error
	Recv() (*UploadResponse, error)
	grpc.ClientStream
}

type fileStorageServiceUploadClient struct {
	grpc.ClientStream
}

func (x *fileStorageServiceUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileStorageServiceUploadClient) Recv() (*UploadResponse, error) {
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileStorageServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/FileStorageService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileStorageServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/FileStorageService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileStorageServiceServer is the server API for FileStorageService service.
// All implementations must embed UnimplementedFileStorageServiceServer
// for forward compatibility
type FileStorageServiceServer interface {
	// Upload uploads file to storage both request and response as stream.
	Upload(FileStorageService_UploadServer) error
	// List lists the files by given bucket name.
	List(context.Context, *ListRequest) (*ListResponse, error)
	// Delete deletes the file by given bucket and object name.
	Delete(context.Context, *DeleteRequest) (*Empty, error)
	mustEmbedUnimplementedFileStorageServiceServer()
}

// UnimplementedFileStorageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileStorageServiceServer struct {
}

func (UnimplementedFileStorageServiceServer) Upload(FileStorageService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFileStorageServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedFileStorageServiceServer) Delete(context.Context, *DeleteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedFileStorageServiceServer) mustEmbedUnimplementedFileStorageServiceServer() {}

// UnsafeFileStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileStorageServiceServer will
// result in compilation errors.
type UnsafeFileStorageServiceServer interface {
	mustEmbedUnimplementedFileStorageServiceServer()
}

func RegisterFileStorageServiceServer(s grpc.ServiceRegistrar, srv FileStorageServiceServer) {
	s.RegisterService(&FileStorageService_ServiceDesc, srv)
}

func _FileStorageService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileStorageServiceServer).Upload(&fileStorageServiceUploadServer{stream})
}

type FileStorageService_UploadServer interface {
	Send(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type fileStorageServiceUploadServer struct {
	grpc.ServerStream
}

func (x *fileStorageServiceUploadServer) Send(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileStorageServiceUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileStorageService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileStorageServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FileStorageService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileStorageServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileStorageService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileStorageServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FileStorageService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileStorageServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileStorageService_ServiceDesc is the grpc.ServiceDesc for FileStorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileStorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileStorageService",
	HandlerType: (*FileStorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _FileStorageService_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _FileStorageService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _FileStorageService_Upload_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service.proto",
}
