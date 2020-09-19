package event

import (
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/pbutil"
	"time"
)

func NewMsgDispatcher(getSender func(id int64) sender.Sender, n, sizePerQuene uint64, name string) *MsgDispatcher {

	s := &MsgDispatcher{}
	s.getSender = getSender
	s.hashQueue = NewHashFuncQueue(n, sizePerQuene, name)

	return s
}

type MsgDispatcher struct {
	getSender func(id int64) sender.Sender

	hashQueue *HashFuncQueue
}

func (s MsgDispatcher) Close() {
	s.hashQueue.Close()
}

func (s MsgDispatcher) TrySend(id int64, toSend pbutil.Buffer) bool {
	if toSend == nil {
		return true
	}

	sender := s.getSender(id)
	if sender == nil {
		return true
	}

	f := func() {
		sender.Send(toSend)
	}

	return s.hashQueue.TryFunc(id, f)
}

func (s MsgDispatcher) TimeoutSend(id int64, toSend pbutil.Buffer, timeout time.Duration) bool {
	if toSend == nil {
		return true
	}

	sender := s.getSender(id)
	if sender == nil {
		return true
	}

	f := func() {
		sender.Send(toSend)
	}

	return s.hashQueue.TimeoutFunc(id, f, timeout)
}

func (s MsgDispatcher) TrySendFunc(id int64, toSendFunc func() pbutil.Buffer) bool {
	if toSendFunc == nil {
		return true
	}

	sender := s.getSender(id)
	if sender == nil {
		return true
	}

	f := func() {
		toSend := toSendFunc()
		if toSend != nil {
			sender.Send(toSend)
		}
	}

	return s.hashQueue.TryFunc(id, f)
}

func (s MsgDispatcher) TimeoutSendFunc(id int64, toSendFunc func() pbutil.Buffer, timeout time.Duration) bool {
	if toSendFunc == nil {
		return true
	}

	sender := s.getSender(id)
	if sender == nil {
		return true
	}

	f := func() {
		toSend := toSendFunc()
		if toSend != nil {
			sender.Send(toSend)
		}
	}

	return s.hashQueue.TimeoutFunc(id, f, timeout)
}

func (s MsgDispatcher) TrySendWithSender(sender sender.Sender, toSend pbutil.Buffer) bool {
	if toSend == nil {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		sender.Send(toSend)
	}

	return s.hashQueue.TryFunc(sender.Id(), f)
}

func (s MsgDispatcher) TimeoutSendWithSender(sender sender.Sender, toSend pbutil.Buffer, timeout time.Duration) bool {
	if toSend == nil {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		sender.Send(toSend)
	}

	return s.hashQueue.TimeoutFunc(sender.Id(), f, timeout)
}

func (s MsgDispatcher) TrySendFuncWithSender(sender sender.Sender, toSendFunc func() pbutil.Buffer) bool {
	if toSendFunc == nil {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		toSend := toSendFunc()
		if toSend != nil {
			sender.Send(toSend)
		}
	}

	return s.hashQueue.TryFunc(sender.Id(), f)
}

func (s MsgDispatcher) TimeoutSendFuncWithSender(sender sender.Sender, toSendFunc func() pbutil.Buffer, timeout time.Duration) bool {
	if toSendFunc == nil {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		toSend := toSendFunc()
		if toSend != nil {
			sender.Send(toSend)
		}
	}

	return s.hashQueue.TimeoutFunc(sender.Id(), f, timeout)
}

func (s MsgDispatcher) TrySendArrayWithSender(sender sender.Sender, toSend []pbutil.Buffer) bool {
	if len(toSend) <= 0 {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		sender.SendAll(toSend)
	}

	return s.hashQueue.TryFunc(sender.Id(), f)
}

func (s MsgDispatcher) TimeoutSendArrayWithSender(sender sender.Sender, toSend []pbutil.Buffer, timeout time.Duration) bool {
	if len(toSend) <= 0 {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		sender.SendAll(toSend)
	}

	return s.hashQueue.TimeoutFunc(sender.Id(), f, timeout)
}

func (s MsgDispatcher) TrySendFuncArrayWithSender(sender sender.Sender, toSendFunc []func() pbutil.Buffer) bool {
	if len(toSendFunc) <= 0 {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		for _, f := range toSendFunc {
			if f != nil {
				if toSend := f(); toSend != nil {
					sender.Send(toSend)
				}
			}
		}
	}

	return s.hashQueue.TryFunc(sender.Id(), f)
}

func (s MsgDispatcher) TimeoutSendFuncArrayWithSender(sender sender.Sender, toSendFunc []func() pbutil.Buffer, timeout time.Duration) bool {
	if len(toSendFunc) <= 0 {
		return true
	}

	if sender == nil {
		return true
	}

	f := func() {
		for _, f := range toSendFunc {
			if f != nil {
				if toSend := f(); toSend != nil {
					sender.Send(toSend)
				}
			}
		}
	}

	return s.hashQueue.TimeoutFunc(sender.Id(), f, timeout)
}
