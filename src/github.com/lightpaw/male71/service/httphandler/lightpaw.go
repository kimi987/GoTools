package httphandler

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"net/http"
	"strconv"
	"github.com/lightpaw/male7/util/i64"
	"crypto/md5"
	"github.com/lightpaw/male7/config/kv"
	"fmt"
	"github.com/lightpaw/male7/service/operate_type"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/util/ctxfunc"
	"context"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/logrus"
	"io/ioutil"
	"encoding/json"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/config/goods"
)

var (
	ok = marshalResponse("")
)

func NewLightpawHandler(dep iface.ServiceDep, time iface.TimeService, heroService iface.HeroDataService, serverConfig *kv.IndividualServerConfig) *LightpawHandler {

	m := &LightpawHandler{
		serverConfig: serverConfig,
		time:         time,
		heroService:  heroService,
		wordService:  dep.World(),
		dep:          dep,
	}

	m.hctx = heromodule.NewContext(dep, operate_type.GM)

	// 战斗调试信息
	http.HandleFunc("/lightpaw/fightdebug", func(writer http.ResponseWriter, request *http.Request) {
		m.serverConfig.IsDebugFight = !m.serverConfig.IsDebugFight
		writer.Write([]byte(fmt.Sprintf("debug: %v", m.serverConfig.IsDebugFight)))
	})
	http.HandleFunc("/lightpaw/add_yuanbao", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "加元宝", m.addYuanbao)
	})
	http.HandleFunc("/lightpaw/add_yinliang", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "加银两", m.addYinliang)
	})
	http.HandleFunc("/lightpaw/add_dianquan", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "加点券", m.addDianquan)
	})
	http.HandleFunc("/lightpaw/ban_chat", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "禁言", m.banChat)
	})
	http.HandleFunc("/lightpaw/ban_chat_unlock", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "禁言解禁", m.banChatUnlock)
	})
	http.HandleFunc("/lightpaw/ban_login", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "封号", m.banLogin)
	})
	http.HandleFunc("/lightpaw/ban_login_unlock", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "封号解封", m.banLoginUnlock)
	})
	http.HandleFunc("/lightpaw/prize_mail", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroBodyFunc(writer, request, "发邮件礼品", m.sendPrizeMail)
	})
	http.HandleFunc("/lightpaw/prize_mail_all", func(writer http.ResponseWriter, request *http.Request) {
		m.doBodyFunc(writer, request, "发邮件礼品（全服）", m.sendPrizeMailAll)
	})
	http.HandleFunc("/lightpaw/broadcast", func(writer http.ResponseWriter, request *http.Request) {
		m.doBodyFunc(writer, request, "全服广播（跑马灯）", m.broadcast0)
	})
	http.HandleFunc("/lightpaw/system_chat", func(writer http.ResponseWriter, request *http.Request) {
		m.doBodyFunc(writer, request, "全服系统频道发送", m.broadcast1)
	})
	http.HandleFunc("/lightpaw/hero_info", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "查询角色基础信息", m.requestHeroInfo)
	})
	http.HandleFunc("/lightpaw/captain_info", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "查询武将基础信息(武将列表)", m.requestCaptainInfo)
	})
	http.HandleFunc("/lightpaw/captain_equip", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "查询装备信息(武将)", m.requestCaptainEquipInfo)
	})
	http.HandleFunc("/lightpaw/captain_gem", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroAmountFunc(writer, request, "查询宝石信息(武将)", m.requestCaptainGemInfo)
	})
	http.HandleFunc("/lightpaw/hero_depot_goods", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "查询角色仓库信息(物品)", m.requestHeroDepotGoods)
	})
	http.HandleFunc("/lightpaw/hero_depot_equip", func(writer http.ResponseWriter, request *http.Request) {
		m.doHeroFunc(writer, request, "查询角色仓库信息(装备)", m.requestHeroDepotEquip)
	})
	// 后面照着上面写，根据不同需求斟选doHeroFunc、doHeroAmountFunc、doBodyFunc、doHeroBodyFunc...
	// 如果还有无法满足以上规则的需求，请另外再增加新的doFunc系列
	return m
}

//gogen:iface
type LightpawHandler struct {
	serverConfig *kv.IndividualServerConfig

	dep         iface.ServiceDep
	time        iface.TimeService
	heroService iface.HeroDataService
	wordService iface.WorldService

	hctx *heromodule.HeroContext
}

