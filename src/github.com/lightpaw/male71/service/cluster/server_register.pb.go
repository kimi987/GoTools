// Code generated by protoc-gen-gogo.
// source: github.com/lightpaw/male7/service/cluster/server_register.proto
// DO NOT EDIT!

/*
	Package cluster is a generated protocol buffer package.

	It is generated from these files:
		github.com/lightpaw/male7/service/cluster/server_register.proto

	It has these top-level messages:
		GameServerInfoProto
*/
package cluster

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GameServerInfoProto struct {
	// 4 byte ip address. eg. 192.168.1.1 (must be intranet address or localhost address)
	Address []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// server port
	Port uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// server id
	Id uint32 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	// server name
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// 连接服地址
	ConnAddr string `protobuf:"bytes,5,opt,name=conn_addr,json=connAddr,proto3" json:"conn_addr,omitempty"`
	// 游戏服监控地址
	MetricsAddr string `protobuf:"bytes,6,opt,name=metrics_addr,json=metricsAddr,proto3" json:"metrics_addr,omitempty"`
	// 游戏服版本号
	Version string `protobuf:"bytes,7,opt,name=version,proto3" json:"version,omitempty"`
	// 游戏服配置版本号
	ConfigVersion string `protobuf:"bytes,8,opt,name=config_version,json=configVersion,proto3" json:"config_version,omitempty"`
	// 客户端热更新版本号，如"stable1_1"
	ClientVersion string `protobuf:"bytes,9,opt,name=client_version,json=clientVersion,proto3" json:"client_version,omitempty"`
	// 游戏服的rpc监听地址
	RpcAddr string `protobuf:"bytes,11,opt,name=rpc_addr,json=rpcAddr,proto3" json:"rpc_addr,omitempty"`
	// 版本编译时间
	BuildTime int32 `protobuf:"varint,12,opt,name=build_time,json=buildTime,proto3" json:"build_time,omitempty"`
}

func (m *GameServerInfoProto) Reset()         { *m = GameServerInfoProto{} }
func (m *GameServerInfoProto) String() string { return proto.CompactTextString(m) }
func (*GameServerInfoProto) ProtoMessage()    {}
func (*GameServerInfoProto) Descriptor() ([]byte, []int) {
	return fileDescriptorServerRegister, []int{0}
}

func (m *GameServerInfoProto) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *GameServerInfoProto) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *GameServerInfoProto) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GameServerInfoProto) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GameServerInfoProto) GetConnAddr() string {
	if m != nil {
		return m.ConnAddr
	}
	return ""
}

func (m *GameServerInfoProto) GetMetricsAddr() string {
	if m != nil {
		return m.MetricsAddr
	}
	return ""
}

func (m *GameServerInfoProto) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *GameServerInfoProto) GetConfigVersion() string {
	if m != nil {
		return m.ConfigVersion
	}
	return ""
}

func (m *GameServerInfoProto) GetClientVersion() string {
	if m != nil {
		return m.ClientVersion
	}
	return ""
}

func (m *GameServerInfoProto) GetRpcAddr() string {
	if m != nil {
		return m.RpcAddr
	}
	return ""
}

func (m *GameServerInfoProto) GetBuildTime() int32 {
	if m != nil {
		return m.BuildTime
	}
	return 0
}

func init() {
	proto.RegisterType((*GameServerInfoProto)(nil), "cluster.GameServerInfoProto")
}
func (m *GameServerInfoProto) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GameServerInfoProto) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	if m.Port != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(m.Port))
	}
	if m.Id != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(m.Id))
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.ConnAddr) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.ConnAddr)))
		i += copy(dAtA[i:], m.ConnAddr)
	}
	if len(m.MetricsAddr) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.MetricsAddr)))
		i += copy(dAtA[i:], m.MetricsAddr)
	}
	if len(m.Version) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	if len(m.ConfigVersion) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.ConfigVersion)))
		i += copy(dAtA[i:], m.ConfigVersion)
	}
	if len(m.ClientVersion) > 0 {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.ClientVersion)))
		i += copy(dAtA[i:], m.ClientVersion)
	}
	if len(m.RpcAddr) > 0 {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(len(m.RpcAddr)))
		i += copy(dAtA[i:], m.RpcAddr)
	}
	if m.BuildTime != 0 {
		dAtA[i] = 0x60
		i++
		i = encodeVarintServerRegister(dAtA, i, uint64(m.BuildTime))
	}
	return i, nil
}

func encodeFixed64ServerRegister(dAtA []byte, offset int, v uint64) int {
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
func encodeFixed32ServerRegister(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintServerRegister(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *GameServerInfoProto) Size() (n int) {
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	if m.Port != 0 {
		n += 1 + sovServerRegister(uint64(m.Port))
	}
	if m.Id != 0 {
		n += 1 + sovServerRegister(uint64(m.Id))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.ConnAddr)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.MetricsAddr)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.ConfigVersion)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.ClientVersion)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	l = len(m.RpcAddr)
	if l > 0 {
		n += 1 + l + sovServerRegister(uint64(l))
	}
	if m.BuildTime != 0 {
		n += 1 + sovServerRegister(uint64(m.BuildTime))
	}
	return n
}

