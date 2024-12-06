// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v5.26.1
// source: envoy/extensions/key_value/file_based/v3/config.proto

package file_basedv3

import (
	_ "github.com/cncf/xds/go/udpa/annotations"
	_ "github.com/cncf/xds/go/xds/annotations/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

// [#extension: envoy.key_value.file_based]
// This is configuration to flush a key value store out to disk.
type FileBasedKeyValueStoreConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The filename to read the keys and values from, and write the keys and
	// values to.
	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	// The interval at which the key value store should be flushed to the file.
	FlushInterval *durationpb.Duration `protobuf:"bytes,2,opt,name=flush_interval,json=flushInterval,proto3" json:"flush_interval,omitempty"`
	// The maximum number of entries to cache, or 0 to allow for unlimited entries.
	// Defaults to 1000 if not present.
	MaxEntries *wrapperspb.UInt32Value `protobuf:"bytes,3,opt,name=max_entries,json=maxEntries,proto3" json:"max_entries,omitempty"`
}

func (x *FileBasedKeyValueStoreConfig) Reset() {
	*x = FileBasedKeyValueStoreConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_key_value_file_based_v3_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileBasedKeyValueStoreConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileBasedKeyValueStoreConfig) ProtoMessage() {}

func (x *FileBasedKeyValueStoreConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_key_value_file_based_v3_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileBasedKeyValueStoreConfig.ProtoReflect.Descriptor instead.
func (*FileBasedKeyValueStoreConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescGZIP(), []int{0}
}

func (x *FileBasedKeyValueStoreConfig) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *FileBasedKeyValueStoreConfig) GetFlushInterval() *durationpb.Duration {
	if x != nil {
		return x.FlushInterval
	}
	return nil
}

func (x *FileBasedKeyValueStoreConfig) GetMaxEntries() *wrapperspb.UInt32Value {
	if x != nil {
		return x.MaxEntries
	}
	return nil
}

var File_envoy_extensions_key_value_file_based_v3_config_proto protoreflect.FileDescriptor

var file_envoy_extensions_key_value_file_based_v3_config_proto_rawDesc = []byte{
	0x0a, 0x35, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2f, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64, 0x2f, 0x76, 0x33, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64, 0x2e, 0x76,
	0x33, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x78, 0x64, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x76, 0x33, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x01, 0x0a, 0x1c, 0x46,
	0x69, 0x6c, 0x65, 0x42, 0x61, 0x73, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x23, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x40, 0x0a, 0x0e, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x66, 0x6c, 0x75, 0x73, 0x68, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76,
	0x61, 0x6c, 0x12, 0x3d, 0x0a, 0x0b, 0x6d, 0x61, 0x78, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x6d, 0x61, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x3a, 0x08, 0xd2, 0xc6, 0xa4, 0xe1, 0x06, 0x02, 0x08, 0x01, 0x42, 0xad, 0x01, 0x0a, 0x36,
	0x69, 0x6f, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65, 0x6e,
	0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6b,
	0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x62, 0x61,
	0x73, 0x65, 0x64, 0x2e, 0x76, 0x33, 0x42, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e,
	0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6b,
	0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x62, 0x61,
	0x73, 0x65, 0x64, 0x2f, 0x76, 0x33, 0x3b, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x62, 0x61, 0x73, 0x65,
	0x64, 0x76, 0x33, 0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02, 0x10, 0x02, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescOnce sync.Once
	file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescData = file_envoy_extensions_key_value_file_based_v3_config_proto_rawDesc
)

func file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescGZIP() []byte {
	file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescOnce.Do(func() {
		file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescData)
	})
	return file_envoy_extensions_key_value_file_based_v3_config_proto_rawDescData
}

var file_envoy_extensions_key_value_file_based_v3_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_envoy_extensions_key_value_file_based_v3_config_proto_goTypes = []interface{}{
	(*FileBasedKeyValueStoreConfig)(nil), // 0: envoy.extensions.key_value.file_based.v3.FileBasedKeyValueStoreConfig
	(*durationpb.Duration)(nil),          // 1: google.protobuf.Duration
	(*wrapperspb.UInt32Value)(nil),       // 2: google.protobuf.UInt32Value
}
var file_envoy_extensions_key_value_file_based_v3_config_proto_depIdxs = []int32{
	1, // 0: envoy.extensions.key_value.file_based.v3.FileBasedKeyValueStoreConfig.flush_interval:type_name -> google.protobuf.Duration
	2, // 1: envoy.extensions.key_value.file_based.v3.FileBasedKeyValueStoreConfig.max_entries:type_name -> google.protobuf.UInt32Value
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_envoy_extensions_key_value_file_based_v3_config_proto_init() }
func file_envoy_extensions_key_value_file_based_v3_config_proto_init() {
	if File_envoy_extensions_key_value_file_based_v3_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_extensions_key_value_file_based_v3_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileBasedKeyValueStoreConfig); i {
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
			RawDescriptor: file_envoy_extensions_key_value_file_based_v3_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_extensions_key_value_file_based_v3_config_proto_goTypes,
		DependencyIndexes: file_envoy_extensions_key_value_file_based_v3_config_proto_depIdxs,
		MessageInfos:      file_envoy_extensions_key_value_file_based_v3_config_proto_msgTypes,
	}.Build()
	File_envoy_extensions_key_value_file_based_v3_config_proto = out.File
	file_envoy_extensions_key_value_file_based_v3_config_proto_rawDesc = nil
	file_envoy_extensions_key_value_file_based_v3_config_proto_goTypes = nil
	file_envoy_extensions_key_value_file_based_v3_config_proto_depIdxs = nil
}
