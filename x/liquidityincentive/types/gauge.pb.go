// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquidityincentive/v1/gauge.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// Gauge
type Gauge struct {
	PreviousEpochId uint64                `protobuf:"varint,1,opt,name=previous_epoch_id,json=previousEpochId,proto3" json:"previous_epoch_id,omitempty"`
	PoolId          uint64                `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Count           cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=count,proto3,customtype=cosmossdk.io/math.Int" json:"count"`
}

func (m *Gauge) Reset()         { *m = Gauge{} }
func (m *Gauge) String() string { return proto.CompactTextString(m) }
func (*Gauge) ProtoMessage()    {}
func (*Gauge) Descriptor() ([]byte, []int) {
	return fileDescriptor_af61db3ea3dd0a1a, []int{0}
}
func (m *Gauge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Gauge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Gauge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Gauge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gauge.Merge(m, src)
}
func (m *Gauge) XXX_Size() int {
	return m.Size()
}
func (m *Gauge) XXX_DiscardUnknown() {
	xxx_messageInfo_Gauge.DiscardUnknown(m)
}

var xxx_messageInfo_Gauge proto.InternalMessageInfo

func (m *Gauge) GetPreviousEpochId() uint64 {
	if m != nil {
		return m.PreviousEpochId
	}
	return 0
}

func (m *Gauge) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

// TallyResult
type TallyResult struct {
	PoolId uint64                `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Count  cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=count,proto3,customtype=cosmossdk.io/math.Int" json:"count"`
}

func (m *TallyResult) Reset()         { *m = TallyResult{} }
func (m *TallyResult) String() string { return proto.CompactTextString(m) }
func (*TallyResult) ProtoMessage()    {}
func (*TallyResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_af61db3ea3dd0a1a, []int{1}
}
func (m *TallyResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TallyResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TallyResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TallyResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TallyResult.Merge(m, src)
}
func (m *TallyResult) XXX_Size() int {
	return m.Size()
}
func (m *TallyResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TallyResult.DiscardUnknown(m)
}

var xxx_messageInfo_TallyResult proto.InternalMessageInfo

func (m *TallyResult) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

// PoolWeight
type PoolWeight struct {
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Weight string `protobuf:"bytes,2,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (m *PoolWeight) Reset()         { *m = PoolWeight{} }
func (m *PoolWeight) String() string { return proto.CompactTextString(m) }
func (*PoolWeight) ProtoMessage()    {}
func (*PoolWeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_af61db3ea3dd0a1a, []int{2}
}
func (m *PoolWeight) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolWeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolWeight.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolWeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolWeight.Merge(m, src)
}
func (m *PoolWeight) XXX_Size() int {
	return m.Size()
}
func (m *PoolWeight) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolWeight.DiscardUnknown(m)
}

var xxx_messageInfo_PoolWeight proto.InternalMessageInfo

func (m *PoolWeight) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *PoolWeight) GetWeight() string {
	if m != nil {
		return m.Weight
	}
	return ""
}

// Vote
type Vote struct {
	Sender      string       `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	PoolWeights []PoolWeight `protobuf:"bytes,2,rep,name=pool_weights,json=poolWeights,proto3" json:"pool_weights"`
}

