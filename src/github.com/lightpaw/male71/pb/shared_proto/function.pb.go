// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/shared_proto/function.proto
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

type FunctionOpenDataProto struct {
	FunctionType int32  `protobuf:"varint,1,opt,name=function_type,json=functionType,proto3" json:"function_type,omitempty"`
	GuanFuLevel  int32  `protobuf:"varint,2,opt,name=guan_fu_level,json=guanFuLevel,proto3" json:"guan_fu_level,omitempty"`
	HeroLevel    int32  `protobuf:"varint,3,opt,name=hero_level,json=heroLevel,proto3" json:"hero_level,omitempty"`
	MainTask     int32  `protobuf:"varint,4,opt,name=main_task,json=mainTask,proto3" json:"main_task,omitempty"`
	BaYeStage    int32  `protobuf:"varint,5,opt,name=ba_ye_stage,json=baYeStage,proto3" json:"ba_ye_stage,omitempty"`
	BuildingId   int32  `protobuf:"varint,6,opt,name=building_id,json=buildingId,proto3" json:"building_id,omitempty"`
	TowerFloor   int32  `protobuf:"varint,7,opt,name=tower_floor,json=towerFloor,proto3" json:"tower_floor,omitempty"`
	Dungeon      int32  `protobuf:"varint,8,opt,name=dungeon,proto3" json:"dungeon,omitempty"`
	Desc         string `protobuf:"bytes,10,opt,name=desc,proto3" json:"desc,omitempty"`
	Icon         string `protobuf:"bytes,11,opt,name=icon,proto3" json:"icon,omitempty"`
	NotifyOrder  int32  `protobuf:"varint,12,opt,name=notify_order,json=notifyOrder,proto3" json:"notify_order,omitempty"`
}

func (m *FunctionOpenDataProto) Reset()                    { *m = FunctionOpenDataProto{} }
func (m *FunctionOpenDataProto) String() string            { return proto.CompactTextString(m) }
func (*FunctionOpenDataProto) ProtoMessage()               {}
func (*FunctionOpenDataProto) Descriptor() ([]byte, []int) { return fileDescriptorFunction, []int{0} }

func (m *FunctionOpenDataProto) GetFunctionType() int32 {
	if m != nil {
		return m.FunctionType
	}
	return 0
}

func (m *FunctionOpenDataProto) GetGuanFuLevel() int32 {
	if m != nil {
		return m.GuanFuLevel
	}
	return 0
}

func (m *FunctionOpenDataProto) GetHeroLevel() int32 {
	if m != nil {
		return m.HeroLevel
	}
	return 0
}

func (m *FunctionOpenDataProto) GetMainTask() int32 {
	if m != nil {
		return m.MainTask
	}
	return 0
}

func (m *FunctionOpenDataProto) GetBaYeStage() int32 {
	if m != nil {
		return m.BaYeStage
	}
	return 0
}

func (m *FunctionOpenDataProto) GetBuildingId() int32 {
	if m != nil {
		return m.BuildingId
	}
	return 0
}

func (m *FunctionOpenDataProto) GetTowerFloor() int32 {
	if m != nil {
		return m.TowerFloor
	}
	return 0
}

func (m *FunctionOpenDataProto) GetDungeon() int32 {
	if m != nil {
		return m.Dungeon
	}
	return 0
}

func (m *FunctionOpenDataProto) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *FunctionOpenDataProto) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *FunctionOpenDataProto) GetNotifyOrder() int32 {
	if m != nil {
		return m.NotifyOrder
	}
	return 0
}

type HeroFunctionProto struct {
	OpenTypes []int32 `protobuf:"varint,1,rep,name=open_types,json=openTypes" json:"open_types,omitempty"`
}

func (m *HeroFunctionProto) Reset()                    { *m = HeroFunctionProto{} }
func (m *HeroFunctionProto) String() string            { return proto.CompactTextString(m) }
func (*HeroFunctionProto) ProtoMessage()               {}
func (*HeroFunctionProto) Descriptor() ([]byte, []int) { return fileDescriptorFunction, []int{1} }

func (m *HeroFunctionProto) GetOpenTypes() []int32 {
	if m != nil {
		return m.OpenTypes
	}
	return nil
}

