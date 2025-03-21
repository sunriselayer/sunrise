// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sunrise/liquidityincentive/v1/genesis.proto

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

// GenesisState defines the liquidityincentive module's genesis state.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params     Params  `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Epochs     []Epoch `protobuf:"bytes,2,rep,name=epochs,proto3" json:"epochs"`
	EpochCount uint64  `protobuf:"varint,3,opt,name=epoch_count,json=epochCount,proto3" json:"epoch_count,omitempty"`
	Gauges     []Gauge `protobuf:"bytes,4,rep,name=gauges,proto3" json:"gauges"`
	Votes      []Vote  `protobuf:"bytes,5,rep,name=votes,proto3" json:"votes"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_45e705b380d53cb3, []int{0}
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

func (m *GenesisState) GetEpochs() []Epoch {
	if m != nil {
		return m.Epochs
	}
	return nil
}

func (m *GenesisState) GetEpochCount() uint64 {
	if m != nil {
		return m.EpochCount
	}
	return 0
}

func (m *GenesisState) GetGauges() []Gauge {
	if m != nil {
		return m.Gauges
	}
	return nil
}

func (m *GenesisState) GetVotes() []Vote {
	if m != nil {
		return m.Votes
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "sunrise.liquidityincentive.v1.GenesisState")
}

func init() {
	proto.RegisterFile("sunrise/liquidityincentive/v1/genesis.proto", fileDescriptor_45e705b380d53cb3)
}

var fileDescriptor_45e705b380d53cb3 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xbd, 0x6a, 0xfb, 0x30,
	0x14, 0xc5, 0xad, 0x7c, 0x0d, 0xca, 0x7f, 0x32, 0xff, 0xc1, 0x04, 0xaa, 0x98, 0x7e, 0x80, 0xdb,
	0x82, 0x44, 0xd2, 0xad, 0x4b, 0xc1, 0xa1, 0x64, 0x2d, 0x69, 0xe9, 0xd0, 0xa5, 0x38, 0xae, 0x50,
	0x04, 0x89, 0xe5, 0x5a, 0xb2, 0xa9, 0xdf, 0xa2, 0x2f, 0xd4, 0x3d, 0x63, 0xc6, 0x4e, 0xa5, 0xd8,
	0x2f, 0x52, 0x2c, 0xd9, 0x53, 0x8b, 0x93, 0x4d, 0xba, 0xfc, 0xce, 0x39, 0xf7, 0x70, 0xe1, 0xa5,
	0x4c, 0xa3, 0x84, 0x4b, 0x4a, 0xd6, 0xfc, 0x35, 0xe5, 0x2f, 0x5c, 0xe5, 0x3c, 0x0a, 0x69, 0xa4,
	0x78, 0x46, 0x49, 0x36, 0x21, 0x8c, 0x46, 0x54, 0x72, 0x89, 0xe3, 0x44, 0x28, 0x61, 0x1f, 0xd5,
	0x30, 0xfe, 0x0d, 0xe3, 0x6c, 0x32, 0xfa, 0xcf, 0x04, 0x13, 0x9a, 0x24, 0xd5, 0xcb, 0x88, 0x46,
	0xe7, 0xed, 0x09, 0x34, 0x16, 0xe1, 0xea, 0x30, 0x94, 0x05, 0x29, 0xa3, 0x35, 0x7a, 0xd1, 0x8e,
	0xc6, 0x41, 0x12, 0x6c, 0xea, 0xb5, 0x8f, 0x3f, 0x3a, 0xf0, 0xdf, 0xdc, 0x14, 0xb9, 0x57, 0x81,
	0xa2, 0xf6, 0x0c, 0x0e, 0x0c, 0xe0, 0x00, 0x17, 0x78, 0xc3, 0xe9, 0x19, 0x6e, 0x2d, 0x86, 0xef,
	0x34, 0xec, 0xf7, 0xb6, 0x5f, 0x63, 0x6b, 0x51, 0x4b, 0x6d, 0x1f, 0x0e, 0xf4, 0xee, 0xd2, 0xe9,
	0xb8, 0x5d, 0x6f, 0x38, 0x3d, 0xdd, 0x63, 0x72, 0x5b, 0xc1, 0x8d, 0x87, 0x51, 0xda, 0x63, 0x38,
	0xd4, 0xaf, 0xe7, 0x50, 0xa4, 0x91, 0x72, 0xba, 0x2e, 0xf0, 0x7a, 0x0b, 0xa8, 0x47, 0xb3, 0x6a,
	0x52, 0x85, 0xe8, 0xd6, 0xd2, 0xe9, 0x1d, 0x14, 0x32, 0xaf, 0xe0, 0x26, 0xc4, 0x28, 0xed, 0x1b,
	0xd8, 0xcf, 0x84, 0xa2, 0xd2, 0xe9, 0x6b, 0x8b, 0x93, 0x3d, 0x16, 0x8f, 0x42, 0x35, 0x0e, 0x46,
	0xe7, 0x3f, 0x6c, 0x0b, 0x04, 0x76, 0x05, 0x02, 0xdf, 0x05, 0x02, 0xef, 0x25, 0xb2, 0x76, 0x25,
	0xb2, 0x3e, 0x4b, 0x64, 0x3d, 0x5d, 0x33, 0xae, 0x56, 0xe9, 0x12, 0x87, 0x62, 0x43, 0x6a, 0xd7,
	0x75, 0x90, 0xd3, 0xa4, 0xf9, 0x90, 0xb7, 0xbf, 0xee, 0xa3, 0xf2, 0x98, 0xca, 0xe5, 0x40, 0x1f,
	0xe7, 0xea, 0x27, 0x00, 0x00, 0xff, 0xff, 0xfb, 0xa2, 0xee, 0x13, 0x82, 0x02, 0x00, 0x00,
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
	if len(m.Votes) > 0 {
		for iNdEx := len(m.Votes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Votes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Gauges) > 0 {
		for iNdEx := len(m.Gauges) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Gauges[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.EpochCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.EpochCount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Epochs) > 0 {
		for iNdEx := len(m.Epochs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Epochs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Epochs) > 0 {
		for _, e := range m.Epochs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.EpochCount != 0 {
		n += 1 + sovGenesis(uint64(m.EpochCount))
	}
	if len(m.Gauges) > 0 {
		for _, e := range m.Gauges {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Votes) > 0 {
		for _, e := range m.Votes {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Epochs", wireType)
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
			m.Epochs = append(m.Epochs, Epoch{})
			if err := m.Epochs[len(m.Epochs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochCount", wireType)
			}
			m.EpochCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gauges", wireType)
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
			m.Gauges = append(m.Gauges, Gauge{})
			if err := m.Gauges[len(m.Gauges)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Votes", wireType)
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
			m.Votes = append(m.Votes, Vote{})
			if err := m.Votes[len(m.Votes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
