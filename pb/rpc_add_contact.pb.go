// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.1
// source: rpc_add_contact.proto

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

type AddContactRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerId  int64 `protobuf:"varint,1,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	TargetId int64 `protobuf:"varint,2,opt,name=target_id,json=targetId,proto3" json:"target_id,omitempty"`
	Type     int32 `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *AddContactRequest) Reset() {
	*x = AddContactRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_add_contact_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddContactRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddContactRequest) ProtoMessage() {}

func (x *AddContactRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_add_contact_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddContactRequest.ProtoReflect.Descriptor instead.
func (*AddContactRequest) Descriptor() ([]byte, []int) {
	return file_rpc_add_contact_proto_rawDescGZIP(), []int{0}
}

func (x *AddContactRequest) GetOwnerId() int64 {
	if x != nil {
		return x.OwnerId
	}
	return 0
}

func (x *AddContactRequest) GetTargetId() int64 {
	if x != nil {
		return x.TargetId
	}
	return 0
}

func (x *AddContactRequest) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

type AddContactResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AddContactResponse) Reset() {
	*x = AddContactResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_add_contact_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddContactResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddContactResponse) ProtoMessage() {}

func (x *AddContactResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_add_contact_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddContactResponse.ProtoReflect.Descriptor instead.
func (*AddContactResponse) Descriptor() ([]byte, []int) {
	return file_rpc_add_contact_proto_rawDescGZIP(), []int{1}
}

func (x *AddContactResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_rpc_add_contact_proto protoreflect.FileDescriptor

var file_rpc_add_contact_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x70, 0x63, 0x5f, 0x61, 0x64, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x5f, 0x0a, 0x11, 0x41,
	0x64, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x2e, 0x0a, 0x12,
	0x41, 0x64, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0b, 0x5a, 0x09,
	0x49, 0x4d, 0x43, 0x68, 0x61, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_rpc_add_contact_proto_rawDescOnce sync.Once
	file_rpc_add_contact_proto_rawDescData = file_rpc_add_contact_proto_rawDesc
)

func file_rpc_add_contact_proto_rawDescGZIP() []byte {
	file_rpc_add_contact_proto_rawDescOnce.Do(func() {
		file_rpc_add_contact_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_add_contact_proto_rawDescData)
	})
	return file_rpc_add_contact_proto_rawDescData
}

var file_rpc_add_contact_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_add_contact_proto_goTypes = []interface{}{
	(*AddContactRequest)(nil),  // 0: pb.AddContactRequest
	(*AddContactResponse)(nil), // 1: pb.AddContactResponse
}
var file_rpc_add_contact_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_add_contact_proto_init() }
func file_rpc_add_contact_proto_init() {
	if File_rpc_add_contact_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_add_contact_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddContactRequest); i {
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
		file_rpc_add_contact_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddContactResponse); i {
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
			RawDescriptor: file_rpc_add_contact_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_add_contact_proto_goTypes,
		DependencyIndexes: file_rpc_add_contact_proto_depIdxs,
		MessageInfos:      file_rpc_add_contact_proto_msgTypes,
	}.Build()
	File_rpc_add_contact_proto = out.File
	file_rpc_add_contact_proto_rawDesc = nil
	file_rpc_add_contact_proto_goTypes = nil
	file_rpc_add_contact_proto_depIdxs = nil
}