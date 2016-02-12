// Code generated by protoc-gen-go.
// source: service.proto
// DO NOT EDIT!

/*
Package discovery is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	Void
	DescribeResponse
*/
package discovery

import prpccommon "github.com/luci/luci-go/common/prpc"
import prpc "github.com/luci/luci-go/server/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google/descriptor"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Void is an empty message.
type Void struct {
}

func (m *Void) Reset()                    { *m = Void{} }
func (m *Void) String() string            { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()               {}
func (*Void) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// DescribeResponse describes services.
type DescribeResponse struct {
	// Description contains descriptions of all services, their types and all
	// transitive dependencies.
	Description *google_protobuf.FileDescriptorSet `protobuf:"bytes,1,opt,name=description" json:"description,omitempty"`
	// Services are service names provided by a server.
	Services []string `protobuf:"bytes,2,rep,name=services" json:"services,omitempty"`
}

func (m *DescribeResponse) Reset()                    { *m = DescribeResponse{} }
func (m *DescribeResponse) String() string            { return proto.CompactTextString(m) }
func (*DescribeResponse) ProtoMessage()               {}
func (*DescribeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DescribeResponse) GetDescription() *google_protobuf.FileDescriptorSet {
	if m != nil {
		return m.Description
	}
	return nil
}

func init() {
	proto.RegisterType((*Void)(nil), "discovery.Void")
	proto.RegisterType((*DescribeResponse)(nil), "discovery.DescribeResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Discovery service

type DiscoveryClient interface {
	// Returns a list of services and a
	// descriptor.FileDescriptorSet that covers them all.
	Describe(ctx context.Context, in *Void, opts ...grpc.CallOption) (*DescribeResponse, error)
}
type discoveryPRPCClient struct {
	client *prpccommon.Client
}

func NewDiscoveryPRPCClient(client *prpccommon.Client) DiscoveryClient {
	return &discoveryPRPCClient{client}
}

func (c *discoveryPRPCClient) Describe(ctx context.Context, in *Void, opts ...grpc.CallOption) (*DescribeResponse, error) {
	out := new(DescribeResponse)
	err := c.client.Call(ctx, "discovery.Discovery", "Describe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type discoveryClient struct {
	cc *grpc.ClientConn
}

func NewDiscoveryClient(cc *grpc.ClientConn) DiscoveryClient {
	return &discoveryClient{cc}
}

func (c *discoveryClient) Describe(ctx context.Context, in *Void, opts ...grpc.CallOption) (*DescribeResponse, error) {
	out := new(DescribeResponse)
	err := grpc.Invoke(ctx, "/discovery.Discovery/Describe", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Discovery service

type DiscoveryServer interface {
	// Returns a list of services and a
	// descriptor.FileDescriptorSet that covers them all.
	Describe(context.Context, *Void) (*DescribeResponse, error)
}

func RegisterDiscoveryServer(s prpc.Registrar, srv DiscoveryServer) {
	s.RegisterService(&_Discovery_serviceDesc, srv)
}

func _Discovery_Describe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(DiscoveryServer).Describe(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Discovery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discovery.Discovery",
	HandlerType: (*DiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Describe",
			Handler:    _Discovery_Describe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x8e, 0xb1, 0x0b, 0x82, 0x40,
	0x18, 0xc5, 0xb3, 0x42, 0xf4, 0x93, 0x48, 0x6e, 0x12, 0x5b, 0xe4, 0x26, 0xa7, 0x13, 0x6c, 0x08,
	0x9a, 0xa5, 0xf6, 0x82, 0xb6, 0x16, 0xf5, 0x4b, 0x0e, 0xc4, 0x4f, 0xee, 0x4c, 0xe8, 0xbf, 0x8f,
	0xce, 0x4e, 0xa2, 0xf5, 0xbd, 0xc7, 0xef, 0xfd, 0x60, 0xa3, 0x51, 0x8d, 0xb2, 0x42, 0xd1, 0x2b,
	0x1a, 0x88, 0xf9, 0xb5, 0xd4, 0x15, 0x8d, 0xa8, 0x5e, 0x71, 0xd2, 0x10, 0x35, 0x2d, 0x66, 0xa6,
	0x28, 0x9f, 0x8f, 0xac, 0x46, 0x5d, 0x29, 0xd9, 0x0f, 0xa4, 0xa6, 0x31, 0x77, 0x61, 0x7d, 0x23,
	0x59, 0xf3, 0x3b, 0x84, 0x85, 0xe9, 0x4a, 0xbc, 0xa0, 0xee, 0xa9, 0xd3, 0xc8, 0x0e, 0x10, 0xd8,
	0xbd, 0xa4, 0x2e, 0x72, 0x12, 0x27, 0x0d, 0x72, 0x2e, 0x26, 0xa6, 0xb0, 0x4c, 0x71, 0x92, 0x2d,
	0x16, 0x33, 0xf7, 0x8a, 0x03, 0x0b, 0xc1, 0xfb, 0x2a, 0xe9, 0x68, 0x99, 0xac, 0x52, 0x3f, 0x3f,
	0x83, 0x5f, 0x58, 0x2b, 0x76, 0x04, 0xcf, 0x7e, 0xb1, 0xad, 0x98, 0x6d, 0xc5, 0x47, 0x24, 0xde,
	0xfd, 0x04, 0xff, 0x46, 0x7c, 0x51, 0xba, 0xe6, 0x76, 0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xf5,
	0xe3, 0xb6, 0x9a, 0xf4, 0x00, 0x00, 0x00,
}
