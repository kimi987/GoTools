package combatx

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
)

func HandleBytes(config *Config, uploader Uploader, data []byte) ([]byte, error) {
	request := &server_proto.CombatXRequestServerProto{}
	if err := request.Unmarshal(data); err != nil {
		return nil, errors.Wrapf(err, "Combat.Unmarshal Request Fail")
	}

	response := Handle(config, uploader, request)
	return response.Marshal()
}

//var localUploader = NewLocalUploader("temp")
//
//func LocalHandle(request *server_proto.CombatXRequestServerProto) *server_proto.CombatXResponseServerProto {
//	return Handle(localUploader, request)
//}

func Handle(config *Config, uploader Uploader, request *server_proto.CombatXRequestServerProto) *server_proto.CombatXResponseServerProto {
	response := &server_proto.CombatXResponseServerProto{}

	combat, err := NewCombat(request, config)
	if err != nil {
		response.ReturnCode = 1
		response.ReturnMsg = err.Error()
		return response
	}

	result := combat.Calculate()
	//if err != nil {
	//	response.ReturnCode = 4
	//	response.ReturnMsg = err.Error()
	//	return response
	//}

	data, err := result.Marshal()
	if err != nil {
		response.ReturnCode = 2
		response.ReturnMsg = err.Error()
		return response
	}

	link, secondLink, err := uploader.Upload(request.UploadFilePath, data)
	if err != nil {
		response.ReturnCode = 3
		response.ReturnMsg = err.Error()
		return response
	}

	if request.ReturnResult {
		response.Result = result
	}
	response.Score = result.Score
	response.TotalFrame = result.MaxFrame

	response.Link = link
	response.AttackerShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_SINGLE_X, Link: response.Link, SecondLink: secondLink, IsAttacker: true}
	response.DefenserShare = &shared_proto.CombatShareProto{Type: shared_proto.CombatType_SINGLE_X, Link: response.Link, SecondLink: secondLink, IsAttacker: false}

	response.AttackerId = request.AttackerId
	response.DefenserId = request.DefenserId
	response.AttackerWin = result.AttackerWin

	aliveSoldierMap := make(map[int32]int32)
	for _, v := range result.AliveSolider {
		aliveSoldierMap[v.Key] = v.Value
	}

	response.AttackerAliveSoldier = make(map[int32]int32)
	for i, v := range result.Attacker.Troops {
		if s := aliveSoldierMap[result.AttackerTroopData[i].Index]; s > 0 {
			response.AttackerAliveSoldier[v.Captain.Id] = s
		}
	}

	response.DefenserAliveSoldier = make(map[int32]int32)
	if s := aliveSoldierMap[0]; s > 0 {
		response.DefenserAliveSoldier[0] = s
	}
	for i, v := range result.Defenser.Troops {
		if s := aliveSoldierMap[result.DefenserTroopData[i].Index]; s > 0 {
			response.DefenserAliveSoldier[v.Captain.Id] = s
		}
	}

	killSoldierMap := make(map[int32]int32)
	for _, v := range result.KillSolider {
		killSoldierMap[v.Key] = v.Value
	}

	response.AttackerKillSoldier = make(map[int32]int32)
	for i, v := range result.Attacker.Troops {
		if s := killSoldierMap[result.AttackerTroopData[i].Index]; s > 0 {
			response.AttackerKillSoldier[v.Captain.Id] = s
		}
	}

	response.DefenserKillSoldier = make(map[int32]int32)
	if s := killSoldierMap[0]; s > 0 {
		// 城墙击杀
		response.DefenserKillSoldier[0] = s
	}
	for i, v := range result.Defenser.Troops {
		if s := killSoldierMap[result.DefenserTroopData[i].Index]; s > 0 {
			response.DefenserKillSoldier[v.Captain.Id] = s
		}
	}

	return response
}
