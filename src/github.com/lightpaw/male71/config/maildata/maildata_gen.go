// AUTO_GEN, DONT MODIFY!!!
package maildata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var _ = strings.ToUpper("")      // import strings
var _ = strconv.IntSize          // import strconv
var _ = shared_proto.Int32Pair{} // import shared_proto
var _ = errors.Errorf("")        // import errors
var _ = time.Second              // import time

// start with MailData ----------------------------------

func LoadMailData(gos *config.GameObjects) (map[string]*MailData, map[*MailData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MailDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*MailData, len(lIsT))
	pArSeRmAp := make(map[*MailData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMailData) {
			continue
		}

		dAtA, err := NewMailData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMailData(dAtAmAp map[*MailData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MailDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMailDataKeyArray(datas []*MailData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMailData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MailData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMailData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MailData{}

	dAtA.Id = pArSeR.String("id")
	dAtA.Icon = pArSeR.Uint64("icon")
	dAtA.Keep = pArSeR.Bool("keep")
	dAtA.Tag = 0
	if pArSeR.KeyExist("tag") {
		dAtA.Tag = pArSeR.Uint64("tag")
	}

	dAtA.Image = pArSeR.String("image")
	dAtA.ImageWord = 0
	if pArSeR.KeyExist("image_word") {
		dAtA.ImageWord = pArSeR.Uint64("image_word")
	}

	// releated field: Prize

	// i18n fields
	dAtA.Title = i18n.NewI18nRef(fIlEnAmE, "title", dAtA.Id, pArSeR.String("title"))
	// i18n fields
	dAtA.SubTitle = i18n.NewI18nRef(fIlEnAmE, "sub_title", dAtA.Id, pArSeR.String("sub_title"))
	// i18n fields
	dAtA.Text = i18n.NewI18nRef(fIlEnAmE, "text", dAtA.Id, pArSeR.String("text"))
	// i18n fields
	dAtA.Desc = i18n.NewI18nRef(fIlEnAmE, "desc", dAtA.Id, pArSeR.String("desc"))

	return dAtA, nil
}

var vAlIdAtOrMailData = map[string]*config.Validator{

	"id":         config.ParseValidator("string", "", false, nil, nil),
	"icon":       config.ParseValidator("uint", "", false, nil, nil),
	"title":      config.ParseValidator("string", "", false, nil, nil),
	"sub_title":  config.ParseValidator("string", "", false, nil, nil),
	"text":       config.ParseValidator("string", "", false, nil, nil),
	"keep":       config.ParseValidator("bool", "", false, nil, nil),
	"tag":        config.ParseValidator("int", "", false, nil, []string{"0"}),
	"desc":       config.ParseValidator("string", "", false, nil, nil),
	"image":      config.ParseValidator("string", "", false, nil, nil),
	"image_word": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"prize":      config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *MailData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil && pArSeR.Int("prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with MailHelp ----------------------------------

func LoadMailHelp(gos *config.GameObjects) (*MailHelp, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MailHelpPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewMailHelp(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMailHelp(gos *config.GameObjects, dAtA *MailHelp, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MailHelpPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewMailHelp(fIlEnAmE string, pArSeR *config.ObjectParser) (*MailHelp, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMailHelp)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MailHelp{}

	// releated field: AssemblyAdaFail
	// releated field: AssemblyAdaSuccess
	// releated field: AssemblyAddFail
	// releated field: AssemblyAddSuccess
	// releated field: AssemblyAsaFail
	// releated field: AssemblyAsdFail
	// releated field: AssemblyAsdSuccess
	// releated field: AssemblyAssFail
	// releated field: AssemblyAssSuccess
	// releated field: AssemblyDoneAttacker
	// releated field: AssemblyDoneDefenser
	// releated field: AssemblyExpelAttackerFail
	// releated field: AssemblyExpelAttackerSuccess
	// releated field: AssemblyExpelDefenserFail
	// releated field: AssemblyExpelDefenserSuccess
	// releated field: AssemblySaaFail
	// releated field: AssemblySaaSuccess
	// releated field: AssemblySadFail
	// releated field: AssemblySadSuccess
	// releated field: AssemblySasFail
	// releated field: AssemblySasSuccess
	// releated field: BaiZhanJunXianLevelDown
	// releated field: BaiZhanJunXianLevelKeep
	// releated field: BaiZhanJunXianLevelMaxKeep
	// releated field: BaiZhanJunXianLevelUp
	// releated field: BuyDaillyBargainSuccess
	// releated field: CopyDefenserBeenKilled
	// releated field: CopyDefenserExpired
	// releated field: CountryChangeNameFail
	// releated field: CountryChangeNameSucc
	// releated field: CountryDestroy
	// releated field: CountryOfficialAppoint
	// releated field: CountryOfficialDepose
	// releated field: FirstBeenRob
	// releated field: FirstProtectRemoved
	// releated field: FirstTowerFail
	// releated field: GuildBeKickedOut
	// releated field: GuildBigBox
	// releated field: GuildChangeCountry
	// releated field: GuildClassLevelDown
	// releated field: GuildClassLevelUp
	// releated field: GuildConveneLeader
	// releated field: GuildConveneOfficer
	// releated field: GuildEventPrize
	// releated field: GuildFirstJoin
	// releated field: GuildFirstJoinNpc
	// releated field: GuildLeaderChanged
	// releated field: GuildLeaderGenerated
	// releated field: GuildMcWarApplyAtkFail
	// releated field: GuildMcWarApplyAtkSucc
	// releated field: GuildPaySalary
	// releated field: GuildReceiveYinliangFromGuild
	// releated field: GuildSendYinliangToMember
	// releated field: GuildTaskEvaluatePrize
	// releated field: HebiBeRobbedPrize
	// releated field: HebiCompletePrize
	// releated field: HebiCopyBeRobbedPrize
	// releated field: HebiCopyCompletePrize
	// releated field: HebiRobPrize
	// releated field: HebiRoomBeenRobbed
	// releated field: McBuildGuildMemberPrize
	// releated field: ReportAdaFail
	// releated field: ReportAdaSuccess
	// releated field: ReportAddFail
	// releated field: ReportAddSuccess
	// releated field: ReportAsaFail
	// releated field: ReportAsaSuccess
	// releated field: ReportAsdFail
	// releated field: ReportAsdSuccess
	// releated field: ReportAssFail
	// releated field: ReportAssSuccess
	// releated field: ReportBaozRepatriateMoving
	// releated field: ReportBaozRepatriateRobber
	// releated field: ReportDoneAttacker
	// releated field: ReportDoneBaoz
	// releated field: ReportDoneBaozBack
	// releated field: ReportDoneDefenser
	// releated field: ReportExpelAttackerFail
	// releated field: ReportExpelAttackerSuccess
	// releated field: ReportExpelDefenserFail
	// releated field: ReportExpelDefenserSuccess
	// releated field: ReportSaaFail
	// releated field: ReportSaaSuccess
	// releated field: ReportSadFail
	// releated field: ReportSadSuccess
	// releated field: ReportSasFail
	// releated field: ReportSasSuccess
	// releated field: ReportWatchAttacker
	// releated field: ReportWatchDefenser
	// releated field: SurveyMail
	// releated field: SystemCompensation
	// releated field: XiongNuResistSuc
	// releated field: XiongNuScore

	return dAtA, nil
}

var vAlIdAtOrMailHelp = map[string]*config.Validator{

	"assembly_ada_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyAdaFail"}),
	"assembly_ada_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblyAdaSuccess"}),
	"assembly_add_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyAddFail"}),
	"assembly_add_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblyAddSuccess"}),
	"assembly_asa_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyAsaFail"}),
	"assembly_asd_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyAsdFail"}),
	"assembly_asd_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblyAsdSuccess"}),
	"assembly_ass_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyAssFail"}),
	"assembly_ass_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblyAssSuccess"}),
	"assembly_done_attacker":            config.ParseValidator("string", "", false, nil, []string{"AssemblyDoneAttacker"}),
	"assembly_done_defenser":            config.ParseValidator("string", "", false, nil, []string{"AssemblyDoneDefenser"}),
	"assembly_expel_attacker_fail":      config.ParseValidator("string", "", false, nil, []string{"AssemblyExpelAttackerFail"}),
	"assembly_expel_attacker_success":   config.ParseValidator("string", "", false, nil, []string{"AssemblyExpelAttackerSuccess"}),
	"assembly_expel_defenser_fail":      config.ParseValidator("string", "", false, nil, []string{"AssemblyExpelDefenserFail"}),
	"assembly_expel_defenser_success":   config.ParseValidator("string", "", false, nil, []string{"AssemblyExpelDefenserSuccess"}),
	"assembly_saa_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblySaaFail"}),
	"assembly_saa_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblySaaSuccess"}),
	"assembly_sad_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblySadFail"}),
	"assembly_sad_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblySadSuccess"}),
	"assembly_sas_fail":                 config.ParseValidator("string", "", false, nil, []string{"AssemblySasFail"}),
	"assembly_sas_success":              config.ParseValidator("string", "", false, nil, []string{"AssemblySasSuccess"}),
	"bai_zhan_jun_xian_level_down":      config.ParseValidator("string", "", false, nil, []string{"BaiZhanJunXianLevelDown"}),
	"bai_zhan_jun_xian_level_keep":      config.ParseValidator("string", "", false, nil, []string{"BaiZhanJunXianLevelKeep"}),
	"bai_zhan_jun_xian_level_max_keep":  config.ParseValidator("string", "", false, nil, []string{"BaiZhanJunXianLevelMaxKeep"}),
	"bai_zhan_jun_xian_level_up":        config.ParseValidator("string", "", false, nil, []string{"BaiZhanJunXianLevelUp"}),
	"buy_dailly_bargain_success":        config.ParseValidator("string", "", false, nil, []string{"BuyDaillyBargainSuccess"}),
	"copy_defenser_been_killed":         config.ParseValidator("string", "", false, nil, []string{"CopyDefenserBeenKilled"}),
	"copy_defenser_expired":             config.ParseValidator("string", "", false, nil, []string{"CopyDefenserExpired"}),
	"country_change_name_fail":          config.ParseValidator("string", "", false, nil, []string{"CountryChangeNameFail"}),
	"country_change_name_succ":          config.ParseValidator("string", "", false, nil, []string{"CountryChangeNameSucc"}),
	"country_destroy":                   config.ParseValidator("string", "", false, nil, []string{"CountryDestroy"}),
	"country_official_appoint":          config.ParseValidator("string", "", false, nil, []string{"CountryOfficialAppoint"}),
	"country_official_depose":           config.ParseValidator("string", "", false, nil, []string{"CountryOfficialDepose"}),
	"first_been_rob":                    config.ParseValidator("string", "", false, nil, []string{"FirstBeenRob"}),
	"first_protect_removed":             config.ParseValidator("string", "", false, nil, []string{"FirstProtectRemoved"}),
	"first_tower_fail":                  config.ParseValidator("string", "", false, nil, []string{"FirstTowerFail"}),
	"guild_be_kicked_out":               config.ParseValidator("string", "", false, nil, []string{"GuildBeKickedOut"}),
	"guild_big_box":                     config.ParseValidator("string", "", false, nil, []string{"GuildBigBox"}),
	"guild_change_country":              config.ParseValidator("string", "", false, nil, []string{"GuildChangeCountry"}),
	"guild_class_level_down":            config.ParseValidator("string", "", false, nil, []string{"GuildClassLevelDown"}),
	"guild_class_level_up":              config.ParseValidator("string", "", false, nil, []string{"GuildClassLevelUp"}),
	"guild_convene_leader":              config.ParseValidator("string", "", false, nil, []string{"GuildConveneLeader"}),
	"guild_convene_officer":             config.ParseValidator("string", "", false, nil, []string{"GuildConveneOfficer"}),
	"guild_event_prize":                 config.ParseValidator("string", "", false, nil, []string{"GuildEventPrize"}),
	"guild_first_join":                  config.ParseValidator("string", "", false, nil, []string{"GuildFirstJoin"}),
	"guild_first_join_npc":              config.ParseValidator("string", "", false, nil, []string{"GuildFirstJoinNpc"}),
	"guild_leader_changed":              config.ParseValidator("string", "", false, nil, []string{"GuildLeaderChanged"}),
	"guild_leader_generated":            config.ParseValidator("string", "", false, nil, []string{"GuildLeaderGenerated"}),
	"guild_mc_war_apply_atk_fail":       config.ParseValidator("string", "", false, nil, []string{"GuildMcWarApplyAtkFail"}),
	"guild_mc_war_apply_atk_succ":       config.ParseValidator("string", "", false, nil, []string{"GuildMcWarApplyAtkSucc"}),
	"guild_pay_salary":                  config.ParseValidator("string", "", false, nil, []string{"GuildPaySalary"}),
	"guild_receive_yinliang_from_guild": config.ParseValidator("string", "", false, nil, []string{"GuildReceiveYinliangFromGuild"}),
	"guild_send_yinliang_to_member":     config.ParseValidator("string", "", false, nil, []string{"GuildSendYinliangToMember"}),
	"guild_task_evaluate_prize":         config.ParseValidator("string", "", false, nil, []string{"GuildTaskEvaluatePrize"}),
	"hebi_be_robbed_prize":              config.ParseValidator("string", "", false, nil, []string{"HebiBeRobbedPrize"}),
	"hebi_complete_prize":               config.ParseValidator("string", "", false, nil, []string{"HebiCompletePrize"}),
	"hebi_copy_be_robbed_prize":         config.ParseValidator("string", "", false, nil, []string{"HebiCopyBeRobbedPrize"}),
	"hebi_copy_complete_prize":          config.ParseValidator("string", "", false, nil, []string{"HebiCopyCompletePrize"}),
	"hebi_rob_prize":                    config.ParseValidator("string", "", false, nil, []string{"HebiRobPrize"}),
	"hebi_room_been_robbed":             config.ParseValidator("string", "", false, nil, []string{"HebiRoomBeenRobbed"}),
	"mc_build_guild_member_prize":       config.ParseValidator("string", "", false, nil, []string{"McBuildGuildMemberPrize"}),
	"report_ada_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportAdaFail"}),
	"report_ada_success":                config.ParseValidator("string", "", false, nil, []string{"ReportAdaSuccess"}),
	"report_add_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportAddFail"}),
	"report_add_success":                config.ParseValidator("string", "", false, nil, []string{"ReportAddSuccess"}),
	"report_asa_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportAsaFail"}),
	"report_asa_success":                config.ParseValidator("string", "", false, nil, []string{"ReportAsaSuccess"}),
	"report_asd_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportAsdFail"}),
	"report_asd_success":                config.ParseValidator("string", "", false, nil, []string{"ReportAsdSuccess"}),
	"report_ass_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportAssFail"}),
	"report_ass_success":                config.ParseValidator("string", "", false, nil, []string{"ReportAssSuccess"}),
	"report_baoz_repatriate_moving":     config.ParseValidator("string", "", false, nil, []string{"ReportBaozRepatriateMoving"}),
	"report_baoz_repatriate_robber":     config.ParseValidator("string", "", false, nil, []string{"ReportBaozRepatriateRobber"}),
	"report_done_attacker":              config.ParseValidator("string", "", false, nil, []string{"ReportDoneAttacker"}),
	"report_done_baoz":                  config.ParseValidator("string", "", false, nil, []string{"ReportDoneBaoz"}),
	"report_done_baoz_back":             config.ParseValidator("string", "", false, nil, []string{"ReportDoneBaozBack"}),
	"report_done_defenser":              config.ParseValidator("string", "", false, nil, []string{"ReportDoneDefenser"}),
	"report_expel_attacker_fail":        config.ParseValidator("string", "", false, nil, []string{"ReportExpelAttackerFail"}),
	"report_expel_attacker_success":     config.ParseValidator("string", "", false, nil, []string{"ReportExpelAttackerSuccess"}),
	"report_expel_defenser_fail":        config.ParseValidator("string", "", false, nil, []string{"ReportExpelDefenserFail"}),
	"report_expel_defenser_success":     config.ParseValidator("string", "", false, nil, []string{"ReportExpelDefenserSuccess"}),
	"report_saa_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportSaaFail"}),
	"report_saa_success":                config.ParseValidator("string", "", false, nil, []string{"ReportSaaSuccess"}),
	"report_sad_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportSadFail"}),
	"report_sad_success":                config.ParseValidator("string", "", false, nil, []string{"ReportSadSuccess"}),
	"report_sas_fail":                   config.ParseValidator("string", "", false, nil, []string{"ReportSasFail"}),
	"report_sas_success":                config.ParseValidator("string", "", false, nil, []string{"ReportSasSuccess"}),
	"report_watch_attacker":             config.ParseValidator("string", "", false, nil, []string{"ReportWatchAttacker"}),
	"report_watch_defenser":             config.ParseValidator("string", "", false, nil, []string{"ReportWatchDefenser"}),
	"survey_mail":                       config.ParseValidator("string", "", false, nil, []string{"SurveyMail"}),
	"system_compensation":               config.ParseValidator("string", "", false, nil, []string{"SystemCompensation"}),
	"xiong_nu_resist_suc":               config.ParseValidator("string", "", false, nil, []string{"XiongNuResistSuc"}),
	"xiong_nu_score":                    config.ParseValidator("string", "", false, nil, []string{"XiongNuScore"}),
}

