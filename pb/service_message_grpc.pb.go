// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// MessageClient is the client API for Message service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageClient interface {
	SendMessage(ctx context.Context, opts ...grpc.CallOption) (Message_SendMessageClient, error)
}

type messageClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageClient(cc grpc.ClientConnInterface) MessageClient {
	return &messageClient{cc}
}

func (c *messageClient) SendMessage(ctx context.Context, opts ...grpc.CallOption) (Message_SendMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &Message_ServiceDesc.Streams[0], "/pb.Message/SendMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageSendMessageClient{stream}
	return x, nil
}

type Message_SendMessageClient interface {
	Send(*SendMessageRequest) error
	Recv() (*SendMessageResponse, error)
	grpc.ClientStream
}

type messageSendMessageClient struct {
	grpc.ClientStream
}

func (x *messageSendMessageClient) Send(m *SendMessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageSendMessageClient) Recv() (*SendMessageResponse, error) {
	m := new(SendMessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServer is the server API for Message service.
// All implementations must embed UnimplementedMessageServer
// for forward compatibility
type MessageServer interface {
	SendMessage(Message_SendMessageServer) error
	mustEmbedUnimplementedMessageServer()
}

// UnimplementedMessageServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServer struct {
}

func (UnimplementedMessageServer) SendMessage(Message_SendMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessageServer) mustEmbedUnimplementedMessageServer() {}

// UnsafeMessageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServer will
// result in compilation errors.
type UnsafeMessageServer interface {
	mustEmbedUnimplementedMessageServer()
}

func RegisterMessageServer(s grpc.ServiceRegistrar, srv MessageServer) {
	s.RegisterService(&Message_ServiceDesc, srv)
}

func _Message_SendMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageServer).SendMessage(&messageSendMessageServer{stream})
}

type Message_SendMessageServer interface {
	Send(*SendMessageResponse) error
	Recv() (*SendMessageRequest, error)
	grpc.ServerStream
}

type messageSendMessageServer struct {
	grpc.ServerStream
}

func (x *messageSendMessageServer) Send(m *SendMessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageSendMessageServer) Recv() (*SendMessageRequest, error) {
	m := new(SendMessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Message_ServiceDesc is the grpc.ServiceDesc for Message service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Message_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Message",
	HandlerType: (*MessageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendMessage",
			Handler:       _Message_SendMessage_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service_message.proto",
}
