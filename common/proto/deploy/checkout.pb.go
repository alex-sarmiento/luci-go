// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/common/proto/deploy/checkout.proto
// DO NOT EDIT!

/*
Package deploy is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/common/proto/deploy/checkout.proto
	github.com/luci/luci-go/common/proto/deploy/component.proto
	github.com/luci/luci-go/common/proto/deploy/config.proto
	github.com/luci/luci-go/common/proto/deploy/internal.proto
	github.com/luci/luci-go/common/proto/deploy/userconfig.proto

It has these top-level messages:
	SourceLayout
	SourceInitResult
	GoPath
	Component
	BuildPath
	AppEngineModule
	AppEngineResources
	ContainerEnginePod
	KubernetesPod
	Layout
	Source
	Application
	Deployment
	FrozenLayout
	UserConfig
*/
package deploy

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

// *
// Source layout configuration file.
//
// Each Source checkout may include a textproto layout file named
// "luci-deploytool.cfg". If present, this file will be loaded and used to
// integrate the source into the deployment.
//
// Uncooperative repositories, or those not owned by the author, may have the
// same effect by specifying the SourceLayout in the Source's layout definition.
type SourceLayout struct {
	// * The source initialization operations to execute, in order.
	Init []*SourceLayout_Init `protobuf:"bytes,1,rep,name=init" json:"init,omitempty"`
	// * Go Paths to add to this repository.
	GoPath []*GoPath `protobuf:"bytes,10,rep,name=go_path,json=goPath" json:"go_path,omitempty"`
}

func (m *SourceLayout) Reset()                    { *m = SourceLayout{} }
func (m *SourceLayout) String() string            { return proto.CompactTextString(m) }
func (*SourceLayout) ProtoMessage()               {}
func (*SourceLayout) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SourceLayout) GetInit() []*SourceLayout_Init {
	if m != nil {
		return m.Init
	}
	return nil
}

func (m *SourceLayout) GetGoPath() []*GoPath {
	if m != nil {
		return m.GoPath
	}
	return nil
}

// *
// Init is a single initialization operation execution.
type SourceLayout_Init struct {
	// Types that are valid to be assigned to Operation:
	//	*SourceLayout_Init_PythonScript_
	Operation isSourceLayout_Init_Operation `protobuf_oneof:"operation"`
}

func (m *SourceLayout_Init) Reset()                    { *m = SourceLayout_Init{} }
func (m *SourceLayout_Init) String() string            { return proto.CompactTextString(m) }
func (*SourceLayout_Init) ProtoMessage()               {}
func (*SourceLayout_Init) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type isSourceLayout_Init_Operation interface {
	isSourceLayout_Init_Operation()
}

type SourceLayout_Init_PythonScript_ struct {
	PythonScript *SourceLayout_Init_PythonScript `protobuf:"bytes,1,opt,name=python_script,json=pythonScript,oneof"`
}

func (*SourceLayout_Init_PythonScript_) isSourceLayout_Init_Operation() {}

func (m *SourceLayout_Init) GetOperation() isSourceLayout_Init_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (m *SourceLayout_Init) GetPythonScript() *SourceLayout_Init_PythonScript {
	if x, ok := m.GetOperation().(*SourceLayout_Init_PythonScript_); ok {
		return x.PythonScript
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SourceLayout_Init) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SourceLayout_Init_OneofMarshaler, _SourceLayout_Init_OneofUnmarshaler, _SourceLayout_Init_OneofSizer, []interface{}{
		(*SourceLayout_Init_PythonScript_)(nil),
	}
}

func _SourceLayout_Init_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SourceLayout_Init)
	// operation
	switch x := m.Operation.(type) {
	case *SourceLayout_Init_PythonScript_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.PythonScript); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SourceLayout_Init.Operation has unexpected type %T", x)
	}
	return nil
}

func _SourceLayout_Init_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SourceLayout_Init)
	switch tag {
	case 1: // operation.python_script
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SourceLayout_Init_PythonScript)
		err := b.DecodeMessage(msg)
		m.Operation = &SourceLayout_Init_PythonScript_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SourceLayout_Init_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SourceLayout_Init)
	// operation
	switch x := m.Operation.(type) {
	case *SourceLayout_Init_PythonScript_:
		s := proto.Size(x.PythonScript)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// *
// A Python initialization script.
//
// The script will be run as follows:
// $ PYTHON PATH SOURCE-ROOT RESULT-PATH
//
// - PYTHON is the deploy tool-resolved Python interpreter.
// - PATH is the absolute path of the script.
// - SOURCE-ROOT is the root of the source that is being initialized.
// - RESULT-PATH is the path where, optionally, a SourceInitResult protobuf
//   may be written. If the file is present, it will be read and linked into
//   the deployment runtime.
type SourceLayout_Init_PythonScript struct {
	// * The source-relative path of the Python script.
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *SourceLayout_Init_PythonScript) Reset()         { *m = SourceLayout_Init_PythonScript{} }
func (m *SourceLayout_Init_PythonScript) String() string { return proto.CompactTextString(m) }
func (*SourceLayout_Init_PythonScript) ProtoMessage()    {}
func (*SourceLayout_Init_PythonScript) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 0, 0}
}

// *
// SourceInitResult is a protobuf that can be emitted from a SourceInit Script
// to describe how to link the results of that initialization into the
// deployment layout.
type SourceInitResult struct {
	// * Source-relative entries to add to GOPATH.
	GoPath []*GoPath `protobuf:"bytes,1,rep,name=go_path,json=goPath" json:"go_path,omitempty"`
}

