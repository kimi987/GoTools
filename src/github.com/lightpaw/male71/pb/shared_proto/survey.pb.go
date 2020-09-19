// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/shared_proto/survey.proto
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

// 问卷调查配置
type SurveyDataProto struct {
	Id        string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string                `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Icon      string                `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
	Url       string                `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Condition *UnlockConditionProto `protobuf:"bytes,5,opt,name=condition" json:"condition,omitempty"`
	Prize     *PrizeProto           `protobuf:"bytes,6,opt,name=prize" json:"prize,omitempty"`
}

func (m *SurveyDataProto) Reset()                    { *m = SurveyDataProto{} }
func (m *SurveyDataProto) String() string            { return proto.CompactTextString(m) }
func (*SurveyDataProto) ProtoMessage()               {}
func (*SurveyDataProto) Descriptor() ([]byte, []int) { return fileDescriptorSurvey, []int{0} }

func (m *SurveyDataProto) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SurveyDataProto) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SurveyDataProto) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *SurveyDataProto) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *SurveyDataProto) GetCondition() *UnlockConditionProto {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (m *SurveyDataProto) GetPrize() *PrizeProto {
	if m != nil {
		return m.Prize
	}
	return nil
}

type HeroSurveyProto struct {
	CompeleteSurvey []string `protobuf:"bytes,1,rep,name=compelete_survey,json=compeleteSurvey" json:"compelete_survey,omitempty"`
}

func (m *HeroSurveyProto) Reset()                    { *m = HeroSurveyProto{} }
func (m *HeroSurveyProto) String() string            { return proto.CompactTextString(m) }
func (*HeroSurveyProto) ProtoMessage()               {}
func (*HeroSurveyProto) Descriptor() ([]byte, []int) { return fileDescriptorSurvey, []int{1} }

func (m *HeroSurveyProto) GetCompeleteSurvey() []string {
	if m != nil {
		return m.CompeleteSurvey
	}
	return nil
}

func init() {
	proto.RegisterType((*SurveyDataProto)(nil), "proto.SurveyDataProto")
	proto.RegisterType((*HeroSurveyProto)(nil), "proto.HeroSurveyProto")
}
func (m *SurveyDataProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SurveyDataProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(len(m.Id)))
		i += copy(dAtA[i:], m.Id)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Icon) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(len(m.Icon)))
		i += copy(dAtA[i:], m.Icon)
	}
	if len(m.Url) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(len(m.Url)))
		i += copy(dAtA[i:], m.Url)
	}
	if m.Condition != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(m.Condition.Size()))
		n1, err := m.Condition.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Prize != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintSurvey(dAtA, i, uint64(m.Prize.Size()))
		n2, err := m.Prize.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *HeroSurveyProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroSurveyProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.CompeleteSurvey) > 0 {
		for _, s := range m.CompeleteSurvey {
			dAtA[i] = 0xa
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func encodeFixed64Survey(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Survey(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintSurvey(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *SurveyDataProto) Size() (n int) {
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovSurvey(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovSurvey(uint64(l))
	}
	l = len(m.Icon)
	if l > 0 {
		n += 1 + l + sovSurvey(uint64(l))
	}
	l = len(m.Url)
	if l > 0 {
		n += 1 + l + sovSurvey(uint64(l))
	}
	if m.Condition != nil {
		l = m.Condition.Size()
		n += 1 + l + sovSurvey(uint64(l))
	}
	if m.Prize != nil {
		l = m.Prize.Size()
		n += 1 + l + sovSurvey(uint64(l))
	}
	return n
}

func (m *HeroSurveyProto) Size() (n int) {
	var l int
	_ = l
	if len(m.CompeleteSurvey) > 0 {
		for _, s := range m.CompeleteSurvey {
			l = len(s)
			n += 1 + l + sovSurvey(uint64(l))
		}
	}
	return n
}

func sovSurvey(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSurvey(x uint64) (n int) {
	return sovSurvey(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SurveyDataProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSurvey
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
			return fmt.Errorf("proto: SurveyDataProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SurveyDataProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Icon", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Icon = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Url", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Url = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Condition", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
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
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Condition == nil {
				m.Condition = &UnlockConditionProto{}
			}
			if err := m.Condition.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prize", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
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
				return ErrInvalidLengthSurvey
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
		default:
			iNdEx = preIndex
			skippy, err := skipSurvey(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSurvey
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
func (m *HeroSurveyProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSurvey
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
			return fmt.Errorf("proto: HeroSurveyProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroSurveyProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CompeleteSurvey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSurvey
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSurvey
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CompeleteSurvey = append(m.CompeleteSurvey, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSurvey(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSurvey
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
func skipSurvey(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSurvey
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
					return 0, ErrIntOverflowSurvey
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
					return 0, ErrIntOverflowSurvey
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
				return 0, ErrInvalidLengthSurvey
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSurvey
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
				next, err := skipSurvey(dAtA[start:])
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
	ErrInvalidLengthSurvey = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSurvey   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/pb/shared_proto/survey.proto", fileDescriptorSurvey)
}

var fileDescriptorSurvey = []byte{
	// 289 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0xc6, 0x4d, 0xbb, 0x5d, 0x68, 0x14, 0xbb, 0xe6, 0x14, 0x14, 0x4a, 0xd9, 0x8b, 0xf5, 0xd2,
	0x82, 0x8a, 0x22, 0x78, 0x52, 0x0f, 0x7b, 0x5c, 0x2a, 0x5e, 0xbc, 0x2c, 0x69, 0x1b, 0xb6, 0xc1,
	0xb6, 0x09, 0x69, 0xaa, 0xe8, 0x93, 0xf8, 0x36, 0x5e, 0x3d, 0xfa, 0x08, 0x52, 0x5f, 0x44, 0x3a,
	0xf1, 0xdf, 0xd1, 0x53, 0xbe, 0xfc, 0xe6, 0xfb, 0x66, 0x98, 0xc1, 0x27, 0x6b, 0x61, 0xaa, 0x3e,
	0x4f, 0x0a, 0xd9, 0xa4, 0xb5, 0x58, 0x57, 0x46, 0xb1, 0x87, 0xb4, 0x61, 0x35, 0x3f, 0x4d, 0x55,
	0x9e, 0x76, 0x15, 0xd3, 0xbc, 0x5c, 0x29, 0x2d, 0x8d, 0x4c, 0xbb, 0x5e, 0xdf, 0xf3, 0xc7, 0x04,
	0x3e, 0xc4, 0x83, 0x67, 0xf7, 0xf8, 0xff, 0xf1, 0x9c, 0x75, 0xdc, 0x86, 0xe7, 0x2f, 0x08, 0x07,
	0xd7, 0xd0, 0xed, 0x8a, 0x19, 0xb6, 0x84, 0x86, 0xdb, 0xd8, 0x11, 0x25, 0x45, 0x11, 0x8a, 0xfd,
	0xcc, 0x11, 0x25, 0x21, 0x78, 0xd2, 0xb2, 0x86, 0x53, 0x07, 0x08, 0xe8, 0x91, 0x89, 0x42, 0xb6,
	0xd4, 0xb5, 0x6c, 0xd4, 0x64, 0x86, 0xdd, 0x5e, 0xd7, 0x74, 0x02, 0x68, 0x94, 0xe4, 0x0c, 0xfb,
	0x85, 0x6c, 0x4b, 0x61, 0x84, 0x6c, 0xa9, 0x17, 0xa1, 0x78, 0xf3, 0x70, 0xcf, 0x0e, 0x4e, 0x6e,
	0xda, 0x5a, 0x16, 0x77, 0x97, 0xdf, 0x55, 0x98, 0x9c, 0xfd, 0xba, 0xc9, 0x3e, 0xf6, 0x94, 0x16,
	0x4f, 0x9c, 0x4e, 0x21, 0xb6, 0xf3, 0x15, 0x5b, 0x8e, 0xcc, 0x9a, 0x6d, 0x7d, 0x7e, 0x8e, 0x83,
	0x05, 0xd7, 0xd2, 0x2e, 0x61, 0x17, 0x38, 0xc0, 0xb3, 0x42, 0x36, 0x8a, 0xd7, 0xdc, 0xf0, 0x95,
	0xbd, 0x15, 0x45, 0x91, 0x1b, 0xfb, 0x59, 0xf0, 0xc3, 0xad, 0xff, 0x22, 0x7a, 0x1d, 0x42, 0xf4,
	0x36, 0x84, 0xe8, 0x7d, 0x08, 0xd1, 0xf3, 0x47, 0xb8, 0xb1, 0x40, 0xb7, 0x5b, 0x7f, 0x8f, 0x95,
	0x4f, 0xe1, 0x39, 0xfa, 0x0c, 0x00, 0x00, 0xff, 0xff, 0xcb, 0xc3, 0x42, 0x9f, 0x9f, 0x01, 0x00,
	0x00,
}
