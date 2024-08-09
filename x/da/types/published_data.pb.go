// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/da/published_data.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

type PublishedData struct {
	MetadataUri        string                                   `protobuf:"bytes,1,opt,name=metadata_uri,json=metadataUri,proto3" json:"metadata_uri,omitempty"`
	ShardDoubleHashes  [][]byte                                 `protobuf:"bytes,2,rep,name=shard_double_hashes,json=shardDoubleHashes,proto3" json:"shard_double_hashes,omitempty"`
	Timestamp          time.Time                                `protobuf:"bytes,3,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	Status             string                                   `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	Publisher          string                                   `protobuf:"bytes,5,opt,name=publisher,proto3" json:"publisher,omitempty"`
	Challenger         string                                   `protobuf:"bytes,6,opt,name=challenger,proto3" json:"challenger,omitempty"`
	Collateral         github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,7,rep,name=collateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"collateral"`
	ChallengeTimestamp time.Time                                `protobuf:"bytes,8,opt,name=challenge_timestamp,json=challengeTimestamp,proto3,stdtime" json:"challenge_timestamp"`
}

func (m *PublishedData) Reset()         { *m = PublishedData{} }
func (m *PublishedData) String() string { return proto.CompactTextString(m) }
func (*PublishedData) ProtoMessage()    {}
func (*PublishedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_43a893e4b5a365a1, []int{0}
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

func (m *PublishedData) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
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

type DataShares struct {
	Indices []int64  `protobuf:"varint,1,rep,packed,name=indices,proto3" json:"indices,omitempty"`
	Shares  [][]byte `protobuf:"bytes,2,rep,name=shares,proto3" json:"shares,omitempty"`
}

func (m *DataShares) Reset()         { *m = DataShares{} }
func (m *DataShares) String() string { return proto.CompactTextString(m) }
func (*DataShares) ProtoMessage()    {}
func (*DataShares) Descriptor() ([]byte, []int) {
	return fileDescriptor_43a893e4b5a365a1, []int{1}
}
func (m *DataShares) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DataShares) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DataShares.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DataShares) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataShares.Merge(m, src)
}
func (m *DataShares) XXX_Size() int {
	return m.Size()
}
func (m *DataShares) XXX_DiscardUnknown() {
	xxx_messageInfo_DataShares.DiscardUnknown(m)
}

var xxx_messageInfo_DataShares proto.InternalMessageInfo

func (m *DataShares) GetIndices() []int64 {
	if m != nil {
		return m.Indices
	}
	return nil
}

func (m *DataShares) GetShares() [][]byte {
	if m != nil {
		return m.Shares
	}
	return nil
}

type Proof struct {
	MetadataUri string `protobuf:"bytes,1,opt,name=metadata_uri,json=metadataUri,proto3" json:"metadata_uri,omitempty"`
	Sender      string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	// TODO: replace indices and shard hashes to proof
	Indices     []uint64 `protobuf:"varint,3,rep,packed,name=indices,proto3" json:"indices,omitempty"`
	ShardHashes [][]byte `protobuf:"bytes,4,rep,name=shard_hashes,json=shardHashes,proto3" json:"shard_hashes,omitempty"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_43a893e4b5a365a1, []int{2}
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

func (m *Proof) GetIndices() []uint64 {
	if m != nil {
		return m.Indices
	}
	return nil
}

func (m *Proof) GetShardHashes() [][]byte {
	if m != nil {
		return m.ShardHashes
	}
	return nil
}

func init() {
	proto.RegisterType((*PublishedData)(nil), "sunrise.da.PublishedData")
	proto.RegisterType((*DataShares)(nil), "sunrise.da.DataShares")
	proto.RegisterType((*Proof)(nil), "sunrise.da.Proof")
}

func init() { proto.RegisterFile("sunrise/da/published_data.proto", fileDescriptor_43a893e4b5a365a1) }

var fileDescriptor_43a893e4b5a365a1 = []byte{
	// 525 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xb1, 0x6e, 0xdb, 0x30,
	0x10, 0xb5, 0x22, 0xc7, 0x89, 0x69, 0x77, 0x08, 0x53, 0xb4, 0x8a, 0x51, 0xc8, 0x8a, 0x27, 0xa1,
	0x40, 0xc5, 0x3a, 0xdd, 0x3b, 0xb8, 0x1e, 0x3a, 0x06, 0x6e, 0xb3, 0x74, 0x11, 0x4e, 0x12, 0x23,
	0x11, 0x91, 0x44, 0x83, 0xa4, 0x82, 0x66, 0xea, 0x2f, 0xe4, 0x3b, 0xfa, 0x15, 0x1d, 0x33, 0x66,
	0xec, 0xd4, 0x14, 0xf6, 0x8f, 0x14, 0xa2, 0x28, 0x5b, 0x63, 0xb3, 0x88, 0x7a, 0xf7, 0x8e, 0x77,
	0xc7, 0xf7, 0x0e, 0x4d, 0x65, 0x55, 0x0a, 0x26, 0x29, 0x49, 0x80, 0xac, 0xab, 0x28, 0x67, 0x32,
	0xa3, 0x49, 0x98, 0x80, 0x82, 0x60, 0x2d, 0xb8, 0xe2, 0x18, 0x99, 0x84, 0x20, 0x81, 0xc9, 0x34,
	0xe5, 0x3c, 0xcd, 0x29, 0xd1, 0x4c, 0x54, 0x5d, 0x13, 0xc5, 0x0a, 0x2a, 0x15, 0x14, 0xeb, 0x26,
	0x79, 0x72, 0x02, 0x05, 0x2b, 0x39, 0xd1, 0x5f, 0x13, 0x7a, 0x1d, 0x73, 0x59, 0x70, 0x49, 0x0a,
	0x99, 0x92, 0xdb, 0x79, 0x7d, 0x18, 0xe2, 0xac, 0x21, 0x42, 0x8d, 0x48, 0x03, 0x0c, 0xf5, 0x32,
	0xe5, 0x29, 0x6f, 0xe2, 0xf5, 0x5f, 0x5b, 0xa9, 0x3b, 0x2a, 0x08, 0x28, 0xda, 0x74, 0xd7, 0xb4,
	0x88, 0x40, 0x52, 0x72, 0x3b, 0x8f, 0xa8, 0x82, 0x39, 0x89, 0x39, 0x2b, 0x1b, 0x7e, 0xf6, 0xcb,
	0x46, 0x2f, 0x2e, 0xdb, 0xb7, 0x2d, 0x41, 0x01, 0x3e, 0x47, 0xe3, 0x82, 0x2a, 0xa8, 0x9f, 0x19,
	0x56, 0x82, 0x39, 0x96, 0x67, 0xf9, 0xc3, 0xd5, 0xa8, 0x8d, 0x5d, 0x09, 0x86, 0x03, 0x74, 0x2a,
	0x33, 0x10, 0x49, 0x98, 0xf0, 0x2a, 0xca, 0x69, 0x98, 0x81, 0xcc, 0xa8, 0x74, 0x0e, 0x3c, 0xdb,
	0x1f, 0xaf, 0x4e, 0x34, 0xb5, 0xd4, 0xcc, 0x67, 0x4d, 0xe0, 0x05, 0x1a, 0xee, 0xd4, 0x70, 0x6c,
	0xcf, 0xf2, 0x47, 0x17, 0x93, 0xa0, 0xd1, 0x2b, 0x68, 0xf5, 0x0a, 0xbe, 0xb6, 0x19, 0x8b, 0xe3,
	0x87, 0x3f, 0xd3, 0xde, 0xfd, 0xd3, 0xd4, 0x5a, 0xed, 0xaf, 0xe1, 0x57, 0x68, 0x20, 0x15, 0xa8,
	0x4a, 0x3a, 0x7d, 0x3d, 0x90, 0x41, 0xf8, 0x0d, 0x1a, 0xb6, 0xde, 0x08, 0xe7, 0x50, 0x53, 0xfb,
	0x00, 0x76, 0x11, 0x8a, 0x33, 0xc8, 0x73, 0x5a, 0xa6, 0x54, 0x38, 0x03, 0x4d, 0x77, 0x22, 0xf8,
	0x06, 0xa1, 0x98, 0xe7, 0x39, 0x28, 0x2a, 0x20, 0x77, 0x8e, 0x3c, 0xdb, 0x1f, 0x5d, 0x9c, 0x05,
	0x46, 0xf0, 0x5a, 0xb3, 0xc0, 0x68, 0x16, 0x7c, 0xe2, 0xac, 0x5c, 0xbc, 0xaf, 0x27, 0xfb, 0xf9,
	0x34, 0xf5, 0x53, 0xa6, 0xb2, 0x2a, 0x0a, 0x62, 0x5e, 0x18, 0x77, 0xcc, 0xf1, 0x4e, 0x26, 0x37,
	0x44, 0xdd, 0xad, 0xa9, 0xd4, 0x17, 0xe4, 0xaa, 0x53, 0x1e, 0x5f, 0xa1, 0xd3, 0x5d, 0xeb, 0x70,
	0x2f, 0xc8, 0xf1, 0x33, 0x04, 0xc1, 0xbb, 0x02, 0x3b, 0x76, 0xf6, 0x11, 0xa1, 0xda, 0xb8, 0x2f,
	0x19, 0x08, 0x2a, 0xb1, 0x83, 0x8e, 0x58, 0x99, 0xb0, 0x98, 0x4a, 0xc7, 0xf2, 0x6c, 0xdf, 0x5e,
	0xb5, 0x50, 0x2b, 0xa8, 0x73, 0x8c, 0x51, 0x06, 0xcd, 0x7e, 0xa0, 0xc3, 0x4b, 0xc1, 0xf9, 0xf5,
	0xff, 0x38, 0x5f, 0xd7, 0xa0, 0x65, 0x42, 0x85, 0x73, 0x60, 0x5c, 0xd0, 0xa8, 0xdb, 0xd5, 0xf6,
	0x6c, 0xbf, 0xbf, 0xef, 0x7a, 0x8e, 0xc6, 0xcd, 0xae, 0x98, 0x25, 0xe9, 0xeb, 0xde, 0x23, 0x1d,
	0x6b, 0xd6, 0x63, 0xb1, 0x7c, 0xd8, 0xb8, 0xd6, 0xe3, 0xc6, 0xb5, 0xfe, 0x6e, 0x5c, 0xeb, 0x7e,
	0xeb, 0xf6, 0x1e, 0xb7, 0x6e, 0xef, 0xf7, 0xd6, 0xed, 0x7d, 0x7b, 0xdb, 0xd1, 0xd9, 0x6c, 0x78,
	0x0e, 0x77, 0x54, 0xb4, 0x80, 0x7c, 0xaf, 0x17, 0x5e, 0xeb, 0x1d, 0x0d, 0xb4, 0x70, 0x1f, 0xfe,
	0x05, 0x00, 0x00, 0xff, 0xff, 0xa9, 0xc8, 0x83, 0xff, 0xb6, 0x03, 0x00, 0x00,
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
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.ChallengeTimestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ChallengeTimestamp):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintPublishedData(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x42
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
			dAtA[i] = 0x3a
		}
	}
	if len(m.Challenger) > 0 {
		i -= len(m.Challenger)
		copy(dAtA[i:], m.Challenger)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Challenger)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Publisher) > 0 {
		i -= len(m.Publisher)
		copy(dAtA[i:], m.Publisher)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Publisher)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x22
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintPublishedData(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	if len(m.ShardDoubleHashes) > 0 {
		for iNdEx := len(m.ShardDoubleHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ShardDoubleHashes[iNdEx])
			copy(dAtA[i:], m.ShardDoubleHashes[iNdEx])
			i = encodeVarintPublishedData(dAtA, i, uint64(len(m.ShardDoubleHashes[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
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

func (m *DataShares) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DataShares) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DataShares) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Shares) > 0 {
		for iNdEx := len(m.Shares) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Shares[iNdEx])
			copy(dAtA[i:], m.Shares[iNdEx])
			i = encodeVarintPublishedData(dAtA, i, uint64(len(m.Shares[iNdEx])))
			i--
			dAtA[i] = 0x12
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
	if len(m.ShardHashes) > 0 {
		for iNdEx := len(m.ShardHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ShardHashes[iNdEx])
			copy(dAtA[i:], m.ShardHashes[iNdEx])
			i = encodeVarintPublishedData(dAtA, i, uint64(len(m.ShardHashes[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Indices) > 0 {
		dAtA6 := make([]byte, len(m.Indices)*10)
		var j5 int
		for _, num := range m.Indices {
			for num >= 1<<7 {
				dAtA6[j5] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j5++
			}
			dAtA6[j5] = uint8(num)
			j5++
		}
		i -= j5
		copy(dAtA[i:], dAtA6[:j5])
		i = encodeVarintPublishedData(dAtA, i, uint64(j5))
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
	if len(m.ShardDoubleHashes) > 0 {
		for _, b := range m.ShardDoubleHashes {
			l = len(b)
			n += 1 + l + sovPublishedData(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovPublishedData(uint64(l))
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovPublishedData(uint64(l))
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
	return n
}

func (m *DataShares) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Indices) > 0 {
		l = 0
		for _, e := range m.Indices {
			l += sovPublishedData(uint64(e))
		}
		n += 1 + sovPublishedData(uint64(l)) + l
	}
	if len(m.Shares) > 0 {
		for _, b := range m.Shares {
			l = len(b)
			n += 1 + l + sovPublishedData(uint64(l))
		}
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
	if len(m.ShardHashes) > 0 {
		for _, b := range m.ShardHashes {
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
		case 3:
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
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
func (m *DataShares) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: DataShares: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DataShares: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
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
			m.Shares = append(m.Shares, make([]byte, postIndex-iNdEx))
			copy(m.Shares[len(m.Shares)-1], dAtA[iNdEx:postIndex])
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
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPublishedData
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
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
					m.Indices = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPublishedData
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field ShardHashes", wireType)
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
			m.ShardHashes = append(m.ShardHashes, make([]byte, postIndex-iNdEx))
			copy(m.ShardHashes[len(m.ShardHashes)-1], dAtA[iNdEx:postIndex])
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
