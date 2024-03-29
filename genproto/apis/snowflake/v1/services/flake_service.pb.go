// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: snowflake/v1/services/flake_service.proto

package services

import (
	resources "github.com/prettykingking/snowflake/genproto/apis/snowflake/v1/resources"
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

type GetFlakeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetFlakeRequest) Reset() {
	*x = GetFlakeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snowflake_v1_services_flake_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlakeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlakeRequest) ProtoMessage() {}

func (x *GetFlakeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snowflake_v1_services_flake_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlakeRequest.ProtoReflect.Descriptor instead.
func (*GetFlakeRequest) Descriptor() ([]byte, []int) {
	return file_snowflake_v1_services_flake_service_proto_rawDescGZIP(), []int{0}
}

var File_snowflake_v1_services_flake_service_proto protoreflect.FileDescriptor

var file_snowflake_v1_services_flake_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x73, 0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73, 0x6e, 0x6f,
	0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x1a, 0x22, 0x73, 0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x66, 0x6c, 0x61, 0x6b, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61,
	0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x32, 0x63, 0x0a, 0x0c, 0x46, 0x6c, 0x61,
	0x6b, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x08, 0x47, 0x65, 0x74,
	0x46, 0x6c, 0x61, 0x6b, 0x65, 0x12, 0x26, 0x2e, 0x73, 0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x47, 0x65,
	0x74, 0x46, 0x6c, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x73, 0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2e, 0x46, 0x6c, 0x61, 0x6b, 0x65, 0x22, 0x00, 0x42, 0x49,
	0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x65,
	0x74, 0x74, 0x79, 0x6b, 0x69, 0x6e, 0x67, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x73, 0x6e, 0x6f, 0x77,
	0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x73, 0x6e, 0x6f, 0x77, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_snowflake_v1_services_flake_service_proto_rawDescOnce sync.Once
	file_snowflake_v1_services_flake_service_proto_rawDescData = file_snowflake_v1_services_flake_service_proto_rawDesc
)

func file_snowflake_v1_services_flake_service_proto_rawDescGZIP() []byte {
	file_snowflake_v1_services_flake_service_proto_rawDescOnce.Do(func() {
		file_snowflake_v1_services_flake_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_snowflake_v1_services_flake_service_proto_rawDescData)
	})
	return file_snowflake_v1_services_flake_service_proto_rawDescData
}

var file_snowflake_v1_services_flake_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_snowflake_v1_services_flake_service_proto_goTypes = []interface{}{
	(*GetFlakeRequest)(nil), // 0: snowflake.v1.services.GetFlakeRequest
	(*resources.Flake)(nil), // 1: snowflake.v1.resources.Flake
}
var file_snowflake_v1_services_flake_service_proto_depIdxs = []int32{
	0, // 0: snowflake.v1.services.FlakeService.GetFlake:input_type -> snowflake.v1.services.GetFlakeRequest
	1, // 1: snowflake.v1.services.FlakeService.GetFlake:output_type -> snowflake.v1.resources.Flake
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_snowflake_v1_services_flake_service_proto_init() }
func file_snowflake_v1_services_flake_service_proto_init() {
	if File_snowflake_v1_services_flake_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_snowflake_v1_services_flake_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlakeRequest); i {
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
			RawDescriptor: file_snowflake_v1_services_flake_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_snowflake_v1_services_flake_service_proto_goTypes,
		DependencyIndexes: file_snowflake_v1_services_flake_service_proto_depIdxs,
		MessageInfos:      file_snowflake_v1_services_flake_service_proto_msgTypes,
	}.Build()
	File_snowflake_v1_services_flake_service_proto = out.File
	file_snowflake_v1_services_flake_service_proto_rawDesc = nil
	file_snowflake_v1_services_flake_service_proto_goTypes = nil
	file_snowflake_v1_services_flake_service_proto_depIdxs = nil
}
