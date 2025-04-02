// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquiditypool/v1/genesis.proto

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

// GenesisState defines the liquiditypool module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params               Params                `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Pools                []Pool                `protobuf:"bytes,2,rep,name=pools,proto3" json:"pools"`
	PoolCount            uint64                `protobuf:"varint,3,opt,name=pool_count,json=poolCount,proto3" json:"pool_count,omitempty"`
	Positions            []Position            `protobuf:"bytes,4,rep,name=positions,proto3" json:"positions"`
	PositionCount        uint64                `protobuf:"varint,5,opt,name=position_count,json=positionCount,proto3" json:"position_count,omitempty"`
	Accumulators         []AccumulatorObject   `protobuf:"bytes,6,rep,name=accumulators,proto3" json:"accumulators"`
	AccumulatorPositions []AccumulatorPosition `protobuf:"bytes,7,rep,name=accumulator_positions,json=accumulatorPositions,proto3" json:"accumulator_positions"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_96f15cd388d63dc3, []int{0}
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

func (m *GenesisState) GetPools() []Pool {
	if m != nil {
		return m.Pools
	}
	return nil
}

func (m *GenesisState) GetPoolCount() uint64 {
	if m != nil {
		return m.PoolCount
	}
	return 0
}

func (m *GenesisState) GetPositions() []Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *GenesisState) GetPositionCount() uint64 {
	if m != nil {
		return m.PositionCount
	}
	return 0
}

func (m *GenesisState) GetAccumulators() []AccumulatorObject {
	if m != nil {
		return m.Accumulators
	}
	return nil
}

func (m *GenesisState) GetAccumulatorPositions() []AccumulatorPosition {
	if m != nil {
		return m.AccumulatorPositions
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "sunrise.liquiditypool.v1.GenesisState")
}

func init() {
	proto.RegisterFile("sunrise/liquiditypool/v1/genesis.proto", fileDescriptor_96f15cd388d63dc3)
}

var fileDescriptor_96f15cd388d63dc3 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4b, 0xfb, 0x30,
	0x18, 0xc7, 0xdb, 0xdf, 0xfe, 0xfc, 0x58, 0x36, 0x3d, 0x94, 0x09, 0x65, 0x60, 0x2c, 0x93, 0xe9,
	0x50, 0x6c, 0xd9, 0xc4, 0x8b, 0x07, 0xc1, 0x09, 0x7a, 0x9c, 0x4c, 0xbc, 0x78, 0x19, 0x59, 0x0d,
	0x5d, 0xa4, 0x6b, 0x6a, 0x93, 0x0e, 0xf7, 0x2e, 0x7c, 0x59, 0x3b, 0xee, 0xa8, 0x17, 0x91, 0xed,
	0x8d, 0x48, 0x93, 0x94, 0xfd, 0x91, 0x0e, 0x6f, 0x79, 0x9e, 0x7c, 0xbe, 0xdf, 0xef, 0xf3, 0x90,
	0x80, 0x23, 0x16, 0x07, 0x11, 0x61, 0xd8, 0xf1, 0xc9, 0x6b, 0x4c, 0x9e, 0x09, 0x9f, 0x84, 0x94,
	0xfa, 0xce, 0xb8, 0xe5, 0x78, 0x38, 0xc0, 0x8c, 0x30, 0x3b, 0x8c, 0x28, 0xa7, 0x86, 0xa9, 0x38,
	0x7b, 0x8d, 0xb3, 0xc7, 0xad, 0x5a, 0xd5, 0xa3, 0x1e, 0x15, 0x90, 0x93, 0x9c, 0x24, 0x5f, 0x3b,
	0xc9, 0xf4, 0x45, 0xae, 0x1b, 0x8f, 0x62, 0x1f, 0x71, 0x1a, 0x29, 0xb6, 0x91, 0xc9, 0x86, 0x28,
	0x42, 0x23, 0x35, 0x42, 0xed, 0x30, 0x1b, 0x4b, 0x46, 0x91, 0xd0, 0xf1, 0x16, 0x88, 0x11, 0x4e,
	0x68, 0x20, 0xc1, 0xfa, 0x67, 0x0e, 0x54, 0xee, 0xe4, 0x8a, 0x0f, 0x1c, 0x71, 0x6c, 0x5c, 0x81,
	0xa2, 0x8c, 0x33, 0x75, 0x4b, 0x6f, 0x96, 0xdb, 0x96, 0x9d, 0xb5, 0xb2, 0x7d, 0x2f, 0xb8, 0x4e,
	0x7e, 0xfa, 0x75, 0xa0, 0xf5, 0x94, 0xca, 0xb8, 0x04, 0x85, 0xe4, 0x9e, 0x99, 0xff, 0xac, 0x5c,
	0xb3, 0xdc, 0x86, 0x5b, 0xe4, 0x94, 0xfa, 0x4a, 0x2c, 0x25, 0xc6, 0x3e, 0x00, 0xc9, 0xa1, 0xef,
	0xd2, 0x38, 0xe0, 0x66, 0xce, 0xd2, 0x9b, 0xf9, 0x5e, 0x29, 0xe9, 0xdc, 0x24, 0x0d, 0xe3, 0x16,
	0x94, 0xd2, 0xe9, 0x99, 0x99, 0x17, 0xf6, 0xf5, 0x6d, 0xf6, 0x12, 0x55, 0x11, 0x4b, 0xa9, 0xd1,
	0x00, 0xbb, 0x69, 0xa1, 0xa2, 0x0a, 0x22, 0x6a, 0x27, 0xed, 0xca, 0xb8, 0x47, 0x50, 0x59, 0x79,
	0x24, 0x66, 0x16, 0x45, 0xe2, 0x69, 0x76, 0xe2, 0xf5, 0x92, 0xee, 0x0e, 0x5e, 0xb0, 0xcb, 0x55,
	0xf4, 0x9a, 0x8d, 0x31, 0x04, 0x7b, 0x2b, 0x75, 0x7f, 0xb9, 0xd1, 0x7f, 0xe1, 0x7f, 0xf6, 0x27,
	0xff, 0x8d, 0xe5, 0xaa, 0xe8, 0xf7, 0x15, 0xeb, 0x74, 0xa7, 0x73, 0xa8, 0xcf, 0xe6, 0x50, 0xff,
	0x9e, 0x43, 0xfd, 0x7d, 0x01, 0xb5, 0xd9, 0x02, 0x6a, 0x1f, 0x0b, 0xa8, 0x3d, 0x5d, 0x78, 0x84,
	0x0f, 0xe3, 0x81, 0xed, 0xd2, 0x91, 0xa3, 0xe2, 0x7c, 0x34, 0xc1, 0x51, 0x5a, 0x38, 0x6f, 0x1b,
	0x1f, 0x87, 0x4f, 0x42, 0xcc, 0x06, 0x45, 0xf1, 0x67, 0xce, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x4f, 0xc5, 0x48, 0x01, 0x2e, 0x03, 0x00, 0x00,
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
	if len(m.AccumulatorPositions) > 0 {
		for iNdEx := len(m.AccumulatorPositions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AccumulatorPositions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Accumulators) > 0 {
		for iNdEx := len(m.Accumulators) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Accumulators[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.PositionCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.PositionCount))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.PoolCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.PoolCount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.PoolCount != 0 {
		n += 1 + sovGenesis(uint64(m.PoolCount))
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.PositionCount != 0 {
		n += 1 + sovGenesis(uint64(m.PositionCount))
	}
	if len(m.Accumulators) > 0 {
		for _, e := range m.Accumulators {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AccumulatorPositions) > 0 {
		for _, e := range m.AccumulatorPositions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
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
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
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
			m.Pools = append(m.Pools, Pool{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolCount", wireType)
			}
			m.PoolCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
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
			m.Positions = append(m.Positions, Position{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PositionCount", wireType)
			}
			m.PositionCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PositionCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Accumulators", wireType)
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
			m.Accumulators = append(m.Accumulators, AccumulatorObject{})
			if err := m.Accumulators[len(m.Accumulators)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccumulatorPositions", wireType)
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
			m.AccumulatorPositions = append(m.AccumulatorPositions, AccumulatorPosition{})
			if err := m.AccumulatorPositions[len(m.AccumulatorPositions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
