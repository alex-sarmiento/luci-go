// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/cipd/client/cipd/internal/messages/messages.proto
// DO NOT EDIT!

/*
Package messages is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/cipd/client/cipd/internal/messages/messages.proto

It has these top-level messages:
	BlobWithSHA1
	TagCache
	InstanceCache
*/
package messages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// BlobWithSHA1 is a wrapper around a binary blob with SHA1 hash to verify
// its integrity.
type BlobWithSHA1 struct {
	Blob []byte `protobuf:"bytes,1,opt,name=blob,proto3" json:"blob,omitempty"`
	Sha1 []byte `protobuf:"bytes,2,opt,name=sha1,proto3" json:"sha1,omitempty"`
}

func (m *BlobWithSHA1) Reset()                    { *m = BlobWithSHA1{} }
func (m *BlobWithSHA1) String() string            { return proto.CompactTextString(m) }
func (*BlobWithSHA1) ProtoMessage()               {}
func (*BlobWithSHA1) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BlobWithSHA1) GetBlob() []byte {
	if m != nil {
		return m.Blob
	}
	return nil
}

func (m *BlobWithSHA1) GetSha1() []byte {
	if m != nil {
		return m.Sha1
	}
	return nil
}

// TagCache stores a mapping (package name, tag) -> instance ID to
// speed up subsequent ResolveVersion calls when tags are used.
//
// It also contains a (instance_id, file_name) -> hash mapping which is used for
// client self-update purposes. file_name is case-senstive and must always use
// POSIX-style slashes.
type TagCache struct {
	// Capped list of entries, most recently resolved is last.
	Entries     []*TagCache_Entry     `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
	FileEntries []*TagCache_FileEntry `protobuf:"bytes,2,rep,name=file_entries,json=fileEntries" json:"file_entries,omitempty"`
}

func (m *TagCache) Reset()                    { *m = TagCache{} }
func (m *TagCache) String() string            { return proto.CompactTextString(m) }
func (*TagCache) ProtoMessage()               {}
func (*TagCache) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TagCache) GetEntries() []*TagCache_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *TagCache) GetFileEntries() []*TagCache_FileEntry {
	if m != nil {
		return m.FileEntries
	}
	return nil
}

type TagCache_Entry struct {
	Package    string `protobuf:"bytes,1,opt,name=package" json:"package,omitempty"`
	Tag        string `protobuf:"bytes,2,opt,name=tag" json:"tag,omitempty"`
	InstanceId string `protobuf:"bytes,3,opt,name=instance_id,json=instanceId" json:"instance_id,omitempty"`
}

func (m *TagCache_Entry) Reset()                    { *m = TagCache_Entry{} }
func (m *TagCache_Entry) String() string            { return proto.CompactTextString(m) }
func (*TagCache_Entry) ProtoMessage()               {}
func (*TagCache_Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *TagCache_Entry) GetPackage() string {
	if m != nil {
		return m.Package
	}
	return ""
}

func (m *TagCache_Entry) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *TagCache_Entry) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

type TagCache_FileEntry struct {
	Package    string `protobuf:"bytes,1,opt,name=package" json:"package,omitempty"`
	InstanceId string `protobuf:"bytes,2,opt,name=instance_id,json=instanceId" json:"instance_id,omitempty"`
	FileName   string `protobuf:"bytes,3,opt,name=file_name,json=fileName" json:"file_name,omitempty"`
	Hash       string `protobuf:"bytes,4,opt,name=hash" json:"hash,omitempty"`
}

func (m *TagCache_FileEntry) Reset()                    { *m = TagCache_FileEntry{} }
func (m *TagCache_FileEntry) String() string            { return proto.CompactTextString(m) }
func (*TagCache_FileEntry) ProtoMessage()               {}
func (*TagCache_FileEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *TagCache_FileEntry) GetPackage() string {
	if m != nil {
		return m.Package
	}
	return ""
}

func (m *TagCache_FileEntry) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

func (m *TagCache_FileEntry) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *TagCache_FileEntry) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

// InstanceCache stores a list of instances in cache
// and their last access time.
type InstanceCache struct {
	// Entries is a map of {instance id -> information about instance}.
	Entries map[string]*InstanceCache_Entry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// LastSynced is timestamp when we synchronized Entries with actual
	// instance files.
	LastSynced *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=last_synced,json=lastSynced" json:"last_synced,omitempty"`
}

