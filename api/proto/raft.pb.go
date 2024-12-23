// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: api/proto/raft.proto

package proto

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

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term     int64   `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	Command  string  `protobuf:"bytes,3,opt,name=Command,proto3" json:"Command,omitempty"`
	Key      string  `protobuf:"bytes,4,opt,name=Key,proto3" json:"Key,omitempty"`
	Value    *string `protobuf:"bytes,5,opt,name=Value,proto3,oneof" json:"Value,omitempty"`
	OldValue *string `protobuf:"bytes,6,opt,name=OldValue,proto3,oneof" json:"OldValue,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_raft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_raft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_api_proto_raft_proto_rawDescGZIP(), []int{0}
}

func (x *LogEntry) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *LogEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LogEntry) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

func (x *LogEntry) GetOldValue() string {
	if x != nil && x.OldValue != nil {
		return *x.OldValue
	}
	return ""
}

type AppendEntriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64       `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	PrevLogIndex int64       `protobuf:"varint,2,opt,name=PrevLogIndex,proto3" json:"PrevLogIndex,omitempty"`
	PrevLogTerm  int64       `protobuf:"varint,3,opt,name=PrevLogTerm,proto3" json:"PrevLogTerm,omitempty"`
	LeaderCommit int64       `protobuf:"varint,4,opt,name=LeaderCommit,proto3" json:"LeaderCommit,omitempty"`
	LeaderID     int64       `protobuf:"varint,5,opt,name=LeaderID,proto3" json:"LeaderID,omitempty"`
	Entries      []*LogEntry `protobuf:"bytes,6,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *AppendEntriesRequest) Reset() {
	*x = AppendEntriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_raft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesRequest) ProtoMessage() {}

func (x *AppendEntriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_raft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesRequest.ProtoReflect.Descriptor instead.
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_raft_proto_rawDescGZIP(), []int{1}
}

func (x *AppendEntriesRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesRequest) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntriesRequest) GetPrevLogTerm() int64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntriesRequest) GetLeaderCommit() int64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

func (x *AppendEntriesRequest) GetLeaderID() int64 {
	if x != nil {
		return x.LeaderID
	}
	return 0
}

func (x *AppendEntriesRequest) GetEntries() []*LogEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type AppendEntriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Success bool  `protobuf:"varint,2,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *AppendEntriesResponse) Reset() {
	*x = AppendEntriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_raft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesResponse) ProtoMessage() {}

func (x *AppendEntriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_raft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesResponse.ProtoReflect.Descriptor instead.
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_raft_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntriesResponse) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	CandidateID  int64 `protobuf:"varint,2,opt,name=CandidateID,proto3" json:"CandidateID,omitempty"`
	LastLogIndex int64 `protobuf:"varint,3,opt,name=LastLogIndex,proto3" json:"LastLogIndex,omitempty"`
	LastLogTerm  int64 `protobuf:"varint,4,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_raft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_raft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_raft_proto_rawDescGZIP(), []int{3}
}

func (x *VoteRequest) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteRequest) GetCandidateID() int64 {
	if x != nil {
		return x.CandidateID
	}
	return 0
}

func (x *VoteRequest) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *VoteRequest) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_raft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_raft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_raft_proto_rawDescGZIP(), []int{4}
}

func (x *VoteResponse) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteResponse) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

var File_api_proto_raft_proto protoreflect.FileDescriptor

var file_api_proto_raft_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x61, 0x66, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x61, 0x66, 0x74, 0x22, 0x9d, 0x01, 0x0a,
	0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a,
	0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x4f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x4f, 0x6c, 0x64, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42,
	0x0b, 0x0a, 0x09, 0x5f, 0x4f, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xda, 0x01, 0x0a,
	0x14, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x65,
	0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0c, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a,
	0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12,
	0x22, 0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x44, 0x12,
	0x28, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x45, 0x0a, 0x15, 0x41, 0x70, 0x70,
	0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x89, 0x01, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4c, 0x61,
	0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61,
	0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x44, 0x0a, 0x0c,
	0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x20, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x65, 0x64, 0x32, 0x86, 0x01, 0x0a, 0x04, 0x52, 0x61, 0x66, 0x74, 0x12, 0x34, 0x0a, 0x0b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x61, 0x66,
	0x74, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x72, 0x61, 0x66, 0x74, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x12, 0x1a, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b,
	0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x13, 0x5a, 0x11, 0x64,
	0x69, 0x73, 0x74, 0x5f, 0x64, 0x62, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_raft_proto_rawDescOnce sync.Once
	file_api_proto_raft_proto_rawDescData = file_api_proto_raft_proto_rawDesc
)

func file_api_proto_raft_proto_rawDescGZIP() []byte {
	file_api_proto_raft_proto_rawDescOnce.Do(func() {
		file_api_proto_raft_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_raft_proto_rawDescData)
	})
	return file_api_proto_raft_proto_rawDescData
}

var file_api_proto_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_proto_raft_proto_goTypes = []interface{}{
	(*LogEntry)(nil),              // 0: raft.LogEntry
	(*AppendEntriesRequest)(nil),  // 1: raft.AppendEntriesRequest
	(*AppendEntriesResponse)(nil), // 2: raft.AppendEntriesResponse
	(*VoteRequest)(nil),           // 3: raft.VoteRequest
	(*VoteResponse)(nil),          // 4: raft.VoteResponse
}
var file_api_proto_raft_proto_depIdxs = []int32{
	0, // 0: raft.AppendEntriesRequest.entries:type_name -> raft.LogEntry
	3, // 1: raft.Raft.RequestVote:input_type -> raft.VoteRequest
	1, // 2: raft.Raft.AppendEntries:input_type -> raft.AppendEntriesRequest
	4, // 3: raft.Raft.RequestVote:output_type -> raft.VoteResponse
	2, // 4: raft.Raft.AppendEntries:output_type -> raft.AppendEntriesResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_proto_raft_proto_init() }
func file_api_proto_raft_proto_init() {
	if File_api_proto_raft_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_raft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
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
		file_api_proto_raft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesRequest); i {
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
		file_api_proto_raft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesResponse); i {
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
		file_api_proto_raft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
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
		file_api_proto_raft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
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
	file_api_proto_raft_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_raft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_raft_proto_goTypes,
		DependencyIndexes: file_api_proto_raft_proto_depIdxs,
		MessageInfos:      file_api_proto_raft_proto_msgTypes,
	}.Build()
	File_api_proto_raft_proto = out.File
	file_api_proto_raft_proto_rawDesc = nil
	file_api_proto_raft_proto_goTypes = nil
	file_api_proto_raft_proto_depIdxs = nil
}
