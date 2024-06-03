package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// 获取网络流量增量
func getNetworkTrafficDelta(interfaceName string) (upload float64, download float64, err error) {
	currentTrafficData := TrafficData{}
	txBytesPath := fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", interfaceName) // 上行流量文件路径
	rxBytesPath := fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", interfaceName) // 下行流量文件路径

	txData, err := ioutil.ReadFile(txBytesPath) // 读取上行流量数据
	if err != nil {
		return 0, 0, err
	}
	currentTrafficData.Upload, err = parseTrafficData(txData) // 解析上行流量数据
	if err != nil {
		return 0, 0, err
	}

	rxData, err := ioutil.ReadFile(rxBytesPath) // 读取下行流量数据
	if err != nil {
		return 0, 0, err
	}
	currentTrafficData.Download, err = parseTrafficData(rxData) // 解析下行流量数据
	if err != nil {
		return 0, 0, err
	}

	lastData, ok := lastTrafficData[interfaceName] // 获取上一次的流量数据
	if !ok {
		lastTrafficData[interfaceName] = currentTrafficData // 如果没有上一次的数据，存储当前流量数据
		return 0, 0, nil
	}

	uploadDelta := currentTrafficData.Upload - lastData.Upload       // 计算上行流量增量
	downloadDelta := currentTrafficData.Download - lastData.Download // 计算下行流量增量

	lastTrafficData[interfaceName] = currentTrafficData // 更新存储的流量数据
	uploadDelta = uploadDelta / 1e9
	downloadDelta = downloadDelta / 1e9
	return uploadDelta, downloadDelta, nil
}

// 解析流量数据
func parseTrafficData(data []byte) (float64, error) {
	strData := strings.TrimSpace(string(data)) // 去掉字符串两端的空格
	return strconv.ParseFloat(strData, 64)     // 将字符串转换为浮点数
}

// 每天每小时都进行流量增量的日志写入
func wdtFlow(config *Config) { // 流量监控狗，每小时计算一下流量增量，并加入到高峰或平峰
	now := time.Now()
	if now.Minute() == 0 && now.Second() == 0 {
		upload, download, err := getNetworkTrafficDelta(config.Interface)
		if err != nil {
			log.Printf("Error getting network traffic: %v", err) // 如果获取网络流量数据出错，记录日志
			WriteError(fmt.Sprintf("%s Error getting network traffic: %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
			return
		}
		// 获取当前日期作为日志文件名
		logFileName := fmt.Sprintf("/var/log/Traffic/%s.log", now.Format("2006-01-02"))
		result := fmt.Sprintf("%s Upload Delta: %.2f, Download Delta: %.2f\n", time.Now().Format("2006-01-02 15:04:05"), upload, download)
		Write(result, logFileName)
		// 睡眠1秒，避免重复触发
		time.Sleep(1 * time.Second)
	}
}
