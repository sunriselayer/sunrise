// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/da/v0/params.proto

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
	PublishDataGas uint64 `protobuf:"varint,1,opt,name=publish_data_gas,json=publishDataGas,proto3" json:"publish_data_gas,omitempty"`
	// Invalid shard threshold required to go to STATUS_CHALLENGING.
	ChallengeThreshold string `protobuf:"bytes,2,opt,name=challenge_threshold,json=challengeThreshold,proto3" json:"challenge_threshold,omitempty"`
	// https://docs.sunriselayer.io/learn/sunrise/data-availability#the-condition-of-data-availability
	ReplicationFactor string `protobuf:"bytes,3,opt,name=replication_factor,json=replicationFactor,proto3" json:"replication_factor,omitempty"`
	// How many blocks of slash are done every
	SlashEpoch uint64 `protobuf:"varint,4,opt,name=slash_epoch,json=slashEpoch,proto3" json:"slash_epoch,omitempty"`
	// (number of challenges a validator did not submit proof / number of all challenge) is over this threshold in an epoch
	// that validator will be slashed
	SlashFaultThreshold string `protobuf:"bytes,5,opt,name=slash_fault_threshold,json=slashFaultThreshold,proto3" json:"slash_fault_threshold,omitempty"`
	// voting power deducted during slash
	SlashFraction              string                                   `protobuf:"bytes,6,opt,name=slash_fraction,json=slashFraction,proto3" json:"slash_fraction,omitempty"`
	ChallengePeriod            time.Duration                            `protobuf:"bytes,7,opt,name=challenge_period,json=challengePeriod,proto3,stdduration" json:"challenge_period"`
	ProofPeriod                time.Duration                            `protobuf:"bytes,8,opt,name=proof_period,json=proofPeriod,proto3,stdduration" json:"proof_period"`
	RejectedRemovalPeriod      time.Duration                            `protobuf:"bytes,9,opt,name=rejected_removal_period,json=rejectedRemovalPeriod,proto3,stdduration" json:"rejected_removal_period"`
	VerifiedRemovalPeriod      time.Duration                            `protobuf:"bytes,10,opt,name=verified_removal_period,json=verifiedRemovalPeriod,proto3,stdduration" json:"verified_removal_period"`
	PublishDataCollateral      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,11,rep,name=publish_data_collateral,json=publishDataCollateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"publish_data_collateral"`
	SubmitInvalidityCollateral github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,12,rep,name=submit_invalidity_collateral,json=submitInvalidityCollateral,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"submit_invalidity_collateral"`
	ZkpVerifyingKey            []byte                                   `protobuf:"bytes,13,opt,name=zkp_verifying_key,json=zkpVerifyingKey,proto3" json:"zkp_verifying_key,omitempty"`
	// proving key used in sunrise-data
	ZkpProvingKey []byte `protobuf:"bytes,14,opt,name=zkp_proving_key,json=zkpProvingKey,proto3" json:"zkp_proving_key,omitempty"`
	// min_shard_count used in sunrise-data
	MinShardCount uint64 `protobuf:"varint,15,opt,name=min_shard_count,json=minShardCount,proto3" json:"min_shard_count,omitempty"`
	// max_shard_count used in sunrise-data
	MaxShardCount uint64 `protobuf:"varint,16,opt,name=max_shard_count,json=maxShardCount,proto3" json:"max_shard_count,omitempty"`
	// max_shard_size used in sunrise-data
	MaxShardSize uint64 `protobuf:"varint,17,opt,name=max_shard_size,json=maxShardSize,proto3" json:"max_shard_size,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_be6249d5e253afb0, []int{0}
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

func (m *Params) GetPublishDataGas() uint64 {
	if m != nil {
		return m.PublishDataGas
	}
	return 0
}

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

func (m *Params) GetVerifiedRemovalPeriod() time.Duration {
	if m != nil {
		return m.VerifiedRemovalPeriod
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
	proto.RegisterType((*Params)(nil), "sunrise.da.v0.Params")
}

func init() { proto.RegisterFile("sunrise/da/v0/params.proto", fileDescriptor_be6249d5e253afb0) }

var fileDescriptor_be6249d5e253afb0 = []byte{
	// 687 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcf, 0x6b, 0xd4, 0x4c,
	0x1c, 0xc6, 0x37, 0x6f, 0xfb, 0xf6, 0x6d, 0x67, 0x7f, 0xb4, 0x4d, 0xdf, 0xd2, 0x74, 0x79, 0xc9,
	0x2e, 0x2f, 0x22, 0x4b, 0xa1, 0x49, 0xab, 0x78, 0x11, 0x44, 0xd8, 0xae, 0x15, 0x11, 0xa4, 0x6c,
	0xc5, 0x83, 0x1e, 0xc2, 0x24, 0x99, 0x4d, 0xc6, 0x9d, 0x64, 0xc2, 0xcc, 0x24, 0x74, 0xf7, 0x2a,
	0x78, 0x13, 0x3c, 0x7a, 0xf4, 0xec, 0xd9, 0x3f, 0xa2, 0xc7, 0xe2, 0xc9, 0x93, 0x95, 0xf6, 0xe2,
	0x9f, 0x21, 0x99, 0x4c, 0x76, 0x53, 0xec, 0xa1, 0x82, 0xa7, 0x64, 0x9e, 0xf9, 0x7c, 0x9f, 0xef,
	0xc3, 0xfc, 0x02, 0x6d, 0x9e, 0xc6, 0x0c, 0x73, 0x64, 0xfb, 0xd0, 0xce, 0xf6, 0xec, 0x04, 0x32,
	0x18, 0x71, 0x2b, 0x61, 0x54, 0x50, 0xbd, 0xa9, 0xe6, 0x2c, 0x1f, 0x5a, 0xd9, 0x5e, 0xdb, 0xf4,
	0x28, 0x8f, 0x28, 0xb7, 0x5d, 0xc8, 0x91, 0x9d, 0xed, 0xbb, 0x48, 0xc0, 0x7d, 0xdb, 0xa3, 0x38,
	0x2e, 0xf0, 0xf6, 0x76, 0x31, 0xef, 0xc8, 0x91, 0x5d, 0x0c, 0xd4, 0xd4, 0xbf, 0x01, 0x0d, 0x68,
	0xa1, 0xe7, 0x7f, 0x4a, 0x35, 0x03, 0x4a, 0x03, 0x82, 0x6c, 0x39, 0x72, 0xd3, 0x91, 0xed, 0xa7,
	0x0c, 0x0a, 0x4c, 0x95, 0xe1, 0xff, 0x6f, 0x57, 0xc0, 0xd2, 0x91, 0x0c, 0xa4, 0xf7, 0xc0, 0x5a,
	0x92, 0xba, 0x04, 0xf3, 0xd0, 0xf1, 0xa1, 0x80, 0x4e, 0x00, 0xb9, 0xa1, 0x75, 0xb5, 0xde, 0xe2,
	0xb0, 0xa5, 0xf4, 0x01, 0x14, 0xf0, 0x31, 0xe4, 0xfa, 0x43, 0xb0, 0xe1, 0x85, 0x90, 0x10, 0x14,
	0x07, 0xc8, 0x11, 0x21, 0x43, 0x3c, 0xa4, 0xc4, 0x37, 0xfe, 0xea, 0x6a, 0xbd, 0x95, 0x7e, 0xeb,
	0xcb, 0xe7, 0x5d, 0xa0, 0x92, 0x0d, 0x90, 0x37, 0xd4, 0x67, 0xe8, 0xf3, 0x92, 0xd4, 0x1f, 0x00,
	0x9d, 0xa1, 0x84, 0x60, 0x4f, 0x46, 0x71, 0x46, 0xd0, 0x13, 0x94, 0x19, 0x0b, 0xd7, 0xd6, 0xaf,
	0x57, 0xc8, 0x43, 0x09, 0xea, 0x1d, 0x50, 0xe7, 0x04, 0xf2, 0xd0, 0x41, 0x09, 0xf5, 0x42, 0x63,
	0x51, 0x86, 0x04, 0x52, 0x7a, 0x94, 0x2b, 0x7a, 0x1f, 0x6c, 0x16, 0xc0, 0x08, 0xa6, 0x44, 0x54,
	0x22, 0xfe, 0x7d, 0x6d, 0x8b, 0x0d, 0x09, 0x1f, 0xe6, 0xec, 0x3c, 0xe3, 0x3d, 0xd0, 0x52, 0x1e,
	0x0c, 0x7a, 0x79, 0x73, 0x63, 0xe9, 0xda, 0xe2, 0x66, 0x51, 0xac, 0x20, 0xfd, 0x19, 0x58, 0x9b,
	0xaf, 0x4d, 0x82, 0x18, 0xa6, 0xbe, 0xf1, 0x4f, 0x57, 0xeb, 0xd5, 0xef, 0x6c, 0x5b, 0xc5, 0x5e,
	0x58, 0xe5, 0x5e, 0x58, 0x03, 0xb5, 0x17, 0xfd, 0xe5, 0xd3, 0x6f, 0x9d, 0xda, 0x87, 0xf3, 0x8e,
	0x36, 0x5c, 0x9d, 0x15, 0x1f, 0xc9, 0x5a, 0xfd, 0x10, 0x34, 0x12, 0x46, 0xe9, 0xa8, 0xf4, 0x5a,
	0xbe, 0xb9, 0x57, 0x5d, 0x16, 0x2a, 0x9f, 0x57, 0x60, 0x8b, 0xa1, 0xd7, 0xc8, 0x13, 0xc8, 0x77,
	0x18, 0x8a, 0x68, 0x06, 0x49, 0x69, 0xb9, 0x72, 0x73, 0xcb, 0xcd, 0xd2, 0x63, 0x58, 0x58, 0xcc,
	0xcd, 0x33, 0xc4, 0xf0, 0x08, 0xff, 0x6a, 0x0e, 0x7e, 0xc3, 0xbc, 0xf4, 0xb8, 0x6a, 0xfe, 0x46,
	0x03, 0x5b, 0x57, 0x0e, 0xa6, 0x47, 0x09, 0x81, 0x02, 0x31, 0x48, 0x8c, 0x7a, 0x77, 0x41, 0xba,
	0xab, 0xfd, 0xc8, 0xaf, 0x8d, 0xa5, 0xae, 0x8d, 0x75, 0x40, 0x71, 0xdc, 0xdf, 0xcb, 0xdd, 0x3f,
	0x9d, 0x77, 0x7a, 0x01, 0x16, 0x61, 0xea, 0x5a, 0x1e, 0x8d, 0xd4, 0xb5, 0x51, 0x9f, 0x5d, 0xee,
	0x8f, 0x6d, 0x31, 0x49, 0x10, 0x97, 0x05, 0x7c, 0xb8, 0x59, 0x39, 0xec, 0x07, 0xb3, 0x4e, 0xfa,
	0x3b, 0x0d, 0xfc, 0xc7, 0x53, 0x37, 0xc2, 0xc2, 0xc1, 0x71, 0x06, 0x09, 0xf6, 0xb1, 0x98, 0x54,
	0xa3, 0x34, 0xfe, 0x7c, 0x94, 0x76, 0xd1, 0xf0, 0xc9, 0xac, 0x5f, 0x25, 0xcf, 0x0e, 0x58, 0x9f,
	0x8e, 0x13, 0x47, 0x2e, 0xd9, 0x04, 0xc7, 0x81, 0x33, 0x46, 0x13, 0xa3, 0xd9, 0xd5, 0x7a, 0x8d,
	0xe1, 0xea, 0x74, 0x9c, 0xbc, 0x28, 0xf5, 0xa7, 0x68, 0xa2, 0xdf, 0x06, 0xb9, 0x94, 0x3f, 0x1a,
	0x59, 0x49, 0xb6, 0x24, 0xd9, 0x9c, 0x8e, 0x93, 0xa3, 0x42, 0x55, 0x5c, 0x84, 0x63, 0x87, 0x87,
	0x90, 0xf9, 0x8e, 0x47, 0xd3, 0x58, 0x18, 0xab, 0xf2, 0x6e, 0x35, 0x23, 0x1c, 0x1f, 0xe7, 0xea,
	0x41, 0x2e, 0x4a, 0x0e, 0x9e, 0x5c, 0xe1, 0xd6, 0x14, 0x07, 0x4f, 0x2a, 0xdc, 0x2d, 0xd0, 0x9a,
	0x73, 0x1c, 0x4f, 0x91, 0xb1, 0x2e, 0xb1, 0x46, 0x89, 0x1d, 0xe3, 0x29, 0xba, 0xbf, 0xf8, 0xe3,
	0x63, 0x47, 0xeb, 0x0f, 0x4e, 0x2f, 0x4c, 0xed, 0xec, 0xc2, 0xd4, 0xbe, 0x5f, 0x98, 0xda, 0xfb,
	0x4b, 0xb3, 0x76, 0x76, 0x69, 0xd6, 0xbe, 0x5e, 0x9a, 0xb5, 0x97, 0x3b, 0x95, 0xf5, 0x52, 0xaf,
	0x25, 0x81, 0x13, 0xc4, 0xca, 0x81, 0x7d, 0x92, 0x3f, 0xac, 0x72, 0xdd, 0xdc, 0x25, 0x79, 0xbe,
	0xee, 0xfe, 0x0c, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x4c, 0xb6, 0xf6, 0x73, 0x05, 0x00, 0x00,
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
	if this.PublishDataGas != that1.PublishDataGas {
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
	if this.VerifiedRemovalPeriod != that1.VerifiedRemovalPeriod {
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
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	if m.MaxShardCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxShardCount))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x80
	}
	if m.MinShardCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MinShardCount))
		i--
		dAtA[i] = 0x78
	}
	if len(m.ZkpProvingKey) > 0 {
		i -= len(m.ZkpProvingKey)
		copy(dAtA[i:], m.ZkpProvingKey)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ZkpProvingKey)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.ZkpVerifyingKey) > 0 {
		i -= len(m.ZkpVerifyingKey)
		copy(dAtA[i:], m.ZkpVerifyingKey)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ZkpVerifyingKey)))
		i--
		dAtA[i] = 0x6a
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
			dAtA[i] = 0x62
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
			dAtA[i] = 0x5a
		}
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.VerifiedRemovalPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.VerifiedRemovalPeriod):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x52
	n2, err2 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.RejectedRemovalPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.RejectedRemovalPeriod):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x4a
	n3, err3 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ProofPeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ProofPeriod):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintParams(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x42
	n4, err4 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ChallengePeriod, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ChallengePeriod):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintParams(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x3a
	if len(m.SlashFraction) > 0 {
		i -= len(m.SlashFraction)
		copy(dAtA[i:], m.SlashFraction)
		i = encodeVarintParams(dAtA, i, uint64(len(m.SlashFraction)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.SlashFaultThreshold) > 0 {
		i -= len(m.SlashFaultThreshold)
		copy(dAtA[i:], m.SlashFaultThreshold)
		i = encodeVarintParams(dAtA, i, uint64(len(m.SlashFaultThreshold)))
		i--
		dAtA[i] = 0x2a
	}
	if m.SlashEpoch != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SlashEpoch))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ReplicationFactor) > 0 {
		i -= len(m.ReplicationFactor)
		copy(dAtA[i:], m.ReplicationFactor)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ReplicationFactor)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ChallengeThreshold) > 0 {
		i -= len(m.ChallengeThreshold)
		copy(dAtA[i:], m.ChallengeThreshold)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ChallengeThreshold)))
		i--
		dAtA[i] = 0x12
	}
	if m.PublishDataGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.PublishDataGas))
		i--
		dAtA[i] = 0x8
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
	if m.PublishDataGas != 0 {
		n += 1 + sovParams(uint64(m.PublishDataGas))
	}
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
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.VerifiedRemovalPeriod)
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
		n += 2 + sovParams(uint64(m.MaxShardCount))
	}
	if m.MaxShardSize != 0 {
		n += 2 + sovParams(uint64(m.MaxShardSize))
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublishDataGas", wireType)
			}
			m.PublishDataGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PublishDataGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
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
		case 3:
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
		case 4:
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
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
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
		case 9:
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
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VerifiedRemovalPeriod", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.VerifiedRemovalPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
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
		case 12:
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
		case 13:
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
		case 14:
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
		case 15:
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
		case 16:
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
		case 17:
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
