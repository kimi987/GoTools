package tss

type MsgCategory uint8

const (
	Mail MsgCategory = 1
	Chat MsgCategory = 2
)

type Callback func(heroId int64, msgResultFlag int32, replaceMsg string, callbackData []byte)
