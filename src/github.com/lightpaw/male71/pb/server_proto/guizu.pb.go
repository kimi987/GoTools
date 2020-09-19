// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/server_proto/guizu.proto
// DO NOT EDIT!

package server_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type HeroGuiZuServerProto struct {
	CollectedLevels []uint64 `protobuf:"varint,1,rep,packed,name=collected_levels,json=collectedLevels" json:"collected_levels,omitempty"`
}

func (m *HeroGuiZuServerProto) Reset()                    { *m = HeroGuiZuServerProto{} }
func (m *HeroGuiZuServerProto) String() string            { return proto.CompactTextString(m) }
func (*HeroGuiZuServerProto) ProtoMessage()               {}
func (*HeroGuiZuServerProto) Descriptor() ([]byte, []int) { return fileDescriptorGuizu, []int{0} }

func (m *HeroGuiZuServerProto) GetCollectedLevels() []uint64 {
	if m != nil {
		return m.CollectedLevels
	}
	return nil
}

func init() {
	proto.RegisterType((*HeroGuiZuServerProto)(nil), "proto.HeroGuiZuServerProto")
}
func (m *HeroGuiZuServerProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroGuiZuServerProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.CollectedLevels) > 0 {
		dAtA2 := make([]byte, len(m.CollectedLevels)*10)
		var j1 int
		for _, num := range m.CollectedLevels {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		dAtA[i] = 0xa
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(j1))
		i += copy(dAtA[i:], dAtA2[:j1])
	}
	return i, nil
}

func encodeFixed64Guizu(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Guizu(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintGuizu(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *HeroGuiZuServerProto) Size() (n int) {
	var l int
	_ = l
	if len(m.CollectedLevels) > 0 {
		l = 0
		for _, e := range m.CollectedLevels {
			l += sovGuizu(uint64(e))
		}
		n += 1 + sovGuizu(uint64(l)) + l
	}
	return n
}

func sovGuizu(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGuizu(x uint64) (n int) {
	return sovGuizu(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HeroGuiZuServerProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGuizu
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HeroGuiZuServerProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroGuiZuServerProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGuizu
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.CollectedLevels = append(m.CollectedLevels, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGuizu
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGuizu
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGuizu
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.CollectedLevels = append(m.CollectedLevels, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field CollectedLevels", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGuizu(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGuizu
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
func skipGuizu(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGuizu
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
					return 0, ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGuizu
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthGuizu
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGuizu
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGuizu(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGuizu = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGuizu   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/pb/server_proto/guizu.proto", fileDescriptorGuizu)
}

var fileDescriptorGuizu = []byte{
	// 160 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4d, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0xc9, 0x4c, 0xcf, 0x28, 0x29, 0x48, 0x2c, 0xd7,
	0xcf, 0x4d, 0xcc, 0x49, 0x35, 0xd7, 0x2f, 0x48, 0xd2, 0x2f, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0x8a,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x2f, 0xcd, 0xac, 0x2a, 0xd5, 0x03, 0xb3, 0x85, 0x58,
	0xc1, 0x94, 0x92, 0x23, 0x97, 0x88, 0x47, 0x6a, 0x51, 0xbe, 0x7b, 0x69, 0x66, 0x54, 0x69, 0x30,
	0x58, 0x6d, 0x00, 0x58, 0x5a, 0x93, 0x4b, 0x20, 0x39, 0x3f, 0x27, 0x27, 0x35, 0xb9, 0x24, 0x35,
	0x25, 0x3e, 0x27, 0xb5, 0x2c, 0x35, 0xa7, 0x58, 0x82, 0x51, 0x81, 0x59, 0x83, 0x25, 0x88, 0x1f,
	0x2e, 0xee, 0x03, 0x16, 0x76, 0x52, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07,
	0x8f, 0xe4, 0x18, 0x67, 0x3c, 0x96, 0x63, 0xf0, 0x60, 0x8c, 0xe2, 0x41, 0xb6, 0x37, 0x89, 0x0d,
	0x4c, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xee, 0xf3, 0xe8, 0xe7, 0xab, 0x00, 0x00, 0x00,
}
