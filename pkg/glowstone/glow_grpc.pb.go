// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package glowstone

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// GlowClient is the client API for Glow service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GlowClient interface {
	Listen(ctx context.Context, opts ...grpc.CallOption) (Glow_ListenClient, error)
}

type glowClient struct {
	cc grpc.ClientConnInterface
}

func NewGlowClient(cc grpc.ClientConnInterface) GlowClient {
	return &glowClient{cc}
}

func (c *glowClient) Listen(ctx context.Context, opts ...grpc.CallOption) (Glow_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Glow_serviceDesc.Streams[0], "/glowstone.Glow/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &glowListenClient{stream}
	return x, nil
}

type Glow_ListenClient interface {
	Send(*Tick) error
	Recv() (*Tick, error)
	grpc.ClientStream
}

type glowListenClient struct {
	grpc.ClientStream
}

func (x *glowListenClient) Send(m *Tick) error {
	return x.ClientStream.SendMsg(m)
}

func (x *glowListenClient) Recv() (*Tick, error) {
	m := new(Tick)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GlowServer is the server API for Glow service.
// All implementations must embed UnimplementedGlowServer
// for forward compatibility
type GlowServer interface {
	Listen(Glow_ListenServer) error
	mustEmbedUnimplementedGlowServer()
}

// UnimplementedGlowServer must be embedded to have forward compatible implementations.
type UnimplementedGlowServer struct {
}

func (UnimplementedGlowServer) Listen(Glow_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedGlowServer) mustEmbedUnimplementedGlowServer() {}

// UnsafeGlowServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GlowServer will
// result in compilation errors.
type UnsafeGlowServer interface {
	mustEmbedUnimplementedGlowServer()
}

func RegisterGlowServer(s grpc.ServiceRegistrar, srv GlowServer) {
	s.RegisterService(&_Glow_serviceDesc, srv)
}

func _Glow_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GlowServer).Listen(&glowListenServer{stream})
}

type Glow_ListenServer interface {
	Send(*Tick) error
	Recv() (*Tick, error)
	grpc.ServerStream
}

type glowListenServer struct {
	grpc.ServerStream
}

func (x *glowListenServer) Send(m *Tick) error {
	return x.ServerStream.SendMsg(m)
}

func (x *glowListenServer) Recv() (*Tick, error) {
	m := new(Tick)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Glow_serviceDesc = grpc.ServiceDesc{
	ServiceName: "glowstone.Glow",
	HandlerType: (*GlowServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Glow_Listen_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/glowstone/pb/glow.proto",
}
