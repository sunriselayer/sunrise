// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/lockup/v1/lockup_account.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
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

// LockupAccount defines the lockup account.
type LockupAccount struct {
	Address           string                `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Owner             string                `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Id                uint64                `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	StartTime         int64                 `protobuf:"varint,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime           int64                 `protobuf:"varint,5,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	OriginalLocking   cosmossdk_io_math.Int `protobuf:"bytes,6,opt,name=original_locking,json=originalLocking,proto3,customtype=cosmossdk.io/math.Int" json:"original_locking"`
	DelegatedFree     cosmossdk_io_math.Int `protobuf:"bytes,7,opt,name=delegated_free,json=delegatedFree,proto3,customtype=cosmossdk.io/math.Int" json:"delegated_free"`
	DelegatedLocking  cosmossdk_io_math.Int `protobuf:"bytes,8,opt,name=delegated_locking,json=delegatedLocking,proto3,customtype=cosmossdk.io/math.Int" json:"delegated_locking"`
	UnbondEntries     *UnbondingEntries     `protobuf:"bytes,9,opt,name=unbond_entries,json=unbondEntries,proto3" json:"unbond_entries,omitempty"`
	AdditionalLocking cosmossdk_io_math.Int `protobuf:"bytes,10,opt,name=additional_locking,json=additionalLocking,proto3,customtype=cosmossdk.io/math.Int" json:"additional_locking"`
}

func (m *LockupAccount) Reset()         { *m = LockupAccount{} }
func (m *LockupAccount) String() string { return proto.CompactTextString(m) }
func (*LockupAccount) ProtoMessage()    {}
func (*LockupAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_453c51fa8fc49e78, []int{0}
}
func (m *LockupAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LockupAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LockupAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LockupAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LockupAccount.Merge(m, src)
}
func (m *LockupAccount) XXX_Size() int {
	return m.Size()
}
func (m *LockupAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_LockupAccount.DiscardUnknown(m)
}

var xxx_messageInfo_LockupAccount proto.InternalMessageInfo

func (m *LockupAccount) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *LockupAccount) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *LockupAccount) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *LockupAccount) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *LockupAccount) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *LockupAccount) GetUnbondEntries() *UnbondingEntries {
	if m != nil {
		return m.UnbondEntries
	}
	return nil
}

// UnbondingEntries list of elements
type UnbondingEntries struct {
	Entries []*UnbondingEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (m *UnbondingEntries) Reset()         { *m = UnbondingEntries{} }
func (m *UnbondingEntries) String() string { return proto.CompactTextString(m) }
func (*UnbondingEntries) ProtoMessage()    {}
func (*UnbondingEntries) Descriptor() ([]byte, []int) {
	return fileDescriptor_453c51fa8fc49e78, []int{1}
}
func (m *UnbondingEntries) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnbondingEntries) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnbondingEntries.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnbondingEntries) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnbondingEntries.Merge(m, src)
}
func (m *UnbondingEntries) XXX_Size() int {
	return m.Size()
}
func (m *UnbondingEntries) XXX_DiscardUnknown() {
	xxx_messageInfo_UnbondingEntries.DiscardUnknown(m)
}

var xxx_messageInfo_UnbondingEntries proto.InternalMessageInfo

func (m *UnbondingEntries) GetEntries() []*UnbondingEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

// UnbondingEntry defines an entry tracking the lockup account unbonding operation.
type UnbondingEntry struct {
	CreationHeight int64 `protobuf:"varint,1,opt,name=creation_height,json=creationHeight,proto3" json:"creation_height,omitempty"`
	// end time of entry
	EndTime int64 `protobuf:"varint,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	// unbond amount
	Amount cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
	// validator address
	ValidatorAddress string `protobuf:"bytes,4,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}

func (m *UnbondingEntry) Reset()         { *m = UnbondingEntry{} }
func (m *UnbondingEntry) String() string { return proto.CompactTextString(m) }
func (*UnbondingEntry) ProtoMessage()    {}
func (*UnbondingEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_453c51fa8fc49e78, []int{2}
}
func (m *UnbondingEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnbondingEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnbondingEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnbondingEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnbondingEntry.Merge(m, src)
}
func (m *UnbondingEntry) XXX_Size() int {
	return m.Size()
}
func (m *UnbondingEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_UnbondingEntry.DiscardUnknown(m)
}

var xxx_messageInfo_UnbondingEntry proto.InternalMessageInfo

func (m *UnbondingEntry) GetCreationHeight() int64 {
	if m != nil {
		return m.CreationHeight
	}
	return 0
}

func (m *UnbondingEntry) GetEndTime() int64 {
	if m != nil {
		return m.EndTime
	}
	return 0
}

func (m *UnbondingEntry) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*LockupAccount)(nil), "sunrise.lockup.v1.LockupAccount")
	proto.RegisterType((*UnbondingEntries)(nil), "sunrise.lockup.v1.UnbondingEntries")
	proto.RegisterType((*UnbondingEntry)(nil), "sunrise.lockup.v1.UnbondingEntry")
}

