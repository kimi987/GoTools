package tssclient

import (
	"github.com/lightpaw/male7/pb/rpcpb/game2tss"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/service/tss"
	"github.com/lightpaw/male7/gen/iface"
	"fmt"
)

func NewTssClient(serverConfig *kv.IndividualServerConfig, gs iface.GameServer) *TssClient {

	var client *rpc7.Client
	if len(serverConfig.TssClusterAddr) > 0 {
		c, err := rpc7.NewClient(serverConfig.TssClusterAddr)
		if err != nil {
			logrus.WithError(err).Panic("初始化TssSdk RPC失败")
		}
		client = c
	}

	callbackAddr := fmt.Sprintf("%s:%d", serverConfig.RpcAddr, gs.GetRpcPort())

	tss := &TssClient{
		client:       client,
		callbackMap:  make(map[tss.MsgCategory]tss.Callback),
		callbackAddr: callbackAddr,
	}

	rpc7.Handle(game2tss.NewUicChatCallbackHandler(tss.callback))

	return tss
}

//gogen:iface
type TssClient struct {
	client       *rpc7.Client
	callbackMap  map[tss.MsgCategory]tss.Callback
	callbackAddr string
}

func (c *TssClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *TssClient) RegisterCallback(t tss.MsgCategory, callback tss.Callback) {
	c.callbackMap[t] = callback
}

func (c *TssClient) CallbackAddr() string {
	return c.callbackAddr
}

func (c *TssClient) Client() *rpc7.Client {
	return c.client
}

func (c *TssClient) IsEnable() bool {
	return c.client != nil
}

var callbackSuccess = &game2tss.S2CUicChatCallbackProto{Success: true}

func (c *TssClient) callback(r *game2tss.C2SUicChatCallbackProto) (*game2tss.S2CUicChatCallbackProto, error) {

	if len(r.CallbackData) <= 0 {
		logrus.Error("Tss收到敏感词回调，没有callback数据")
		return callbackSuccess, nil
	}

	mc := tss.MsgCategory(r.CallbackData[0])
	callback := c.callbackMap[mc]
	if callback == nil {
		logrus.WithField("msg_category", mc).Error("Tss收到敏感词回调，未知的Callback类型")
		return callbackSuccess, nil
	}

	// 调用callback
	callback(r.RoldId, r.MsgResultFlag, r.ReplaceMsg, r.CallbackData[1:])
	return callbackSuccess, nil
}

var checkNameSuccess = &game2tss.S2CUicJudgeUserInputNameV2Proto{}

func (c *TssClient) CheckName(name string, ifReplace bool) (*game2tss.S2CUicJudgeUserInputNameV2Proto, error) {
	if c.client == nil {
		return checkNameSuccess, nil
	}

	// 需要走这个
	var resp *game2tss.S2CUicJudgeUserInputNameV2Proto
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		r, err := game2tss.UicJudgeUserInputNameV2(c.client, ctx, name, false)
		if err != nil {
			return errors.Wrapf(err, "tss 敏感词查询失败")
		}
		resp = r
		return nil
	}); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *TssClient) TryCheckName(funcName string, sender sender.Sender, name string, sensitiveWordsMsg, serverErrorMsg pbutil.Buffer) bool {

	resp, err := c.CheckName(name, false)
	if err != nil {
		logrus.WithError(err).Errorf("%s，查询失败", funcName)
		sender.Send(serverErrorMsg)
		return false
	}

	if resp.Ret != 0 {
		logrus.WithField("ret", resp.Ret).WithField("msg", resp.RetMsg).Error("%s，敏感词查询失败", funcName)
		sender.Send(serverErrorMsg)
		return false
	}

	// 0: 合法； 1：不合法，不能显示；  2：合法，但包含敏感词
	if resp.MsgResultFlag != 0 {
		logrus.WithField("ret", resp.MsgResultFlag).Error("%s，内容中包含敏感词", funcName)
		sender.Send(sensitiveWordsMsg)
		return false
	}

	return true
}

var errClientNotEnable = errors.Errorf("TssClient not enable")

func (c *TssClient) JudgeChat(proto *game2tss.C2SUicJudgeUserInputChatV2Proto) (resp *game2tss.S2CUicJudgeUserInputChatV2Proto, err error) {
	if c.client == nil {
		return nil, errClientNotEnable
	}

	b := make([]byte, 1+len(proto.CallbackData))
	b[0] = uint8(proto.MsgCategory)
	copy(b[1:], proto.CallbackData)
	proto.CallbackData = b

	if len(proto.CallbackAddr) <= 0 {
		proto.CallbackAddr = c.callbackAddr
	}

	ctxfunc.Timeout3s(func(ctx context.Context) (error) {
		resp, err = game2tss.UicJudgeUserInputChatV2Proto(c.client, ctx, proto)
		return nil
	})
	return
}
