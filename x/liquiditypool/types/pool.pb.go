// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquiditypool/v1/pool.proto

package types

import (
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

// Pool
type Pool struct {
	Id                   uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DenomBase            string     `protobuf:"bytes,2,opt,name=denom_base,json=denomBase,proto3" json:"denom_base,omitempty"`
	DenomQuote           string     `protobuf:"bytes,3,opt,name=denom_quote,json=denomQuote,proto3" json:"denom_quote,omitempty"`
	FeeRate              string     `protobuf:"bytes,4,opt,name=fee_rate,json=feeRate,proto3" json:"fee_rate,omitempty"`
	TickParams           TickParams `protobuf:"bytes,5,opt,name=tick_params,json=tickParams,proto3" json:"tick_params"`
	CurrentTick          int64      `protobuf:"varint,6,opt,name=current_tick,json=currentTick,proto3" json:"current_tick,omitempty"`
	CurrentTickLiquidity string     `protobuf:"bytes,7,opt,name=current_tick_liquidity,json=currentTickLiquidity,proto3" json:"current_tick_liquidity,omitempty"`
	CurrentSqrtPrice     string     `protobuf:"bytes,8,opt,name=current_sqrt_price,json=currentSqrtPrice,proto3" json:"current_sqrt_price,omitempty"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f2038cccc8b5a8, []int{0}
}
func (m *Pool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pool.Merge(m, src)
}
func (m *Pool) XXX_Size() int {
	return m.Size()
}
func (m *Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_Pool proto.InternalMessageInfo

func (m *Pool) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Pool) GetDenomBase() string {
	if m != nil {
		return m.DenomBase
	}
	return ""
}

func (m *Pool) GetDenomQuote() string {
	if m != nil {
		return m.DenomQuote
	}
	return ""
}

func (m *Pool) GetFeeRate() string {
	if m != nil {
		return m.FeeRate
	}
	return ""
}

func (m *Pool) GetTickParams() TickParams {
	if m != nil {
		return m.TickParams
	}
	return TickParams{}
}

func (m *Pool) GetCurrentTick() int64 {
	if m != nil {
		return m.CurrentTick
	}
	return 0
}

func (m *Pool) GetCurrentTickLiquidity() string {
	if m != nil {
		return m.CurrentTickLiquidity
	}
	return ""
}

func (m *Pool) GetCurrentSqrtPrice() string {
	if m != nil {
		return m.CurrentSqrtPrice
	}
	return ""
}

// TickParams
// PriceRatio^(Tick - BaseOffSet)
type TickParams struct {
	// Basically 1.0001
	PriceRatio string `protobuf:"bytes,1,opt,name=price_ratio,json=priceRatio,proto3" json:"price_ratio,omitempty"`
	// basically 0 and (-1,0]. In the 1:1 stable pair, -0.5 would work
	BaseOffset string `protobuf:"bytes,2,opt,name=base_offset,json=baseOffset,proto3" json:"base_offset,omitempty"`
}

func (m *TickParams) Reset()         { *m = TickParams{} }
func (m *TickParams) String() string { return proto.CompactTextString(m) }
func (*TickParams) ProtoMessage()    {}
func (*TickParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f2038cccc8b5a8, []int{1}
}
func (m *TickParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TickParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TickParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TickParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TickParams.Merge(m, src)
}
func (m *TickParams) XXX_Size() int {
	return m.Size()
}
func (m *TickParams) XXX_DiscardUnknown() {
	xxx_messageInfo_TickParams.DiscardUnknown(m)
}

var xxx_messageInfo_TickParams proto.InternalMessageInfo

func (m *TickParams) GetPriceRatio() string {
	if m != nil {
		return m.PriceRatio
	}
	return ""
}

func (m *TickParams) GetBaseOffset() string {
	if m != nil {
		return m.BaseOffset
	}
	return ""
}

func init() {
	proto.RegisterType((*Pool)(nil), "sunrise.liquiditypool.v1.Pool")
	proto.RegisterType((*TickParams)(nil), "sunrise.liquiditypool.v1.TickParams")
}

func init() {
	proto.RegisterFile("sunrise/liquiditypool/v1/pool.proto", fileDescriptor_69f2038cccc8b5a8)
}

var fileDescriptor_69f2038cccc8b5a8 = []byte{
	// 422 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcd, 0x6e, 0xd4, 0x30,
	0x10, 0xc7, 0xd7, 0xdb, 0xa5, 0x1f, 0x13, 0x54, 0x21, 0xab, 0x42, 0xa6, 0x12, 0xe9, 0x52, 0x38,
	0x2c, 0x07, 0x12, 0x15, 0xc4, 0x8d, 0x53, 0xd4, 0x1b, 0x48, 0x5d, 0x02, 0x27, 0x2e, 0x56, 0x36,
	0x99, 0x2c, 0x56, 0x93, 0x38, 0xb1, 0x9d, 0x8a, 0x7d, 0x0b, 0x1e, 0x86, 0x87, 0xe8, 0xb1, 0xe2,
	0xc4, 0x09, 0xa1, 0xdd, 0x67, 0xe0, 0x8e, 0xec, 0x64, 0x4b, 0xa1, 0xdb, 0x53, 0xe2, 0xff, 0xfc,
	0xfe, 0x63, 0xcf, 0x07, 0x3c, 0xd5, 0x6d, 0xa5, 0x84, 0xc6, 0xb0, 0x10, 0x4d, 0x2b, 0x32, 0x61,
	0x16, 0xb5, 0x94, 0x45, 0x78, 0x71, 0x12, 0xda, 0x6f, 0x50, 0x2b, 0x69, 0x24, 0x65, 0x3d, 0x14,
	0xfc, 0x03, 0x05, 0x17, 0x27, 0x87, 0x8f, 0x52, 0xa9, 0x4b, 0xa9, 0xb9, 0xe3, 0xc2, 0xee, 0xd0,
	0x99, 0x0e, 0x0f, 0xe6, 0x72, 0x2e, 0x3b, 0xdd, 0xfe, 0x75, 0xea, 0xf1, 0xef, 0x21, 0x8c, 0xa6,
	0x52, 0x16, 0x74, 0x1f, 0x86, 0x22, 0x63, 0x64, 0x4c, 0x26, 0xa3, 0x78, 0x28, 0x32, 0xfa, 0x18,
	0x20, 0xc3, 0x4a, 0x96, 0x7c, 0x96, 0x68, 0x64, 0xc3, 0x31, 0x99, 0xec, 0xc5, 0x7b, 0x4e, 0x89,
	0x12, 0x8d, 0xf4, 0x08, 0xbc, 0x2e, 0xdc, 0xb4, 0xd2, 0x20, 0xdb, 0x72, 0xf1, 0xce, 0xf1, 0xde,
	0x2a, 0xf4, 0x39, 0xec, 0xe6, 0x88, 0x5c, 0x25, 0x06, 0xd9, 0xc8, 0x46, 0xa3, 0xfd, 0xef, 0xdf,
	0x5e, 0x40, 0xff, 0xa4, 0x53, 0x4c, 0xe3, 0x9d, 0x1c, 0x31, 0x4e, 0x0c, 0xd2, 0xb7, 0xe0, 0x19,
	0x91, 0x9e, 0xf3, 0x3a, 0x51, 0x49, 0xa9, 0xd9, 0xbd, 0x31, 0x99, 0x78, 0x2f, 0x9f, 0x05, 0x77,
	0x15, 0x19, 0x7c, 0x14, 0xe9, 0xf9, 0xd4, 0xb1, 0xd1, 0xe8, 0xf2, 0xe7, 0xd1, 0x20, 0x06, 0x73,
	0xad, 0xd0, 0x27, 0x70, 0x3f, 0x6d, 0x95, 0xc2, 0xca, 0x70, 0xab, 0xb2, 0xed, 0x31, 0x99, 0x6c,
	0xc5, 0x5e, 0xaf, 0x59, 0x2b, 0x3d, 0x85, 0x87, 0x37, 0x11, 0x7e, 0x7d, 0x01, 0xdb, 0xd9, 0xf8,
	0xd0, 0x83, 0x1b, 0xe6, 0x77, 0x6b, 0x96, 0xbe, 0x01, 0xba, 0xce, 0xa2, 0x1b, 0x65, 0x78, 0xad,
	0x44, 0x8a, 0x6c, 0x77, 0x63, 0x86, 0x07, 0x3d, 0xf9, 0xa1, 0x51, 0x66, 0x6a, 0xb9, 0xe3, 0x0a,
	0xe0, 0x6f, 0x19, 0x34, 0x04, 0xcf, 0xd9, 0x6d, 0xbb, 0x84, 0x74, 0x53, 0xb8, 0x9d, 0x04, 0x1c,
	0x12, 0x5b, 0xc2, 0x1a, 0xec, 0x5c, 0xb8, 0xcc, 0x73, 0x8d, 0xa6, 0x1b, 0xcf, 0x6d, 0x83, 0x45,
	0xce, 0x1c, 0x11, 0x9d, 0x5d, 0x2e, 0x7d, 0x72, 0xb5, 0xf4, 0xc9, 0xaf, 0xa5, 0x4f, 0xbe, 0xae,
	0xfc, 0xc1, 0xd5, 0xca, 0x1f, 0xfc, 0x58, 0xf9, 0x83, 0x4f, 0xaf, 0xe7, 0xc2, 0x7c, 0x6e, 0x67,
	0x41, 0x2a, 0xcb, 0xb0, 0x6f, 0x79, 0x91, 0x2c, 0x50, 0xad, 0x0f, 0xe1, 0x97, 0xff, 0x76, 0xd1,
	0x2c, 0x6a, 0xd4, 0xb3, 0x6d, 0xb7, 0x3f, 0xaf, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0xf6, 0xd6,
	0x86, 0xfd, 0xb1, 0x02, 0x00, 0x00,
}

func (m *Pool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.CurrentSqrtPrice) > 0 {
		i -= len(m.CurrentSqrtPrice)
		copy(dAtA[i:], m.CurrentSqrtPrice)
		i = encodeVarintPool(dAtA, i, uint64(len(m.CurrentSqrtPrice)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.CurrentTickLiquidity) > 0 {
		i -= len(m.CurrentTickLiquidity)
		copy(dAtA[i:], m.CurrentTickLiquidity)
		i = encodeVarintPool(dAtA, i, uint64(len(m.CurrentTickLiquidity)))
		i--
		dAtA[i] = 0x3a
	}
	if m.CurrentTick != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.CurrentTick))
		i--
		dAtA[i] = 0x30
	}
	{
		size, err := m.TickParams.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.FeeRate) > 0 {
		i -= len(m.FeeRate)
		copy(dAtA[i:], m.FeeRate)
		i = encodeVarintPool(dAtA, i, uint64(len(m.FeeRate)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.DenomQuote) > 0 {
		i -= len(m.DenomQuote)
		copy(dAtA[i:], m.DenomQuote)
		i = encodeVarintPool(dAtA, i, uint64(len(m.DenomQuote)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DenomBase) > 0 {
		i -= len(m.DenomBase)
		copy(dAtA[i:], m.DenomBase)
		i = encodeVarintPool(dAtA, i, uint64(len(m.DenomBase)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TickParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TickParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TickParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BaseOffset) > 0 {
		i -= len(m.BaseOffset)
		copy(dAtA[i:], m.BaseOffset)
		i = encodeVarintPool(dAtA, i, uint64(len(m.BaseOffset)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PriceRatio) > 0 {
		i -= len(m.PriceRatio)
		copy(dAtA[i:], m.PriceRatio)
		i = encodeVarintPool(dAtA, i, uint64(len(m.PriceRatio)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPool(uint64(m.Id))
	}
	l = len(m.DenomBase)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = len(m.DenomQuote)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = len(m.FeeRate)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.TickParams.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.CurrentTick != 0 {
		n += 1 + sovPool(uint64(m.CurrentTick))
	}
	l = len(m.CurrentTickLiquidity)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = len(m.CurrentSqrtPrice)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	return n
}

func (m *TickParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PriceRatio)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = len(m.BaseOffset)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Pool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: Pool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomBase", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomBase = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DenomQuote", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DenomQuote = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeRate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TickParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TickParams.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentTick", wireType)
			}
			m.CurrentTick = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentTick |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentTickLiquidity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentTickLiquidity = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentSqrtPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentSqrtPrice = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func (m *TickParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: TickParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TickParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PriceRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PriceRatio = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseOffset", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseOffset = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func skipPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
				return 0, ErrInvalidLengthPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPool = fmt.Errorf("proto: unexpected end of group")
)
