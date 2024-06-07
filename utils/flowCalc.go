package utils

import (
	"fmt"
	"log"
	"time"
)

// 从日志数据库中读取上月流量并返回上行和下行总数

func getLastFlow() (UploadDelta, DownloadDelta float64, err error) {
	var flows []PeakFlow
	// 获取上个月的时间戳
	now := time.Now()
	// 获取上个月第一天零点
	startOfMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	// 获取这个月的第一天零点
	endOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 查询上个月的流量数据
	result := db.Where("time >= ? AND time < ?", startOfMonth, endOfMonth).Find(&flows)
	if result.Error != nil {
		return 0, 0, fmt.Errorf("failed to query database: %v", result.Error)
	}
	for _, flow := range flows {
		if err != nil {
			log.Printf("failed to parse timestamp: %v", err)
			WriteError("failed to parse timestamp: " + err.Error() + "\n")
			continue
		}
		UploadDelta += flow.Upload
		DownloadDelta += flow.Download
	}
	return UploadDelta, DownloadDelta, nil
}

// 从日志中读取流量数据并返回平峰和高峰流量
func getPeakOrOff() (peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta float64, err error) {
	var flows []PeakFlow
	// 获取当前时间的前一天时间戳
	yesterday := time.Now().AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	// 查询前一天的流量数据
	result := db.Where("time >= ? AND time < ?", startOfDay, endOfDay).Find(&flows)
	if result.Error != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to query database: %v", result.Error)
	}
	for _, flow := range flows {
		timestamp, err := time.Parse("2006-01-02 15:04:05", flow.Time)
		if err != nil {
			log.Printf("failed to parse timestamp: %v", err)
			WriteError("failed to parse timestamp: " + err.Error() + "\n")
			continue
		}
		if isOffPeakTime(timestamp) {
			peakUploadDelta += flow.Upload
			peakDownloadDelta += flow.Download
		} else {
			offpeakUploadDelta += flow.Upload
			offpeakDownloadDelta += flow.Download
		}
	}
	return peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta, nil
}

// 判断当前时间是否为早晚高峰
func isOffPeakTime(timestamp time.Time) bool {
	// 读取文件的小时部分,读取的结果与设置的全局变量进行判断,高峰期为true,低峰期为false
	currentHour := timestamp.Hour() // 获取当前时间的小时部分
	return currentHour > peakHoursStart && currentHour <= peakHoursEnd
}

// 记录每日流量总和
func getPeakAndOff() (upload, download float64, err error) {
	peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta, err := getPeakOrOff()
	if err != nil {
		return 0, 0, err
	}
	upload = peakUploadDelta + offpeakUploadDelta
	download = peakDownloadDelta + offpeakDownloadDelta
	return upload, download, nil
}
