// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/shared_proto/guizu.proto
// DO NOT EDIT!

package shared_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GuiZuLevelDataProto struct {
	Level       int32       `protobuf:"varint,1,opt,name=level,proto3" json:"level,omitempty"`
	Name        int32       `protobuf:"varint,4,opt,name=name,proto3" json:"name,omitempty"`
	Icon        int32       `protobuf:"varint,5,opt,name=icon,proto3" json:"icon,omitempty"`
	HistoryJade int32       `protobuf:"varint,2,opt,name=history_jade,json=historyJade,proto3" json:"history_jade,omitempty"`
	Prize       *PrizeProto `protobuf:"bytes,3,opt,name=prize" json:"prize,omitempty"`
}

func (m *GuiZuLevelDataProto) Reset()                    { *m = GuiZuLevelDataProto{} }
func (m *GuiZuLevelDataProto) String() string            { return proto.CompactTextString(m) }
func (*GuiZuLevelDataProto) ProtoMessage()               {}
func (*GuiZuLevelDataProto) Descriptor() ([]byte, []int) { return fileDescriptorGuizu, []int{0} }

func (m *GuiZuLevelDataProto) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *GuiZuLevelDataProto) GetName() int32 {
	if m != nil {
		return m.Name
	}
	return 0
}

func (m *GuiZuLevelDataProto) GetIcon() int32 {
	if m != nil {
		return m.Icon
	}
	return 0
}

func (m *GuiZuLevelDataProto) GetHistoryJade() int32 {
	if m != nil {
		return m.HistoryJade
	}
	return 0
}

func (m *GuiZuLevelDataProto) GetPrize() *PrizeProto {
	if m != nil {
		return m.Prize
	}
	return nil
}

type HeroGuiZuProto struct {
	CollectedLevels []int32 `protobuf:"varint,1,rep,name=collected_levels,json=collectedLevels" json:"collected_levels,omitempty"`
}

func (m *HeroGuiZuProto) Reset()                    { *m = HeroGuiZuProto{} }
func (m *HeroGuiZuProto) String() string            { return proto.CompactTextString(m) }
func (*HeroGuiZuProto) ProtoMessage()               {}
func (*HeroGuiZuProto) Descriptor() ([]byte, []int) { return fileDescriptorGuizu, []int{1} }

func (m *HeroGuiZuProto) GetCollectedLevels() []int32 {
	if m != nil {
		return m.CollectedLevels
	}
	return nil
}

func init() {
	proto.RegisterType((*GuiZuLevelDataProto)(nil), "proto.GuiZuLevelDataProto")
	proto.RegisterType((*HeroGuiZuProto)(nil), "proto.HeroGuiZuProto")
}
func (m *GuiZuLevelDataProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GuiZuLevelDataProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Level != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(m.Level))
	}
	if m.HistoryJade != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(m.HistoryJade))
	}
	if m.Prize != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(m.Prize.Size()))
		n1, err := m.Prize.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Name != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(m.Name))
	}
	if m.Icon != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintGuizu(dAtA, i, uint64(m.Icon))
	}
	return i, nil
}

