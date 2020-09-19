package gm

import (
	"github.com/lightpaw/male7/gen/iface"
	activitypb "github.com/lightpaw/male7/gen/pb/activity"
	bai_zhanpb "github.com/lightpaw/male7/gen/pb/bai_zhan"
	chatpb "github.com/lightpaw/male7/gen/pb/chat"
	client_configpb "github.com/lightpaw/male7/gen/pb/client_config"
	countrypb "github.com/lightpaw/male7/gen/pb/country"
	depotpb "github.com/lightpaw/male7/gen/pb/depot"
	dianquanpb "github.com/lightpaw/male7/gen/pb/dianquan"
	domesticpb "github.com/lightpaw/male7/gen/pb/domestic"
	dungeonpb "github.com/lightpaw/male7/gen/pb/dungeon"
	equipmentpb "github.com/lightpaw/male7/gen/pb/equipment"
	farmpb "github.com/lightpaw/male7/gen/pb/farm"
	fishingpb "github.com/lightpaw/male7/gen/pb/fishing"
	gardenpb "github.com/lightpaw/male7/gen/pb/garden"
	gempb "github.com/lightpaw/male7/gen/pb/gem"
	guildpb "github.com/lightpaw/male7/gen/pb/guild"
	hebipb "github.com/lightpaw/male7/gen/pb/hebi"
	mailpb "github.com/lightpaw/male7/gen/pb/mail"
	militarypb "github.com/lightpaw/male7/gen/pb/military"
	mingcpb "github.com/lightpaw/male7/gen/pb/mingc"
	mingc_warpb "github.com/lightpaw/male7/gen/pb/mingc_war"
	miscpb "github.com/lightpaw/male7/gen/pb/misc"
	promotionpb "github.com/lightpaw/male7/gen/pb/promotion"
	questionpb "github.com/lightpaw/male7/gen/pb/question"
	random_eventpb "github.com/lightpaw/male7/gen/pb/random_event"
	rankpb "github.com/lightpaw/male7/gen/pb/rank"
	red_packetpb "github.com/lightpaw/male7/gen/pb/red_packet"
	regionpb "github.com/lightpaw/male7/gen/pb/region"
	relationpb "github.com/lightpaw/male7/gen/pb/relation"
	secret_towerpb "github.com/lightpaw/male7/gen/pb/secret_tower"
	shoppb "github.com/lightpaw/male7/gen/pb/shop"
	strategypb "github.com/lightpaw/male7/gen/pb/strategy"
	stresspb "github.com/lightpaw/male7/gen/pb/stress"
	surveypb "github.com/lightpaw/male7/gen/pb/survey"
	tagpb "github.com/lightpaw/male7/gen/pb/tag"
	taskpb "github.com/lightpaw/male7/gen/pb/task"
	teachpb "github.com/lightpaw/male7/gen/pb/teach"
	towerpb "github.com/lightpaw/male7/gen/pb/tower"
	vippb "github.com/lightpaw/male7/gen/pb/vip"
	xiongnupb "github.com/lightpaw/male7/gen/pb/xiongnu"
	xuanyuanpb "github.com/lightpaw/male7/gen/pb/xuanyuan"
	zhanjiangpb "github.com/lightpaw/male7/gen/pb/zhanjiang"
	zhengwupb "github.com/lightpaw/male7/gen/pb/zhengwu"

	"strings"
)

func (gm *GmModule) initModuleHandler() []*gm_group {
	return []*gm_group{
		gm.initHandlerActivity(),
		gm.initHandlerBai_zhan(),
		gm.initHandlerChat(),
		gm.initHandlerClient_config(),
		gm.initHandlerCountry(),
		gm.initHandlerDepot(),
		gm.initHandlerDianquan(),
		gm.initHandlerDomestic(),
		gm.initHandlerDungeon(),
		gm.initHandlerEquipment(),
		gm.initHandlerFarm(),
		gm.initHandlerFishing(),
		gm.initHandlerGarden(),
		gm.initHandlerGem(),
		gm.initHandlerGuild(),
		gm.initHandlerHebi(),
		gm.initHandlerMail(),
		gm.initHandlerMilitary(),
		gm.initHandlerMingc(),
		gm.initHandlerMingc_war(),
		gm.initHandlerMisc(),
		gm.initHandlerPromotion(),
		gm.initHandlerQuestion(),
		gm.initHandlerRandom_event(),
		gm.initHandlerRank(),
		gm.initHandlerRed_packet(),
		gm.initHandlerRegion(),
		gm.initHandlerRelation(),
		gm.initHandlerSecret_tower(),
		gm.initHandlerShop(),
		gm.initHandlerStrategy(),
		gm.initHandlerStress(),
		gm.initHandlerSurvey(),
		gm.initHandlerTag(),
		gm.initHandlerTask(),
		gm.initHandlerTeach(),
		gm.initHandlerTower(),
		gm.initHandlerVip(),
		gm.initHandlerXiongnu(),
		gm.initHandlerXuanyuan(),
		gm.initHandlerZhanjiang(),
		gm.initHandlerZhengwu(),
	}
}

func (gm *GmModule) initHandlerActivity() *gm_group {
	group := &gm_group{
		tab: "activity",
		handler: []*gm_handler{
			newStringHandler("collect_collection", " ", gm.handleActivity_c2s_collect_collection),
		},
	}

	return group
}

