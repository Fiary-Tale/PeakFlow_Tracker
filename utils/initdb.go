package utils

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PeakFlow struct {
	Time     string `gorm:"primaryKey"`
	Upload   float64
	Download float64
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("/etc/Traffic/peakflow.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
		WriteError("failed to connect database" + err.Error() + "\n")
	}
	// 自动迁移表结构
	if err := db.AutoMigrate(&PeakFlow{}); err != nil {
		fmt.Println("failed to migrate database")
		WriteError("failed to migrate database: " + err.Error() + "\n")
	}
}

// 将数据插入到数据库中
func InsertSampleData(currentTime string, upload, download float64) error {
	// 创建一个新的PeakFlow记录
	sampleData := PeakFlow{
		Time:     currentTime,
		Upload:   upload,
		Download: download,
	}
	result := db.Create(&sampleData)
	return result.Error
}
