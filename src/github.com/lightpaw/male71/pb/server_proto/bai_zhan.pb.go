// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/pb/server_proto/bai_zhan.proto
// DO NOT EDIT!

package server_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto5 "github.com/lightpaw/male7/pb/shared_proto"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type BaiZhanObjServerProto struct {
	Id                          int64                     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	JunXianLevel                uint64                    `protobuf:"varint,2,opt,name=jun_xian_level,json=junXianLevel,proto3" json:"jun_xian_level,omitempty"`
	LastJunXianLevel            uint64                    `protobuf:"varint,3,opt,name=last_jun_xian_level,json=lastJunXianLevel,proto3" json:"last_jun_xian_level,omitempty"`
	IsJunXianBeenRemoved        bool                      `protobuf:"varint,4,opt,name=is_jun_xian_been_removed,json=isJunXianBeenRemoved,proto3" json:"is_jun_xian_been_removed,omitempty"`
	ChallengeTimes              uint64                    `protobuf:"varint,5,opt,name=challenge_times,json=challengeTimes,proto3" json:"challenge_times,omitempty"`
	Point                       uint64                    `protobuf:"varint,6,opt,name=point,proto3" json:"point,omitempty"`
	LastPointChangeTime         int64                     `protobuf:"varint,7,opt,name=last_point_change_time,json=lastPointChangeTime,proto3" json:"last_point_change_time,omitempty"`
	IsCollectSalary             bool                      `protobuf:"varint,8,opt,name=is_collect_salary,json=isCollectSalary,proto3" json:"is_collect_salary,omitempty"`
	LastCollectedJunXianPrizeId uint64                    `protobuf:"varint,9,opt,name=last_collected_jun_xian_prize_id,json=lastCollectedJunXianPrizeId,proto3" json:"last_collected_jun_xian_prize_id,omitempty"`
	Mirror                      *proto5.CombatPlayerProto `protobuf:"bytes,10,opt,name=mirror" json:"mirror,omitempty"`
	MirrorFightAmount           uint64                    `protobuf:"varint,13,opt,name=mirror_fight_amount,json=mirrorFightAmount,proto3" json:"mirror_fight_amount,omitempty"`
	HistoryMaxJunXianLevel      uint64                    `protobuf:"varint,11,opt,name=history_max_jun_xian_level,json=historyMaxJunXianLevel,proto3" json:"history_max_jun_xian_level,omitempty"`
	HistoryMaxPoints            []uint64                  `protobuf:"varint,12,rep,packed,name=history_max_points,json=historyMaxPoints" json:"history_max_points,omitempty"`
}

func (m *BaiZhanObjServerProto) Reset()                    { *m = BaiZhanObjServerProto{} }
func (m *BaiZhanObjServerProto) String() string            { return proto.CompactTextString(m) }
func (*BaiZhanObjServerProto) ProtoMessage()               {}
func (*BaiZhanObjServerProto) Descriptor() ([]byte, []int) { return fileDescriptorBaiZhan, []int{0} }

