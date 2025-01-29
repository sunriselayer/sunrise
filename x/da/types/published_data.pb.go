// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/da/v1/published_data.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// PublishedDataStatus
type Status int32

const (
	// Default value
	Status_STATUS_UNSPECIFIED Status = 0
	// Verified
	Status_STATUS_VERIFIED Status = 1
	// Rejected
	Status_STATUS_REJECTED Status = 2
	// After processing in the msg_server
	Status_STATUS_VOTING Status = 3
	// Verified the votes from the validators. Challenge can be received (after preBlocker)
	Status_STATUS_CHALLENGE_PERIOD Status = 4
	// reported as fraud. SubmitProof tx can be received (after received ChallengeForFraud tx)
	Status_STATUS_CHALLENGING Status = 5
)

var Status_name = map[int32]string{
	0: "STATUS_UNSPECIFIED",
	1: "STATUS_VERIFIED",
	2: "STATUS_REJECTED",
	3: "STATUS_VOTING",
	4: "STATUS_CHALLENGE_PERIOD",
	5: "STATUS_CHALLENGING",
}

var Status_value = map[string]int32{
	"STATUS_UNSPECIFIED":      0,
	"STATUS_VERIFIED":         1,
	"STATUS_REJECTED":         2,
	"STATUS_VOTING":           3,
	"STATUS_CHALLENGE_PERIOD": 4,
	"STATUS_CHALLENGING":      5,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dbb00f0e2104e63f, []int{0}
}

// PublishedData
type PublishedData struct {
	MetadataUri        string                                   `protobuf:"bytes,1,opt,name=metadata_uri,json=metadataUri,proto3" json:"metadata_uri,omitempty"`
	ParityShardCount   uint64                                   `protobuf:"varint,2,opt,name=parity_shard_count,json=parityShardCount,proto3" json:"parity_shard_count,omitempty"`
	ShardDoubleHashes  [][]byte                                 `protobuf:"bytes,3,rep,name=shard_double_hashes,json=shardDoubleHashes,proto3" json:"shard_double_hashes,omitempty"`
	Timestamp          time.Time                                `protobuf:"bytes,4,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	Status             Status                                   `protobuf:"varint,5,opt,name=status,proto3,enum=sunrise.da.v1.Status" json:"status,omitempty"`
	Publisher          string                                   `protobuf:"bytes,6,opt,name=publisher,proto3" json:"publisher,omitempty"`
	Challenger         string                                   `protobuf:"bytes,7,opt,name=challenger,proto3" json:"challenger,omitempty"`
	Collateral         github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,8,rep,name=collateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"collateral"`
	ChallengeTimestamp time.Time                                `protobuf:"bytes,9,opt,name=challenge_timestamp,json=challengeTimestamp,proto3,stdtime" json:"challenge_timestamp"`
	DataSourceInfo     string                                   `protobuf:"bytes,10,opt,name=data_source_info,json=dataSourceInfo,proto3" json:"data_source_info,omitempty"`
}

func (m *PublishedData) Reset()         { *m = PublishedData{} }
func (m *PublishedData) String() string { return proto.CompactTextString(m) }
func (*PublishedData) ProtoMessage()    {}
func (*PublishedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbb00f0e2104e63f, []int{0}
}
func (m *PublishedData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PublishedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PublishedData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PublishedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublishedData.Merge(m, src)
}
func (m *PublishedData) XXX_Size() int {
	return m.Size()
}
func (m *PublishedData) XXX_DiscardUnknown() {
	xxx_messageInfo_PublishedData.DiscardUnknown(m)
}

var xxx_messageInfo_PublishedData proto.InternalMessageInfo

func (m *PublishedData) GetMetadataUri() string {
	if m != nil {
		return m.MetadataUri
	}
	return ""
}

func (m *PublishedData) GetParityShardCount() uint64 {
	if m != nil {
		return m.ParityShardCount
	}
	return 0
}

func (m *PublishedData) GetShardDoubleHashes() [][]byte {
	if m != nil {
		return m.ShardDoubleHashes
	}
	return nil
}

