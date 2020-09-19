package client_config

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/client_config"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/pbutil"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
	"github.com/lightpaw/male7/util/u64"
)

func NewClientConfigModule(serverConfig iface.IndividualServerConfig) *ClientConfigModule {
	m := &ClientConfigModule{
		isDebug: serverConfig.GetIsDebug(),
	}

	if serverConfig.GetIsDebug() {
		allTableProto := loadTables("conf_client")
		if allTableProto != nil {
			allTableMap := make(map[string]pbutil.Buffer, len(allTableProto.Tables))
			for _, table := range allTableProto.Tables {
				allTableMap[table.GetTablePath()] = client_config.NewS2cConfigMsg(must.Marshal(table))
				logrus.WithField("path", table.GetTablePath()).Infoln("能够响应客户端请求的路径")
			}

			m.allTableMap = allTableMap
			m.allTableProto = allTableProto
		}
	}

	return m
}

//gogen:iface
type ClientConfigModule struct {
	allTableProto *shared_proto.AllTable // 所有表
	allTableMap   map[string]pbutil.Buffer
	isDebug       bool // 是不是debug模式
}

func loadTables(dir string) (allTable *shared_proto.AllTable) {
	datas, err := loadTableDatas(dir)
	if err != nil {
		if !os.IsNotExist(errors.Cause(err)) {
			logrus.WithField("dir", dir).WithError(err).Error("loadTables 报错了")
		}
		return
	}

	allTable = &shared_proto.AllTable{
		Tables: make([]*shared_proto.Table, 0, len(datas)),
	}

	for path, content := range datas {
		table, err := parseTable(path, content)
		if err != nil {
			logrus.WithError(err).Panicln("解析客户端表，报错了")
		}

		allTable.Tables = append(allTable.Tables, table)
	}

	return
}

func loadTableDatas(dir string) (map[string]string, error) {

	fmt.Println(filepath.Abs(dir))
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	gos := make(map[string]string)

	fmt.Println(filepath.Abs(dir))

	decoder := mahonia.NewDecoder("gbk")

	err1 := filepath.Walk(dir, func(path0 string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".txt") {
			// 不是.txt结尾的文件不读取
			return nil
		}

		data, err1 := ioutil.ReadFile(path0)
		if err1 != nil {
			return errors.Wrapf(err1, "read config fail, %s", path0)
		}

		dp := strings.Replace(path0, "\\", "/", -1) // 更新windows支持
		dp = strings.Replace(dp, dir, "", -1)

		// 字符编码转换
		if !utf8.Valid(data) {
			gos[dp] = decoder.ConvertString(string(data))
		} else {
			gos[dp] = string(data)
		}

		return nil
	})

	return gos, err1
}

func parseTable(filePath, fileContent string) (result *shared_proto.Table, err error) {
	result = &shared_proto.Table{
		TablePath: filePath,
	}

	if len(fileContent) == 0 {
		return
	}

	fileContent = deleteHeadRN(fileContent)

	as := strings.Split(fileContent, "\r\n")
	if len(as) <= 1 {
		as = strings.Split(fileContent, "\n")
		if len(as) <= 1 {
			as = strings.Split(fileContent, "\r")
			if len(as) <= 1 {
				return nil, errors.Errorf("%s 格式不正确，请复制一个正常文件，然后使用excel来编辑保存", filePath)
			}
		}
	}

	heads := strings.Split(as[1], "\t")

	for i := 0; i < len(heads); i++ {
		heads[i] = strings.TrimSpace(heads[i])
		heads[i] = strings.Trim(heads[i], "\"")
	}

	result.Rows = make([]*shared_proto.Row, 0, len(as)-2)
	for i := 2; i < len(as); i++ {
		if len(strings.TrimSpace(as[i])) == 0 {
			// empty line
			continue
		}

		line := i + 1
		fields := strings.Split(as[i], "\t")
		if len(heads) < len(fields) {
			return nil, errors.Errorf("%s 存在head之外的行，line: %d", filePath, line)
		}

		row := &shared_proto.Row{
			Cell: make([]*shared_proto.Cell, 0, len(fields)),
		}
		for idx, field := range fields {
			if len(field) <= 0 {
				continue
			}

			head := heads[idx]
			if len(head) <= 0 {
				continue
			}

			row.Cell = append(row.Cell, &shared_proto.Cell{
				Title: head,
				Data:  field,
			})
		}

		if len(row.Cell) > 0 {
			result.Rows = append(result.Rows, row)
		}
	}

	return
}

func deleteHeadRN(origin string) string {
	// 将双引号中间的换行符删掉
	array := strings.Split(origin, "\"")
	if len(array) <= 1 {
		return origin
	}

	buf := bytes.Buffer{}

	for i := 0; i < len(array); i++ {
		s := array[i]
		if len(s) == 0 {
			continue
		}

		if i%2 == 1 {
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, "\n", "", -1)
		}

		if i > 0 {
			buf.WriteString("\"")
		}

		buf.WriteString(s)
	}

	return buf.String()

}

//gogen:iface
func (m *ClientConfigModule) ProcessConfig(proto *client_config.C2SConfigProto, hc iface.HeroController) {
	if !m.isDebug {
		// 不管
		return
	}

	result := m.allTableMap[proto.Path]
	if result == nil {
		logrus.WithField("path", proto.GetPath()).Debugln("没找到该路径下面的数据")
		hc.Send(client_config.ERR_CONFIG_FAIL_PATH_NOT_EXIST)
		return
	}

	hc.Send(result)
}

//gogen:iface
func (m *ClientConfigModule) ProcessSetClientData(proto *client_config.C2SSetClientDataProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.ClientDatas().SetBool(int(proto.Index), proto.ToSetBool)
		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *ClientConfigModule) ProcessSetClientKey(proto *client_config.C2SSetClientKeyProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.ClientDatas().SetClientKey(u64.FromInt32(proto.KeyType), u64.FromInt32(proto.KeyValue))
		result.Add(client_config.NewS2cSetClientKeyMsg(proto.KeyType, proto.KeyValue))
		result.Changed()
		result.Ok()
	})
}
