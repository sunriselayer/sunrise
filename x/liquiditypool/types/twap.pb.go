// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquiditypool/v1/twap.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type Twap struct {
	BaseDenom       string                      `protobuf:"bytes,1,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	QuoteDenom      string                      `protobuf:"bytes,2,opt,name=quote_denom,json=quoteDenom,proto3" json:"quote_denom,omitempty"`
	Value           cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=value,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"value"`
	LatestTimestamp *time.Time                  `protobuf:"bytes,4,opt,name=latest_timestamp,json=latestTimestamp,proto3,stdtime" json:"latest_timestamp,omitempty"`
}

func (m *Twap) Reset()         { *m = Twap{} }
func (m *Twap) String() string { return proto.CompactTextString(m) }
func (*Twap) ProtoMessage()    {}
func (*Twap) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef7adb8b10610411, []int{0}
}
func (m *Twap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Twap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Twap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Twap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Twap.Merge(m, src)
}
func (m *Twap) XXX_Size() int {
	return m.Size()
}
func (m *Twap) XXX_DiscardUnknown() {
	xxx_messageInfo_Twap.DiscardUnknown(m)
}

var xxx_messageInfo_Twap proto.InternalMessageInfo

func (m *Twap) GetBaseDenom() string {
	if m != nil {
		return m.BaseDenom
	}
	return ""
}

func (m *Twap) GetQuoteDenom() string {
	if m != nil {
		return m.QuoteDenom
	}
	return ""
}

func (m *Twap) GetLatestTimestamp() *time.Time {
	if m != nil {
		return m.LatestTimestamp
	}
	return nil
}