func (m *PublishedData) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func (m *PublishedData) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_STATUS_UNSPECIFIED
}

func (m *PublishedData) GetPublisher() string {
	if m != nil {
		return m.Publisher
	}
	return ""
}

func (m *PublishedData) GetChallenger() string {
	if m != nil {
		return m.Challenger
	}
	return ""
}

func (m *PublishedData) GetCollateral() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Collateral
	}
	return nil
}

func (m *PublishedData) GetChallengeTimestamp() time.Time {
	if m != nil {
		return m.ChallengeTimestamp
	}
	return time.Time{}
}

func (m *PublishedData) GetDataSourceInfo() string {
	if m != nil {
		return m.DataSourceInfo
	}
	return ""
}

// Proof
type Proof struct {
	MetadataUri string   `protobuf:"bytes,1,opt,name=metadata_uri,json=metadataUri,proto3" json:"metadata_uri,omitempty"`
	Sender      string   `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Indices     []int64  `protobuf:"varint,3,rep,packed,name=indices,proto3" json:"indices,omitempty"`
	Proofs      [][]byte `protobuf:"bytes,4,rep,name=proofs,proto3" json:"proofs,omitempty"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbb00f0e2104e63f, []int{1}
}
func (m *Proof) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Proof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Proof.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Proof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proof.Merge(m, src)
}
func (m *Proof) XXX_Size() int {
	return m.Size()
}
func (m *Proof) XXX_DiscardUnknown() {
	xxx_messageInfo_Proof.DiscardUnknown(m)
}

var xxx_messageInfo_Proof proto.InternalMessageInfo

func (m *Proof) GetMetadataUri() string {
	if m != nil {
		return m.MetadataUri
	}
	return ""
}

func (m *Proof) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *Proof) GetIndices() []int64 {
	if m != nil {
		return m.Indices
	}
	return nil
}

func (m *Proof) GetProofs() [][]byte {
	if m != nil {
		return m.Proofs
	}
	return nil
}

func init() {
	proto.RegisterEnum("sunrise.da.v1.Status", Status_name, Status_value)
	proto.RegisterType((*PublishedData)(nil), "sunrise.da.v1.PublishedData")
	proto.RegisterType((*Proof)(nil), "sunrise.da.v1.Proof")
}

func init() {
	proto.RegisterFile("sunrise/da/v1/published_data.proto", fileDescriptor_dbb00f0e2104e63f)
}

