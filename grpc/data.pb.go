// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/data.proto

package grpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FormType int32

const (
	FormType_SQL FormType = 0
	FormType_Doc FormType = 1
)

var FormType_name = map[int32]string{
	0: "SQL",
	1: "Doc",
}

var FormType_value = map[string]int32{
	"SQL": 0,
	"Doc": 1,
}

func (x FormType) String() string {
	return proto.EnumName(FormType_name, int32(x))
}

func (FormType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_86020f0d45cc739c, []int{0}
}

// Lily 数据库引擎对象
type Lily struct {
	Databases            map[string]*Database `protobuf:"bytes,1,rep,name=databases,proto3" json:"databases,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Lily) Reset()         { *m = Lily{} }
func (m *Lily) String() string { return proto.CompactTextString(m) }
func (*Lily) ProtoMessage()    {}
func (*Lily) Descriptor() ([]byte, []int) {
	return fileDescriptor_86020f0d45cc739c, []int{0}
}

func (m *Lily) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lily.Unmarshal(m, b)
}
func (m *Lily) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lily.Marshal(b, m, deterministic)
}
func (m *Lily) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lily.Merge(m, src)
}
func (m *Lily) XXX_Size() int {
	return xxx_messageInfo_Lily.Size(m)
}
func (m *Lily) XXX_DiscardUnknown() {
	xxx_messageInfo_Lily.DiscardUnknown(m)
}

var xxx_messageInfo_Lily proto.InternalMessageInfo

func (m *Lily) GetDatabases() map[string]*Database {
	if m != nil {
		return m.Databases
	}
	return nil
}

// Database 数据库对象
type Database struct {
	Id                   string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Comment              string           `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
	Forms                map[string]*Form `protobuf:"bytes,4,rep,name=forms,proto3" json:"forms,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Database) Reset()         { *m = Database{} }
func (m *Database) String() string { return proto.CompactTextString(m) }
func (*Database) ProtoMessage()    {}
func (*Database) Descriptor() ([]byte, []int) {
	return fileDescriptor_86020f0d45cc739c, []int{1}
}

func (m *Database) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Database.Unmarshal(m, b)
}
func (m *Database) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Database.Marshal(b, m, deterministic)
}
func (m *Database) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Database.Merge(m, src)
}
func (m *Database) XXX_Size() int {
	return xxx_messageInfo_Database.Size(m)
}
func (m *Database) XXX_DiscardUnknown() {
	xxx_messageInfo_Database.DiscardUnknown(m)
}

var xxx_messageInfo_Database proto.InternalMessageInfo

func (m *Database) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Database) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Database) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *Database) GetForms() map[string]*Form {
	if m != nil {
		return m.Forms
	}
	return nil
}

// Form 数据库表对象
type Form struct {
	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Comment              string            `protobuf:"bytes,3,opt,name=comment,proto3" json:"comment,omitempty"`
	FormType             FormType          `protobuf:"varint,4,opt,name=formType,proto3,enum=grpc.FormType" json:"formType,omitempty"`
	Indexes              map[string]*Index `protobuf:"bytes,5,rep,name=indexes,proto3" json:"indexes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Form) Reset()         { *m = Form{} }
func (m *Form) String() string { return proto.CompactTextString(m) }
func (*Form) ProtoMessage()    {}
func (*Form) Descriptor() ([]byte, []int) {
	return fileDescriptor_86020f0d45cc739c, []int{2}
}

func (m *Form) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Form.Unmarshal(m, b)
}
func (m *Form) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Form.Marshal(b, m, deterministic)
}
func (m *Form) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Form.Merge(m, src)
}
func (m *Form) XXX_Size() int {
	return xxx_messageInfo_Form.Size(m)
}
func (m *Form) XXX_DiscardUnknown() {
	xxx_messageInfo_Form.DiscardUnknown(m)
}

var xxx_messageInfo_Form proto.InternalMessageInfo

func (m *Form) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Form) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Form) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *Form) GetFormType() FormType {
	if m != nil {
		return m.FormType
	}
	return FormType_SQL
}

func (m *Form) GetIndexes() map[string]*Index {
	if m != nil {
		return m.Indexes
	}
	return nil
}