func (m *LightpawHandler) checkSign(w http.ResponseWriter, r *http.Request) bool {
	timeStr := r.FormValue("time")
	if len(timeStr) <= 0 {
		w.Write(marshalResponse("empty time"))
		return false
	}

	time, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		w.Write(marshalResponse("invalid time, " + timeStr))
		return false
	}

	ctime := m.time.CurrentTime()
	if i64.Abs(ctime.Unix()-time) > 300 {
		w.Write(marshalResponse("time expired, " + timeStr))
		return false
	}

	sign := r.FormValue("sign")
	if len(timeStr) <= 0 {
		w.Write(marshalResponse("empty sign"))
		return false
	}

	if sign != m.computeHash(timeStr) {
		w.Write(marshalResponse("sign not match"))
		return false
	}

	return true
}

func (m *LightpawHandler) computeHash(time string) string {
	return computeHash(time, m.serverConfig.LightpawKey)
}

func computeHash(time string, key []byte) string {
	sum := md5.New()
	sum.Write([]byte(time))
	sum.Write(key)
	return fmt.Sprintf("%x", sum.Sum(nil))
}

// 传出核心数据(hereId)
func (m *LightpawHandler) extractHero(w http.ResponseWriter, r *http.Request) (heroId int64) {
	heroIdStr := r.FormValue("id")
	emptyHeroId := len(heroIdStr) <= 0 || heroIdStr == "0"
	if emptyHeroId {
		// 查看有没有name属性
		heroNameStr := r.FormValue("name")
		if len(heroNameStr) <= 0 {
			w.Write(marshalResponse("empty id"))
			return
		} else if err := ctxfunc.Timeout2s(func(ctx context.Context) error {
			// 根据name从数据库中查找ID
			id, e := m.dep.Db().HeroId(ctx, heroNameStr)
			heroId = id
			return e
		}); err != nil {
			w.Write(marshalResponse(err.Error()))
			return
		}
	}
	if heroId == 0 {
		if emptyHeroId {
			w.Write(marshalResponse("invalid id or name, " + heroIdStr))
			return
		}
		id, err := strconv.ParseInt(heroIdStr, 10, 64)
		if err != nil {
			w.Write(marshalResponse("invalid id, " + heroIdStr))
			return
		}
		if m.wordService.IsOnline(id) {
			heroId = id
		} else if err = ctxfunc.Timeout2s(func(ctx context.Context) error {
			// 从数据库中查找是否存在该id
			ok, e := m.dep.Db().HeroIdExist(ctx, id)
			if ok {
				heroId = id
			}
			return e
		}); err != nil {
			w.Write(marshalResponse(err.Error()))
			return
		}
		// 还是找不到，再看有没有name，走到这里绝B没有查过name
		if heroId == 0 {
			heroNameStr := r.FormValue("name")
			if len(heroNameStr) <= 0 {
				w.Write(marshalResponse("hero not exist"))
				return
			} else if err := ctxfunc.Timeout2s(func(ctx context.Context) error {
				// 根据name从数据库中查找ID
				id, e := m.dep.Db().HeroId(ctx, heroNameStr)
				heroId = id
				return e
			}); err != nil {
				w.Write(marshalResponse(err.Error()))
				return
			}
			if heroId == 0 {
				w.Write(marshalResponse("hero not exist"))
			}
		}
	}
	return
}

// 传出核心数据(body)
func extractBody(w http.ResponseWriter, r *http.Request) (body []byte) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(marshalResponse("no body"))
		return
	}
	body = b
	return
}

// 传出核心数据(amount)
func extractAmount(w http.ResponseWriter, r *http.Request) (amount uint64) {
	amountStr := r.FormValue("amount")
	if len(amountStr) <= 0 {
		w.Write(marshalResponse("empty amount"))
		return
	}
	a, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		w.Write(marshalResponse("invalid amount, " + amountStr))
		return
	}
	amount = a
	return
}

func (m *LightpawHandler) doHeroFunc(w http.ResponseWriter, r *http.Request, name string, f func(int64) []byte) {
	if !m.checkSign(w, r) {
		return
	}
	heroId := m.extractHero(w, r)
	if heroId == 0 {
		return
	}
	logrus.WithField("name", name).WithField("heroId", heroId).Debug("收到后台Gm命令")
	w.Write(f(heroId))
}

