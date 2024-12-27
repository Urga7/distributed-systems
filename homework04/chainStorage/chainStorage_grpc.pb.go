// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.2
// source: chainStorage.proto

package chainStorage

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
	ChainReplication_Put_FullMethodName    = "/chainStorage.ChainReplication/Put"
	ChainReplication_Get_FullMethodName    = "/chainStorage.ChainReplication/Get"
	ChainReplication_Commit_FullMethodName = "/chainStorage.ChainReplication/Commit"
)

// ChainReplicationClient is the client API for ChainReplication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChainReplicationClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Commit(ctx context.Context, in *Todo, opts ...grpc.CallOption) (*PutResponse, error)
}

type chainReplicationClient struct {
	cc grpc.ClientConnInterface
}

func NewChainReplicationClient(cc grpc.ClientConnInterface) ChainReplicationClient {
	return &chainReplicationClient{cc}
}

func (c *chainReplicationClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, ChainReplication_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainReplicationClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, ChainReplication_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainReplicationClient) Commit(ctx context.Context, in *Todo, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, ChainReplication_Commit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChainReplicationServer is the server API for ChainReplication service.
// All implementations must embed UnimplementedChainReplicationServer
// for forward compatibility
type ChainReplicationServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Commit(context.Context, *Todo) (*PutResponse, error)
	mustEmbedUnimplementedChainReplicationServer()
}

// UnimplementedChainReplicationServer must be embedded to have forward compatible implementations.
type UnimplementedChainReplicationServer struct {
}

func (UnimplementedChainReplicationServer) Put(context.Context, *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedChainReplicationServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedChainReplicationServer) Commit(context.Context, *Todo) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedChainReplicationServer) mustEmbedUnimplementedChainReplicationServer() {}

// UnsafeChainReplicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChainReplicationServer will
// result in compilation errors.
type UnsafeChainReplicationServer interface {
	mustEmbedUnimplementedChainReplicationServer()
}

func RegisterChainReplicationServer(s grpc.ServiceRegistrar, srv ChainReplicationServer) {
	s.RegisterService(&ChainReplication_ServiceDesc, srv)
}

func _ChainReplication_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainReplicationServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChainReplication_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainReplicationServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChainReplication_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainReplicationServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChainReplication_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainReplicationServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChainReplication_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Todo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainReplicationServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChainReplication_Commit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainReplicationServer).Commit(ctx, req.(*Todo))
	}
	return interceptor(ctx, in, info, handler)
}

// ChainReplication_ServiceDesc is the grpc.ServiceDesc for ChainReplication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChainReplication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chainStorage.ChainReplication",
	HandlerType: (*ChainReplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _ChainReplication_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ChainReplication_Get_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _ChainReplication_Commit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chainStorage.proto",
}
