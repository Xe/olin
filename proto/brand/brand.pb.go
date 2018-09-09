// Code generated by protoc-gen-go. DO NOT EDIT.
// source: brand.proto

package brand

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

type Kind int32

const (
	Kind_UNKNOWN     Kind = 0
	Kind_CWAGI       Kind = 1
	Kind_DAEMON      Kind = 2
	Kind_INTERACTIVE Kind = 3
)

var Kind_name = map[int32]string{
	0: "UNKNOWN",
	1: "CWAGI",
	2: "DAEMON",
	3: "INTERACTIVE",
}
var Kind_value = map[string]int32{
	"UNKNOWN":     0,
	"CWAGI":       1,
	"DAEMON":      2,
	"INTERACTIVE": 3,
}

func (x Kind) String() string {
	return proto.EnumName(Kind_name, int32(x))
}
func (Kind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_brand_141de6d19d21dc5e, []int{0}
}

type Brand struct {
	Opts                 *VMOptions `protobuf:"bytes,1,opt,name=opts,proto3" json:"opts,omitempty"`
	Meta                 *Metadata  `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Brand) Reset()         { *m = Brand{} }
func (m *Brand) String() string { return proto.CompactTextString(m) }
func (*Brand) ProtoMessage()    {}
func (*Brand) Descriptor() ([]byte, []int) {
	return fileDescriptor_brand_141de6d19d21dc5e, []int{0}
}
func (m *Brand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Brand.Unmarshal(m, b)
}
func (m *Brand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Brand.Marshal(b, m, deterministic)
}
func (dst *Brand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Brand.Merge(dst, src)
}
func (m *Brand) XXX_Size() int {
	return xxx_messageInfo_Brand.Size(m)
}
func (m *Brand) XXX_DiscardUnknown() {
	xxx_messageInfo_Brand.DiscardUnknown(m)
}

var xxx_messageInfo_Brand proto.InternalMessageInfo

func (m *Brand) GetOpts() *VMOptions {
	if m != nil {
		return m.Opts
	}
	return nil
}

func (m *Brand) GetMeta() *Metadata {
	if m != nil {
		return m.Meta
	}
	return nil
}

type VMOptions struct {
	EnableJit            bool     `protobuf:"varint,1,opt,name=enable_jit,json=enableJit,proto3" json:"enable_jit,omitempty"`
	DefaultPages         int32    `protobuf:"varint,2,opt,name=default_pages,json=defaultPages,proto3" json:"default_pages,omitempty"`
	MaxPages             int32    `protobuf:"varint,3,opt,name=max_pages,json=maxPages,proto3" json:"max_pages,omitempty"`
	MainFunc             string   `protobuf:"bytes,4,opt,name=main_func,json=mainFunc,proto3" json:"main_func,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VMOptions) Reset()         { *m = VMOptions{} }
func (m *VMOptions) String() string { return proto.CompactTextString(m) }
func (*VMOptions) ProtoMessage()    {}
func (*VMOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_brand_141de6d19d21dc5e, []int{1}
}
func (m *VMOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VMOptions.Unmarshal(m, b)
}
func (m *VMOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VMOptions.Marshal(b, m, deterministic)
}
func (dst *VMOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VMOptions.Merge(dst, src)
}
func (m *VMOptions) XXX_Size() int {
	return xxx_messageInfo_VMOptions.Size(m)
}
func (m *VMOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_VMOptions.DiscardUnknown(m)
}

var xxx_messageInfo_VMOptions proto.InternalMessageInfo

func (m *VMOptions) GetEnableJit() bool {
	if m != nil {
		return m.EnableJit
	}
	return false
}

func (m *VMOptions) GetDefaultPages() int32 {
	if m != nil {
		return m.DefaultPages
	}
	return 0
}

func (m *VMOptions) GetMaxPages() int32 {
	if m != nil {
		return m.MaxPages
	}
	return 0
}

func (m *VMOptions) GetMainFunc() string {
	if m != nil {
		return m.MainFunc
	}
	return ""
}

type Metadata struct {
	ExpectedRuntime      string   `protobuf:"bytes,1,opt,name=expected_runtime,json=expectedRuntime,proto3" json:"expected_runtime,omitempty"`
	BinaryKind           Kind     `protobuf:"varint,2,opt,name=binary_kind,json=binaryKind,proto3,enum=within.olin.brand.Kind" json:"binary_kind,omitempty"`
	Author               string   `protobuf:"bytes,3,opt,name=author,proto3" json:"author,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Metadata) Reset()         { *m = Metadata{} }
func (m *Metadata) String() string { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()    {}
func (*Metadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_brand_141de6d19d21dc5e, []int{2}
}
func (m *Metadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metadata.Unmarshal(m, b)
}
func (m *Metadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metadata.Marshal(b, m, deterministic)
}
func (dst *Metadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metadata.Merge(dst, src)
}
func (m *Metadata) XXX_Size() int {
	return xxx_messageInfo_Metadata.Size(m)
}
func (m *Metadata) XXX_DiscardUnknown() {
	xxx_messageInfo_Metadata.DiscardUnknown(m)
}

var xxx_messageInfo_Metadata proto.InternalMessageInfo

func (m *Metadata) GetExpectedRuntime() string {
	if m != nil {
		return m.ExpectedRuntime
	}
	return ""
}

func (m *Metadata) GetBinaryKind() Kind {
	if m != nil {
		return m.BinaryKind
	}
	return Kind_UNKNOWN
}

func (m *Metadata) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *Metadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Brand)(nil), "within.olin.brand.Brand")
	proto.RegisterType((*VMOptions)(nil), "within.olin.brand.VMOptions")
	proto.RegisterType((*Metadata)(nil), "within.olin.brand.Metadata")
	proto.RegisterEnum("within.olin.brand.Kind", Kind_name, Kind_value)
}

func init() { proto.RegisterFile("brand.proto", fileDescriptor_brand_141de6d19d21dc5e) }

var fileDescriptor_brand_141de6d19d21dc5e = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x4f, 0xc2, 0x30,
	0x14, 0xc7, 0x1d, 0x6c, 0x40, 0xdf, 0x54, 0x66, 0x0f, 0x4a, 0x82, 0x26, 0x04, 0x2f, 0xe8, 0x61,
	0x1a, 0xbc, 0x98, 0x78, 0x02, 0x44, 0x83, 0x84, 0x61, 0x1a, 0x84, 0xc4, 0xcb, 0xd2, 0xb1, 0x22,
	0x45, 0xd6, 0x2d, 0xa3, 0x8b, 0xf8, 0x11, 0xfc, 0x0c, 0x7e, 0x59, 0xb3, 0x32, 0x8c, 0x89, 0xdc,
	0x5e, 0xff, 0xef, 0xf7, 0xfa, 0x7e, 0x4d, 0xc1, 0xf4, 0x62, 0x2a, 0x7c, 0x3b, 0x8a, 0x43, 0x19,
	0xe2, 0xa3, 0x0f, 0x2e, 0xe7, 0x5c, 0xd8, 0xe1, 0x92, 0x0b, 0x5b, 0x35, 0xea, 0x0b, 0x30, 0xda,
	0x69, 0x81, 0xaf, 0x41, 0x0f, 0x23, 0xb9, 0xaa, 0x68, 0x35, 0xad, 0x61, 0x36, 0x4f, 0xed, 0x7f,
	0xa8, 0x3d, 0x1e, 0x0c, 0x23, 0xc9, 0x43, 0xb1, 0x22, 0x8a, 0xc4, 0x57, 0xa0, 0x07, 0x4c, 0xd2,
	0x4a, 0x4e, 0x4d, 0x54, 0x77, 0x4c, 0x0c, 0x98, 0xa4, 0x3e, 0x95, 0x94, 0x28, 0xb0, 0xfe, 0xa5,
	0x01, 0xfa, 0xbd, 0x04, 0x9f, 0x01, 0x30, 0x41, 0xbd, 0x25, 0x73, 0x17, 0x5c, 0xaa, 0xb5, 0x25,
	0x82, 0x36, 0xc9, 0x13, 0x97, 0xf8, 0x1c, 0x0e, 0x7c, 0x36, 0xa3, 0xc9, 0x52, 0xba, 0x11, 0x7d,
	0x63, 0x2b, 0xb5, 0xc6, 0x20, 0xfb, 0x59, 0xf8, 0x9c, 0x66, 0xb8, 0x0a, 0x28, 0xa0, 0xeb, 0x0c,
	0xc8, 0x2b, 0xa0, 0x14, 0xd0, 0xf5, 0x9f, 0x26, 0x17, 0xee, 0x2c, 0x11, 0xd3, 0x8a, 0x5e, 0xd3,
	0x1a, 0x28, 0x6d, 0x72, 0xf1, 0x90, 0x88, 0x69, 0xfd, 0x5b, 0x83, 0xd2, 0x56, 0x0f, 0x5f, 0x80,
	0xc5, 0xd6, 0x11, 0x9b, 0x4a, 0xe6, 0xbb, 0x71, 0x22, 0x24, 0x0f, 0x98, 0x12, 0x42, 0xa4, 0xbc,
	0xcd, 0xc9, 0x26, 0xc6, 0xb7, 0x60, 0x7a, 0x5c, 0xd0, 0xf8, 0xd3, 0x7d, 0xe7, 0xc2, 0x57, 0x52,
	0x87, 0xcd, 0x93, 0x1d, 0x6f, 0xef, 0x73, 0xe1, 0x13, 0xd8, 0xb0, 0x69, 0x8d, 0x8f, 0xa1, 0x40,
	0x13, 0x39, 0x0f, 0x63, 0x25, 0x8a, 0x48, 0x76, 0xc2, 0x18, 0x74, 0x41, 0x03, 0x96, 0x19, 0xaa,
	0xfa, 0xf2, 0x0e, 0x74, 0x35, 0x63, 0x42, 0xf1, 0xc5, 0xe9, 0x3b, 0xc3, 0x89, 0x63, 0xed, 0x61,
	0x04, 0x46, 0x67, 0xd2, 0x7a, 0xec, 0x59, 0x1a, 0x06, 0x28, 0xdc, 0xb7, 0xba, 0x83, 0xa1, 0x63,
	0xe5, 0x70, 0x19, 0xcc, 0x9e, 0x33, 0xea, 0x92, 0x56, 0x67, 0xd4, 0x1b, 0x77, 0xad, 0x7c, 0xbb,
	0xf8, 0x6a, 0x28, 0x05, 0xaf, 0xa0, 0x7e, 0xfd, 0xe6, 0x27, 0x00, 0x00, 0xff, 0xff, 0x79, 0xad,
	0x24, 0x3e, 0x04, 0x02, 0x00, 0x00,
}