func (dAtA *MailHelp) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("assembly_ada_fail") {
		dAtA.AssemblyAdaFail = cOnFigS.GetMailData(pArSeR.String("assembly_ada_fail"))
	} else {
		dAtA.AssemblyAdaFail = cOnFigS.GetMailData("AssemblyAdaFail")
	}
	if dAtA.AssemblyAdaFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_ada_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_ada_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_ada_success") {
		dAtA.AssemblyAdaSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_ada_success"))
	} else {
		dAtA.AssemblyAdaSuccess = cOnFigS.GetMailData("AssemblyAdaSuccess")
	}
	if dAtA.AssemblyAdaSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_ada_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_ada_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_add_fail") {
		dAtA.AssemblyAddFail = cOnFigS.GetMailData(pArSeR.String("assembly_add_fail"))
	} else {
		dAtA.AssemblyAddFail = cOnFigS.GetMailData("AssemblyAddFail")
	}
	if dAtA.AssemblyAddFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_add_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_add_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_add_success") {
		dAtA.AssemblyAddSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_add_success"))
	} else {
		dAtA.AssemblyAddSuccess = cOnFigS.GetMailData("AssemblyAddSuccess")
	}
	if dAtA.AssemblyAddSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_add_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_add_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_asa_fail") {
		dAtA.AssemblyAsaFail = cOnFigS.GetMailData(pArSeR.String("assembly_asa_fail"))
	} else {
		dAtA.AssemblyAsaFail = cOnFigS.GetMailData("AssemblyAsaFail")
	}
	if dAtA.AssemblyAsaFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_asa_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_asa_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_asd_fail") {
		dAtA.AssemblyAsdFail = cOnFigS.GetMailData(pArSeR.String("assembly_asd_fail"))
	} else {
		dAtA.AssemblyAsdFail = cOnFigS.GetMailData("AssemblyAsdFail")
	}
	if dAtA.AssemblyAsdFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_asd_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_asd_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_asd_success") {
		dAtA.AssemblyAsdSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_asd_success"))
	} else {
		dAtA.AssemblyAsdSuccess = cOnFigS.GetMailData("AssemblyAsdSuccess")
	}
	if dAtA.AssemblyAsdSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_asd_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_asd_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_ass_fail") {
		dAtA.AssemblyAssFail = cOnFigS.GetMailData(pArSeR.String("assembly_ass_fail"))
	} else {
		dAtA.AssemblyAssFail = cOnFigS.GetMailData("AssemblyAssFail")
	}
	if dAtA.AssemblyAssFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_ass_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_ass_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_ass_success") {
		dAtA.AssemblyAssSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_ass_success"))
	} else {
		dAtA.AssemblyAssSuccess = cOnFigS.GetMailData("AssemblyAssSuccess")
	}
	if dAtA.AssemblyAssSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_ass_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_ass_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_done_attacker") {
		dAtA.AssemblyDoneAttacker = cOnFigS.GetMailData(pArSeR.String("assembly_done_attacker"))
	} else {
		dAtA.AssemblyDoneAttacker = cOnFigS.GetMailData("AssemblyDoneAttacker")
	}
	if dAtA.AssemblyDoneAttacker == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_done_attacker] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_done_attacker"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_done_defenser") {
		dAtA.AssemblyDoneDefenser = cOnFigS.GetMailData(pArSeR.String("assembly_done_defenser"))
	} else {
		dAtA.AssemblyDoneDefenser = cOnFigS.GetMailData("AssemblyDoneDefenser")
	}
	if dAtA.AssemblyDoneDefenser == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_done_defenser] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_done_defenser"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_expel_attacker_fail") {
		dAtA.AssemblyExpelAttackerFail = cOnFigS.GetMailData(pArSeR.String("assembly_expel_attacker_fail"))
	} else {
		dAtA.AssemblyExpelAttackerFail = cOnFigS.GetMailData("AssemblyExpelAttackerFail")
	}
	if dAtA.AssemblyExpelAttackerFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_expel_attacker_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_expel_attacker_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_expel_attacker_success") {
		dAtA.AssemblyExpelAttackerSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_expel_attacker_success"))
	} else {
		dAtA.AssemblyExpelAttackerSuccess = cOnFigS.GetMailData("AssemblyExpelAttackerSuccess")
	}
	if dAtA.AssemblyExpelAttackerSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_expel_attacker_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_expel_attacker_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_expel_defenser_fail") {
		dAtA.AssemblyExpelDefenserFail = cOnFigS.GetMailData(pArSeR.String("assembly_expel_defenser_fail"))
	} else {
		dAtA.AssemblyExpelDefenserFail = cOnFigS.GetMailData("AssemblyExpelDefenserFail")
	}
	if dAtA.AssemblyExpelDefenserFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_expel_defenser_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_expel_defenser_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_expel_defenser_success") {
		dAtA.AssemblyExpelDefenserSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_expel_defenser_success"))
	} else {
		dAtA.AssemblyExpelDefenserSuccess = cOnFigS.GetMailData("AssemblyExpelDefenserSuccess")
	}
	if dAtA.AssemblyExpelDefenserSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_expel_defenser_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_expel_defenser_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_saa_fail") {
		dAtA.AssemblySaaFail = cOnFigS.GetMailData(pArSeR.String("assembly_saa_fail"))
	} else {
		dAtA.AssemblySaaFail = cOnFigS.GetMailData("AssemblySaaFail")
	}
	if dAtA.AssemblySaaFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_saa_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_saa_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_saa_success") {
		dAtA.AssemblySaaSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_saa_success"))
	} else {
		dAtA.AssemblySaaSuccess = cOnFigS.GetMailData("AssemblySaaSuccess")
	}
	if dAtA.AssemblySaaSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_saa_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_saa_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_sad_fail") {
		dAtA.AssemblySadFail = cOnFigS.GetMailData(pArSeR.String("assembly_sad_fail"))
	} else {
		dAtA.AssemblySadFail = cOnFigS.GetMailData("AssemblySadFail")
	}
	if dAtA.AssemblySadFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_sad_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_sad_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_sad_success") {
		dAtA.AssemblySadSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_sad_success"))
	} else {
		dAtA.AssemblySadSuccess = cOnFigS.GetMailData("AssemblySadSuccess")
	}
	if dAtA.AssemblySadSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_sad_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_sad_success"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_sas_fail") {
		dAtA.AssemblySasFail = cOnFigS.GetMailData(pArSeR.String("assembly_sas_fail"))
	} else {
		dAtA.AssemblySasFail = cOnFigS.GetMailData("AssemblySasFail")
	}
	if dAtA.AssemblySasFail == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_sas_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_sas_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_sas_success") {
		dAtA.AssemblySasSuccess = cOnFigS.GetMailData(pArSeR.String("assembly_sas_success"))
	} else {
		dAtA.AssemblySasSuccess = cOnFigS.GetMailData("AssemblySasSuccess")
	}
	if dAtA.AssemblySasSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_sas_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_sas_success"), *pArSeR)
	}

	if pArSeR.KeyExist("bai_zhan_jun_xian_level_down") {
		dAtA.BaiZhanJunXianLevelDown = cOnFigS.GetMailData(pArSeR.String("bai_zhan_jun_xian_level_down"))
	} else {
		dAtA.BaiZhanJunXianLevelDown = cOnFigS.GetMailData("BaiZhanJunXianLevelDown")
	}
	if dAtA.BaiZhanJunXianLevelDown == nil {
		return errors.Errorf("%s 配置的关联字段[bai_zhan_jun_xian_level_down] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bai_zhan_jun_xian_level_down"), *pArSeR)
	}

	if pArSeR.KeyExist("bai_zhan_jun_xian_level_keep") {
		dAtA.BaiZhanJunXianLevelKeep = cOnFigS.GetMailData(pArSeR.String("bai_zhan_jun_xian_level_keep"))
	} else {
		dAtA.BaiZhanJunXianLevelKeep = cOnFigS.GetMailData("BaiZhanJunXianLevelKeep")
	}
	if dAtA.BaiZhanJunXianLevelKeep == nil {
		return errors.Errorf("%s 配置的关联字段[bai_zhan_jun_xian_level_keep] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bai_zhan_jun_xian_level_keep"), *pArSeR)
	}

	if pArSeR.KeyExist("bai_zhan_jun_xian_level_max_keep") {
		dAtA.BaiZhanJunXianLevelMaxKeep = cOnFigS.GetMailData(pArSeR.String("bai_zhan_jun_xian_level_max_keep"))
	} else {
		dAtA.BaiZhanJunXianLevelMaxKeep = cOnFigS.GetMailData("BaiZhanJunXianLevelMaxKeep")
	}
	if dAtA.BaiZhanJunXianLevelMaxKeep == nil {
		return errors.Errorf("%s 配置的关联字段[bai_zhan_jun_xian_level_max_keep] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bai_zhan_jun_xian_level_max_keep"), *pArSeR)
	}

	if pArSeR.KeyExist("bai_zhan_jun_xian_level_up") {
		dAtA.BaiZhanJunXianLevelUp = cOnFigS.GetMailData(pArSeR.String("bai_zhan_jun_xian_level_up"))
	} else {
		dAtA.BaiZhanJunXianLevelUp = cOnFigS.GetMailData("BaiZhanJunXianLevelUp")
	}
	if dAtA.BaiZhanJunXianLevelUp == nil {
		return errors.Errorf("%s 配置的关联字段[bai_zhan_jun_xian_level_up] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bai_zhan_jun_xian_level_up"), *pArSeR)
	}

	if pArSeR.KeyExist("buy_dailly_bargain_success") {
		dAtA.BuyDaillyBargainSuccess = cOnFigS.GetMailData(pArSeR.String("buy_dailly_bargain_success"))
	} else {
		dAtA.BuyDaillyBargainSuccess = cOnFigS.GetMailData("BuyDaillyBargainSuccess")
	}
	if dAtA.BuyDaillyBargainSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[buy_dailly_bargain_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("buy_dailly_bargain_success"), *pArSeR)
	}

	if pArSeR.KeyExist("copy_defenser_been_killed") {
		dAtA.CopyDefenserBeenKilled = cOnFigS.GetMailData(pArSeR.String("copy_defenser_been_killed"))
	} else {
		dAtA.CopyDefenserBeenKilled = cOnFigS.GetMailData("CopyDefenserBeenKilled")
	}
	if dAtA.CopyDefenserBeenKilled == nil {
		return errors.Errorf("%s 配置的关联字段[copy_defenser_been_killed] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("copy_defenser_been_killed"), *pArSeR)
	}

	if pArSeR.KeyExist("copy_defenser_expired") {
		dAtA.CopyDefenserExpired = cOnFigS.GetMailData(pArSeR.String("copy_defenser_expired"))
	} else {
		dAtA.CopyDefenserExpired = cOnFigS.GetMailData("CopyDefenserExpired")
	}
	if dAtA.CopyDefenserExpired == nil {
		return errors.Errorf("%s 配置的关联字段[copy_defenser_expired] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("copy_defenser_expired"), *pArSeR)
	}

	if pArSeR.KeyExist("country_change_name_fail") {
		dAtA.CountryChangeNameFail = cOnFigS.GetMailData(pArSeR.String("country_change_name_fail"))
	} else {
		dAtA.CountryChangeNameFail = cOnFigS.GetMailData("CountryChangeNameFail")
	}
	if dAtA.CountryChangeNameFail == nil {
		return errors.Errorf("%s 配置的关联字段[country_change_name_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_change_name_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("country_change_name_succ") {
		dAtA.CountryChangeNameSucc = cOnFigS.GetMailData(pArSeR.String("country_change_name_succ"))
	} else {
		dAtA.CountryChangeNameSucc = cOnFigS.GetMailData("CountryChangeNameSucc")
	}
	if dAtA.CountryChangeNameSucc == nil {
		return errors.Errorf("%s 配置的关联字段[country_change_name_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_change_name_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("country_destroy") {
		dAtA.CountryDestroy = cOnFigS.GetMailData(pArSeR.String("country_destroy"))
	} else {
		dAtA.CountryDestroy = cOnFigS.GetMailData("CountryDestroy")
	}
	if dAtA.CountryDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[country_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("country_official_appoint") {
		dAtA.CountryOfficialAppoint = cOnFigS.GetMailData(pArSeR.String("country_official_appoint"))
	} else {
		dAtA.CountryOfficialAppoint = cOnFigS.GetMailData("CountryOfficialAppoint")
	}
	if dAtA.CountryOfficialAppoint == nil {
		return errors.Errorf("%s 配置的关联字段[country_official_appoint] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_official_appoint"), *pArSeR)
	}

	if pArSeR.KeyExist("country_official_depose") {
		dAtA.CountryOfficialDepose = cOnFigS.GetMailData(pArSeR.String("country_official_depose"))
	} else {
		dAtA.CountryOfficialDepose = cOnFigS.GetMailData("CountryOfficialDepose")
	}
	if dAtA.CountryOfficialDepose == nil {
		return errors.Errorf("%s 配置的关联字段[country_official_depose] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_official_depose"), *pArSeR)
	}

	if pArSeR.KeyExist("first_been_rob") {
		dAtA.FirstBeenRob = cOnFigS.GetMailData(pArSeR.String("first_been_rob"))
	} else {
		dAtA.FirstBeenRob = cOnFigS.GetMailData("FirstBeenRob")
	}
	if dAtA.FirstBeenRob == nil {
		return errors.Errorf("%s 配置的关联字段[first_been_rob] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_been_rob"), *pArSeR)
	}

	if pArSeR.KeyExist("first_protect_removed") {
		dAtA.FirstProtectRemoved = cOnFigS.GetMailData(pArSeR.String("first_protect_removed"))
	} else {
		dAtA.FirstProtectRemoved = cOnFigS.GetMailData("FirstProtectRemoved")
	}
	if dAtA.FirstProtectRemoved == nil {
		return errors.Errorf("%s 配置的关联字段[first_protect_removed] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_protect_removed"), *pArSeR)
	}

	if pArSeR.KeyExist("first_tower_fail") {
		dAtA.FirstTowerFail = cOnFigS.GetMailData(pArSeR.String("first_tower_fail"))
	} else {
		dAtA.FirstTowerFail = cOnFigS.GetMailData("FirstTowerFail")
	}
	if dAtA.FirstTowerFail == nil {
		return errors.Errorf("%s 配置的关联字段[first_tower_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_tower_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_be_kicked_out") {
		dAtA.GuildBeKickedOut = cOnFigS.GetMailData(pArSeR.String("guild_be_kicked_out"))
	} else {
		dAtA.GuildBeKickedOut = cOnFigS.GetMailData("GuildBeKickedOut")
	}
	if dAtA.GuildBeKickedOut == nil {
		return errors.Errorf("%s 配置的关联字段[guild_be_kicked_out] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_be_kicked_out"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_big_box") {
		dAtA.GuildBigBox = cOnFigS.GetMailData(pArSeR.String("guild_big_box"))
	} else {
		dAtA.GuildBigBox = cOnFigS.GetMailData("GuildBigBox")
	}
	if dAtA.GuildBigBox == nil {
		return errors.Errorf("%s 配置的关联字段[guild_big_box] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_big_box"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_change_country") {
		dAtA.GuildChangeCountry = cOnFigS.GetMailData(pArSeR.String("guild_change_country"))
	} else {
		dAtA.GuildChangeCountry = cOnFigS.GetMailData("GuildChangeCountry")
	}
	if dAtA.GuildChangeCountry == nil {
		return errors.Errorf("%s 配置的关联字段[guild_change_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_change_country"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_class_level_down") {
		dAtA.GuildClassLevelDown = cOnFigS.GetMailData(pArSeR.String("guild_class_level_down"))
	} else {
		dAtA.GuildClassLevelDown = cOnFigS.GetMailData("GuildClassLevelDown")
	}
	if dAtA.GuildClassLevelDown == nil {
		return errors.Errorf("%s 配置的关联字段[guild_class_level_down] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_class_level_down"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_class_level_up") {
		dAtA.GuildClassLevelUp = cOnFigS.GetMailData(pArSeR.String("guild_class_level_up"))
	} else {
		dAtA.GuildClassLevelUp = cOnFigS.GetMailData("GuildClassLevelUp")
	}
	if dAtA.GuildClassLevelUp == nil {
		return errors.Errorf("%s 配置的关联字段[guild_class_level_up] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_class_level_up"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_convene_leader") {
		dAtA.GuildConveneLeader = cOnFigS.GetMailData(pArSeR.String("guild_convene_leader"))
	} else {
		dAtA.GuildConveneLeader = cOnFigS.GetMailData("GuildConveneLeader")
	}
	if dAtA.GuildConveneLeader == nil {
		return errors.Errorf("%s 配置的关联字段[guild_convene_leader] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_convene_leader"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_convene_officer") {
		dAtA.GuildConveneOfficer = cOnFigS.GetMailData(pArSeR.String("guild_convene_officer"))
	} else {
		dAtA.GuildConveneOfficer = cOnFigS.GetMailData("GuildConveneOfficer")
	}
	if dAtA.GuildConveneOfficer == nil {
		return errors.Errorf("%s 配置的关联字段[guild_convene_officer] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_convene_officer"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_event_prize") {
		dAtA.GuildEventPrize = cOnFigS.GetMailData(pArSeR.String("guild_event_prize"))
	} else {
		dAtA.GuildEventPrize = cOnFigS.GetMailData("GuildEventPrize")
	}
	if dAtA.GuildEventPrize == nil {
		return errors.Errorf("%s 配置的关联字段[guild_event_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_event_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_first_join") {
		dAtA.GuildFirstJoin = cOnFigS.GetMailData(pArSeR.String("guild_first_join"))
	} else {
		dAtA.GuildFirstJoin = cOnFigS.GetMailData("GuildFirstJoin")
	}
	if dAtA.GuildFirstJoin == nil {
		return errors.Errorf("%s 配置的关联字段[guild_first_join] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_first_join"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_first_join_npc") {
		dAtA.GuildFirstJoinNpc = cOnFigS.GetMailData(pArSeR.String("guild_first_join_npc"))
	} else {
		dAtA.GuildFirstJoinNpc = cOnFigS.GetMailData("GuildFirstJoinNpc")
	}
	if dAtA.GuildFirstJoinNpc == nil {
		return errors.Errorf("%s 配置的关联字段[guild_first_join_npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_first_join_npc"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_leader_changed") {
		dAtA.GuildLeaderChanged = cOnFigS.GetMailData(pArSeR.String("guild_leader_changed"))
	} else {
		dAtA.GuildLeaderChanged = cOnFigS.GetMailData("GuildLeaderChanged")
	}
	if dAtA.GuildLeaderChanged == nil {
		return errors.Errorf("%s 配置的关联字段[guild_leader_changed] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_leader_changed"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_leader_generated") {
		dAtA.GuildLeaderGenerated = cOnFigS.GetMailData(pArSeR.String("guild_leader_generated"))
	} else {
		dAtA.GuildLeaderGenerated = cOnFigS.GetMailData("GuildLeaderGenerated")
	}
	if dAtA.GuildLeaderGenerated == nil {
		return errors.Errorf("%s 配置的关联字段[guild_leader_generated] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_leader_generated"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_mc_war_apply_atk_fail") {
		dAtA.GuildMcWarApplyAtkFail = cOnFigS.GetMailData(pArSeR.String("guild_mc_war_apply_atk_fail"))
	} else {
		dAtA.GuildMcWarApplyAtkFail = cOnFigS.GetMailData("GuildMcWarApplyAtkFail")
	}
	if dAtA.GuildMcWarApplyAtkFail == nil {
		return errors.Errorf("%s 配置的关联字段[guild_mc_war_apply_atk_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_mc_war_apply_atk_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_mc_war_apply_atk_succ") {
		dAtA.GuildMcWarApplyAtkSucc = cOnFigS.GetMailData(pArSeR.String("guild_mc_war_apply_atk_succ"))
	} else {
		dAtA.GuildMcWarApplyAtkSucc = cOnFigS.GetMailData("GuildMcWarApplyAtkSucc")
	}
	if dAtA.GuildMcWarApplyAtkSucc == nil {
		return errors.Errorf("%s 配置的关联字段[guild_mc_war_apply_atk_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_mc_war_apply_atk_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_pay_salary") {
		dAtA.GuildPaySalary = cOnFigS.GetMailData(pArSeR.String("guild_pay_salary"))
	} else {
		dAtA.GuildPaySalary = cOnFigS.GetMailData("GuildPaySalary")
	}
	if dAtA.GuildPaySalary == nil {
		return errors.Errorf("%s 配置的关联字段[guild_pay_salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_pay_salary"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_receive_yinliang_from_guild") {
		dAtA.GuildReceiveYinliangFromGuild = cOnFigS.GetMailData(pArSeR.String("guild_receive_yinliang_from_guild"))
	} else {
		dAtA.GuildReceiveYinliangFromGuild = cOnFigS.GetMailData("GuildReceiveYinliangFromGuild")
	}
	if dAtA.GuildReceiveYinliangFromGuild == nil {
		return errors.Errorf("%s 配置的关联字段[guild_receive_yinliang_from_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_receive_yinliang_from_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_send_yinliang_to_member") {
		dAtA.GuildSendYinliangToMember = cOnFigS.GetMailData(pArSeR.String("guild_send_yinliang_to_member"))
	} else {
		dAtA.GuildSendYinliangToMember = cOnFigS.GetMailData("GuildSendYinliangToMember")
	}
	if dAtA.GuildSendYinliangToMember == nil {
		return errors.Errorf("%s 配置的关联字段[guild_send_yinliang_to_member] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_send_yinliang_to_member"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_task_evaluate_prize") {
		dAtA.GuildTaskEvaluatePrize = cOnFigS.GetMailData(pArSeR.String("guild_task_evaluate_prize"))
	} else {
		dAtA.GuildTaskEvaluatePrize = cOnFigS.GetMailData("GuildTaskEvaluatePrize")
	}
	if dAtA.GuildTaskEvaluatePrize == nil {
		return errors.Errorf("%s 配置的关联字段[guild_task_evaluate_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_task_evaluate_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_be_robbed_prize") {
		dAtA.HebiBeRobbedPrize = cOnFigS.GetMailData(pArSeR.String("hebi_be_robbed_prize"))
	} else {
		dAtA.HebiBeRobbedPrize = cOnFigS.GetMailData("HebiBeRobbedPrize")
	}
	if dAtA.HebiBeRobbedPrize == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_be_robbed_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_be_robbed_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_complete_prize") {
		dAtA.HebiCompletePrize = cOnFigS.GetMailData(pArSeR.String("hebi_complete_prize"))
	} else {
		dAtA.HebiCompletePrize = cOnFigS.GetMailData("HebiCompletePrize")
	}
	if dAtA.HebiCompletePrize == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_complete_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_complete_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_copy_be_robbed_prize") {
		dAtA.HebiCopyBeRobbedPrize = cOnFigS.GetMailData(pArSeR.String("hebi_copy_be_robbed_prize"))
	} else {
		dAtA.HebiCopyBeRobbedPrize = cOnFigS.GetMailData("HebiCopyBeRobbedPrize")
	}
	if dAtA.HebiCopyBeRobbedPrize == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_copy_be_robbed_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_copy_be_robbed_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_copy_complete_prize") {
		dAtA.HebiCopyCompletePrize = cOnFigS.GetMailData(pArSeR.String("hebi_copy_complete_prize"))
	} else {
		dAtA.HebiCopyCompletePrize = cOnFigS.GetMailData("HebiCopyCompletePrize")
	}
	if dAtA.HebiCopyCompletePrize == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_copy_complete_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_copy_complete_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_rob_prize") {
		dAtA.HebiRobPrize = cOnFigS.GetMailData(pArSeR.String("hebi_rob_prize"))
	} else {
		dAtA.HebiRobPrize = cOnFigS.GetMailData("HebiRobPrize")
	}
	if dAtA.HebiRobPrize == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_rob_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_rob_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("hebi_room_been_robbed") {
		dAtA.HebiRoomBeenRobbed = cOnFigS.GetMailData(pArSeR.String("hebi_room_been_robbed"))
	} else {
		dAtA.HebiRoomBeenRobbed = cOnFigS.GetMailData("HebiRoomBeenRobbed")
	}
	if dAtA.HebiRoomBeenRobbed == nil {
		return errors.Errorf("%s 配置的关联字段[hebi_room_been_robbed] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hebi_room_been_robbed"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_build_guild_member_prize") {
		dAtA.McBuildGuildMemberPrize = cOnFigS.GetMailData(pArSeR.String("mc_build_guild_member_prize"))
	} else {
		dAtA.McBuildGuildMemberPrize = cOnFigS.GetMailData("McBuildGuildMemberPrize")
	}
	if dAtA.McBuildGuildMemberPrize == nil {
		return errors.Errorf("%s 配置的关联字段[mc_build_guild_member_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_build_guild_member_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("report_ada_fail") {
		dAtA.ReportAdaFail = cOnFigS.GetMailData(pArSeR.String("report_ada_fail"))
	} else {
		dAtA.ReportAdaFail = cOnFigS.GetMailData("ReportAdaFail")
	}
	if dAtA.ReportAdaFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_ada_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_ada_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_ada_success") {
		dAtA.ReportAdaSuccess = cOnFigS.GetMailData(pArSeR.String("report_ada_success"))
	} else {
		dAtA.ReportAdaSuccess = cOnFigS.GetMailData("ReportAdaSuccess")
	}
	if dAtA.ReportAdaSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_ada_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_ada_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_add_fail") {
		dAtA.ReportAddFail = cOnFigS.GetMailData(pArSeR.String("report_add_fail"))
	} else {
		dAtA.ReportAddFail = cOnFigS.GetMailData("ReportAddFail")
	}
	if dAtA.ReportAddFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_add_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_add_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_add_success") {
		dAtA.ReportAddSuccess = cOnFigS.GetMailData(pArSeR.String("report_add_success"))
	} else {
		dAtA.ReportAddSuccess = cOnFigS.GetMailData("ReportAddSuccess")
	}
	if dAtA.ReportAddSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_add_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_add_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_asa_fail") {
		dAtA.ReportAsaFail = cOnFigS.GetMailData(pArSeR.String("report_asa_fail"))
	} else {
		dAtA.ReportAsaFail = cOnFigS.GetMailData("ReportAsaFail")
	}
	if dAtA.ReportAsaFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_asa_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_asa_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_asa_success") {
		dAtA.ReportAsaSuccess = cOnFigS.GetMailData(pArSeR.String("report_asa_success"))
	} else {
		dAtA.ReportAsaSuccess = cOnFigS.GetMailData("ReportAsaSuccess")
	}
	if dAtA.ReportAsaSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_asa_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_asa_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_asd_fail") {
		dAtA.ReportAsdFail = cOnFigS.GetMailData(pArSeR.String("report_asd_fail"))
	} else {
		dAtA.ReportAsdFail = cOnFigS.GetMailData("ReportAsdFail")
	}
	if dAtA.ReportAsdFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_asd_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_asd_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_asd_success") {
		dAtA.ReportAsdSuccess = cOnFigS.GetMailData(pArSeR.String("report_asd_success"))
	} else {
		dAtA.ReportAsdSuccess = cOnFigS.GetMailData("ReportAsdSuccess")
	}
	if dAtA.ReportAsdSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_asd_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_asd_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_ass_fail") {
		dAtA.ReportAssFail = cOnFigS.GetMailData(pArSeR.String("report_ass_fail"))
	} else {
		dAtA.ReportAssFail = cOnFigS.GetMailData("ReportAssFail")
	}
	if dAtA.ReportAssFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_ass_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_ass_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_ass_success") {
		dAtA.ReportAssSuccess = cOnFigS.GetMailData(pArSeR.String("report_ass_success"))
	} else {
		dAtA.ReportAssSuccess = cOnFigS.GetMailData("ReportAssSuccess")
	}
	if dAtA.ReportAssSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_ass_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_ass_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_baoz_repatriate_moving") {
		dAtA.ReportBaozRepatriateMoving = cOnFigS.GetMailData(pArSeR.String("report_baoz_repatriate_moving"))
	} else {
		dAtA.ReportBaozRepatriateMoving = cOnFigS.GetMailData("ReportBaozRepatriateMoving")
	}
	if dAtA.ReportBaozRepatriateMoving == nil {
		return errors.Errorf("%s 配置的关联字段[report_baoz_repatriate_moving] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_baoz_repatriate_moving"), *pArSeR)
	}

	if pArSeR.KeyExist("report_baoz_repatriate_robber") {
		dAtA.ReportBaozRepatriateRobber = cOnFigS.GetMailData(pArSeR.String("report_baoz_repatriate_robber"))
	} else {
		dAtA.ReportBaozRepatriateRobber = cOnFigS.GetMailData("ReportBaozRepatriateRobber")
	}
	if dAtA.ReportBaozRepatriateRobber == nil {
		return errors.Errorf("%s 配置的关联字段[report_baoz_repatriate_robber] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_baoz_repatriate_robber"), *pArSeR)
	}

	if pArSeR.KeyExist("report_done_attacker") {
		dAtA.ReportDoneAttacker = cOnFigS.GetMailData(pArSeR.String("report_done_attacker"))
	} else {
		dAtA.ReportDoneAttacker = cOnFigS.GetMailData("ReportDoneAttacker")
	}
	if dAtA.ReportDoneAttacker == nil {
		return errors.Errorf("%s 配置的关联字段[report_done_attacker] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_done_attacker"), *pArSeR)
	}

	if pArSeR.KeyExist("report_done_baoz") {
		dAtA.ReportDoneBaoz = cOnFigS.GetMailData(pArSeR.String("report_done_baoz"))
	} else {
		dAtA.ReportDoneBaoz = cOnFigS.GetMailData("ReportDoneBaoz")
	}
	if dAtA.ReportDoneBaoz == nil {
		return errors.Errorf("%s 配置的关联字段[report_done_baoz] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_done_baoz"), *pArSeR)
	}

	if pArSeR.KeyExist("report_done_baoz_back") {
		dAtA.ReportDoneBaozBack = cOnFigS.GetMailData(pArSeR.String("report_done_baoz_back"))
	} else {
		dAtA.ReportDoneBaozBack = cOnFigS.GetMailData("ReportDoneBaozBack")
	}
	if dAtA.ReportDoneBaozBack == nil {
		return errors.Errorf("%s 配置的关联字段[report_done_baoz_back] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_done_baoz_back"), *pArSeR)
	}

	if pArSeR.KeyExist("report_done_defenser") {
		dAtA.ReportDoneDefenser = cOnFigS.GetMailData(pArSeR.String("report_done_defenser"))
	} else {
		dAtA.ReportDoneDefenser = cOnFigS.GetMailData("ReportDoneDefenser")
	}
	if dAtA.ReportDoneDefenser == nil {
		return errors.Errorf("%s 配置的关联字段[report_done_defenser] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_done_defenser"), *pArSeR)
	}

	if pArSeR.KeyExist("report_expel_attacker_fail") {
		dAtA.ReportExpelAttackerFail = cOnFigS.GetMailData(pArSeR.String("report_expel_attacker_fail"))
	} else {
		dAtA.ReportExpelAttackerFail = cOnFigS.GetMailData("ReportExpelAttackerFail")
	}
	if dAtA.ReportExpelAttackerFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_expel_attacker_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_expel_attacker_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_expel_attacker_success") {
		dAtA.ReportExpelAttackerSuccess = cOnFigS.GetMailData(pArSeR.String("report_expel_attacker_success"))
	} else {
		dAtA.ReportExpelAttackerSuccess = cOnFigS.GetMailData("ReportExpelAttackerSuccess")
	}
	if dAtA.ReportExpelAttackerSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_expel_attacker_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_expel_attacker_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_expel_defenser_fail") {
		dAtA.ReportExpelDefenserFail = cOnFigS.GetMailData(pArSeR.String("report_expel_defenser_fail"))
	} else {
		dAtA.ReportExpelDefenserFail = cOnFigS.GetMailData("ReportExpelDefenserFail")
	}
	if dAtA.ReportExpelDefenserFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_expel_defenser_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_expel_defenser_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_expel_defenser_success") {
		dAtA.ReportExpelDefenserSuccess = cOnFigS.GetMailData(pArSeR.String("report_expel_defenser_success"))
	} else {
		dAtA.ReportExpelDefenserSuccess = cOnFigS.GetMailData("ReportExpelDefenserSuccess")
	}
	if dAtA.ReportExpelDefenserSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_expel_defenser_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_expel_defenser_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_saa_fail") {
		dAtA.ReportSaaFail = cOnFigS.GetMailData(pArSeR.String("report_saa_fail"))
	} else {
		dAtA.ReportSaaFail = cOnFigS.GetMailData("ReportSaaFail")
	}
	if dAtA.ReportSaaFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_saa_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_saa_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_saa_success") {
		dAtA.ReportSaaSuccess = cOnFigS.GetMailData(pArSeR.String("report_saa_success"))
	} else {
		dAtA.ReportSaaSuccess = cOnFigS.GetMailData("ReportSaaSuccess")
	}
	if dAtA.ReportSaaSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_saa_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_saa_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_sad_fail") {
		dAtA.ReportSadFail = cOnFigS.GetMailData(pArSeR.String("report_sad_fail"))
	} else {
		dAtA.ReportSadFail = cOnFigS.GetMailData("ReportSadFail")
	}
	if dAtA.ReportSadFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_sad_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_sad_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_sad_success") {
		dAtA.ReportSadSuccess = cOnFigS.GetMailData(pArSeR.String("report_sad_success"))
	} else {
		dAtA.ReportSadSuccess = cOnFigS.GetMailData("ReportSadSuccess")
	}
	if dAtA.ReportSadSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_sad_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_sad_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_sas_fail") {
		dAtA.ReportSasFail = cOnFigS.GetMailData(pArSeR.String("report_sas_fail"))
	} else {
		dAtA.ReportSasFail = cOnFigS.GetMailData("ReportSasFail")
	}
	if dAtA.ReportSasFail == nil {
		return errors.Errorf("%s 配置的关联字段[report_sas_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_sas_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("report_sas_success") {
		dAtA.ReportSasSuccess = cOnFigS.GetMailData(pArSeR.String("report_sas_success"))
	} else {
		dAtA.ReportSasSuccess = cOnFigS.GetMailData("ReportSasSuccess")
	}
	if dAtA.ReportSasSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[report_sas_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_sas_success"), *pArSeR)
	}

	if pArSeR.KeyExist("report_watch_attacker") {
		dAtA.ReportWatchAttacker = cOnFigS.GetMailData(pArSeR.String("report_watch_attacker"))
	} else {
		dAtA.ReportWatchAttacker = cOnFigS.GetMailData("ReportWatchAttacker")
	}
	if dAtA.ReportWatchAttacker == nil {
		return errors.Errorf("%s 配置的关联字段[report_watch_attacker] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_watch_attacker"), *pArSeR)
	}

	if pArSeR.KeyExist("report_watch_defenser") {
		dAtA.ReportWatchDefenser = cOnFigS.GetMailData(pArSeR.String("report_watch_defenser"))
	} else {
		dAtA.ReportWatchDefenser = cOnFigS.GetMailData("ReportWatchDefenser")
	}
	if dAtA.ReportWatchDefenser == nil {
		return errors.Errorf("%s 配置的关联字段[report_watch_defenser] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("report_watch_defenser"), *pArSeR)
	}

	if pArSeR.KeyExist("survey_mail") {
		dAtA.SurveyMail = cOnFigS.GetMailData(pArSeR.String("survey_mail"))
	} else {
		dAtA.SurveyMail = cOnFigS.GetMailData("SurveyMail")
	}
	if dAtA.SurveyMail == nil {
		return errors.Errorf("%s 配置的关联字段[survey_mail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("survey_mail"), *pArSeR)
	}

	if pArSeR.KeyExist("system_compensation") {
		dAtA.SystemCompensation = cOnFigS.GetMailData(pArSeR.String("system_compensation"))
	} else {
		dAtA.SystemCompensation = cOnFigS.GetMailData("SystemCompensation")
	}
	if dAtA.SystemCompensation == nil {
		return errors.Errorf("%s 配置的关联字段[system_compensation] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("system_compensation"), *pArSeR)
	}

	if pArSeR.KeyExist("xiong_nu_resist_suc") {
		dAtA.XiongNuResistSuc = cOnFigS.GetMailData(pArSeR.String("xiong_nu_resist_suc"))
	} else {
		dAtA.XiongNuResistSuc = cOnFigS.GetMailData("XiongNuResistSuc")
	}
	if dAtA.XiongNuResistSuc == nil {
		return errors.Errorf("%s 配置的关联字段[xiong_nu_resist_suc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("xiong_nu_resist_suc"), *pArSeR)
	}

	if pArSeR.KeyExist("xiong_nu_score") {
		dAtA.XiongNuScore = cOnFigS.GetMailData(pArSeR.String("xiong_nu_score"))
	} else {
		dAtA.XiongNuScore = cOnFigS.GetMailData("XiongNuScore")
	}
	if dAtA.XiongNuScore == nil {
		return errors.Errorf("%s 配置的关联字段[xiong_nu_score] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("xiong_nu_score"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetMailData(string) *MailData
	GetPrize(int) *resdata.Prize
}
