// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: shop_error.proto

package shop

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type ShopAdminErrorReason int32

const (
	ShopAdminErrorReason_UNKNOWN_ERROR ShopAdminErrorReason = 0
	ShopAdminErrorReason_LOGIN_FAILED  ShopAdminErrorReason = 1
)

// Enum value maps for ShopAdminErrorReason.
var (
	ShopAdminErrorReason_name = map[int32]string{
		0: "UNKNOWN_ERROR",
		1: "LOGIN_FAILED",
	}
	ShopAdminErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR": 0,
		"LOGIN_FAILED":  1,
	}
)

func (x ShopAdminErrorReason) Enum() *ShopAdminErrorReason {
	p := new(ShopAdminErrorReason)
	*p = x
	return p
}

func (x ShopAdminErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ShopAdminErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_shop_error_proto_enumTypes[0].Descriptor()
}

func (ShopAdminErrorReason) Type() protoreflect.EnumType {
	return &file_shop_error_proto_enumTypes[0]
}

func (x ShopAdminErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ShopAdminErrorReason.Descriptor instead.
func (ShopAdminErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_shop_error_proto_rawDescGZIP(), []int{0}
}

var File_shop_error_proto protoreflect.FileDescriptor

var file_shop_error_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x68, 0x6f, 0x70, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x1a, 0x13, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2a, 0x47, 0x0a, 0x14, 0x53, 0x68, 0x6f, 0x70, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x0c,
	0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x1a, 0x04,
	0xa8, 0x45, 0xf4, 0x03, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x24, 0x50, 0x01, 0x5a, 0x20,
	0x6d, 0x61, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x72,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x64, 0x2f, 0x73, 0x68, 0x6f, 0x70, 0x3b, 0x73, 0x68, 0x6f, 0x70,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shop_error_proto_rawDescOnce sync.Once
	file_shop_error_proto_rawDescData = file_shop_error_proto_rawDesc
)

func file_shop_error_proto_rawDescGZIP() []byte {
	file_shop_error_proto_rawDescOnce.Do(func() {
		file_shop_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_shop_error_proto_rawDescData)
	})
	return file_shop_error_proto_rawDescData
}

var file_shop_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shop_error_proto_goTypes = []interface{}{
	(ShopAdminErrorReason)(0), // 0: api.shop.ShopAdminErrorReason
}
var file_shop_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shop_error_proto_init() }
func file_shop_error_proto_init() {
	if File_shop_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shop_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shop_error_proto_goTypes,
		DependencyIndexes: file_shop_error_proto_depIdxs,
		EnumInfos:         file_shop_error_proto_enumTypes,
	}.Build()
	File_shop_error_proto = out.File
	file_shop_error_proto_rawDesc = nil
	file_shop_error_proto_goTypes = nil
	file_shop_error_proto_depIdxs = nil
}
