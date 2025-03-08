// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: loadbalancer.proto

package __

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

type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ServerId      string                 `protobuf:"bytes,1,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"` // Unique ID of the server.
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`                   // Address (host:port) of the server.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_loadbalancer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetServerId() string {
	if x != nil {
		return x.ServerId
	}
	return ""
}

func (x *RegisterRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_loadbalancer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RegisterResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type LoadReportRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ServerId      string                 `protobuf:"bytes,1,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	CurrentLoad   float32                `protobuf:"fixed32,2,opt,name=current_load,json=currentLoad,proto3" json:"current_load,omitempty"` // Example: number of active requests.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoadReportRequest) Reset() {
	*x = LoadReportRequest{}
	mi := &file_loadbalancer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoadReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadReportRequest) ProtoMessage() {}

func (x *LoadReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadReportRequest.ProtoReflect.Descriptor instead.
func (*LoadReportRequest) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{2}
}

func (x *LoadReportRequest) GetServerId() string {
	if x != nil {
		return x.ServerId
	}
	return ""
}

func (x *LoadReportRequest) GetCurrentLoad() float32 {
	if x != nil {
		return x.CurrentLoad
	}
	return 0
}

type LoadReportResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoadReportResponse) Reset() {
	*x = LoadReportResponse{}
	mi := &file_loadbalancer_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoadReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadReportResponse) ProtoMessage() {}

func (x *LoadReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadReportResponse.ProtoReflect.Descriptor instead.
func (*LoadReportResponse) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{3}
}

func (x *LoadReportResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LoadReportResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type BestServerRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Policy        string                 `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"` // "pick_first", "round_robin", or "least_load".
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BestServerRequest) Reset() {
	*x = BestServerRequest{}
	mi := &file_loadbalancer_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BestServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestServerRequest) ProtoMessage() {}

func (x *BestServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestServerRequest.ProtoReflect.Descriptor instead.
func (*BestServerRequest) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{4}
}

func (x *BestServerRequest) GetPolicy() string {
	if x != nil {
		return x.Policy
	}
	return ""
}

type BestServerResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"` // Chosen server's address.
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BestServerResponse) Reset() {
	*x = BestServerResponse{}
	mi := &file_loadbalancer_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BestServerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BestServerResponse) ProtoMessage() {}

func (x *BestServerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BestServerResponse.ProtoReflect.Descriptor instead.
func (*BestServerResponse) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{5}
}

func (x *BestServerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *BestServerResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *BestServerResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ComputeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskData      string                 `protobuf:"bytes,1,opt,name=task_data,json=taskData,proto3" json:"task_data,omitempty"` // Data for the computational task.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ComputeRequest) Reset() {
	*x = ComputeRequest{}
	mi := &file_loadbalancer_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ComputeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComputeRequest) ProtoMessage() {}

func (x *ComputeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComputeRequest.ProtoReflect.Descriptor instead.
func (*ComputeRequest) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{6}
}

func (x *ComputeRequest) GetTaskData() string {
	if x != nil {
		return x.TaskData
	}
	return ""
}

type ComputeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` // Result returned by the server.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ComputeResponse) Reset() {
	*x = ComputeResponse{}
	mi := &file_loadbalancer_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ComputeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComputeResponse) ProtoMessage() {}

func (x *ComputeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loadbalancer_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComputeResponse.ProtoReflect.Descriptor instead.
func (*ComputeResponse) Descriptor() ([]byte, []int) {
	return file_loadbalancer_proto_rawDescGZIP(), []int{7}
}

func (x *ComputeResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_loadbalancer_proto protoreflect.FileDescriptor

var file_loadbalancer_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x72, 0x22, 0x48, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x46, 0x0a, 0x10,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x53, 0x0a, 0x11, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x61, 0x64, 0x22, 0x48, 0x0a, 0x12, 0x4c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x2b, 0x0a, 0x11, 0x42, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x22, 0x62, 0x0a, 0x12, 0x42, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x2d, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x44,
	0x61, 0x74, 0x61, 0x22, 0x29, 0x0a, 0x0f, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x91,
	0x02, 0x0a, 0x13, 0x4c, 0x6f, 0x61, 0x64, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0a, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x4c, 0x6f, 0x61, 0x64, 0x12, 0x1f, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x0d,
	0x47, 0x65, 0x74, 0x42, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1f, 0x2e,
	0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x42, 0x65, 0x73,
	0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x42, 0x65,
	0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x32, 0x5e, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x54,
	0x61, 0x73, 0x6b, 0x12, 0x1c, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x72,
	0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_loadbalancer_proto_rawDescOnce sync.Once
	file_loadbalancer_proto_rawDescData []byte
)

func file_loadbalancer_proto_rawDescGZIP() []byte {
	file_loadbalancer_proto_rawDescOnce.Do(func() {
		file_loadbalancer_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_loadbalancer_proto_rawDesc), len(file_loadbalancer_proto_rawDesc)))
	})
	return file_loadbalancer_proto_rawDescData
}

var file_loadbalancer_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_loadbalancer_proto_goTypes = []any{
	(*RegisterRequest)(nil),    // 0: loadbalancer.RegisterRequest
	(*RegisterResponse)(nil),   // 1: loadbalancer.RegisterResponse
	(*LoadReportRequest)(nil),  // 2: loadbalancer.LoadReportRequest
	(*LoadReportResponse)(nil), // 3: loadbalancer.LoadReportResponse
	(*BestServerRequest)(nil),  // 4: loadbalancer.BestServerRequest
	(*BestServerResponse)(nil), // 5: loadbalancer.BestServerResponse
	(*ComputeRequest)(nil),     // 6: loadbalancer.ComputeRequest
	(*ComputeResponse)(nil),    // 7: loadbalancer.ComputeResponse
}
var file_loadbalancer_proto_depIdxs = []int32{
	0, // 0: loadbalancer.LoadBalancerService.RegisterServer:input_type -> loadbalancer.RegisterRequest
	2, // 1: loadbalancer.LoadBalancerService.ReportLoad:input_type -> loadbalancer.LoadReportRequest
	4, // 2: loadbalancer.LoadBalancerService.GetBestServer:input_type -> loadbalancer.BestServerRequest
	6, // 3: loadbalancer.ComputeService.ComputeTask:input_type -> loadbalancer.ComputeRequest
	1, // 4: loadbalancer.LoadBalancerService.RegisterServer:output_type -> loadbalancer.RegisterResponse
	3, // 5: loadbalancer.LoadBalancerService.ReportLoad:output_type -> loadbalancer.LoadReportResponse
	5, // 6: loadbalancer.LoadBalancerService.GetBestServer:output_type -> loadbalancer.BestServerResponse
	7, // 7: loadbalancer.ComputeService.ComputeTask:output_type -> loadbalancer.ComputeResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_loadbalancer_proto_init() }
func file_loadbalancer_proto_init() {
	if File_loadbalancer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_loadbalancer_proto_rawDesc), len(file_loadbalancer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_loadbalancer_proto_goTypes,
		DependencyIndexes: file_loadbalancer_proto_depIdxs,
		MessageInfos:      file_loadbalancer_proto_msgTypes,
	}.Build()
	File_loadbalancer_proto = out.File
	file_loadbalancer_proto_goTypes = nil
	file_loadbalancer_proto_depIdxs = nil
}