func init() {
	proto.RegisterFile("sunrise/lockup/v1/lockup_account.proto", fileDescriptor_453c51fa8fc49e78)
}

var fileDescriptor_453c51fa8fc49e78 = []byte{
	// 568 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0x6d, 0xda, 0xae, 0x5d, 0x3d, 0xb5, 0x6b, 0xad, 0x4d, 0xca, 0x26, 0x2d, 0xeb, 0x8a, 0x04,
	0x15, 0xa8, 0x89, 0x56, 0x8e, 0x9c, 0x56, 0x04, 0xda, 0xd0, 0x04, 0x52, 0x80, 0x09, 0xed, 0x12,
	0xb9, 0xb1, 0x49, 0xad, 0x35, 0x76, 0xe5, 0xb8, 0x85, 0xfe, 0x0b, 0x7e, 0x06, 0x47, 0x0e, 0xfd,
	0x11, 0x3b, 0x4e, 0x3d, 0x21, 0x0e, 0x13, 0x6a, 0x0f, 0x88, 0x2b, 0xbf, 0x00, 0xd5, 0x8e, 0x37,
	0x3a, 0x0e, 0x53, 0x2f, 0x91, 0xbf, 0xf7, 0xde, 0xf7, 0xfc, 0x62, 0x7d, 0x36, 0x78, 0x98, 0x0c,
	0x99, 0xa0, 0x09, 0xf1, 0xfa, 0x3c, 0xbc, 0x18, 0x0e, 0xbc, 0xd1, 0x61, 0xba, 0x0a, 0x50, 0x18,
	0xf2, 0x21, 0x93, 0xee, 0x40, 0x70, 0xc9, 0x61, 0x2d, 0xd5, 0xb9, 0x9a, 0x75, 0x47, 0x87, 0xbb,
	0x35, 0x14, 0x53, 0xc6, 0x3d, 0xf5, 0xd5, 0xaa, 0xdd, 0x9d, 0x90, 0x27, 0x31, 0x4f, 0x02, 0x55,
	0x79, 0xba, 0x48, 0xa9, 0xad, 0x88, 0x47, 0x5c, 0xe3, 0x8b, 0x95, 0x46, 0x1b, 0x7f, 0xf2, 0xa0,
	0x7c, 0xaa, 0x1c, 0x8f, 0xf4, 0x76, 0xb0, 0x0d, 0x8a, 0x08, 0x63, 0x41, 0x92, 0xc4, 0xb6, 0xea,
	0x56, 0xb3, 0xd4, 0xb1, 0xa7, 0x93, 0xd6, 0x56, 0x6a, 0x75, 0xa4, 0x99, 0xb7, 0x52, 0x50, 0x16,
	0xf9, 0x46, 0x08, 0x5d, 0xb0, 0xc6, 0x3f, 0x31, 0x22, 0xec, 0xec, 0x3d, 0x1d, 0x5a, 0x06, 0xb7,
	0x41, 0x96, 0x62, 0x3b, 0x57, 0xb7, 0x9a, 0xf9, 0xce, 0xda, 0xd7, 0x5f, 0xdf, 0x1e, 0x5b, 0x7e,
	0x96, 0x62, 0xb8, 0x07, 0x40, 0x22, 0x91, 0x90, 0x81, 0xa4, 0x31, 0xb1, 0xf3, 0x75, 0xab, 0x99,
	0xf3, 0x4b, 0x0a, 0x79, 0x47, 0x63, 0x02, 0x77, 0xc0, 0x3a, 0x61, 0x58, 0x93, 0x6b, 0x8a, 0x2c,
	0x12, 0x86, 0x15, 0x75, 0x06, 0xaa, 0x5c, 0xd0, 0x88, 0x32, 0xd4, 0x0f, 0x16, 0x07, 0x44, 0x59,
	0x64, 0x17, 0x54, 0x96, 0x27, 0x97, 0xd7, 0xfb, 0x99, 0x1f, 0xd7, 0xfb, 0xdb, 0x3a, 0x4f, 0x82,
	0x2f, 0x5c, 0xca, 0xbd, 0x18, 0xc9, 0x9e, 0x7b, 0xc2, 0xe4, 0x74, 0xd2, 0x02, 0x69, 0xd0, 0x13,
	0x26, 0xfd, 0x4d, 0x63, 0x72, 0xaa, 0x3d, 0xa0, 0x0f, 0x2a, 0x98, 0xf4, 0x49, 0x84, 0x24, 0xc1,
	0xc1, 0x47, 0x41, 0x88, 0x5d, 0x5c, 0xdd, 0xb5, 0x7c, 0x63, 0xf1, 0x52, 0x10, 0x02, 0x3f, 0x80,
	0xda, 0xad, 0xa7, 0x09, 0xbb, 0xbe, 0xba, 0x6d, 0xf5, 0xc6, 0xc5, 0xa4, 0x7d, 0x05, 0x2a, 0x43,
	0xd6, 0xe5, 0x0c, 0x07, 0x84, 0x49, 0x41, 0x49, 0x62, 0x97, 0xea, 0x56, 0x73, 0xa3, 0xfd, 0xc0,
	0xfd, 0x6f, 0x78, 0xdc, 0xf7, 0x4a, 0x48, 0x59, 0xf4, 0x42, 0x4b, 0xfd, 0xb2, 0x6e, 0x4d, 0x4b,
	0x78, 0x0e, 0x20, 0xc2, 0x98, 0x4a, 0xca, 0xff, 0x3d, 0x53, 0xb0, 0x7a, 0xcc, 0xda, 0xad, 0x4d,
	0x9a, 0xb3, 0xf1, 0x06, 0x54, 0xef, 0x6e, 0x0f, 0x9f, 0x81, 0xa2, 0x09, 0x6d, 0xd5, 0x73, 0xcd,
	0x8d, 0xf6, 0xc1, 0x7d, 0xa1, 0xc7, 0xbe, 0xe9, 0x68, 0xfc, 0xb6, 0x40, 0x65, 0x99, 0x83, 0x8f,
	0xc0, 0x66, 0x28, 0x08, 0x5a, 0x6c, 0x1c, 0xf4, 0x08, 0x8d, 0x7a, 0x52, 0x8d, 0x73, 0xce, 0xaf,
	0x18, 0xf8, 0x58, 0xa1, 0x4b, 0x53, 0x95, 0x5d, 0x9e, 0xaa, 0xe7, 0xa0, 0x80, 0xe2, 0xc5, 0xa5,
	0x50, 0xa3, 0xba, 0xe2, 0x7f, 0xa7, 0xad, 0xf0, 0x35, 0xa8, 0x8d, 0x50, 0x9f, 0x62, 0x24, 0xb9,
	0x08, 0xcc, 0xcd, 0xca, 0x2b, 0xbf, 0x83, 0xe9, 0xa4, 0xb5, 0x97, 0xb6, 0x9c, 0x19, 0xcd, 0xf2,
	0x85, 0xa9, 0x8e, 0xee, 0xe0, 0x9d, 0xe3, 0xcb, 0x99, 0x63, 0x5d, 0xcd, 0x1c, 0xeb, 0xe7, 0xcc,
	0xb1, 0xbe, 0xcc, 0x9d, 0xcc, 0xd5, 0xdc, 0xc9, 0x7c, 0x9f, 0x3b, 0x99, 0x73, 0x37, 0xa2, 0xb2,
	0x37, 0xec, 0xba, 0x21, 0x8f, 0xbd, 0xf4, 0xec, 0xfa, 0x68, 0x4c, 0x84, 0x29, 0xbc, 0xcf, 0xe6,
	0x91, 0x91, 0xe3, 0x01, 0x49, 0xba, 0x05, 0xf5, 0x04, 0x3c, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff,
	0x7c, 0x69, 0xae, 0xae, 0x83, 0x04, 0x00, 0x00,
}

