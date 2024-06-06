package commond

import (
	"PeakFlow_Tracker/utils"
	"flag"
	"fmt"
)

func version() {
	version := "V1.0.4"
	fmt.Println(version)
}

func Flag() {
	file := flag.String("f", "", "-f 指定加载配置文件,如 -f config.yaml,必须是yaml格式")
	ver := flag.Bool("V", false, "-V 查看当前软件版本")
	flag.Parse()
	if *file != "" {
		utils.Exec(*file)
	} else if *ver != false {
		version()
	} else {
		fmt.Println("未指定文件将尝试使用默认配置")
		utils.Exec("/etc/Traffic/config.yaml")
	}
}
