package model

// {"conversationId":"cidxOW0hkFkDWO3bgD95/GooOhl50xvdPJ1m4wxM+94QO8=",
//"chatbotCorpId":"ding0ea97de30886f4f535c2f4657eb6378f",
//"chatbotUserId":"$:LWCP_v1:$txU3IoSuj5tsO6XbYIoY+OUVHsmaBHUt",
//"msgId":"msgthh69D/QQjyPZRtbAwGIMg==",
//"senderNick":"Kimi(汪津津)",
//"isAdmin":true,"senderStaffId":"11470135002042552237",
//"sessionWebhookExpiredTime":1597407747402,
//"createAt":1597402347315,
//"senderCorpId":"ding0ea97de30886f4f535c2f4657eb6378f",
//"conversationType":"1",
//"senderId":"$:LWCP_v1:$QljqiS1t/7c/Dh+q+xL8rw==",
//"sessionWebhook":"https://oapi.dingtalk.com/robot/sendBySession?session=3f479c3f3b785012fcbe1d143db4ff9c","text":{"content":"123"},"msgtype":"text"}

type ReqData struct {
	ConversationId            string          `json:"conversationId"`
	ChatbotCorpId             string          `json:"chatbotCorpId"`
	ChatbotUserId             string          `json:"chatbotUserId"`
	MsgId                     string          `json:"msgId"`
	SenderNick                string          `json:"senderNick"`
	IsAdmin                   string          `json:"isAdmin"`
	SenderStaffId             string          `json:"senderStaffId"`
	SessionWebhookExpiredTime int64           `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64           `json:"createAt"`
	SenderCorpId              string          `json:"senderCorpId"`
	ConversationType          string          `json:"conversationType"`
	SenderId                  string          `json:"senderId"`
	SessionWebhook            string          `json:"sessionWebhook"`
	Text                      *ReqDataContent `json:"text"`
	Msgtype                   string          `json:"msgtype"`
}

type ReqDataContent struct {
	Content string `json:"content"`
}

// type ResData struct {
// 	msgtype
// }

type ReqGitData struct {
	Token     string `json:"token"`
	Ref       string `json:"ref"`
	Variables string `json:""`
}

type GitBKData struct {
	ID              string           `json:"id"`
	Sha             string           `json:"sha"`
	Ref             string           `json:"ref"`
	Status          string           `json:"status"`
	Created_at      string           `json:"Created_at"`
	Updated_at      string           `json:"updated_at"`
	Web_url         string           `json:"web_url"`
	Before_sha      string           `json:"before_sha"`
	Tag             bool             `json:"tag"`
	Yaml_errors     string           `json:"yaml_errors"`
	User            *GitBKUser       `json:"user"`
	Started_at      string           `json:"started_at"`
	Finished_at     string           `json:"finished_at"`
	Committed_at    string           `json:"committed_at"`
	Duration        string           `json:"duration"`
	Detailed_status *GitDetailStatus `json:"detailed_status"`
}

type GitBKUser struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"uIsername"`
	State      string `json:"state"`
	Avatar_url string `json:"avatar_url"`
	Web_url    string `json:"web_url"`
}

type GitDetailStatus struct {
	Icon         string `json:"icon"`
	Text         string `json:"text"`
	Label        string `json:"label"`
	Group        string `json:"group"`
	Tooltip      string `json:"tooltip"`
	Has_details  bool   `json:"has_details"`
	Details_path string `json:"details_path"`
	Illustration string `json:"illustration"`
	Favicon      string `json:"favicon"`
}

type DingTalkBKMsg struct {
	Msgtype  string        `json:"msgtype"`
	Markdown *DingTalkBKMD `json:"markdown"`
}

type DingTalkBKMD struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