func (m *SourceInitResult) Reset()                    { *m = SourceInitResult{} }
func (m *SourceInitResult) String() string            { return proto.CompactTextString(m) }
func (*SourceInitResult) ProtoMessage()               {}
func (*SourceInitResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SourceInitResult) GetGoPath() []*GoPath {
	if m != nil {
		return m.GoPath
	}
	return nil
}

// *
// Describes how to link a source-relative directory into the generated GOPATH.
type GoPath struct {
	// *
	// The source-relative path to add to GOPATH. If empty, this is the source
	// root.
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	// *
	// The name of the Go package to bind to "path".
	//
	// For example, given checkout:
	//   path: gosrc/my/thing
	//   go_package: github.com/example/mything
	//
	// This will add a GOPATH entry:
	// src/github.com/example/mything => <root>/gosrc/my/thing
	GoPackage string `protobuf:"bytes,2,opt,name=go_package,json=goPackage" json:"go_package,omitempty"`
}

func (m *GoPath) Reset()                    { *m = GoPath{} }
func (m *GoPath) String() string            { return proto.CompactTextString(m) }
func (*GoPath) ProtoMessage()               {}
func (*GoPath) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*SourceLayout)(nil), "deploy.SourceLayout")
	proto.RegisterType((*SourceLayout_Init)(nil), "deploy.SourceLayout.Init")
	proto.RegisterType((*SourceLayout_Init_PythonScript)(nil), "deploy.SourceLayout.Init.PythonScript")
	proto.RegisterType((*SourceInitResult)(nil), "deploy.SourceInitResult")
	proto.RegisterType((*GoPath)(nil), "deploy.GoPath")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/common/proto/deploy/checkout.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x90, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0xc6, 0xad, 0x96, 0x4a, 0x67, 0xab, 0x48, 0x4e, 0x75, 0x41, 0x90, 0x1e, 0xd4, 0xcb, 0x26,
	0xa0, 0x37, 0xbd, 0x79, 0x51, 0x41, 0x61, 0xe9, 0x3e, 0xc0, 0xd2, 0x8d, 0x21, 0x0d, 0xdb, 0xed,
	0x84, 0x6e, 0x72, 0xe8, 0x0b, 0xf8, 0xba, 0xbe, 0x82, 0xe9, 0x14, 0xa1, 0x82, 0x7a, 0xc9, 0x9f,
	0x6f, 0xbe, 0xdf, 0x7c, 0x99, 0xc0, 0xbd, 0x36, 0xae, 0xf6, 0x1b, 0x2e, 0x71, 0x27, 0x1a, 0x2f,
	0x0d, 0x2d, 0x0b, 0x8d, 0x22, 0x08, 0x3b, 0x6c, 0x85, 0xed, 0xd0, 0xa1, 0x78, 0x57, 0xb6, 0xc1,
	0x5e, 0xc8, 0x5a, 0xc9, 0x2d, 0x7a, 0xc7, 0x49, 0x65, 0xc9, 0x28, 0x17, 0x9f, 0x11, 0x64, 0x2b,
	0xf4, 0x9d, 0x54, 0xaf, 0x55, 0x1f, 0xca, 0x6c, 0x01, 0xb1, 0x69, 0x8d, 0xcb, 0xa3, 0xcb, 0xa3,
	0x9b, 0xd9, 0xed, 0x39, 0x1f, 0x7d, 0x7c, 0xea, 0xe1, 0x2f, 0xc1, 0x50, 0x92, 0x8d, 0x5d, 0xc3,
	0xb1, 0xc6, 0xb5, 0xad, 0x5c, 0x9d, 0x03, 0x11, 0xa7, 0xdf, 0xc4, 0x13, 0x2e, 0x83, 0x5a, 0x26,
	0x9a, 0xf6, 0xf9, 0x47, 0x04, 0xf1, 0xc0, 0xb1, 0x37, 0x38, 0xb1, 0xbd, 0xab, 0xb1, 0x5d, 0xef,
	0x65, 0x67, 0xec, 0x90, 0x14, 0x05, 0xee, 0xea, 0xcf, 0x24, 0xbe, 0x24, 0xfb, 0x8a, 0xdc, 0xcf,
	0x07, 0x65, 0x66, 0x27, 0xf7, 0x79, 0x01, 0xd9, 0xb4, 0xce, 0x18, 0xc4, 0xf4, 0x9a, 0xa1, 0x6b,
	0x5a, 0xd2, 0xf9, 0x71, 0x06, 0x29, 0x5a, 0xd5, 0x55, 0xce, 0x60, 0x5b, 0x3c, 0xc0, 0xd9, 0x18,
	0x41, 0x53, 0xa8, 0xbd, 0x6f, 0x7e, 0x4c, 0x11, 0xfd, 0x37, 0x45, 0x80, 0x93, 0x51, 0xf9, 0x2d,
	0x87, 0x5d, 0x00, 0x50, 0x1b, 0xb9, 0xad, 0xb4, 0xca, 0x0f, 0xa9, 0x92, 0x0e, 0x24, 0x09, 0x9b,
	0x84, 0xbe, 0xfe, 0xee, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xe0, 0x9a, 0x59, 0x56, 0xb8, 0x01, 0x00,
	0x00,
}