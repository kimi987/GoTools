package controllers

import (
	"bytes"
	"dingTalkRobotGS/model"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var DingTalkRobotAddress = "https://oapi.dingtalk.com/robot/send?access_token=145e42d9d9154db13999692ec0ad25914fbd991f3d6e56485c5b1b4fda8f35e4"
var HelpMsg = ">通过@机器人给机器人发送消息 [构建资源|构建整包] [版署|测试|正式] [TAG] \n 构建资源:默认为测试资源,版署无资源 \n 构建整包:版署为默认版本,其他为构建资源版本 \n 版本号为当前的资源版本"
var GitURL = "http://git7.lightpaw.com/api/v4/projects/75/trigger/pipeline"

type MainController struct {
	beego.Controller
}

//Get Method
func (c *MainController) Get() {
	fmt.Println("Get")
	c.TplName = "index.tpl"
	c.Ctx.WriteString("OK")
}

//Post method
func (c *MainController) Post() {
	var rd model.ReqData
	json.Unmarshal(c.Ctx.Input.RequestBody, &rd)
	c.TplName = "index.tpl"
	c.Ctx.WriteString("OK")

	parseCommand(rd.SenderNick, rd.SessionWebhook, rd.Text.Content)
	// Post(rd.SessionWebhook, rd, "application/json")
}

func parseCommand(nickname, webhook, command string) string {
	commands := strings.Fields(command)

	for _, v := range commands {
		fmt.Println("v = ", v)
	}
	if len(commands) == 0 {
		PostDingTalkHelpMsg(webhook, nickname)
		return "消息为空, 请输入[帮助]来获取帮助信息"
	}

	cmd := commands[0]

	switch cmd {
	case "帮助":
		PostDingTalkHelpMsg(webhook, nickname)
	case "构建资源":
		if len(commands) < 2 {
			PostDingTalkDefineMsg(webhook, nickname, "参数错误, 请输入 [测试|正式] 选择构建资源分类")
			return "参数错误"
		}
		switch commands[1] {
		// case "测试":
		// 	PostUpdateMsg(nickname, webhook, "master", "TE1", "TE1", "", "")
		// case "GP" :
		// 	PostUpdateMsg(nickname, webhook, "releaseGP", "GP", "GP", "", "GP")
		// case "GPInit" :
		// 	PostUpdateMsg(nickname, webhook, "releaseGP", "GPInit", "GPInit", "", "GPInit")	
		// case "GPIOS":
		// 	PostUpdateMsg(nickname, webhook, "releaseGP", "GPIOS", "GPIOS", "", "GPIOS")
		// case "GP1":
		// 	PostUpdateMsg(nickname, webhook, "releaseGP.1", "GP1", "GP1", "", "GP1")
		// case "GPIOS1":
		// 	PostUpdateMsg(nickname, webhook, "releaseGP.1", "GPIOS1", "GPIOS1", "", "GPIOS1")
		// case "TW":
		// 	PostUpdateMsg(nickname, webhook, "releaseTW", "TW", "TW", "", "TW")
		// case "TWIOS":
		// 	PostUpdateMsg(nickname, webhook, "releaseTW", "TWIOS", "TWIOS", "", "TWIOS")
		// case "CN":
		// 	PostUpdateMsg(nickname, webhook, "releaseCN", "CN", "CN", "", "CN")
		// case "CNIOS":
		// 	PostUpdateMsg(nickname, webhook, "releaseCN", "CNIOS", "CNIOS", "", "CNIOS")
		// case "BC":
		// 	PostUpdateMsg(nickname, webhook, "releaseBaiCai", "BC", "BC", "", "BC")
		// case "ENG":
		// 	PostUpdateMsg(nickname, webhook, "releaseENG", "ENG", "ENG", "", "ENG")
		// case "ENGIOS":
		// 	PostUpdateMsg(nickname, webhook, "releaseENG", "ENGIOS", "ENGIOS", "", "ENGIOS")
		// case "正式":
		// 	PostUpdateMsg(nickname, webhook, "develop1", "CN", "CN", "", "")
		case "版署":
			PostDingTalkDefineMsg(webhook, nickname, "参数错误, [版署]无可用构建资源 请输入[测试|正式]")
			return "参数错误"
		default:
			PostDingTalkDefineMsg(webhook, nickname, "参数错误, 请输入 [测试|正式] 选择渠道")
			return "参数错误"
		}
	case "构建整包":
		if len(commands) < 2 {
			PostDingTalkDefineMsg(webhook, nickname, "参数错误, 请输入 [版署|测试|正式|GP] 选择包体")
			return "参数错误"
		}
		// if len(commands) >= 4 {
		// 	sdk = commands[3]

		// 	if sdk != "taptap" && sdk != "hykb" || sdk != "default" || sdk != "tt" {
		// 		PostDingTalkDefineMsg(webhook, nickname, "参数错误, 请输入 现在只支持 [taptap|hykb|tt|default] 渠道构建")
		// 		return "参数错误"
		// 	}
		// }

		switch commands[1] {
		case "测试Fake":
			PostUpdateMsg(nickname, webhook, "master", "devFake", "[测试][channel=dev]")
		case "测试Real":
			PostUpdateMsg(nickname, webhook, "master", "devReal", "[测试][channel=dev]")
		}

	default:
		PostDingTalkHelpMsg(webhook, nickname)
	}

	return ""
}

func PostBuildCommand() {

}

func PostUpdateMsg(nickname, dingTalkUrl, ref, triggle, eventName string) {
	postData := make(map[string]string)
	postData["token"] = "41e02c841bd2c916ef0c3b8411d016"
	postData["ref"] = ref
	postData["variables[TRIGGER]"] = triggle
	// postData["variables[TAG]"] = tag
	// postData["variables[SDK]"] = sdk
	var data, err = PostWithFormData1("POST", GitURL, &postData)
	if err != nil {
		fmt.Println("PostWithFormData error = ", err)
		PostDingTalkFailMsg(dingTalkUrl, eventName, nickname, fmt.Sprintf("%v", err))
		return
	}
	fmt.Println("data = ", string(data))
	var result model.GitBKData
	json.Unmarshal(data, &result)
	if result.Status == "pending" {
		PostDingTalkSuccMsg(dingTalkUrl, eventName, nickname, &result)
	} else {
		PostDingTalkFailMsg(dingTalkUrl, eventName, nickname, fmt.Sprintf("%v", data))
	}

}

func PostDingTalkDefineMsg(dingTalkUrl, nickname, msg string) {
	md := &model.DingTalkBKMsg{
		Msgtype: "text",
		Markdown: &model.DingTalkBKMD{
			Title: "更新帮助",
			Text:  msg,
		},
	}

	Post(dingTalkUrl, md, "application/json")
}

func PostDingTalkHelpMsg(dingTalkUrl, nickname string) {
	md := &model.DingTalkBKMsg{
		Msgtype: "text",
		Markdown: &model.DingTalkBKMD{
			Title: "更新帮助",
			Text:  HelpMsg,
		},
	}

	Post(dingTalkUrl, md, "application/json")
}

func PostDingTalkFailMsg(dingTalkUrl, eventName, nickname, reason string) {

	msg := fmt.Sprintf("[@%s] \n # **创建任务失败 %s# \n  > %s", nickname, eventName, reason)

	md := &model.DingTalkBKMsg{
		Msgtype: "markdown",
		Markdown: &model.DingTalkBKMD{
			Title: "更新事件",
			Text:  msg,
		},
	}

	Post(dingTalkUrl, md, "application/json")
}

func PostDingTalkSuccMsg(dingTalkUrl, eventName, nickname string, gitBKData *model.GitBKData) {

	msg := fmt.Sprintf("[@%s] \n # **开始构建%s# \n 当前分支%s[日志](%s)", nickname, eventName, gitBKData.Ref, gitBKData.Web_url)

	md := &model.DingTalkBKMsg{
		Msgtype: "markdown",
		Markdown: &model.DingTalkBKMD{
			Title: "更新事件",
			Text:  msg,
		},
	}

	Post(dingTalkUrl, md, "application/json")
}

//发送POST请求
//url:请求地址		data:POST请求提交的数据		contentType:请求体格式，如：application/json
//content:请求返回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, _ := json.Marshal(data)

	fmt.Println("url = " , url)
	fmt.Println("jsonStr = " , jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	fmt.Println("content = ", content)
	return
}

func PostWithFormData(method, url string, postData *map[string]string) ([]byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	fmt.Println("method = ", method)
	fmt.Println("url = ", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return data, nil
	}
	return nil, fmt.Errorf("发送构建请求失败")
}

func PostWithFormData1(method, gitUrl string, postData *map[string]string) ([]byte, error) {
	val := url.Values{}
	for k, v := range *postData {
		val.Add(k, v)
	}
	resp, err := http.PostForm(gitUrl, val)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}
	fmt.Println("resp.StatusCode = ", resp.StatusCode)
	fmt.Println("resp.body = ", string(body))
	return nil, fmt.Errorf("发送构建请求失败")
}
