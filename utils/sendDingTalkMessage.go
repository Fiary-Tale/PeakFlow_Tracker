package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	lastTrafficData map[string]TrafficData      // 用于存储上一次的流量数据
	peakHoursStart  int                    = 18 // 晚高峰开始时间，18点
	peakHoursEnd    int                    = 24 // 晚高峰结束时间，24点
)

func init() {
	lastTrafficData = make(map[string]TrafficData) // 初始化流量数据存储
}

func SendMonthNetworkMessage(config *Config) {
	// 获取上月流量
	monthUpload, monthDownload, err := getLastFlow()
	if err != nil {
		log.Printf("Error getting peak upload delta: %v", err)
		WriteError("Error getting peak upload delta")
		return
	}
	// 获取当前时间
	triggerTime := time.Now().Format("2006-01-02 15:04:05")
	// 构造消息
	message := fmt.Sprintf(
		"NAS设备月流量监控\n\n- **IP地址:**\n\n  %s\n- **上月流量消耗情况:**\n\n  下行流量总消耗：%.2fG\n\n  上行流量总消耗：%.2fG\n\n- **网络流量监控推送时间:**\n  %s",
		GetIPAddress(config.Interface),
		monthDownload,
		monthUpload,
		triggerTime,
	)
	data := DingTalkMessage{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "NAS设备流量监控",
			Text:  message,
		},
	}
	switch config.Method {
	case "dingtalk":
		sendDingTalkMessage(config.Token, data)
	case "serverchan":
		sendServerChatMessage(config.Token, data)
	default:
		log.Fatalf("未知的推送方法: %s", config.Method)
	}
}

func SendNetworkMesssage(config *Config) {
	// 获取当天峰值及平峰流量情况
	peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta, err := getPeakOrOff()
	if err != nil {
		log.Printf("Error getting peak upload delta: %v", err)
		WriteError("Error getting peak upload delta")
		return
	}
	// 获取当天上行及下行流量总量
	upload, download, err := getPeakAndOff()
	// 获取当前时间
	triggerTime := time.Now().Format("2006-01-02 15:04:05")
	// 构造消息
	message := fmt.Sprintf(
		"NAS设备日流量监控\n\n- **IP地址:**\n\n  %s\n- **昨日流量消耗情况:**\n\n  下行流量总消耗：%.2fG\n\n  上行流量总消耗：%.2fG\n\n  晚高峰流量上行消耗：%.2fG\n\n  晚高峰流量下行消耗：%.2fG\n\n  平峰流量上行消耗：%.2fG\n\n  平峰流量下行消耗：%.2fG\n\n- **网络流量监控推送时间:**\n  %s",
		GetIPAddress(config.Interface),
		download,             // 转换为GB
		upload,               // 转换为GB
		peakUploadDelta,      // 转换为GB
		peakDownloadDelta,    // 转换为GB
		offpeakUploadDelta,   // 转换为GB
		offpeakDownloadDelta, // 转换为GB
		triggerTime,
	)
	data := DingTalkMessage{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "NAS设备流量监控",
			Text:  message,
		},
	}
	switch config.Method {
	case "dingtalk":
		sendDingTalkMessage(config.Token, data)
	case "serverchan":
		sendServerChatMessage(config.Token, data)
	default:
		log.Fatalf("未知的推送方法: %s", config.Method)
	}
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("2006-01-02")
	Write(fmt.Sprintf(
		"%s - - [%s] 下行流量总数: %.2f, 上行流量总数: %.2f",
		GetIPAddress(config.Interface),
		yesterdayStr,
		download,
		upload,
	))
}

// 发送钉钉消息
func sendDingTalkMessage(token string, message DingTalkMessage) {
	data, err := json.Marshal(message) // 将消息结构体转换为JSON
	if err != nil {
		log.Printf("Error sending DingTalk message: %v", err)
		WriteError(fmt.Sprintf("Error sending DingTalk message: %v", err))
		return
	}
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token) // 构建Webhook URL
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))            // 发送HTTP POST请求
	if err != nil {
		log.Printf("Error sending DingTalk message: %v", err)
		WriteError(fmt.Sprintf("Error sending DingTalk message: %v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error sending DingTalk message: %v", resp.Status)              // 打印错误状态码
		WriteError(fmt.Sprintf("Error sending DingTalk message: %v", resp.Status)) // 写入错误日志
	}
}

func sendServerChatMessage(token string, message DingTalkMessage) {
	data := url.Values{}
	data.Add("title", message.Markdown.Title)
	data.Add("text", message.Markdown.Text)
	webhookURL := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", token) // 构建Webhook URL
	resp, err := http.PostForm(webhookURL, data)                        // 发送HTTP POST请求
	if err != nil {
		log.Printf("Error sending ServerChat message: %v", err)
		WriteError(fmt.Sprintf("Error sending ServerChat message: %v", resp.Status))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error sending ServerChat message: %v", resp.Status)
		WriteError(fmt.Sprintf("Error sending ServerChat message: %v", resp.Status))
	}
}
