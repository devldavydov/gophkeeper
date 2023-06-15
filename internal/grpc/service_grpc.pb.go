// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: internal/grpc/proto/service.proto

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

const (
	GophKeeperService_UserCreate_FullMethodName    = "/proto.GophKeeperService/UserCreate"
	GophKeeperService_UserLogin_FullMethodName     = "/proto.GophKeeperService/UserLogin"
	GophKeeperService_SecretGetList_FullMethodName = "/proto.GophKeeperService/SecretGetList"
	GophKeeperService_SecretGet_FullMethodName     = "/proto.GophKeeperService/SecretGet"
	GophKeeperService_SecretSet_FullMethodName     = "/proto.GophKeeperService/SecretSet"
	GophKeeperService_SecretDel_FullMethodName     = "/proto.GophKeeperService/SecretDel"
)

// GophKeeperServiceClient is the client API for GophKeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperServiceClient interface {
	// User
	UserCreate(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserAuthToken, error)
	UserLogin(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserAuthToken, error)
	// Secret
	SecretGetList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SecretList, error)
	SecretGet(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Secret, error)
	SecretSet(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Empty, error)
	SecretDel(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Empty, error)
}

type gophKeeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperServiceClient(cc grpc.ClientConnInterface) GophKeeperServiceClient {
	return &gophKeeperServiceClient{cc}
}

func (c *gophKeeperServiceClient) UserCreate(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserAuthToken, error) {
	out := new(UserAuthToken)
	err := c.cc.Invoke(ctx, GophKeeperService_UserCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) UserLogin(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserAuthToken, error) {
	out := new(UserAuthToken)
	err := c.cc.Invoke(ctx, GophKeeperService_UserLogin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SecretGetList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*SecretList, error) {
	out := new(SecretList)
	err := c.cc.Invoke(ctx, GophKeeperService_SecretGetList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SecretGet(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Secret, error) {
	out := new(Secret)
	err := c.cc.Invoke(ctx, GophKeeperService_SecretGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SecretSet(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, GophKeeperService_SecretSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) SecretDel(ctx context.Context, in *Secret, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, GophKeeperService_SecretDel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServiceServer is the server API for GophKeeperService service.
// All implementations must embed UnimplementedGophKeeperServiceServer
// for forward compatibility
type GophKeeperServiceServer interface {
	// User
	UserCreate(context.Context, *User) (*UserAuthToken, error)
	UserLogin(context.Context, *User) (*UserAuthToken, error)
	// Secret
	SecretGetList(context.Context, *Empty) (*SecretList, error)
	SecretGet(context.Context, *Secret) (*Secret, error)
	SecretSet(context.Context, *Secret) (*Empty, error)
	SecretDel(context.Context, *Secret) (*Empty, error)
	mustEmbedUnimplementedGophKeeperServiceServer()
}

// UnimplementedGophKeeperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServiceServer struct {
}

func (UnimplementedGophKeeperServiceServer) UserCreate(context.Context, *User) (*UserAuthToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreate not implemented")
}
func (UnimplementedGophKeeperServiceServer) UserLogin(context.Context, *User) (*UserAuthToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLogin not implemented")
}
func (UnimplementedGophKeeperServiceServer) SecretGetList(context.Context, *Empty) (*SecretList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretGetList not implemented")
}
func (UnimplementedGophKeeperServiceServer) SecretGet(context.Context, *Secret) (*Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretGet not implemented")
}
func (UnimplementedGophKeeperServiceServer) SecretSet(context.Context, *Secret) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretSet not implemented")
}
func (UnimplementedGophKeeperServiceServer) SecretDel(context.Context, *Secret) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretDel not implemented")
}
func (UnimplementedGophKeeperServiceServer) mustEmbedUnimplementedGophKeeperServiceServer() {}

// UnsafeGophKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServiceServer will
// result in compilation errors.
type UnsafeGophKeeperServiceServer interface {
	mustEmbedUnimplementedGophKeeperServiceServer()
}

func RegisterGophKeeperServiceServer(s grpc.ServiceRegistrar, srv GophKeeperServiceServer) {
	s.RegisterService(&GophKeeperService_ServiceDesc, srv)
}

func _GophKeeperService_UserCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).UserCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_UserCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).UserCreate(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_UserLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).UserLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_UserLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).UserLogin(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SecretGetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SecretGetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SecretGetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SecretGetList(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SecretGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Secret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SecretGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SecretGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SecretGet(ctx, req.(*Secret))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SecretSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Secret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SecretSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SecretSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SecretSet(ctx, req.(*Secret))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_SecretDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Secret)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).SecretDel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_SecretDel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).SecretDel(ctx, req.(*Secret))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeperService_ServiceDesc is the grpc.ServiceDesc for GophKeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GophKeeperService",
	HandlerType: (*GophKeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreate",
			Handler:    _GophKeeperService_UserCreate_Handler,
		},
		{
			MethodName: "UserLogin",
			Handler:    _GophKeeperService_UserLogin_Handler,
		},
		{
			MethodName: "SecretGetList",
			Handler:    _GophKeeperService_SecretGetList_Handler,
		},
		{
			MethodName: "SecretGet",
			Handler:    _GophKeeperService_SecretGet_Handler,
		},
		{
			MethodName: "SecretSet",
			Handler:    _GophKeeperService_SecretSet_Handler,
		},
		{
			MethodName: "SecretDel",
			Handler:    _GophKeeperService_SecretDel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/service.proto",
}
