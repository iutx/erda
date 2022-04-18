// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: action.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID            string                 `protobuf:"bytes,1,opt,name=ID,json=id,proto3" json:"ID,omitempty"`
	TimeCreated   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timeCreated,proto3" json:"timeCreated,omitempty"`
	TimeUpdated   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timeUpdated,proto3" json:"timeUpdated,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Category      string                 `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	DisplayName   string                 `protobuf:"bytes,6,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	LogoUrl       string                 `protobuf:"bytes,7,opt,name=logo_url,json=logoUrl,proto3" json:"logo_url,omitempty"`
	Desc          string                 `protobuf:"bytes,8,opt,name=desc,proto3" json:"desc,omitempty"`
	Readme        string                 `protobuf:"bytes,9,opt,name=readme,proto3" json:"readme,omitempty"`
	Dice          *structpb.Value        `protobuf:"bytes,10,opt,name=dice,proto3" json:"dice,omitempty"`
	Spec          *structpb.Value        `protobuf:"bytes,11,opt,name=spec,proto3" json:"spec,omitempty"`
	Version       string                 `protobuf:"bytes,12,opt,name=version,proto3" json:"version,omitempty"`
	Location      string                 `protobuf:"bytes,13,opt,name=location,proto3" json:"location,omitempty"`
	IsPublic      bool                   `protobuf:"varint,14,opt,name=isPublic,proto3" json:"isPublic,omitempty"`
	IsDefault     bool                   `protobuf:"varint,15,opt,name=is_default,json=isDefault,proto3" json:"is_default,omitempty"`
	IsDelete      bool                   `protobuf:"varint,16,opt,name=is_delete,json=isDelete,proto3" json:"is_delete,omitempty"`
	SoftDeletedAt *timestamppb.Timestamp `protobuf:"bytes,17,opt,name=softDeletedAt,proto3" json:"softDeletedAt,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{0}
}

func (x *Action) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Action) GetTimeCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeCreated
	}
	return nil
}

func (x *Action) GetTimeUpdated() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeUpdated
	}
	return nil
}

func (x *Action) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Action) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Action) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Action) GetLogoUrl() string {
	if x != nil {
		return x.LogoUrl
	}
	return ""
}

func (x *Action) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *Action) GetReadme() string {
	if x != nil {
		return x.Readme
	}
	return ""
}

func (x *Action) GetDice() *structpb.Value {
	if x != nil {
		return x.Dice
	}
	return nil
}

func (x *Action) GetSpec() *structpb.Value {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *Action) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Action) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Action) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *Action) GetIsDefault() bool {
	if x != nil {
		return x.IsDefault
	}
	return false
}

func (x *Action) GetIsDelete() bool {
	if x != nil {
		return x.IsDelete
	}
	return false
}

func (x *Action) GetSoftDeletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.SoftDeletedAt
	}
	return nil
}

type ActionNameWithVersionQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version        string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	IsDefault      bool   `protobuf:"varint,3,opt,name=is_default,json=isDefault,proto3" json:"is_default,omitempty"`
	LocationFilter string `protobuf:"bytes,4,opt,name=locationFilter,proto3" json:"locationFilter,omitempty"`
}

func (x *ActionNameWithVersionQuery) Reset() {
	*x = ActionNameWithVersionQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionNameWithVersionQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionNameWithVersionQuery) ProtoMessage() {}

func (x *ActionNameWithVersionQuery) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionNameWithVersionQuery.ProtoReflect.Descriptor instead.
func (*ActionNameWithVersionQuery) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{1}
}

func (x *ActionNameWithVersionQuery) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ActionNameWithVersionQuery) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *ActionNameWithVersionQuery) GetIsDefault() bool {
	if x != nil {
		return x.IsDefault
	}
	return false
}

func (x *ActionNameWithVersionQuery) GetLocationFilter() string {
	if x != nil {
		return x.LocationFilter
	}
	return ""
}

// list
type PipelineActionListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActionNameWithVersionQuery []*ActionNameWithVersionQuery `protobuf:"bytes,1,rep,name=actionNameWithVersionQuery,proto3" json:"actionNameWithVersionQuery,omitempty"`
	Categories                 []string                      `protobuf:"bytes,2,rep,name=categories,proto3" json:"categories,omitempty"`
	Locations                  []string                      `protobuf:"bytes,3,rep,name=locations,proto3" json:"locations,omitempty"`
	IsPublic                   bool                          `protobuf:"varint,4,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	YamlFormat                 bool                          `protobuf:"varint,5,opt,name=yamlFormat,proto3" json:"yamlFormat,omitempty"`
}

func (x *PipelineActionListRequest) Reset() {
	*x = PipelineActionListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionListRequest) ProtoMessage() {}

func (x *PipelineActionListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionListRequest.ProtoReflect.Descriptor instead.
func (*PipelineActionListRequest) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{2}
}

