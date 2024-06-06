package utils

import (
	"fmt"
	"log"
	"time"
)

// 计算晚高峰和平峰流量新版(通过读取文件)

// 获取前一天的日志文件名
func getYesterdayLogFileName() string {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return fmt.Sprintf("/var/log/Traffic/%s.log", yesterday)
}

// 从日志中读取流量数据并返回平峰和高峰流量
func getPeakOrOff() (peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta float64, err error) {
	var flows []PeakFlow
	result := db.Find(&flows)
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