func (m *LightpawHandler) doHeroAmountFunc(w http.ResponseWriter, r *http.Request, name string, f func(int64, uint64) []byte) {
	if !m.checkSign(w, r) {
		return
	}
	heroId := m.extractHero(w, r)
	if heroId == 0 {
		return
	}
	amount := extractAmount(w, r)
	if amount == 0 {
		return
	}
	logrus.WithField("name", name).WithField("heroId", heroId).WithField("amount", amount).Debug("收到后台Gm命令")
	w.Write(f(heroId, amount))
}

func (m *LightpawHandler) doHeroBodyFunc(w http.ResponseWriter, r *http.Request, name string, f func(int64, []byte) []byte) {
	if !m.checkSign(w, r) {
		return
	}
	heroId := m.extractHero(w, r)
	if heroId == 0 {
		return
	}
	body := extractBody(w, r)
	if len(body) <= 0 {
		return
	}
	logrus.WithField("name", name).WithField("heroId", heroId).Debug("收到后台Gm命令")
	w.Write(f(heroId, body))
}

func (m *LightpawHandler) doBodyFunc(w http.ResponseWriter, r *http.Request, name string, f func([]byte) []byte) {
	if !m.checkSign(w, r) {
		return
	}
	body := extractBody(w, r)
	if len(body) <= 0 {
		return
	}
	logrus.WithField("name", name).Debug("收到后台Gm命令")
	w.Write(f(body))
}

func (m *LightpawHandler) broadcast0(body []byte) []byte {
	m.dep.Broadcast().Broadcast(string(body), false)
	return ok
}

func (m *LightpawHandler) broadcast1(body []byte) []byte {
	m.dep.Broadcast().Broadcast(string(body), true)
	return ok
}

func marshalResponse(resp interface{}) []byte {
	if resp != nil {
		str, ok := resp.(string)
		if ok {
			resp = newGMResponse(str, nil)
		} else {
			resp = newGMResponse("", resp)
		}
	} else {
		resp = newGMResponse("hero not found", nil)
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		return []byte(err.Error())
	}
	return bytes
}

func (m *LightpawHandler) requestHeroInfo(heroId int64) []byte {
	var resp *GMHeroInfoResponse
	var guildId int64
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		resp = newGMHeroInfoResponse(hero)
		guildId = hero.GuildId()
		return
	})
	if guildId != 0 {
		m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g != nil {
				resp.setGuildInfo(g)
			}
		})
	}
	return marshalResponse(resp)
}

func (m *LightpawHandler) requestCaptainInfo(heroId int64) []byte {
	var resp *GMCaptainInfoResponse
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		resp = &GMCaptainInfoResponse{}
		for _, captain := range hero.Military().Captains() {
			resp.Captains = append(resp.Captains, newGMCaptianInfo(captain))
		}
		return
	})
	return marshalResponse(resp)
}

func (m *LightpawHandler) requestCaptainEquipInfo(heroId int64, captainId uint64) []byte {
	var resp *GMCaptainEquipInfoResponse
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			return
		}
		resp = &GMCaptainEquipInfoResponse {
			CaptainId: captainId,
		}
		captain.WalkEquipment(func(e *entity.Equipment) (walkEnd bool) {
			resp.Equips = append(resp.Equips, newGMEquipInfo(e))
			return
		})
		return
	})
	return marshalResponse(resp)
}

func (m *LightpawHandler) requestCaptainGemInfo(heroId int64, captainId uint64) []byte {
	var resp *GMCaptainGemInfoResponse
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			return
		}
		resp = &GMCaptainGemInfoResponse {
			CaptainId: captainId,
		}
		for _, gem := range captain.Gems() {
			if gem != nil {
				resp.Gems = append(resp.Gems, newGMGemInfo(gem))
			}
		}
		return
	})
	return marshalResponse(resp)
}

func (m *LightpawHandler) addYuanbao(heroId int64, toAdd uint64) []byte {
	if data := m.dep.Datas().MailHelp().SystemCompensation; data != nil {
		proto := data.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Prize = resdata.NewPrizeBuilder().AddYuanbao(toAdd).Build().Encode()
		m.dep.Mail().SendProtoMail(heroId, proto, m.time.CurrentTime())
	}
	return ok
}

func (m *LightpawHandler) addYinliang(heroId int64, toAdd uint64) []byte {
	if data := m.dep.Datas().MailHelp().SystemCompensation; data != nil {
		proto := data.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Prize = resdata.NewPrizeBuilder().AddYinliang(toAdd).Build().Encode()
		m.dep.Mail().SendProtoMail(heroId, proto, m.time.CurrentTime())
	}
	return ok
}