var fileDescriptor_dbb00f0e2104e63f = []byte{
	// 626 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xcd, 0x52, 0xdb, 0x30,
	0x10, 0x8e, 0x49, 0x08, 0x44, 0xfc, 0x34, 0x88, 0x96, 0xba, 0xb4, 0xe3, 0xa4, 0x9c, 0x3c, 0x4c,
	0x91, 0x1b, 0xfa, 0x04, 0x24, 0x71, 0x21, 0x1d, 0x06, 0x32, 0x4e, 0xd2, 0x43, 0x2f, 0x1e, 0xd9,
	0x56, 0x12, 0x0d, 0x8e, 0xe5, 0x91, 0x64, 0xa6, 0xbc, 0x05, 0xd3, 0xc7, 0xe8, 0xa5, 0xaf, 0xc1,
	0x91, 0x63, 0x4f, 0xa5, 0x03, 0x2f, 0xd2, 0xb1, 0x6c, 0x27, 0xe9, 0xad, 0x3d, 0x59, 0xfb, 0x7d,
	0xdf, 0x4a, 0xbb, 0xfb, 0xad, 0xc1, 0x81, 0x48, 0x22, 0x4e, 0x05, 0xb1, 0x02, 0x6c, 0x5d, 0xb7,
	0xac, 0x38, 0xf1, 0x42, 0x2a, 0xa6, 0x24, 0x70, 0x03, 0x2c, 0x31, 0x8a, 0x39, 0x93, 0x0c, 0x6e,
	0xe5, 0x1a, 0x14, 0x60, 0x74, 0xdd, 0xda, 0x37, 0x7c, 0x26, 0x66, 0x4c, 0x58, 0x1e, 0x16, 0xc4,
	0xba, 0x6e, 0x79, 0x44, 0xe2, 0x96, 0xe5, 0x33, 0x1a, 0x65, 0xf2, 0xfd, 0xe7, 0x13, 0x36, 0x61,
	0xea, 0x68, 0xa5, 0xa7, 0x1c, 0x6d, 0x4c, 0x18, 0x9b, 0x84, 0xc4, 0x52, 0x91, 0x97, 0x8c, 0x2d,
	0x49, 0x67, 0x44, 0x48, 0x3c, 0x8b, 0x33, 0xc1, 0xc1, 0x8f, 0x0a, 0xd8, 0xea, 0x17, 0xcf, 0x77,
	0xb1, 0xc4, 0xf0, 0x2d, 0xd8, 0x9c, 0x11, 0x89, 0xd3, 0x4a, 0xdc, 0x84, 0x53, 0x5d, 0x6b, 0x6a,
	0x66, 0xcd, 0xd9, 0x28, 0xb0, 0x11, 0xa7, 0xf0, 0x1d, 0x80, 0x31, 0xe6, 0x54, 0xde, 0xb8, 0x62,
	0x8a, 0x79, 0xe0, 0xfa, 0x2c, 0x89, 0xa4, 0xbe, 0xd2, 0xd4, 0xcc, 0x8a, 0x53, 0xcf, 0x98, 0x41,
	0x4a, 0x74, 0x52, 0x1c, 0x22, 0xb0, 0x9b, 0xc9, 0x02, 0x96, 0x78, 0x21, 0x71, 0xa7, 0x58, 0x4c,
	0x89, 0xd0, 0xcb, 0xcd, 0xb2, 0xb9, 0xe9, 0xec, 0x28, 0xaa, 0xab, 0x98, 0x33, 0x45, 0xc0, 0x36,
	0xa8, 0xcd, 0xab, 0xd4, 0x2b, 0x4d, 0xcd, 0xdc, 0x38, 0xde, 0x47, 0x59, 0x1f, 0xa8, 0xe8, 0x03,
	0x0d, 0x0b, 0x45, 0x7b, 0xfd, 0xee, 0x57, 0xa3, 0x74, 0xfb, 0xd0, 0xd0, 0x9c, 0x45, 0x1a, 0x3c,
	0x02, 0x55, 0x21, 0xb1, 0x4c, 0x84, 0xbe, 0xda, 0xd4, 0xcc, 0xed, 0xe3, 0x17, 0xe8, 0xaf, 0x69,
	0xa2, 0x81, 0x22, 0x9d, 0x5c, 0x04, 0xdf, 0x80, 0x5a, 0xe1, 0x01, 0xd7, 0xab, 0xaa, 0xe1, 0x05,
	0x00, 0x0d, 0x00, 0xfc, 0x29, 0x0e, 0x43, 0x12, 0x4d, 0x08, 0xd7, 0xd7, 0x14, 0xbd, 0x84, 0xc0,
	0x2b, 0x00, 0x7c, 0x16, 0x86, 0x58, 0x12, 0x8e, 0x43, 0x7d, 0xbd, 0x59, 0x36, 0x37, 0x8e, 0x5f,
	0xa1, 0xcc, 0x2f, 0x94, 0xfa, 0x85, 0x72, 0xbf, 0x50, 0x87, 0xd1, 0xa8, 0xfd, 0x3e, 0x2d, 0xf8,
	0xfb, 0x43, 0xc3, 0x9c, 0x50, 0x39, 0x4d, 0x3c, 0xe4, 0xb3, 0x99, 0x95, 0x9b, 0x9b, 0x7d, 0x8e,
	0x44, 0x70, 0x65, 0xc9, 0x9b, 0x98, 0x08, 0x95, 0x20, 0x9c, 0xa5, 0xeb, 0xe1, 0x08, 0xec, 0xce,
	0x9f, 0x76, 0x17, 0x73, 0xaa, 0xfd, 0xc7, 0x9c, 0xe0, 0xfc, 0x82, 0x39, 0x0b, 0x4d, 0x50, 0x57,
	0x8e, 0x0b, 0x96, 0x70, 0x9f, 0xb8, 0x34, 0x1a, 0x33, 0x1d, 0xa8, 0x4e, 0xb7, 0x53, 0x7c, 0xa0,
	0xe0, 0x5e, 0x34, 0x66, 0x07, 0x12, 0xac, 0xf6, 0x39, 0x63, 0xe3, 0x7f, 0x59, 0x94, 0x3d, 0x50,
	0x15, 0x24, 0x0a, 0x08, 0x57, 0xcb, 0x51, 0x73, 0xf2, 0x08, 0xea, 0x60, 0x8d, 0x46, 0x01, 0xf5,
	0xf3, 0x35, 0x28, 0x3b, 0x45, 0x98, 0x66, 0xc4, 0xe9, 0xed, 0x42, 0xaf, 0xa8, 0xfd, 0xc8, 0xa3,
	0xc3, 0x6f, 0x1a, 0xa8, 0x66, 0xa6, 0xc1, 0x3d, 0x00, 0x07, 0xc3, 0x93, 0xe1, 0x68, 0xe0, 0x8e,
	0x2e, 0x06, 0x7d, 0xbb, 0xd3, 0xfb, 0xd8, 0xb3, 0xbb, 0xf5, 0x12, 0xdc, 0x05, 0xcf, 0x72, 0xfc,
	0xb3, 0xed, 0x64, 0xa0, 0xb6, 0x04, 0x3a, 0xf6, 0x27, 0xbb, 0x33, 0xb4, 0xbb, 0xf5, 0x15, 0xb8,
	0x03, 0xb6, 0x0a, 0xe5, 0xe5, 0xb0, 0x77, 0x71, 0x5a, 0x2f, 0xc3, 0xd7, 0xe0, 0x65, 0x0e, 0x75,
	0xce, 0x4e, 0xce, 0xcf, 0xed, 0x8b, 0x53, 0xdb, 0xed, 0xdb, 0x4e, 0xef, 0xb2, 0x5b, 0xaf, 0x2c,
	0xbd, 0x58, 0x90, 0x69, 0xd2, 0x6a, 0xbb, 0x7b, 0xf7, 0x68, 0x68, 0xf7, 0x8f, 0x86, 0xf6, 0xfb,
	0xd1, 0xd0, 0x6e, 0x9f, 0x8c, 0xd2, 0xfd, 0x93, 0x51, 0xfa, 0xf9, 0x64, 0x94, 0xbe, 0x1c, 0x2e,
	0x79, 0x9b, 0x6f, 0x5e, 0x88, 0x6f, 0x08, 0x2f, 0x02, 0xeb, 0x6b, 0xfa, 0xeb, 0x2b, 0x8f, 0xbd,
	0xaa, 0x32, 0xeb, 0xc3, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x56, 0x43, 0x41, 0xe6, 0x15, 0x04,
	0x00, 0x00,
}

