// Code generated by protoc-gen-go.
// source: sub.proto
// DO NOT EDIT!

/*
Package sub is a generated protocol buffer package.

It is generated from these files:
	sub.proto

It has these top-level messages:
	Sub
*/
package sub

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Sub struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *Sub) Reset()                    { *m = Sub{} }
func (m *Sub) String() string            { return proto.CompactTextString(m) }
func (*Sub) ProtoMessage()               {}
func (*Sub) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Sub)(nil), "sub.Sub")
}

var fileDescriptor0 = []byte{
	// 65 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x2e, 0x4d, 0xd2,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x06, 0x32, 0x95, 0x84, 0xb9, 0x98, 0x83, 0x4b, 0x93,
	0x84, 0x78, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x93, 0xd8,
	0xc0, 0x0a, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x63, 0x53, 0xa4, 0xae, 0x2d, 0x00, 0x00,
	0x00,
}