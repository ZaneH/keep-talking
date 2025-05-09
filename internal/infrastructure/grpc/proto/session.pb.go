// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: proto/session.proto

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

type GameSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionId     string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameSessionRequest) Reset() {
	*x = GameSessionRequest{}
	mi := &file_proto_session_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameSessionRequest) ProtoMessage() {}

func (x *GameSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_session_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameSessionRequest.ProtoReflect.Descriptor instead.
func (*GameSessionRequest) Descriptor() ([]byte, []int) {
	return file_proto_session_proto_rawDescGZIP(), []int{0}
}

func (x *GameSessionRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type GameSession struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionId     string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	StrikeCount   int32                  `protobuf:"varint,2,opt,name=strike_count,json=strikeCount,proto3" json:"strike_count,omitempty"`
	MaxStrikes    int32                  `protobuf:"varint,3,opt,name=max_strikes,json=maxStrikes,proto3" json:"max_strikes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GameSession) Reset() {
	*x = GameSession{}
	mi := &file_proto_session_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GameSession) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameSession) ProtoMessage() {}

func (x *GameSession) ProtoReflect() protoreflect.Message {
	mi := &file_proto_session_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameSession.ProtoReflect.Descriptor instead.
func (*GameSession) Descriptor() ([]byte, []int) {
	return file_proto_session_proto_rawDescGZIP(), []int{1}
}

func (x *GameSession) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *GameSession) GetStrikeCount() int32 {
	if x != nil {
		return x.StrikeCount
	}
	return 0
}

func (x *GameSession) GetMaxStrikes() int32 {
	if x != nil {
		return x.MaxStrikes
	}
	return 0
}

type GetBombsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionId     string                 `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBombsRequest) Reset() {
	*x = GetBombsRequest{}
	mi := &file_proto_session_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBombsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBombsRequest) ProtoMessage() {}

func (x *GetBombsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_session_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBombsRequest.ProtoReflect.Descriptor instead.
func (*GetBombsRequest) Descriptor() ([]byte, []int) {
	return file_proto_session_proto_rawDescGZIP(), []int{2}
}

func (x *GetBombsRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type GetBombsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Bombs         []*Bomb                `protobuf:"bytes,1,rep,name=bombs,proto3" json:"bombs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBombsResponse) Reset() {
	*x = GetBombsResponse{}
	mi := &file_proto_session_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBombsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBombsResponse) ProtoMessage() {}

func (x *GetBombsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_session_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBombsResponse.ProtoReflect.Descriptor instead.
func (*GetBombsResponse) Descriptor() ([]byte, []int) {
	return file_proto_session_proto_rawDescGZIP(), []int{3}
}

func (x *GetBombsResponse) GetBombs() []*Bomb {
	if x != nil {
		return x.Bombs
	}
	return nil
}

var File_proto_session_proto protoreflect.FileDescriptor

const file_proto_session_proto_rawDesc = "" +
	"\n" +
	"\x13proto/session.proto\x12\asession\x1a\x10proto/bomb.proto\"3\n" +
	"\x12GameSessionRequest\x12\x1d\n" +
	"\n" +
	"session_id\x18\x01 \x01(\tR\tsessionId\"p\n" +
	"\vGameSession\x12\x1d\n" +
	"\n" +
	"session_id\x18\x01 \x01(\tR\tsessionId\x12!\n" +
	"\fstrike_count\x18\x02 \x01(\x05R\vstrikeCount\x12\x1f\n" +
	"\vmax_strikes\x18\x03 \x01(\x05R\n" +
	"maxStrikes\"0\n" +
	"\x0fGetBombsRequest\x12\x1d\n" +
	"\n" +
	"session_id\x18\x01 \x01(\tR\tsessionId\"4\n" +
	"\x10GetBombsResponse\x12 \n" +
	"\x05bombs\x18\x01 \x03(\v2\n" +
	".bomb.BombR\x05bombsB\tZ\a./protob\x06proto3"

var (
	file_proto_session_proto_rawDescOnce sync.Once
	file_proto_session_proto_rawDescData []byte
)

func file_proto_session_proto_rawDescGZIP() []byte {
	file_proto_session_proto_rawDescOnce.Do(func() {
		file_proto_session_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_session_proto_rawDesc), len(file_proto_session_proto_rawDesc)))
	})
	return file_proto_session_proto_rawDescData
}

var file_proto_session_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_session_proto_goTypes = []any{
	(*GameSessionRequest)(nil), // 0: session.GameSessionRequest
	(*GameSession)(nil),        // 1: session.GameSession
	(*GetBombsRequest)(nil),    // 2: session.GetBombsRequest
	(*GetBombsResponse)(nil),   // 3: session.GetBombsResponse
	(*Bomb)(nil),               // 4: bomb.Bomb
}
var file_proto_session_proto_depIdxs = []int32{
	4, // 0: session.GetBombsResponse.bombs:type_name -> bomb.Bomb
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_session_proto_init() }
func file_proto_session_proto_init() {
	if File_proto_session_proto != nil {
		return
	}
	file_proto_bomb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_session_proto_rawDesc), len(file_proto_session_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_session_proto_goTypes,
		DependencyIndexes: file_proto_session_proto_depIdxs,
		MessageInfos:      file_proto_session_proto_msgTypes,
	}.Build()
	File_proto_session_proto = out.File
	file_proto_session_proto_goTypes = nil
	file_proto_session_proto_depIdxs = nil
}
