// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: proto/common.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IncrementDecrement int32

const (
	IncrementDecrement_INCREMENT IncrementDecrement = 0
	IncrementDecrement_DECREMENT IncrementDecrement = 1
)

// Enum value maps for IncrementDecrement.
var (
	IncrementDecrement_name = map[int32]string{
		0: "INCREMENT",
		1: "DECREMENT",
	}
	IncrementDecrement_value = map[string]int32{
		"INCREMENT": 0,
		"DECREMENT": 1,
	}
)

func (x IncrementDecrement) Enum() *IncrementDecrement {
	p := new(IncrementDecrement)
	*p = x
	return p
}

func (x IncrementDecrement) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IncrementDecrement) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_common_proto_enumTypes[0].Descriptor()
}

func (IncrementDecrement) Type() protoreflect.EnumType {
	return &file_proto_common_proto_enumTypes[0]
}

func (x IncrementDecrement) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IncrementDecrement.Descriptor instead.
func (IncrementDecrement) EnumDescriptor() ([]byte, []int) {
	return file_proto_common_proto_rawDescGZIP(), []int{0}
}

type PressType int32

const (
	PressType_TAP     PressType = 0
	PressType_HOLD    PressType = 1
	PressType_RELEASE PressType = 2
)

// Enum value maps for PressType.
var (
	PressType_name = map[int32]string{
		0: "TAP",
		1: "HOLD",
		2: "RELEASE",
	}
	PressType_value = map[string]int32{
		"TAP":     0,
		"HOLD":    1,
		"RELEASE": 2,
	}
)

func (x PressType) Enum() *PressType {
	p := new(PressType)
	*p = x
	return p
}

func (x PressType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PressType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_common_proto_enumTypes[1].Descriptor()
}

func (PressType) Type() protoreflect.EnumType {
	return &file_proto_common_proto_enumTypes[1]
}

func (x PressType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PressType.Descriptor instead.
func (PressType) EnumDescriptor() ([]byte, []int) {
	return file_proto_common_proto_rawDescGZIP(), []int{1}
}

type Color int32

const (
	Color_RED     Color = 0
	Color_BLUE    Color = 1
	Color_WHITE   Color = 2
	Color_BLACK   Color = 3
	Color_YELLOW  Color = 4
	Color_GREEN   Color = 5
	Color_ORANGE  Color = 6
	Color_PINK    Color = 7
	Color_UNKNOWN Color = 99
)

// Enum value maps for Color.
var (
	Color_name = map[int32]string{
		0:  "RED",
		1:  "BLUE",
		2:  "WHITE",
		3:  "BLACK",
		4:  "YELLOW",
		5:  "GREEN",
		6:  "ORANGE",
		7:  "PINK",
		99: "UNKNOWN",
	}
	Color_value = map[string]int32{
		"RED":     0,
		"BLUE":    1,
		"WHITE":   2,
		"BLACK":   3,
		"YELLOW":  4,
		"GREEN":   5,
		"ORANGE":  6,
		"PINK":    7,
		"UNKNOWN": 99,
	}
)

func (x Color) Enum() *Color {
	p := new(Color)
	*p = x
	return p
}

func (x Color) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Color) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_common_proto_enumTypes[2].Descriptor()
}

func (Color) Type() protoreflect.EnumType {
	return &file_proto_common_proto_enumTypes[2]
}

func (x Color) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Color.Descriptor instead.
func (Color) EnumDescriptor() ([]byte, []int) {
	return file_proto_common_proto_rawDescGZIP(), []int{2}
}

var File_proto_common_proto protoreflect.FileDescriptor

const file_proto_common_proto_rawDesc = "" +
	"\n" +
	"\x12proto/common.proto\x12\x06common*2\n" +
	"\x12IncrementDecrement\x12\r\n" +
	"\tINCREMENT\x10\x00\x12\r\n" +
	"\tDECREMENT\x10\x01*+\n" +
	"\tPressType\x12\a\n" +
	"\x03TAP\x10\x00\x12\b\n" +
	"\x04HOLD\x10\x01\x12\v\n" +
	"\aRELEASE\x10\x02*j\n" +
	"\x05Color\x12\a\n" +
	"\x03RED\x10\x00\x12\b\n" +
	"\x04BLUE\x10\x01\x12\t\n" +
	"\x05WHITE\x10\x02\x12\t\n" +
	"\x05BLACK\x10\x03\x12\n" +
	"\n" +
	"\x06YELLOW\x10\x04\x12\t\n" +
	"\x05GREEN\x10\x05\x12\n" +
	"\n" +
	"\x06ORANGE\x10\x06\x12\b\n" +
	"\x04PINK\x10\a\x12\v\n" +
	"\aUNKNOWN\x10cB\tZ\a./protob\x06proto3"

var (
	file_proto_common_proto_rawDescOnce sync.Once
	file_proto_common_proto_rawDescData []byte
)

func file_proto_common_proto_rawDescGZIP() []byte {
	file_proto_common_proto_rawDescOnce.Do(func() {
		file_proto_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_common_proto_rawDesc), len(file_proto_common_proto_rawDesc)))
	})
	return file_proto_common_proto_rawDescData
}

var file_proto_common_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_common_proto_goTypes = []any{
	(IncrementDecrement)(0), // 0: common.IncrementDecrement
	(PressType)(0),          // 1: common.PressType
	(Color)(0),              // 2: common.Color
}
var file_proto_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_common_proto_init() }
func file_proto_common_proto_init() {
	if File_proto_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_common_proto_rawDesc), len(file_proto_common_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_common_proto_goTypes,
		DependencyIndexes: file_proto_common_proto_depIdxs,
		EnumInfos:         file_proto_common_proto_enumTypes,
	}.Build()
	File_proto_common_proto = out.File
	file_proto_common_proto_goTypes = nil
	file_proto_common_proto_depIdxs = nil
}