type PriceFootprint struct {
	Price     cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=price,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"price"`
	Timestamp *time.Time                  `protobuf:"bytes,2,opt,name=timestamp,proto3,stdtime" json:"timestamp,omitempty"`
}

func (m *PriceFootprint) Reset()         { *m = PriceFootprint{} }
func (m *PriceFootprint) String() string { return proto.CompactTextString(m) }
func (*PriceFootprint) ProtoMessage()    {}
func (*PriceFootprint) Descriptor() ([]byte, []int) {
	return fileDescriptor_ef7adb8b10610411, []int{1}
}
func (m *PriceFootprint) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PriceFootprint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PriceFootprint.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PriceFootprint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PriceFootprint.Merge(m, src)
}
func (m *PriceFootprint) XXX_Size() int {
	return m.Size()
}
func (m *PriceFootprint) XXX_DiscardUnknown() {
	xxx_messageInfo_PriceFootprint.DiscardUnknown(m)
}

var xxx_messageInfo_PriceFootprint proto.InternalMessageInfo

func (m *PriceFootprint) GetTimestamp() *time.Time {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*Twap)(nil), "sunrise.liquiditypool.v1.Twap")
	proto.RegisterType((*PriceFootprint)(nil), "sunrise.liquiditypool.v1.PriceFootprint")
}

func init() {
	proto.RegisterFile("sunrise/liquiditypool/v1/twap.proto", fileDescriptor_ef7adb8b10610411)
}

var fileDescriptor_ef7adb8b10610411 = []byte{
	// 399 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xbf, 0xee, 0xd3, 0x30,
	0x10, 0x8e, 0x7f, 0x14, 0xa4, 0xba, 0x12, 0x7f, 0x22, 0x86, 0x10, 0x44, 0x52, 0x95, 0xa5, 0x42,
	0xaa, 0xad, 0x82, 0x84, 0xc4, 0xc2, 0x50, 0x55, 0x2c, 0x74, 0x40, 0xa5, 0x13, 0x4b, 0xe5, 0xa4,
	0x26, 0xb5, 0x48, 0x72, 0x6e, 0xec, 0xb4, 0xe4, 0x2d, 0xfa, 0x04, 0xcc, 0x8c, 0x0c, 0x3c, 0x44,
	0xc7, 0x8a, 0x09, 0x31, 0x14, 0x68, 0x07, 0x5e, 0x03, 0x25, 0x4e, 0x5a, 0xc1, 0x84, 0x58, 0x2c,
	0xdf, 0xf7, 0x7d, 0x77, 0xf7, 0xdd, 0xe9, 0xf0, 0x43, 0x95, 0xa7, 0x99, 0x50, 0x9c, 0xc6, 0x62,
	0x95, 0x8b, 0x85, 0xd0, 0x85, 0x04, 0x88, 0xe9, 0x7a, 0x48, 0xf5, 0x86, 0x49, 0x22, 0x33, 0xd0,
	0x60, 0x3b, 0xb5, 0x88, 0xfc, 0x21, 0x22, 0xeb, 0xa1, 0x7b, 0x87, 0x25, 0x22, 0x05, 0x5a, 0xbd,
	0x46, 0xec, 0xfa, 0x11, 0x40, 0x14, 0x73, 0x5a, 0x45, 0x41, 0xfe, 0x96, 0x6a, 0x91, 0x70, 0xa5,
	0x59, 0x52, 0x57, 0x73, 0xef, 0x46, 0x10, 0x41, 0xf5, 0xa5, 0xe5, 0xaf, 0x46, 0xef, 0x85, 0xa0,
	0x12, 0x50, 0x73, 0x43, 0x98, 0xc0, 0x50, 0xbd, 0x9f, 0x08, 0xb7, 0x66, 0x1b, 0x26, 0xed, 0x07,
	0x18, 0x07, 0x4c, 0xf1, 0xf9, 0x82, 0xa7, 0x90, 0x38, 0xa8, 0x8b, 0xfa, 0xed, 0x69, 0xbb, 0x44,
	0xc6, 0x25, 0x60, 0xfb, 0xb8, 0xb3, 0xca, 0x41, 0x37, 0xfc, 0x55, 0xc5, 0xe3, 0x0a, 0x32, 0x82,
	0x09, 0xbe, 0xbe, 0x66, 0x71, 0xce, 0x9d, 0x6b, 0x25, 0x35, 0x7a, 0xba, 0x3b, 0xf8, 0xd6, 0xb7,
	0x83, 0x7f, 0xdf, 0x74, 0x53, 0x8b, 0x77, 0x44, 0x00, 0x4d, 0x98, 0x5e, 0x92, 0x09, 0x8f, 0x58,
	0x58, 0x8c, 0x79, 0xf8, 0xe5, 0xf3, 0x00, 0xd7, 0x66, 0xc6, 0x3c, 0xfc, 0xf8, 0xeb, 0xd3, 0x23,
	0x34, 0x35, 0x45, 0xec, 0x97, 0xf8, 0x76, 0xcc, 0x34, 0x57, 0x7a, 0x7e, 0x9e, 0xd0, 0x69, 0x75,
	0x51, 0xbf, 0xf3, 0xd8, 0x25, 0x66, 0x07, 0xa4, 0xd9, 0x01, 0x99, 0x35, 0x8a, 0x51, 0x6b, 0xfb,
	0xdd, 0x47, 0xd3, 0x5b, 0x26, 0xf3, 0x0c, 0xf7, 0x3e, 0x20, 0x7c, 0xf3, 0x55, 0x26, 0x42, 0xfe,
	0x02, 0x40, 0xcb, 0x4c, 0xa4, 0xba, 0x74, 0x2b, 0x4b, 0xc4, 0x0c, 0xfa, 0xff, 0x6e, 0xab, 0x22,
	0xf6, 0x73, 0xdc, 0xbe, 0xd8, 0xbc, 0xfa, 0x47, 0x9b, 0x97, 0x94, 0xd1, 0xeb, 0xdd, 0xd1, 0x43,
	0xfb, 0xa3, 0x87, 0x7e, 0x1c, 0x3d, 0xb4, 0x3d, 0x79, 0xd6, 0xfe, 0xe4, 0x59, 0x5f, 0x4f, 0x9e,
	0xf5, 0xe6, 0x59, 0x24, 0xf4, 0x32, 0x0f, 0x48, 0x08, 0x09, 0xad, 0x0f, 0x25, 0x66, 0x05, 0xcf,
	0x9a, 0x60, 0xc0, 0xa4, 0xa4, 0xef, 0xff, 0x3a, 0x30, 0x5d, 0x48, 0xae, 0x82, 0x1b, 0x55, 0xe7,
	0x27, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xae, 0x13, 0x55, 0x1e, 0x86, 0x02, 0x00, 0x00,
}

func (m *Twap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Twap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Twap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LatestTimestamp != nil {
		n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(*m.LatestTimestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LatestTimestamp):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintTwap(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.Value.Size()
		i -= size
		if _, err := m.Value.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.QuoteDenom) > 0 {
		i -= len(m.QuoteDenom)
		copy(dAtA[i:], m.QuoteDenom)
		i = encodeVarintTwap(dAtA, i, uint64(len(m.QuoteDenom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintTwap(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PriceFootprint) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PriceFootprint) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PriceFootprint) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(*m.Timestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.Timestamp):])
		if err2 != nil {
			return 0, err2
		}
		i -= n2
		i = encodeVarintTwap(dAtA, i, uint64(n2))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintTwap(dAtA []byte, offset int, v uint64) int {
	offset -= sovTwap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Twap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovTwap(uint64(l))
	}
	l = len(m.QuoteDenom)
	if l > 0 {
		n += 1 + l + sovTwap(uint64(l))
	}
	l = m.Value.Size()
	n += 1 + l + sovTwap(uint64(l))
	if m.LatestTimestamp != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LatestTimestamp)
		n += 1 + l + sovTwap(uint64(l))
	}
	return n
}

func (m *PriceFootprint) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Price.Size()
	n += 1 + l + sovTwap(uint64(l))
	if m.Timestamp != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.Timestamp)
		n += 1 + l + sovTwap(uint64(l))
	}
	return n
}

func sovTwap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTwap(x uint64) (n int) {
	return sovTwap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Twap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTwap
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
			return fmt.Errorf("proto: Twap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Twap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.QuoteDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LatestTimestamp == nil {
				m.LatestTimestamp = new(time.Time)
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(m.LatestTimestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTwap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTwap
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
func (m *PriceFootprint) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTwap
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
			return fmt.Errorf("proto: PriceFootprint: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PriceFootprint: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTwap
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
				return ErrInvalidLengthTwap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = new(time.Time)
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTwap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTwap
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
func skipTwap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTwap
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
					return 0, ErrIntOverflowTwap
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
					return 0, ErrIntOverflowTwap
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
				return 0, ErrInvalidLengthTwap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTwap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTwap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTwap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTwap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTwap = fmt.Errorf("proto: unexpected end of group")
)
