# PeakFlow_Tracker

## 项目描述
为了更好地管理和监控网络附加存储（NAS）或其他设备的流量消耗，我们开发了一套流量监控系统。该系统能够实时监控NAS设备及其他常见设备的网络流量，统计晚高峰和平峰时段的流量消耗，并通过钉钉机器人发送流量统计信息。系统还会在每天午夜重置流量统计，以便生成每日流量报告，帮助管理员及时了解网络使用情况并作出相应的管理决策。

## 项目功能
1. 实时流量监控：
+ 监控设备的上行和下行流量。
+ 计算流量增量，即两次监控之间的流量变化。
2. 高峰和平峰流量统计：
+ 定义晚高峰时段（18:00 - 24:00）。
+ 统计并累积晚高峰和平峰时段的上行和下行流量。
3. 流量数据发送：
+ 定时发送流量统计信息至钉钉。
+ 程序首次运行时立即发送一次流量统计信息，并打印日志。
4. 流量重置：
+ 每天午夜（00:00:00）重置高峰和平峰流量统计。

## 更新记录

+ 2024-5-20 PeakFlow_TrackerV1.0.1
    + 程序发布
+ 2024-5-22 PeakFlow_TrackerV1.0.2
    + 新增高峰、平峰流量统计
    + 修复第一次运行统计流量为0G的BUG(增加程序运行第一次即发送消息)
    + 新增高峰、平峰流量统计定时重置功能(限制在每日午夜00:00:00重置流量统计)
    + 新增默认寻找配置文件(无需指定配置文件,也可进行指定,默认配置文件位置为`/etc/Traffic/config.yaml`)
    + 新增日志记录功能(默认日志存储位置`/var/log/Traffic.log`)
+ 2024-6-2 PeakFlow_TrackerV1.0.3
    + 重构流量增量统计代码(每小时计算一次增量)
    + 修复程序启动只能获取平峰或高峰流量增量异常BUG
    + 修改日志记录(记录正常流量增量日志及error日志)
    + 修改日志记录文件名称(access日志设置为当前日期.log,error日志设置为error.log)
    + 修改日志记录目录(修改为/var/log/Traffic/目录)
    + 新增当前版本查看详情(输入 -V 即可输出当前版本)
+ 2024-6-6 PeakFlow_TrackerV1.0.4
    + 重构流量增量统计存储代码(现存入到/etc/Traffic/peakflow.db中)
    + 重构流量增量统计计算代码(现从/etc/Traffic/peakflow.db获取存储数据计算)
    + 删除流量增量统计日志存储及计算方法
    + 增加存储数据库自动生成(无需手动创建)
+ 2024-6-7 PeakFlow_TrackerV1.0.5
    + 修复读取流量统计结果统计所有的BUG(现统计只有前一日推送)
    + 修复添加参数执行卡死的BUG
    + 添加Server酱推送功能(可选钉钉推送或Server酱推送)
+ 2024-10-26 PeakFlow_TrackerV1.0.6
    + 新增月流量统计(每月1日在每天设定的时间节点会统计上行及下行总流量)
    + 新增流量统计日志记录(将每日流量统计写入日志/var/log/Traffic/Flow.log，以方便后续复查)
## 特性

1. 程序简单,傻瓜式操作,可快速上手使用
2. 进行钉钉消息推送,完美契合免费原则(钉钉可免费使用WebHook无次数限制,只有频率限制)
3. 速度快,体积小,跨平台使用,支持使用systemctl开机自启动或/etc/init.d/开机自启动

## 使用教程
1. 下载对应架构的程序,或使用源码自行编译(前提配置好编译环境,将二进制程序及配置文件放入到/etc/Traffic/目录下)
2. 编写配置文件,在config.yaml中添加自己的token及网卡名称,time时间可不更改,默认为24小时及30天各统计推送一次
3. 添加运行权限,`chmod 777 PeakFlow_Tracker`
4. 运行程序放入后台`./PeakFlow_Tracker &`(或是放入到/etc/init.d/目录下,使用命令/etc/init.d/PeakFlow_Tracker start)