func (x *PipelineActionListRequest) GetActionNameWithVersionQuery() []*ActionNameWithVersionQuery {
	if x != nil {
		return x.ActionNameWithVersionQuery
	}
	return nil
}

func (x *PipelineActionListRequest) GetCategories() []string {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *PipelineActionListRequest) GetLocations() []string {
	if x != nil {
		return x.Locations
	}
	return nil
}

func (x *PipelineActionListRequest) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *PipelineActionListRequest) GetYamlFormat() bool {
	if x != nil {
		return x.YamlFormat
	}
	return false
}

type PipelineActionListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Action `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *PipelineActionListResponse) Reset() {
	*x = PipelineActionListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionListResponse) ProtoMessage() {}

func (x *PipelineActionListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionListResponse.ProtoReflect.Descriptor instead.
func (*PipelineActionListResponse) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{3}
}

func (x *PipelineActionListResponse) GetData() []*Action {
	if x != nil {
		return x.Data
	}
	return nil
}

// save
type PipelineActionSaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Readme    string `protobuf:"bytes,1,opt,name=readme,proto3" json:"readme,omitempty"`
	Dice      string `protobuf:"bytes,2,opt,name=dice,proto3" json:"dice,omitempty"`
	Spec      string `protobuf:"bytes,3,opt,name=spec,proto3" json:"spec,omitempty"`
	Location  string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	IsDefault bool   `protobuf:"varint,5,opt,name=is_default,json=isDefault,proto3" json:"is_default,omitempty"`
}

func (x *PipelineActionSaveRequest) Reset() {
	*x = PipelineActionSaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionSaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionSaveRequest) ProtoMessage() {}

func (x *PipelineActionSaveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionSaveRequest.ProtoReflect.Descriptor instead.
func (*PipelineActionSaveRequest) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{4}
}

func (x *PipelineActionSaveRequest) GetReadme() string {
	if x != nil {
		return x.Readme
	}
	return ""
}

func (x *PipelineActionSaveRequest) GetDice() string {
	if x != nil {
		return x.Dice
	}
	return ""
}

func (x *PipelineActionSaveRequest) GetSpec() string {
	if x != nil {
		return x.Spec
	}
	return ""
}

func (x *PipelineActionSaveRequest) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *PipelineActionSaveRequest) GetIsDefault() bool {
	if x != nil {
		return x.IsDefault
	}
	return false
}

type PipelineActionSaveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action *Action `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
}

func (x *PipelineActionSaveResponse) Reset() {
	*x = PipelineActionSaveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionSaveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionSaveResponse) ProtoMessage() {}

func (x *PipelineActionSaveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionSaveResponse.ProtoReflect.Descriptor instead.
func (*PipelineActionSaveResponse) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{5}
}

func (x *PipelineActionSaveResponse) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

// delete
type PipelineActionDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version  string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Location string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *PipelineActionDeleteRequest) Reset() {
	*x = PipelineActionDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionDeleteRequest) ProtoMessage() {}

func (x *PipelineActionDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionDeleteRequest.ProtoReflect.Descriptor instead.
func (*PipelineActionDeleteRequest) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{6}
}

func (x *PipelineActionDeleteRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PipelineActionDeleteRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PipelineActionDeleteRequest) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

type PipelineActionDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PipelineActionDeleteResponse) Reset() {
	*x = PipelineActionDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_action_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PipelineActionDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PipelineActionDeleteResponse) ProtoMessage() {}

func (x *PipelineActionDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_action_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PipelineActionDeleteResponse.ProtoReflect.Descriptor instead.
func (*PipelineActionDeleteResponse) Descriptor() ([]byte, []int) {
	return file_action_proto_rawDescGZIP(), []int{7}
}

var File_action_proto protoreflect.FileDescriptor

var file_action_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19,
	0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd6, 0x04, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x3c, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12,
	0x3c, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x21, 0x0a,
	0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x6f, 0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x64, 0x69, 0x63, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x64,
	0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x40, 0x0a,
	0x0d, 0x73, 0x6f, 0x66, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x11,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0d, 0x73, 0x6f, 0x66, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x91, 0x01, 0x0a, 0x1a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x57, 0x69,
	0x74, 0x68, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a,
	0x69, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x22, 0x8d, 0x02, 0x0a, 0x19, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x75, 0x0a, 0x1a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x57,
	0x69, 0x74, 0x68, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x57, 0x69, 0x74, 0x68,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x1a, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x57, 0x69, 0x74, 0x68, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65,
	0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x79, 0x61, 0x6d, 0x6c, 0x46, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x79, 0x61, 0x6d, 0x6c, 0x46, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x22, 0x53, 0x0a, 0x1a, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x35, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x96, 0x01, 0x0a, 0x19, 0x50, 0x69, 0x70,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x61, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x69,
	0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x22, 0x57, 0x0a, 0x1a, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x39, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x67, 0x0a, 0x1b, 0x50, 0x69,
	0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x1e, 0x0a, 0x1c, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xe0, 0x03, 0x0a, 0x0d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x92, 0x01, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x34,
	0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x2d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x9e, 0x01, 0x0a, 0x04, 0x53,
	0x61, 0x76, 0x65, 0x12, 0x34, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x61,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x35, 0x2e, 0x65, 0x72, 0x64, 0x61,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x29, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x23, 0x22, 0x21, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x61, 0x76, 0x65, 0x12, 0x98, 0x01, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x36, 0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x37,
	0x2e, 0x65, 0x72, 0x64, 0x61, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x69, 0x70, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x2a,
	0x15, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2d, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x72, 0x64, 0x61, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x2f, 0x65, 0x72, 0x64, 0x61, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d, 0x67, 0x6f, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2f, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_action_proto_rawDescOnce sync.Once
	file_action_proto_rawDescData = file_action_proto_rawDesc
)

