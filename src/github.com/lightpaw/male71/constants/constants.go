package constants

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
)

const PID = 1

// db
const (
	DBBusyCallingTimes = 30
)

// region
const (
	BaseConflictRange = 4
)

// limit
const (
	RebuildMaxTroopMoveDuration    = 8 * time.Hour
	RebuildMaxRecurringRobDuration = 8 * time.Hour
)

// domestic
const (
	SeekTypeWorker = 1
	SeekTypeTech   = 2
)

// military
const (
	CaptainCountPerTroop = 5
)

// misc
const (
	Ik = 1000
	Iw = 10000

	BaowuLogPreviewCount = 4
)

// function
const (
	FunctionType_TYPE_FEN_CHENG   = 29  // 分城
	FunctionType_TYPE_FEN_CHENG_2 = 56; // 分城2
	FunctionType_TYPE_FEN_CHENG_3 = 57; // 分城3
	FunctionType_TYPE_FEN_CHENG_4 = 58; // 分城4

	FunctionType_TYPE_MULTI_LEVEL_MONSTER = 38; // 讨伐野怪
)

var (
	PlayerNamePrefix = "Player_"
)

const (
	RealmIndexBlockSize = 25
)

type ChatFunc func(proto *shared_proto.ChatMsgProto)
