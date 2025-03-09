// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: mapreduce.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MasterService_RegisterWorker_FullMethodName = "/mapreduce.MasterService/RegisterWorker"
	MasterService_GetTask_FullMethodName        = "/mapreduce.MasterService/GetTask"
	MasterService_ReportTask_FullMethodName     = "/mapreduce.MasterService/ReportTask"
)

// MasterServiceClient is the client API for MasterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterServiceClient interface {
	RegisterWorker(ctx context.Context, in *RegisterWorkerRequest, opts ...grpc.CallOption) (*RegisterWorkerResponse, error)
	GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error)
	ReportTask(ctx context.Context, in *ReportTaskRequest, opts ...grpc.CallOption) (*ReportTaskResponse, error)
}

type masterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterServiceClient(cc grpc.ClientConnInterface) MasterServiceClient {
	return &masterServiceClient{cc}
}

func (c *masterServiceClient) RegisterWorker(ctx context.Context, in *RegisterWorkerRequest, opts ...grpc.CallOption) (*RegisterWorkerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterWorkerResponse)
	err := c.cc.Invoke(ctx, MasterService_RegisterWorker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterServiceClient) GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTaskResponse)
	err := c.cc.Invoke(ctx, MasterService_GetTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterServiceClient) ReportTask(ctx context.Context, in *ReportTaskRequest, opts ...grpc.CallOption) (*ReportTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportTaskResponse)
	err := c.cc.Invoke(ctx, MasterService_ReportTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServiceServer is the server API for MasterService service.
// All implementations must embed UnimplementedMasterServiceServer
// for forward compatibility.
type MasterServiceServer interface {
	RegisterWorker(context.Context, *RegisterWorkerRequest) (*RegisterWorkerResponse, error)
	GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error)
	ReportTask(context.Context, *ReportTaskRequest) (*ReportTaskResponse, error)
	mustEmbedUnimplementedMasterServiceServer()
}

// UnimplementedMasterServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMasterServiceServer struct{}

func (UnimplementedMasterServiceServer) RegisterWorker(context.Context, *RegisterWorkerRequest) (*RegisterWorkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterWorker not implemented")
}
func (UnimplementedMasterServiceServer) GetTask(context.Context, *GetTaskRequest) (*GetTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}
func (UnimplementedMasterServiceServer) ReportTask(context.Context, *ReportTaskRequest) (*ReportTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportTask not implemented")
}
func (UnimplementedMasterServiceServer) mustEmbedUnimplementedMasterServiceServer() {}
func (UnimplementedMasterServiceServer) testEmbeddedByValue()                       {}

// UnsafeMasterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterServiceServer will
// result in compilation errors.
type UnsafeMasterServiceServer interface {
	mustEmbedUnimplementedMasterServiceServer()
}

func RegisterMasterServiceServer(s grpc.ServiceRegistrar, srv MasterServiceServer) {
	// If the following call pancis, it indicates UnimplementedMasterServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MasterService_ServiceDesc, srv)
}

func _MasterService_RegisterWorker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterWorkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServiceServer).RegisterWorker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterService_RegisterWorker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServiceServer).RegisterWorker(ctx, req.(*RegisterWorkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterService_GetTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServiceServer).GetTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterService_GetTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServiceServer).GetTask(ctx, req.(*GetTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterService_ReportTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServiceServer).ReportTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MasterService_ReportTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServiceServer).ReportTask(ctx, req.(*ReportTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MasterService_ServiceDesc is the grpc.ServiceDesc for MasterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mapreduce.MasterService",
	HandlerType: (*MasterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterWorker",
			Handler:    _MasterService_RegisterWorker_Handler,
		},
		{
			MethodName: "GetTask",
			Handler:    _MasterService_GetTask_Handler,
		},
		{
			MethodName: "ReportTask",
			Handler:    _MasterService_ReportTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mapreduce.proto",
}
