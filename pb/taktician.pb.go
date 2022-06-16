// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: tak/proto/taktician.proto

package pb

import (
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

type AnalyzeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Position string `protobuf:"bytes,1,opt,name=position,proto3" json:"position,omitempty"`
	Depth    int32  `protobuf:"varint,2,opt,name=depth,proto3" json:"depth,omitempty"`
	Precise  bool   `protobuf:"varint,3,opt,name=precise,proto3" json:"precise,omitempty"`
}

func (x *AnalyzeRequest) Reset() {
	*x = AnalyzeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeRequest) ProtoMessage() {}

func (x *AnalyzeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeRequest.ProtoReflect.Descriptor instead.
func (*AnalyzeRequest) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{0}
}

func (x *AnalyzeRequest) GetPosition() string {
	if x != nil {
		return x.Position
	}
	return ""
}

func (x *AnalyzeRequest) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

func (x *AnalyzeRequest) GetPrecise() bool {
	if x != nil {
		return x.Precise
	}
	return false
}

type AnalyzeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pv    []string `protobuf:"bytes,1,rep,name=pv,proto3" json:"pv,omitempty"`
	Value int64    `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Depth int32    `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *AnalyzeResponse) Reset() {
	*x = AnalyzeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeResponse) ProtoMessage() {}

func (x *AnalyzeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeResponse.ProtoReflect.Descriptor instead.
func (*AnalyzeResponse) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{1}
}

func (x *AnalyzeResponse) GetPv() []string {
	if x != nil {
		return x.Pv
	}
	return nil
}