func init() {
	proto.RegisterType((*FunctionOpenDataProto)(nil), "proto.FunctionOpenDataProto")
	proto.RegisterType((*HeroFunctionProto)(nil), "proto.HeroFunctionProto")
}
func (m *FunctionOpenDataProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FunctionOpenDataProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FunctionType != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.FunctionType))
	}
	if m.GuanFuLevel != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.GuanFuLevel))
	}
	if m.HeroLevel != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.HeroLevel))
	}
	if m.MainTask != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.MainTask))
	}
	if m.BaYeStage != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.BaYeStage))
	}
	if m.BuildingId != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.BuildingId))
	}
	if m.TowerFloor != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.TowerFloor))
	}
	if m.Dungeon != 0 {
		dAtA[i] = 0x40
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.Dungeon))
	}
	if len(m.Desc) > 0 {
		dAtA[i] = 0x52
		i++
		i = encodeVarintFunction(dAtA, i, uint64(len(m.Desc)))
		i += copy(dAtA[i:], m.Desc)
	}
	if len(m.Icon) > 0 {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintFunction(dAtA, i, uint64(len(m.Icon)))
		i += copy(dAtA[i:], m.Icon)
	}
	if m.NotifyOrder != 0 {
		dAtA[i] = 0x60
		i++
		i = encodeVarintFunction(dAtA, i, uint64(m.NotifyOrder))
	}
	return i, nil
}

func (m *HeroFunctionProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HeroFunctionProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.OpenTypes) > 0 {
		for _, num := range m.OpenTypes {
			dAtA[i] = 0x8
			i++
			i = encodeVarintFunction(dAtA, i, uint64(num))
		}
	}
	return i, nil
}

