// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ocea/crypto/multisig/keys.proto

package multisig

import (
	fmt "fmt"
	types "github.com/ocea/sdk/codec/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// LegacyAminoPubKey specifies a public key type
// which nests multiple public keys and a threshold,
// it uses legacy amino address rules.
type LegacyAminoPubKey struct {
	Threshold uint32       `protobuf:"varint,1,opt,name=threshold,proto3" json:"threshold,omitempty" yaml:"threshold"`
	PubKeys   []*types.Any `protobuf:"bytes,2,rep,name=public_keys,json=publicKeys,proto3" json:"public_keys,omitempty" yaml:"pubkeys"`
}

func (m *LegacyAminoPubKey) Reset()         { *m = LegacyAminoPubKey{} }
func (m *LegacyAminoPubKey) String() string { return proto.CompactTextString(m) }
func (*LegacyAminoPubKey) ProtoMessage()    {}
func (*LegacyAminoPubKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_46b57537e097d47d, []int{0}
}
func (m *LegacyAminoPubKey) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyAminoPubKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyAminoPubKey.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyAminoPubKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyAminoPubKey.Merge(m, src)
}
func (m *LegacyAminoPubKey) XXX_Size() int {
	return m.Size()
}
func (m *LegacyAminoPubKey) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyAminoPubKey.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyAminoPubKey proto.InternalMessageInfo

func init() {
	proto.RegisterType((*LegacyAminoPubKey)(nil), "ocea.crypto.multisig.LegacyAminoPubKey")
}

func init() { proto.RegisterFile("ocea/crypto/multisig/keys.proto", fileDescriptor_46b57537e097d47d) }

var fileDescriptor_46b57537e097d47d = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0xce, 0x2f, 0xce,
	0xcd, 0x2f, 0xd6, 0x4f, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0xd7, 0xcf, 0x2d, 0xcd, 0x29, 0xc9, 0x2c,
	0xce, 0x4c, 0xd7, 0xcf, 0x4e, 0xad, 0x2c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x83,
	0x28, 0xd1, 0x83, 0x28, 0xd1, 0x83, 0x29, 0x91, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x2b, 0xd1,
	0x07, 0xb1, 0x20, 0xaa, 0xa5, 0x24, 0xd3, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xbc, 0xa4,
	0xd2, 0x34, 0xfd, 0xc4, 0xbc, 0x4a, 0x88, 0x94, 0xd2, 0x62, 0x46, 0x2e, 0x41, 0x9f, 0xd4, 0xf4,
	0xc4, 0xe4, 0x4a, 0xc7, 0xdc, 0xcc, 0xbc, 0xfc, 0x80, 0xd2, 0x24, 0xef, 0xd4, 0x4a, 0x21, 0x23,
	0x2e, 0xce, 0x92, 0x8c, 0xa2, 0xd4, 0xe2, 0x8c, 0xfc, 0x9c, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d,
	0x5e, 0x27, 0x91, 0x4f, 0xf7, 0xe4, 0x05, 0x2a, 0x13, 0x73, 0x73, 0xac, 0x94, 0xe0, 0x52, 0x4a,
	0x41, 0x08, 0x65, 0x42, 0x21, 0x5c, 0xdc, 0x05, 0xa5, 0x49, 0x39, 0x99, 0xc9, 0xf1, 0x20, 0x77,
	0x4a, 0x30, 0x29, 0x30, 0x6b, 0x70, 0x1b, 0x89, 0xe8, 0x41, 0xac, 0xd6, 0x83, 0x59, 0xad, 0xe7,
	0x98, 0x57, 0xe9, 0x24, 0xfb, 0xe8, 0x9e, 0x3c, 0x3b, 0xc4, 0xaa, 0xe2, 0x4f, 0xf7, 0xe4, 0xf9,
	0x20, 0xc6, 0x16, 0x94, 0x26, 0x81, 0x74, 0x2a, 0x05, 0x71, 0x41, 0xcc, 0x01, 0xc9, 0x5a, 0xb1,
	0x74, 0x2c, 0x90, 0x67, 0x70, 0xf2, 0x3e, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07,
	0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86,
	0x28, 0xc3, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0x58, 0xb0, 0x81,
	0x29, 0xdd, 0xe2, 0x94, 0x6c, 0x58, 0x08, 0x82, 0x8c, 0x85, 0x07, 0x63, 0x12, 0x1b, 0xd8, 0x2d,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x76, 0x6e, 0xbf, 0x3c, 0x67, 0x01, 0x00, 0x00,
}

func (m *LegacyAminoPubKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyAminoPubKey) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyAminoPubKey) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PubKeys) > 0 {
		for iNdEx := len(m.PubKeys) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PubKeys[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintKeys(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Threshold != 0 {
		i = encodeVarintKeys(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintKeys(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeys(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LegacyAminoPubKey) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Threshold != 0 {
		n += 1 + sovKeys(uint64(m.Threshold))
	}
	if len(m.PubKeys) > 0 {
		for _, e := range m.PubKeys {
			l = e.Size()
			n += 1 + l + sovKeys(uint64(l))
		}
	}
	return n
}

func sovKeys(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeys(x uint64) (n int) {
	return sovKeys(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LegacyAminoPubKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeys
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
			return fmt.Errorf("proto: LegacyAminoPubKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyAminoPubKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKeys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeys
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
				return ErrInvalidLengthKeys
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthKeys
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKeys = append(m.PubKeys, &types.Any{})
			if err := m.PubKeys[len(m.PubKeys)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipKeys(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthKeys
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthKeys
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
func skipKeys(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeys
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
					return 0, ErrIntOverflowKeys
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
					return 0, ErrIntOverflowKeys
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
				return 0, ErrInvalidLengthKeys
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeys
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeys
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthKeys        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowKeys          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupKeys = fmt.Errorf("proto: unexpected end of group")
)