func (gm *GmModule) handleActivity_c2s_collect_collection(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ActivityModule().(interface {
			ProcessCollectCollection(*activitypb.C2SCollectCollectionProto, iface.HeroController)
		}).ProcessCollectCollection(&activitypb.C2SCollectCollectionProto{

			Id: parseInt32(strArray[0]),

			ExchangeId: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerBai_zhan() *gm_group {
	group := &gm_group{
		tab: "bai_zhan",
		handler: []*gm_handler{
			newStringHandler("query_bai_zhan_info", " ", gm.handleBai_zhan_c2s_query_bai_zhan_info),
			newStringHandler("clear_last_jun_xian", " ", gm.handleBai_zhan_c2s_clear_last_jun_xian),
			newStringHandler("bai_zhan_challenge", " ", gm.handleBai_zhan_c2s_bai_zhan_challenge),
			newStringHandler("collect_salary", " ", gm.handleBai_zhan_c2s_collect_salary),
			newStringHandler("collect_jun_xian_prize", " ", gm.handleBai_zhan_c2s_collect_jun_xian_prize),
			newStringHandler("self_record", " ", gm.handleBai_zhan_c2s_self_record),
			newStringHandler("request_rank", " ", gm.handleBai_zhan_c2s_request_rank),
			newStringHandler("request_self_rank", " ", gm.handleBai_zhan_c2s_request_self_rank),
		},
	}

	return group
}

func (gm *GmModule) handleBai_zhan_c2s_query_bai_zhan_info(amount string, hc iface.HeroController) {
	gm.modules.BaiZhanModule().(interface {
		ProcessQueryBaiZhanInfo(iface.HeroController)
	}).ProcessQueryBaiZhanInfo(hc)
}
func (gm *GmModule) handleBai_zhan_c2s_clear_last_jun_xian(amount string, hc iface.HeroController) {
	gm.modules.BaiZhanModule().(interface {
		ProcesClearLastJunXian(iface.HeroController)
	}).ProcesClearLastJunXian(hc)
}
func (gm *GmModule) handleBai_zhan_c2s_bai_zhan_challenge(amount string, hc iface.HeroController) {
	gm.modules.BaiZhanModule().(interface {
		ProcessChallenge(iface.HeroController)
	}).ProcessChallenge(hc)
}
func (gm *GmModule) handleBai_zhan_c2s_collect_salary(amount string, hc iface.HeroController) {
	gm.modules.BaiZhanModule().(interface {
		ProcessCollectSalary(iface.HeroController)
	}).ProcessCollectSalary(hc)
}
func (gm *GmModule) handleBai_zhan_c2s_collect_jun_xian_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.BaiZhanModule().(interface {
			ProcessCollectJunXianPrize(*bai_zhanpb.C2SCollectJunXianPrizeProto, iface.HeroController)
		}).ProcessCollectJunXianPrize(&bai_zhanpb.C2SCollectJunXianPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleBai_zhan_c2s_self_record(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.BaiZhanModule().(interface {
			ProcessSelfRecord(*bai_zhanpb.C2SSelfRecordProto, iface.HeroController)
		}).ProcessSelfRecord(&bai_zhanpb.C2SSelfRecordProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleBai_zhan_c2s_request_rank(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.BaiZhanModule().(interface {
			ProcessRequestRank(*bai_zhanpb.C2SRequestRankProto, iface.HeroController)
		}).ProcessRequestRank(&bai_zhanpb.C2SRequestRankProto{

			Self: parseBool(strArray[0]),

			StartRank: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleBai_zhan_c2s_request_self_rank(amount string, hc iface.HeroController) {
	gm.modules.BaiZhanModule().(interface {
		ProcessRequestSelfRank(iface.HeroController)
	}).ProcessRequestSelfRank(hc)
}

func (gm *GmModule) initHandlerChat() *gm_group {
	group := &gm_group{
		tab: "chat",
		handler: []*gm_handler{
			newStringHandler("world_chat", " ", gm.handleChat_c2s_world_chat),
			newStringHandler("guild_chat", " ", gm.handleChat_c2s_guild_chat),
			newStringHandler("self_chat_window", " ", gm.handleChat_c2s_self_chat_window),
			newStringHandler("create_self_chat_window", " ", gm.handleChat_c2s_create_self_chat_window),
			newStringHandler("remove_chat_window", " ", gm.handleChat_c2s_remove_chat_window),
			newStringHandler("list_history_chat", " ", gm.handleChat_c2s_list_history_chat),
			newStringHandler("send_chat", " ", gm.handleChat_c2s_send_chat),
			newStringHandler("read_chat_msg", " ", gm.handleChat_c2s_read_chat_msg),
			newStringHandler("get_hero_chat_info", " ", gm.handleChat_c2s_get_hero_chat_info),
		},
	}

	return group
}

func (gm *GmModule) handleChat_c2s_world_chat(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessWorldChat(*chatpb.C2SWorldChatProto, iface.HeroController)
		}).ProcessWorldChat(&chatpb.C2SWorldChatProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_guild_chat(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessGuildChat(*chatpb.C2SGuildChatProto, iface.HeroController)
		}).ProcessGuildChat(&chatpb.C2SGuildChatProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_self_chat_window(amount string, hc iface.HeroController) {
	gm.modules.ChatModule().(interface {
		ProcessSelfChatWindow(iface.HeroController)
	}).ProcessSelfChatWindow(hc)
}
func (gm *GmModule) handleChat_c2s_create_self_chat_window(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessCreateSelfChatWindow(*chatpb.C2SCreateSelfChatWindowProto, iface.HeroController)
		}).ProcessCreateSelfChatWindow(&chatpb.C2SCreateSelfChatWindowProto{

			Target: parseBytes(strArray[0]),

			SetUp: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_remove_chat_window(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessRemoveChatWindow(*chatpb.C2SRemoveChatWindowProto, iface.HeroController)
		}).ProcessRemoveChatWindow(&chatpb.C2SRemoveChatWindowProto{

			ChatType: parseInt32(strArray[0]),

			ChatTarget: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_list_history_chat(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessListHistoryChat(*chatpb.C2SListHistoryChatProto, iface.HeroController)
		}).ProcessListHistoryChat(&chatpb.C2SListHistoryChatProto{

			ChatType: parseInt32(strArray[0]),

			ChatTarget: parseBytes(strArray[1]),

			MinChatId: parseBytes(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_send_chat(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessSendChat(*chatpb.C2SSendChatProto, iface.HeroController)
		}).ProcessSendChat(&chatpb.C2SSendChatProto{

			ChatMsg: parseBytes(strArray[0]),

			Receiver: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_read_chat_msg(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessReadChatMsg(*chatpb.C2SReadChatMsgProto, iface.HeroController)
		}).ProcessReadChatMsg(&chatpb.C2SReadChatMsgProto{

			ChatType: parseInt32(strArray[0]),

			ChatTarget: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleChat_c2s_get_hero_chat_info(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ChatModule().(interface {
			ProcessGetHeroChatInfo(*chatpb.C2SGetHeroChatInfoProto, iface.HeroController)
		}).ProcessGetHeroChatInfo(&chatpb.C2SGetHeroChatInfoProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerClient_config() *gm_group {
	group := &gm_group{
		tab: "client_config",
		handler: []*gm_handler{
			newStringHandler("config", " ", gm.handleClient_config_c2s_config),
			newStringHandler("set_client_data", " ", gm.handleClient_config_c2s_set_client_data),
			newStringHandler("set_client_key", " ", gm.handleClient_config_c2s_set_client_key),
		},
	}

	return group
}

func (gm *GmModule) handleClient_config_c2s_config(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ClientConfigModule().(interface {
			ProcessConfig(*client_configpb.C2SConfigProto, iface.HeroController)
		}).ProcessConfig(&client_configpb.C2SConfigProto{

			Path: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleClient_config_c2s_set_client_data(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ClientConfigModule().(interface {
			ProcessSetClientData(*client_configpb.C2SSetClientDataProto, iface.HeroController)
		}).ProcessSetClientData(&client_configpb.C2SSetClientDataProto{

			Index: parseInt32(strArray[0]),

			ToSetBool: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleClient_config_c2s_set_client_key(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ClientConfigModule().(interface {
			ProcessSetClientKey(*client_configpb.C2SSetClientKeyProto, iface.HeroController)
		}).ProcessSetClientKey(&client_configpb.C2SSetClientKeyProto{

			KeyType: parseInt32(strArray[0]),

			KeyValue: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerCountry() *gm_group {
	group := &gm_group{
		tab: "country",
		handler: []*gm_handler{
			newStringHandler("request_country_prestige", " ", gm.handleCountry_c2s_request_country_prestige),
			newStringHandler("request_countries", " ", gm.handleCountry_c2s_request_countries),
			newStringHandler("hero_change_country", " ", gm.handleCountry_c2s_hero_change_country),
			newStringHandler("country_detail", " ", gm.handleCountry_c2s_country_detail),
			newStringHandler("official_appoint", " ", gm.handleCountry_c2s_official_appoint),
			newStringHandler("official_depose", " ", gm.handleCountry_c2s_official_depose),
			newStringHandler("official_leave", " ", gm.handleCountry_c2s_official_leave),
			newStringHandler("collect_official_salary", " ", gm.handleCountry_c2s_collect_official_salary),
			newStringHandler("change_name_start", " ", gm.handleCountry_c2s_change_name_start),
			newStringHandler("change_name_vote", " ", gm.handleCountry_c2s_change_name_vote),
			newStringHandler("search_to_appoint_hero_list", " ", gm.handleCountry_c2s_search_to_appoint_hero_list),
			newStringHandler("default_to_appoint_hero_list", " ", gm.handleCountry_c2s_default_to_appoint_hero_list),
		},
	}

	return group
}

func (gm *GmModule) handleCountry_c2s_request_country_prestige(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessRequestCountryPrestige(*countrypb.C2SRequestCountryPrestigeProto, iface.HeroController)
		}).ProcessRequestCountryPrestige(&countrypb.C2SRequestCountryPrestigeProto{

			Vsn: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_request_countries(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessRequestCountries(*countrypb.C2SRequestCountriesProto, iface.HeroController)
		}).ProcessRequestCountries(&countrypb.C2SRequestCountriesProto{

			Vsn: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_hero_change_country(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessHeroChangeCountry(*countrypb.C2SHeroChangeCountryProto, iface.HeroController)
		}).ProcessHeroChangeCountry(&countrypb.C2SHeroChangeCountryProto{

			NewCountry: parseInt32(strArray[0]),

			Buy: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_country_detail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessCountryDetail(*countrypb.C2SCountryDetailProto, iface.HeroController)
		}).ProcessCountryDetail(&countrypb.C2SCountryDetailProto{

			CountryId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_official_appoint(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessOfficialAppoint(*countrypb.C2SOfficialAppointProto, iface.HeroController)
		}).ProcessOfficialAppoint(&countrypb.C2SOfficialAppointProto{

			HeroId: parseBytes(strArray[1]),

			Pos: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_official_depose(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessOfficialDepose(*countrypb.C2SOfficialDeposeProto, iface.HeroController)
		}).ProcessOfficialDepose(&countrypb.C2SOfficialDeposeProto{

			HeroId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_official_leave(amount string, hc iface.HeroController) {
	gm.modules.CountryModule().(interface {
		ProcessOfficialLeave(iface.HeroController)
	}).ProcessOfficialLeave(hc)
}
func (gm *GmModule) handleCountry_c2s_collect_official_salary(amount string, hc iface.HeroController) {
	gm.modules.CountryModule().(interface {
		ProcessCollectOfficialSalary(iface.HeroController)
	}).ProcessCollectOfficialSalary(hc)
}
func (gm *GmModule) handleCountry_c2s_change_name_start(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessChangeNameStart(*countrypb.C2SChangeNameStartProto, iface.HeroController)
		}).ProcessChangeNameStart(&countrypb.C2SChangeNameStartProto{

			NewName: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_change_name_vote(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessChangeNameVote(*countrypb.C2SChangeNameVoteProto, iface.HeroController)
		}).ProcessChangeNameVote(&countrypb.C2SChangeNameVoteProto{

			Agree: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_search_to_appoint_hero_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.CountryModule().(interface {
			ProcessSearchToAppointHeroList(*countrypb.C2SSearchToAppointHeroListProto, iface.HeroController)
		}).ProcessSearchToAppointHeroList(&countrypb.C2SSearchToAppointHeroListProto{

			Name: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleCountry_c2s_default_to_appoint_hero_list(amount string, hc iface.HeroController) {
	gm.modules.CountryModule().(interface {
		ProcessDefaultToAppointHeroList(iface.HeroController)
	}).ProcessDefaultToAppointHeroList(hc)
}

func (gm *GmModule) initHandlerDepot() *gm_group {
	group := &gm_group{
		tab: "depot",
		handler: []*gm_handler{
			newStringHandler("use_goods", " ", gm.handleDepot_c2s_use_goods),
			newStringHandler("use_cdr_goods", " ", gm.handleDepot_c2s_use_cdr_goods),
			newStringHandler("goods_combine", " ", gm.handleDepot_c2s_goods_combine),
			newStringHandler("goods_parts_combine", " ", gm.handleDepot_c2s_goods_parts_combine),
			newStringHandler("unlock_baowu", " ", gm.handleDepot_c2s_unlock_baowu),
			newStringHandler("collect_baowu", " ", gm.handleDepot_c2s_collect_baowu),
			newStringHandler("list_baowu_log", " ", gm.handleDepot_c2s_list_baowu_log),
			newStringHandler("decompose_baowu", " ", gm.handleDepot_c2s_decompose_baowu),
		},
	}

	return group
}

func (gm *GmModule) handleDepot_c2s_use_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessUseGoods(*depotpb.C2SUseGoodsProto, iface.HeroController)
		}).ProcessUseGoods(&depotpb.C2SUseGoodsProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_use_cdr_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessUseCdrGoods(*depotpb.C2SUseCdrGoodsProto, iface.HeroController)
		}).ProcessUseCdrGoods(&depotpb.C2SUseCdrGoodsProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),

			CdrType: parseInt32(strArray[2]),

			Index: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_goods_combine(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessGoodsCombine(*depotpb.C2SGoodsCombineProto, iface.HeroController)
		}).ProcessGoodsCombine(&depotpb.C2SGoodsCombineProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_goods_parts_combine(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessGoodsPartCombine(*depotpb.C2SGoodsPartsCombineProto, iface.HeroController)
		}).ProcessGoodsPartCombine(&depotpb.C2SGoodsPartsCombineProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),

			SelectIndex: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_unlock_baowu(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessUnlockBaowu(*depotpb.C2SUnlockBaowuProto, iface.HeroController)
		}).ProcessUnlockBaowu(&depotpb.C2SUnlockBaowuProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_collect_baowu(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessCollectBaowu(*depotpb.C2SCollectBaowuProto, iface.HeroController)
		}).ProcessCollectBaowu(&depotpb.C2SCollectBaowuProto{

			Miao: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDepot_c2s_list_baowu_log(amount string, hc iface.HeroController) {
	gm.modules.DepotModule().(interface {
		ProcessListBaowuLog(iface.HeroController)
	}).ProcessListBaowuLog(hc)
}
func (gm *GmModule) handleDepot_c2s_decompose_baowu(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DepotModule().(interface {
			ProcessDecomposeBaowu(*depotpb.C2SDecomposeBaowuProto, iface.HeroController)
		}).ProcessDecomposeBaowu(&depotpb.C2SDecomposeBaowuProto{

			BaowuId: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerDianquan() *gm_group {
	group := &gm_group{
		tab: "dianquan",
		handler: []*gm_handler{
			newStringHandler("exchange", " ", gm.handleDianquan_c2s_exchange),
		},
	}

	return group
}

func (gm *GmModule) handleDianquan_c2s_exchange(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DianquanModule().(interface {
			ProcessExchange(*dianquanpb.C2SExchangeProto, iface.HeroController)
		}).ProcessExchange(&dianquanpb.C2SExchangeProto{

			Times: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerDomestic() *gm_group {
	group := &gm_group{
		tab: "domestic",
		handler: []*gm_handler{
			newStringHandler("create_building", " ", gm.handleDomestic_c2s_create_building),
			newStringHandler("upgrade_building", " ", gm.handleDomestic_c2s_upgrade_building),
			newStringHandler("rebuild_resource_building", " ", gm.handleDomestic_c2s_rebuild_resource_building),
			newStringHandler("unlock_outer_city", " ", gm.handleDomestic_c2s_unlock_outer_city),
			newStringHandler("update_outer_city_type", " ", gm.handleDomestic_c2s_update_outer_city_type),
			newStringHandler("upgrade_outer_city_building", " ", gm.handleDomestic_c2s_upgrade_outer_city_building),
			newStringHandler("collect_resource", " ", gm.handleDomestic_c2s_collect_resource),
			newStringHandler("collect_resource_v2", " ", gm.handleDomestic_c2s_collect_resource_v2),
			newStringHandler("request_resource_conflict", " ", gm.handleDomestic_c2s_request_resource_conflict),
			newStringHandler("learn_technology", " ", gm.handleDomestic_c2s_learn_technology),
			newStringHandler("unlock_stable_building", " ", gm.handleDomestic_c2s_unlock_stable_building),
			newStringHandler("upgrade_stable_building", " ", gm.handleDomestic_c2s_upgrade_stable_building),
			newStringHandler("is_hero_name_exist", " ", gm.handleDomestic_c2s_is_hero_name_exist),
			newStringHandler("change_hero_name", " ", gm.handleDomestic_c2s_change_hero_name),
			newStringHandler("list_old_name", " ", gm.handleDomestic_c2s_list_old_name),
			newStringHandler("view_other_hero", " ", gm.handleDomestic_c2s_view_other_hero),
			newStringHandler("view_fight_info", " ", gm.handleDomestic_c2s_view_fight_info),
			newStringHandler("miao_building_worker_cd", " ", gm.handleDomestic_c2s_miao_building_worker_cd),
			newStringHandler("miao_tech_worker_cd", " ", gm.handleDomestic_c2s_miao_tech_worker_cd),
			newStringHandler("forging_equip", " ", gm.handleDomestic_c2s_forging_equip),
			newStringHandler("sign", " ", gm.handleDomestic_c2s_sign),
			newStringHandler("voice", " ", gm.handleDomestic_c2s_voice),
			newStringHandler("request_city_exchange_event", " ", gm.handleDomestic_c2s_request_city_exchange_event),
			newStringHandler("city_event_exchange", " ", gm.handleDomestic_c2s_city_event_exchange),
			newStringHandler("change_head", " ", gm.handleDomestic_c2s_change_head),
			newStringHandler("change_body", " ", gm.handleDomestic_c2s_change_body),
			newStringHandler("collect_countdown_prize", " ", gm.handleDomestic_c2s_collect_countdown_prize),
			newStringHandler("start_workshop", " ", gm.handleDomestic_c2s_start_workshop),
			newStringHandler("collect_workshop", " ", gm.handleDomestic_c2s_collect_workshop),
			newStringHandler("workshop_miao_cd", " ", gm.handleDomestic_c2s_workshop_miao_cd),
			newStringHandler("refresh_workshop", " ", gm.handleDomestic_c2s_refresh_workshop),
			newStringHandler("collect_season_prize", " ", gm.handleDomestic_c2s_collect_season_prize),
			newStringHandler("buy_sp", " ", gm.handleDomestic_c2s_buy_sp),
			newStringHandler("use_buf_effect", " ", gm.handleDomestic_c2s_use_buf_effect),
			newStringHandler("open_buf_effect_ui", " ", gm.handleDomestic_c2s_open_buf_effect_ui),
			newStringHandler("use_advantage", " ", gm.handleDomestic_c2s_use_advantage),
			newStringHandler("worker_unlock", " ", gm.handleDomestic_c2s_worker_unlock),
		},
	}

	return group
}

func (gm *GmModule) handleDomestic_c2s_create_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessCreateResourceBuilding(*domesticpb.C2SCreateBuildingProto, iface.HeroController)
		}).ProcessCreateResourceBuilding(&domesticpb.C2SCreateBuildingProto{

			Id: parseInt32(strArray[0]),

			Type: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_upgrade_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUpgradeResourceBuilding(*domesticpb.C2SUpgradeBuildingProto, iface.HeroController)
		}).ProcessUpgradeResourceBuilding(&domesticpb.C2SUpgradeBuildingProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_rebuild_resource_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessRebuildBuilding(*domesticpb.C2SRebuildResourceBuildingProto, iface.HeroController)
		}).ProcessRebuildBuilding(&domesticpb.C2SRebuildResourceBuildingProto{

			Id: parseInt32(strArray[0]),

			Type: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_unlock_outer_city(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUnlockOuterCity(*domesticpb.C2SUnlockOuterCityProto, iface.HeroController)
		}).ProcessUnlockOuterCity(&domesticpb.C2SUnlockOuterCityProto{

			Id: parseInt32(strArray[0]),

			T: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_update_outer_city_type(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUpdateOuterCityType(*domesticpb.C2SUpdateOuterCityTypeProto, iface.HeroController)
		}).ProcessUpdateOuterCityType(&domesticpb.C2SUpdateOuterCityTypeProto{

			Id: parseInt32(strArray[0]),

			T: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_upgrade_outer_city_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUpgradeOuterCityBuilding(*domesticpb.C2SUpgradeOuterCityBuildingProto, iface.HeroController)
		}).ProcessUpgradeOuterCityBuilding(&domesticpb.C2SUpgradeOuterCityBuildingProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_collect_resource(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessCollectResource(*domesticpb.C2SCollectResourceProto, iface.HeroController)
		}).ProcessCollectResource(&domesticpb.C2SCollectResourceProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_collect_resource_v2(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessCollectResourceV2(*domesticpb.C2SCollectResourceV2Proto, iface.HeroController)
		}).ProcessCollectResourceV2(&domesticpb.C2SCollectResourceV2Proto{

			ResType: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_request_resource_conflict(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessRequestResourceConflict(iface.HeroController)
	}).ProcessRequestResourceConflict(hc)
}
func (gm *GmModule) handleDomestic_c2s_learn_technology(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessLearnTechnology(*domesticpb.C2SLearnTechnologyProto, iface.HeroController)
		}).ProcessLearnTechnology(&domesticpb.C2SLearnTechnologyProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_unlock_stable_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUnlockStableBuilding(*domesticpb.C2SUnlockStableBuildingProto, iface.HeroController)
		}).ProcessUnlockStableBuilding(&domesticpb.C2SUnlockStableBuildingProto{

			Type: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_upgrade_stable_building(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUpgradeStableBuilding(*domesticpb.C2SUpgradeStableBuildingProto, iface.HeroController)
		}).ProcessUpgradeStableBuilding(&domesticpb.C2SUpgradeStableBuildingProto{

			Type: parseInt32(strArray[0]),

			Level: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_is_hero_name_exist(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessIsHeroNameExist(*domesticpb.C2SIsHeroNameExistProto, iface.HeroController)
		}).ProcessIsHeroNameExist(&domesticpb.C2SIsHeroNameExistProto{

			Name: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_change_hero_name(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessChangeHeroName(*domesticpb.C2SChangeHeroNameProto, iface.HeroController)
		}).ProcessChangeHeroName(&domesticpb.C2SChangeHeroNameProto{

			Name: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_list_old_name(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessListOldName(*domesticpb.C2SListOldNameProto, iface.HeroController)
		}).ProcessListOldName(&domesticpb.C2SListOldNameProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_view_other_hero(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessViewOtherHero(*domesticpb.C2SViewOtherHeroProto, iface.HeroController)
		}).ProcessViewOtherHero(&domesticpb.C2SViewOtherHeroProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_view_fight_info(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessViewFightInfo(*domesticpb.C2SViewFightInfoProto, iface.HeroController)
		}).ProcessViewFightInfo(&domesticpb.C2SViewFightInfoProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_miao_building_worker_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessMiaoBuildingWorkerCd(*domesticpb.C2SMiaoBuildingWorkerCdProto, iface.HeroController)
		}).ProcessMiaoBuildingWorkerCd(&domesticpb.C2SMiaoBuildingWorkerCdProto{

			WorkerPos: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_miao_tech_worker_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessMiaoTechWorkerCd(*domesticpb.C2SMiaoTechWorkerCdProto, iface.HeroController)
		}).ProcessMiaoTechWorkerCd(&domesticpb.C2SMiaoTechWorkerCdProto{

			WorkerPos: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_forging_equip(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessForgingEquip(*domesticpb.C2SForgingEquipProto, iface.HeroController)
		}).ProcessForgingEquip(&domesticpb.C2SForgingEquipProto{

			Slot: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_sign(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessSign(*domesticpb.C2SSignProto, iface.HeroController)
		}).ProcessSign(&domesticpb.C2SSignProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_voice(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessVoice(*domesticpb.C2SVoiceProto, iface.HeroController)
		}).ProcessVoice(&domesticpb.C2SVoiceProto{

			Content: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_request_city_exchange_event(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessRequestCityExchangeEvent(iface.HeroController)
	}).ProcessRequestCityExchangeEvent(hc)
}
func (gm *GmModule) handleDomestic_c2s_city_event_exchange(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessCityEventExchange(*domesticpb.C2SCityEventExchangeProto, iface.HeroController)
		}).ProcessCityEventExchange(&domesticpb.C2SCityEventExchangeProto{

			GiveUp: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_change_head(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessChangeHead(*domesticpb.C2SChangeHeadProto, iface.HeroController)
		}).ProcessChangeHead(&domesticpb.C2SChangeHeadProto{

			HeadId: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_change_body(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessChangeBody(*domesticpb.C2SChangeBodyProto, iface.HeroController)
		}).ProcessChangeBody(&domesticpb.C2SChangeBodyProto{

			BodyId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_collect_countdown_prize(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessCollectCountdownPrize(iface.HeroController)
	}).ProcessCollectCountdownPrize(hc)
}
func (gm *GmModule) handleDomestic_c2s_start_workshop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessWorkshopStartForge(*domesticpb.C2SStartWorkshopProto, iface.HeroController)
		}).ProcessWorkshopStartForge(&domesticpb.C2SStartWorkshopProto{

			Index: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_collect_workshop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessWorkshopCollect(*domesticpb.C2SCollectWorkshopProto, iface.HeroController)
		}).ProcessWorkshopCollect(&domesticpb.C2SCollectWorkshopProto{

			Index: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_workshop_miao_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessWorkshopMiaoCd(*domesticpb.C2SWorkshopMiaoCdProto, iface.HeroController)
		}).ProcessWorkshopMiaoCd(&domesticpb.C2SWorkshopMiaoCdProto{

			Index: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_refresh_workshop(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessRefreshWorkshop(iface.HeroController)
	}).ProcessRefreshWorkshop(hc)
}
func (gm *GmModule) handleDomestic_c2s_collect_season_prize(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessCollectSeasonPrize(iface.HeroController)
	}).ProcessCollectSeasonPrize(hc)
}
func (gm *GmModule) handleDomestic_c2s_buy_sp(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessBuySp(*domesticpb.C2SBuySpProto, iface.HeroController)
		}).ProcessBuySp(&domesticpb.C2SBuySpProto{

			BuyTimes: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_use_buf_effect(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUseBufEffect(*domesticpb.C2SUseBufEffectProto, iface.HeroController)
		}).ProcessUseBufEffect(&domesticpb.C2SUseBufEffectProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_open_buf_effect_ui(amount string, hc iface.HeroController) {
	gm.modules.DomesticModule().(interface {
		ProcessOpenBufEffectUi(iface.HeroController)
	}).ProcessOpenBufEffectUi(hc)
}
func (gm *GmModule) handleDomestic_c2s_use_advantage(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessUseAdvantage(*domesticpb.C2SUseAdvantageProto, iface.HeroController)
		}).ProcessUseAdvantage(&domesticpb.C2SUseAdvantageProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDomestic_c2s_worker_unlock(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DomesticModule().(interface {
			ProcessWorkerUnlock(*domesticpb.C2SWorkerUnlockProto, iface.HeroController)
		}).ProcessWorkerUnlock(&domesticpb.C2SWorkerUnlockProto{

			Pos: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerDungeon() *gm_group {
	group := &gm_group{
		tab: "dungeon",
		handler: []*gm_handler{
			newStringHandler("challenge", " ", gm.handleDungeon_c2s_challenge),
			newStringHandler("collect_chapter_prize", " ", gm.handleDungeon_c2s_collect_chapter_prize),
			newStringHandler("collect_pass_dungeon_prize", " ", gm.handleDungeon_c2s_collect_pass_dungeon_prize),
			newStringHandler("auto_challenge", " ", gm.handleDungeon_c2s_auto_challenge),
			newStringHandler("collect_chapter_star_prize", " ", gm.handleDungeon_c2s_collect_chapter_star_prize),
		},
	}

	return group
}

func (gm *GmModule) handleDungeon_c2s_challenge(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DungeonModule().(interface {
			ProcessChallenge(*dungeonpb.C2SChallengeProto, iface.HeroController)
		}).ProcessChallenge(&dungeonpb.C2SChallengeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDungeon_c2s_collect_chapter_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DungeonModule().(interface {
			ProcessCollectChapterPrize(*dungeonpb.C2SCollectChapterPrizeProto, iface.HeroController)
		}).ProcessCollectChapterPrize(&dungeonpb.C2SCollectChapterPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDungeon_c2s_collect_pass_dungeon_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DungeonModule().(interface {
			ProcessCollectPassDungeonPrize(*dungeonpb.C2SCollectPassDungeonPrizeProto, iface.HeroController)
		}).ProcessCollectPassDungeonPrize(&dungeonpb.C2SCollectPassDungeonPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleDungeon_c2s_auto_challenge(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DungeonModule().(interface {
			ProcessAutoChallenge(*dungeonpb.C2SAutoChallengeProto, iface.HeroController)
		}).ProcessAutoChallenge(&dungeonpb.C2SAutoChallengeProto{

			Id: parseInt32(strArray[0]),

			Times: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleDungeon_c2s_collect_chapter_star_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.DungeonModule().(interface {
			ProcessCollectChapterStarPrize(*dungeonpb.C2SCollectChapterStarPrizeProto, iface.HeroController)
		}).ProcessCollectChapterStarPrize(&dungeonpb.C2SCollectChapterStarPrizeProto{

			Id: parseInt32(strArray[0]),

			CollectN: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerEquipment() *gm_group {
	group := &gm_group{
		tab: "equipment",
		handler: []*gm_handler{
			newStringHandler("view_chat_equip", " ", gm.handleEquipment_c2s_view_chat_equip),
			newStringHandler("wear_equipment", " ", gm.handleEquipment_c2s_wear_equipment),
			newStringHandler("upgrade_equipment", " ", gm.handleEquipment_c2s_upgrade_equipment),
			newStringHandler("upgrade_equipment_all", " ", gm.handleEquipment_c2s_upgrade_equipment_all),
			newStringHandler("refined_equipment", " ", gm.handleEquipment_c2s_refined_equipment),
			newStringHandler("smelt_equipment", " ", gm.handleEquipment_c2s_smelt_equipment),
			newStringHandler("rebuild_equipment", " ", gm.handleEquipment_c2s_rebuild_equipment),
			newStringHandler("one_key_take_off", " ", gm.handleEquipment_c2s_one_key_take_off),
		},
	}

	return group
}

func (gm *GmModule) handleEquipment_c2s_view_chat_equip(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessViewChatEquip(*equipmentpb.C2SViewChatEquipProto, iface.HeroController)
		}).ProcessViewChatEquip(&equipmentpb.C2SViewChatEquipProto{

			DataId: parseInt32(strArray[0]),

			Level: parseInt32(strArray[1]),

			Refined: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_wear_equipment(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessWearEquipment(*equipmentpb.C2SWearEquipmentProto, iface.HeroController)
		}).ProcessWearEquipment(&equipmentpb.C2SWearEquipmentProto{

			CaptainId: parseInt32(strArray[0]),

			EquipmentId: parseInt32(strArray[1]),

			Down: parseBool(strArray[2]),

			Inhert: parseBool(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_upgrade_equipment(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessUpgradeEquipment(*equipmentpb.C2SUpgradeEquipmentProto, iface.HeroController)
		}).ProcessUpgradeEquipment(&equipmentpb.C2SUpgradeEquipmentProto{

			CaptainId: parseInt32(strArray[0]),

			EquipmentId: parseInt32(strArray[1]),

			UpgradeTimes: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_upgrade_equipment_all(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessUpgradeEquipmentAll(*equipmentpb.C2SUpgradeEquipmentAllProto, iface.HeroController)
		}).ProcessUpgradeEquipmentAll(&equipmentpb.C2SUpgradeEquipmentAllProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_refined_equipment(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessRefinedEquipment(*equipmentpb.C2SRefinedEquipmentProto, iface.HeroController)
		}).ProcessRefinedEquipment(&equipmentpb.C2SRefinedEquipmentProto{

			CaptainId: parseInt32(strArray[0]),

			EquipmentId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_smelt_equipment(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessSmeltEquipment(*equipmentpb.C2SSmeltEquipmentProto, iface.HeroController)
		}).ProcessSmeltEquipment(&equipmentpb.C2SSmeltEquipmentProto{

			EquipmentId: parseInt32Array(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_rebuild_equipment(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessRebuildEquipment(*equipmentpb.C2SRebuildEquipmentProto, iface.HeroController)
		}).ProcessRebuildEquipment(&equipmentpb.C2SRebuildEquipmentProto{

			EquipmentId: parseInt32Array(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleEquipment_c2s_one_key_take_off(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.EquipmentModule().(interface {
			ProcessOneKeyTakeOff(*equipmentpb.C2SOneKeyTakeOffProto, iface.HeroController)
		}).ProcessOneKeyTakeOff(&equipmentpb.C2SOneKeyTakeOffProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerFarm() *gm_group {
	group := &gm_group{
		tab: "farm",
		handler: []*gm_handler{
			newStringHandler("plant", " ", gm.handleFarm_c2s_plant),
			newStringHandler("harvest", " ", gm.handleFarm_c2s_harvest),
			newStringHandler("change", " ", gm.handleFarm_c2s_change),
			newStringHandler("one_key_plant", " ", gm.handleFarm_c2s_one_key_plant),
			newStringHandler("one_key_harvest", " ", gm.handleFarm_c2s_one_key_harvest),
			newStringHandler("one_key_reset", " ", gm.handleFarm_c2s_one_key_reset),
			newStringHandler("view_farm", " ", gm.handleFarm_c2s_view_farm),
			newStringHandler("steal", " ", gm.handleFarm_c2s_steal),
			newStringHandler("one_key_steal", " ", gm.handleFarm_c2s_one_key_steal),
			newStringHandler("steal_log_list", " ", gm.handleFarm_c2s_steal_log_list),
			newStringHandler("can_steal_list", " ", gm.handleFarm_c2s_can_steal_list),
		},
	}

	return group
}

func (gm *GmModule) handleFarm_c2s_plant(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessPlant(*farmpb.C2SPlantProto, iface.HeroController)
		}).ProcessPlant(&farmpb.C2SPlantProto{

			CubeX: parseInt32(strArray[0]),

			CubeY: parseInt32(strArray[1]),

			ResId: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_harvest(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessHarvest(*farmpb.C2SHarvestProto, iface.HeroController)
		}).ProcessHarvest(&farmpb.C2SHarvestProto{

			CubeX: parseInt32(strArray[0]),

			CubeY: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_change(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessChange(*farmpb.C2SChangeProto, iface.HeroController)
		}).ProcessChange(&farmpb.C2SChangeProto{

			CubeX: parseInt32(strArray[0]),

			CubeY: parseInt32(strArray[1]),

			ResId: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_one_key_plant(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessOneKeyPlant(*farmpb.C2SOneKeyPlantProto, iface.HeroController)
		}).ProcessOneKeyPlant(&farmpb.C2SOneKeyPlantProto{

			GoldConfId: parseInt32(strArray[0]),

			StoneConfId: parseInt32(strArray[1]),

			GoldCount: parseInt32(strArray[2]),

			StoneCount: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_one_key_harvest(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessOneKeyHarvest(*farmpb.C2SOneKeyHarvestProto, iface.HeroController)
		}).ProcessOneKeyHarvest(&farmpb.C2SOneKeyHarvestProto{

			ResType: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_one_key_reset(amount string, hc iface.HeroController) {
	gm.modules.FarmModule().(interface {
		ProcessOneKeyReset(iface.HeroController)
	}).ProcessOneKeyReset(hc)
}
func (gm *GmModule) handleFarm_c2s_view_farm(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessViewFarm(*farmpb.C2SViewFarmProto, iface.HeroController)
		}).ProcessViewFarm(&farmpb.C2SViewFarmProto{

			Target: parseBytes(strArray[0]),

			OpenWin: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_steal(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessSteal(*farmpb.C2SStealProto, iface.HeroController)
		}).ProcessSteal(&farmpb.C2SStealProto{

			Target: parseBytes(strArray[0]),

			CubeX: parseInt32(strArray[1]),

			CubeY: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_one_key_steal(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessOneKeySteal(*farmpb.C2SOneKeyStealProto, iface.HeroController)
		}).ProcessOneKeySteal(&farmpb.C2SOneKeyStealProto{

			Target: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_steal_log_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FarmModule().(interface {
			ProcessStealLogList(*farmpb.C2SStealLogListProto, iface.HeroController)
		}).ProcessStealLogList(&farmpb.C2SStealLogListProto{

			Target: parseBytes(strArray[0]),

			Newest: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleFarm_c2s_can_steal_list(amount string, hc iface.HeroController) {
	gm.modules.FarmModule().(interface {
		ProcessCanStealList(iface.HeroController)
	}).ProcessCanStealList(hc)
}

func (gm *GmModule) initHandlerFishing() *gm_group {
	group := &gm_group{
		tab: "fishing",
		handler: []*gm_handler{
			newStringHandler("fishing", " ", gm.handleFishing_c2s_fishing),
			newStringHandler("fish_point_exchange", " ", gm.handleFishing_c2s_fish_point_exchange),
			newStringHandler("set_fishing_captain", " ", gm.handleFishing_c2s_set_fishing_captain),
		},
	}

	return group
}

func (gm *GmModule) handleFishing_c2s_fishing(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FishingModule().(interface {
			ProcessFishing(*fishingpb.C2SFishingProto, iface.HeroController)
		}).ProcessFishing(&fishingpb.C2SFishingProto{

			Times: parseInt32(strArray[0]),

			FishType: parseInt32(strArray[1]),

			UseGoods: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleFishing_c2s_fish_point_exchange(amount string, hc iface.HeroController) {
	gm.modules.FishingModule().(interface {
		ProcessFishPointExchange(iface.HeroController)
	}).ProcessFishPointExchange(hc)
}
func (gm *GmModule) handleFishing_c2s_set_fishing_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.FishingModule().(interface {
			ProcessSetFishingCaptain(*fishingpb.C2SSetFishingCaptainProto, iface.HeroController)
		}).ProcessSetFishingCaptain(&fishingpb.C2SSetFishingCaptainProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerGarden() *gm_group {
	group := &gm_group{
		tab: "garden",
		handler: []*gm_handler{
			newStringHandler("list_treasury_tree_hero", " ", gm.handleGarden_c2s_list_treasury_tree_hero),
			newStringHandler("list_help_me", " ", gm.handleGarden_c2s_list_help_me),
			newStringHandler("list_treasury_tree_times", " ", gm.handleGarden_c2s_list_treasury_tree_times),
			newStringHandler("water_treasury_tree", " ", gm.handleGarden_c2s_water_treasury_tree),
			newStringHandler("collect_treasury_tree_prize", " ", gm.handleGarden_c2s_collect_treasury_tree_prize),
		},
	}

	return group
}

func (gm *GmModule) handleGarden_c2s_list_treasury_tree_hero(amount string, hc iface.HeroController) {
	gm.modules.GardenModule().(interface {
		ProcessListTreasuryTreeHero(iface.HeroController)
	}).ProcessListTreasuryTreeHero(hc)
}
func (gm *GmModule) handleGarden_c2s_list_help_me(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GardenModule().(interface {
			ProcessListHelpMe(*gardenpb.C2SListHelpMeProto, iface.HeroController)
		}).ProcessListHelpMe(&gardenpb.C2SListHelpMeProto{

			TargetId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGarden_c2s_list_treasury_tree_times(amount string, hc iface.HeroController) {
	gm.modules.GardenModule().(interface {
		ProcessListTreasuryTreeTimes(iface.HeroController)
	}).ProcessListTreasuryTreeTimes(hc)
}
func (gm *GmModule) handleGarden_c2s_water_treasury_tree(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GardenModule().(interface {
			ProcessWaterTreasuryTree(*gardenpb.C2SWaterTreasuryTreeProto, iface.HeroController)
		}).ProcessWaterTreasuryTree(&gardenpb.C2SWaterTreasuryTreeProto{

			Target: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGarden_c2s_collect_treasury_tree_prize(amount string, hc iface.HeroController) {
	gm.modules.GardenModule().(interface {
		ProcessCollectTreasureTreePrize(iface.HeroController)
	}).ProcessCollectTreasureTreePrize(hc)
}

func (gm *GmModule) initHandlerGem() *gm_group {
	group := &gm_group{
		tab: "gem",
		handler: []*gm_handler{
			newStringHandler("use_gem", " ", gm.handleGem_c2s_use_gem),
			newStringHandler("inlay_gem", " ", gm.handleGem_c2s_inlay_gem),
			newStringHandler("combine_gem", " ", gm.handleGem_c2s_combine_gem),
			newStringHandler("one_key_use_gem", " ", gm.handleGem_c2s_one_key_use_gem),
			newStringHandler("one_key_combine_gem", " ", gm.handleGem_c2s_one_key_combine_gem),
			newStringHandler("request_one_key_combine_cost", " ", gm.handleGem_c2s_request_one_key_combine_cost),
			newStringHandler("one_key_combine_depot_gem", " ", gm.handleGem_c2s_one_key_combine_depot_gem),
		},
	}

	return group
}

func (gm *GmModule) handleGem_c2s_use_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessUseGem(*gempb.C2SUseGemProto, iface.HeroController)
		}).ProcessUseGem(&gempb.C2SUseGemProto{

			CaptainId: parseInt32(strArray[0]),

			SlotIdx: parseInt32(strArray[1]),

			Down: parseBool(strArray[2]),

			GemId: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_inlay_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessInlayGem(*gempb.C2SInlayGemProto, iface.HeroController)
		}).ProcessInlayGem(&gempb.C2SInlayGemProto{

			CaptainId: parseInt32(strArray[0]),

			SlotIdx: parseInt32(strArray[1]),

			GemId: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_combine_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessCombineGem(*gempb.C2SCombineGemProto, iface.HeroController)
		}).ProcessCombineGem(&gempb.C2SCombineGemProto{

			CaptainId: parseInt32(strArray[0]),

			SlotIdx: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_one_key_use_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessOneKeyUseGem(*gempb.C2SOneKeyUseGemProto, iface.HeroController)
		}).ProcessOneKeyUseGem(&gempb.C2SOneKeyUseGemProto{

			CaptainId: parseInt32(strArray[0]),

			DownAll: parseBool(strArray[1]),

			EquipType: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_one_key_combine_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessOneKeyCombineGem(*gempb.C2SOneKeyCombineGemProto, iface.HeroController)
		}).ProcessOneKeyCombineGem(&gempb.C2SOneKeyCombineGemProto{

			CaptainId: parseInt32(strArray[0]),

			SlotIdx: parseInt32(strArray[1]),

			Buy: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_request_one_key_combine_cost(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessRequestOneKeyCombineGemCost(*gempb.C2SRequestOneKeyCombineCostProto, iface.HeroController)
		}).ProcessRequestOneKeyCombineGemCost(&gempb.C2SRequestOneKeyCombineCostProto{

			CaptainId: parseInt32(strArray[0]),

			SlotIdx: parseInt32(strArray[1]),

			GemId: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGem_c2s_one_key_combine_depot_gem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GemModule().(interface {
			ProcessOneKeyCombineDepotGem(*gempb.C2SOneKeyCombineDepotGemProto, iface.HeroController)
		}).ProcessOneKeyCombineDepotGem(&gempb.C2SOneKeyCombineDepotGemProto{

			GemId: parseInt32(strArray[0]),

			NewGemCount: parseInt32(strArray[1]),

			Buy: parseBool(strArray[2]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerGuild() *gm_group {
	group := &gm_group{
		tab: "guild",
		handler: []*gm_handler{
			newStringHandler("list_guild", " ", gm.handleGuild_c2s_list_guild),
			newStringHandler("search_guild", " ", gm.handleGuild_c2s_search_guild),
			newStringHandler("create_guild", " ", gm.handleGuild_c2s_create_guild),
			newStringHandler("self_guild", " ", gm.handleGuild_c2s_self_guild),
			newStringHandler("leave_guild", " ", gm.handleGuild_c2s_leave_guild),
			newStringHandler("kick_other", " ", gm.handleGuild_c2s_kick_other),
			newStringHandler("update_text", " ", gm.handleGuild_c2s_update_text),
			newStringHandler("update_internal_text", " ", gm.handleGuild_c2s_update_internal_text),
			newStringHandler("update_class_names", " ", gm.handleGuild_c2s_update_class_names),
			newStringHandler("update_class_title", " ", gm.handleGuild_c2s_update_class_title),
			newStringHandler("update_flag_type", " ", gm.handleGuild_c2s_update_flag_type),
			newStringHandler("update_member_class_level", " ", gm.handleGuild_c2s_update_member_class_level),
			newStringHandler("cancel_change_leader", " ", gm.handleGuild_c2s_cancel_change_leader),
			newStringHandler("update_join_condition", " ", gm.handleGuild_c2s_update_join_condition),
			newStringHandler("update_guild_name", " ", gm.handleGuild_c2s_update_guild_name),
			newStringHandler("update_guild_label", " ", gm.handleGuild_c2s_update_guild_label),
			newStringHandler("donate", " ", gm.handleGuild_c2s_donate),
			newStringHandler("upgrade_level", " ", gm.handleGuild_c2s_upgrade_level),
			newStringHandler("reduce_upgrade_level_cd", " ", gm.handleGuild_c2s_reduce_upgrade_level_cd),
			newStringHandler("impeach_leader", " ", gm.handleGuild_c2s_impeach_leader),
			newStringHandler("impeach_leader_vote", " ", gm.handleGuild_c2s_impeach_leader_vote),
			newStringHandler("list_guild_by_ids", " ", gm.handleGuild_c2s_list_guild_by_ids),
			newStringHandler("user_request_join", " ", gm.handleGuild_c2s_user_request_join),
			newStringHandler("user_cancel_join_request", " ", gm.handleGuild_c2s_user_cancel_join_request),
			newStringHandler("guild_reply_join_request", " ", gm.handleGuild_c2s_guild_reply_join_request),
			newStringHandler("guild_invate_other", " ", gm.handleGuild_c2s_guild_invate_other),
			newStringHandler("guild_cancel_invate_other", " ", gm.handleGuild_c2s_guild_cancel_invate_other),
			newStringHandler("user_reply_invate_request", " ", gm.handleGuild_c2s_user_reply_invate_request),
			newStringHandler("list_invite_me_guild", " ", gm.handleGuild_c2s_list_invite_me_guild),
			newStringHandler("update_friend_guild", " ", gm.handleGuild_c2s_update_friend_guild),
			newStringHandler("update_enemy_guild", " ", gm.handleGuild_c2s_update_enemy_guild),
			newStringHandler("update_guild_prestige", " ", gm.handleGuild_c2s_update_guild_prestige),
			newStringHandler("place_guild_statue", " ", gm.handleGuild_c2s_place_guild_statue),
			newStringHandler("take_back_guild_statue", " ", gm.handleGuild_c2s_take_back_guild_statue),
			newStringHandler("collect_first_join_guild_prize", " ", gm.handleGuild_c2s_collect_first_join_guild_prize),
			newStringHandler("seek_help", " ", gm.handleGuild_c2s_seek_help),
			newStringHandler("help_guild_member", " ", gm.handleGuild_c2s_help_guild_member),
			newStringHandler("help_all_guild_member", " ", gm.handleGuild_c2s_help_all_guild_member),
			newStringHandler("collect_guild_event_prize", " ", gm.handleGuild_c2s_collect_guild_event_prize),
			newStringHandler("collect_full_big_box", " ", gm.handleGuild_c2s_collect_full_big_box),
			newStringHandler("upgrade_technology", " ", gm.handleGuild_c2s_upgrade_technology),
			newStringHandler("reduce_technology_cd", " ", gm.handleGuild_c2s_reduce_technology_cd),
			newStringHandler("list_guild_logs", " ", gm.handleGuild_c2s_list_guild_logs),
			newStringHandler("request_recommend_guild", " ", gm.handleGuild_c2s_request_recommend_guild),
			newStringHandler("help_tech", " ", gm.handleGuild_c2s_help_tech),
			newStringHandler("recommend_invite_heros", " ", gm.handleGuild_c2s_recommend_invite_heros),
			newStringHandler("search_no_guild_heros", " ", gm.handleGuild_c2s_search_no_guild_heros),
			newStringHandler("view_mc_war_record", " ", gm.handleGuild_c2s_view_mc_war_record),
			newStringHandler("update_guild_mark", " ", gm.handleGuild_c2s_update_guild_mark),
			newStringHandler("view_yinliang_record", " ", gm.handleGuild_c2s_view_yinliang_record),
			newStringHandler("send_yinliang_to_other_guild", " ", gm.handleGuild_c2s_send_yinliang_to_other_guild),
			newStringHandler("send_yinliang_to_member", " ", gm.handleGuild_c2s_send_yinliang_to_member),
			newStringHandler("pay_salary", " ", gm.handleGuild_c2s_pay_salary),
			newStringHandler("set_salary", " ", gm.handleGuild_c2s_set_salary),
			newStringHandler("view_send_yinliang_to_guild", " ", gm.handleGuild_c2s_view_send_yinliang_to_guild),
			newStringHandler("convene", " ", gm.handleGuild_c2s_convene),
			newStringHandler("collect_daily_guild_rank_prize", " ", gm.handleGuild_c2s_collect_daily_guild_rank_prize),
			newStringHandler("view_daily_guild_rank", " ", gm.handleGuild_c2s_view_daily_guild_rank),
			newStringHandler("add_recommend_mc_build", " ", gm.handleGuild_c2s_add_recommend_mc_build),
			newStringHandler("view_task_progress", " ", gm.handleGuild_c2s_view_task_progress),
			newStringHandler("collect_task_prize", " ", gm.handleGuild_c2s_collect_task_prize),
			newStringHandler("guild_change_country", " ", gm.handleGuild_c2s_guild_change_country),
			newStringHandler("cancel_guild_change_country", " ", gm.handleGuild_c2s_cancel_guild_change_country),
			newStringHandler("show_workshop_not_exist", " ", gm.handleGuild_c2s_show_workshop_not_exist),
		},
	}

	return group
}

func (gm *GmModule) handleGuild_c2s_list_guild(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessListGuild(iface.HeroController)
	}).ProcessListGuild(hc)
}
func (gm *GmModule) handleGuild_c2s_search_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSearchGuild(*guildpb.C2SSearchGuildProto, iface.HeroController)
		}).ProcessSearchGuild(&guildpb.C2SSearchGuildProto{

			Name: parseString(strArray[0]),

			Num: parseInt32(strArray[1]),

			ShowSelfGuild: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_create_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessCreateGuild(*guildpb.C2SCreateGuildProto, iface.HeroController)
		}).ProcessCreateGuild(&guildpb.C2SCreateGuildProto{

			Name: parseString(strArray[0]),

			FlagName: parseString(strArray[1]),

			Country: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_self_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSelfGuild(*guildpb.C2SSelfGuildProto, iface.HeroController)
		}).ProcessSelfGuild(&guildpb.C2SSelfGuildProto{

			Version: parseInt32(strArray[0]),

			GuildId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_leave_guild(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessLeaveGuild(iface.HeroController)
	}).ProcessLeaveGuild(hc)
}
func (gm *GmModule) handleGuild_c2s_kick_other(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessKickOther(*guildpb.C2SKickOtherProto, iface.HeroController)
		}).ProcessKickOther(&guildpb.C2SKickOtherProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_text(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateText(*guildpb.C2SUpdateTextProto, iface.HeroController)
		}).ProcessUpdateText(&guildpb.C2SUpdateTextProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_internal_text(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateInternalText(*guildpb.C2SUpdateInternalTextProto, iface.HeroController)
		}).ProcessUpdateInternalText(&guildpb.C2SUpdateInternalTextProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_class_names(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateClassNames(*guildpb.C2SUpdateClassNamesProto, iface.HeroController)
		}).ProcessUpdateClassNames(&guildpb.C2SUpdateClassNamesProto{

			Name: parseStringArray(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_class_title(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateClassTitle(*guildpb.C2SUpdateClassTitleProto, iface.HeroController)
		}).ProcessUpdateClassTitle(&guildpb.C2SUpdateClassTitleProto{

			Proto: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_flag_type(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateFlagType(*guildpb.C2SUpdateFlagTypeProto, iface.HeroController)
		}).ProcessUpdateFlagType(&guildpb.C2SUpdateFlagTypeProto{

			FlagType: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_member_class_level(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateMemberClassLevel(*guildpb.C2SUpdateMemberClassLevelProto, iface.HeroController)
		}).ProcessUpdateMemberClassLevel(&guildpb.C2SUpdateMemberClassLevelProto{

			Id: parseBytes(strArray[0]),

			ClassLevel: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_cancel_change_leader(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessCancelChangeLeader(iface.HeroController)
	}).ProcessCancelChangeLeader(hc)
}
func (gm *GmModule) handleGuild_c2s_update_join_condition(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateJoinCondition(*guildpb.C2SUpdateJoinConditionProto, iface.HeroController)
		}).ProcessUpdateJoinCondition(&guildpb.C2SUpdateJoinConditionProto{

			RejectAutoJoin: parseBool(strArray[0]),

			RequiredHeroLevel: parseInt32(strArray[1]),

			RequiredJunXianLevel: parseInt32(strArray[2]),

			RequiredTowerMaxFloor: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_guild_name(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateGuildName(*guildpb.C2SUpdateGuildNameProto, iface.HeroController)
		}).ProcessUpdateGuildName(&guildpb.C2SUpdateGuildNameProto{

			Name: parseString(strArray[0]),

			FlagName: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_guild_label(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateLabels(*guildpb.C2SUpdateGuildLabelProto, iface.HeroController)
		}).ProcessUpdateLabels(&guildpb.C2SUpdateGuildLabelProto{

			Label: parseStringArray(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_donate(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessDonation(*guildpb.C2SDonateProto, iface.HeroController)
		}).ProcessDonation(&guildpb.C2SDonateProto{

			Sequence: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_upgrade_level(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessUpgradeLevel(iface.HeroController)
	}).ProcessUpgradeLevel(hc)
}
func (gm *GmModule) handleGuild_c2s_reduce_upgrade_level_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessReduceUpgradeLevelCd(*guildpb.C2SReduceUpgradeLevelCdProto, iface.HeroController)
		}).ProcessReduceUpgradeLevelCd(&guildpb.C2SReduceUpgradeLevelCdProto{

			Times: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_impeach_leader(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessImpeachLeader(iface.HeroController)
	}).ProcessImpeachLeader(hc)
}
func (gm *GmModule) handleGuild_c2s_impeach_leader_vote(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessImpeachLeaderVote(*guildpb.C2SImpeachLeaderVoteProto, iface.HeroController)
		}).ProcessImpeachLeaderVote(&guildpb.C2SImpeachLeaderVoteProto{

			Target: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_list_guild_by_ids(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessListGuildByIds(*guildpb.C2SListGuildByIdsProto, iface.HeroController)
		}).ProcessListGuildByIds(&guildpb.C2SListGuildByIdsProto{

			Ids: parseInt32Array(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_user_request_join(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUserRequestJoin(*guildpb.C2SUserRequestJoinProto, iface.HeroController)
		}).ProcessUserRequestJoin(&guildpb.C2SUserRequestJoinProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_user_cancel_join_request(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUserCancelJoinRequest(*guildpb.C2SUserCancelJoinRequestProto, iface.HeroController)
		}).ProcessUserCancelJoinRequest(&guildpb.C2SUserCancelJoinRequestProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_guild_reply_join_request(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessGuildReplyJoinRequest(*guildpb.C2SGuildReplyJoinRequestProto, iface.HeroController)
		}).ProcessGuildReplyJoinRequest(&guildpb.C2SGuildReplyJoinRequestProto{

			Id: parseBytes(strArray[0]),

			Agree: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_guild_invate_other(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessInvateOtherRequest(*guildpb.C2SGuildInvateOtherProto, iface.HeroController)
		}).ProcessInvateOtherRequest(&guildpb.C2SGuildInvateOtherProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_guild_cancel_invate_other(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessCancelInvateOtherRequest(*guildpb.C2SGuildCancelInvateOtherProto, iface.HeroController)
		}).ProcessCancelInvateOtherRequest(&guildpb.C2SGuildCancelInvateOtherProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_user_reply_invate_request(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUserReplyInvateRequest(*guildpb.C2SUserReplyInvateRequestProto, iface.HeroController)
		}).ProcessUserReplyInvateRequest(&guildpb.C2SUserReplyInvateRequestProto{

			Id: parseInt32(strArray[0]),

			Agree: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_list_invite_me_guild(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessListInviteMeGuild(iface.HeroController)
	}).ProcessListInviteMeGuild(hc)
}
func (gm *GmModule) handleGuild_c2s_update_friend_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateFriendGuild(*guildpb.C2SUpdateFriendGuildProto, iface.HeroController)
		}).ProcessUpdateFriendGuild(&guildpb.C2SUpdateFriendGuildProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_enemy_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateEnemyGuild(*guildpb.C2SUpdateEnemyGuildProto, iface.HeroController)
		}).ProcessUpdateEnemyGuild(&guildpb.C2SUpdateEnemyGuildProto{

			Text: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_update_guild_prestige(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateGuildPrestige(*guildpb.C2SUpdateGuildPrestigeProto, iface.HeroController)
		}).ProcessUpdateGuildPrestige(&guildpb.C2SUpdateGuildPrestigeProto{

			Target: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_place_guild_statue(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessPlaceGuildStatue(*guildpb.C2SPlaceGuildStatueProto, iface.HeroController)
		}).ProcessPlaceGuildStatue(&guildpb.C2SPlaceGuildStatueProto{

			RealmId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_take_back_guild_statue(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessTakeBackGuildStatue(iface.HeroController)
	}).ProcessTakeBackGuildStatue(hc)
}
func (gm *GmModule) handleGuild_c2s_collect_first_join_guild_prize(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessCollectFirstJoinGuildPrize(iface.HeroController)
	}).ProcessCollectFirstJoinGuildPrize(hc)
}
func (gm *GmModule) handleGuild_c2s_seek_help(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSeekHelp(*guildpb.C2SSeekHelpProto, iface.HeroController)
		}).ProcessSeekHelp(&guildpb.C2SSeekHelpProto{

			HelpType: parseInt32(strArray[0]),

			WorkerPos: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_help_guild_member(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessHelpGuildMember(*guildpb.C2SHelpGuildMemberProto, iface.HeroController)
		}).ProcessHelpGuildMember(&guildpb.C2SHelpGuildMemberProto{

			Id: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_help_all_guild_member(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessHelpAllGuildMember(iface.HeroController)
	}).ProcessHelpAllGuildMember(hc)
}
func (gm *GmModule) handleGuild_c2s_collect_guild_event_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessCollectGuildEventPrize(*guildpb.C2SCollectGuildEventPrizeProto, iface.HeroController)
		}).ProcessCollectGuildEventPrize(&guildpb.C2SCollectGuildEventPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_collect_full_big_box(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessCollectFullBigBox(iface.HeroController)
	}).ProcessCollectFullBigBox(hc)
}
func (gm *GmModule) handleGuild_c2s_upgrade_technology(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpgradeTechnology(*guildpb.C2SUpgradeTechnologyProto, iface.HeroController)
		}).ProcessUpgradeTechnology(&guildpb.C2SUpgradeTechnologyProto{

			Group: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_reduce_technology_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessReduceTechnologyCd(*guildpb.C2SReduceTechnologyCdProto, iface.HeroController)
		}).ProcessReduceTechnologyCd(&guildpb.C2SReduceTechnologyCdProto{

			Times: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_list_guild_logs(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessListGuildLogs(*guildpb.C2SListGuildLogsProto, iface.HeroController)
		}).ProcessListGuildLogs(&guildpb.C2SListGuildLogsProto{

			LogType: parseInt32(strArray[0]),

			MinId: parseInt32(strArray[1]),

			Count: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_request_recommend_guild(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessRequestRecommendGuild(iface.HeroController)
	}).ProcessRequestRecommendGuild(hc)
}
func (gm *GmModule) handleGuild_c2s_help_tech(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessHelpTech(iface.HeroController)
	}).ProcessHelpTech(hc)
}
func (gm *GmModule) handleGuild_c2s_recommend_invite_heros(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessRecommendInviteHeros(iface.HeroController)
	}).ProcessRecommendInviteHeros(hc)
}
func (gm *GmModule) handleGuild_c2s_search_no_guild_heros(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSearchNoGuildHeros(*guildpb.C2SSearchNoGuildHerosProto, iface.HeroController)
		}).ProcessSearchNoGuildHeros(&guildpb.C2SSearchNoGuildHerosProto{

			Name: parseString(strArray[0]),

			Page: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_view_mc_war_record(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessViewMcWarRecord(iface.HeroController)
	}).ProcessViewMcWarRecord(hc)
}
func (gm *GmModule) handleGuild_c2s_update_guild_mark(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessUpdateGuildMark(*guildpb.C2SUpdateGuildMarkProto, iface.HeroController)
		}).ProcessUpdateGuildMark(&guildpb.C2SUpdateGuildMarkProto{

			Index: parseInt32(strArray[0]),

			PosX: parseInt32(strArray[1]),

			PosY: parseInt32(strArray[2]),

			Msg: parseString(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_view_yinliang_record(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessViewYinliangRecord(iface.HeroController)
	}).ProcessViewYinliangRecord(hc)
}
func (gm *GmModule) handleGuild_c2s_send_yinliang_to_other_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSendYinliangToOtherGuild(*guildpb.C2SSendYinliangToOtherGuildProto, iface.HeroController)
		}).ProcessSendYinliangToOtherGuild(&guildpb.C2SSendYinliangToOtherGuildProto{

			Gid: parseInt32(strArray[0]),

			Amount: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_send_yinliang_to_member(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSendYinliangToMember(*guildpb.C2SSendYinliangToMemberProto, iface.HeroController)
		}).ProcessSendYinliangToMember(&guildpb.C2SSendYinliangToMemberProto{

			MemId: parseBytes(strArray[0]),

			Amount: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_pay_salary(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessPaySalary(iface.HeroController)
	}).ProcessPaySalary(hc)
}
func (gm *GmModule) handleGuild_c2s_set_salary(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessSetSalary(*guildpb.C2SSetSalaryProto, iface.HeroController)
		}).ProcessSetSalary(&guildpb.C2SSetSalaryProto{

			MemId: parseBytes(strArray[0]),

			Salary: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_view_send_yinliang_to_guild(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessViewSendYinliangToGuild(iface.HeroController)
	}).ProcessViewSendYinliangToGuild(hc)
}
func (gm *GmModule) handleGuild_c2s_convene(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessConvene(*guildpb.C2SConveneProto, iface.HeroController)
		}).ProcessConvene(&guildpb.C2SConveneProto{

			Target: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_collect_daily_guild_rank_prize(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessCollectDailyGuildRankPrize(iface.HeroController)
	}).ProcessCollectDailyGuildRankPrize(hc)
}
func (gm *GmModule) handleGuild_c2s_view_daily_guild_rank(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessViewDailyGuildRank(iface.HeroController)
	}).ProcessViewDailyGuildRank(hc)
}
func (gm *GmModule) handleGuild_c2s_add_recommend_mc_build(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessAddRecommendMcBuild(*guildpb.C2SAddRecommendMcBuildProto, iface.HeroController)
		}).ProcessAddRecommendMcBuild(&guildpb.C2SAddRecommendMcBuildProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_view_task_progress(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessViewTaskProgress(*guildpb.C2SViewTaskProgressProto, iface.HeroController)
		}).ProcessViewTaskProgress(&guildpb.C2SViewTaskProgressProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_collect_task_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessCollectTaskPrizeProgress(*guildpb.C2SCollectTaskPrizeProto, iface.HeroController)
		}).ProcessCollectTaskPrizeProgress(&guildpb.C2SCollectTaskPrizeProto{

			TaskId: parseInt32(strArray[0]),

			Stage: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_guild_change_country(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.GuildModule().(interface {
			ProcessGuildChangeCountry(*guildpb.C2SGuildChangeCountryProto, iface.HeroController)
		}).ProcessGuildChangeCountry(&guildpb.C2SGuildChangeCountryProto{

			Country: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleGuild_c2s_cancel_guild_change_country(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessCancelGuildChangeCountry(iface.HeroController)
	}).ProcessCancelGuildChangeCountry(hc)
}
func (gm *GmModule) handleGuild_c2s_show_workshop_not_exist(amount string, hc iface.HeroController) {
	gm.modules.GuildModule().(interface {
		ProcessShowWorkshopNotExist(iface.HeroController)
	}).ProcessShowWorkshopNotExist(hc)
}

func (gm *GmModule) initHandlerHebi() *gm_group {
	group := &gm_group{
		tab: "hebi",
		handler: []*gm_handler{
			newStringHandler("room_list", " ", gm.handleHebi_c2s_room_list),
			newStringHandler("hero_record_list", " ", gm.handleHebi_c2s_hero_record_list),
			newStringHandler("change_captain", " ", gm.handleHebi_c2s_change_captain),
			newStringHandler("check_in_room", " ", gm.handleHebi_c2s_check_in_room),
			newStringHandler("copy_self", " ", gm.handleHebi_c2s_copy_self),
			newStringHandler("join_room", " ", gm.handleHebi_c2s_join_room),
			newStringHandler("rob_pos", " ", gm.handleHebi_c2s_rob_pos),
			newStringHandler("leave_room", " ", gm.handleHebi_c2s_leave_room),
			newStringHandler("rob", " ", gm.handleHebi_c2s_rob),
			newStringHandler("view_show_prize", " ", gm.handleHebi_c2s_view_show_prize),
		},
	}

	return group
}

func (gm *GmModule) handleHebi_c2s_room_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessRoomList(*hebipb.C2SRoomListProto, iface.HeroController)
		}).ProcessRoomList(&hebipb.C2SRoomListProto{

			V: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_hero_record_list(amount string, hc iface.HeroController) {
	gm.modules.HebiModule().(interface {
		ProcessHebiHeroRecordList(iface.HeroController)
	}).ProcessHebiHeroRecordList(hc)
}
func (gm *GmModule) handleHebi_c2s_change_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessChangeCaptain(*hebipb.C2SChangeCaptainProto, iface.HeroController)
		}).ProcessChangeCaptain(&hebipb.C2SChangeCaptainProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_check_in_room(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessCheckInRoom(*hebipb.C2SCheckInRoomProto, iface.HeroController)
		}).ProcessCheckInRoom(&hebipb.C2SCheckInRoomProto{

			RoomId: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_copy_self(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessCopySelf(*hebipb.C2SCopySelfProto, iface.HeroController)
		}).ProcessCopySelf(&hebipb.C2SCopySelfProto{

			RoomId: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_join_room(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessJoinRoom(*hebipb.C2SJoinRoomProto, iface.HeroController)
		}).ProcessJoinRoom(&hebipb.C2SJoinRoomProto{

			RoomId: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_rob_pos(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessRobPos(*hebipb.C2SRobPosProto, iface.HeroController)
		}).ProcessRobPos(&hebipb.C2SRobPosProto{

			RoomId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_leave_room(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessLeave(*hebipb.C2SLeaveRoomProto, iface.HeroController)
		}).ProcessLeave(&hebipb.C2SLeaveRoomProto{

			RoomId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_rob(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessRob(*hebipb.C2SRobProto, iface.HeroController)
		}).ProcessRob(&hebipb.C2SRobProto{

			RoomId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleHebi_c2s_view_show_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.HebiModule().(interface {
			ProcessViewHebiShowPrize(*hebipb.C2SViewShowPrizeProto, iface.HeroController)
		}).ProcessViewHebiShowPrize(&hebipb.C2SViewShowPrizeProto{

			HeroLevel: parseInt32(strArray[0]),

			Goods: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerMail() *gm_group {
	group := &gm_group{
		tab: "mail",
		handler: []*gm_handler{
			newStringHandler("list_mail", " ", gm.handleMail_c2s_list_mail),
			newStringHandler("delete_mail", " ", gm.handleMail_c2s_delete_mail),
			newStringHandler("keep_mail", " ", gm.handleMail_c2s_keep_mail),
			newStringHandler("collect_mail_prize", " ", gm.handleMail_c2s_collect_mail_prize),
			newStringHandler("read_mail", " ", gm.handleMail_c2s_read_mail),
			newStringHandler("read_multi", " ", gm.handleMail_c2s_read_multi),
			newStringHandler("delete_multi", " ", gm.handleMail_c2s_delete_multi),
			newStringHandler("get_mail", " ", gm.handleMail_c2s_get_mail),
		},
	}

	return group
}

func (gm *GmModule) handleMail_c2s_list_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ListMail(*mailpb.C2SListMailProto, iface.HeroController)
		}).ListMail(&mailpb.C2SListMailProto{

			Read: parseInt32(strArray[0]),

			Keep: parseInt32(strArray[1]),

			Report: parseInt32(strArray[2]),

			ReportTag: parseInt32(strArray[3]),

			HasPrize: parseInt32(strArray[4]),

			Collected: parseInt32(strArray[5]),

			MinId: parseBytes(strArray[6]),

			Count: parseInt32(strArray[7]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_delete_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			DeleteMail(*mailpb.C2SDeleteMailProto, iface.HeroController)
		}).DeleteMail(&mailpb.C2SDeleteMailProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_keep_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			KeepMail(*mailpb.C2SKeepMailProto, iface.HeroController)
		}).KeepMail(&mailpb.C2SKeepMailProto{

			Id: parseBytes(strArray[0]),

			Keep: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_collect_mail_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ProcessCollectMailPrize(*mailpb.C2SCollectMailPrizeProto, iface.HeroController)
		}).ProcessCollectMailPrize(&mailpb.C2SCollectMailPrizeProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_read_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ReadMail(*mailpb.C2SReadMailProto, iface.HeroController)
		}).ReadMail(&mailpb.C2SReadMailProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_read_multi(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ProcessReadMulti(*mailpb.C2SReadMultiProto, iface.HeroController)
		}).ProcessReadMulti(&mailpb.C2SReadMultiProto{

			Ids: parseBytesArray(strArray[0]),

			Report: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_delete_multi(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ProcessDeleteMulti(*mailpb.C2SDeleteMultiProto, iface.HeroController)
		}).ProcessDeleteMulti(&mailpb.C2SDeleteMultiProto{

			Ids: parseBytesArray(strArray[0]),

			Report: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMail_c2s_get_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MailModule().(interface {
			ProcessGetMail(*mailpb.C2SGetMailProto, iface.HeroController)
		}).ProcessGetMail(&mailpb.C2SGetMailProto{

			Bid: parseBytes(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerMilitary() *gm_group {
	group := &gm_group{
		tab: "military",
		handler: []*gm_handler{
			newStringHandler("recruit_soldier", " ", gm.handleMilitary_c2s_recruit_soldier),
			newStringHandler("recruit_soldier_v2", " ", gm.handleMilitary_c2s_recruit_soldier_v2),
			newStringHandler("heal_wounded_soldier", " ", gm.handleMilitary_c2s_heal_wounded_soldier),
			newStringHandler("captain_change_soldier", " ", gm.handleMilitary_c2s_captain_change_soldier),
			newStringHandler("captain_full_soldier", " ", gm.handleMilitary_c2s_captain_full_soldier),
			newStringHandler("force_add_soldier", " ", gm.handleMilitary_c2s_force_add_soldier),
			newStringHandler("fight", " ", gm.handleMilitary_c2s_fight),
			newStringHandler("multi_fight", " ", gm.handleMilitary_c2s_multi_fight),
			newStringHandler("fightx", " ", gm.handleMilitary_c2s_fightx),
			newStringHandler("upgrade_soldier_level", " ", gm.handleMilitary_c2s_upgrade_soldier_level),
			newStringHandler("recruit_captain_v2", " ", gm.handleMilitary_c2s_recruit_captain_v2),
			newStringHandler("random_captain_head", " ", gm.handleMilitary_c2s_random_captain_head),
			newStringHandler("recruit_captain_seeker", " ", gm.handleMilitary_c2s_recruit_captain_seeker),
			newStringHandler("set_defense_troop", " ", gm.handleMilitary_c2s_set_defense_troop),
			newStringHandler("clear_defense_troop_defeated_mail", " ", gm.handleMilitary_c2s_clear_defense_troop_defeated_mail),
			newStringHandler("set_defenser_auto_full_soldier", " ", gm.handleMilitary_c2s_set_defenser_auto_full_soldier),
			newStringHandler("use_copy_defenser_goods", " ", gm.handleMilitary_c2s_use_copy_defenser_goods),
			newStringHandler("sell_seek_captain", " ", gm.handleMilitary_c2s_sell_seek_captain),
			newStringHandler("set_multi_captain_index", " ", gm.handleMilitary_c2s_set_multi_captain_index),
			newStringHandler("set_pve_captain", " ", gm.handleMilitary_c2s_set_pve_captain),
			newStringHandler("fire_captain", " ", gm.handleMilitary_c2s_fire_captain),
			newStringHandler("captain_refined", " ", gm.handleMilitary_c2s_captain_refined),
			newStringHandler("captain_enhance", " ", gm.handleMilitary_c2s_captain_enhance),
			newStringHandler("change_captain_name", " ", gm.handleMilitary_c2s_change_captain_name),
			newStringHandler("change_captain_race", " ", gm.handleMilitary_c2s_change_captain_race),
			newStringHandler("captain_rebirth_preview", " ", gm.handleMilitary_c2s_captain_rebirth_preview),
			newStringHandler("captain_rebirth", " ", gm.handleMilitary_c2s_captain_rebirth),
			newStringHandler("captain_progress", " ", gm.handleMilitary_c2s_captain_progress),
			newStringHandler("captain_rebirth_miao_cd", " ", gm.handleMilitary_c2s_captain_rebirth_miao_cd),
			newStringHandler("collect_captain_training_exp", " ", gm.handleMilitary_c2s_collect_captain_training_exp),
			newStringHandler("captain_train_exp", " ", gm.handleMilitary_c2s_captain_train_exp),
			newStringHandler("captain_can_collect_exp", " ", gm.handleMilitary_c2s_captain_can_collect_exp),
			newStringHandler("use_training_exp_goods", " ", gm.handleMilitary_c2s_use_training_exp_goods),
			newStringHandler("use_level_exp_goods", " ", gm.handleMilitary_c2s_use_level_exp_goods),
			newStringHandler("use_level_exp_goods2", " ", gm.handleMilitary_c2s_use_level_exp_goods2),
			newStringHandler("auto_use_goods_until_captain_levelup", " ", gm.handleMilitary_c2s_auto_use_goods_until_captain_levelup),
			newStringHandler("get_max_recruit_soldier", " ", gm.handleMilitary_c2s_get_max_recruit_soldier),
			newStringHandler("get_max_heal_soldier", " ", gm.handleMilitary_c2s_get_max_heal_soldier),
			newStringHandler("jiu_guan_consult", " ", gm.handleMilitary_c2s_jiu_guan_consult),
			newStringHandler("jiu_guan_refresh", " ", gm.handleMilitary_c2s_jiu_guan_refresh),
			newStringHandler("unlock_captain_restraint_spell", " ", gm.handleMilitary_c2s_unlock_captain_restraint_spell),
			newStringHandler("get_captain_stat_details", " ", gm.handleMilitary_c2s_get_captain_stat_details),
			newStringHandler("captain_stat_details", " ", gm.handleMilitary_c2s_captain_stat_details),
			newStringHandler("update_captain_official", " ", gm.handleMilitary_c2s_update_captain_official),
			newStringHandler("set_captain_official", " ", gm.handleMilitary_c2s_set_captain_official),
			newStringHandler("leave_captain_official", " ", gm.handleMilitary_c2s_leave_captain_official),
			newStringHandler("use_gong_xun_goods", " ", gm.handleMilitary_c2s_use_gong_xun_goods),
			newStringHandler("use_gongxun_goods", " ", gm.handleMilitary_c2s_use_gongxun_goods),
			newStringHandler("close_fight_guide", " ", gm.handleMilitary_c2s_close_fight_guide),
			newStringHandler("view_other_hero_captain", " ", gm.handleMilitary_c2s_view_other_hero_captain),
			newStringHandler("captain_born", " ", gm.handleMilitary_c2s_captain_born),
			newStringHandler("captain_upstar", " ", gm.handleMilitary_c2s_captain_upstar),
			newStringHandler("captain_exchange", " ", gm.handleMilitary_c2s_captain_exchange),
			newStringHandler("notice_captain_has_viewed", " ", gm.handleMilitary_c2s_notice_captain_has_viewed),
			newStringHandler("activate_captain_friendship", " ", gm.handleMilitary_c2s_activate_captain_friendship),
			newStringHandler("notice_official_has_viewed", " ", gm.handleMilitary_c2s_notice_official_has_viewed),
		},
	}

	return group
}

func (gm *GmModule) handleMilitary_c2s_recruit_soldier(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SRecruitSoldierMsg_deprecated(*militarypb.C2SRecruitSoldierProto, iface.HeroController)
		}).ProcessC2SRecruitSoldierMsg_deprecated(&militarypb.C2SRecruitSoldierProto{

			Count: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_recruit_soldier_v2(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SRecruitSoldierV2Msg_deprecated(*militarypb.C2SRecruitSoldierV2Proto, iface.HeroController)
		}).ProcessC2SRecruitSoldierV2Msg_deprecated(&militarypb.C2SRecruitSoldierV2Proto{

			All: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_heal_wounded_soldier(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SHealWoundedSoldierMsg_deprecated(*militarypb.C2SHealWoundedSoldierProto, iface.HeroController)
		}).ProcessC2SHealWoundedSoldierMsg_deprecated(&militarypb.C2SHealWoundedSoldierProto{

			Count: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_change_soldier(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SCaptainChangeSoldierMsg(*militarypb.C2SCaptainChangeSoldierProto, iface.HeroController)
		}).ProcessC2SCaptainChangeSoldierMsg(&militarypb.C2SCaptainChangeSoldierProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_full_soldier(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SCaptainFullSoldierMsg(*militarypb.C2SCaptainFullSoldierProto, iface.HeroController)
		}).ProcessC2SCaptainFullSoldierMsg(&militarypb.C2SCaptainFullSoldierProto{

			Id: parseInt32Array(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_force_add_soldier(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessForceAddSoldierMsg(iface.HeroController)
	}).ProcessForceAddSoldierMsg(hc)
}
func (gm *GmModule) handleMilitary_c2s_fight(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SFightMsg(*militarypb.C2SFightProto, iface.HeroController)
		}).ProcessC2SFightMsg(&militarypb.C2SFightProto{

			Wall: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_multi_fight(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessC2SMultiFightMsg(iface.HeroController)
	}).ProcessC2SMultiFightMsg(hc)
}
func (gm *GmModule) handleMilitary_c2s_fightx(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessC2SFightxMsg(*militarypb.C2SFightxProto, iface.HeroController)
		}).ProcessC2SFightxMsg(&militarypb.C2SFightxProto{

			Attacker: parseInt32Array(strArray[0]),

			Defenser: parseInt32Array(strArray[1]),

			Wall: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_upgrade_soldier_level(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessC2SUpgradeSoldierMsg(iface.HeroController)
	}).ProcessC2SUpgradeSoldierMsg(hc)
}
func (gm *GmModule) handleMilitary_c2s_recruit_captain_v2(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessRecruitCaptainV2_deprecated(iface.HeroController)
	}).ProcessRecruitCaptainV2_deprecated(hc)
}
func (gm *GmModule) handleMilitary_c2s_random_captain_head(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessRandomCaptainHead_deprecated(iface.HeroController)
	}).ProcessRandomCaptainHead_deprecated(hc)
}
func (gm *GmModule) handleMilitary_c2s_recruit_captain_seeker(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessRecruitCaptainSeeker_deprecated(*militarypb.C2SRecruitCaptainSeekerProto, iface.HeroController)
		}).ProcessRecruitCaptainSeeker_deprecated(&militarypb.C2SRecruitCaptainSeekerProto{

			Index: parseInt32(strArray[0]),

			Head: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_set_defense_troop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSetDefenseTroop(*militarypb.C2SSetDefenseTroopProto, iface.HeroController)
		}).ProcessSetDefenseTroop(&militarypb.C2SSetDefenseTroopProto{

			IsTent: parseBool(strArray[0]),

			TroopIndex: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_clear_defense_troop_defeated_mail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessClearDefenseTroopDefeatedMail(*militarypb.C2SClearDefenseTroopDefeatedMailProto, iface.HeroController)
		}).ProcessClearDefenseTroopDefeatedMail(&militarypb.C2SClearDefenseTroopDefeatedMailProto{

			IsTent: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_set_defenser_auto_full_soldier(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSetDefenserAutoFullSoldier(*militarypb.C2SSetDefenserAutoFullSoldierProto, iface.HeroController)
		}).ProcessSetDefenserAutoFullSoldier(&militarypb.C2SSetDefenserAutoFullSoldierProto{

			Dont: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_use_copy_defenser_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCopyDefenserGoods(*militarypb.C2SUseCopyDefenserGoodsProto, iface.HeroController)
		}).ProcessCopyDefenserGoods(&militarypb.C2SUseCopyDefenserGoodsProto{

			Goods: parseInt32(strArray[0]),

			AutoBuy: parseBool(strArray[1]),

			TroopIndex: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_sell_seek_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSellSeekCaptain_deprecated(*militarypb.C2SSellSeekCaptainProto, iface.HeroController)
		}).ProcessSellSeekCaptain_deprecated(&militarypb.C2SSellSeekCaptainProto{

			Index: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_set_multi_captain_index(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSetMultiCaptainIndex(*militarypb.C2SSetMultiCaptainIndexProto, iface.HeroController)
		}).ProcessSetMultiCaptainIndex(&militarypb.C2SSetMultiCaptainIndexProto{

			Index: parseInt32(strArray[0]),

			Id: parseInt32Array(strArray[1]),

			XIndex: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_set_pve_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSetPveCaptain(*militarypb.C2SSetPveCaptainProto, iface.HeroController)
		}).ProcessSetPveCaptain(&militarypb.C2SSetPveCaptainProto{

			PveType: parseInt32(strArray[0]),

			Id: parseInt32Array(strArray[1]),

			XIndex: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_fire_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessFireCaptain_deprecated(*militarypb.C2SFireCaptainProto, iface.HeroController)
		}).ProcessFireCaptain_deprecated(&militarypb.C2SFireCaptainProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_refined(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainRefined_deprecated(*militarypb.C2SCaptainRefinedProto, iface.HeroController)
		}).ProcessCaptainRefined_deprecated(&militarypb.C2SCaptainRefinedProto{

			Captain: parseInt32(strArray[0]),

			GoodsId: parseInt32Array(strArray[1]),

			Count: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_enhance(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainEnhance(*militarypb.C2SCaptainEnhanceProto, iface.HeroController)
		}).ProcessCaptainEnhance(&militarypb.C2SCaptainEnhanceProto{

			Captain: parseInt32(strArray[0]),

			GoodsId: parseInt32Array(strArray[1]),

			Count: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_change_captain_name(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessChangeCaptainName_deprecated(*militarypb.C2SChangeCaptainNameProto, iface.HeroController)
		}).ProcessChangeCaptainName_deprecated(&militarypb.C2SChangeCaptainNameProto{

			Id: parseInt32(strArray[0]),

			Name: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_change_captain_race(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessChangeCaptainRace_deprecated(*militarypb.C2SChangeCaptainRaceProto, iface.HeroController)
		}).ProcessChangeCaptainRace_deprecated(&militarypb.C2SChangeCaptainRaceProto{

			Id: parseInt32(strArray[0]),

			Race: parseInt32(strArray[1]),

			Money: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_rebirth_preview(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainRevirthPreview(*militarypb.C2SCaptainRebirthPreviewProto, iface.HeroController)
		}).ProcessCaptainRevirthPreview(&militarypb.C2SCaptainRebirthPreviewProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_rebirth(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainRebirth_deprecated(*militarypb.C2SCaptainRebirthProto, iface.HeroController)
		}).ProcessCaptainRebirth_deprecated(&militarypb.C2SCaptainRebirthProto{

			Id: parseInt32(strArray[0]),

			Miao: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_progress(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainProgress(*militarypb.C2SCaptainProgressProto, iface.HeroController)
		}).ProcessCaptainProgress(&militarypb.C2SCaptainProgressProto{

			Id: parseInt32(strArray[0]),

			Miao: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_rebirth_miao_cd(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainRebirthMiaoCd(*militarypb.C2SCaptainRebirthMiaoCdProto, iface.HeroController)
		}).ProcessCaptainRebirthMiaoCd(&militarypb.C2SCaptainRebirthMiaoCdProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_collect_captain_training_exp(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessCollectCaptainTrainingExp(iface.HeroController)
	}).ProcessCollectCaptainTrainingExp(hc)
}
func (gm *GmModule) handleMilitary_c2s_captain_train_exp(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessCaptainTrainExp(iface.HeroController)
	}).ProcessCaptainTrainExp(hc)
}
func (gm *GmModule) handleMilitary_c2s_captain_can_collect_exp(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessCaptainCanCollectExp(iface.HeroController)
	}).ProcessCaptainCanCollectExp(hc)
}
func (gm *GmModule) handleMilitary_c2s_use_training_exp_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUseTrainingExpGoods_deprecated(*militarypb.C2SUseTrainingExpGoodsProto, iface.HeroController)
		}).ProcessUseTrainingExpGoods_deprecated(&militarypb.C2SUseTrainingExpGoodsProto{

			CaptainId: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),

			Count: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_use_level_exp_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUseLevelExpGoods_deprecated(*militarypb.C2SUseLevelExpGoodsProto, iface.HeroController)
		}).ProcessUseLevelExpGoods_deprecated(&militarypb.C2SUseLevelExpGoodsProto{

			CaptainId: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),

			Count: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_use_level_exp_goods2(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUseLevelExpGoods2(*militarypb.C2SUseLevelExpGoods2Proto, iface.HeroController)
		}).ProcessUseLevelExpGoods2(&militarypb.C2SUseLevelExpGoods2Proto{

			Captain: parseInt32(strArray[0]),

			GoodsId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_auto_use_goods_until_captain_levelup(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessAutoUseGoodsUntilCaptainLevelup(*militarypb.C2SAutoUseGoodsUntilCaptainLevelupProto, iface.HeroController)
		}).ProcessAutoUseGoodsUntilCaptainLevelup(&militarypb.C2SAutoUseGoodsUntilCaptainLevelupProto{

			Captain: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_get_max_recruit_soldier(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessGetMaxRecruitSoldier_deprecated(iface.HeroController)
	}).ProcessGetMaxRecruitSoldier_deprecated(hc)
}
func (gm *GmModule) handleMilitary_c2s_get_max_heal_soldier(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessGetMaxHealSoldier_deprecated(iface.HeroController)
	}).ProcessGetMaxHealSoldier_deprecated(hc)
}
func (gm *GmModule) handleMilitary_c2s_jiu_guan_consult(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessJiuGuanConsult(iface.HeroController)
	}).ProcessJiuGuanConsult(hc)
}
func (gm *GmModule) handleMilitary_c2s_jiu_guan_refresh(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessJiuGuanRefresh(*militarypb.C2SJiuGuanRefreshProto, iface.HeroController)
		}).ProcessJiuGuanRefresh(&militarypb.C2SJiuGuanRefreshProto{

			AutoMax: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_unlock_captain_restraint_spell(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUnlockCaptainRestraintSpell_deprecated(*militarypb.C2SUnlockCaptainRestraintSpellProto, iface.HeroController)
		}).ProcessUnlockCaptainRestraintSpell_deprecated(&militarypb.C2SUnlockCaptainRestraintSpellProto{

			Captain: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_get_captain_stat_details(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessGetCaptainStatDetails_deprecated(*militarypb.C2SGetCaptainStatDetailsProto, iface.HeroController)
		}).ProcessGetCaptainStatDetails_deprecated(&militarypb.C2SGetCaptainStatDetailsProto{

			Captain: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_stat_details(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainStatDetails(*militarypb.C2SCaptainStatDetailsProto, iface.HeroController)
		}).ProcessCaptainStatDetails(&militarypb.C2SCaptainStatDetailsProto{

			Captain: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_update_captain_official(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUpdateCaptainOfficial_deprecated(*militarypb.C2SUpdateCaptainOfficialProto, iface.HeroController)
		}).ProcessUpdateCaptainOfficial_deprecated(&militarypb.C2SUpdateCaptainOfficialProto{

			Captain: parseInt32(strArray[0]),

			Official: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_set_captain_official(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessSetCaptainOfficial(*militarypb.C2SSetCaptainOfficialProto, iface.HeroController)
		}).ProcessSetCaptainOfficial(&militarypb.C2SSetCaptainOfficialProto{

			Captain: parseInt32Array(strArray[0]),

			Official: parseInt32Array(strArray[1]),

			OfficialIdx: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_leave_captain_official(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessLeaveCaptainOfficial_deprecated(*militarypb.C2SLeaveCaptainOfficialProto, iface.HeroController)
		}).ProcessLeaveCaptainOfficial_deprecated(&militarypb.C2SLeaveCaptainOfficialProto{

			Captain: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_use_gong_xun_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUseGongXunGoods(*militarypb.C2SUseGongXunGoodsProto, iface.HeroController)
		}).ProcessUseGongXunGoods(&militarypb.C2SUseGongXunGoodsProto{

			Captain: parseInt32(strArray[0]),

			GoodsId: parseInt32Array(strArray[1]),

			Count: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_use_gongxun_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessUseGongxunGoods_deprecated(*militarypb.C2SUseGongxunGoodsProto, iface.HeroController)
		}).ProcessUseGongxunGoods_deprecated(&militarypb.C2SUseGongxunGoodsProto{

			Captain: parseInt32(strArray[0]),

			GoodsId: parseInt32Array(strArray[1]),

			Count: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_close_fight_guide(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCloseFightGuide(*militarypb.C2SCloseFightGuideProto, iface.HeroController)
		}).ProcessCloseFightGuide(&militarypb.C2SCloseFightGuideProto{

			Close: parseBool(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_view_other_hero_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessViewOtherHeroCaptain(*militarypb.C2SViewOtherHeroCaptainProto, iface.HeroController)
		}).ProcessViewOtherHeroCaptain(&militarypb.C2SViewOtherHeroCaptainProto{

			HeroId: parseBytes(strArray[0]),

			CaptainId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_born(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainBorn(*militarypb.C2SCaptainBornProto, iface.HeroController)
		}).ProcessCaptainBorn(&militarypb.C2SCaptainBornProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_upstar(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainUpstar(*militarypb.C2SCaptainUpstarProto, iface.HeroController)
		}).ProcessCaptainUpstar(&militarypb.C2SCaptainUpstarProto{

			CaptainId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_captain_exchange(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessCaptainExchange(*militarypb.C2SCaptainExchangeProto, iface.HeroController)
		}).ProcessCaptainExchange(&militarypb.C2SCaptainExchangeProto{

			Cap1Id: parseInt32(strArray[0]),

			Cap2Id: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_notice_captain_has_viewed(amount string, hc iface.HeroController) {
	gm.modules.MilitaryModule().(interface {
		ProcessNoticeCaptainHasViewed(iface.HeroController)
	}).ProcessNoticeCaptainHasViewed(hc)
}
func (gm *GmModule) handleMilitary_c2s_activate_captain_friendship(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessActivateCaptainFetter(*militarypb.C2SActivateCaptainFriendshipProto, iface.HeroController)
		}).ProcessActivateCaptainFetter(&militarypb.C2SActivateCaptainFriendshipProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMilitary_c2s_notice_official_has_viewed(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MilitaryModule().(interface {
			ProcessNoticeOfficialHasViewed(*militarypb.C2SNoticeOfficialHasViewedProto, iface.HeroController)
		}).ProcessNoticeOfficialHasViewed(&militarypb.C2SNoticeOfficialHasViewedProto{

			OfficialId: parseInt32(strArray[0]),

			OfficialIdx: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerMingc() *gm_group {
	group := &gm_group{
		tab: "mingc",
		handler: []*gm_handler{
			newStringHandler("mingc_list", " ", gm.handleMingc_c2s_mingc_list),
			newStringHandler("view_mingc", " ", gm.handleMingc_c2s_view_mingc),
			newStringHandler("mc_build", " ", gm.handleMingc_c2s_mc_build),
			newStringHandler("mc_build_log", " ", gm.handleMingc_c2s_mc_build_log),
			newStringHandler("mingc_host_guild", " ", gm.handleMingc_c2s_mingc_host_guild),
		},
	}

	return group
}

func (gm *GmModule) handleMingc_c2s_mingc_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcModule().(interface {
			ProcessMingcList(*mingcpb.C2SMingcListProto, iface.HeroController)
		}).ProcessMingcList(&mingcpb.C2SMingcListProto{

			Ver: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_c2s_view_mingc(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcModule().(interface {
			ProcessViewMingc(*mingcpb.C2SViewMingcProto, iface.HeroController)
		}).ProcessViewMingc(&mingcpb.C2SViewMingcProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_c2s_mc_build(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcModule().(interface {
			ProcessMcBuild(*mingcpb.C2SMcBuildProto, iface.HeroController)
		}).ProcessMcBuild(&mingcpb.C2SMcBuildProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_c2s_mc_build_log(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcModule().(interface {
			ProcessMcBuildLog(*mingcpb.C2SMcBuildLogProto, iface.HeroController)
		}).ProcessMcBuildLog(&mingcpb.C2SMcBuildLogProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_c2s_mingc_host_guild(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcModule().(interface {
			ProcessMingcHostGuild(*mingcpb.C2SMingcHostGuildProto, iface.HeroController)
		}).ProcessMingcHostGuild(&mingcpb.C2SMingcHostGuildProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerMingc_war() *gm_group {
	group := &gm_group{
		tab: "mingc_war",
		handler: []*gm_handler{
			newStringHandler("view_mc_war_self_guild", " ", gm.handleMingc_war_c2s_view_mc_war_self_guild),
			newStringHandler("view_mc_war", " ", gm.handleMingc_war_c2s_view_mc_war),
			newStringHandler("apply_atk", " ", gm.handleMingc_war_c2s_apply_atk),
			newStringHandler("apply_ast", " ", gm.handleMingc_war_c2s_apply_ast),
			newStringHandler("cancel_apply_ast", " ", gm.handleMingc_war_c2s_cancel_apply_ast),
			newStringHandler("reply_apply_ast", " ", gm.handleMingc_war_c2s_reply_apply_ast),
			newStringHandler("view_mingc_war_mc", " ", gm.handleMingc_war_c2s_view_mingc_war_mc),
			newStringHandler("join_fight", " ", gm.handleMingc_war_c2s_join_fight),
			newStringHandler("quit_fight", " ", gm.handleMingc_war_c2s_quit_fight),
			newStringHandler("scene_move", " ", gm.handleMingc_war_c2s_scene_move),
			newStringHandler("scene_back", " ", gm.handleMingc_war_c2s_scene_back),
			newStringHandler("scene_speed_up", " ", gm.handleMingc_war_c2s_scene_speed_up),
			newStringHandler("scene_troop_relive", " ", gm.handleMingc_war_c2s_scene_troop_relive),
			newStringHandler("view_mc_war_scene", " ", gm.handleMingc_war_c2s_view_mc_war_scene),
			newStringHandler("watch", " ", gm.handleMingc_war_c2s_watch),
			newStringHandler("quit_watch", " ", gm.handleMingc_war_c2s_quit_watch),
			newStringHandler("view_mc_war_record", " ", gm.handleMingc_war_c2s_view_mc_war_record),
			newStringHandler("view_mc_war_troop_record", " ", gm.handleMingc_war_c2s_view_mc_war_troop_record),
			newStringHandler("view_scene_troop_record", " ", gm.handleMingc_war_c2s_view_scene_troop_record),
			newStringHandler("apply_refresh_rank", " ", gm.handleMingc_war_c2s_apply_refresh_rank),
			newStringHandler("view_my_guild_member_rank", " ", gm.handleMingc_war_c2s_view_my_guild_member_rank),
			newStringHandler("scene_change_mode", " ", gm.handleMingc_war_c2s_scene_change_mode),
			newStringHandler("scene_tou_shi_building_turn_to", " ", gm.handleMingc_war_c2s_scene_tou_shi_building_turn_to),
			newStringHandler("scene_tou_shi_building_fire", " ", gm.handleMingc_war_c2s_scene_tou_shi_building_fire),
			newStringHandler("scene_drum", " ", gm.handleMingc_war_c2s_scene_drum),
		},
	}

	return group
}

func (gm *GmModule) handleMingc_war_c2s_view_mc_war_self_guild(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessViewMcWarSelfGuild(iface.HeroController)
	}).ProcessViewMcWarSelfGuild(hc)
}
func (gm *GmModule) handleMingc_war_c2s_view_mc_war(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMcWar(*mingc_warpb.C2SViewMcWarProto, iface.HeroController)
		}).ProcessViewMcWar(&mingc_warpb.C2SViewMcWarProto{

			Ver: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_apply_atk(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessApplyAtk(*mingc_warpb.C2SApplyAtkProto, iface.HeroController)
		}).ProcessApplyAtk(&mingc_warpb.C2SApplyAtkProto{

			Mcid: parseInt32(strArray[0]),

			Cost: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_apply_ast(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessApplyAst(*mingc_warpb.C2SApplyAstProto, iface.HeroController)
		}).ProcessApplyAst(&mingc_warpb.C2SApplyAstProto{

			Mcid: parseInt32(strArray[0]),

			Atk: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_cancel_apply_ast(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessCancelApplyAst(*mingc_warpb.C2SCancelApplyAstProto, iface.HeroController)
		}).ProcessCancelApplyAst(&mingc_warpb.C2SCancelApplyAstProto{

			Mcid: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_reply_apply_ast(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessReplyApplyAst(*mingc_warpb.C2SReplyApplyAstProto, iface.HeroController)
		}).ProcessReplyApplyAst(&mingc_warpb.C2SReplyApplyAstProto{

			Mcid: parseInt32(strArray[0]),

			Gid: parseInt32(strArray[1]),

			Agree: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_view_mingc_war_mc(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMingcWarMc(*mingc_warpb.C2SViewMingcWarMcProto, iface.HeroController)
		}).ProcessViewMingcWarMc(&mingc_warpb.C2SViewMingcWarMcProto{

			Mcid: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_join_fight(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessJoinFight(*mingc_warpb.C2SJoinFightProto, iface.HeroController)
		}).ProcessJoinFight(&mingc_warpb.C2SJoinFightProto{

			Mcid: parseInt32(strArray[0]),

			CaptainId: parseInt32Array(strArray[1]),

			XIndex: parseInt32Array(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_quit_fight(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessQuitFight(iface.HeroController)
	}).ProcessQuitFight(hc)
}
func (gm *GmModule) handleMingc_war_c2s_scene_move(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessSceneMove(*mingc_warpb.C2SSceneMoveProto, iface.HeroController)
		}).ProcessSceneMove(&mingc_warpb.C2SSceneMoveProto{

			DestPosX: parseInt32(strArray[0]),

			DestPosY: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_back(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessSceneBack(iface.HeroController)
	}).ProcessSceneBack(hc)
}
func (gm *GmModule) handleMingc_war_c2s_scene_speed_up(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessSceneSpeedUp(*mingc_warpb.C2SSceneSpeedUpProto, iface.HeroController)
		}).ProcessSceneSpeedUp(&mingc_warpb.C2SSceneSpeedUpProto{

			Id: parseBytes(strArray[0]),

			GoodsId: parseInt32(strArray[1]),

			Money: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_troop_relive(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessSceneTroopRelive(iface.HeroController)
	}).ProcessSceneTroopRelive(hc)
}
func (gm *GmModule) handleMingc_war_c2s_view_mc_war_scene(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMcWarScene(*mingc_warpb.C2SViewMcWarSceneProto, iface.HeroController)
		}).ProcessViewMcWarScene(&mingc_warpb.C2SViewMcWarSceneProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_watch(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessWatch(*mingc_warpb.C2SWatchProto, iface.HeroController)
		}).ProcessWatch(&mingc_warpb.C2SWatchProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_quit_watch(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessQuitWatch(*mingc_warpb.C2SQuitWatchProto, iface.HeroController)
		}).ProcessQuitWatch(&mingc_warpb.C2SQuitWatchProto{

			McId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_view_mc_war_record(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMcWarRecord(*mingc_warpb.C2SViewMcWarRecordProto, iface.HeroController)
		}).ProcessViewMcWarRecord(&mingc_warpb.C2SViewMcWarRecordProto{

			WarId: parseInt32(strArray[0]),

			McId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_view_mc_war_troop_record(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMcWarTroopRecord(*mingc_warpb.C2SViewMcWarTroopRecordProto, iface.HeroController)
		}).ProcessViewMcWarTroopRecord(&mingc_warpb.C2SViewMcWarTroopRecordProto{

			WarId: parseInt32(strArray[0]),

			McId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_view_scene_troop_record(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessViewSceneTroopRecord(iface.HeroController)
	}).ProcessViewSceneTroopRecord(hc)
}
func (gm *GmModule) handleMingc_war_c2s_apply_refresh_rank(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessApplyRefreshRank(*mingc_warpb.C2SApplyRefreshRankProto, iface.HeroController)
		}).ProcessApplyRefreshRank(&mingc_warpb.C2SApplyRefreshRankProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_view_my_guild_member_rank(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessViewMyGuildMemberRank(*mingc_warpb.C2SViewMyGuildMemberRankProto, iface.HeroController)
		}).ProcessViewMyGuildMemberRank(&mingc_warpb.C2SViewMyGuildMemberRankProto{

			WarId: parseInt32(strArray[0]),

			McId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_change_mode(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessSceneChangeMode(*mingc_warpb.C2SSceneChangeModeProto, iface.HeroController)
		}).ProcessSceneChangeMode(&mingc_warpb.C2SSceneChangeModeProto{

			Mode: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_tou_shi_building_turn_to(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessSceneTouShiBuildingTurnTo(*mingc_warpb.C2SSceneTouShiBuildingTurnToProto, iface.HeroController)
		}).ProcessSceneTouShiBuildingTurnTo(&mingc_warpb.C2SSceneTouShiBuildingTurnToProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),

			Left: parseBool(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_tou_shi_building_fire(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MingcWarModule().(interface {
			ProcessSceneTouShiBuildingFire(*mingc_warpb.C2SSceneTouShiBuildingFireProto, iface.HeroController)
		}).ProcessSceneTouShiBuildingFire(&mingc_warpb.C2SSceneTouShiBuildingFireProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMingc_war_c2s_scene_drum(amount string, hc iface.HeroController) {
	gm.modules.MingcWarModule().(interface {
		ProcessSceneDrum(iface.HeroController)
	}).ProcessSceneDrum(hc)
}

func (gm *GmModule) initHandlerMisc() *gm_group {
	group := &gm_group{
		tab: "misc",
		handler: []*gm_handler{
			newStringHandler("heart_beat", " ", gm.handleMisc_c2s_heart_beat),
			newStringHandler("background_heart_beat", " ", gm.handleMisc_c2s_background_heart_beat),
			newStringHandler("background_weakup", " ", gm.handleMisc_c2s_background_weakup),
			newStringHandler("config", " ", gm.handleMisc_c2s_config),
			newStringHandler("configlua", " ", gm.handleMisc_c2s_configlua),
			newStringHandler("client_log", " ", gm.handleMisc_c2s_client_log),
			newStringHandler("sync_time", " ", gm.handleMisc_c2s_sync_time),
			newStringHandler("block", " ", gm.handleMisc_c2s_block),
			newStringHandler("ping", " ", gm.handleMisc_c2s_ping),
			newStringHandler("client_version", " ", gm.handleMisc_c2s_client_version),
			newStringHandler("update_pf_token", " ", gm.handleMisc_c2s_update_pf_token),
			newStringHandler("settings", " ", gm.handleMisc_c2s_settings),
			newStringHandler("settings_to_default", " ", gm.handleMisc_c2s_settings_to_default),
			newStringHandler("update_location", " ", gm.handleMisc_c2s_update_location),
			newStringHandler("collect_charge_prize", " ", gm.handleMisc_c2s_collect_charge_prize),
			newStringHandler("collect_daily_bargain", " ", gm.handleMisc_c2s_collect_daily_bargain),
			newStringHandler("activate_duration_card", " ", gm.handleMisc_c2s_activate_duration_card),
			newStringHandler("collect_duration_card_daily_prize", " ", gm.handleMisc_c2s_collect_duration_card_daily_prize),
			newStringHandler("set_privacy_setting", " ", gm.handleMisc_c2s_set_privacy_setting),
			newStringHandler("set_default_privacy_settings", " ", gm.handleMisc_c2s_set_default_privacy_settings),
			newStringHandler("get_product_info", " ", gm.handleMisc_c2s_get_product_info),
		},
	}

	return group
}

func (gm *GmModule) handleMisc_c2s_heart_beat(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessHeartBeat(iface.HeroController)
	}).ProcessHeartBeat(hc)
}
func (gm *GmModule) handleMisc_c2s_background_heart_beat(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessBackgroudHeartBeat(iface.HeroController)
	}).ProcessBackgroudHeartBeat(hc)
}
func (gm *GmModule) handleMisc_c2s_background_weakup(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessBackgroudWeakup(iface.HeroController)
	}).ProcessBackgroudWeakup(hc)
}
func (gm *GmModule) handleMisc_c2s_config(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessRequestConfig(*miscpb.C2SConfigProto, iface.HeroController)
		}).ProcessRequestConfig(&miscpb.C2SConfigProto{

			Version: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_configlua(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessRequestLuaConfig(*miscpb.C2SConfigluaProto, iface.HeroController)
		}).ProcessRequestLuaConfig(&miscpb.C2SConfigluaProto{

			Version: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_client_log(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessClientLog(*miscpb.C2SClientLogProto, iface.HeroController)
		}).ProcessClientLog(&miscpb.C2SClientLogProto{

			Level: parseString(strArray[0]),

			Text: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_sync_time(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessSyncTime(*miscpb.C2SSyncTimeProto, iface.HeroController)
		}).ProcessSyncTime(&miscpb.C2SSyncTimeProto{

			ClientTime: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_block(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessGetBlock(iface.HeroController)
	}).ProcessGetBlock(hc)
}
func (gm *GmModule) handleMisc_c2s_ping(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessPing(iface.HeroController)
	}).ProcessPing(hc)
}
func (gm *GmModule) handleMisc_c2s_client_version(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessClientVersion(*miscpb.C2SClientVersionProto, iface.HeroController)
		}).ProcessClientVersion(&miscpb.C2SClientVersionProto{

			Os: parseString(strArray[0]),

			T: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_update_pf_token(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessUpdatePfToken(*miscpb.C2SUpdatePfTokenProto, iface.HeroController)
		}).ProcessUpdatePfToken(&miscpb.C2SUpdatePfTokenProto{

			Token: parseString(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_settings(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessSettings(*miscpb.C2SSettingsProto, iface.HeroController)
		}).ProcessSettings(&miscpb.C2SSettingsProto{

			SettingType: parseInt32(strArray[0]),

			Open: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_settings_to_default(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessSettingsToDefault(iface.HeroController)
	}).ProcessSettingsToDefault(hc)
}
func (gm *GmModule) handleMisc_c2s_update_location(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessUpdateLocation(*miscpb.C2SUpdateLocationProto, iface.HeroController)
		}).ProcessUpdateLocation(&miscpb.C2SUpdateLocationProto{

			Location: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_collect_charge_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessCollectChargePrize(*miscpb.C2SCollectChargePrizeProto, iface.HeroController)
		}).ProcessCollectChargePrize(&miscpb.C2SCollectChargePrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_collect_daily_bargain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessCollectDailyBargain(*miscpb.C2SCollectDailyBargainProto, iface.HeroController)
		}).ProcessCollectDailyBargain(&miscpb.C2SCollectDailyBargainProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_activate_duration_card(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessActivateDurationCard(*miscpb.C2SActivateDurationCardProto, iface.HeroController)
		}).ProcessActivateDurationCard(&miscpb.C2SActivateDurationCardProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_collect_duration_card_daily_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessCollectDurationCardDailyPrize(*miscpb.C2SCollectDurationCardDailyPrizeProto, iface.HeroController)
		}).ProcessCollectDurationCardDailyPrize(&miscpb.C2SCollectDurationCardDailyPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_set_privacy_setting(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessSetPrivacySetting(*miscpb.C2SSetPrivacySettingProto, iface.HeroController)
		}).ProcessSetPrivacySetting(&miscpb.C2SSetPrivacySettingProto{

			SettingId: parseInt32(strArray[0]),

			OpenOrClose: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleMisc_c2s_set_default_privacy_settings(amount string, hc iface.HeroController) {
	gm.modules.MiscModule().(interface {
		ProcessSetDefaultPrivacySettings(iface.HeroController)
	}).ProcessSetDefaultPrivacySettings(hc)
}
func (gm *GmModule) handleMisc_c2s_get_product_info(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.MiscModule().(interface {
			ProcessGetProductInfo(*miscpb.C2SGetProductInfoProto, iface.HeroController)
		}).ProcessGetProductInfo(&miscpb.C2SGetProductInfoProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerPromotion() *gm_group {
	group := &gm_group{
		tab: "promotion",
		handler: []*gm_handler{
			newStringHandler("collect_login_day_prize", " ", gm.handlePromotion_c2s_collect_login_day_prize),
			newStringHandler("buy_level_fund", " ", gm.handlePromotion_c2s_buy_level_fund),
			newStringHandler("collect_level_fund", " ", gm.handlePromotion_c2s_collect_level_fund),
			newStringHandler("collect_daily_sp", " ", gm.handlePromotion_c2s_collect_daily_sp),
			newStringHandler("collect_free_gift", " ", gm.handlePromotion_c2s_collect_free_gift),
			newStringHandler("buy_time_limit_gift", " ", gm.handlePromotion_c2s_buy_time_limit_gift),
			newStringHandler("buy_event_limit_gift", " ", gm.handlePromotion_c2s_buy_event_limit_gift),
		},
	}

	return group
}

func (gm *GmModule) handlePromotion_c2s_collect_login_day_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessCollectLogin7DayPrize(*promotionpb.C2SCollectLoginDayPrizeProto, iface.HeroController)
		}).ProcessCollectLogin7DayPrize(&promotionpb.C2SCollectLoginDayPrizeProto{

			Day: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handlePromotion_c2s_buy_level_fund(amount string, hc iface.HeroController) {
	gm.modules.PromotionModule().(interface {
		ProcessBuyHeroLevelFund(iface.HeroController)
	}).ProcessBuyHeroLevelFund(hc)
}
func (gm *GmModule) handlePromotion_c2s_collect_level_fund(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessCollectHeroLevelFund(*promotionpb.C2SCollectLevelFundProto, iface.HeroController)
		}).ProcessCollectHeroLevelFund(&promotionpb.C2SCollectLevelFundProto{

			Level: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handlePromotion_c2s_collect_daily_sp(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessCollectDailySp(*promotionpb.C2SCollectDailySpProto, iface.HeroController)
		}).ProcessCollectDailySp(&promotionpb.C2SCollectDailySpProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handlePromotion_c2s_collect_free_gift(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessCollectFreeGift(*promotionpb.C2SCollectFreeGiftProto, iface.HeroController)
		}).ProcessCollectFreeGift(&promotionpb.C2SCollectFreeGiftProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handlePromotion_c2s_buy_time_limit_gift(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessBuyTimeLimitGift(*promotionpb.C2SBuyTimeLimitGiftProto, iface.HeroController)
		}).ProcessBuyTimeLimitGift(&promotionpb.C2SBuyTimeLimitGiftProto{

			GrpId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handlePromotion_c2s_buy_event_limit_gift(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.PromotionModule().(interface {
			ProcessBuyEventLimitGift(*promotionpb.C2SBuyEventLimitGiftProto, iface.HeroController)
		}).ProcessBuyEventLimitGift(&promotionpb.C2SBuyEventLimitGiftProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerQuestion() *gm_group {
	group := &gm_group{
		tab: "question",
		handler: []*gm_handler{
			newStringHandler("start", " ", gm.handleQuestion_c2s_start),
			newStringHandler("answer", " ", gm.handleQuestion_c2s_answer),
			newStringHandler("next", " ", gm.handleQuestion_c2s_next),
			newStringHandler("get_prize", " ", gm.handleQuestion_c2s_get_prize),
		},
	}

	return group
}

func (gm *GmModule) handleQuestion_c2s_start(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.QuestionModule().(interface {
			ProcessStart(*questionpb.C2SStartProto, iface.HeroController)
		}).ProcessStart(&questionpb.C2SStartProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleQuestion_c2s_answer(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.QuestionModule().(interface {
			ProcessAnswer(*questionpb.C2SAnswerProto, iface.HeroController)
		}).ProcessAnswer(&questionpb.C2SAnswerProto{

			Id: parseInt32(strArray[0]),

			Right: parseBool(strArray[1]),

			Answer: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleQuestion_c2s_next(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.QuestionModule().(interface {
			ProcessNext(*questionpb.C2SNextProto, iface.HeroController)
		}).ProcessNext(&questionpb.C2SNextProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleQuestion_c2s_get_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.QuestionModule().(interface {
			ProcessGetPrize(*questionpb.C2SGetPrizeProto, iface.HeroController)
		}).ProcessGetPrize(&questionpb.C2SGetPrizeProto{

			Score: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerRandom_event() *gm_group {
	group := &gm_group{
		tab: "random_event",
		handler: []*gm_handler{
			newStringHandler("choose_option", " ", gm.handleRandom_event_c2s_choose_option),
			newStringHandler("open_event", " ", gm.handleRandom_event_c2s_open_event),
		},
	}

	return group
}

func (gm *GmModule) handleRandom_event_c2s_choose_option(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RandomEventModule().(interface {
			ProcessChooseOption(*random_eventpb.C2SChooseOptionProto, iface.HeroController)
		}).ProcessChooseOption(&random_eventpb.C2SChooseOptionProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),

			Option: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleRandom_event_c2s_open_event(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RandomEventModule().(interface {
			ProcessOpenEvent(*random_eventpb.C2SOpenEventProto, iface.HeroController)
		}).ProcessOpenEvent(&random_eventpb.C2SOpenEventProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerRank() *gm_group {
	group := &gm_group{
		tab: "rank",
		handler: []*gm_handler{
			newStringHandler("request_rank", " ", gm.handleRank_c2s_request_rank),
		},
	}

	return group
}

func (gm *GmModule) handleRank_c2s_request_rank(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RankModule().(interface {
			ProcessRequestRank(*rankpb.C2SRequestRankProto, iface.HeroController)
		}).ProcessRequestRank(&rankpb.C2SRequestRankProto{

			RankType: parseInt32(strArray[0]),

			Name: parseString(strArray[1]),

			Self: parseBool(strArray[2]),

			StartCount: parseInt32(strArray[3]),

			JunXianLevel: parseInt32(strArray[4]),

			SubType: parseInt32(strArray[5]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerRed_packet() *gm_group {
	group := &gm_group{
		tab: "red_packet",
		handler: []*gm_handler{
			newStringHandler("buy", " ", gm.handleRed_packet_c2s_buy),
			newStringHandler("create", " ", gm.handleRed_packet_c2s_create),
			newStringHandler("grab", " ", gm.handleRed_packet_c2s_grab),
		},
	}

	return group
}

func (gm *GmModule) handleRed_packet_c2s_buy(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RedPacketModule().(interface {
			ProcessBuy(*red_packetpb.C2SBuyProto, iface.HeroController)
		}).ProcessBuy(&red_packetpb.C2SBuyProto{

			DataId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRed_packet_c2s_create(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RedPacketModule().(interface {
			ProcessCreate(*red_packetpb.C2SCreateProto, iface.HeroController)
		}).ProcessCreate(&red_packetpb.C2SCreateProto{

			DataId: parseInt32(strArray[0]),

			Count: parseInt32(strArray[2]),

			Text: parseString(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleRed_packet_c2s_grab(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RedPacketModule().(interface {
			ProcessGrab(*red_packetpb.C2SGrabProto, iface.HeroController)
		}).ProcessGrab(&red_packetpb.C2SGrabProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerRegion() *gm_group {
	group := &gm_group{
		tab: "region",
		handler: []*gm_handler{
			newStringHandler("update_self_view", " ", gm.handleRegion_c2s_update_self_view),
			newStringHandler("close_view", " ", gm.handleRegion_c2s_close_view),
			newStringHandler("pre_invasion_target", " ", gm.handleRegion_c2s_pre_invasion_target),
			newStringHandler("watch_base_unit", " ", gm.handleRegion_c2s_watch_base_unit),
			newStringHandler("request_troop_unit", " ", gm.handleRegion_c2s_request_troop_unit),
			newStringHandler("request_ruins_base", " ", gm.handleRegion_c2s_request_ruins_base),
			newStringHandler("use_mian_goods", " ", gm.handleRegion_c2s_use_mian_goods),
			newStringHandler("upgrade_base", " ", gm.handleRegion_c2s_upgrade_base),
			newStringHandler("white_flag_detail", " ", gm.handleRegion_c2s_white_flag_detail),
			newStringHandler("get_buy_prosperity_cost", " ", gm.handleRegion_c2s_get_buy_prosperity_cost),
			newStringHandler("buy_prosperity", " ", gm.handleRegion_c2s_buy_prosperity),
			newStringHandler("switch_action", " ", gm.handleRegion_c2s_switch_action),
			newStringHandler("request_military_push", " ", gm.handleRegion_c2s_request_military_push),
			newStringHandler("create_base", " ", gm.handleRegion_c2s_create_base),
			newStringHandler("fast_move_base", " ", gm.handleRegion_c2s_fast_move_base),
			newStringHandler("invasion", " ", gm.handleRegion_c2s_invasion),
			newStringHandler("cancel_invasion", " ", gm.handleRegion_c2s_cancel_invasion),
			newStringHandler("repatriate", " ", gm.handleRegion_c2s_repatriate),
			newStringHandler("baoz_repatriate", " ", gm.handleRegion_c2s_baoz_repatriate),
			newStringHandler("speed_up", " ", gm.handleRegion_c2s_speed_up),
			newStringHandler("expel", " ", gm.handleRegion_c2s_expel),
			newStringHandler("favorite_pos", " ", gm.handleRegion_c2s_favorite_pos),
			newStringHandler("favorite_pos_list", " ", gm.handleRegion_c2s_favorite_pos_list),
			newStringHandler("get_prev_investigate", " ", gm.handleRegion_c2s_get_prev_investigate),
			newStringHandler("investigate", " ", gm.handleRegion_c2s_investigate),
			newStringHandler("investigate_invade", " ", gm.handleRegion_c2s_investigate_invade),
			newStringHandler("use_multi_level_npc_times_goods", " ", gm.handleRegion_c2s_use_multi_level_npc_times_goods),
			newStringHandler("use_invase_hero_times_goods", " ", gm.handleRegion_c2s_use_invase_hero_times_goods),
			newStringHandler("calc_move_speed", " ", gm.handleRegion_c2s_calc_move_speed),
			newStringHandler("list_enemy_pos", " ", gm.handleRegion_c2s_list_enemy_pos),
			newStringHandler("search_baoz_npc", " ", gm.handleRegion_c2s_search_baoz_npc),
			newStringHandler("home_ast_defending_info", " ", gm.handleRegion_c2s_home_ast_defending_info),
			newStringHandler("guild_please_help_me", " ", gm.handleRegion_c2s_guild_please_help_me),
			newStringHandler("create_assembly", " ", gm.handleRegion_c2s_create_assembly),
			newStringHandler("show_assembly", " ", gm.handleRegion_c2s_show_assembly),
			newStringHandler("join_assembly", " ", gm.handleRegion_c2s_join_assembly),
			newStringHandler("create_guild_workshop", " ", gm.handleRegion_c2s_create_guild_workshop),
			newStringHandler("show_guild_workshop", " ", gm.handleRegion_c2s_show_guild_workshop),
			newStringHandler("hurt_guild_workshop", " ", gm.handleRegion_c2s_hurt_guild_workshop),
			newStringHandler("remove_guild_workshop", " ", gm.handleRegion_c2s_remove_guild_workshop),
			newStringHandler("catch_guild_workshop_logs", " ", gm.handleRegion_c2s_catch_guild_workshop_logs),
			newStringHandler("get_self_baoz", " ", gm.handleRegion_c2s_get_self_baoz),
		},
	}

	return group
}

func (gm *GmModule) handleRegion_c2s_update_self_view(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessUpdateSelfView(*regionpb.C2SUpdateSelfViewProto, iface.HeroController)
		}).ProcessUpdateSelfView(&regionpb.C2SUpdateSelfViewProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),

			LenX: parseInt32(strArray[2]),

			LenY: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_close_view(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessCloseView(iface.HeroController)
	}).ProcessCloseView(hc)
}
func (gm *GmModule) handleRegion_c2s_pre_invasion_target(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessPreInvasionTarget(*regionpb.C2SPreInvasionTargetProto, iface.HeroController)
		}).ProcessPreInvasionTarget(&regionpb.C2SPreInvasionTargetProto{

			MapId: parseInt32(strArray[0]),

			Target: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_watch_base_unit(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessWatchBaseUnit(*regionpb.C2SWatchBaseUnitProto, iface.HeroController)
		}).ProcessWatchBaseUnit(&regionpb.C2SWatchBaseUnitProto{

			Target: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_request_troop_unit(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessRequestTroopUnit(*regionpb.C2SRequestTroopUnitProto, iface.HeroController)
		}).ProcessRequestTroopUnit(&regionpb.C2SRequestTroopUnitProto{

			TroopId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_request_ruins_base(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessRequestRuinsBase(*regionpb.C2SRequestRuinsBaseProto, iface.HeroController)
		}).ProcessRequestRuinsBase(&regionpb.C2SRequestRuinsBaseProto{

			RealmId: parseInt32(strArray[0]),

			PosX: parseInt32(strArray[1]),

			PosY: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_use_mian_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessUseMianGoods(*regionpb.C2SUseMianGoodsProto, iface.HeroController)
		}).ProcessUseMianGoods(&regionpb.C2SUseMianGoodsProto{

			Id: parseInt32(strArray[0]),

			Buy: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_upgrade_base(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessUpgradeBase(iface.HeroController)
	}).ProcessUpgradeBase(hc)
}
func (gm *GmModule) handleRegion_c2s_white_flag_detail(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessGetWhiteFlagDetail(*regionpb.C2SWhiteFlagDetailProto, iface.HeroController)
		}).ProcessGetWhiteFlagDetail(&regionpb.C2SWhiteFlagDetailProto{

			HeroId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_get_buy_prosperity_cost(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessGetBuyProsperityCost(iface.HeroController)
	}).ProcessGetBuyProsperityCost(hc)
}
func (gm *GmModule) handleRegion_c2s_buy_prosperity(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessBuyProsperity(iface.HeroController)
	}).ProcessBuyProsperity(hc)
}
func (gm *GmModule) handleRegion_c2s_switch_action(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessSwitchAction(*regionpb.C2SSwitchActionProto, iface.HeroController)
		}).ProcessSwitchAction(&regionpb.C2SSwitchActionProto{

			Open: parseBool(strArray[0]),

			Condition: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_request_military_push(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessRequestMilitaryPush(*regionpb.C2SRequestMilitaryPushProto, iface.HeroController)
		}).ProcessRequestMilitaryPush(&regionpb.C2SRequestMilitaryPushProto{

			MainMilitary: parseBool(strArray[0]),

			GuildMilitary: parseBool(strArray[1]),

			ToTarget: parseBytes(strArray[2]),

			ToTargetBase: parseBool(strArray[3]),

			FromTarget: parseBytes(strArray[4]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_create_base(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCreateBase(*regionpb.C2SCreateBaseProto, iface.HeroController)
		}).ProcessCreateBase(&regionpb.C2SCreateBaseProto{

			MapId: parseInt32(strArray[0]),

			NewX: parseInt32(strArray[1]),

			NewY: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_fast_move_base(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessFastMoveBase(*regionpb.C2SFastMoveBaseProto, iface.HeroController)
		}).ProcessFastMoveBase(&regionpb.C2SFastMoveBaseProto{

			MapId: parseInt32(strArray[0]),

			NewX: parseInt32(strArray[1]),

			NewY: parseInt32(strArray[2]),

			GoodsId: parseInt32(strArray[3]),

			IsTent: parseBool(strArray[4]),

			Free: parseBool(strArray[5]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_invasion(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessInvasion(*regionpb.C2SInvasionProto, iface.HeroController)
		}).ProcessInvasion(&regionpb.C2SInvasionProto{

			Operate: parseInt32(strArray[0]),

			Target: parseBytes(strArray[1]),

			TroopIndex: parseInt32(strArray[2]),

			TargetLevel: parseInt32(strArray[3]),

			GoodsId: parseInt32(strArray[4]),

			AutoBuy: parseBool(strArray[5]),

			MultiLevelMonsterCount: parseInt32(strArray[6]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_cancel_invasion(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCancelInvasion(*regionpb.C2SCancelInvasionProto, iface.HeroController)
		}).ProcessCancelInvasion(&regionpb.C2SCancelInvasionProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_repatriate(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessRepatriate(*regionpb.C2SRepatriateProto, iface.HeroController)
		}).ProcessRepatriate(&regionpb.C2SRepatriateProto{

			Id: parseBytes(strArray[0]),

			IsTent: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_baoz_repatriate(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessBaozRepatriate(*regionpb.C2SBaozRepatriateProto, iface.HeroController)
		}).ProcessBaozRepatriate(&regionpb.C2SBaozRepatriateProto{

			BaseId: parseBytes(strArray[0]),

			TroopId: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_speed_up(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessSpeedUp(*regionpb.C2SSpeedUpProto, iface.HeroController)
		}).ProcessSpeedUp(&regionpb.C2SSpeedUpProto{

			Id: parseBytes(strArray[0]),

			OtherId: parseBytes(strArray[1]),

			GoodsId: parseInt32(strArray[2]),

			Money: parseBool(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_expel(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessExpel(*regionpb.C2SExpelProto, iface.HeroController)
		}).ProcessExpel(&regionpb.C2SExpelProto{

			Id: parseBytes(strArray[0]),

			Mapid: parseInt32(strArray[1]),

			TroopIndex: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_favorite_pos(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessFavoritePos(*regionpb.C2SFavoritePosProto, iface.HeroController)
		}).ProcessFavoritePos(&regionpb.C2SFavoritePosProto{

			Add: parseBool(strArray[0]),

			Id: parseInt32(strArray[1]),

			PosX: parseInt32(strArray[2]),

			PosY: parseInt32(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_favorite_pos_list(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessFavoritePosList(iface.HeroController)
	}).ProcessFavoritePosList(hc)
}
func (gm *GmModule) handleRegion_c2s_get_prev_investigate(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessGetPrevInvestigate(*regionpb.C2SGetPrevInvestigateProto, iface.HeroController)
		}).ProcessGetPrevInvestigate(&regionpb.C2SGetPrevInvestigateProto{

			HeroId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_investigate(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessInvestigate(*regionpb.C2SInvestigateProto, iface.HeroController)
		}).ProcessInvestigate(&regionpb.C2SInvestigateProto{

			HeroId: parseBytes(strArray[0]),

			Cost: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_investigate_invade(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessInvestigateInvade(*regionpb.C2SInvestigateInvadeProto, iface.HeroController)
		}).ProcessInvestigateInvade(&regionpb.C2SInvestigateInvadeProto{

			HeroId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_use_multi_level_npc_times_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessUseMultiLevelNpcTimesGoods(*regionpb.C2SUseMultiLevelNpcTimesGoodsProto, iface.HeroController)
		}).ProcessUseMultiLevelNpcTimesGoods(&regionpb.C2SUseMultiLevelNpcTimesGoodsProto{

			Id: parseInt32(strArray[0]),

			Buy: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_use_invase_hero_times_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessUseInvaseHeroTimesGoods(*regionpb.C2SUseInvaseHeroTimesGoodsProto, iface.HeroController)
		}).ProcessUseInvaseHeroTimesGoods(&regionpb.C2SUseInvaseHeroTimesGoodsProto{

			Id: parseInt32(strArray[0]),

			Buy: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_calc_move_speed(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCalcMoveSpeed(*regionpb.C2SCalcMoveSpeedProto, iface.HeroController)
		}).ProcessCalcMoveSpeed(&regionpb.C2SCalcMoveSpeedProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_list_enemy_pos(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessListEnemyPos(iface.HeroController)
	}).ProcessListEnemyPos(hc)
}
func (gm *GmModule) handleRegion_c2s_search_baoz_npc(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessSearchBaozNpc(*regionpb.C2SSearchBaozNpcProto, iface.HeroController)
		}).ProcessSearchBaozNpc(&regionpb.C2SSearchBaozNpcProto{

			DataId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_home_ast_defending_info(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessHomeAstDefendingInfo(iface.HeroController)
	}).ProcessHomeAstDefendingInfo(hc)
}
func (gm *GmModule) handleRegion_c2s_guild_please_help_me(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessGuildPleaseHelpMe(iface.HeroController)
	}).ProcessGuildPleaseHelpMe(hc)
}
func (gm *GmModule) handleRegion_c2s_create_assembly(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCreateAssembly(*regionpb.C2SCreateAssemblyProto, iface.HeroController)
		}).ProcessCreateAssembly(&regionpb.C2SCreateAssemblyProto{

			TroopIndex: parseInt32(strArray[0]),

			Target: parseBytes(strArray[1]),

			TargetLevel: parseInt32(strArray[2]),

			WaitIndex: parseInt32(strArray[3]),

			GoodsId: parseInt32(strArray[4]),

			AutoBuy: parseBool(strArray[5]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_show_assembly(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessShowAssembly(*regionpb.C2SShowAssemblyProto, iface.HeroController)
		}).ProcessShowAssembly(&regionpb.C2SShowAssemblyProto{

			Id: parseBytes(strArray[0]),

			Version: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_join_assembly(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessJoinAssembly(*regionpb.C2SJoinAssemblyProto, iface.HeroController)
		}).ProcessJoinAssembly(&regionpb.C2SJoinAssemblyProto{

			Id: parseBytes(strArray[0]),

			TroopIndex: parseInt32(strArray[1]),

			GoodsId: parseInt32(strArray[2]),

			AutoBuy: parseBool(strArray[3]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_create_guild_workshop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCreateGuildWorkshop(*regionpb.C2SCreateGuildWorkshopProto, iface.HeroController)
		}).ProcessCreateGuildWorkshop(&regionpb.C2SCreateGuildWorkshopProto{

			PosX: parseInt32(strArray[0]),

			PosY: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_show_guild_workshop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessShowGuildWorkshop(*regionpb.C2SShowGuildWorkshopProto, iface.HeroController)
		}).ProcessShowGuildWorkshop(&regionpb.C2SShowGuildWorkshopProto{

			BaseId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_hurt_guild_workshop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessHurtGuildWorkshop(*regionpb.C2SHurtGuildWorkshopProto, iface.HeroController)
		}).ProcessHurtGuildWorkshop(&regionpb.C2SHurtGuildWorkshopProto{

			BaseId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_remove_guild_workshop(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessRemoveGuildWorkshop(iface.HeroController)
	}).ProcessRemoveGuildWorkshop(hc)
}
func (gm *GmModule) handleRegion_c2s_catch_guild_workshop_logs(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RegionModule().(interface {
			ProcessCatchGuildWorkshopLogs(*regionpb.C2SCatchGuildWorkshopLogsProto, iface.HeroController)
		}).ProcessCatchGuildWorkshopLogs(&regionpb.C2SCatchGuildWorkshopLogsProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRegion_c2s_get_self_baoz(amount string, hc iface.HeroController) {
	gm.modules.RegionModule().(interface {
		ProcessGetSelfBaoz(iface.HeroController)
	}).ProcessGetSelfBaoz(hc)
}

func (gm *GmModule) initHandlerRelation() *gm_group {
	group := &gm_group{
		tab: "relation",
		handler: []*gm_handler{
			newStringHandler("add_relation", " ", gm.handleRelation_c2s_add_relation),
			newStringHandler("remove_enemy", " ", gm.handleRelation_c2s_remove_enemy),
			newStringHandler("remove_relation", " ", gm.handleRelation_c2s_remove_relation),
			newStringHandler("list_relation", " ", gm.handleRelation_c2s_list_relation),
			newStringHandler("new_list_relation", " ", gm.handleRelation_c2s_new_list_relation),
			newStringHandler("recommend_hero_list", " ", gm.handleRelation_c2s_recommend_hero_list),
			newStringHandler("search_heros", " ", gm.handleRelation_c2s_search_heros),
			newStringHandler("search_hero_by_id", " ", gm.handleRelation_c2s_search_hero_by_id),
			newStringHandler("set_important_friend", " ", gm.handleRelation_c2s_set_important_friend),
			newStringHandler("cancel_important_friend", " ", gm.handleRelation_c2s_cancel_important_friend),
		},
	}

	return group
}

func (gm *GmModule) handleRelation_c2s_add_relation(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessAddRelation(*relationpb.C2SAddRelationProto, iface.HeroController)
		}).ProcessAddRelation(&relationpb.C2SAddRelationProto{

			Friend: parseBool(strArray[0]),

			Id: parseBytes(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_remove_enemy(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessRemoveEnemy(*relationpb.C2SRemoveEnemyProto, iface.HeroController)
		}).ProcessRemoveEnemy(&relationpb.C2SRemoveEnemyProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_remove_relation(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessRemoveRelation(*relationpb.C2SRemoveRelationProto, iface.HeroController)
		}).ProcessRemoveRelation(&relationpb.C2SRemoveRelationProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_list_relation(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessListRelation(*relationpb.C2SListRelationProto, iface.HeroController)
		}).ProcessListRelation(&relationpb.C2SListRelationProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_new_list_relation(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessNewListRelation(*relationpb.C2SNewListRelationProto, iface.HeroController)
		}).ProcessNewListRelation(&relationpb.C2SNewListRelationProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_recommend_hero_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessRecommendHeroList(*relationpb.C2SRecommendHeroListProto, iface.HeroController)
		}).ProcessRecommendHeroList(&relationpb.C2SRecommendHeroListProto{

			Page: parseInt32(strArray[0]),

			NeedLoc: parseBool(strArray[1]),

			Loc: parseInt32(strArray[2]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_search_heros(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessSearchHeros(*relationpb.C2SSearchHerosProto, iface.HeroController)
		}).ProcessSearchHeros(&relationpb.C2SSearchHerosProto{

			Name: parseString(strArray[0]),

			Page: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_search_hero_by_id(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessSearchHeroById(*relationpb.C2SSearchHeroByIdProto, iface.HeroController)
		}).ProcessSearchHeroById(&relationpb.C2SSearchHeroByIdProto{

			HeroId: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_set_important_friend(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessSetImportantFriend(*relationpb.C2SSetImportantFriendProto, iface.HeroController)
		}).ProcessSetImportantFriend(&relationpb.C2SSetImportantFriendProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleRelation_c2s_cancel_important_friend(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.RelationModule().(interface {
			ProcessCancelImportantFriend(*relationpb.C2SCancelImportantFriendProto, iface.HeroController)
		}).ProcessCancelImportantFriend(&relationpb.C2SCancelImportantFriendProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerSecret_tower() *gm_group {
	group := &gm_group{
		tab: "secret_tower",
		handler: []*gm_handler{
			newStringHandler("request_team_count", " ", gm.handleSecret_tower_c2s_request_team_count),
			newStringHandler("request_team_list", " ", gm.handleSecret_tower_c2s_request_team_list),
			newStringHandler("create_team", " ", gm.handleSecret_tower_c2s_create_team),
			newStringHandler("join_team", " ", gm.handleSecret_tower_c2s_join_team),
			newStringHandler("leave_team", " ", gm.handleSecret_tower_c2s_leave_team),
			newStringHandler("kick_member", " ", gm.handleSecret_tower_c2s_kick_member),
			newStringHandler("move_member", " ", gm.handleSecret_tower_c2s_move_member),
			newStringHandler("update_member_pos", " ", gm.handleSecret_tower_c2s_update_member_pos),
			newStringHandler("change_mode", " ", gm.handleSecret_tower_c2s_change_mode),
			newStringHandler("invite", " ", gm.handleSecret_tower_c2s_invite),
			newStringHandler("invite_all", " ", gm.handleSecret_tower_c2s_invite_all),
			newStringHandler("request_invite_list", " ", gm.handleSecret_tower_c2s_request_invite_list),
			newStringHandler("request_team_detail", " ", gm.handleSecret_tower_c2s_request_team_detail),
			newStringHandler("start_challenge", " ", gm.handleSecret_tower_c2s_start_challenge),
			newStringHandler("quick_query_team_basic", " ", gm.handleSecret_tower_c2s_quick_query_team_basic),
			newStringHandler("change_guild_mode", " ", gm.handleSecret_tower_c2s_change_guild_mode),
			newStringHandler("list_record", " ", gm.handleSecret_tower_c2s_list_record),
			newStringHandler("team_talk", " ", gm.handleSecret_tower_c2s_team_talk),
		},
	}

	return group
}

func (gm *GmModule) handleSecret_tower_c2s_request_team_count(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessRequestTeamCount(iface.HeroController)
	}).ProcessRequestTeamCount(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_request_team_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessRequestTeamList(*secret_towerpb.C2SRequestTeamListProto, iface.HeroController)
		}).ProcessRequestTeamList(&secret_towerpb.C2SRequestTeamListProto{

			SecretTowerId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_create_team(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessCreateTeam(*secret_towerpb.C2SCreateTeamProto, iface.HeroController)
		}).ProcessCreateTeam(&secret_towerpb.C2SCreateTeamProto{

			SecretTowerId: parseInt32(strArray[0]),

			IsGuild: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_join_team(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessJoinTeam(*secret_towerpb.C2SJoinTeamProto, iface.HeroController)
		}).ProcessJoinTeam(&secret_towerpb.C2SJoinTeamProto{

			TeamId: parseInt32(strArray[0]),

			SecretTowerId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_leave_team(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessLeaveTeam(iface.HeroController)
	}).ProcessLeaveTeam(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_kick_member(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessKickMember(*secret_towerpb.C2SKickMemberProto, iface.HeroController)
		}).ProcessKickMember(&secret_towerpb.C2SKickMemberProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_move_member(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessMoveMember(*secret_towerpb.C2SMoveMemberProto, iface.HeroController)
		}).ProcessMoveMember(&secret_towerpb.C2SMoveMemberProto{

			Id: parseBytes(strArray[0]),

			Up: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_update_member_pos(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessUpdateMemberPos(*secret_towerpb.C2SUpdateMemberPosProto, iface.HeroController)
		}).ProcessUpdateMemberPos(&secret_towerpb.C2SUpdateMemberPosProto{

			Id: parseBytesArray(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_change_mode(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessChangeMode(*secret_towerpb.C2SChangeModeProto, iface.HeroController)
		}).ProcessChangeMode(&secret_towerpb.C2SChangeModeProto{

			Mode: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_invite(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessInvite(*secret_towerpb.C2SInviteProto, iface.HeroController)
		}).ProcessInvite(&secret_towerpb.C2SInviteProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_invite_all(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessInviteAll(*secret_towerpb.C2SInviteAllProto, iface.HeroController)
		}).ProcessInviteAll(&secret_towerpb.C2SInviteAllProto{

			Id: parseBytesArray(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_request_invite_list(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessRequestInviteList(iface.HeroController)
	}).ProcessRequestInviteList(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_request_team_detail(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessRequestTeamDetail(iface.HeroController)
	}).ProcessRequestTeamDetail(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_start_challenge(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessStartChallenge(iface.HeroController)
	}).ProcessStartChallenge(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_quick_query_team_basic(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessQuickQueryTeamBasic(*secret_towerpb.C2SQuickQueryTeamBasicProto, iface.HeroController)
		}).ProcessQuickQueryTeamBasic(&secret_towerpb.C2SQuickQueryTeamBasicProto{

			Ids: parseInt32Array(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleSecret_tower_c2s_change_guild_mode(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessChangeGuildMode(iface.HeroController)
	}).ProcessChangeGuildMode(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_list_record(amount string, hc iface.HeroController) {
	gm.modules.SecretTowerModule().(interface {
		ProcessListRecord(iface.HeroController)
	}).ProcessListRecord(hc)
}
func (gm *GmModule) handleSecret_tower_c2s_team_talk(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SecretTowerModule().(interface {
			ProcessTeamTalk(*secret_towerpb.C2STeamTalkProto, iface.HeroController)
		}).ProcessTeamTalk(&secret_towerpb.C2STeamTalkProto{

			WordsId: parseInt32(strArray[0]),

			Text: parseString(strArray[1]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerShop() *gm_group {
	group := &gm_group{
		tab: "shop",
		handler: []*gm_handler{
			newStringHandler("buy_goods", " ", gm.handleShop_c2s_buy_goods),
			newStringHandler("buy_black_market_goods", " ", gm.handleShop_c2s_buy_black_market_goods),
			newStringHandler("refresh_black_market_goods", " ", gm.handleShop_c2s_refresh_black_market_goods),
		},
	}

	return group
}

func (gm *GmModule) handleShop_c2s_buy_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ShopModule().(interface {
			ProcessBuyGoods(*shoppb.C2SBuyGoodsProto, iface.HeroController)
		}).ProcessBuyGoods(&shoppb.C2SBuyGoodsProto{

			Id: parseInt32(strArray[0]),

			Count: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleShop_c2s_buy_black_market_goods(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ShopModule().(interface {
			ProcessBuyBlackMarketGoods(*shoppb.C2SBuyBlackMarketGoodsProto, iface.HeroController)
		}).ProcessBuyBlackMarketGoods(&shoppb.C2SBuyBlackMarketGoodsProto{

			Index: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleShop_c2s_refresh_black_market_goods(amount string, hc iface.HeroController) {
	gm.modules.ShopModule().(interface {
		ProcessRefreshBlackMarketGoods(iface.HeroController)
	}).ProcessRefreshBlackMarketGoods(hc)
}

func (gm *GmModule) initHandlerStrategy() *gm_group {
	group := &gm_group{
		tab: "strategy",
		handler: []*gm_handler{
			newStringHandler("use_stratagem", " ", gm.handleStrategy_c2s_use_stratagem),
		},
	}

	return group
}

func (gm *GmModule) handleStrategy_c2s_use_stratagem(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.StrategyModule().(interface {
			ProcessUseStratagem(*strategypb.C2SUseStratagemProto, iface.HeroController)
		}).ProcessUseStratagem(&strategypb.C2SUseStratagemProto{

			Id: parseInt32(strArray[0]),

			Target: parseBytes(strArray[1]),

			DataId: parseInt32(strArray[2]),

			PosX: parseInt32(strArray[3]),

			PosY: parseInt32(strArray[4]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerStress() *gm_group {
	group := &gm_group{
		tab: "stress",
		handler: []*gm_handler{
			newStringHandler("robot_ping", " ", gm.handleStress_c2s_robot_ping),
		},
	}

	return group
}

func (gm *GmModule) handleStress_c2s_robot_ping(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.StressModule().(interface {
			Ping(*stresspb.C2SRobotPingProto, iface.HeroController)
		}).Ping(&stresspb.C2SRobotPingProto{

			Time: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerSurvey() *gm_group {
	group := &gm_group{
		tab: "survey",
		handler: []*gm_handler{
			newStringHandler("complete", " ", gm.handleSurvey_c2s_complete),
		},
	}

	return group
}

func (gm *GmModule) handleSurvey_c2s_complete(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.SurveyModule().(interface {
			ProcessComplete(*surveypb.C2SCompleteProto, iface.HeroController)
		}).ProcessComplete(&surveypb.C2SCompleteProto{

			ToDel: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerTag() *gm_group {
	group := &gm_group{
		tab: "tag",
		handler: []*gm_handler{
			newStringHandler("add_or_update_tag", " ", gm.handleTag_c2s_add_or_update_tag),
			newStringHandler("delete_tag", " ", gm.handleTag_c2s_delete_tag),
		},
	}

	return group
}

func (gm *GmModule) handleTag_c2s_add_or_update_tag(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TagModule().(interface {
			ProcessAddOrUpdateTag(*tagpb.C2SAddOrUpdateTagProto, iface.HeroController)
		}).ProcessAddOrUpdateTag(&tagpb.C2SAddOrUpdateTagProto{

			Id: parseBytes(strArray[0]),

			Tag: parseString(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleTag_c2s_delete_tag(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TagModule().(interface {
			ProcessDeleteTag(*tagpb.C2SDeleteTagProto, iface.HeroController)
		}).ProcessDeleteTag(&tagpb.C2SDeleteTagProto{

			Tags: parseStringArray(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerTask() *gm_group {
	group := &gm_group{
		tab: "task",
		handler: []*gm_handler{
			newStringHandler("collect_task_prize", " ", gm.handleTask_c2s_collect_task_prize),
			newStringHandler("collect_task_box_prize", " ", gm.handleTask_c2s_collect_task_box_prize),
			newStringHandler("collect_ba_ye_stage_prize", " ", gm.handleTask_c2s_collect_ba_ye_stage_prize),
			newStringHandler("collect_active_degree_prize", " ", gm.handleTask_c2s_collect_active_degree_prize),
			newStringHandler("collect_achieve_star_prize", " ", gm.handleTask_c2s_collect_achieve_star_prize),
			newStringHandler("change_select_show_achieve", " ", gm.handleTask_c2s_change_select_show_achieve),
			newStringHandler("collect_bwzl_prize", " ", gm.handleTask_c2s_collect_bwzl_prize),
			newStringHandler("view_other_achieve_task_list", " ", gm.handleTask_c2s_view_other_achieve_task_list),
			newStringHandler("get_troop_title_fight_amount", " ", gm.handleTask_c2s_get_troop_title_fight_amount),
			newStringHandler("upgrade_title", " ", gm.handleTask_c2s_upgrade_title),
			newStringHandler("complete_bool_task", " ", gm.handleTask_c2s_complete_bool_task),
		},
	}

	return group
}

func (gm *GmModule) handleTask_c2s_collect_task_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCollectTaskPrize(*taskpb.C2SCollectTaskPrizeProto, iface.HeroController)
		}).ProcessCollectTaskPrize(&taskpb.C2SCollectTaskPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_collect_task_box_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCollectTaskBoxPrize(*taskpb.C2SCollectTaskBoxPrizeProto, iface.HeroController)
		}).ProcessCollectTaskBoxPrize(&taskpb.C2SCollectTaskBoxPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_collect_ba_ye_stage_prize(amount string, hc iface.HeroController) {
	gm.modules.TaskModule().(interface {
		ProcessCollectBaYeStagePrize(iface.HeroController)
	}).ProcessCollectBaYeStagePrize(hc)
}
func (gm *GmModule) handleTask_c2s_collect_active_degree_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCollectActiveDegreePrize(*taskpb.C2SCollectActiveDegreePrizeProto, iface.HeroController)
		}).ProcessCollectActiveDegreePrize(&taskpb.C2SCollectActiveDegreePrizeProto{

			CollectIndex: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_collect_achieve_star_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCollectAchieveStarPrize(*taskpb.C2SCollectAchieveStarPrizeProto, iface.HeroController)
		}).ProcessCollectAchieveStarPrize(&taskpb.C2SCollectAchieveStarPrizeProto{

			StarCount: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_change_select_show_achieve(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessChangeSelectShowAchieve(*taskpb.C2SChangeSelectShowAchieveProto, iface.HeroController)
		}).ProcessChangeSelectShowAchieve(&taskpb.C2SChangeSelectShowAchieveProto{

			AchieveType: parseInt32(strArray[0]),

			AddOrRemove: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_collect_bwzl_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCollectBwzlPrize(*taskpb.C2SCollectBwzlPrizeProto, iface.HeroController)
		}).ProcessCollectBwzlPrize(&taskpb.C2SCollectBwzlPrizeProto{

			CompleteCount: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_view_other_achieve_task_list(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessViewOtherAchieveTaskList(*taskpb.C2SViewOtherAchieveTaskListProto, iface.HeroController)
		}).ProcessViewOtherAchieveTaskList(&taskpb.C2SViewOtherAchieveTaskListProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_get_troop_title_fight_amount(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessGetUpgradeTitleFightAmount(*taskpb.C2SGetTroopTitleFightAmountProto, iface.HeroController)
		}).ProcessGetUpgradeTitleFightAmount(&taskpb.C2SGetTroopTitleFightAmountProto{

			TroopIndex: parseInt32(strArray[0]),

			TitleId: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleTask_c2s_upgrade_title(amount string, hc iface.HeroController) {
	gm.modules.TaskModule().(interface {
		ProcessUpgradeTitle(iface.HeroController)
	}).ProcessUpgradeTitle(hc)
}
func (gm *GmModule) handleTask_c2s_complete_bool_task(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TaskModule().(interface {
			ProcessCompleteBoolTask(*taskpb.C2SCompleteBoolTaskProto, iface.HeroController)
		}).ProcessCompleteBoolTask(&taskpb.C2SCompleteBoolTaskProto{

			BoolType: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerTeach() *gm_group {
	group := &gm_group{
		tab: "teach",
		handler: []*gm_handler{
			newStringHandler("fight", " ", gm.handleTeach_c2s_fight),
			newStringHandler("collect_prize", " ", gm.handleTeach_c2s_collect_prize),
		},
	}

	return group
}

func (gm *GmModule) handleTeach_c2s_fight(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TeachModule().(interface {
			ProcessFight(*teachpb.C2SFightProto, iface.HeroController)
		}).ProcessFight(&teachpb.C2SFightProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTeach_c2s_collect_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TeachModule().(interface {
			ProcessCollectPrize(*teachpb.C2SCollectPrizeProto, iface.HeroController)
		}).ProcessCollectPrize(&teachpb.C2SCollectPrizeProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerTower() *gm_group {
	group := &gm_group{
		tab: "tower",
		handler: []*gm_handler{
			newStringHandler("challenge", " ", gm.handleTower_c2s_challenge),
			newStringHandler("auto_challenge", " ", gm.handleTower_c2s_auto_challenge),
			newStringHandler("collect_box", " ", gm.handleTower_c2s_collect_box),
			newStringHandler("list_pass_replay", " ", gm.handleTower_c2s_list_pass_replay),
		},
	}

	return group
}

func (gm *GmModule) handleTower_c2s_challenge(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TowerModule().(interface {
			ProcessChallenge(*towerpb.C2SChallengeProto, iface.HeroController)
		}).ProcessChallenge(&towerpb.C2SChallengeProto{

			Floor: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTower_c2s_auto_challenge(amount string, hc iface.HeroController) {
	gm.modules.TowerModule().(interface {
		ProcessAutoChallenge(iface.HeroController)
	}).ProcessAutoChallenge(hc)
}
func (gm *GmModule) handleTower_c2s_collect_box(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TowerModule().(interface {
			ProcessCollectBox(*towerpb.C2SCollectBoxProto, iface.HeroController)
		}).ProcessCollectBox(&towerpb.C2SCollectBoxProto{

			BoxFloor: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleTower_c2s_list_pass_replay(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.TowerModule().(interface {
			ProcessListPassReplay(*towerpb.C2SListPassReplayProto, iface.HeroController)
		}).ProcessListPassReplay(&towerpb.C2SListPassReplayProto{

			Floor: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerVip() *gm_group {
	group := &gm_group{
		tab: "vip",
		handler: []*gm_handler{
			newStringHandler("vip_collect_daily_prize", " ", gm.handleVip_c2s_vip_collect_daily_prize),
			newStringHandler("vip_collect_level_prize", " ", gm.handleVip_c2s_vip_collect_level_prize),
			newStringHandler("vip_buy_dungeon_times", " ", gm.handleVip_c2s_vip_buy_dungeon_times),
		},
	}

	return group
}

func (gm *GmModule) handleVip_c2s_vip_collect_daily_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.VipModule().(interface {
			ProcessVipCollectDailyPrize(*vippb.C2SVipCollectDailyPrizeProto, iface.HeroController)
		}).ProcessVipCollectDailyPrize(&vippb.C2SVipCollectDailyPrizeProto{

			VipLevel: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleVip_c2s_vip_collect_level_prize(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.VipModule().(interface {
			ProcessVipCollectLevelPrize(*vippb.C2SVipCollectLevelPrizeProto, iface.HeroController)
		}).ProcessVipCollectLevelPrize(&vippb.C2SVipCollectLevelPrizeProto{

			VipLevel: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleVip_c2s_vip_buy_dungeon_times(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.VipModule().(interface {
			ProcessVipBuyDungeonTimes(*vippb.C2SVipBuyDungeonTimesProto, iface.HeroController)
		}).ProcessVipBuyDungeonTimes(&vippb.C2SVipBuyDungeonTimesProto{

			DungeonId: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerXiongnu() *gm_group {
	group := &gm_group{
		tab: "xiongnu",
		handler: []*gm_handler{
			newStringHandler("set_defender", " ", gm.handleXiongnu_c2s_set_defender),
			newStringHandler("start", " ", gm.handleXiongnu_c2s_start),
			newStringHandler("troop_info", " ", gm.handleXiongnu_c2s_troop_info),
			newStringHandler("get_xiong_nu_npc_base_info", " ", gm.handleXiongnu_c2s_get_xiong_nu_npc_base_info),
			newStringHandler("get_defenser_fight_amount", " ", gm.handleXiongnu_c2s_get_defenser_fight_amount),
			newStringHandler("get_xiong_nu_fight_info", " ", gm.handleXiongnu_c2s_get_xiong_nu_fight_info),
		},
	}

	return group
}

func (gm *GmModule) handleXiongnu_c2s_set_defender(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XiongNuModule().(interface {
			ProcessSetDefender(*xiongnupb.C2SSetDefenderProto, iface.HeroController)
		}).ProcessSetDefender(&xiongnupb.C2SSetDefenderProto{

			Id: parseBytes(strArray[0]),

			ToSet: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleXiongnu_c2s_start(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XiongNuModule().(interface {
			ProcessStart(*xiongnupb.C2SStartProto, iface.HeroController)
		}).ProcessStart(&xiongnupb.C2SStartProto{

			Level: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleXiongnu_c2s_troop_info(amount string, hc iface.HeroController) {
	gm.modules.XiongNuModule().(interface {
		ProcessTroopInfo(iface.HeroController)
	}).ProcessTroopInfo(hc)
}
func (gm *GmModule) handleXiongnu_c2s_get_xiong_nu_npc_base_info(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XiongNuModule().(interface {
			ProcessGetXiongNuNpcBaseInfo(*xiongnupb.C2SGetXiongNuNpcBaseInfoProto, iface.HeroController)
		}).ProcessGetXiongNuNpcBaseInfo(&xiongnupb.C2SGetXiongNuNpcBaseInfoProto{

			GuildId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleXiongnu_c2s_get_defenser_fight_amount(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XiongNuModule().(interface {
			ProcessGetDefenserFightAmount(*xiongnupb.C2SGetDefenserFightAmountProto, iface.HeroController)
		}).ProcessGetDefenserFightAmount(&xiongnupb.C2SGetDefenserFightAmountProto{

			Version: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleXiongnu_c2s_get_xiong_nu_fight_info(amount string, hc iface.HeroController) {
	gm.modules.XiongNuModule().(interface {
		ProcessGetXiongNuFightInfo(iface.HeroController)
	}).ProcessGetXiongNuFightInfo(hc)
}

func (gm *GmModule) initHandlerXuanyuan() *gm_group {
	group := &gm_group{
		tab: "xuanyuan",
		handler: []*gm_handler{
			newStringHandler("self_info", " ", gm.handleXuanyuan_c2s_self_info),
			newStringHandler("list_target", " ", gm.handleXuanyuan_c2s_list_target),
			newStringHandler("query_target_troop", " ", gm.handleXuanyuan_c2s_query_target_troop),
			newStringHandler("challenge", " ", gm.handleXuanyuan_c2s_challenge),
			newStringHandler("list_record", " ", gm.handleXuanyuan_c2s_list_record),
			newStringHandler("collect_rank_prize", " ", gm.handleXuanyuan_c2s_collect_rank_prize),
		},
	}

	return group
}

func (gm *GmModule) handleXuanyuan_c2s_self_info(amount string, hc iface.HeroController) {
	gm.modules.XuanyuanModule().(interface {
		ProcessSelfInfo(iface.HeroController)
	}).ProcessSelfInfo(hc)
}
func (gm *GmModule) handleXuanyuan_c2s_list_target(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XuanyuanModule().(interface {
			ProcessListTarget(*xuanyuanpb.C2SListTargetProto, iface.HeroController)
		}).ProcessListTarget(&xuanyuanpb.C2SListTargetProto{

			RangeId: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleXuanyuan_c2s_query_target_troop(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XuanyuanModule().(interface {
			ProcessQueryTargetTroop(*xuanyuanpb.C2SQueryTargetTroopProto, iface.HeroController)
		}).ProcessQueryTargetTroop(&xuanyuanpb.C2SQueryTargetTroopProto{

			Id: parseBytes(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleXuanyuan_c2s_challenge(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XuanyuanModule().(interface {
			ProcessChallenge(*xuanyuanpb.C2SChallengeProto, iface.HeroController)
		}).ProcessChallenge(&xuanyuanpb.C2SChallengeProto{

			Id: parseBytes(strArray[0]),

			Version: parseInt32(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleXuanyuan_c2s_list_record(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.XuanyuanModule().(interface {
			ProcessListRecord(*xuanyuanpb.C2SListRecordProto, iface.HeroController)
		}).ProcessListRecord(&xuanyuanpb.C2SListRecordProto{

			Id: parseInt32(strArray[0]),

			Up: parseBool(strArray[1]),
		}, hc)
	}
}
func (gm *GmModule) handleXuanyuan_c2s_collect_rank_prize(amount string, hc iface.HeroController) {
	gm.modules.XuanyuanModule().(interface {
		ProcessCollectRankPrize(iface.HeroController)
	}).ProcessCollectRankPrize(hc)
}

func (gm *GmModule) initHandlerZhanjiang() *gm_group {
	group := &gm_group{
		tab: "zhanjiang",
		handler: []*gm_handler{
			newStringHandler("open", " ", gm.handleZhanjiang_c2s_open),
			newStringHandler("give_up", " ", gm.handleZhanjiang_c2s_give_up),
			newStringHandler("update_captain", " ", gm.handleZhanjiang_c2s_update_captain),
			newStringHandler("challenge", " ", gm.handleZhanjiang_c2s_challenge),
		},
	}

	return group
}

func (gm *GmModule) handleZhanjiang_c2s_open(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ZhanJiangModule().(interface {
			ProcessOpen(*zhanjiangpb.C2SOpenProto, iface.HeroController)
		}).ProcessOpen(&zhanjiangpb.C2SOpenProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleZhanjiang_c2s_give_up(amount string, hc iface.HeroController) {
	gm.modules.ZhanJiangModule().(interface {
		ProcessGiveUp(iface.HeroController)
	}).ProcessGiveUp(hc)
}
func (gm *GmModule) handleZhanjiang_c2s_update_captain(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ZhanJiangModule().(interface {
			ProcessUpdateCaptain(*zhanjiangpb.C2SUpdateCaptainProto, iface.HeroController)
		}).ProcessUpdateCaptain(&zhanjiangpb.C2SUpdateCaptainProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleZhanjiang_c2s_challenge(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ZhanJiangModule().(interface {
			ProcessChallenge(*zhanjiangpb.C2SChallengeProto, iface.HeroController)
		}).ProcessChallenge(&zhanjiangpb.C2SChallengeProto{

			PassCount: parseInt32(strArray[0]),
		}, hc)
	}
}

func (gm *GmModule) initHandlerZhengwu() *gm_group {
	group := &gm_group{
		tab: "zhengwu",
		handler: []*gm_handler{
			newStringHandler("start", " ", gm.handleZhengwu_c2s_start),
			newStringHandler("collect", " ", gm.handleZhengwu_c2s_collect),
			newStringHandler("yuanbao_complete", " ", gm.handleZhengwu_c2s_yuanbao_complete),
			newStringHandler("yuanbao_refresh", " ", gm.handleZhengwu_c2s_yuanbao_refresh),
			newStringHandler("vip_collect", " ", gm.handleZhengwu_c2s_vip_collect),
		},
	}

	return group
}

func (gm *GmModule) handleZhengwu_c2s_start(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ZhengWuModule().(interface {
			ProcessStart(*zhengwupb.C2SStartProto, iface.HeroController)
		}).ProcessStart(&zhengwupb.C2SStartProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}
func (gm *GmModule) handleZhengwu_c2s_collect(amount string, hc iface.HeroController) {
	gm.modules.ZhengWuModule().(interface {
		ProcessCollect(iface.HeroController)
	}).ProcessCollect(hc)
}
func (gm *GmModule) handleZhengwu_c2s_yuanbao_complete(amount string, hc iface.HeroController) {
	gm.modules.ZhengWuModule().(interface {
		ProcessYuanBaoComplete(iface.HeroController)
	}).ProcessYuanBaoComplete(hc)
}
func (gm *GmModule) handleZhengwu_c2s_yuanbao_refresh(amount string, hc iface.HeroController) {
	gm.modules.ZhengWuModule().(interface {
		ProcessYuanBaoRefresh(iface.HeroController)
	}).ProcessYuanBaoRefresh(hc)
}
func (gm *GmModule) handleZhengwu_c2s_vip_collect(amount string, hc iface.HeroController) {
	strArray := strings.Split(strings.TrimSpace(amount), " ")

	if len(strArray) > 0 {
		gm.modules.ZhengWuModule().(interface {
			ProcessVipCollect(*zhengwupb.C2SVipCollectProto, iface.HeroController)
		}).ProcessVipCollect(&zhengwupb.C2SVipCollectProto{

			Id: parseInt32(strArray[0]),
		}, hc)
	}
}

// 有问题时, 注释上面的, 反注释下面的
/*
----------- 自动识别分割线 -----------
package gm

func (gm *GmModule) initModuleHandler() []*gm_group {
	return []*gm_group{}
}
----------- 自动识别分割线 -----------
*/