func (m *BaiZhanObjServerProto) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetJunXianLevel() uint64 {
	if m != nil {
		return m.JunXianLevel
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetLastJunXianLevel() uint64 {
	if m != nil {
		return m.LastJunXianLevel
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetIsJunXianBeenRemoved() bool {
	if m != nil {
		return m.IsJunXianBeenRemoved
	}
	return false
}

func (m *BaiZhanObjServerProto) GetChallengeTimes() uint64 {
	if m != nil {
		return m.ChallengeTimes
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetPoint() uint64 {
	if m != nil {
		return m.Point
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetLastPointChangeTime() int64 {
	if m != nil {
		return m.LastPointChangeTime
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetIsCollectSalary() bool {
	if m != nil {
		return m.IsCollectSalary
	}
	return false
}

func (m *BaiZhanObjServerProto) GetLastCollectedJunXianPrizeId() uint64 {
	if m != nil {
		return m.LastCollectedJunXianPrizeId
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetMirror() *proto5.CombatPlayerProto {
	if m != nil {
		return m.Mirror
	}
	return nil
}

func (m *BaiZhanObjServerProto) GetMirrorFightAmount() uint64 {
	if m != nil {
		return m.MirrorFightAmount
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetHistoryMaxJunXianLevel() uint64 {
	if m != nil {
		return m.HistoryMaxJunXianLevel
	}
	return 0
}

func (m *BaiZhanObjServerProto) GetHistoryMaxPoints() []uint64 {
	if m != nil {
		return m.HistoryMaxPoints
	}
	return nil
}

type BaiZhanServerProto struct {
	LastResetTime int64                    `protobuf:"varint,1,opt,name=last_reset_time,json=lastResetTime,proto3" json:"last_reset_time,omitempty"`
	BaiZhanObjs   []*BaiZhanObjServerProto `protobuf:"bytes,2,rep,name=bai_zhan_objs,json=baiZhanObjs" json:"bai_zhan_objs,omitempty"`
}

func (m *BaiZhanServerProto) Reset()                    { *m = BaiZhanServerProto{} }
func (m *BaiZhanServerProto) String() string            { return proto.CompactTextString(m) }
func (*BaiZhanServerProto) ProtoMessage()               {}
func (*BaiZhanServerProto) Descriptor() ([]byte, []int) { return fileDescriptorBaiZhan, []int{1} }

func (m *BaiZhanServerProto) GetLastResetTime() int64 {
	if m != nil {
		return m.LastResetTime
	}
	return 0
}

func (m *BaiZhanServerProto) GetBaiZhanObjs() []*BaiZhanObjServerProto {
	if m != nil {
		return m.BaiZhanObjs
	}
	return nil
}

func init() {
	proto.RegisterType((*BaiZhanObjServerProto)(nil), "proto.BaiZhanObjServerProto")
	proto.RegisterType((*BaiZhanServerProto)(nil), "proto.BaiZhanServerProto")
}
func (m *BaiZhanObjServerProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BaiZhanObjServerProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.Id))
	}
	if m.JunXianLevel != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.JunXianLevel))
	}
	if m.LastJunXianLevel != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.LastJunXianLevel))
	}
	if m.IsJunXianBeenRemoved {
		dAtA[i] = 0x20
		i++
		if m.IsJunXianBeenRemoved {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.ChallengeTimes != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.ChallengeTimes))
	}
	if m.Point != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.Point))
	}
	if m.LastPointChangeTime != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.LastPointChangeTime))
	}
	if m.IsCollectSalary {
		dAtA[i] = 0x40
		i++
		if m.IsCollectSalary {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.LastCollectedJunXianPrizeId != 0 {
		dAtA[i] = 0x48
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.LastCollectedJunXianPrizeId))
	}
	if m.Mirror != nil {
		dAtA[i] = 0x52
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.Mirror.Size()))
		n1, err := m.Mirror.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.HistoryMaxJunXianLevel != 0 {
		dAtA[i] = 0x58
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.HistoryMaxJunXianLevel))
	}
	if len(m.HistoryMaxPoints) > 0 {
		dAtA3 := make([]byte, len(m.HistoryMaxPoints)*10)
		var j2 int
		for _, num := range m.HistoryMaxPoints {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		dAtA[i] = 0x62
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(j2))
		i += copy(dAtA[i:], dAtA3[:j2])
	}
	if m.MirrorFightAmount != 0 {
		dAtA[i] = 0x68
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.MirrorFightAmount))
	}
	return i, nil
}

