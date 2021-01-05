// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: pkg/glowstone/pb/glow.proto

package glowstone

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Tick struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Src     string `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dest    string `protobuf:"bytes,2,opt,name=dest,proto3" json:"dest,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Tick) Reset() {
	*x = Tick{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_glowstone_pb_glow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tick) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tick) ProtoMessage() {}

func (x *Tick) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_glowstone_pb_glow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tick.ProtoReflect.Descriptor instead.
func (*Tick) Descriptor() ([]byte, []int) {
	return file_pkg_glowstone_pb_glow_proto_rawDescGZIP(), []int{0}
}

func (x *Tick) GetSrc() string {
	if x != nil {
		return x.Src
	}
	return ""
}

func (x *Tick) GetDest() string {
	if x != nil {
		return x.Dest
	}
	return ""
}

func (x *Tick) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_pkg_glowstone_pb_glow_proto protoreflect.FileDescriptor

var file_pkg_glowstone_pb_glow_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x2f,
	0x70, 0x62, 0x2f, 0x67, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67,
	0x6c, 0x6f, 0x77, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x22, 0x46, 0x0a, 0x04, 0x54, 0x69, 0x63, 0x6b,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x72, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73,
	0x72, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x64, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x42, 0x0f, 0x5a, 0x0d, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x6c, 0x6f, 0x77, 0x73, 0x74, 0x6f, 0x6e,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_glowstone_pb_glow_proto_rawDescOnce sync.Once
	file_pkg_glowstone_pb_glow_proto_rawDescData = file_pkg_glowstone_pb_glow_proto_rawDesc
)

func file_pkg_glowstone_pb_glow_proto_rawDescGZIP() []byte {
	file_pkg_glowstone_pb_glow_proto_rawDescOnce.Do(func() {
		file_pkg_glowstone_pb_glow_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_glowstone_pb_glow_proto_rawDescData)
	})
	return file_pkg_glowstone_pb_glow_proto_rawDescData
}

var file_pkg_glowstone_pb_glow_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_glowstone_pb_glow_proto_goTypes = []interface{}{
	(*Tick)(nil), // 0: glowstone.Tick
}
var file_pkg_glowstone_pb_glow_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_glowstone_pb_glow_proto_init() }
func file_pkg_glowstone_pb_glow_proto_init() {
	if File_pkg_glowstone_pb_glow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_glowstone_pb_glow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tick); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_glowstone_pb_glow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_glowstone_pb_glow_proto_goTypes,
		DependencyIndexes: file_pkg_glowstone_pb_glow_proto_depIdxs,
		MessageInfos:      file_pkg_glowstone_pb_glow_proto_msgTypes,
	}.Build()
	File_pkg_glowstone_pb_glow_proto = out.File
	file_pkg_glowstone_pb_glow_proto_rawDesc = nil
	file_pkg_glowstone_pb_glow_proto_goTypes = nil
	file_pkg_glowstone_pb_glow_proto_depIdxs = nil
}
