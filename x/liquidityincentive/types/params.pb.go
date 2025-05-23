// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquidityincentive/v1/params.proto

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

// Params defines the parameters for the module.
type Params struct {
	EpochBlocks        int64  `protobuf:"varint,1,opt,name=epoch_blocks,json=epochBlocks,proto3" json:"epoch_blocks,omitempty"`
	StakingRewardRatio string `protobuf:"bytes,2,opt,name=staking_reward_ratio,json=stakingRewardRatio,proto3" json:"staking_reward_ratio,omitempty"`
	BribeClaimEpochs   uint64 `protobuf:"varint,3,opt,name=bribe_claim_epochs,json=bribeClaimEpochs,proto3" json:"bribe_claim_epochs,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_3be0a9cd60f422db, []int{0}
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

func (m *Params) GetEpochBlocks() int64 {
	if m != nil {
		return m.EpochBlocks
	}
	return 0
}

func (m *Params) GetStakingRewardRatio() string {
	if m != nil {
		return m.StakingRewardRatio
	}
	return ""
}

func (m *Params) GetBribeClaimEpochs() uint64 {
	if m != nil {
		return m.BribeClaimEpochs
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "sunrise.liquidityincentive.v1.Params")
}

func init() {
	proto.RegisterFile("sunrise/liquidityincentive/v1/params.proto", fileDescriptor_3be0a9cd60f422db)
}

var fileDescriptor_3be0a9cd60f422db = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x31, 0x4e, 0xc3, 0x30,
	0x14, 0x86, 0x6b, 0x5a, 0x55, 0xc2, 0x20, 0x84, 0xac, 0x0e, 0xa5, 0x12, 0xa6, 0x30, 0x55, 0x08,
	0x62, 0x55, 0x6c, 0x9d, 0x50, 0x80, 0x1d, 0x45, 0x4c, 0x2c, 0x96, 0xe3, 0x5a, 0xa9, 0xd5, 0x24,
	0x0e, 0xb6, 0x13, 0xc8, 0x2d, 0x38, 0x02, 0xdc, 0x81, 0x43, 0x30, 0x56, 0x4c, 0x8c, 0x28, 0x59,
	0x38, 0x06, 0x8a, 0x13, 0x26, 0xd8, 0xfc, 0xbe, 0xf7, 0xf9, 0xfd, 0xd2, 0x0f, 0x4f, 0x4d, 0x9e,
	0x6a, 0x69, 0x04, 0x89, 0xe5, 0x43, 0x2e, 0x97, 0xd2, 0x96, 0x32, 0xe5, 0x22, 0xb5, 0xb2, 0x10,
	0xa4, 0x98, 0x93, 0x8c, 0x69, 0x96, 0x18, 0x2f, 0xd3, 0xca, 0x2a, 0x74, 0xd8, 0xb9, 0xde, 0x5f,
	0xd7, 0x2b, 0xe6, 0x93, 0x03, 0xae, 0x4c, 0xa2, 0x0c, 0x75, 0x32, 0x69, 0x87, 0xf6, 0xe7, 0x64,
	0x14, 0xa9, 0x48, 0xb5, 0xbc, 0x79, 0xb5, 0xf4, 0xe4, 0x15, 0xc0, 0xe1, 0xad, 0x0b, 0x40, 0xc7,
	0x70, 0x57, 0x64, 0x8a, 0xaf, 0x68, 0x18, 0x2b, 0xbe, 0x36, 0x63, 0x30, 0x05, 0xb3, 0x7e, 0xb0,
	0xe3, 0x98, 0xef, 0x10, 0xba, 0x84, 0x23, 0x63, 0xd9, 0x5a, 0xa6, 0x11, 0xd5, 0xe2, 0x91, 0xe9,
	0x25, 0xd5, 0xcc, 0x4a, 0x35, 0xde, 0x9a, 0x82, 0xd9, 0xb6, 0xbf, 0xf7, 0xf1, 0x76, 0x0e, 0xbb,
	0xcc, 0x6b, 0xc1, 0x03, 0xd4, 0xb9, 0x81, 0x53, 0x83, 0xc6, 0x44, 0x67, 0x10, 0x85, 0x5a, 0x86,
	0x82, 0xf2, 0x98, 0xc9, 0x84, 0xba, 0xe3, 0x66, 0xdc, 0x9f, 0x82, 0xd9, 0x20, 0xd8, 0x77, 0x9b,
	0xab, 0x66, 0x71, 0xe3, 0xf8, 0x62, 0xf0, 0xfd, 0x72, 0x04, 0xfc, 0xbb, 0xf7, 0x0a, 0x83, 0x4d,
	0x85, 0xc1, 0x57, 0x85, 0xc1, 0x73, 0x8d, 0x7b, 0x9b, 0x1a, 0xf7, 0x3e, 0x6b, 0xdc, 0xbb, 0x5f,
	0x44, 0xd2, 0xae, 0xf2, 0xd0, 0xe3, 0x2a, 0x21, 0x5d, 0x31, 0x31, 0x2b, 0x85, 0xfe, 0x1d, 0xc8,
	0xd3, 0x7f, 0x9d, 0xda, 0x32, 0x13, 0x26, 0x1c, 0xba, 0x02, 0x2e, 0x7e, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xcb, 0xe8, 0x59, 0x0a, 0x7e, 0x01, 0x00, 0x00,
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
	if this.EpochBlocks != that1.EpochBlocks {
		return false
	}
	if this.StakingRewardRatio != that1.StakingRewardRatio {
		return false
	}
	if this.BribeClaimEpochs != that1.BribeClaimEpochs {
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
	if m.BribeClaimEpochs != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.BribeClaimEpochs))
		i--
		dAtA[i] = 0x18
	}
	if len(m.StakingRewardRatio) > 0 {
		i -= len(m.StakingRewardRatio)
		copy(dAtA[i:], m.StakingRewardRatio)
		i = encodeVarintParams(dAtA, i, uint64(len(m.StakingRewardRatio)))
		i--
		dAtA[i] = 0x12
	}
	if m.EpochBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EpochBlocks))
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
	if m.EpochBlocks != 0 {
		n += 1 + sovParams(uint64(m.EpochBlocks))
	}
	l = len(m.StakingRewardRatio)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.BribeClaimEpochs != 0 {
		n += 1 + sovParams(uint64(m.BribeClaimEpochs))
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
				return fmt.Errorf("proto: wrong wireType = %d for field EpochBlocks", wireType)
			}
			m.EpochBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochBlocks |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakingRewardRatio", wireType)
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
			m.StakingRewardRatio = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BribeClaimEpochs", wireType)
			}
			m.BribeClaimEpochs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BribeClaimEpochs |= uint64(b&0x7F) << shift
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
