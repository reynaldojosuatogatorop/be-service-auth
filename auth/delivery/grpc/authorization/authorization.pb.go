// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.0--rc1
// source: auth/repository/grpc/proto/authorization.proto

package authorization

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

type AuthorizationAuthServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AuthorizationAuthServiceRequest) Reset() {
	*x = AuthorizationAuthServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_repository_grpc_proto_authorization_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizationAuthServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationAuthServiceRequest) ProtoMessage() {}

func (x *AuthorizationAuthServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_repository_grpc_proto_authorization_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationAuthServiceRequest.ProtoReflect.Descriptor instead.
func (*AuthorizationAuthServiceRequest) Descriptor() ([]byte, []int) {
	return file_auth_repository_grpc_proto_authorization_proto_rawDescGZIP(), []int{0}
}

func (x *AuthorizationAuthServiceRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AuthorizationAuthServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Email           string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Role            string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	Token           string `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
	ExpiredDatetime string `protobuf:"bytes,5,opt,name=expired_datetime,json=expiredDatetime,proto3" json:"expired_datetime,omitempty"`
}

func (x *AuthorizationAuthServiceResponse) Reset() {
	*x = AuthorizationAuthServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_repository_grpc_proto_authorization_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorizationAuthServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizationAuthServiceResponse) ProtoMessage() {}

func (x *AuthorizationAuthServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_repository_grpc_proto_authorization_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizationAuthServiceResponse.ProtoReflect.Descriptor instead.
func (*AuthorizationAuthServiceResponse) Descriptor() ([]byte, []int) {
	return file_auth_repository_grpc_proto_authorization_proto_rawDescGZIP(), []int{1}
}

func (x *AuthorizationAuthServiceResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AuthorizationAuthServiceResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *AuthorizationAuthServiceResponse) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *AuthorizationAuthServiceResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *AuthorizationAuthServiceResponse) GetExpiredDatetime() string {
	if x != nil {
		return x.ExpiredDatetime
	}
	return ""
}

var File_auth_repository_grpc_proto_authorization_proto protoreflect.FileDescriptor

var file_auth_repository_grpc_proto_authorization_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x37, 0x0a, 0x1f, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x9d, 0x01, 0x0a, 0x20, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x5f,
	0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x32,
	0x80, 0x01, 0x0a, 0x14, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x68, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x75, 0x74,
	0x68, 0x12, 0x25, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x75, 0x74,
	0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_repository_grpc_proto_authorization_proto_rawDescOnce sync.Once
	file_auth_repository_grpc_proto_authorization_proto_rawDescData = file_auth_repository_grpc_proto_authorization_proto_rawDesc
)

func file_auth_repository_grpc_proto_authorization_proto_rawDescGZIP() []byte {
	file_auth_repository_grpc_proto_authorization_proto_rawDescOnce.Do(func() {
		file_auth_repository_grpc_proto_authorization_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_repository_grpc_proto_authorization_proto_rawDescData)
	})
	return file_auth_repository_grpc_proto_authorization_proto_rawDescData
}

var file_auth_repository_grpc_proto_authorization_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_auth_repository_grpc_proto_authorization_proto_goTypes = []interface{}{
	(*AuthorizationAuthServiceRequest)(nil),  // 0: auth.AuthorizationAuthServiceRequest
	(*AuthorizationAuthServiceResponse)(nil), // 1: auth.AuthorizationAuthServiceResponse
}
var file_auth_repository_grpc_proto_authorization_proto_depIdxs = []int32{
	0, // 0: auth.AuthorizationService.GetSessionServiceAuth:input_type -> auth.AuthorizationAuthServiceRequest
	1, // 1: auth.AuthorizationService.GetSessionServiceAuth:output_type -> auth.AuthorizationAuthServiceResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_repository_grpc_proto_authorization_proto_init() }
func file_auth_repository_grpc_proto_authorization_proto_init() {
	if File_auth_repository_grpc_proto_authorization_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_repository_grpc_proto_authorization_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizationAuthServiceRequest); i {
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
		file_auth_repository_grpc_proto_authorization_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorizationAuthServiceResponse); i {
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
			RawDescriptor: file_auth_repository_grpc_proto_authorization_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_repository_grpc_proto_authorization_proto_goTypes,
		DependencyIndexes: file_auth_repository_grpc_proto_authorization_proto_depIdxs,
		MessageInfos:      file_auth_repository_grpc_proto_authorization_proto_msgTypes,
	}.Build()
	File_auth_repository_grpc_proto_authorization_proto = out.File
	file_auth_repository_grpc_proto_authorization_proto_rawDesc = nil
	file_auth_repository_grpc_proto_authorization_proto_goTypes = nil
	file_auth_repository_grpc_proto_authorization_proto_depIdxs = nil
}
