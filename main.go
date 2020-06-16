package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_agent/conf"
	"log_agent/etcd"
	"log_agent/kafka"
	"log_agent/tailLog"
	"sync"
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
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.ChanMaxSize)
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

	// 写入etcdConf
	//err = etcd.PutConf(cfg.EtcdConf.Key, `[{"path":"c:/tmp/nginx.log","topic":"web_log"},{"path":"d:/xxx/redis.log","topic":"redis_log"}]`)
	//if err != nil {
	//	fmt.Printf("etcd PutConf err :%v\n", err)
	//}

	// 2.1 从etcd 中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success :%v \n", logEntryConf)


	// 3.打开日志文件准备收集日志
	// 3.1 循环收集每一个收集项 创建tailObj
	tailLog.Init(logEntryConf)	// 因为NewConfChan 访问了 tskMgr 的newConfChan, 这个channel是在此步操作执行初始化

	// 哨兵监视 etcd 变化
	newConfChan := tailLog.NewConfChan()	// 从taillog 包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key, newConfChan) 	// 哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()

}
