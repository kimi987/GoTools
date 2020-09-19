package sharedguilddata

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewMember(id int64,
	classLevelData *guild_data.GuildClassLevelData, createTime time.Time) *GuildMember {
	return newMember(id, classLevelData, createTime, nil, nil, nil)
}

func newMember(id int64, classLevelData *guild_data.GuildClassLevelData,
	createTime time.Time, proto *server_proto.GuildMemberServerProto,
	datas *config.ConfigDatas, npcProto *shared_proto.HeroBasicSnapshotProto) *GuildMember {

	m := &GuildMember{}
	m.IdHolder = idbytes.NewIdHolder(id)
	m.classLevelData = classLevelData
	m.createTime = createTime
	m.npcProto = npcProto

	if proto != nil {
		if datas != nil {
			m.classTitleData = datas.GetGuildClassTitleData(proto.ClassTitle)
		}

		m.hufuAmount = proto.HufuAmount
		m.hufuTotalAmount = proto.HufuTotalAmount
		m.hufuAmountPerDay = proto.HufuAmountPerDay
		m.hufuAmount7 = proto.HufuAmount
		for _, a := range m.hufuAmountPerDay {
			m.hufuAmount7 += a
		}

		m.contributionAmount = proto.ContributionAmount
		m.contributionTotalAmount = proto.ContributionTotalAmount
		m.contributionAmountPerDay = proto.ContributionAmountPerDay
		m.contributionAmount7 = proto.ContributionAmount
		for _, a := range m.contributionAmountPerDay {
			m.contributionAmount7 += a
		}

		m.donationAmount = proto.DonationAmount
		m.donationTotalAmount = proto.DonationTotalAmount
		m.donationAmountPerDay = proto.DonationAmountPerDay
		m.donationAmount7 = proto.DonationAmount
		for _, a := range m.donationAmountPerDay {
			m.donationAmount7 += a
		}

		m.donationTotalYuanbao = proto.DonationTotalYuanbao

		m.isTechHelpable = proto.IsTechHelpable

		m.salary = proto.Salary

		m.historySalary = proto.HistorySalary

		m.workshopPrizeCount = proto.WorkshopPrizeCount
	}

	return m
}

type GuildMember struct {
	idbytes.IdHolder

	classLevelData *guild_data.GuildClassLevelData
	classTitleData *guild_data.GuildClassTitleData

	createTime time.Time

	// 虎符捐献
	hufuAmount       uint64
	hufuTotalAmount  uint64
	hufuAmountPerDay []uint64
	hufuAmount7      uint64

	// 联盟捐献
	contributionAmount       uint64
	contributionTotalAmount  uint64
	contributionAmountPerDay []uint64
	contributionAmount7      uint64

	donationAmount       uint64
	donationTotalAmount  uint64
	donationAmountPerDay []uint64
	donationAmount7      uint64

	donationTotalYuanbao uint64

	isTechHelpable bool

	npcProto *shared_proto.HeroBasicSnapshotProto

	salary uint64

	historySalary uint64

	// 不保存也没关系
	nextConveneTime time.Time

	// 联盟工坊奖励个数
	workshopPrizeCount uint64

	// 展示联盟工坊不存在弹窗
	isShowWorkshopNotExist bool
}

func (m *GuildMember) SetShowWorkshopNotExist(toSet bool) {
	m.isShowWorkshopNotExist = toSet
}

func (m *GuildMember) GetShowWorkshopNotExist() bool {
	return m.isShowWorkshopNotExist
}

func (m *GuildMember) IncWorkshopPrizeCount() {
	m.workshopPrizeCount++
}

func (m *GuildMember) GetWorkshopPrizeCount() uint64 {
	return m.workshopPrizeCount
}

func (m *GuildMember) ClearAndGetWorkshopPrizeCount() uint64 {
	count := m.workshopPrizeCount
	m.workshopPrizeCount = 0
	return count
}

