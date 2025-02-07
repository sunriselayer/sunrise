// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/da/v1/params.proto

package types

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
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

// Params defines the parameters for the module.
type Params struct {
	// Invalid shard threshold required to go to STATUS_CHALLENGING.
	ChallengeThreshold string `protobuf:"bytes,1,opt,name=challenge_threshold,json=challengeThreshold,proto3" json:"challenge_threshold,omitempty"`
	// https://docs.sunriselayer.io/learn/sunrise/data-availability#the-condition-of-data-availability
	ReplicationFactor string `protobuf:"bytes,2,opt,name=replication_factor,json=replicationFactor,proto3" json:"replication_factor,omitempty"`
	// How many blocks of slash are done every
	SlashEpoch uint64 `protobuf:"varint,3,opt,name=slash_epoch,json=slashEpoch,proto3" json:"slash_epoch,omitempty"`
	// (number of challenges a validator did not submit proof / number of all challenge) is over this threshold in an epoch
	// that validator will be slashed
	SlashFaultThreshold string `protobuf:"bytes,4,opt,name=slash_fault_threshold,json=slashFaultThreshold,proto3" json:"slash_fault_threshold,omitempty"`
	// voting power deducted during slash
	SlashFraction              string                                   `protobuf:"bytes,5,opt,name=slash_fraction,json=slashFraction,proto3" json:"slash_fraction,omitempty"`
	ChallengePeriod            time.Duration                            `protobuf:"bytes,6,opt,name=challenge_period,json=challengePeriod,proto3,stdduration" json:"challenge_period"`
	ProofPeriod                time.Duration                            `protobuf:"bytes,7,opt,name=proof_period,json=proofPeriod,proto3,stdduration" json:"proof_period"`
	RejectedRemovalPeriod      time.Duration                            `protobuf:"bytes,8,opt,name=rejected_removal_period,json=rejectedRemovalPeriod,proto3,stdduration" json:"rejected_removal_period"`
	PublishDataCollateral      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,9,rep,name=publish_data_collateral,json=publishDataCollateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"publish_data_collateral"`
	SubmitInvalidityCollateral github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,10,rep,name=submit_invalidity_collateral,json=submitInvalidityCollateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"submit_invalidity_collateral"`
	ZkpVerifyingKey            []byte                                   `protobuf:"bytes,11,opt,name=zkp_verifying_key,json=zkpVerifyingKey,proto3" json:"zkp_verifying_key,omitempty"`
	// proving key used in sunrise-data
	ZkpProvingKey []byte `protobuf:"bytes,12,opt,name=zkp_proving_key,json=zkpProvingKey,proto3" json:"zkp_proving_key,omitempty"`
	// min_shard_count used in sunrise-data
	MinShardCount uint64 `protobuf:"varint,13,opt,name=min_shard_count,json=minShardCount,proto3" json:"min_shard_count,omitempty"`
	// max_shard_count used in sunrise-data
	MaxShardCount uint64 `protobuf:"varint,14,opt,name=max_shard_count,json=maxShardCount,proto3" json:"max_shard_count,omitempty"`
	// max_shard_size used in sunrise-data
	MaxShardSize uint64 `protobuf:"varint,15,opt,name=max_shard_size,json=maxShardSize,proto3" json:"max_shard_size,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_70863a5b1732e9c8, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetChallengeThreshold() string {
	if m != nil {
		return m.ChallengeThreshold
	}
	return ""
}

func (m *Params) GetReplicationFactor() string {
	if m != nil {
		return m.ReplicationFactor
	}
	return ""
}

func (m *Params) GetSlashEpoch() uint64 {
	if m != nil {
		return m.SlashEpoch
	}
	return 0
}

func (m *Params) GetSlashFaultThreshold() string {
	if m != nil {
		return m.SlashFaultThreshold
	}
	return ""
}

func (m *Params) GetSlashFraction() string {
	if m != nil {
		return m.SlashFraction
	}
	return ""
}

func (m *Params) GetChallengePeriod() time.Duration {
	if m != nil {
		return m.ChallengePeriod
	}
	return 0
}

func (m *Params) GetProofPeriod() time.Duration {
	if m != nil {
		return m.ProofPeriod
	}
	return 0
}

func (m *Params) GetRejectedRemovalPeriod() time.Duration {
	if m != nil {
		return m.RejectedRemovalPeriod
	}
	return 0
}

func (m *Params) GetPublishDataCollateral() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.PublishDataCollateral
	}
	return nil
}

func (m *Params) GetSubmitInvalidityCollateral() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.SubmitInvalidityCollateral
	}
	return nil
}

func (m *Params) GetZkpVerifyingKey() []byte {
	if m != nil {
		return m.ZkpVerifyingKey
	}
	return nil
}

func (m *Params) GetZkpProvingKey() []byte {
	if m != nil {
		return m.ZkpProvingKey
	}
	return nil
}

func (m *Params) GetMinShardCount() uint64 {
	if m != nil {
		return m.MinShardCount
	}
	return 0
}

func (m *Params) GetMaxShardCount() uint64 {
	if m != nil {
		return m.MaxShardCount
	}
	return 0
}

func (m *Params) GetMaxShardSize() uint64 {
	if m != nil {
		return m.MaxShardSize
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "sunrise.da.v1.Params")
}

func init() { proto.RegisterFile("sunrise/da/v1/params.proto", fileDescriptor_70863a5b1732e9c8) }

var fileDescriptor_70863a5b1732e9c8 = []byte{
	// 652 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x4f, 0x6b, 0xd4, 0x40,
	0x18, 0xc6, 0x37, 0xb6, 0xf6, 0xcf, 0xec, 0x9f, 0xda, 0xd4, 0xd2, 0x74, 0x91, 0xec, 0x22, 0x22,
	0x4b, 0xa1, 0x89, 0xab, 0x78, 0x11, 0x44, 0xd8, 0xae, 0x05, 0x11, 0xa4, 0x6c, 0xc5, 0x83, 0x1e,
	0xc2, 0x24, 0x99, 0x4d, 0xc6, 0x9d, 0x64, 0xc2, 0xcc, 0x24, 0x74, 0xf7, 0xea, 0x59, 0xf0, 0xe8,
	0xd1, 0xb3, 0x67, 0x3f, 0x44, 0x8f, 0xc5, 0x93, 0x27, 0x2b, 0xed, 0xc5, 0x6f, 0xa1, 0x64, 0x32,
	0xd9, 0x46, 0xd8, 0x43, 0x0f, 0x9e, 0x76, 0xe7, 0x99, 0xdf, 0xfb, 0x3c, 0x2f, 0x2f, 0x6f, 0x06,
	0xb4, 0x79, 0x1a, 0x33, 0xcc, 0x91, 0xed, 0x43, 0x3b, 0xeb, 0xdb, 0x09, 0x64, 0x30, 0xe2, 0x56,
	0xc2, 0xa8, 0xa0, 0x7a, 0x53, 0xdd, 0x59, 0x3e, 0xb4, 0xb2, 0x7e, 0xdb, 0xf4, 0x28, 0x8f, 0x28,
	0xb7, 0x5d, 0xc8, 0x91, 0x9d, 0xf5, 0x5d, 0x24, 0x60, 0xdf, 0xf6, 0x28, 0x8e, 0x0b, 0xbc, 0xbd,
	0x5b, 0xdc, 0x3b, 0xf2, 0x64, 0x17, 0x07, 0x75, 0x75, 0x3b, 0xa0, 0x01, 0x2d, 0xf4, 0xfc, 0x9f,
	0x52, 0xcd, 0x80, 0xd2, 0x80, 0x20, 0x5b, 0x9e, 0xdc, 0x74, 0x6c, 0xfb, 0x29, 0x83, 0x02, 0x53,
	0x65, 0x78, 0xf7, 0xcf, 0x2a, 0x58, 0x39, 0x92, 0x0d, 0xe9, 0xcf, 0xc0, 0x96, 0x17, 0x42, 0x42,
	0x50, 0x1c, 0x20, 0x47, 0x84, 0x0c, 0xf1, 0x90, 0x12, 0xdf, 0xd0, 0xba, 0x5a, 0x6f, 0x7d, 0xd0,
	0xfa, 0xfe, 0x6d, 0x1f, 0xa8, 0xbc, 0x21, 0xf2, 0x46, 0xfa, 0x1c, 0x7d, 0x5d, 0x92, 0xfa, 0x53,
	0xa0, 0x33, 0x94, 0x10, 0xec, 0xc9, 0x00, 0x67, 0x0c, 0x3d, 0x41, 0x99, 0x71, 0x63, 0x61, 0xfd,
	0x66, 0x85, 0x3c, 0x94, 0xa0, 0xde, 0x01, 0x75, 0x4e, 0x20, 0x0f, 0x1d, 0x94, 0x50, 0x2f, 0x34,
	0x96, 0xba, 0x5a, 0x6f, 0x79, 0x04, 0xa4, 0xf4, 0x3c, 0x57, 0xf4, 0x01, 0xd8, 0x2e, 0x80, 0x31,
	0x4c, 0x89, 0xa8, 0xb4, 0xb8, 0xbc, 0x30, 0x62, 0x4b, 0xc2, 0x87, 0x39, 0x7b, 0xd5, 0xe3, 0x63,
	0xd0, 0x52, 0x1e, 0x0c, 0x7a, 0x79, 0xb8, 0x71, 0x73, 0x61, 0x71, 0xb3, 0x28, 0x56, 0x90, 0xfe,
	0x0a, 0xdc, 0xba, 0x9a, 0x4d, 0x82, 0x18, 0xa6, 0xbe, 0xb1, 0xd2, 0xd5, 0x7a, 0xf5, 0x87, 0xbb,
	0x56, 0x31, 0x61, 0xab, 0x9c, 0xb0, 0x35, 0x54, 0x13, 0x1e, 0xac, 0x9d, 0xfe, 0xec, 0xd4, 0x3e,
	0x9f, 0x77, 0xb4, 0xd1, 0xc6, 0xbc, 0xf8, 0x48, 0xd6, 0xea, 0x87, 0xa0, 0x91, 0x30, 0x4a, 0xc7,
	0xa5, 0xd7, 0xea, 0xf5, 0xbd, 0xea, 0xb2, 0x50, 0xf9, 0xbc, 0x03, 0x3b, 0x0c, 0xbd, 0x47, 0x9e,
	0x40, 0xbe, 0xc3, 0x50, 0x44, 0x33, 0x48, 0x4a, 0xcb, 0xb5, 0xeb, 0x5b, 0x6e, 0x97, 0x1e, 0xa3,
	0xc2, 0x42, 0x99, 0x7f, 0xd0, 0xc0, 0x4e, 0x92, 0xba, 0x04, 0xf3, 0xd0, 0xf1, 0xa1, 0x80, 0x8e,
	0x47, 0x09, 0x81, 0x02, 0x31, 0x48, 0x8c, 0xf5, 0xee, 0x92, 0x74, 0x57, 0x23, 0xcb, 0xf7, 0xd5,
	0x52, 0xfb, 0x6a, 0x1d, 0x50, 0x1c, 0x0f, 0x1e, 0xe4, 0xee, 0x5f, 0xcf, 0x3b, 0xbd, 0x00, 0x8b,
	0x30, 0x75, 0x2d, 0x8f, 0x46, 0x6a, 0x5f, 0xd5, 0xcf, 0x3e, 0xf7, 0x27, 0xb6, 0x98, 0x26, 0x88,
	0xcb, 0x02, 0x3e, 0xda, 0x56, 0x59, 0x43, 0x28, 0xe0, 0xc1, 0x3c, 0x49, 0xff, 0xa8, 0x81, 0x3b,
	0x3c, 0x75, 0x23, 0x2c, 0x1c, 0x1c, 0x67, 0x90, 0x60, 0x1f, 0x8b, 0x69, 0xb5, 0x15, 0xf0, 0xff,
	0x5b, 0x69, 0x17, 0x81, 0x2f, 0xe6, 0x79, 0x95, 0x7e, 0xf6, 0xc0, 0xe6, 0x6c, 0x92, 0x38, 0x19,
	0x62, 0x78, 0x3c, 0xc5, 0x71, 0xe0, 0x4c, 0xd0, 0xd4, 0xa8, 0x77, 0xb5, 0x5e, 0x63, 0xb4, 0x31,
	0x9b, 0x24, 0x6f, 0x4a, 0xfd, 0x25, 0x9a, 0xea, 0xf7, 0x41, 0x2e, 0xe5, 0x5f, 0x6b, 0x56, 0x92,
	0x0d, 0x49, 0x36, 0x67, 0x93, 0xe4, 0xa8, 0x50, 0x15, 0x17, 0xe1, 0xd8, 0xe1, 0x21, 0x64, 0xbe,
	0xe3, 0xd1, 0x34, 0x16, 0x46, 0x53, 0xae, 0x7f, 0x33, 0xc2, 0xf1, 0x71, 0xae, 0x1e, 0xe4, 0xa2,
	0xe4, 0xe0, 0xc9, 0x3f, 0x5c, 0x4b, 0x71, 0xf0, 0xa4, 0xc2, 0xdd, 0x03, 0xad, 0x2b, 0x8e, 0xe3,
	0x19, 0x32, 0x36, 0x24, 0xd6, 0x28, 0xb1, 0x63, 0x3c, 0x43, 0x4f, 0x96, 0x7f, 0x7f, 0xe9, 0x68,
	0x83, 0xe1, 0xe9, 0x85, 0xa9, 0x9d, 0x5d, 0x98, 0xda, 0xaf, 0x0b, 0x53, 0xfb, 0x74, 0x69, 0xd6,
	0xce, 0x2e, 0xcd, 0xda, 0x8f, 0x4b, 0xb3, 0xf6, 0x76, 0xaf, 0x32, 0x2f, 0xf5, 0x4c, 0x11, 0x38,
	0x45, 0xac, 0x3c, 0xd8, 0x27, 0xf9, 0x8b, 0x26, 0xe7, 0xe6, 0xae, 0xc8, 0xfd, 0x7a, 0xf4, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x87, 0xa8, 0xc6, 0xb4, 0xec, 0x04, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ChallengeThreshold != that1.ChallengeThreshold {
		return false
	}
	if this.ReplicationFactor != that1.ReplicationFactor {
		return false
	}
	if this.SlashEpoch != that1.SlashEpoch {
		return false
	}
	if this.SlashFaultThreshold != that1.SlashFaultThreshold {
		return false
	}
	if this.SlashFraction != that1.SlashFraction {
		return false
	}
	if this.ChallengePeriod != that1.ChallengePeriod {
		return false
	}
	if this.ProofPeriod != that1.ProofPeriod {
		return false
	}
	if this.RejectedRemovalPeriod != that1.RejectedRemovalPeriod {
		return false
	}
	if len(this.PublishDataCollateral) != len(that1.PublishDataCollateral) {
		return false
	}
	for i := range this.PublishDataCollateral {
		if !this.PublishDataCollateral[i].Equal(&that1.PublishDataCollateral[i]) {
			return false
		}
	}
	if len(this.SubmitInvalidityCollateral) != len(that1.SubmitInvalidityCollateral) {
		return false
	}
	for i := range this.SubmitInvalidityCollateral {
		if !this.SubmitInvalidityCollateral[i].Equal(&that1.SubmitInvalidityCollateral[i]) {
			return false
		}
	}
	if !bytes.Equal(this.ZkpVerifyingKey, that1.ZkpVerifyingKey) {
		return false
	}
	if !bytes.Equal(this.ZkpProvingKey, that1.ZkpProvingKey) {
		return false
	}
	if this.MinShardCount != that1.MinShardCount {
		return false
	}
	if this.MaxShardCount != that1.MaxShardCount {
		return false
	}
	if this.MaxShardSize != that1.MaxShardSize {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxShardSize != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxShardSize))
		i--
		dAtA[i] = 0x78
	}
	if m.MaxShardCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxShardCount))
		i--
		dAtA[i] = 0x70
	}
	if m.MinShardCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MinShardCount))
		i--
		dAtA[i] = 0x68
	}
	if len(m.ZkpProvingKey) > 0 {
		i -= len(m.ZkpProvingKey)
		copy(dAtA[i:], m.ZkpProvingKey)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ZkpProvingKey)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.ZkpVerifyingKey) > 0 {
		i -= len(m.ZkpVerifyingKey)
		copy(dAtA[i:], m.ZkpVerifyingKey)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ZkpVerifyingKey)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.SubmitInvalidityCollateral) > 0 {
		for iNdEx := len(m.SubmitInvalidityCollateral) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SubmitInvalidityCollateral[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x52
		}
	}
	if len(m.PublishDataCollateral) > 0 {
		for iNdEx := len(m.PublishDataCollateral) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PublishDataCollateral[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.RejectedRemovalPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.RejectedRemovalPeriod):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x42
	n2, err2 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ProofPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ProofPeriod):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x3a
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ChallengePeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ChallengePeriod):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintParams(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x32
	if len(m.SlashFraction) > 0 {
		i -= len(m.SlashFraction)
		copy(dAtA[i:], m.SlashFraction)
		i = encodeVarintParams(dAtA, i, uint64(len(m.SlashFraction)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SlashFaultThreshold) > 0 {
		i -= len(m.SlashFaultThreshold)
		copy(dAtA[i:], m.SlashFaultThreshold)
		i = encodeVarintParams(dAtA, i, uint64(len(m.SlashFaultThreshold)))
		i--
		dAtA[i] = 0x22
	}
	if m.SlashEpoch != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SlashEpoch))
		i--
		dAtA[i] = 0x18
	}
	if len(m.ReplicationFactor) > 0 {
		i -= len(m.ReplicationFactor)
		copy(dAtA[i:], m.ReplicationFactor)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ReplicationFactor)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ChallengeThreshold) > 0 {
		i -= len(m.ChallengeThreshold)
		copy(dAtA[i:], m.ChallengeThreshold)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ChallengeThreshold)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChallengeThreshold)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.ReplicationFactor)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.SlashEpoch != 0 {
		n += 1 + sovParams(uint64(m.SlashEpoch))
	}
	l = len(m.SlashFaultThreshold)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.SlashFraction)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ChallengePeriod)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ProofPeriod)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.RejectedRemovalPeriod)
	n += 1 + l + sovParams(uint64(l))
	if len(m.PublishDataCollateral) > 0 {
		for _, e := range m.PublishDataCollateral {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.SubmitInvalidityCollateral) > 0 {
		for _, e := range m.SubmitInvalidityCollateral {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	l = len(m.ZkpVerifyingKey)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.ZkpProvingKey)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.MinShardCount != 0 {
		n += 1 + sovParams(uint64(m.MinShardCount))
	}
	if m.MaxShardCount != 0 {
		n += 1 + sovParams(uint64(m.MaxShardCount))
	}
	if m.MaxShardSize != 0 {
		n += 1 + sovParams(uint64(m.MaxShardSize))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChallengeThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChallengeThreshold = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicationFactor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReplicationFactor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashEpoch", wireType)
			}
			m.SlashEpoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SlashEpoch |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFaultThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SlashFaultThreshold = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFraction", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SlashFraction = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChallengePeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.ChallengePeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProofPeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.ProofPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectedRemovalPeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.RejectedRemovalPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublishDataCollateral", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PublishDataCollateral = append(m.PublishDataCollateral, types.Coin{})
			if err := m.PublishDataCollateral[len(m.PublishDataCollateral)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubmitInvalidityCollateral", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubmitInvalidityCollateral = append(m.SubmitInvalidityCollateral, types.Coin{})
			if err := m.SubmitInvalidityCollateral[len(m.SubmitInvalidityCollateral)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZkpVerifyingKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ZkpVerifyingKey = append(m.ZkpVerifyingKey[:0], dAtA[iNdEx:postIndex]...)
			if m.ZkpVerifyingKey == nil {
				m.ZkpVerifyingKey = []byte{}
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZkpProvingKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ZkpProvingKey = append(m.ZkpProvingKey[:0], dAtA[iNdEx:postIndex]...)
			if m.ZkpProvingKey == nil {
				m.ZkpProvingKey = []byte{}
			}
			iNdEx = postIndex
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinShardCount", wireType)
			}
			m.MinShardCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinShardCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxShardCount", wireType)
			}
			m.MaxShardCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxShardCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxShardSize", wireType)
			}
			m.MaxShardSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxShardSize |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
