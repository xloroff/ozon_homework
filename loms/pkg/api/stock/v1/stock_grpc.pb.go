// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.0
// source: stock.proto

package stock

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	StockAPI_Info_FullMethodName = "/gitlab.ozon.dev.xloroff.ozon_hw.loms.pkg.api.stock.v1.StockAPI/Info"
)

// StockAPIClient is the client API for StockAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StockAPIClient interface {
	Info(ctx context.Context, in *StockInfoRequest, opts ...grpc.CallOption) (*StockInfoResponse, error)
}

type stockAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewStockAPIClient(cc grpc.ClientConnInterface) StockAPIClient {
	return &stockAPIClient{cc}
}

func (c *stockAPIClient) Info(ctx context.Context, in *StockInfoRequest, opts ...grpc.CallOption) (*StockInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StockInfoResponse)
	err := c.cc.Invoke(ctx, StockAPI_Info_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StockAPIServer is the server API for StockAPI service.
// All implementations must embed UnimplementedStockAPIServer
// for forward compatibility
type StockAPIServer interface {
	Info(context.Context, *StockInfoRequest) (*StockInfoResponse, error)
	mustEmbedUnimplementedStockAPIServer()
}

// UnimplementedStockAPIServer must be embedded to have forward compatible implementations.
type UnimplementedStockAPIServer struct {
}

func (UnimplementedStockAPIServer) Info(context.Context, *StockInfoRequest) (*StockInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedStockAPIServer) mustEmbedUnimplementedStockAPIServer() {}

// UnsafeStockAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StockAPIServer will
// result in compilation errors.
type UnsafeStockAPIServer interface {
	mustEmbedUnimplementedStockAPIServer()
}

func RegisterStockAPIServer(s grpc.ServiceRegistrar, srv StockAPIServer) {
	s.RegisterService(&StockAPI_ServiceDesc, srv)
}

func _StockAPI_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockAPIServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StockAPI_Info_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockAPIServer).Info(ctx, req.(*StockInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StockAPI_ServiceDesc is the grpc.ServiceDesc for StockAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StockAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gitlab.ozon.dev.xloroff.ozon_hw.loms.pkg.api.stock.v1.StockAPI",
	HandlerType: (*StockAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Info",
			Handler:    _StockAPI_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stock.proto",
}
