package chat

import (
	"github.com/lightpaw/male7/util/collection"
	//"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"sync"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/util/i64"
)

func NewChatRecord(capcity, firstSendNum int) *ChatRecord {
	cr := &ChatRecord{
		RingList: collection.NewRingList(capcity),
		firstSendNum: firstSendNum,
		mux: &sync.RWMutex{},
	}

	return cr
}

type ChatRecord struct {
	*collection.RingList

	firstSendNum int
	msgId uint64 // 消息自增计数ID

	mux *sync.RWMutex
}

func (r *ChatRecord) AddChat(proto *shared_proto.ChatMsgProto)  {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.msgId++
	proto.ChatId = i64.ToBytesU64(r.msgId)
	r.Add(proto)
}

func (r *ChatRecord) GetFirstChatRecord() pbutil.Buffer {

	r.mux.RLock()
	defer r.mux.RUnlock()

	if len := r.Length(); len > 0 {
		var size int
		if len > r.firstSendNum {
			size = r.firstSendNum
		} else {
			size = len
		}
		data := make([][]byte, size)
		index := len - 1
		for i := 0; i < size; i++ {
			proto := r.Get(index).(*shared_proto.ChatMsgProto)
			data[i] = must.Marshal(proto)
			index--
		}
		return chat.NewS2cListHistoryChatMarshalMsg(data)
	}
	return nil
}

func (r *ChatRecord) GetChatRecored(minChatId int64) pbutil.Buffer {
	r.mux.RLock()
	defer r.mux.RUnlock()

	len := r.Length()
	if len > 0 {
		proto := r.Get(0).(*shared_proto.ChatMsgProto)
		chatId, _ := i64.FromBytes(proto.ChatId)
		if chatId == minChatId {
			return nil
		}
	} else {
		return nil
	}

	index := len - 1
	findIndex := -1
	for i := 0; i < len; i++ {
		proto := r.Get(index).(*shared_proto.ChatMsgProto)
		chatId, _ := i64.FromBytes(proto.ChatId)
		if minChatId == chatId {
			findIndex = index - 1
			break
		}
		index--
	}
	if findIndex != -1 {
		size := findIndex + 1
		if size > 10 { // 最多推10条过去
			size = 10
		}
		data := make([][]byte, size)
		for i := 0; i < size; i++ {
			proto := r.Get(findIndex).(*shared_proto.ChatMsgProto)
			data[i] = must.Marshal(proto)
			findIndex--
		}
		return chat.NewS2cListHistoryChatMarshalMsg(data)
	}

	return nil
}