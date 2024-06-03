package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// 发送网络流量消息

func SendNetworkTrafficMessage(config *Config) {
	// 获取当天峰值及平峰流量情况
	peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta, err := getPeakOrOff()
	// 获取当天上行及下行流量总量
	upload, download, err := getPeakAndOff()
	// 获取当前时间
	triggerTime := time.Now().Format("2006-01-02 15:04:05")
	// 构造消息
	message := fmt.Sprintf(
		"NAS设备流量监控\n\n- **IP地址:**\n\n  %s\n- **昨日流量消耗情况:**\n\n  下行流量消耗：%.2fG\n\n  上行流量消耗：%.2fG\n\n  晚高峰流量上行消耗：%.2fG\n\n  晚高峰流量下行消耗：%.2fG\n\n  平峰流量上行消耗：%.2fG\n\n  平峰流量下行消耗：%.2fG\n\n- **Webhook触发时间:**\n  %s",
		GetIPAddress(config.Interface),
		download,             // 转换为GB
		upload,               // 转换为GB
		peakUploadDelta,      // 转换为GB
		peakDownloadDelta,    // 转换为GB
		offpeakUploadDelta,   // 转换为GB
		offpeakDownloadDelta, // 转换为GB
		triggerTime,
	)

	// 发送消息
	err = sendDingTalkMessage(config.Token, message)
	if err != nil {
		log.Printf("Error sending DingTalk message: %v", err) // 如果发送消息出错，记录日志
		WriteError(fmt.Sprintf("%s Error sending DingTalk message: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
		return
	}
	fmt.Println("Network traffic message sent successfully.") // 成功发送消息后打印提示
}

// 发送钉钉消息
func sendDingTalkMessage(token, content string) error {
	message := DingTalkMessage{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "NAS设备流量监控",
			Text:  content,
		},
	}

	data, err := json.Marshal(message) // 将消息结构体转换为JSON
	if err != nil {
		return err
	}

	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", token) // 构建Webhook URL
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))            // 发送HTTP POST请求
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode) // 检查HTTP响应状态码
	}

	return nil
}
