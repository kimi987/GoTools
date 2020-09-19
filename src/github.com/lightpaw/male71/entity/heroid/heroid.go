package heroid

import "math"

// 英雄id = (sid << 32) | accountId

func NewHeroId(accountId int64, sid uint32) int64 {
	return int64(uint64(sid)<<32 | uint64(accountId))
}

func GetAccountId(heroId int64) int64 {
	return int64(uint64(heroId) & math.MaxUint32)
}

func GetSid(heroId int64) uint32 {
	return uint32(uint64(heroId) >> 32)
}