// Index 索引对象
type Index struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Primary              bool     `protobuf:"varint,2,opt,name=primary,proto3" json:"primary,omitempty"`
	KeyStructure         string   `protobuf:"bytes,3,opt,name=keyStructure,proto3" json:"keyStructure,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Index) Reset()         { *m = Index{} }
func (m *Index) String() string { return proto.CompactTextString(m) }
func (*Index) ProtoMessage()    {}
func (*Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_86020f0d45cc739c, []int{3}
}

func (m *Index) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Index.Unmarshal(m, b)
}
func (m *Index) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Index.Marshal(b, m, deterministic)
}
func (m *Index) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Index.Merge(m, src)
}
func (m *Index) XXX_Size() int {
	return xxx_messageInfo_Index.Size(m)
}
func (m *Index) XXX_DiscardUnknown() {
	xxx_messageInfo_Index.DiscardUnknown(m)
}

var xxx_messageInfo_Index proto.InternalMessageInfo

func (m *Index) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Index) GetPrimary() bool {
	if m != nil {
		return m.Primary
	}
	return false
}

func (m *Index) GetKeyStructure() string {
	if m != nil {
		return m.KeyStructure
	}
	return ""
}

func init() {
	proto.RegisterEnum("grpc.FormType", FormType_name, FormType_value)
	proto.RegisterType((*Lily)(nil), "grpc.Lily")
	proto.RegisterMapType((map[string]*Database)(nil), "grpc.Lily.DatabasesEntry")
	proto.RegisterType((*Database)(nil), "grpc.Database")
	proto.RegisterMapType((map[string]*Form)(nil), "grpc.Database.FormsEntry")
	proto.RegisterType((*Form)(nil), "grpc.Form")
	proto.RegisterMapType((map[string]*Index)(nil), "grpc.Form.IndexesEntry")
	proto.RegisterType((*Index)(nil), "grpc.Index")
}

func init() { proto.RegisterFile("grpc/data.proto", fileDescriptor_86020f0d45cc739c) }

var fileDescriptor_86020f0d45cc739c = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x4d, 0x4e, 0xe3, 0x40,
	0x10, 0x85, 0xa7, 0xfd, 0x33, 0xb6, 0x2b, 0x91, 0x27, 0xea, 0x59, 0x4c, 0x4f, 0x34, 0x0b, 0x4f,
	0xc4, 0xc2, 0xca, 0xa2, 0x23, 0xc2, 0x02, 0xc4, 0x32, 0x0a, 0x41, 0x48, 0x59, 0x80, 0x03, 0x07,
	0x70, 0xec, 0x06, 0x59, 0x89, 0x7f, 0xd4, 0x76, 0x10, 0x3e, 0x04, 0x67, 0xe2, 0x42, 0x1c, 0x02,
	0xb5, 0x3b, 0x4d, 0x6c, 0x91, 0x1d, 0xbb, 0xea, 0xf7, 0xaa, 0x9e, 0xea, 0xb3, 0x0b, 0x7e, 0x3d,
	0xf1, 0x22, 0x9a, 0xc4, 0x61, 0x15, 0xd2, 0x82, 0xe7, 0x55, 0x8e, 0x0d, 0x21, 0x8c, 0x5e, 0x11,
	0x18, 0xcb, 0x64, 0x5b, 0xe3, 0x73, 0x70, 0x84, 0xb9, 0x0e, 0x4b, 0x56, 0x12, 0xe4, 0xe9, 0x7e,
	0x6f, 0xfa, 0x97, 0x8a, 0x16, 0x2a, 0x6c, 0x3a, 0x57, 0xde, 0x55, 0x56, 0xf1, 0x3a, 0x38, 0xf4,
	0x0e, 0x97, 0xe0, 0x76, 0x4d, 0x3c, 0x00, 0x7d, 0xc3, 0x6a, 0x82, 0x3c, 0xe4, 0x3b, 0x81, 0x28,
	0xf1, 0x09, 0x98, 0xcf, 0xe1, 0x76, 0xc7, 0x88, 0xe6, 0x21, 0xbf, 0x37, 0x75, 0x65, 0xb0, 0x1a,
	0x0b, 0xa4, 0x79, 0xa9, 0x5d, 0xa0, 0xd1, 0x1b, 0x02, 0x5b, 0xe9, 0xd8, 0x05, 0x2d, 0x89, 0xf7,
	0x39, 0x5a, 0x12, 0x63, 0x0c, 0x46, 0x16, 0xa6, 0x32, 0xc5, 0x09, 0x9a, 0x1a, 0x13, 0xb0, 0xa2,
	0x3c, 0x4d, 0x59, 0x56, 0x11, 0xbd, 0x91, 0xd5, 0x13, 0x4f, 0xc0, 0x7c, 0xcc, 0x79, 0x5a, 0x12,
	0xa3, 0x4d, 0xa3, 0xc2, 0xe9, 0x42, 0x78, 0x92, 0x46, 0xf6, 0x0d, 0xe7, 0x00, 0x07, 0xf1, 0x08,
	0x85, 0xd7, 0xa5, 0x00, 0x19, 0x28, 0x46, 0xda, 0x04, 0xef, 0x08, 0x0c, 0xa1, 0x7d, 0x73, 0xfb,
	0x31, 0xd8, 0x62, 0xab, 0xfb, 0xba, 0x60, 0xc4, 0xf0, 0x90, 0xef, 0xaa, 0xaf, 0xb6, 0xd8, 0xab,
	0xc1, 0xa7, 0x8f, 0x4f, 0xc1, 0x4a, 0xb2, 0x98, 0xbd, 0xb0, 0x92, 0x98, 0x0d, 0xeb, 0x9f, 0x43,
	0x2b, 0xbd, 0x91, 0x8e, 0x24, 0x55, 0x7d, 0xc3, 0x6b, 0xe8, 0xb7, 0x8d, 0x23, 0xb4, 0xff, 0xbb,
	0xb4, 0x3d, 0x19, 0xd9, 0x0c, 0xb5, 0x71, 0x1f, 0xc0, 0x6c, 0xb4, 0x2f, 0xb8, 0x04, 0xac, 0x82,
	0x27, 0x69, 0xc8, 0xeb, 0x26, 0xc1, 0x0e, 0xd4, 0x13, 0x8f, 0xa0, 0xbf, 0x61, 0xf5, 0xaa, 0xe2,
	0xbb, 0xa8, 0xda, 0x71, 0xb6, 0x27, 0xef, 0x68, 0xe3, 0x7f, 0x60, 0x2b, 0x50, 0x6c, 0x81, 0xbe,
	0xba, 0x5b, 0x0e, 0x7e, 0x88, 0x62, 0x9e, 0x47, 0x03, 0x34, 0xf3, 0xe0, 0x77, 0x94, 0xd1, 0x70,
	0xcd, 0x78, 0x12, 0xd1, 0xad, 0xb8, 0x4f, 0xb1, 0xdc, 0xcc, 0x11, 0x3f, 0xf7, 0x56, 0x5c, 0xf7,
	0xfa, 0x67, 0x73, 0xe4, 0x67, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa2, 0xa5, 0x89, 0x5e, 0xf7,
	0x02, 0x00, 0x00,
}