### 开机自启教程
1. systemctl 设置开机自启动

将程序及配置文件放入到`/etc/Traffic/`下,若没有该目录则创建该目录
```bash
[Unit]
Description=PeakFlow_Tracker server
After=network.target
Wants = network.target

[Service]
Type=simple

ExecStart=/etc/Traffic/PeakFlow_Tracker -f /etc/Traffic/config.yaml
ExecReload=/bin/kill -s HUP $MAINPID
ExecStop=/bin/kill -s QUIT $MAINPID

[Install]
WantedBy=multi-user.target
```
将上述写入到`/etc/systemd/system/PeakFlow_Tracker.service`中

命令设置开机自启
```bash
systemctl start PeakFlow_Tracker.service      # 启动程序
systemctl stop PeakFlow_Tracker.service       # 停止程序
systemctl status PeakFlow_Tracker.service     # 查看程序运行状态
systemctl enable PeakFlow_Tracker.service     # 添加开机自启动
```
2. `/etc/init.d/`设置开机自启

适用于没有systemctl及service命令的情况
将程序及配置文件添加权限放入到`/etc/Traffic/`目录下,在`/etc/init.d/`目录下添加如下脚本,并添加执行权限:
```bash
#!/bin/sh
#### BEGIN INIT INFO
# Provides:          PeakFlow_Tracker
# Required-Start:    $all
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:
# Short-Description: Start PeakFlow_Tracker at boot time
### END INIT INFO

# 设置二进制程序的路径
PROG="/etc/Traffic/PeakFlow_Tracker"
PIDFILE="/var/run/PeakFlow_Tracker.pid"

case "$1" in
  start)
    echo "Starting PeakFlow_Tracker"
    $PROG &
    echo $! > $PIDFILE
    ;;
  stop)
    echo "Stopping PeakFlow_Tracker"
    kill $(cat $PIDFILE)
    rm -f $PIDFILE
    ;;
  restart)
    $0 stop
    $0 start
    ;;
  status)
    if [ -e $PIDFILE ]; then
      echo "PeakFlow_Tracker is running, PID=$(cat $PIDFILE)"
    else
      echo "PeakFlow_Tracker is not running"
    fi
    ;;
  *)
    echo "Usage: /etc/init.d/PeakFlow_Tracker {start|stop|restart|status}"
    exit 1
    ;;
esac

exit 0
```
使用方法：
```bash
/etc/init.d/ start          # 启动程序
/etc/init.d/ stop           # 停止程序
/etc/init.d/ status         # 查看程序运行状态
/etc/init.d/ enable         # 添加开机自启动
```

3. 绿联开机自启动设置方法

```
#!/bin/sh /etc/rc.common
# Copyright (C) 2015 OpenWrt.org

START=99

USE_PROCD=1

PROGRAM="/etc/Traffic/PeakFlow_Tracker"
PIDFILE="/var/run/PeakFlow_Tracker.pid"

start_service() {
    # 检查程序是否存在
    if [ ! -f "$PROGRAM" ]; then
        logger -t "PeakFlow_Tracker" -p "daemon.err" "Program not found: $PROGRAM"
        exit 1
    fi

    procd_open_instance
    procd_set_param command "$PROGRAM"
    procd_set_param stdout 1
    procd_set_param stderr 1
    procd_set_param pidfile "$PIDFILE"
    procd_set_param respawn
    procd_close_instance
}

stop_service() {
    if [ -f "$PIDFILE" ]; then
        kill $(cat $PIDFILE)
        rm -f $PIDFILE
    fi
}

restart_service() {
    stop_service
    start_service
}

reload_service() {
    restart_service
}

status_service() {
    if [ -f "$PIDFILE" ]; then
        echo "PeakFlow_Tracker is running, PID=$(cat $PIDFILE)"
    else
        echo "PeakFlow_Tracker is not running"
    fi
}
```

