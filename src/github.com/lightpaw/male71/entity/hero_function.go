package entity

import (
	"github.com/lightpaw/male7/config/function"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func newHeroFunction(functionOpenDataArray []*function.FunctionOpenData) *HeroFunction {
	hero := &HeroFunction{
		FunctionOpenDataArray: functionOpenDataArray,
	}

	return hero
}

// 玩家功能
type HeroFunction struct {
	FunctionOpenDataArray []*function.FunctionOpenData
	openTypes             []uint64 // 解锁了的功能
}

// 功能是否开启了
func (h *HeroFunction) IsFunctionOpened(funcType uint64) bool {
	arrayIndex := funcType / 64
	shiftAmount := funcType % 64

	if arrayIndex >= uint64(len(h.openTypes)) {
		return false
	}

	return h.openTypes[arrayIndex]&(1<<shiftAmount) != 0
}

// 开启功能
func (h *HeroFunction) OpenFunction(funcType uint64) {
	intType := int(funcType)
	arrayIndex := intType / 64
	shiftAmount := uint64(intType % 64)
	if arrayIndex >= len(h.openTypes) {
		openTypes := make([]uint64, arrayIndex+1)
		copy(openTypes, h.openTypes)
		h.openTypes = openTypes
	}

	h.openTypes[arrayIndex] |= 1 << shiftAmount
}

// 序列化给客户端
func (h *HeroFunction) encodeClient() *shared_proto.HeroFunctionProto {
	proto := &shared_proto.HeroFunctionProto{}

	for i := 0; i < len(h.openTypes); i++ {
		startIndex := i * 64
		for i := 0; i < 64; i++ {
			funcType := uint64(startIndex + i)
			if h.IsFunctionOpened(funcType) {
				proto.OpenTypes = append(proto.OpenTypes, u64.Int32(funcType))
			}
		}
	}

	return proto
}

// 序列化给服务器
func (h *HeroFunction) encodeServer() *server_proto.HeroFunctionServerProto {
	proto := &server_proto.HeroFunctionServerProto{}

	proto.OpenTypes = h.openTypes

	return proto
}

// 反序列化
func (h *HeroFunction) unmarshal(proto *server_proto.HeroFunctionServerProto) {
	if proto == nil {
		return
	}

	h.openTypes = proto.OpenTypes
}
