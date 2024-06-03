package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	fileName := getYesterdayLogFileName()
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("error opening log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		timestamp, upload, download, err := parseLogLine(line)
		if err != nil {
			log.Printf("Error parsing log line: %v", err)
			continue
		}
		if isOffPeakTime(timestamp) {
			peakUploadDelta += upload
			peakDownloadDelta += download
		} else {
			offpeakUploadDelta += upload
			offpeakDownloadDelta += download
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("error reading log file: %v", err)
	}

	return peakUploadDelta, peakDownloadDelta, offpeakUploadDelta, offpeakDownloadDelta, nil
}

// 解析日志行中的时间和流量数据
func parseLogLine(line string) (timestamp time.Time, upload, download float64, err error) {
	// 示例日志行: "2006-01-02 15:04:05 Upload Delta: 12, Download Delta: 13"
	parts := strings.Split(line, " ")
	if len(parts) < 7 {
		return time.Time{}, 0, 0, fmt.Errorf("invalid log line format")
	}

	// 解析时间戳
	timeStr := fmt.Sprintf("%s %s", parts[0], parts[1])
	timestamp, err = time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Time{}, 0, 0, fmt.Errorf("invalid timestamp format")
	}

	// 解析流量数据
	for i := 2; i < len(parts); i++ {
		if strings.Contains(parts[i], "Upload") {
			fmt.Sscanf(parts[i+2], "%f", &upload)
		}
		if strings.Contains(parts[i], "Download") {
			fmt.Sscanf(parts[i+2], "%f", &download)
		}
	}
	return timestamp, upload, download, nil
}

// 判断当前时间是否为早晚高峰
func isOffPeakTime(timestamp time.Time) bool {
	// 读取文件的小时部分,读取的结果与设置的全局变量进行判断,高峰期为true,低峰期为false
	currentHour := timestamp.Hour() // 获取当前时间的小时部分
	return currentHour > peakHoursStart && currentHour < peakHoursEnd
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