func (x *AnalyzeResponse) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *AnalyzeResponse) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type CanonicalizeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size  int32    `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	Moves []string `protobuf:"bytes,2,rep,name=moves,proto3" json:"moves,omitempty"`
}

func (x *CanonicalizeRequest) Reset() {
	*x = CanonicalizeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalizeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalizeRequest) ProtoMessage() {}

func (x *CanonicalizeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalizeRequest.ProtoReflect.Descriptor instead.
func (*CanonicalizeRequest) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{2}
}

func (x *CanonicalizeRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *CanonicalizeRequest) GetMoves() []string {
	if x != nil {
		return x.Moves
	}
	return nil
}

type CanonicalizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Moves []string `protobuf:"bytes,1,rep,name=moves,proto3" json:"moves,omitempty"`
}

func (x *CanonicalizeResponse) Reset() {
	*x = CanonicalizeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalizeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalizeResponse) ProtoMessage() {}

func (x *CanonicalizeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalizeResponse.ProtoReflect.Descriptor instead.
func (*CanonicalizeResponse) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{3}
}

func (x *CanonicalizeResponse) GetMoves() []string {
	if x != nil {
		return x.Moves
	}
	return nil
}

type IsPositionInTakRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Position string `protobuf:"bytes,1,opt,name=position,proto3" json:"position,omitempty"`
}

func (x *IsPositionInTakRequest) Reset() {
	*x = IsPositionInTakRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsPositionInTakRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsPositionInTakRequest) ProtoMessage() {}

func (x *IsPositionInTakRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsPositionInTakRequest.ProtoReflect.Descriptor instead.
func (*IsPositionInTakRequest) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{4}
}

func (x *IsPositionInTakRequest) GetPosition() string {
	if x != nil {
		return x.Position
	}
	return ""
}

type IsPositionInTakResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InTak   bool   `protobuf:"varint,1,opt,name=inTak,proto3" json:"inTak,omitempty"`
	TakMove string `protobuf:"bytes,2,opt,name=takMove,proto3" json:"takMove,omitempty"`
}

func (x *IsPositionInTakResponse) Reset() {
	*x = IsPositionInTakResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tak_proto_taktician_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsPositionInTakResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsPositionInTakResponse) ProtoMessage() {}

func (x *IsPositionInTakResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tak_proto_taktician_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsPositionInTakResponse.ProtoReflect.Descriptor instead.
func (*IsPositionInTakResponse) Descriptor() ([]byte, []int) {
	return file_tak_proto_taktician_proto_rawDescGZIP(), []int{5}
}

func (x *IsPositionInTakResponse) GetInTak() bool {
	if x != nil {
		return x.InTak
	}
	return false
}

func (x *IsPositionInTakResponse) GetTakMove() string {
	if x != nil {
		return x.TakMove
	}
	return ""
}

var File_tak_proto_taktician_proto protoreflect.FileDescriptor

var file_tak_proto_taktician_proto_rawDesc = []byte{
	0x0a, 0x19, 0x74, 0x61, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x6b, 0x74,
	0x69, 0x63, 0x69, 0x61, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x61, 0x6b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x0e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x65, 0x63, 0x69, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x65,
	0x63, 0x69, 0x73, 0x65, 0x22, 0x4d, 0x0a, 0x0f, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x76, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x02, 0x70, 0x76, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x64, 0x65,
	0x70, 0x74, 0x68, 0x22, 0x3f, 0x0a, 0x13, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c,
	0x69, 0x7a, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x6f, 0x76, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x6f, 0x76, 0x65, 0x73, 0x22, 0x2c, 0x0a, 0x14, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61,
	0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x6d, 0x6f, 0x76, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x76,
	0x65, 0x73, 0x22, 0x34, 0x0a, 0x16, 0x49, 0x73, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x54, 0x61, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x49, 0x0a, 0x17, 0x49, 0x73, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x54, 0x61, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x05, 0x69, 0x6e, 0x54, 0x61, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x61, 0x6b,
	0x4d, 0x6f, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x61, 0x6b, 0x4d,
	0x6f, 0x76, 0x65, 0x32, 0xfe, 0x01, 0x0a, 0x09, 0x54, 0x61, 0x6b, 0x74, 0x69, 0x63, 0x69, 0x61,
	0x6e, 0x12, 0x42, 0x0a, 0x07, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x12, 0x19, 0x2e, 0x74,
	0x61, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0c, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0f, 0x49, 0x73, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x12, 0x21, 0x2e, 0x74, 0x61,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x73, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22,
	0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x73, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x54, 0x61, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x6c, 0x68, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x61, 0x6b, 0x74, 0x69,
	0x63, 0x69, 0x61, 0x6e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tak_proto_taktician_proto_rawDescOnce sync.Once
	file_tak_proto_taktician_proto_rawDescData = file_tak_proto_taktician_proto_rawDesc
)

func file_tak_proto_taktician_proto_rawDescGZIP() []byte {
	file_tak_proto_taktician_proto_rawDescOnce.Do(func() {
		file_tak_proto_taktician_proto_rawDescData = protoimpl.X.CompressGZIP(file_tak_proto_taktician_proto_rawDescData)
	})
	return file_tak_proto_taktician_proto_rawDescData
}

var file_tak_proto_taktician_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tak_proto_taktician_proto_goTypes = []interface{}{
	(*AnalyzeRequest)(nil),          // 0: tak.proto.AnalyzeRequest
	(*AnalyzeResponse)(nil),         // 1: tak.proto.AnalyzeResponse
	(*CanonicalizeRequest)(nil),     // 2: tak.proto.CanonicalizeRequest
	(*CanonicalizeResponse)(nil),    // 3: tak.proto.CanonicalizeResponse
	(*IsPositionInTakRequest)(nil),  // 4: tak.proto.IsPositionInTakRequest
	(*IsPositionInTakResponse)(nil), // 5: tak.proto.IsPositionInTakResponse
}
var file_tak_proto_taktician_proto_depIdxs = []int32{
	0, // 0: tak.proto.Taktician.Analyze:input_type -> tak.proto.AnalyzeRequest
	2, // 1: tak.proto.Taktician.Canonicalize:input_type -> tak.proto.CanonicalizeRequest
	4, // 2: tak.proto.Taktician.IsPositionInTak:input_type -> tak.proto.IsPositionInTakRequest
	1, // 3: tak.proto.Taktician.Analyze:output_type -> tak.proto.AnalyzeResponse
	3, // 4: tak.proto.Taktician.Canonicalize:output_type -> tak.proto.CanonicalizeResponse
	5, // 5: tak.proto.Taktician.IsPositionInTak:output_type -> tak.proto.IsPositionInTakResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_tak_proto_taktician_proto_init() }
func file_tak_proto_taktician_proto_init() {
	if File_tak_proto_taktician_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tak_proto_taktician_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeRequest); i {
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
		file_tak_proto_taktician_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeResponse); i {
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
		file_tak_proto_taktician_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalizeRequest); i {
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
		file_tak_proto_taktician_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalizeResponse); i {
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
		file_tak_proto_taktician_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsPositionInTakRequest); i {
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
		file_tak_proto_taktician_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsPositionInTakResponse); i {
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
			RawDescriptor: file_tak_proto_taktician_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tak_proto_taktician_proto_goTypes,
		DependencyIndexes: file_tak_proto_taktician_proto_depIdxs,
		MessageInfos:      file_tak_proto_taktician_proto_msgTypes,
	}.Build()
	File_tak_proto_taktician_proto = out.File
	file_tak_proto_taktician_proto_rawDesc = nil
	file_tak_proto_taktician_proto_goTypes = nil
	file_tak_proto_taktician_proto_depIdxs = nil
}
