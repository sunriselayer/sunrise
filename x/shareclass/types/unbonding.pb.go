// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/shareclass/v1/unbonding.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"

	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// Unbonding
type Unbonding struct {
	Address        string     `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	CompletionTime time.Time  `protobuf:"bytes,2,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
	Amount         types.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *Unbonding) Reset()         { *m = Unbonding{} }
func (m *Unbonding) String() string { return proto.CompactTextString(m) }
func (*Unbonding) ProtoMessage()    {}
func (*Unbonding) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca9cb0aa09d0e4a3, []int{0}
}
func (m *Unbonding) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Unbonding) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Unbonding.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Unbonding) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Unbonding.Merge(m, src)
}
func (m *Unbonding) XXX_Size() int {
	return m.Size()
}
func (m *Unbonding) XXX_DiscardUnknown() {
	xxx_messageInfo_Unbonding.DiscardUnknown(m)
}

var xxx_messageInfo_Unbonding proto.InternalMessageInfo

func (m *Unbonding) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Unbonding) GetCompletionTime() time.Time {
	if m != nil {
		return m.CompletionTime
	}
	return time.Time{}
}

func (m *Unbonding) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*Unbonding)(nil), "sunrise.shareclass.v1.Unbonding")
}

func init() {
	proto.RegisterFile("sunrise/shareclass/v1/unbonding.proto", fileDescriptor_ca9cb0aa09d0e4a3)
}

var fileDescriptor_ca9cb0aa09d0e4a3 = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x51, 0xb1, 0x4e, 0xeb, 0x30,
	0x14, 0x8d, 0xdf, 0x7b, 0xea, 0x7b, 0xcd, 0x93, 0x40, 0x8a, 0x8a, 0x94, 0x76, 0x70, 0x2b, 0x24,
	0xa4, 0x2e, 0xd8, 0x4a, 0x19, 0x98, 0x09, 0x2b, 0x2c, 0x05, 0x16, 0x96, 0xca, 0x49, 0x8d, 0x6b,
	0x29, 0xf1, 0x8d, 0x62, 0xa7, 0xa2, 0x7f, 0xd1, 0x8f, 0xe1, 0x1f, 0xe8, 0x58, 0x31, 0x31, 0x01,
	0x6a, 0x7f, 0x04, 0x25, 0x76, 0x54, 0x36, 0x9f, 0x7b, 0xce, 0x3d, 0xe7, 0xea, 0xd8, 0x3f, 0xd3,
	0x95, 0x2a, 0xa5, 0xe6, 0x54, 0x2f, 0x58, 0xc9, 0xd3, 0x8c, 0x69, 0x4d, 0x97, 0x11, 0xad, 0x54,
	0x02, 0x6a, 0x2e, 0x95, 0x20, 0x45, 0x09, 0x06, 0x82, 0x13, 0x27, 0x23, 0x07, 0x19, 0x59, 0x46,
	0x03, 0x9c, 0x82, 0xce, 0x41, 0xd3, 0x84, 0x69, 0x4e, 0x97, 0x51, 0xc2, 0x0d, 0x8b, 0x68, 0x0a,
	0x52, 0xd9, 0xb5, 0x41, 0xdf, 0xf2, 0xb3, 0x06, 0x51, 0x0b, 0x1c, 0xd5, 0x13, 0x20, 0xc0, 0xce,
	0xeb, 0x97, 0x9b, 0x0e, 0x05, 0x80, 0xc8, 0x38, 0x6d, 0x50, 0x52, 0x3d, 0x51, 0x23, 0x73, 0xae,
	0x0d, 0xcb, 0x0b, 0x2b, 0x38, 0x7d, 0x45, 0x7e, 0xf7, 0xa1, 0x3d, 0x2e, 0x98, 0xf8, 0x7f, 0xd9,
	0x7c, 0x5e, 0x72, 0xad, 0x43, 0x34, 0x42, 0xe3, 0x6e, 0x1c, 0xbe, 0xbd, 0x9c, 0xf7, 0x5c, 0xce,
	0x95, 0x65, 0xee, 0x4c, 0x29, 0x95, 0x98, 0xb6, 0xc2, 0xe0, 0xd6, 0x3f, 0x4e, 0x21, 0x2f, 0x32,
	0x6e, 0x24, 0xa8, 0x59, 0xed, 0x1f, 0xfe, 0x1a, 0xa1, 0xf1, 0xff, 0xc9, 0x80, 0xd8, 0x70, 0xd2,
	0x86, 0x93, 0xfb, 0x36, 0x3c, 0xfe, 0xb7, 0xf9, 0x18, 0x7a, 0xeb, 0xcf, 0x21, 0x9a, 0x1e, 0x1d,
	0x96, 0x6b, 0x3a, 0xb8, 0xf4, 0x3b, 0x2c, 0x87, 0x4a, 0x99, 0xf0, 0x77, 0xe3, 0xd2, 0x27, 0x2e,
	0xbe, 0xee, 0x84, 0xb8, 0x4e, 0xc8, 0x35, 0x48, 0x15, 0xff, 0xa9, 0x4d, 0xa6, 0x4e, 0x1e, 0xdf,
	0x6c, 0x76, 0x18, 0x6d, 0x77, 0x18, 0x7d, 0xed, 0x30, 0x5a, 0xef, 0xb1, 0xb7, 0xdd, 0x63, 0xef,
	0x7d, 0x8f, 0xbd, 0xc7, 0x89, 0x90, 0x66, 0x51, 0x25, 0x24, 0x85, 0x9c, 0xba, 0xde, 0x33, 0xb6,
	0xe2, 0x65, 0x0b, 0xe8, 0xf3, 0xcf, 0xdf, 0x32, 0xab, 0x82, 0xeb, 0xa4, 0xd3, 0x1c, 0x7d, 0xf1,
	0x1d, 0x00, 0x00, 0xff, 0xff, 0xca, 0x22, 0x9b, 0x2d, 0xd0, 0x01, 0x00, 0x00,
}

func (m *Unbonding) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Unbonding) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Unbonding) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintUnbonding(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.CompletionTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.CompletionTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintUnbonding(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintUnbonding(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintUnbonding(dAtA []byte, offset int, v uint64) int {
	offset -= sovUnbonding(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Unbonding) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovUnbonding(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.CompletionTime)
	n += 1 + l + sovUnbonding(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovUnbonding(uint64(l))
	return n
}

func sovUnbonding(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUnbonding(x uint64) (n int) {
	return sovUnbonding(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Unbonding) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUnbonding
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
			return fmt.Errorf("proto: Unbonding: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Unbonding: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbonding
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
				return ErrInvalidLengthUnbonding
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUnbonding
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CompletionTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbonding
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
				return ErrInvalidLengthUnbonding
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUnbonding
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.CompletionTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbonding
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
				return ErrInvalidLengthUnbonding
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUnbonding
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUnbonding(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUnbonding
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
func skipUnbonding(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUnbonding
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
					return 0, ErrIntOverflowUnbonding
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
					return 0, ErrIntOverflowUnbonding
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
				return 0, ErrInvalidLengthUnbonding
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUnbonding
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUnbonding
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUnbonding        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUnbonding          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUnbonding = fmt.Errorf("proto: unexpected end of group")
)
