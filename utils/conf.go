package utils

type DingTalkMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

type TrafficData struct {
	Upload   float64
	Download float64
}