func (m *HeroGuiZuProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroGuiZuProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.CollectedLevels) > 0 {
		for _, num := range m.CollectedLevels {
			dAtA[i] = 0x8
			i++
			i = encodeVarintGuizu(dAtA, i, uint64(num))
		}
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
func (m *GuiZuLevelDataProto) Size() (n int) {
	var l int
	_ = l
	if m.Level != 0 {
		n += 1 + sovGuizu(uint64(m.Level))
	}
	if m.HistoryJade != 0 {
		n += 1 + sovGuizu(uint64(m.HistoryJade))
	}
	if m.Prize != nil {
		l = m.Prize.Size()
		n += 1 + l + sovGuizu(uint64(l))
	}
	if m.Name != 0 {
		n += 1 + sovGuizu(uint64(m.Name))
	}
	if m.Icon != 0 {
		n += 1 + sovGuizu(uint64(m.Icon))
	}
	return n
}

func (m *HeroGuiZuProto) Size() (n int) {
	var l int
	_ = l
	if len(m.CollectedLevels) > 0 {
		for _, e := range m.CollectedLevels {
			n += 1 + sovGuizu(uint64(e))
		}
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
func (m *GuiZuLevelDataProto) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GuiZuLevelDataProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GuiZuLevelDataProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Level", wireType)
			}
			m.Level = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Level |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HistoryJade", wireType)
			}
			m.HistoryJade = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HistoryJade |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prize", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGuizu
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Prize == nil {
				m.Prize = &PrizeProto{}
			}
			if err := m.Prize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			m.Name = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Name |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Icon", wireType)
			}
			m.Icon = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGuizu
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Icon |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
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
func (m *HeroGuiZuProto) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HeroGuiZuProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroGuiZuProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGuizu
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
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
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGuizu
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
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
	proto.RegisterFile("github.com/lightpaw/male7/pb/shared_proto/guizu.proto", fileDescriptorGuizu)
}

var fileDescriptorGuizu = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4d, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0xc9, 0x4c, 0xcf, 0x28, 0x29, 0x48, 0x2c, 0xd7,
	0xcf, 0x4d, 0xcc, 0x49, 0x35, 0xd7, 0x2f, 0x48, 0xd2, 0x2f, 0xce, 0x48, 0x2c, 0x4a, 0x4d, 0x89,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x2f, 0xcd, 0xac, 0x2a, 0xd5, 0x03, 0xb3, 0x85, 0x58,
	0xc1, 0x94, 0x94, 0x09, 0xf1, 0xba, 0x93, 0x12, 0x8b, 0x53, 0x21, 0x9a, 0x95, 0xe6, 0x33, 0x72,
	0x09, 0xbb, 0x97, 0x66, 0x46, 0x95, 0xfa, 0xa4, 0x96, 0xa5, 0xe6, 0xb8, 0x24, 0x96, 0x24, 0x06,
	0x80, 0x0d, 0x15, 0xe1, 0x62, 0xcd, 0x01, 0x89, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x41,
	0x38, 0x42, 0x8a, 0x5c, 0x3c, 0x19, 0x99, 0xc5, 0x25, 0xf9, 0x45, 0x95, 0xf1, 0x59, 0x89, 0x29,
	0xa9, 0x12, 0x4c, 0x60, 0x49, 0x6e, 0xa8, 0x98, 0x57, 0x62, 0x4a, 0xaa, 0x90, 0x3a, 0x17, 0x6b,
	0x41, 0x51, 0x66, 0x55, 0xaa, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x20, 0xc4, 0x1e, 0xbd,
	0x00, 0x90, 0x18, 0xd8, 0xe8, 0x20, 0x88, 0xbc, 0x90, 0x10, 0x17, 0x4b, 0x5e, 0x62, 0x6e, 0xaa,
	0x04, 0x0b, 0xd8, 0x0c, 0x30, 0x1b, 0x24, 0x96, 0x99, 0x9c, 0x9f, 0x27, 0xc1, 0x0a, 0x11, 0x03,
	0xb1, 0x95, 0xec, 0xb9, 0xf8, 0x3c, 0x52, 0x8b, 0xf2, 0xc1, 0x8e, 0x84, 0xb8, 0x4d, 0x97, 0x4b,
	0x20, 0x39, 0x3f, 0x27, 0x27, 0x35, 0xb9, 0x24, 0x35, 0x25, 0x1e, 0xec, 0xb0, 0x62, 0x09, 0x46,
	0x05, 0x66, 0x0d, 0x56, 0x27, 0x26, 0x01, 0x86, 0x20, 0x7e, 0xb8, 0x1c, 0xd8, 0x4b, 0xc5, 0x4e,
	0x0a, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x8c, 0xc7,
	0x72, 0x0c, 0x1e, 0x8c, 0x51, 0x3c, 0xc8, 0xe1, 0x91, 0xc4, 0x06, 0xa6, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x74, 0xc5, 0x77, 0x17, 0x81, 0x01, 0x00, 0x00,
}