func (m *GuildMember) SetWorkshopPrizeCount(toSet uint64) {
	m.workshopPrizeCount = toSet
}

func (m *GuildMember) Salary() uint64 {
	return m.salary
}

func (m *GuildMember) SetSalary(new uint64) {
	m.salary = new
}

func (m *GuildMember) GetNextConveneTime() time.Time {
	return m.nextConveneTime
}

func (m *GuildMember) SetNextConveneTime(toSet time.Time) {
	m.nextConveneTime = toSet
}

func (m *GuildMember) AddHistorySalary(toAdd uint64) (new uint64) {
	m.historySalary += toAdd
	return m.historySalary
}

func (m *GuildMember) IsNpc() bool {
	return npcid.IsNpcId(m.Id())
}

func (m *GuildMember) GetCreateTime() time.Time {
	return m.createTime
}

func (m *GuildMember) GetIsTechHelpable() bool {
	return m.isTechHelpable
}

func (m *GuildMember) SetIsTechHelpable(b bool) {
	m.isTechHelpable = b
}

func (m *GuildMember) NpcProto() *shared_proto.HeroBasicSnapshotProto {
	return m.npcProto
}

func (m *GuildMember) resetDaily(contributionDay int) {
	// 虎符贡献
	if len(m.hufuAmountPerDay)+1 >= contributionDay {
		m.hufuAmountPerDay = u64.RemoveHead(m.hufuAmountPerDay)
	}
	m.hufuAmountPerDay = append(m.hufuAmountPerDay, m.hufuAmount)
	m.hufuAmount = 0
	m.hufuAmount7 = 0
	for _, a := range m.hufuAmountPerDay {
		m.hufuAmount7 += a
	}

	// 贡献值
	if len(m.contributionAmountPerDay)+1 >= contributionDay {
		m.contributionAmountPerDay = u64.RemoveHead(m.contributionAmountPerDay)
	}
	m.contributionAmountPerDay = append(m.contributionAmountPerDay, m.contributionAmount)
	m.contributionAmount = 0
	m.contributionAmount7 = 0
	for _, a := range m.contributionAmountPerDay {
		m.contributionAmount7 += a
	}

	// 捐献值
	if len(m.donationAmountPerDay)+1 >= contributionDay {
		m.donationAmountPerDay = u64.RemoveHead(m.donationAmountPerDay)
	}
	m.donationAmountPerDay = append(m.donationAmountPerDay, m.donationAmount)
	m.donationAmount = 0
	m.donationAmount7 = 0
	for _, a := range m.donationAmountPerDay {
		m.donationAmount7 += a
	}
}

func (m *GuildMember) HasPermission(f func(permission *guild_data.GuildPermissionData) bool) bool {

	if m.classTitleData != nil {
		if f(m.classTitleData.Permission) {
			return true
		}
	}

	return f(m.classLevelData.Permission)
}

func (m *GuildMember) ClassLevelData() *guild_data.GuildClassLevelData {
	return m.classLevelData
}

func (m *GuildMember) SetClassLevelData(toSet *guild_data.GuildClassLevelData) {
	m.classLevelData = toSet
}

func (m *GuildMember) ClassTitleData() *guild_data.GuildClassTitleData {
	return m.classTitleData
}

func (m *GuildMember) SetClassTitleData(toSet *guild_data.GuildClassTitleData) {
	m.classTitleData = toSet
}

func (m *GuildMember) AddHufu(toAdd uint64) {
	m.hufuAmount += toAdd
	m.hufuTotalAmount += toAdd
	m.hufuAmount7 += toAdd
}

func (m *GuildMember) HufuAmount() uint64 {
	return m.hufuAmount
}

func (m *GuildMember) HufuTotalAmount() uint64 {
	return m.hufuTotalAmount
}

func (m *GuildMember) HufuAmount7() uint64 {
	return m.hufuAmount7
}

