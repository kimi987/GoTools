package combat

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/pkg/errors"
)

type Uploader interface {
	Upload(filename string, data []byte) (path string, err error)
}

type UploaderFunc func(filename string, data []byte) (path string, err error)

func (f UploaderFunc) Upload(filename string, data []byte) (path string, err error) {
	return f(filename, data)
}

func NewLocalUploader(dir string) Uploader {
	return UploaderFunc(func(filename string, data []byte) (path string, err error) {
		// marshal数据，保存在本地
		// 没找到好的uuid库，先临时使用时间戳作为文件名
		//filename := fmt.Sprintf("%v.txt", time.Now().UnixNano())
		err = ioutil.WriteFile(filepath.Join(dir, filename), data, os.ModePerm)
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(path, "/") {
			return "{{local}}" + filename, nil
		} else {
			return "{{local}}/" + filename, nil
		}
	})
}

var localUploader = NewLocalUploader("temp")

func HandleBytes(uploader Uploader, data []byte) ([]byte, error) {
	request := &server_proto.CombatRequestServerProto{}
	if err := request.Unmarshal(data); err != nil {
		return nil, errors.Wrapf(err, "Combat.Unmarshal Request Fail")
	}

	response := Handle(uploader, request)
	return response.Marshal()
}

func LocalHandle(request *server_proto.CombatRequestServerProto) *server_proto.CombatResponseServerProto {
	return Handle(localUploader, request)
}

func Handle(uploader Uploader, request *server_proto.CombatRequestServerProto) *server_proto.CombatResponseServerProto {
	response := &server_proto.CombatResponseServerProto{}

	combat, err := NewCombat(request)
	if err != nil {
		response.ReturnCode = 1
		response.ReturnMsg = err.Error()
		return response
	}

	result, err := combat.Calculate()
	if err != nil {
		response.ReturnCode = 4
		response.ReturnMsg = err.Error()
		return response
	}

	data, err := result.Marshal()
	if err != nil {
		response.ReturnCode = 2
		response.ReturnMsg = err.Error()
		return response
	}

	link, err := uploader.Upload(request.UploadFilePath, data)
	if err != nil {
		response.ReturnCode = 3
		response.ReturnMsg = err.Error()
		return response
	}

	if request.ReturnResult {
		response.Result = result
	}
	response.Score = result.Score

	response.Link = link
	response.AttackerShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_SINGLE, Link: response.Link, IsAttacker: true}
	response.DefenserShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_SINGLE, Link: response.Link, IsAttacker: false}

	response.AttackerId = request.AttackerId
	response.DefenserId = request.DefenserId
	response.AttackerWin = result.AttackerWin

	aliveSoldierMap := make(map[int32]int32)
	for _, v := range result.AliveSolider {
		aliveSoldierMap[v.Key] = v.Value
	}

	response.AttackerAliveSoldier = make(map[int32]int32)
	for i, v := range result.Attacker.Troops {
		if s := aliveSoldierMap[result.AttackerTroopPos[i].Index]; s > 0 {
			response.AttackerAliveSoldier[v.Captain.Id] = s
		}
	}

	response.DefenserAliveSoldier = make(map[int32]int32)
	if s := aliveSoldierMap[0]; s > 0 {
		response.DefenserAliveSoldier[0] = s
	}
	for i, v := range result.Defenser.Troops {
		if s := aliveSoldierMap[result.DefenserTroopPos[i].Index]; s > 0 {
			response.DefenserAliveSoldier[v.Captain.Id] = s
		}
	}

	killSoldierMap := make(map[int32]int32)
	for _, v := range result.KillSolider {
		killSoldierMap[v.Key] = v.Value
	}

	response.AttackerKillSoldier = make(map[int32]int32)
	for i, v := range result.Attacker.Troops {
		if s := killSoldierMap[result.AttackerTroopPos[i].Index]; s > 0 {
			response.AttackerKillSoldier[v.Captain.Id] = s
		}
	}

	response.DefenserKillSoldier = make(map[int32]int32)
	if s := killSoldierMap[0]; s > 0 {
		response.DefenserKillSoldier[0] = s
	}
	for i, v := range result.Defenser.Troops {
		if s := killSoldierMap[result.DefenserTroopPos[i].Index]; s > 0 {
			response.DefenserKillSoldier[v.Captain.Id] = s
		}
	}

	return response
}

func HandleMultiBytes(uploader Uploader, data []byte) ([]byte, error) {
	request := &server_proto.MultiCombatRequestServerProto{}
	if err := request.Unmarshal(data); err != nil {
		return nil, errors.Wrapf(err, "Combat.Unmarshal Request Fail")
	}

	response := HandleMulti(uploader, request)
	return response.Marshal()
}

func LocalHandleMulti(request *server_proto.MultiCombatRequestServerProto) *server_proto.MultiCombatResponseServerProto {
	return HandleMulti(localUploader, request)
}

func HandleMulti(uploader Uploader, request *server_proto.MultiCombatRequestServerProto) *server_proto.MultiCombatResponseServerProto {

	response := &server_proto.MultiCombatResponseServerProto{}

	combat, err := NewMultiCombat(request)
	if err != nil {
		response.ReturnCode = 1
		response.ReturnMsg = err.Error()
		return response
	}

	result, err := combat.Calculate()
	if err != nil {
		response.ReturnCode = 4
		response.ReturnMsg = err.Error()
		return response
	}

	data, err := result.Marshal()
	if err != nil {
		response.ReturnCode = 2
		response.ReturnMsg = err.Error()
		return response
	}

	link, err := uploader.Upload(request.UploadFilePath, data)
	if err != nil {
		response.ReturnCode = 3
		response.ReturnMsg = err.Error()
		return response
	}

	if request.ReturnResult {
		response.Result = result
	}
	response.Score = result.Score

	response.Link = link
	response.AttackerShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_MULTI, Link: response.Link, IsAttacker: true}
	response.DefenserShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_MULTI, Link: response.Link, IsAttacker: false}

	response.AttackerId = request.AttackerId
	response.DefenserId = request.DefenserId
	response.AttackerWin = result.AttackerWin
	response.AliveSoldiers = combat.getAliveSoldiers()
	response.WinTimesMap = combat.getWinTimesMap()

	return response
}
