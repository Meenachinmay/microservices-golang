// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: enquiries.proto

package enquiries

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

// EnquiryServiceClient is the client API for EnquiryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnquiryServiceClient interface {
	HandleCustomerEnquiry(ctx context.Context, in *CustomerEnquiryRequest, opts ...grpc.CallOption) (*CustomerEnquiryResponse, error)
}

type enquiryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEnquiryServiceClient(cc grpc.ClientConnInterface) EnquiryServiceClient {
	return &enquiryServiceClient{cc}
}

func (c *enquiryServiceClient) HandleCustomerEnquiry(ctx context.Context, in *CustomerEnquiryRequest, opts ...grpc.CallOption) (*CustomerEnquiryResponse, error) {
	out := new(CustomerEnquiryResponse)
	err := c.cc.Invoke(ctx, "/enquiries.EnquiryService/HandleCustomerEnquiry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnquiryServiceServer is the server API for EnquiryService service.
// All implementations must embed UnimplementedEnquiryServiceServer
// for forward compatibility
type EnquiryServiceServer interface {
	HandleCustomerEnquiry(context.Context, *CustomerEnquiryRequest) (*CustomerEnquiryResponse, error)
	mustEmbedUnimplementedEnquiryServiceServer()
}

// UnimplementedEnquiryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEnquiryServiceServer struct {
}

func (UnimplementedEnquiryServiceServer) HandleCustomerEnquiry(context.Context, *CustomerEnquiryRequest) (*CustomerEnquiryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleCustomerEnquiry not implemented")
}
func (UnimplementedEnquiryServiceServer) mustEmbedUnimplementedEnquiryServiceServer() {}

// UnsafeEnquiryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnquiryServiceServer will
// result in compilation errors.
type UnsafeEnquiryServiceServer interface {
	mustEmbedUnimplementedEnquiryServiceServer()
}

func RegisterEnquiryServiceServer(s grpc.ServiceRegistrar, srv EnquiryServiceServer) {
	s.RegisterService(&EnquiryService_ServiceDesc, srv)
}

func _EnquiryService_HandleCustomerEnquiry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomerEnquiryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnquiryServiceServer).HandleCustomerEnquiry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/enquiries.EnquiryService/HandleCustomerEnquiry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnquiryServiceServer).HandleCustomerEnquiry(ctx, req.(*CustomerEnquiryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EnquiryService_ServiceDesc is the grpc.ServiceDesc for EnquiryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EnquiryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "enquiries.EnquiryService",
	HandlerType: (*EnquiryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleCustomerEnquiry",
			Handler:    _EnquiryService_HandleCustomerEnquiry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "enquiries.proto",
}
