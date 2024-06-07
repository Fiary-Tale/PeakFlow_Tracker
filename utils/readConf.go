package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
)

// Config 结构体，用于存储从配置文件读取的配置信息
type Config struct {
	Token     string `yaml:"token"`     // DingTalk 机器人access_token
	Interface string `yaml:"interface"` // 网络接口名称
	Time      string `yaml:"time"`      // 固定时间,示例：09:35:00
	Method    string `yaml:"method"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var Conf Config
	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		return nil, err
	}
	return &Conf, nil
}

var isFirstRun = true // 初始标记为 true

// 执行死循环发送程序

func Exec(name string) {
	conf, err := ReadConfig(name)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		WriteError(fmt.Sprintf("%s Error reading config file: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
	}

	if isFirstRun {
		_, _, err := getNetworkTrafficDelta(conf.Interface)
		if err != nil {
			log.Printf("获取网卡流量失败: %v", err)
			WriteError(fmt.Sprintf("%s 获取网卡流量失败: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
		}
		isFirstRun = false
	}

	triggerTime, err := time.Parse("15:04:05", conf.Time)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		WriteError(fmt.Sprintf("%s Error parsing time: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
	}

	// 定时任务，监控流量并发送消息
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for { // 死循环
		select {
		case <-ticker.C:
			now := time.Now()
			// 检查设定时间
			if now.Hour() == triggerTime.Hour() && now.Minute() == triggerTime.Minute() && now.Second() == triggerTime.Second() {
				SendNetworkMesssage(conf)
				time.Sleep(1 * time.Second) // 睡眠1秒，避免重复发送
			}
		default:
			wdtFlow(conf)
		}
	}
}
