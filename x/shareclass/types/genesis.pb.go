// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/shareclass/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the shareclass module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params                    Params                            `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Unbondings                []Unbonding                       `protobuf:"bytes,2,rep,name=unbondings,proto3" json:"unbondings"`
	UnbondingCount            uint64                            `protobuf:"varint,3,opt,name=unbonding_count,json=unbondingCount,proto3" json:"unbonding_count,omitempty"`
	RewardMultipliers         []GenesisRewardMultiplier         `protobuf:"bytes,4,rep,name=reward_multipliers,json=rewardMultipliers,proto3" json:"reward_multipliers"`
	UserLastRewardMultipliers []GenesisUserLastRewardMultiplier `protobuf:"bytes,5,rep,name=user_last_reward_multipliers,json=userLastRewardMultipliers,proto3" json:"user_last_reward_multipliers"`
	LastRewardHandlingTimes   []GenesisLastRewardHandlingTime   `protobuf:"bytes,6,rep,name=last_reward_handling_times,json=lastRewardHandlingTimes,proto3" json:"last_reward_handling_times"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5ff9906d61cff1d, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetUnbondings() []Unbonding {
	if m != nil {
		return m.Unbondings
	}
	return nil
}

func (m *GenesisState) GetUnbondingCount() uint64 {
	if m != nil {
		return m.UnbondingCount
	}
	return 0
}

func (m *GenesisState) GetRewardMultipliers() []GenesisRewardMultiplier {
	if m != nil {
		return m.RewardMultipliers
	}
	return nil
}

func (m *GenesisState) GetUserLastRewardMultipliers() []GenesisUserLastRewardMultiplier {
	if m != nil {
		return m.UserLastRewardMultipliers
	}
	return nil
}

func (m *GenesisState) GetLastRewardHandlingTimes() []GenesisLastRewardHandlingTime {
	if m != nil {
		return m.LastRewardHandlingTimes
	}
	return nil
}

// GenesisRewardMultiplier
type GenesisRewardMultiplier struct {
	Validator        string `protobuf:"bytes,1,opt,name=validator,proto3" json:"validator,omitempty"`
	Denom            string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	RewardMultiplier string `protobuf:"bytes,3,opt,name=reward_multiplier,json=rewardMultiplier,proto3" json:"reward_multiplier,omitempty"`
}

func (m *GenesisRewardMultiplier) Reset()         { *m = GenesisRewardMultiplier{} }
func (m *GenesisRewardMultiplier) String() string { return proto.CompactTextString(m) }
func (*GenesisRewardMultiplier) ProtoMessage()    {}
func (*GenesisRewardMultiplier) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5ff9906d61cff1d, []int{1}
}
func (m *GenesisRewardMultiplier) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisRewardMultiplier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisRewardMultiplier.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisRewardMultiplier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisRewardMultiplier.Merge(m, src)
}
func (m *GenesisRewardMultiplier) XXX_Size() int {
	return m.Size()
}
func (m *GenesisRewardMultiplier) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisRewardMultiplier.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisRewardMultiplier proto.InternalMessageInfo

func (m *GenesisRewardMultiplier) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *GenesisRewardMultiplier) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *GenesisRewardMultiplier) GetRewardMultiplier() string {
	if m != nil {
		return m.RewardMultiplier
	}
	return ""
}

// GenesisUserLastRewardMultiplier
type GenesisUserLastRewardMultiplier struct {
	User             string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Validator        string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	Denom            string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
	RewardMultiplier string `protobuf:"bytes,4,opt,name=reward_multiplier,json=rewardMultiplier,proto3" json:"reward_multiplier,omitempty"`
}

func (m *GenesisUserLastRewardMultiplier) Reset()         { *m = GenesisUserLastRewardMultiplier{} }
func (m *GenesisUserLastRewardMultiplier) String() string { return proto.CompactTextString(m) }
func (*GenesisUserLastRewardMultiplier) ProtoMessage()    {}
func (*GenesisUserLastRewardMultiplier) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5ff9906d61cff1d, []int{2}
}
func (m *GenesisUserLastRewardMultiplier) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisUserLastRewardMultiplier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisUserLastRewardMultiplier.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisUserLastRewardMultiplier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisUserLastRewardMultiplier.Merge(m, src)
}
func (m *GenesisUserLastRewardMultiplier) XXX_Size() int {
	return m.Size()
}
func (m *GenesisUserLastRewardMultiplier) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisUserLastRewardMultiplier.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisUserLastRewardMultiplier proto.InternalMessageInfo

func (m *GenesisUserLastRewardMultiplier) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *GenesisUserLastRewardMultiplier) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *GenesisUserLastRewardMultiplier) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *GenesisUserLastRewardMultiplier) GetRewardMultiplier() string {
	if m != nil {
		return m.RewardMultiplier
	}
	return ""
}

// GenesisLastRewardHandlingTime
type GenesisLastRewardHandlingTime struct {
	Validator              string `protobuf:"bytes,1,opt,name=validator,proto3" json:"validator,omitempty"`
	LastRewardHandlingTime int64  `protobuf:"varint,2,opt,name=last_reward_handling_time,json=lastRewardHandlingTime,proto3" json:"last_reward_handling_time,omitempty"`
}

func (m *GenesisLastRewardHandlingTime) Reset()         { *m = GenesisLastRewardHandlingTime{} }
func (m *GenesisLastRewardHandlingTime) String() string { return proto.CompactTextString(m) }
func (*GenesisLastRewardHandlingTime) ProtoMessage()    {}
func (*GenesisLastRewardHandlingTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5ff9906d61cff1d, []int{3}
}
func (m *GenesisLastRewardHandlingTime) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisLastRewardHandlingTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisLastRewardHandlingTime.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisLastRewardHandlingTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisLastRewardHandlingTime.Merge(m, src)
}
func (m *GenesisLastRewardHandlingTime) XXX_Size() int {
	return m.Size()
}
func (m *GenesisLastRewardHandlingTime) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisLastRewardHandlingTime.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisLastRewardHandlingTime proto.InternalMessageInfo

func (m *GenesisLastRewardHandlingTime) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *GenesisLastRewardHandlingTime) GetLastRewardHandlingTime() int64 {
	if m != nil {
		return m.LastRewardHandlingTime
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "sunrise.shareclass.v1.GenesisState")
	proto.RegisterType((*GenesisRewardMultiplier)(nil), "sunrise.shareclass.v1.GenesisRewardMultiplier")
	proto.RegisterType((*GenesisUserLastRewardMultiplier)(nil), "sunrise.shareclass.v1.GenesisUserLastRewardMultiplier")
	proto.RegisterType((*GenesisLastRewardHandlingTime)(nil), "sunrise.shareclass.v1.GenesisLastRewardHandlingTime")
}

func init() {
	proto.RegisterFile("sunrise/shareclass/v1/genesis.proto", fileDescriptor_f5ff9906d61cff1d)
}

var fileDescriptor_f5ff9906d61cff1d = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x1b, 0x37, 0x52, 0xa6, 0x08, 0xe8, 0xaa, 0x50, 0x37, 0x6a, 0x5d, 0x2b, 0x1c, 0x88,
	0x84, 0x64, 0xab, 0x01, 0x21, 0x21, 0x6e, 0x45, 0x02, 0x0e, 0x45, 0x42, 0x86, 0x5e, 0xb8, 0x44,
	0x9b, 0x78, 0xe5, 0xac, 0x64, 0x7b, 0xa3, 0x9d, 0x75, 0xda, 0x22, 0xf1, 0x0f, 0x9c, 0x90, 0xf8,
	0xa3, 0x1e, 0x7b, 0xe4, 0x84, 0x50, 0xf2, 0x23, 0xc8, 0xeb, 0x6d, 0x12, 0x8a, 0x9d, 0xdc, 0x76,
	0x67, 0xde, 0xbc, 0xf7, 0xf2, 0x66, 0x63, 0x78, 0x82, 0x79, 0x26, 0x39, 0xb2, 0x00, 0xc7, 0x54,
	0xb2, 0x51, 0x42, 0x11, 0x83, 0xe9, 0x49, 0x10, 0xb3, 0x8c, 0x21, 0x47, 0x7f, 0x22, 0x85, 0x12,
	0xe4, 0x91, 0x01, 0xf9, 0x4b, 0x90, 0x3f, 0x3d, 0xe9, 0xec, 0xc5, 0x22, 0x16, 0x1a, 0x11, 0x14,
	0xa7, 0x12, 0xdc, 0xe9, 0x56, 0x33, 0x4e, 0xa8, 0xa4, 0xa9, 0x21, 0xec, 0x78, 0xd5, 0x98, 0xa1,
	0xc8, 0xa2, 0x12, 0xd1, 0xfd, 0x69, 0xc3, 0xbd, 0x77, 0xa5, 0x89, 0x4f, 0x8a, 0x2a, 0x46, 0x5e,
	0x43, 0xab, 0xa4, 0x70, 0x2c, 0xcf, 0xea, 0xed, 0xf4, 0x8f, 0xfc, 0x4a, 0x53, 0xfe, 0x47, 0x0d,
	0x3a, 0xb5, 0xaf, 0x7f, 0x1f, 0x37, 0x42, 0x33, 0x42, 0xde, 0x02, 0xe4, 0x59, 0xc1, 0xce, 0xb3,
	0x18, 0x9d, 0x2d, 0xaf, 0xd9, 0xdb, 0xe9, 0x7b, 0x35, 0x04, 0xe7, 0xb7, 0x40, 0xc3, 0xb1, 0x32,
	0x49, 0x9e, 0xc2, 0x83, 0xc5, 0x6d, 0x30, 0x12, 0x79, 0xa6, 0x9c, 0xa6, 0x67, 0xf5, 0xec, 0xf0,
	0xfe, 0xa2, 0xfc, 0xa6, 0xa8, 0x92, 0x11, 0x10, 0xc9, 0x2e, 0xa8, 0x8c, 0x06, 0x69, 0x9e, 0x28,
	0x3e, 0x49, 0x38, 0x93, 0xe8, 0xd8, 0x5a, 0xd8, 0xaf, 0x11, 0x36, 0x3f, 0x37, 0xd4, 0x73, 0x1f,
	0x16, 0x63, 0xc6, 0xc6, 0xae, 0xbc, 0x53, 0x47, 0xf2, 0x0d, 0x0e, 0x73, 0x64, 0x72, 0x90, 0x50,
	0x54, 0x83, 0x0a, 0xb9, 0x6d, 0x2d, 0xf7, 0x72, 0xbd, 0xdc, 0x39, 0x32, 0x79, 0x46, 0x51, 0xd5,
	0xc8, 0x1e, 0xe4, 0x35, 0x7d, 0x24, 0x17, 0xd0, 0x59, 0x55, 0x1e, 0xd3, 0x2c, 0x4a, 0x8a, 0x5c,
	0x14, 0x4f, 0x19, 0x3a, 0x2d, 0x2d, 0xfe, 0x62, 0xbd, 0xf8, 0x92, 0xf8, 0xbd, 0x99, 0xfe, 0xcc,
	0x53, 0x66, 0xa4, 0xf7, 0x93, 0xca, 0x2e, 0x76, 0xbf, 0xc2, 0x7e, 0x4d, 0x56, 0xe4, 0x10, 0xda,
	0x53, 0x9a, 0xf0, 0x88, 0x2a, 0x21, 0xf5, 0x43, 0x69, 0x87, 0xcb, 0x02, 0xd9, 0x83, 0xed, 0x88,
	0x65, 0x22, 0x75, 0xb6, 0x74, 0xa7, 0xbc, 0x90, 0x67, 0xb0, 0xfb, 0x5f, 0x78, 0x7a, 0xad, 0xed,
	0xf0, 0xe1, 0xdd, 0xd0, 0xbb, 0x3f, 0x2c, 0x38, 0xde, 0x90, 0x1c, 0x21, 0x60, 0x17, 0xa9, 0x19,
	0x7d, 0x7d, 0xfe, 0xd7, 0xd8, 0x56, 0xad, 0xb1, 0xe6, 0x46, 0x63, 0x76, 0x8d, 0xb1, 0x4b, 0x38,
	0x5a, 0x1b, 0xea, 0x86, 0x68, 0x5e, 0xc1, 0x41, 0xed, 0x32, 0xb5, 0xdf, 0x66, 0xf8, 0xb8, 0x7a,
	0x1f, 0xa7, 0x67, 0xd7, 0x33, 0xd7, 0xba, 0x99, 0xb9, 0xd6, 0x9f, 0x99, 0x6b, 0x7d, 0x9f, 0xbb,
	0x8d, 0x9b, 0xb9, 0xdb, 0xf8, 0x35, 0x77, 0x1b, 0x5f, 0xfa, 0x31, 0x57, 0xe3, 0x7c, 0xe8, 0x8f,
	0x44, 0x1a, 0x98, 0x77, 0x90, 0xd0, 0x2b, 0x26, 0x6f, 0x2f, 0xc1, 0xe5, 0xea, 0x07, 0x40, 0x5d,
	0x4d, 0x18, 0x0e, 0x5b, 0xfa, 0xff, 0xff, 0xfc, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe8, 0x1a,
	0xaf, 0xcf, 0x99, 0x04, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.LastRewardHandlingTimes) > 0 {
		for iNdEx := len(m.LastRewardHandlingTimes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LastRewardHandlingTimes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.UserLastRewardMultipliers) > 0 {
		for iNdEx := len(m.UserLastRewardMultipliers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UserLastRewardMultipliers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.RewardMultipliers) > 0 {
		for iNdEx := len(m.RewardMultipliers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RewardMultipliers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.UnbondingCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.UnbondingCount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Unbondings) > 0 {
		for iNdEx := len(m.Unbondings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Unbondings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *GenesisRewardMultiplier) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisRewardMultiplier) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisRewardMultiplier) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RewardMultiplier) > 0 {
		i -= len(m.RewardMultiplier)
		copy(dAtA[i:], m.RewardMultiplier)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.RewardMultiplier)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GenesisUserLastRewardMultiplier) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisUserLastRewardMultiplier) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisUserLastRewardMultiplier) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RewardMultiplier) > 0 {
		i -= len(m.RewardMultiplier)
		copy(dAtA[i:], m.RewardMultiplier)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.RewardMultiplier)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GenesisLastRewardHandlingTime) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisLastRewardHandlingTime) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisLastRewardHandlingTime) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastRewardHandlingTime != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.LastRewardHandlingTime))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Unbondings) > 0 {
		for _, e := range m.Unbondings {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.UnbondingCount != 0 {
		n += 1 + sovGenesis(uint64(m.UnbondingCount))
	}
	if len(m.RewardMultipliers) > 0 {
		for _, e := range m.RewardMultipliers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.UserLastRewardMultipliers) > 0 {
		for _, e := range m.UserLastRewardMultipliers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LastRewardHandlingTimes) > 0 {
		for _, e := range m.LastRewardHandlingTimes {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisRewardMultiplier) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.RewardMultiplier)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *GenesisUserLastRewardMultiplier) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.RewardMultiplier)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *GenesisLastRewardHandlingTime) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.LastRewardHandlingTime != 0 {
		n += 1 + sovGenesis(uint64(m.LastRewardHandlingTime))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unbondings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Unbondings = append(m.Unbondings, Unbonding{})
			if err := m.Unbondings[len(m.Unbondings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingCount", wireType)
			}
			m.UnbondingCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondingCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardMultipliers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardMultipliers = append(m.RewardMultipliers, GenesisRewardMultiplier{})
			if err := m.RewardMultipliers[len(m.RewardMultipliers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserLastRewardMultipliers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserLastRewardMultipliers = append(m.UserLastRewardMultipliers, GenesisUserLastRewardMultiplier{})
			if err := m.UserLastRewardMultipliers[len(m.UserLastRewardMultipliers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastRewardHandlingTimes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastRewardHandlingTimes = append(m.LastRewardHandlingTimes, GenesisLastRewardHandlingTime{})
			if err := m.LastRewardHandlingTimes[len(m.LastRewardHandlingTimes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisRewardMultiplier) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisRewardMultiplier: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisRewardMultiplier: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardMultiplier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardMultiplier = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisUserLastRewardMultiplier) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisUserLastRewardMultiplier: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisUserLastRewardMultiplier: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardMultiplier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardMultiplier = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisLastRewardHandlingTime) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisLastRewardHandlingTime: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisLastRewardHandlingTime: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastRewardHandlingTime", wireType)
			}
			m.LastRewardHandlingTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastRewardHandlingTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