func (m *Vote) Reset()         { *m = Vote{} }
func (m *Vote) String() string { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()    {}
func (*Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_af61db3ea3dd0a1a, []int{3}
}
func (m *Vote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vote.Merge(m, src)
}
func (m *Vote) XXX_Size() int {
	return m.Size()
}
func (m *Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Vote proto.InternalMessageInfo

func (m *Vote) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *Vote) GetPoolWeights() []PoolWeight {
	if m != nil {
		return m.PoolWeights
	}
	return nil
}

func init() {
	proto.RegisterType((*Gauge)(nil), "sunrise.liquidityincentive.v1.Gauge")
	proto.RegisterType((*TallyResult)(nil), "sunrise.liquidityincentive.v1.TallyResult")
	proto.RegisterType((*PoolWeight)(nil), "sunrise.liquidityincentive.v1.PoolWeight")
	proto.RegisterType((*Vote)(nil), "sunrise.liquidityincentive.v1.Vote")
}

func init() {
	proto.RegisterFile("sunrise/liquidityincentive/v1/gauge.proto", fileDescriptor_af61db3ea3dd0a1a)
}

var fileDescriptor_af61db3ea3dd0a1a = []byte{
	// 412 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0xaa, 0xd3, 0x40,
	0x14, 0xc6, 0x33, 0xbd, 0xbd, 0x91, 0x3b, 0x15, 0xc5, 0x70, 0xc5, 0x78, 0xc1, 0xdc, 0xd2, 0x85,
	0xb4, 0x4a, 0x27, 0x56, 0x77, 0xee, 0x1a, 0x14, 0xc9, 0x42, 0x90, 0x58, 0x14, 0xdc, 0x94, 0x34,
	0x33, 0x24, 0x83, 0xe9, 0x9c, 0x98, 0x99, 0x44, 0xf3, 0x0e, 0x82, 0x3e, 0x4c, 0x1f, 0xa2, 0xcb,
	0xd2, 0x95, 0xb8, 0x28, 0xd2, 0xbe, 0x88, 0xe4, 0x4f, 0xb5, 0x60, 0xe9, 0xc2, 0x5d, 0xce, 0x77,
	0xbe, 0xf3, 0xfd, 0x4e, 0x98, 0x83, 0x07, 0x32, 0x13, 0x29, 0x97, 0xcc, 0x8e, 0xf9, 0xa7, 0x8c,
	0x53, 0xae, 0x0a, 0x2e, 0x02, 0x26, 0x14, 0xcf, 0x99, 0x9d, 0x8f, 0xec, 0xd0, 0xcf, 0x42, 0x46,
	0x92, 0x14, 0x14, 0x18, 0x0f, 0x1a, 0x2b, 0xf9, 0xd7, 0x4a, 0xf2, 0xd1, 0xd5, 0xfd, 0x00, 0xe4,
	0x1c, 0xe4, 0xb4, 0x32, 0xdb, 0x75, 0x51, 0x4f, 0x5e, 0x5d, 0x86, 0x10, 0x42, 0xad, 0x97, 0x5f,
	0xb5, 0xda, 0xfb, 0x86, 0xf0, 0xf9, 0xab, 0x32, 0xdf, 0x78, 0x84, 0xef, 0x24, 0x29, 0xcb, 0x39,
	0x64, 0x72, 0xca, 0x12, 0x08, 0xa2, 0x29, 0xa7, 0x26, 0xea, 0xa2, 0x7e, 0xdb, 0xbb, 0xbd, 0x6f,
	0xbc, 0x2c, 0x75, 0x97, 0x1a, 0xf7, 0xf0, 0x8d, 0x04, 0x20, 0x2e, 0x1d, 0xad, 0xca, 0xa1, 0x97,
	0xa5, 0x4b, 0x8d, 0x31, 0x3e, 0x0f, 0x20, 0x13, 0xca, 0x3c, 0xeb, 0xa2, 0xfe, 0x85, 0xf3, 0x78,
	0xb9, 0xb9, 0xd6, 0x7e, 0x6e, 0xae, 0xef, 0xd6, 0x9b, 0x48, 0xfa, 0x91, 0x70, 0xb0, 0xe7, 0xbe,
	0x8a, 0x88, 0x2b, 0xd4, 0x7a, 0x31, 0xc4, 0xcd, 0x8a, 0xae, 0x50, 0x5e, 0x3d, 0xd9, 0xe3, 0xb8,
	0x33, 0xf1, 0xe3, 0xb8, 0xf0, 0x98, 0xcc, 0x62, 0x75, 0x88, 0x42, 0xc7, 0x51, 0xad, 0xff, 0x46,
	0xbd, 0xc6, 0xf8, 0x0d, 0x40, 0xfc, 0x9e, 0xf1, 0x30, 0x3a, 0x41, 0x7a, 0x88, 0xf5, 0xcf, 0x95,
	0xa5, 0x41, 0xdd, 0x3a, 0x48, 0x7b, 0xc1, 0x02, 0xaf, 0xe9, 0xf6, 0xbe, 0x22, 0xdc, 0x7e, 0x07,
	0x8a, 0x19, 0x4f, 0xb0, 0x2e, 0x99, 0xa0, 0x2c, 0xad, 0x82, 0x2e, 0x1c, 0x73, 0xbd, 0x18, 0x5e,
	0x36, 0x03, 0x63, 0x4a, 0x53, 0x26, 0xe5, 0x5b, 0x95, 0x72, 0x11, 0x7a, 0x8d, 0xcf, 0xf0, 0xf0,
	0xcd, 0x8a, 0x5d, 0x27, 0x49, 0xb3, 0xd5, 0x3d, 0xeb, 0x77, 0x9e, 0x0e, 0xc8, 0xc9, 0xd7, 0x26,
	0x7f, 0x97, 0x77, 0xda, 0xe5, 0xef, 0x7b, 0x9d, 0xe4, 0x8f, 0x22, 0x9d, 0xc9, 0x72, 0x6b, 0xa1,
	0xd5, 0xd6, 0x42, 0xbf, 0xb6, 0x16, 0xfa, 0xbe, 0xb3, 0xb4, 0xd5, 0xce, 0xd2, 0x7e, 0xec, 0x2c,
	0xed, 0xc3, 0xf3, 0x90, 0xab, 0x28, 0x9b, 0x91, 0x00, 0xe6, 0x76, 0x43, 0x88, 0xfd, 0x82, 0xa5,
	0xfb, 0xc2, 0xfe, 0x72, 0xec, 0x12, 0x55, 0x91, 0x30, 0x39, 0xd3, 0xab, 0xbb, 0x79, 0xf6, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0xad, 0x1d, 0xdc, 0x2c, 0xb4, 0x02, 0x00, 0x00,
}

func (m *Gauge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Gauge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Gauge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Count.Size()
		i -= size
		if _, err := m.Count.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGauge(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.PoolId != 0 {
		i = encodeVarintGauge(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x10
	}
	if m.PreviousEpochId != 0 {
		i = encodeVarintGauge(dAtA, i, uint64(m.PreviousEpochId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TallyResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TallyResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TallyResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Count.Size()
		i -= size
		if _, err := m.Count.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGauge(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PoolId != 0 {
		i = encodeVarintGauge(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PoolWeight) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolWeight) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolWeight) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Weight) > 0 {
		i -= len(m.Weight)
		copy(dAtA[i:], m.Weight)
		i = encodeVarintGauge(dAtA, i, uint64(len(m.Weight)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintGauge(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Vote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PoolWeights) > 0 {
		for iNdEx := len(m.PoolWeights) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolWeights[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGauge(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintGauge(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGauge(dAtA []byte, offset int, v uint64) int {
	offset -= sovGauge(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Gauge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PreviousEpochId != 0 {
		n += 1 + sovGauge(uint64(m.PreviousEpochId))
	}
	if m.PoolId != 0 {
		n += 1 + sovGauge(uint64(m.PoolId))
	}
	l = m.Count.Size()
	n += 1 + l + sovGauge(uint64(l))
	return n
}

func (m *TallyResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovGauge(uint64(m.PoolId))
	}
	l = m.Count.Size()
	n += 1 + l + sovGauge(uint64(l))
	return n
}

func (m *PoolWeight) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovGauge(uint64(m.PoolId))
	}
	l = len(m.Weight)
	if l > 0 {
		n += 1 + l + sovGauge(uint64(l))
	}
	return n
}

func (m *Vote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovGauge(uint64(l))
	}
	if len(m.PoolWeights) > 0 {
		for _, e := range m.PoolWeights {
			l = e.Size()
			n += 1 + l + sovGauge(uint64(l))
		}
	}
	return n
}

func sovGauge(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGauge(x uint64) (n int) {
	return sovGauge(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Gauge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGauge
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
			return fmt.Errorf("proto: Gauge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Gauge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousEpochId", wireType)
			}
			m.PreviousEpochId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PreviousEpochId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
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
				return ErrInvalidLengthGauge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGauge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Count.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGauge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGauge
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
func (m *TallyResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGauge
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
			return fmt.Errorf("proto: TallyResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TallyResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
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
				return ErrInvalidLengthGauge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGauge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Count.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGauge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGauge
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
func (m *PoolWeight) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGauge
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
			return fmt.Errorf("proto: PoolWeight: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolWeight: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
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
				return ErrInvalidLengthGauge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGauge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Weight = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGauge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGauge
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
func (m *Vote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGauge
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
			return fmt.Errorf("proto: Vote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
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
				return ErrInvalidLengthGauge
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGauge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolWeights", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGauge
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
				return ErrInvalidLengthGauge
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGauge
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolWeights = append(m.PoolWeights, PoolWeight{})
			if err := m.PoolWeights[len(m.PoolWeights)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGauge(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGauge
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
func skipGauge(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGauge
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
					return 0, ErrIntOverflowGauge
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
					return 0, ErrIntOverflowGauge
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
				return 0, ErrInvalidLengthGauge
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGauge
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGauge
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGauge        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGauge          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGauge = fmt.Errorf("proto: unexpected end of group")
)
