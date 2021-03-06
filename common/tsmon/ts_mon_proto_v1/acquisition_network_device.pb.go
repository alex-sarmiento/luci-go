// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/common/tsmon/ts_mon_proto_v1/acquisition_network_device.proto
// DO NOT EDIT!

/*
Package ts_mon_proto_v1 is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v1/acquisition_network_device.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v1/acquisition_task.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v1/metrics.proto

It has these top-level messages:
	NetworkDevice
	Task
	MetricsCollection
	MetricsField
	PrecomputedDistribution
	MetricsData
*/
package ts_mon_proto_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NetworkDevice_TypeId int32

const (
	NetworkDevice_MESSAGE_TYPE_ID NetworkDevice_TypeId = 34049749
)

var NetworkDevice_TypeId_name = map[int32]string{
	34049749: "MESSAGE_TYPE_ID",
}
var NetworkDevice_TypeId_value = map[string]int32{
	"MESSAGE_TYPE_ID": 34049749,
}

func (x NetworkDevice_TypeId) Enum() *NetworkDevice_TypeId {
	p := new(NetworkDevice_TypeId)
	*p = x
	return p
}
func (x NetworkDevice_TypeId) String() string {
	return proto.EnumName(NetworkDevice_TypeId_name, int32(x))
}
func (x *NetworkDevice_TypeId) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(NetworkDevice_TypeId_value, data, "NetworkDevice_TypeId")
	if err != nil {
		return err
	}
	*x = NetworkDevice_TypeId(value)
	return nil
}
func (NetworkDevice_TypeId) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type NetworkDevice struct {
	Alertable        *bool   `protobuf:"varint,101,opt,name=alertable" json:"alertable,omitempty"`
	Realm            *string `protobuf:"bytes,102,opt,name=realm" json:"realm,omitempty"`
	Metro            *string `protobuf:"bytes,104,opt,name=metro" json:"metro,omitempty"`
	Role             *string `protobuf:"bytes,105,opt,name=role" json:"role,omitempty"`
	Hostname         *string `protobuf:"bytes,106,opt,name=hostname" json:"hostname,omitempty"`
	Hostgroup        *string `protobuf:"bytes,108,opt,name=hostgroup" json:"hostgroup,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NetworkDevice) Reset()                    { *m = NetworkDevice{} }
func (m *NetworkDevice) String() string            { return proto.CompactTextString(m) }
func (*NetworkDevice) ProtoMessage()               {}
func (*NetworkDevice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NetworkDevice) GetAlertable() bool {
	if m != nil && m.Alertable != nil {
		return *m.Alertable
	}
	return false
}

func (m *NetworkDevice) GetRealm() string {
	if m != nil && m.Realm != nil {
		return *m.Realm
	}
	return ""
}

func (m *NetworkDevice) GetMetro() string {
	if m != nil && m.Metro != nil {
		return *m.Metro
	}
	return ""
}

func (m *NetworkDevice) GetRole() string {
	if m != nil && m.Role != nil {
		return *m.Role
	}
	return ""
}

func (m *NetworkDevice) GetHostname() string {
	if m != nil && m.Hostname != nil {
		return *m.Hostname
	}
	return ""
}

func (m *NetworkDevice) GetHostgroup() string {
	if m != nil && m.Hostgroup != nil {
		return *m.Hostgroup
	}
	return ""
}

func init() {
	proto.RegisterType((*NetworkDevice)(nil), "ts_mon.proto.v1.NetworkDevice")
	proto.RegisterEnum("ts_mon.proto.v1.NetworkDevice_TypeId", NetworkDevice_TypeId_name, NetworkDevice_TypeId_value)
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/common/tsmon/ts_mon_proto_v1/acquisition_network_device.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x44, 0x8d, 0xd1, 0x4a, 0xc3, 0x30,
	0x14, 0x86, 0x29, 0xa8, 0x6c, 0x01, 0xd9, 0x08, 0x22, 0x41, 0xbc, 0x28, 0xbb, 0xda, 0x8d, 0x2d,
	0x7b, 0x04, 0x61, 0x45, 0x76, 0xa1, 0xc8, 0x36, 0x2f, 0xbc, 0x0a, 0x59, 0x76, 0x6c, 0xa3, 0x49,
	0x4e, 0x4d, 0x93, 0x8a, 0x0f, 0xa3, 0xef, 0xe3, 0x03, 0xf8, 0x3e, 0xd2, 0x04, 0xec, 0xcd, 0xcf,
	0xff, 0x7d, 0xf9, 0xc3, 0x21, 0x4f, 0xb5, 0xf2, 0x4d, 0x38, 0x14, 0x12, 0x4d, 0xa9, 0x83, 0x54,
	0x31, 0x6e, 0x6a, 0x2c, 0x25, 0x1a, 0x83, 0xb6, 0xf4, 0x5d, 0x4a, 0x6e, 0xd0, 0xf2, 0xd6, 0xa1,
	0x47, 0xde, 0xaf, 0x4a, 0x21, 0xdf, 0x83, 0xea, 0x94, 0x57, 0x68, 0xb9, 0x05, 0xff, 0x81, 0xee,
	0x8d, 0x1f, 0xa1, 0x57, 0x12, 0x8a, 0xb8, 0xa1, 0xb3, 0xf4, 0x23, 0x51, 0xd1, 0xaf, 0x16, 0x3f,
	0x19, 0x39, 0x7f, 0x48, 0xcb, 0x75, 0x1c, 0xd2, 0x6b, 0x32, 0x15, 0x1a, 0x9c, 0x17, 0x07, 0x0d,
	0x0c, 0xf2, 0x6c, 0x39, 0xd9, 0x8e, 0x82, 0x5e, 0x90, 0x53, 0x07, 0x42, 0x1b, 0xf6, 0x92, 0x67,
	0xcb, 0xe9, 0x36, 0xc1, 0x60, 0x0d, 0x78, 0x87, 0xac, 0x49, 0x36, 0x02, 0xa5, 0xe4, 0xc4, 0xa1,
	0x06, 0xa6, 0xa2, 0x8c, 0x9d, 0x5e, 0x91, 0x49, 0x83, 0x9d, 0xb7, 0xc2, 0x00, 0x7b, 0x8d, 0xfe,
	0x9f, 0x87, 0xcb, 0x43, 0xaf, 0x1d, 0x86, 0x96, 0xe9, 0xf8, 0x38, 0x8a, 0x45, 0x4e, 0xce, 0xf6,
	0x9f, 0x2d, 0x6c, 0x8e, 0xf4, 0x92, 0xcc, 0xee, 0xab, 0xdd, 0xee, 0xf6, 0xae, 0xe2, 0xfb, 0xe7,
	0xc7, 0x8a, 0x6f, 0xd6, 0xf3, 0xdf, 0xaf, 0xef, 0xf9, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6e,
	0xe6, 0xaa, 0xbf, 0x34, 0x01, 0x00, 0x00,
}
