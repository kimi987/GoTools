package gen

import (
	"encoding/xml"
	"io/ioutil"
	"github.com/lightpaw/logrus"
	"strings"
	"bytes"
	"fmt"
)

/*
	======== xml ===========
 */

func Unmarshal(xmlPath string) (*MetaLib, error) {
	data, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		logrus.WithError(err).Panic("tlogBaseService 加载 tlog.xml 失败")
		return nil, err
	}
	v := MetaLib{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		logrus.WithError(err).Panic("tlogBaseService 解析 tlog.xml 失败")
		return nil, err
	}

	for _, s := range v.Structs {
		s.init()
	}

	return &v, nil
}

type MetaLib struct {
	TagSetVersion int            `xml:"tagsetversion,attr"`
	Name          string         `xml:"name,attr"`
	Version       int            `xml:"version,attr"`
	Structs       []*XmlStruct   `xml:"struct"`
	MacrosGroups  []*MacrosGroup `xml:"macrosgroup"`
}

type XmlStruct struct {
	Name    string  `xml:"name,attr"`
	Filter  int     `xml:"filter,attr"`
	Version int     `xml:"version,attr"`
	Desc    string  `xml:"desc,attr"`
	Entrys  []Entry `xml:"entry"`

	HasHeroField bool
	HasTxField   bool

	TlogPrefix      string
	TlogPrefixPass  string
	BuildPrefix     string
	BuildPrefixPass string

	BuildMethod string

	DefineFieldString string
	PassFieldString   string

	BuildFields string
}

var heroIdFields = []string{"vRoleID"}

var systemFields = map[string][3]string{
	"GameSvrId":        {"int", "s.config.GetServerID()", "server"},
	"GameSvrID":        {"int", "s.config.GetServerID()", "server"},
	"dtEventTime":      {"time", "s.timeService.CurrentTime()", "server"},
	"vGameAppid":       {"string", "s.config.GetGameAppID()", "server"},
	"iZoneAreaID":      {"int", "s.config.GetZoneAreaID()", "server"},
	"vGameIP":          {"string", "s.config.GetLocalAddStr()", "server"},
	"vRoleID":          {"i64", "heroInfo.Id()", "hero"},
	"vRoleName":        {"string", "heroInfo.Name()", "hero"},
	"Level":            {"u64", "heroInfo.Level()", "hero"},
	"iLevel":           {"u64", "heroInfo.Level()", "hero"},
	"KingLevel":        {"u64", "heroInfo.Level()", "hero"},
	"CityLevel":        {"u64", "heroInfo.BaseLevel()", "hero"},
	"iVipLevel":        {"u64", "heroInfo.VipLevel()", "hero"},
	"PlayerFriendsNum": {"u64", "heroInfo.FriendsCount()", "hero"},
	"TotalOnlineTime":  {"u64", "uint64(heroInfo.TotalOnlineTime().Minutes())", "hero"},
	"PlatID":           {"i32", "tencentInfo.PlatID", "tx"},
	"vopenid":          {"string", "tencentInfo.OpenID", "tx"},
	"vClientIP":        {"string", "tencentInfo.ClientIP", "tx"},
	"ClientVersion":    {"string", "tencentInfo.ClientVersion", "tx"},
	"SystemSoftware":   {"string", "tencentInfo.ClientSoftware", "tx"},
	"SystemHardware":   {"string", "tencentInfo.ClientHardware", "tx"},
	"TelecomOper":      {"string", "tencentInfo.ClientTelecom", "tx"},
	"Network":          {"string", "tencentInfo.ClientNetwork", "tx"},
	"ScreenWidth":      {"i32", "tencentInfo.ScreenWidth", "tx"},
	"ScreenHight":      {"i32", "tencentInfo.ScreenHight", "tx"},
	"Density":          {"f32", "tencentInfo.Density", "tx"},
	"LoginChannel":     {"string", "tencentInfo.LoginChannel", "tx"},
	"RegChannel":       {"string", "tencentInfo.RegChannel", "tx"},
	"CpuHardware":      {"string", "tencentInfo.CpuHardware", "tx"},
	"Memory":           {"i32", "tencentInfo.Memory", "tx"},
	"GLRender":         {"string", "tencentInfo.GLRender", "tx"},
	"GLVersion":        {"string", "tencentInfo.GLVersion", "tx"},
	"DeviceId":         {"string", "tencentInfo.DeviceId", "tx"},
}

func isHeroIdField(name string) bool {
	for _, v := range heroIdFields {
		if v == name {
			return true
		}
	}
	return false
}

func isSystemField(name string) bool {
	_, exist := systemFields[name]
	return exist
}