func (m *InstanceCache) Reset()                    { *m = InstanceCache{} }
func (m *InstanceCache) String() string            { return proto.CompactTextString(m) }
func (*InstanceCache) ProtoMessage()               {}
func (*InstanceCache) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *InstanceCache) GetEntries() map[string]*InstanceCache_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *InstanceCache) GetLastSynced() *google_protobuf.Timestamp {
	if m != nil {
		return m.LastSynced
	}
	return nil
}

// Entry stores info about an instance.
type InstanceCache_Entry struct {
	// LastAccess is last time this instance was retrieved from or put to the
	// cache.
	LastAccess *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=last_access,json=lastAccess" json:"last_access,omitempty"`
}

func (m *InstanceCache_Entry) Reset()                    { *m = InstanceCache_Entry{} }
func (m *InstanceCache_Entry) String() string            { return proto.CompactTextString(m) }
func (*InstanceCache_Entry) ProtoMessage()               {}
func (*InstanceCache_Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *InstanceCache_Entry) GetLastAccess() *google_protobuf.Timestamp {
	if m != nil {
		return m.LastAccess
	}
	return nil
}

func init() {
	proto.RegisterType((*BlobWithSHA1)(nil), "messages.BlobWithSHA1")
	proto.RegisterType((*TagCache)(nil), "messages.TagCache")
	proto.RegisterType((*TagCache_Entry)(nil), "messages.TagCache.Entry")
	proto.RegisterType((*TagCache_FileEntry)(nil), "messages.TagCache.FileEntry")
	proto.RegisterType((*InstanceCache)(nil), "messages.InstanceCache")
	proto.RegisterType((*InstanceCache_Entry)(nil), "messages.InstanceCache.Entry")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/cipd/client/cipd/internal/messages/messages.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0xa5, 0xe9, 0xae, 0xdb, 0xde, 0x54, 0x90, 0x79, 0x0a, 0x51, 0xd9, 0xa5, 0xf8, 0xb0, 0x2f,
	0x26, 0x6c, 0x17, 0x44, 0x14, 0x94, 0xf5, 0x0b, 0xfb, 0xe2, 0x43, 0xb6, 0x20, 0x3e, 0x95, 0x9b,
	0xe9, 0x6d, 0x32, 0xec, 0x64, 0x52, 0x3a, 0x53, 0xa5, 0x3f, 0xca, 0xbf, 0xe2, 0x6f, 0x92, 0x99,
	0xc9, 0x54, 0x0d, 0x52, 0xf6, 0x25, 0x9c, 0x39, 0x39, 0xf7, 0xdc, 0x3b, 0x67, 0x2e, 0xcc, 0x2b,
	0x61, 0xea, 0x5d, 0x99, 0xf1, 0xb6, 0xc9, 0xe5, 0x8e, 0x0b, 0xf7, 0x79, 0x5e, 0xb5, 0x39, 0x17,
	0x9b, 0x55, 0xce, 0xa5, 0x20, 0x65, 0x3c, 0x16, 0xca, 0xd0, 0x56, 0xa1, 0xcc, 0x1b, 0xd2, 0x1a,
	0x2b, 0xd2, 0x07, 0x90, 0x6d, 0xb6, 0xad, 0x69, 0xd9, 0x28, 0x9c, 0xd3, 0xf3, 0xaa, 0x6d, 0x2b,
	0x49, 0xb9, 0xe3, 0xcb, 0xdd, 0x3a, 0x37, 0xa2, 0x21, 0x6d, 0xb0, 0xd9, 0x78, 0xe9, 0xf4, 0x05,
	0x4c, 0xde, 0xc9, 0xb6, 0xfc, 0x2a, 0x4c, 0x7d, 0xfb, 0xf9, 0xe6, 0x8a, 0x31, 0x38, 0x29, 0x65,
	0x5b, 0x26, 0x83, 0x8b, 0xc1, 0xe5, 0xa4, 0x70, 0xd8, 0x72, 0xba, 0xc6, 0xab, 0x24, 0xf2, 0x9c,
	0xc5, 0xd3, 0x5f, 0x11, 0x8c, 0x16, 0x58, 0xbd, 0x47, 0x5e, 0x13, 0x9b, 0xc1, 0x19, 0x29, 0xb3,
	0x15, 0xa4, 0x93, 0xc1, 0xc5, 0xf0, 0x32, 0x9e, 0x25, 0xd9, 0x61, 0xa2, 0x20, 0xca, 0x3e, 0x2a,
	0xb3, 0xdd, 0x17, 0x41, 0xc8, 0xde, 0xc2, 0x64, 0x2d, 0x24, 0x2d, 0x43, 0x61, 0xe4, 0x0a, 0x9f,
	0xfc, 0xa7, 0xf0, 0x93, 0x90, 0xe4, 0x8b, 0xe3, 0x75, 0x07, 0x05, 0xe9, 0x74, 0x01, 0xa7, 0x8e,
	0x65, 0x09, 0x9c, 0x6d, 0x90, 0xdf, 0x61, 0x45, 0x6e, 0xea, 0x71, 0x11, 0x8e, 0xec, 0x11, 0x0c,
	0x0d, 0x56, 0x6e, 0xee, 0x71, 0x61, 0x21, 0x3b, 0x87, 0x58, 0x28, 0x6d, 0x50, 0x71, 0x5a, 0x8a,
	0x55, 0x32, 0x74, 0x7f, 0x20, 0x50, 0xf3, 0x55, 0xfa, 0x03, 0xc6, 0x87, 0x7e, 0x47, 0x9c, 0x7b,
	0x3e, 0x51, 0xdf, 0x87, 0x3d, 0x86, 0xb1, 0xbb, 0x9e, 0xc2, 0x86, 0xba, 0x36, 0x23, 0x4b, 0x7c,
	0xc1, 0x86, 0x6c, 0xa0, 0x35, 0xea, 0x3a, 0x39, 0x71, 0xbc, 0xc3, 0xd3, 0x9f, 0x11, 0x3c, 0x9c,
	0x77, 0xf5, 0x3e, 0xd5, 0x37, 0xfd, 0x54, 0x9f, 0xfd, 0x09, 0xe7, 0x1f, 0x65, 0xd6, 0x45, 0xd2,
	0x4b, 0xf8, 0x35, 0xc4, 0x12, 0xb5, 0x59, 0xea, 0xbd, 0xe2, 0xe4, 0x67, 0x8c, 0x67, 0x69, 0xe6,
	0x37, 0x22, 0x0b, 0x1b, 0x91, 0x2d, 0xc2, 0x46, 0x14, 0x60, 0xe5, 0xb7, 0x4e, 0x9d, 0x7e, 0x08,
	0xe9, 0x06, 0x17, 0xe4, 0x9c, 0xb4, 0xbe, 0xaf, 0xcb, 0x8d, 0x53, 0xa7, 0xdf, 0x60, 0xf2, 0xf7,
	0x6c, 0xf6, 0x41, 0xee, 0x68, 0xdf, 0x85, 0x69, 0x21, 0xbb, 0x86, 0xd3, 0xef, 0x28, 0x77, 0xd4,
	0x19, 0x3f, 0x3d, 0x76, 0xc5, 0x7d, 0xe1, 0xb5, 0xaf, 0xa2, 0x97, 0x83, 0xf2, 0x81, 0x6b, 0x7d,
	0xfd, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x27, 0x58, 0x90, 0xbf, 0x37, 0x03, 0x00, 0x00,
}
