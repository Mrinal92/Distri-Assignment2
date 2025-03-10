// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: stripe.proto

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
	PaymentGateway_Register_FullMethodName       = "/stripe.PaymentGateway/Register"
	PaymentGateway_Authenticate_FullMethodName   = "/stripe.PaymentGateway/Authenticate"
	PaymentGateway_ProcessPayment_FullMethodName = "/stripe.PaymentGateway/ProcessPayment"
)

// PaymentGatewayClient is the client API for PaymentGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ---------------------
// Service definitions
// ---------------------
type PaymentGatewayClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Authenticate(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error)
	ProcessPayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
}

type paymentGatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentGatewayClient(cc grpc.ClientConnInterface) PaymentGatewayClient {
	return &paymentGatewayClient{cc}
}

func (c *paymentGatewayClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, PaymentGateway_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) Authenticate(ctx context.Context, in *AuthRequest, opts ...grpc.CallOption) (*AuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResponse)
	err := c.cc.Invoke(ctx, PaymentGateway_Authenticate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentGatewayClient) ProcessPayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, PaymentGateway_ProcessPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentGatewayServer is the server API for PaymentGateway service.
// All implementations must embed UnimplementedPaymentGatewayServer
// for forward compatibility.
//
// ---------------------
// Service definitions
// ---------------------
type PaymentGatewayServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Authenticate(context.Context, *AuthRequest) (*AuthResponse, error)
	ProcessPayment(context.Context, *PaymentRequest) (*PaymentResponse, error)
	mustEmbedUnimplementedPaymentGatewayServer()
}

// UnimplementedPaymentGatewayServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPaymentGatewayServer struct{}

func (UnimplementedPaymentGatewayServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedPaymentGatewayServer) Authenticate(context.Context, *AuthRequest) (*AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedPaymentGatewayServer) ProcessPayment(context.Context, *PaymentRequest) (*PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessPayment not implemented")
}
func (UnimplementedPaymentGatewayServer) mustEmbedUnimplementedPaymentGatewayServer() {}
func (UnimplementedPaymentGatewayServer) testEmbeddedByValue()                        {}

// UnsafePaymentGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentGatewayServer will
// result in compilation errors.
type UnsafePaymentGatewayServer interface {
	mustEmbedUnimplementedPaymentGatewayServer()
}

func RegisterPaymentGatewayServer(s grpc.ServiceRegistrar, srv PaymentGatewayServer) {
	// If the following call pancis, it indicates UnimplementedPaymentGatewayServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PaymentGateway_ServiceDesc, srv)
}

func _PaymentGateway_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentGateway_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentGateway_Authenticate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).Authenticate(ctx, req.(*AuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentGateway_ProcessPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentGatewayServer).ProcessPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentGateway_ProcessPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentGatewayServer).ProcessPayment(ctx, req.(*PaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentGateway_ServiceDesc is the grpc.ServiceDesc for PaymentGateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentGateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stripe.PaymentGateway",
	HandlerType: (*PaymentGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _PaymentGateway_Register_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _PaymentGateway_Authenticate_Handler,
		},
		{
			MethodName: "ProcessPayment",
			Handler:    _PaymentGateway_ProcessPayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stripe.proto",
}

const (
	Bank_PreparePayment_FullMethodName  = "/stripe.Bank/PreparePayment"
	Bank_CommitPayment_FullMethodName   = "/stripe.Bank/CommitPayment"
	Bank_RollbackPayment_FullMethodName = "/stripe.Bank/RollbackPayment"
)

// BankClient is the client API for Bank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BankClient interface {
	PreparePayment(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error)
	CommitPayment(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error)
	RollbackPayment(ctx context.Context, in *RollbackRequest, opts ...grpc.CallOption) (*RollbackResponse, error)
}

type bankClient struct {
	cc grpc.ClientConnInterface
}

func NewBankClient(cc grpc.ClientConnInterface) BankClient {
	return &bankClient{cc}
}

func (c *bankClient) PreparePayment(ctx context.Context, in *PrepareRequest, opts ...grpc.CallOption) (*PrepareResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrepareResponse)
	err := c.cc.Invoke(ctx, Bank_PreparePayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankClient) CommitPayment(ctx context.Context, in *CommitRequest, opts ...grpc.CallOption) (*CommitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommitResponse)
	err := c.cc.Invoke(ctx, Bank_CommitPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankClient) RollbackPayment(ctx context.Context, in *RollbackRequest, opts ...grpc.CallOption) (*RollbackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RollbackResponse)
	err := c.cc.Invoke(ctx, Bank_RollbackPayment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BankServer is the server API for Bank service.
// All implementations must embed UnimplementedBankServer
// for forward compatibility.
type BankServer interface {
	PreparePayment(context.Context, *PrepareRequest) (*PrepareResponse, error)
	CommitPayment(context.Context, *CommitRequest) (*CommitResponse, error)
	RollbackPayment(context.Context, *RollbackRequest) (*RollbackResponse, error)
	mustEmbedUnimplementedBankServer()
}

// UnimplementedBankServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBankServer struct{}

func (UnimplementedBankServer) PreparePayment(context.Context, *PrepareRequest) (*PrepareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreparePayment not implemented")
}
func (UnimplementedBankServer) CommitPayment(context.Context, *CommitRequest) (*CommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitPayment not implemented")
}
func (UnimplementedBankServer) RollbackPayment(context.Context, *RollbackRequest) (*RollbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RollbackPayment not implemented")
}
func (UnimplementedBankServer) mustEmbedUnimplementedBankServer() {}
func (UnimplementedBankServer) testEmbeddedByValue()              {}

// UnsafeBankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BankServer will
// result in compilation errors.
type UnsafeBankServer interface {
	mustEmbedUnimplementedBankServer()
}

func RegisterBankServer(s grpc.ServiceRegistrar, srv BankServer) {
	// If the following call pancis, it indicates UnimplementedBankServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Bank_ServiceDesc, srv)
}

func _Bank_PreparePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrepareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServer).PreparePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bank_PreparePayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServer).PreparePayment(ctx, req.(*PrepareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bank_CommitPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServer).CommitPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bank_CommitPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServer).CommitPayment(ctx, req.(*CommitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bank_RollbackPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RollbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankServer).RollbackPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bank_RollbackPayment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankServer).RollbackPayment(ctx, req.(*RollbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bank_ServiceDesc is the grpc.ServiceDesc for Bank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stripe.Bank",
	HandlerType: (*BankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PreparePayment",
			Handler:    _Bank_PreparePayment_Handler,
		},
		{
			MethodName: "CommitPayment",
			Handler:    _Bank_CommitPayment_Handler,
		},
		{
			MethodName: "RollbackPayment",
			Handler:    _Bank_RollbackPayment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stripe.proto",
}