func file_action_proto_rawDescGZIP() []byte {
	file_action_proto_rawDescOnce.Do(func() {
		file_action_proto_rawDescData = protoimpl.X.CompressGZIP(file_action_proto_rawDescData)
	})
	return file_action_proto_rawDescData
}

var file_action_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_action_proto_goTypes = []interface{}{
	(*Action)(nil),                       // 0: erda.core.pipeline.action.Action
	(*ActionNameWithVersionQuery)(nil),   // 1: erda.core.pipeline.action.ActionNameWithVersionQuery
	(*PipelineActionListRequest)(nil),    // 2: erda.core.pipeline.action.PipelineActionListRequest
	(*PipelineActionListResponse)(nil),   // 3: erda.core.pipeline.action.PipelineActionListResponse
	(*PipelineActionSaveRequest)(nil),    // 4: erda.core.pipeline.action.PipelineActionSaveRequest
	(*PipelineActionSaveResponse)(nil),   // 5: erda.core.pipeline.action.PipelineActionSaveResponse
	(*PipelineActionDeleteRequest)(nil),  // 6: erda.core.pipeline.action.PipelineActionDeleteRequest
	(*PipelineActionDeleteResponse)(nil), // 7: erda.core.pipeline.action.PipelineActionDeleteResponse
	(*timestamppb.Timestamp)(nil),        // 8: google.protobuf.Timestamp
	(*structpb.Value)(nil),               // 9: google.protobuf.Value
}
var file_action_proto_depIdxs = []int32{
	8,  // 0: erda.core.pipeline.action.Action.timeCreated:type_name -> google.protobuf.Timestamp
	8,  // 1: erda.core.pipeline.action.Action.timeUpdated:type_name -> google.protobuf.Timestamp
	9,  // 2: erda.core.pipeline.action.Action.dice:type_name -> google.protobuf.Value
	9,  // 3: erda.core.pipeline.action.Action.spec:type_name -> google.protobuf.Value
	8,  // 4: erda.core.pipeline.action.Action.softDeletedAt:type_name -> google.protobuf.Timestamp
	1,  // 5: erda.core.pipeline.action.PipelineActionListRequest.actionNameWithVersionQuery:type_name -> erda.core.pipeline.action.ActionNameWithVersionQuery
	0,  // 6: erda.core.pipeline.action.PipelineActionListResponse.data:type_name -> erda.core.pipeline.action.Action
	0,  // 7: erda.core.pipeline.action.PipelineActionSaveResponse.action:type_name -> erda.core.pipeline.action.Action
	2,  // 8: erda.core.pipeline.action.ActionService.List:input_type -> erda.core.pipeline.action.PipelineActionListRequest
	4,  // 9: erda.core.pipeline.action.ActionService.Save:input_type -> erda.core.pipeline.action.PipelineActionSaveRequest
	6,  // 10: erda.core.pipeline.action.ActionService.Delete:input_type -> erda.core.pipeline.action.PipelineActionDeleteRequest
	3,  // 11: erda.core.pipeline.action.ActionService.List:output_type -> erda.core.pipeline.action.PipelineActionListResponse
	5,  // 12: erda.core.pipeline.action.ActionService.Save:output_type -> erda.core.pipeline.action.PipelineActionSaveResponse
	7,  // 13: erda.core.pipeline.action.ActionService.Delete:output_type -> erda.core.pipeline.action.PipelineActionDeleteResponse
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_action_proto_init() }
func file_action_proto_init() {
	if File_action_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_action_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
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
		file_action_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionNameWithVersionQuery); i {
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
		file_action_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionListRequest); i {
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
		file_action_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionListResponse); i {
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
		file_action_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionSaveRequest); i {
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
		file_action_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionSaveResponse); i {
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
		file_action_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionDeleteRequest); i {
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
		file_action_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PipelineActionDeleteResponse); i {
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
			RawDescriptor: file_action_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_action_proto_goTypes,
		DependencyIndexes: file_action_proto_depIdxs,
		MessageInfos:      file_action_proto_msgTypes,
	}.Build()
	File_action_proto = out.File
	file_action_proto_rawDesc = nil
	file_action_proto_goTypes = nil
	file_action_proto_depIdxs = nil
}