func (m *LightpawHandler) addDianquan(heroId int64, toAdd uint64) []byte {
	if data := m.dep.Datas().MailHelp().SystemCompensation; data != nil {
		proto := data.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Prize = resdata.NewPrizeBuilder().AddDianquan(toAdd).Build().Encode()
		m.dep.Mail().SendProtoMail(heroId, proto, m.time.CurrentTime())
	}
	return ok
}

// 禁言
func (m *LightpawHandler) banChat(heroId int64, banTimeAdd uint64) []byte {
	m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		hero.MiscData().SetBanChatEndTime(m.time.CurrentTime().Add(time.Second * time.Duration(banTimeAdd)))
		result.Add(chat.NewS2cBanChatMsg(timeutil.Marshal32(hero.MiscData().GetBanChatEndTime())))
		result.Ok()
	})
	return ok
}

// 禁言解除
func (m *LightpawHandler) banChatUnlock(heroId int64) []byte {
	m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		hero.MiscData().SetBanChatEndTime(m.time.CurrentTime())
		result.Add(chat.NewS2cBanChatMsg(timeutil.Marshal32(hero.MiscData().GetBanChatEndTime())))
		result.Ok()
	})
	return ok
}

// 封号
func (m *LightpawHandler) banLogin(heroId int64, banTimeAdd uint64) []byte {
	m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		addTime := time.Second * time.Duration(banTimeAdd)
		hero.MiscData().SetBanLoginEndTime(m.time.CurrentTime().Add(addTime))
		result.Add(login.NewS2cBanLoginMsg(timeutil.DurationMarshal32(addTime)))
		result.Ok()
	})
	if m.wordService.IsOnline(heroId) {
		m.wordService.FuncHero(heroId, func(id int64, hc iface.HeroController) {
			hc.Disconnect(misc.ErrDisconectReasonFailGm)
		})
	}
	return ok
}

// 解封
func (m *LightpawHandler) banLoginUnlock(heroId int64) []byte {
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		hero.MiscData().SetBanLoginEndTime(m.time.CurrentTime())
		return
	})
	return ok
}

func (m *LightpawHandler) sendPrizeMail(heroId int64, body []byte) []byte {
	var req GMSendPrizeMailRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return marshalResponse("sendPrizeMail:request struct not correct")
	}
	if err := checkGMSendPrizeMailRequest(&req, m.dep.Datas()); err != nil {
		return marshalResponse("sendPrizeMail:" + err.Error())
	}
	proto := req.encodeMail(m.dep.Datas())
	m.dep.Mail().SendProtoMail(heroId, proto, m.time.CurrentTime())
	return ok
}

func (m *LightpawHandler) sendPrizeMailAll(body []byte) []byte {
	var req GMSendPrizeMailRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return marshalResponse("sendPrizeMailAll:request struct not correct")
	}
	if err := checkGMSendPrizeMailRequest(&req, m.dep.Datas()); err != nil {
		return marshalResponse("sendPrizeMailAll:" + err.Error())
	}
	proto := req.encodeMail(m.dep.Datas())
	var heroIds []int64
	err := ctxfunc.Timeout2s(func(ctx context.Context) error {
		ids, e := m.dep.Db().HeroIds(ctx)
		if e == nil {
			heroIds = ids
		}
		return e
	})
	if err != nil {
		return marshalResponse(err.Error())
	}
	ctime := m.time.CurrentTime()
	for _, heroId := range heroIds {
		m.dep.Mail().SendProtoMail(heroId, proto, ctime)
	}
	return ok
}

func (m *LightpawHandler) requestHeroDepotGoods(heroId int64) []byte {
	var resp *GMHeroDepotGoodsResponse
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		resp = &GMHeroDepotGoodsResponse{}
		for id, count := range hero.Depot().GoodsMap() {
			resp.Goods = append(resp.Goods, newGMGoodsInfo(id, count, m.dep.Datas()))
		}
		return
	})
	return marshalResponse(resp)
}

func (m *LightpawHandler) requestHeroDepotEquip(heroId int64) []byte {
	var resp *GMHeroDepotEquipResponse
	m.heroService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		resp = &GMHeroDepotEquipResponse{}
		hero.Depot().WalkGenIdGoods(func(g goods.GenIdGoods) {
			if g.GoodsData().GoodsType() == goods.EQUIPMENT {
				resp.Equips = append(resp.Equips, newGMEquipInfo(g.(*entity.Equipment)))
			}
		})
		return
	})
	return marshalResponse(resp)
}