func sovServerRegister(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozServerRegister(x uint64) (n int) {
	return sovServerRegister(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GameServerInfoProto) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerRegister
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
			return fmt.Errorf("proto: GameServerInfoProto: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GameServerInfoProto: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Port", wireType)
			}
			m.Port = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Port |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetricsAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MetricsAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfigVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConfigVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RpcAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
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
				return ErrInvalidLengthServerRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RpcAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BuildTime", wireType)
			}
			m.BuildTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerRegister
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BuildTime |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipServerRegister(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerRegister
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
func skipServerRegister(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowServerRegister
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
					return 0, ErrIntOverflowServerRegister
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
					return 0, ErrIntOverflowServerRegister
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
				return 0, ErrInvalidLengthServerRegister
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowServerRegister
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
				next, err := skipServerRegister(dAtA[start:])
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
	ErrInvalidLengthServerRegister = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowServerRegister   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/lightpaw/male7/service/cluster/server_register.proto", fileDescriptorServerRegister)
}

var fileDescriptorServerRegister = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xcf, 0x4e, 0xe3, 0x30,
	0x10, 0xc6, 0xd7, 0xd9, 0xb6, 0x69, 0xdc, 0x3f, 0x5a, 0x79, 0x2f, 0x5e, 0xad, 0x88, 0x02, 0x12,
	0x52, 0x4e, 0xcd, 0x81, 0x03, 0x47, 0x04, 0x17, 0xc4, 0x0d, 0x05, 0xc4, 0x35, 0x4a, 0xed, 0x69,
	0x6a, 0x29, 0xb6, 0x23, 0xc7, 0x2d, 0xaf, 0xc2, 0x23, 0xc1, 0x8d, 0x47, 0x40, 0xe5, 0x45, 0x50,
	0x26, 0x0d, 0xe2, 0x36, 0xdf, 0xef, 0xfb, 0xc9, 0x1e, 0x0d, 0xbd, 0xaa, 0x94, 0xdf, 0xee, 0xd6,
	0x2b, 0x61, 0x75, 0x56, 0xab, 0x6a, 0xeb, 0x9b, 0xf2, 0x39, 0xd3, 0x65, 0x0d, 0x97, 0x59, 0x0b,
	0x6e, 0xaf, 0x04, 0x64, 0xa2, 0xde, 0xb5, 0x1e, 0x1c, 0x66, 0x70, 0x85, 0x83, 0x4a, 0x75, 0x79,
	0xd5, 0x38, 0xeb, 0x2d, 0x0b, 0x8f, 0xf5, 0xd9, 0x5b, 0x40, 0xff, 0xde, 0x96, 0x1a, 0x1e, 0x50,
	0xbb, 0x33, 0x1b, 0x7b, 0x8f, 0x02, 0xa7, 0x61, 0x29, 0xa5, 0x83, 0xb6, 0xe5, 0x24, 0x21, 0xe9,
	0x3c, 0x1f, 0x22, 0x63, 0x74, 0xd4, 0x58, 0xe7, 0x79, 0x90, 0x90, 0x74, 0x91, 0xe3, 0xcc, 0x96,
	0x34, 0x50, 0x92, 0xff, 0x46, 0x12, 0x28, 0xd9, 0x39, 0xa6, 0xd4, 0xc0, 0x47, 0x09, 0x49, 0xa3,
	0x1c, 0x67, 0xf6, 0x9f, 0x46, 0xc2, 0x1a, 0x53, 0x74, 0xef, 0xf0, 0x31, 0x16, 0xd3, 0x0e, 0x5c,
	0x4b, 0xe9, 0xd8, 0x29, 0x9d, 0x6b, 0xf0, 0x4e, 0x89, 0xb6, 0xef, 0x27, 0xd8, 0xcf, 0x8e, 0x0c,
	0x15, 0x4e, 0xc3, 0x3d, 0xb8, 0x56, 0x59, 0xc3, 0x43, 0x6c, 0x87, 0xc8, 0xce, 0xe9, 0x52, 0x58,
	0xb3, 0x51, 0x55, 0x31, 0x08, 0x53, 0x14, 0x16, 0x3d, 0x7d, 0xfa, 0xa1, 0xd5, 0x0a, 0x8c, 0xff,
	0xd6, 0xa2, 0xa3, 0x86, 0x74, 0xd0, 0xfe, 0xd1, 0xa9, 0x6b, 0x44, 0xbf, 0xc6, 0xac, 0xff, 0xc8,
	0x35, 0x02, 0x57, 0x38, 0xa1, 0x74, 0xbd, 0x53, 0xb5, 0x2c, 0xbc, 0xd2, 0xc0, 0xe7, 0x09, 0x49,
	0xc7, 0x79, 0x84, 0xe4, 0x51, 0x69, 0xb8, 0xf9, 0xf3, 0x7a, 0x88, 0xc9, 0xfb, 0x21, 0x26, 0x1f,
	0x87, 0x98, 0xbc, 0x7c, 0xc6, 0xbf, 0xd6, 0x13, 0xbc, 0xf6, 0xc5, 0x57, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x95, 0xa0, 0x19, 0x5e, 0xb0, 0x01, 0x00, 0x00,
}
