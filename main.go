package main

import (
	"examples/log_agent/conf"
	"examples/log_agent/kafka"
	"examples/log_agent/tailLog"
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

var cfg = new(conf.AppConf)

func main() {
	// 读取配置文件
	err := ini.MapTo(cfg, "./conf/app.ini")
	if err != nil {
		fmt.Printf("load ini filed, err %v\n", err)
	}
	// 1.初始化kafka 链接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Printf("init Kafka failed, err:%v\n", err)
		return
	}
	fmt.Printf("init Kafka \n")
	// 2.打开日志文件准备收集日志
	err = tailLog.Init(cfg.FileName)
	if err != nil {
		fmt.Printf("init tailLog failed, err:%v\n", err)
		return
	}
	fmt.Printf("init tailLog Success\n")

	run()

}

func run() {
	/************************ 读取日志 ******************/
	for {
		select {
		/************************ 发送到kafka ******************/
		case line := <-tailLog.ReadChan():
			kafka.Send2Kafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(1 * time.Second)
		}

	}
}