func (s *XmlStruct) init() {

	hasHeroId := false
	for _, e := range s.Entrys {
		if info, exist := systemFields[e.Name]; exist {
			switch info[2] {
			case "hero":
				if isHeroIdField(e.Name) {
					hasHeroId = true
				} else {
					s.HasHeroField = true
				}
			case "tx":
				s.HasTxField = true
			}
		}
	}

	if s.HasHeroField {
		if s.HasTxField {
			s.TlogPrefix = "heroInfo entity.TlogHero"
			s.TlogPrefixPass = "heroInfo"
			s.BuildPrefix = "heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto"
			s.BuildPrefixPass = "heroInfo, tencentInfo"
			s.BuildMethod = "buildLogHeroTx(heroInfo, "
		} else {
			s.TlogPrefix = "heroInfo entity.TlogHero"
			s.TlogPrefixPass = "heroInfo"
			s.BuildPrefix = "heroInfo entity.TlogHero"
			s.BuildPrefixPass = "heroInfo"
			s.BuildMethod = "buildLogHero(heroInfo, "
		}
	} else {
		if s.HasTxField || hasHeroId {
			s.TlogPrefix = "heroId int64"
			s.TlogPrefixPass = "heroId"
			if s.HasTxField {
				s.BuildPrefix = "heroId int64, tencentInfo *shared_proto.TencentInfoProto"
				s.BuildPrefixPass = "heroId, tencentInfo"
				s.BuildMethod = "buildLogHeroIdTx(heroId, "
			} else {
				s.BuildPrefix = "heroId int64"
				s.BuildPrefixPass = "heroId"
				s.BuildMethod = "buildLogHeroId(heroId, "
			}
		} else {
			// 都是空值
			s.BuildMethod = "buildLog("
		}
	}

	skipFirstComma := !s.HasHeroField && !s.HasTxField && !hasHeroId
	b := &bytes.Buffer{}
	for _, e := range s.Entrys {
		if !isSystemField(e.Name) {
			if skipFirstComma {
				skipFirstComma = false
			} else {
				b.WriteString(", ")
			}
			b.WriteString(e.Name)
			b.WriteString(" ")

			// 类型
			switch strings.ToLower(e.Type) {
			case "int":
				if e.Name == "Sequence" || e.Name == "ItemDeltaCount" {
					b.WriteString("int64")
				} else {
					b.WriteString("uint64")
				}
			case "float":
				b.WriteString("float64")
			case "int数组":
				b.WriteString("[]uint64")
			case "string":
				b.WriteString("string")
			case "datetime":
				b.WriteString("time.Time")
			default:
				b.WriteString(e.Type)
			}
		}
	}

	s.DefineFieldString = b.String()

	skipFirstComma = !s.HasHeroField && !s.HasTxField && !hasHeroId
	b = &bytes.Buffer{}
	for _, e := range s.Entrys {
		if !isSystemField(e.Name) {
			if skipFirstComma {
				skipFirstComma = false
			} else {
				b.WriteString(", ")
			}
			b.WriteString(e.Name)
		}
	}

	s.PassFieldString = b.String()

	b = &bytes.Buffer{}
	for _, e := range s.Entrys {
		b.WriteString("buf.WriteString(sep)\n")
		if info, exist := systemFields[e.Name]; exist {
			// 系统字段
			switch info[0] {
			case "string":
				b.WriteString(fmt.Sprintf("writeString(buf, %s)\n", info[1]))
			default:
				if isHeroIdField(e.Name) && !s.HasHeroField {
					b.WriteString("writeI64(buf, heroId)\n")
				} else {
					b.WriteString(fmt.Sprintf("write%s(buf, %s)\n", strings.Title(info[0]), info[1]))
				}
			}
		} else {
			// 自定义字段
			// 类型
			switch strings.ToLower(e.Type) {
			case "int":
				if e.Name == "Sequence" || e.Name == "ItemDeltaCount" {
					b.WriteString(fmt.Sprintf("writeI64(buf, %s)\n", e.Name))
				} else {
					b.WriteString(fmt.Sprintf("writeU64(buf, %s)\n", e.Name))
				}
			case "int64":
				b.WriteString(fmt.Sprintf("writeI64(buf, %s)\n", e.Name))
			case "uint64":
				b.WriteString(fmt.Sprintf("writeU64(buf, %s)\n", e.Name))
			case "float":
				b.WriteString(fmt.Sprintf("writeF64(buf, %s)\n", e.Name))
			case "int数组":
				b.WriteString(fmt.Sprintf("writeU64Array(buf, %s)\n", e.Name))
			case "string":
				b.WriteString(fmt.Sprintf("writeString(buf, %s)\n", e.Name))
			case "bool":
				b.WriteString(fmt.Sprintf("writeBool(buf, %s)\n", e.Name))
			case "datetime":
				b.WriteString(fmt.Sprintf("writeTime(buf, %s)\n", e.Name))
			default:
				b.WriteString(fmt.Sprintf("writeString(buf, %s)\n", e.Name))
			}
		}
	}
	b.WriteString("buf.WriteString(line)\n")

	s.BuildFields = b.String()

}

type Entry struct {
	Name         string `xml:"name,attr"`
	Type         string `xml:"type,attr"`
	Size         int    `xml:"size,attr"`
	Desc         string `xml:"desc,attr"`
	Index        int    `xml:"index,attr,omitempty"`
	DefaultValue string `xml:"defaultvalue,attr,omitempty"`
}

type MacrosGroup struct {
	Name  string  `xml:"name,attr"`
	Macro []Macro `xml:"macro"`
}

type Macro struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
	Desc  string `xml:"desc,attr"`
}