func encodeFixed64Function(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32Function(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintFunction(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *FunctionOpenDataProto) Size() (n int) {
	var l int
	_ = l
	if m.FunctionType != 0 {
		n += 1 + sovFunction(uint64(m.FunctionType))
	}
	if m.GuanFuLevel != 0 {
		n += 1 + sovFunction(uint64(m.GuanFuLevel))
	}
	if m.HeroLevel != 0 {
		n += 1 + sovFunction(uint64(m.HeroLevel))
	}
	if m.MainTask != 0 {
		n += 1 + sovFunction(uint64(m.MainTask))
	}
	if m.BaYeStage != 0 {
		n += 1 + sovFunction(uint64(m.BaYeStage))
	}
	if m.BuildingId != 0 {
		n += 1 + sovFunction(uint64(m.BuildingId))
	}
	if m.TowerFloor != 0 {
		n += 1 + sovFunction(uint64(m.TowerFloor))
	}
	if m.Dungeon != 0 {
		n += 1 + sovFunction(uint64(m.Dungeon))
	}
	l = len(m.Desc)
	if l > 0 {
		n += 1 + l + sovFunction(uint64(l))
	}
	l = len(m.Icon)
	if l > 0 {
		n += 1 + l + sovFunction(uint64(l))
	}
	if m.NotifyOrder != 0 {
		n += 1 + sovFunction(uint64(m.NotifyOrder))
	}
	return n
}

func (m *HeroFunctionProto) Size() (n int) {
	var l int
	_ = l
	if len(m.OpenTypes) > 0 {
		for _, e := range m.OpenTypes {
			n += 1 + sovFunction(uint64(e))
		}
	}
	return n
}

func sovFunction(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozFunction(x uint64) (n int) {
	return sovFunction(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FunctionOpenDataProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFunction
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
			return fmt.Errorf("proto: FunctionOpenDataProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FunctionOpenDataProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FunctionType", wireType)
			}
			m.FunctionType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FunctionType |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GuanFuLevel", wireType)
			}
			m.GuanFuLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GuanFuLevel |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeroLevel", wireType)
			}
			m.HeroLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HeroLevel |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MainTask", wireType)
			}
			m.MainTask = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MainTask |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaYeStage", wireType)
			}
			m.BaYeStage = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BaYeStage |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BuildingId", wireType)
			}
			m.BuildingId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BuildingId |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TowerFloor", wireType)
			}
			m.TowerFloor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TowerFloor |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Dungeon", wireType)
			}
			m.Dungeon = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Dungeon |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Desc", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
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
				return ErrInvalidLengthFunction
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Desc = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Icon", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
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
				return ErrInvalidLengthFunction
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Icon = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NotifyOrder", wireType)
			}
			m.NotifyOrder = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFunction
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NotifyOrder |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFunction(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFunction
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
func (m *HeroFunctionProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFunction
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
			return fmt.Errorf("proto: HeroFunctionProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HeroFunctionProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFunction
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
				m.OpenTypes = append(m.OpenTypes, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFunction
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
					return ErrInvalidLengthFunction
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFunction
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
					m.OpenTypes = append(m.OpenTypes, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field OpenTypes", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFunction(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFunction
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
func skipFunction(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFunction
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
					return 0, ErrIntOverflowFunction
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
					return 0, ErrIntOverflowFunction
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
				return 0, ErrInvalidLengthFunction
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowFunction
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
				next, err := skipFunction(dAtA[start:])
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
	ErrInvalidLengthFunction = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFunction   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/pb/shared_proto/function.proto", fileDescriptorFunction)
}

var fileDescriptorFunction = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xc1, 0x6a, 0xdc, 0x30,
	0x10, 0x86, 0xa3, 0x4d, 0x9c, 0xc4, 0x63, 0x07, 0x5a, 0x41, 0x41, 0x50, 0xea, 0x3a, 0xdb, 0xcb,
	0x9e, 0xe2, 0x43, 0xa1, 0xed, 0x39, 0x94, 0x25, 0x85, 0x42, 0xca, 0x36, 0x97, 0xf6, 0x22, 0x64,
	0x7b, 0x6c, 0x8b, 0x78, 0x25, 0x23, 0xcb, 0x0d, 0x7e, 0x93, 0x3e, 0x52, 0x8f, 0x7d, 0x84, 0xb2,
	0x79, 0x91, 0x22, 0x29, 0x86, 0x9c, 0x34, 0xf3, 0xfd, 0xdf, 0x5c, 0x7e, 0xc1, 0xa7, 0x56, 0xda,
	0x6e, 0x2a, 0xaf, 0x2a, 0xbd, 0x2f, 0x7a, 0xd9, 0x76, 0x76, 0x10, 0x0f, 0xc5, 0x5e, 0xf4, 0xf8,
	0xb1, 0x18, 0xca, 0x62, 0xec, 0x84, 0xc1, 0x9a, 0x0f, 0x46, 0x5b, 0x5d, 0x34, 0x93, 0xaa, 0xac,
	0xd4, 0xea, 0xca, 0xaf, 0x34, 0xf2, 0xcf, 0xfa, 0x71, 0x05, 0xaf, 0xb6, 0x4f, 0xc9, 0xed, 0x80,
	0xea, 0xb3, 0xb0, 0xe2, 0x9b, 0x17, 0xde, 0xc1, 0xc5, 0x72, 0xc2, 0xed, 0x3c, 0x20, 0x23, 0x39,
	0xd9, 0x44, 0xbb, 0x74, 0x81, 0x77, 0xf3, 0x80, 0x74, 0x0d, 0x17, 0xed, 0x24, 0x14, 0x6f, 0x26,
	0xde, 0xe3, 0x2f, 0xec, 0xd9, 0xca, 0x4b, 0x89, 0x83, 0xdb, 0xe9, 0xab, 0x43, 0xf4, 0x0d, 0x40,
	0x87, 0x46, 0x3f, 0x09, 0xc7, 0x5e, 0x88, 0x1d, 0x09, 0xf1, 0x6b, 0x88, 0xf7, 0x42, 0x2a, 0x6e,
	0xc5, 0x78, 0xcf, 0x4e, 0x7c, 0x7a, 0xee, 0xc0, 0x9d, 0x18, 0xef, 0x69, 0x06, 0x49, 0x29, 0xf8,
	0x8c, 0x7c, 0xb4, 0xa2, 0x45, 0x16, 0x85, 0xe3, 0x52, 0xfc, 0xc0, 0xef, 0x0e, 0xd0, 0xb7, 0x90,
	0x94, 0x93, 0xec, 0x6b, 0xa9, 0x5a, 0x2e, 0x6b, 0x76, 0xea, 0x73, 0x58, 0xd0, 0x97, 0xda, 0x09,
	0x56, 0x3f, 0xa0, 0xe1, 0x4d, 0xaf, 0xb5, 0x61, 0x67, 0x41, 0xf0, 0x68, 0xeb, 0x08, 0x65, 0x70,
	0x56, 0x4f, 0xaa, 0x45, 0xad, 0xd8, 0xb9, 0x0f, 0x97, 0x95, 0x52, 0x38, 0xa9, 0x71, 0xac, 0x18,
	0xe4, 0x64, 0x13, 0xef, 0xfc, 0xec, 0x98, 0xac, 0xb4, 0x62, 0x49, 0x60, 0x6e, 0xa6, 0x97, 0x90,
	0x2a, 0x6d, 0x65, 0x33, 0x73, 0x6d, 0x6a, 0x34, 0x2c, 0x0d, 0x15, 0x04, 0x76, 0xeb, 0xd0, 0xfa,
	0x03, 0xbc, 0xbc, 0x41, 0xa3, 0x97, 0xa2, 0x43, 0xc1, 0x97, 0x00, 0x7a, 0xc0, 0x50, 0xee, 0xc8,
	0x48, 0x7e, 0xbc, 0x89, 0xae, 0x57, 0x2f, 0x8e, 0x76, 0xb1, 0xa3, 0xae, 0xdd, 0xf1, 0x3a, 0xff,
	0x73, 0xc8, 0xc8, 0xdf, 0x43, 0x46, 0xfe, 0x1d, 0x32, 0xf2, 0xfb, 0x31, 0x3b, 0xba, 0x21, 0x3f,
	0xd3, 0xe7, 0xbf, 0x5a, 0x9e, 0xfa, 0xe7, 0xfd, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf0, 0xaf,
	0x5b, 0x61, 0x09, 0x02, 0x00, 0x00,
}