func (m *GuildMember) AddContribution(toAdd uint64) {
	m.contributionAmount += toAdd
	m.contributionTotalAmount += toAdd
	m.contributionAmount7 += toAdd
}

func (m *GuildMember) ContributionAmount() uint64 {
	return m.contributionAmount
}

func (m *GuildMember) ContributionTotalAmount() uint64 {
	return m.contributionTotalAmount
}

func (m *GuildMember) ContributionAmount7() uint64 {
	return m.contributionAmount7
}

func (m *GuildMember) AddDonation(toAdd uint64) {
	m.donationAmount += toAdd
	m.donationTotalAmount += toAdd
	m.donationAmount7 += toAdd
}

func (m *GuildMember) DonationAmount() uint64 {
	return m.donationAmount
}

func (m *GuildMember) DonationTotalAmount() uint64 {
	return m.donationTotalAmount
}

func (m *GuildMember) DonationAmount7() uint64 {
	return m.donationAmount7
}

func (m *GuildMember) DonateTotalYuanbao() uint64 {
	return m.donationTotalYuanbao
}

func (m *GuildMember) AddDonateYuanbao(toAdd uint64) {
	m.donationTotalYuanbao += toAdd
}

func (m *GuildMember) encodeClient(snapshot *shared_proto.HeroBasicSnapshotProto, isTodayJoinResistXiongNu IsTodayJoinXiongNuFunc) *shared_proto.GuildMemberProto {

	proto := &shared_proto.GuildMemberProto{}

	proto.Hero = snapshot
	proto.IsTodayJoinStart = isTodayJoinResistXiongNu(m.Id())

	proto.ClassLevel = u64.Int32(m.classLevelData.Level)
	proto.CreateTime = timeutil.Marshal32(m.createTime)

	proto.HufuAmount = u64.Int32(m.hufuAmount)
	proto.HufuTotalAmount = u64.Int32(m.hufuTotalAmount)
	proto.HufuAmount7 = u64.Int32(m.hufuAmount7)

	proto.ContributionAmount = u64.Int32(m.contributionAmount)
	proto.ContributionTotalAmount = u64.Int32(m.contributionTotalAmount)
	proto.ContributionAmount7 = u64.Int32(m.contributionAmount7)

	proto.DonationAmount = u64.Int32(m.donationAmount)
	proto.DonationTotalAmount = u64.Int32(m.donationTotalAmount)
	proto.DonationAmount7 = u64.Int32(m.donationAmount7)

	proto.DonationTotalYuanbao = u64.Int32(m.donationTotalYuanbao)

	proto.Salary = u64.Int32(m.salary)
	proto.HistorySalary = u64.Int32(m.historySalary)

	return proto
}

func (m *GuildMember) encodeServer() *server_proto.GuildMemberServerProto {
	proto := &server_proto.GuildMemberServerProto{}
	proto.Id = m.Id()

	proto.ClassLevel = m.classLevelData.Level
	if m.classTitleData != nil {
		proto.ClassTitle = m.classTitleData.Id
	}
	proto.CreateTime = timeutil.Marshal64(m.createTime)

	proto.HufuAmount = m.hufuAmount
	proto.HufuTotalAmount = m.hufuTotalAmount
	proto.HufuAmountPerDay = m.hufuAmountPerDay

	proto.ContributionAmount = m.contributionAmount
	proto.ContributionTotalAmount = m.contributionTotalAmount
	proto.ContributionAmountPerDay = m.contributionAmountPerDay

	proto.DonationAmount = m.donationAmount
	proto.DonationTotalAmount = m.donationTotalAmount
	proto.DonationAmountPerDay = m.donationAmountPerDay

	proto.DonationTotalYuanbao = m.donationTotalYuanbao

	proto.IsTechHelpable = m.isTechHelpable

	proto.Salary = m.salary
	proto.HistorySalary = m.historySalary

	proto.WorkshopPrizeCount = m.workshopPrizeCount

	return proto
}
