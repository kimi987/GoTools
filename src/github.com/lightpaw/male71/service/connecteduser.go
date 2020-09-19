package service

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity/heroid"
)

func NewConnectedUser(uid int64, sender sender.ClosableSender, tencentInfoProto *shared_proto.TencentInfoProto) *ConnectedUser {
	c := &ConnectedUser{sender: sender, uid: uid, misc: &server_proto.UserMiscProto{}, tencentInfoProto: tencentInfoProto}
	return c
}

// 已在线的用户
//gogen:iface entity
type ConnectedUser struct {
	// 用户的连接
	sender sender.ClosableSender

	// 用户id
	uid int64

	hc iface.HeroController

	loaded bool // 需要atomic吗?

	misc            *server_proto.UserMiscProto // 玩家杂项
	needOfflineSave bool                        // 杂项是否需要离线保存

	tencentInfoProto *shared_proto.TencentInfoProto

	logoutType uint64
}

func (cu *ConnectedUser) SetLogoutType(t uint64) {
	cu.logoutType = t
}

func (cu *ConnectedUser) LogoutType() uint64 {
	return cu.logoutType
}

func (cu *ConnectedUser) TencentInfo() *shared_proto.TencentInfoProto {
	return cu.tencentInfoProto
}

func (cu *ConnectedUser) IsLoaded() bool {
	return cu.loaded
}

func (cu *ConnectedUser) SetLoaded() {
	cu.loaded = true
}

func (cu *ConnectedUser) SetHeroController(hc iface.HeroController) {
	cu.hc = hc
}

func (cu *ConnectedUser) GetHeroController() iface.HeroController {
	return cu.hc
}

func (cu *ConnectedUser) Id() int64 {
	return cu.uid
}

func (cu *ConnectedUser) Sid() uint32 {
	return heroid.GetSid(cu.uid)
}

func (cu *ConnectedUser) Disconnect(err msg.ErrMsg) {
	cu.sender.Disconnect(err)
}

func (cu *ConnectedUser) DisconnectAndWait(err msg.ErrMsg) {
	cu.sender.DisconnectAndWait(err)
}

func (cu *ConnectedUser) IsClosed() bool {
	return cu.sender.IsClosed()
}

// 发送消息.
func (cu *ConnectedUser) Send(msg pbutil.Buffer) {
	cu.sender.Send(msg)
}

// 发送在线路繁忙时可以被丢掉的消息
func (cu *ConnectedUser) SendIfFree(msg pbutil.Buffer) {
	cu.sender.SendIfFree(msg)
}

// 发送在线路繁忙时可以被丢掉的消息
func (cu *ConnectedUser) SendAll(msgs []pbutil.Buffer) {
	cu.sender.SendAll(msgs)
}

// 玩家杂项
func (cu *ConnectedUser) Misc() *server_proto.UserMiscProto {
	return cu.misc
}

// 玩家杂项
func (cu *ConnectedUser) SetMisc(toSet *server_proto.UserMiscProto) {
	cu.misc = toSet
	cu.needOfflineSave = true
}

// 玩家杂项
func (cu *ConnectedUser) MiscNeedOfflineSave() bool {
	return cu.needOfflineSave
}