func (m *PublishedData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PublishedData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PublishedData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DataSourceInfo) > 0 {
		i -= len(m.DataSourceInfo)
		copy(dAtA[i:], m.DataSourceInfo)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.DataSourceInfo)))
		i--
		dAtA[i] = 0x52
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.ChallengeTimestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ChallengeTimestamp):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintPublishedData(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x4a
	if len(m.Collateral) > 0 {
		for iNdEx := len(m.Collateral) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Collateral[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPublishedData(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.Challenger) > 0 {
		i -= len(m.Challenger)
		copy(dAtA[i:], m.Challenger)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Challenger)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Publisher) > 0 {
		i -= len(m.Publisher)
		copy(dAtA[i:], m.Publisher)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Publisher)))
		i--
		dAtA[i] = 0x32
	}
	if m.Status != 0 {
		i = encodeVarintPublishedData(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x28
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintPublishedData(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x22
	if len(m.ShardDoubleHashes) > 0 {
		for iNdEx := len(m.ShardDoubleHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ShardDoubleHashes[iNdEx])
			copy(dAtA[i:], m.ShardDoubleHashes[iNdEx])
			i = encodeVarintPublishedData(dAtA, i, uint64(len(m.ShardDoubleHashes[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.ParityShardCount != 0 {
		i = encodeVarintPublishedData(dAtA, i, uint64(m.ParityShardCount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.MetadataUri) > 0 {
		i -= len(m.MetadataUri)
		copy(dAtA[i:], m.MetadataUri)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.MetadataUri)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Proof) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Proof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Proof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Proofs) > 0 {
		for iNdEx := len(m.Proofs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Proofs[iNdEx])
			copy(dAtA[i:], m.Proofs[iNdEx])
			i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Proofs[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Indices) > 0 {
		dAtA4 := make([]byte, len(m.Indices)*10)
		var j3 int
		for _, num1 := range m.Indices {
			num := uint64(num1)
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		i -= j3
		copy(dAtA[i:], dAtA4[:j3])
		i = encodeVarintPublishedData(dAtA, i, uint64(j3))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.MetadataUri) > 0 {
		i -= len(m.MetadataUri)
		copy(dAtA[i:], m.MetadataUri)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.MetadataUri)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPublishedData(dAtA []byte, offset int, v uint64) int {
	offset -= sovPublishedData(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PublishedData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MetadataUri)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	if m.ParityShardCount != 0 {
		n += 1 + sovPublishedData(uint64(m.ParityShardCount))
	}
	if len(m.ShardDoubleHashes) > 0 {
		for _, b := range m.ShardDoubleHashes {
			l = len(b)
			n += 1 + l + sovPublishedData(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovPublishedData(uint64(l))
	if m.Status != 0 {
		n += 1 + sovPublishedData(uint64(m.Status))
	}
	l = len(m.Publisher)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	l = len(m.Challenger)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	if len(m.Collateral) > 0 {
		for _, e := range m.Collateral {
			l = e.Size()
			n += 1 + l + sovPublishedData(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ChallengeTimestamp)
	n += 1 + l + sovPublishedData(uint64(l))
	l = len(m.DataSourceInfo)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	return n
}

func (m *Proof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MetadataUri)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
	}
	if len(m.Indices) > 0 {
		l = 0
		for _, e := range m.Indices {
			l += sovPublishedData(uint64(e))
		}
		n += 1 + sovPublishedData(uint64(l)) + l
	}
	if len(m.Proofs) > 0 {
		for _, b := range m.Proofs {
			l = len(b)
			n += 1 + l + sovPublishedData(uint64(l))
		}
	}
	return n
}

func sovPublishedData(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPublishedData(x uint64) (n int) {
	return sovPublishedData(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PublishedData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPublishedData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PublishedData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PublishedData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetadataUri", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MetadataUri = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParityShardCount", wireType)
			}
			m.ParityShardCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParityShardCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShardDoubleHashes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ShardDoubleHashes = append(m.ShardDoubleHashes, make([]byte, postIndex-iNdEx))
			copy(m.ShardDoubleHashes[len(m.ShardDoubleHashes)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= Status(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Publisher", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Publisher = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Challenger", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Challenger = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Collateral", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Collateral = append(m.Collateral, types.Coin{})
			if err := m.Collateral[len(m.Collateral)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChallengeTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.ChallengeTimestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataSourceInfo", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataSourceInfo = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPublishedData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPublishedData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Proof) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPublishedData
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Proof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Proof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetadataUri", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MetadataUri = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v int64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPublishedData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= int64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Indices = append(m.Indices, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPublishedData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPublishedData
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPublishedData
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Indices) == 0 {
					m.Indices = make([]int64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPublishedData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= int64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Indices = append(m.Indices, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Indices", wireType)
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proofs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthPublishedData
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPublishedData
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proofs = append(m.Proofs, make([]byte, postIndex-iNdEx))
			copy(m.Proofs[len(m.Proofs)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPublishedData(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPublishedData
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPublishedData(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPublishedData
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPublishedData
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPublishedData
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPublishedData
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPublishedData
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPublishedData        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPublishedData          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPublishedData = fmt.Errorf("proto: unexpected end of group")
)
