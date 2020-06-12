package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_agent/conf"
	"log_agent/etcd"
	"log_agent/kafka"
	"log_agent/tailLog"
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

	// 2.初始化ETCD
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.TimeOut)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err :%v \n", err)
		return
	}
	fmt.Println("init etcd success")

	// 2.1 从etcd 中获取日志收集项的配置信息

	// 2.2 派一个哨兵去监视日志手机的变化

	// 3.打开日志文件准备收集日志
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
