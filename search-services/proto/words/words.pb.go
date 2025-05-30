// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v4.25.3
// source: proto/words/words.proto

package words

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type WordsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Phrase        string                 `protobuf:"bytes,1,opt,name=phrase,proto3" json:"phrase,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WordsRequest) Reset() {
	*x = WordsRequest{}
	mi := &file_proto_words_words_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WordsRequest) ProtoMessage() {}

func (x *WordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_words_words_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WordsRequest.ProtoReflect.Descriptor instead.
func (*WordsRequest) Descriptor() ([]byte, []int) {
	return file_proto_words_words_proto_rawDescGZIP(), []int{0}
}

func (x *WordsRequest) GetPhrase() string {
	if x != nil {
		return x.Phrase
	}
	return ""
}

type WordsReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Words         []string               `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WordsReply) Reset() {
	*x = WordsReply{}
	mi := &file_proto_words_words_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WordsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WordsReply) ProtoMessage() {}

func (x *WordsReply) ProtoReflect() protoreflect.Message {
	mi := &file_proto_words_words_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WordsReply.ProtoReflect.Descriptor instead.
func (*WordsReply) Descriptor() ([]byte, []int) {
	return file_proto_words_words_proto_rawDescGZIP(), []int{1}
}

func (x *WordsReply) GetWords() []string {
	if x != nil {
		return x.Words
	}
	return nil
}

var File_proto_words_words_proto protoreflect.FileDescriptor

var file_proto_words_words_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2f, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a,
	0x0c, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70,
	0x68, 0x72, 0x61, 0x73, 0x65, 0x22, 0x22, 0x0a, 0x0a, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x32, 0x73, 0x0a, 0x05, 0x57, 0x6f, 0x72,
	0x64, 0x73, 0x12, 0x38, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x04,
	0x4e, 0x6f, 0x72, 0x6d, 0x12, 0x13, 0x2e, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x2e, 0x57, 0x6f, 0x72,
	0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x2e, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x1e,
	0x5a, 0x1c, 0x79, 0x61, 0x64, 0x72, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x75, 0x72,
	0x73, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_words_words_proto_rawDescOnce sync.Once
	file_proto_words_words_proto_rawDescData []byte
)

func file_proto_words_words_proto_rawDescGZIP() []byte {
	file_proto_words_words_proto_rawDescOnce.Do(func() {
		file_proto_words_words_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_words_words_proto_rawDesc), len(file_proto_words_words_proto_rawDesc)))
	})
	return file_proto_words_words_proto_rawDescData
}

var file_proto_words_words_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_words_words_proto_goTypes = []any{
	(*WordsRequest)(nil),  // 0: words.WordsRequest
	(*WordsReply)(nil),    // 1: words.WordsReply
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_proto_words_words_proto_depIdxs = []int32{
	2, // 0: words.Words.Ping:input_type -> google.protobuf.Empty
	0, // 1: words.Words.Norm:input_type -> words.WordsRequest
	2, // 2: words.Words.Ping:output_type -> google.protobuf.Empty
	1, // 3: words.Words.Norm:output_type -> words.WordsReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_words_words_proto_init() }
func file_proto_words_words_proto_init() {
	if File_proto_words_words_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_words_words_proto_rawDesc), len(file_proto_words_words_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_words_words_proto_goTypes,
		DependencyIndexes: file_proto_words_words_proto_depIdxs,
		MessageInfos:      file_proto_words_words_proto_msgTypes,
	}.Build()
	File_proto_words_words_proto = out.File
	file_proto_words_words_proto_goTypes = nil
	file_proto_words_words_proto_depIdxs = nil
}