func (m *BaiZhanServerProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BaiZhanServerProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.LastResetTime != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintBaiZhan(dAtA, i, uint64(m.LastResetTime))
	}
	if len(m.BaiZhanObjs) > 0 {
		for _, msg := range m.BaiZhanObjs {
			dAtA[i] = 0x12
			i++
			i = encodeVarintBaiZhan(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeFixed64BaiZhan(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32BaiZhan(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintBaiZhan(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *BaiZhanObjServerProto) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovBaiZhan(uint64(m.Id))
	}
	if m.JunXianLevel != 0 {
		n += 1 + sovBaiZhan(uint64(m.JunXianLevel))
	}
	if m.LastJunXianLevel != 0 {
		n += 1 + sovBaiZhan(uint64(m.LastJunXianLevel))
	}
	if m.IsJunXianBeenRemoved {
		n += 2
	}
	if m.ChallengeTimes != 0 {
		n += 1 + sovBaiZhan(uint64(m.ChallengeTimes))
	}
	if m.Point != 0 {
		n += 1 + sovBaiZhan(uint64(m.Point))
	}
	if m.LastPointChangeTime != 0 {
		n += 1 + sovBaiZhan(uint64(m.LastPointChangeTime))
	}
	if m.IsCollectSalary {
		n += 2
	}
	if m.LastCollectedJunXianPrizeId != 0 {
		n += 1 + sovBaiZhan(uint64(m.LastCollectedJunXianPrizeId))
	}
	if m.Mirror != nil {
		l = m.Mirror.Size()
		n += 1 + l + sovBaiZhan(uint64(l))
	}
	if m.HistoryMaxJunXianLevel != 0 {
		n += 1 + sovBaiZhan(uint64(m.HistoryMaxJunXianLevel))
	}
	if len(m.HistoryMaxPoints) > 0 {
		l = 0
		for _, e := range m.HistoryMaxPoints {
			l += sovBaiZhan(uint64(e))
		}
		n += 1 + sovBaiZhan(uint64(l)) + l
	}
	if m.MirrorFightAmount != 0 {
		n += 1 + sovBaiZhan(uint64(m.MirrorFightAmount))
	}
	return n
}

func (m *BaiZhanServerProto) Size() (n int) {
	var l int
	_ = l
	if m.LastResetTime != 0 {
		n += 1 + sovBaiZhan(uint64(m.LastResetTime))
	}
	if len(m.BaiZhanObjs) > 0 {
		for _, e := range m.BaiZhanObjs {
			l = e.Size()
			n += 1 + l + sovBaiZhan(uint64(l))
		}
	}
	return n
}

func sovBaiZhan(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozBaiZhan(x uint64) (n int) {
	return sovBaiZhan(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BaiZhanObjServerProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBaiZhan
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
			return fmt.Errorf("proto: BaiZhanObjServerProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BaiZhanObjServerProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JunXianLevel", wireType)
			}
			m.JunXianLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JunXianLevel |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastJunXianLevel", wireType)
			}
			m.LastJunXianLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastJunXianLevel |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsJunXianBeenRemoved", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsJunXianBeenRemoved = bool(v != 0)
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChallengeTimes", wireType)
			}
			m.ChallengeTimes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChallengeTimes |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Point", wireType)
			}
			m.Point = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Point |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastPointChangeTime", wireType)
			}
			m.LastPointChangeTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastPointChangeTime |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsCollectSalary", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsCollectSalary = bool(v != 0)
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastCollectedJunXianPrizeId", wireType)
			}
			m.LastCollectedJunXianPrizeId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastCollectedJunXianPrizeId |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mirror", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
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
				return ErrInvalidLengthBaiZhan
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Mirror == nil {
				m.Mirror = &proto5.CombatPlayerProto{}
			}
			if err := m.Mirror.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HistoryMaxJunXianLevel", wireType)
			}
			m.HistoryMaxJunXianLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HistoryMaxJunXianLevel |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBaiZhan
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
				m.HistoryMaxPoints = append(m.HistoryMaxPoints, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowBaiZhan
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
					return ErrInvalidLengthBaiZhan
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowBaiZhan
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
					m.HistoryMaxPoints = append(m.HistoryMaxPoints, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field HistoryMaxPoints", wireType)
			}
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MirrorFightAmount", wireType)
			}
			m.MirrorFightAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MirrorFightAmount |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBaiZhan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBaiZhan
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
func (m *BaiZhanServerProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBaiZhan
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
			return fmt.Errorf("proto: BaiZhanServerProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BaiZhanServerProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastResetTime", wireType)
			}
			m.LastResetTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastResetTime |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaiZhanObjs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBaiZhan
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
				return ErrInvalidLengthBaiZhan
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaiZhanObjs = append(m.BaiZhanObjs, &BaiZhanObjServerProto{})
			if err := m.BaiZhanObjs[len(m.BaiZhanObjs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBaiZhan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBaiZhan
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
func skipBaiZhan(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBaiZhan
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
					return 0, ErrIntOverflowBaiZhan
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
					return 0, ErrIntOverflowBaiZhan
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
				return 0, ErrInvalidLengthBaiZhan
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowBaiZhan
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
				next, err := skipBaiZhan(dAtA[start:])
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
	ErrInvalidLengthBaiZhan = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBaiZhan   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/pb/server_proto/bai_zhan.proto", fileDescriptorBaiZhan)
}

var fileDescriptorBaiZhan = []byte{
	// 536 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xdd, 0x6a, 0x13, 0x41,
	0x18, 0x75, 0xf3, 0x67, 0x9d, 0xfc, 0xd9, 0x69, 0x2d, 0x43, 0x95, 0xb0, 0x14, 0xd1, 0x20, 0x9a,
	0x48, 0x0b, 0x55, 0xbc, 0xd2, 0x04, 0x45, 0x45, 0x31, 0x6c, 0xbd, 0x90, 0xde, 0x0c, 0xb3, 0xbb,
	0x63, 0x76, 0xc2, 0xee, 0x4c, 0x98, 0x99, 0xc4, 0xa4, 0x17, 0x3e, 0x87, 0x8f, 0xe4, 0xa5, 0x8f,
	0x20, 0xf1, 0xc2, 0xd7, 0x90, 0xf9, 0x76, 0xd3, 0xc4, 0x22, 0x5e, 0x2d, 0x73, 0xce, 0xf9, 0x7e,
	0xcf, 0xb7, 0xe8, 0xe9, 0x58, 0xd8, 0x64, 0x16, 0xf6, 0x22, 0x95, 0xf5, 0x53, 0x31, 0x4e, 0xec,
	0x94, 0x7d, 0xe9, 0x67, 0x2c, 0xe5, 0x4f, 0xfa, 0xd3, 0xb0, 0x6f, 0xb8, 0x9e, 0x73, 0x4d, 0xa7,
	0x5a, 0x59, 0xd5, 0x0f, 0x99, 0xa0, 0x17, 0x09, 0x93, 0x3d, 0x78, 0xe2, 0x2a, 0x7c, 0x0e, 0x4f,
	0xff, 0x9f, 0x20, 0x61, 0x9a, 0xc7, 0x45, 0x82, 0x48, 0x65, 0x21, 0xb3, 0x79, 0xf8, 0xd1, 0xef,
	0x0a, 0xba, 0x35, 0x60, 0xe2, 0x3c, 0x61, 0xf2, 0x43, 0x38, 0x39, 0x83, 0x42, 0x23, 0x48, 0xdc,
	0x42, 0x25, 0x11, 0x13, 0xcf, 0xf7, 0xba, 0xe5, 0xa0, 0x24, 0x62, 0x7c, 0x17, 0xb5, 0x26, 0x33,
	0x49, 0x17, 0x82, 0x49, 0x9a, 0xf2, 0x39, 0x4f, 0x49, 0xc9, 0xf7, 0xba, 0x95, 0xa0, 0x31, 0x99,
	0xc9, 0x4f, 0x82, 0xc9, 0x77, 0x0e, 0xc3, 0x8f, 0xd0, 0x5e, 0xca, 0x8c, 0xa5, 0x57, 0xa4, 0x65,
	0x90, 0xde, 0x74, 0xd4, 0xdb, 0x6d, 0xf9, 0x29, 0x22, 0xc2, 0x6c, 0xc4, 0x21, 0xe7, 0x92, 0x6a,
	0x9e, 0xa9, 0x39, 0x8f, 0x49, 0xc5, 0xf7, 0xba, 0x3b, 0xc1, 0xbe, 0x30, 0x45, 0xc4, 0x80, 0x73,
	0x19, 0xe4, 0x1c, 0xbe, 0x8f, 0xda, 0x51, 0xc2, 0xd2, 0x94, 0xcb, 0x31, 0xa7, 0x56, 0x64, 0xdc,
	0x90, 0x2a, 0x94, 0x68, 0x5d, 0xc2, 0x1f, 0x1d, 0x8a, 0xf7, 0x51, 0x75, 0xaa, 0x84, 0xb4, 0xa4,
	0x06, 0x74, 0xfe, 0xc0, 0x27, 0xe8, 0x00, 0xba, 0x84, 0x17, 0x8d, 0x12, 0xb6, 0x4e, 0x43, 0xae,
	0xc3, 0xbc, 0x30, 0xc3, 0xc8, 0x91, 0x43, 0xe0, 0x5c, 0x2e, 0xfc, 0x00, 0xed, 0x0a, 0x43, 0x23,
	0x95, 0xa6, 0x3c, 0xb2, 0xd4, 0xb0, 0x94, 0xe9, 0x25, 0xd9, 0x81, 0x26, 0xdb, 0xc2, 0x0c, 0x73,
	0xfc, 0x0c, 0x60, 0xfc, 0x12, 0xf9, 0x50, 0xa0, 0x50, 0xf3, 0x78, 0x33, 0xe3, 0x54, 0x8b, 0x0b,
	0x4e, 0x45, 0x4c, 0x6e, 0x40, 0x47, 0xb7, 0x9d, 0x6e, 0xb8, 0x96, 0x15, 0xa3, 0x8e, 0x9c, 0xe6,
	0x4d, 0x8c, 0x1f, 0xa3, 0x5a, 0x26, 0xb4, 0x56, 0x9a, 0x20, 0xdf, 0xeb, 0xd6, 0x8f, 0x49, 0xee,
	0x5a, 0x6f, 0x08, 0x16, 0x8e, 0x52, 0xb6, 0x2c, 0xdc, 0x0a, 0x0a, 0x1d, 0x7e, 0x86, 0x0e, 0x13,
	0x61, 0xac, 0xd2, 0x4b, 0x9a, 0xb1, 0xc5, 0x55, 0x1b, 0xea, 0x50, 0xf2, 0xa0, 0x50, 0xbc, 0x67,
	0x8b, 0xbf, 0xcc, 0x78, 0x88, 0xf0, 0x76, 0x2c, 0x2c, 0xc7, 0x90, 0x86, 0x5f, 0x76, 0xd6, 0x6d,
	0x62, 0x60, 0x2f, 0x06, 0xf7, 0xd0, 0x5e, 0x5e, 0x93, 0x7e, 0x76, 0xf7, 0x46, 0x59, 0xa6, 0x66,
	0xd2, 0x92, 0x26, 0x94, 0xd8, 0xcd, 0xa9, 0x57, 0x8e, 0x79, 0x01, 0xc4, 0xd1, 0x57, 0x84, 0x8b,
	0x43, 0xdb, 0xbe, 0xb2, 0x7b, 0xa8, 0x0d, 0x8b, 0xd2, 0xdc, 0x70, 0x9b, 0x5b, 0x90, 0x9f, 0x5c,
	0xd3, 0xc1, 0x81, 0x43, 0x61, 0xf9, 0xcf, 0x51, 0x73, 0x7d, 0xf8, 0x54, 0x85, 0x13, 0x43, 0x4a,
	0x7e, 0xb9, 0x5b, 0x3f, 0xbe, 0x53, 0x2c, 0xe4, 0x9f, 0x27, 0x1c, 0xd4, 0xc3, 0x4b, 0xd8, 0x0c,
	0xfc, 0xef, 0xab, 0x8e, 0xf7, 0x63, 0xd5, 0xf1, 0x7e, 0xae, 0x3a, 0xde, 0xb7, 0x5f, 0x9d, 0x6b,
	0xaf, 0xbd, 0xf3, 0xc6, 0xf6, 0x9f, 0x15, 0xd6, 0xe0, 0x73, 0xf2, 0x27, 0x00, 0x00, 0xff, 0xff,
	0x53, 0xf5, 0x0c, 0x5e, 0x8d, 0x03, 0x00, 0x00,
}