func (m *LockupAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LockupAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LockupAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AdditionalLocking.Size()
		i -= size
		if _, err := m.AdditionalLocking.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLockupAccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if m.UnbondEntries != nil {
		{
			size, err := m.UnbondEntries.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintLockupAccount(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	{
		size := m.DelegatedLocking.Size()
		i -= size
		if _, err := m.DelegatedLocking.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLockupAccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.DelegatedFree.Size()
		i -= size
		if _, err := m.DelegatedFree.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLockupAccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.OriginalLocking.Size()
		i -= size
		if _, err := m.OriginalLocking.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLockupAccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.EndTime != 0 {
		i = encodeVarintLockupAccount(dAtA, i, uint64(m.EndTime))
		i--
		dAtA[i] = 0x28
	}
	if m.StartTime != 0 {
		i = encodeVarintLockupAccount(dAtA, i, uint64(m.StartTime))
		i--
		dAtA[i] = 0x20
	}
	if m.Id != 0 {
		i = encodeVarintLockupAccount(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintLockupAccount(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintLockupAccount(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UnbondingEntries) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnbondingEntries) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnbondingEntries) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Entries) > 0 {
		for iNdEx := len(m.Entries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Entries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintLockupAccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *UnbondingEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnbondingEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnbondingEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintLockupAccount(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintLockupAccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.EndTime != 0 {
		i = encodeVarintLockupAccount(dAtA, i, uint64(m.EndTime))
		i--
		dAtA[i] = 0x10
	}
	if m.CreationHeight != 0 {
		i = encodeVarintLockupAccount(dAtA, i, uint64(m.CreationHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintLockupAccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovLockupAccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LockupAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovLockupAccount(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovLockupAccount(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovLockupAccount(uint64(m.Id))
	}
	if m.StartTime != 0 {
		n += 1 + sovLockupAccount(uint64(m.StartTime))
	}
	if m.EndTime != 0 {
		n += 1 + sovLockupAccount(uint64(m.EndTime))
	}
	l = m.OriginalLocking.Size()
	n += 1 + l + sovLockupAccount(uint64(l))
	l = m.DelegatedFree.Size()
	n += 1 + l + sovLockupAccount(uint64(l))
	l = m.DelegatedLocking.Size()
	n += 1 + l + sovLockupAccount(uint64(l))
	if m.UnbondEntries != nil {
		l = m.UnbondEntries.Size()
		n += 1 + l + sovLockupAccount(uint64(l))
	}
	l = m.AdditionalLocking.Size()
	n += 1 + l + sovLockupAccount(uint64(l))
	return n
}

func (m *UnbondingEntries) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Entries) > 0 {
		for _, e := range m.Entries {
			l = e.Size()
			n += 1 + l + sovLockupAccount(uint64(l))
		}
	}
	return n
}

func (m *UnbondingEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CreationHeight != 0 {
		n += 1 + sovLockupAccount(uint64(m.CreationHeight))
	}
	if m.EndTime != 0 {
		n += 1 + sovLockupAccount(uint64(m.EndTime))
	}
	l = m.Amount.Size()
	n += 1 + l + sovLockupAccount(uint64(l))
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovLockupAccount(uint64(l))
	}
	return n
}

func sovLockupAccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLockupAccount(x uint64) (n int) {
	return sovLockupAccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LockupAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLockupAccount
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
			return fmt.Errorf("proto: LockupAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LockupAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTime", wireType)
			}
			m.StartTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTime", wireType)
			}
			m.EndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginalLocking", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OriginalLocking.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatedFree", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DelegatedFree.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatedLocking", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DelegatedLocking.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondEntries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.UnbondEntries == nil {
				m.UnbondEntries = &UnbondingEntries{}
			}
			if err := m.UnbondEntries.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AdditionalLocking", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AdditionalLocking.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLockupAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLockupAccount
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
func (m *UnbondingEntries) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLockupAccount
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
			return fmt.Errorf("proto: UnbondingEntries: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnbondingEntries: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Entries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Entries = append(m.Entries, &UnbondingEntry{})
			if err := m.Entries[len(m.Entries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLockupAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLockupAccount
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
func (m *UnbondingEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLockupAccount
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
			return fmt.Errorf("proto: UnbondingEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnbondingEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreationHeight", wireType)
			}
			m.CreationHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreationHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndTime", wireType)
			}
			m.EndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLockupAccount
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
				return ErrInvalidLengthLockupAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthLockupAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLockupAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLockupAccount
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
func skipLockupAccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLockupAccount
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
					return 0, ErrIntOverflowLockupAccount
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
					return 0, ErrIntOverflowLockupAccount
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
				return 0, ErrInvalidLengthLockupAccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLockupAccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLockupAccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLockupAccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLockupAccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLockupAccount = fmt.Errorf("proto: unexpected end of group")
)
