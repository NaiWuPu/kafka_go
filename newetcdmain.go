package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_agent/conf"
	"log_agent/etcd"
	"time"
)

var cfg1 = new(conf.AppConf)

func main() {
	// 读取配置文件
	err := ini.MapTo(cfg1, "./conf/app.ini")
	if err != nil {
		fmt.Printf("load ini filed, err %v\n", err)
	}


	// 2.初始化ETCD
	err = etcd.Init(cfg1.EtcdConf.Address, time.Duration(cfg1.EtcdConf.TimeOut)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err :%v \n", err)
		return
	}
	fmt.Println("init etcd success")

	// 写入etcdConf
	key := cfg1.EtcdConf.Key
	key = "/logagent/172.16.1.218/collect_config"
	value := `[{"path":"D:/gopath/src/log_agent/log/my.log","topic":"web_log"},{"path":"d:/xxx/redis.log","topic":"redis_log"}]`
	err = etcd.PutConf(key, value)
	if err != nil {
		fmt.Printf("etcd PutConf err :%v\n", err)
	}



}
