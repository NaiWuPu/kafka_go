package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`  // 日志存放的路径
	Topic string `json:"topic"` // 日志要法网kafka中的topic
}

func Init(addr string, timeOut time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeOut,
	})

	// watch 操作 用来获取未来更改的通知
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
	}
	fmt.Println("connect to etcd success")
	return err
}

func GetConf(key string) (logEntrys []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd faild, err:%v\n", err)
		return
	}
	for _, v := range resp.Kvs {
		err := json.Unmarshal(v.Value, &logEntrys)
		if err != nil {
			fmt.Printf("unmarshal etcd value failed, err:%v\n", err)
		}
	}
	return logEntrys, err
}

func PutConf(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Put(ctx, key, value)
	cancel()

	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
	}
	return err
}

func WatchConf(key string, newConfch chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)

	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v key:%v value:%v \n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			// 通知别人	taillog.tskMgr
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				// 如果是删除操作，传一个空的配置项
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("unmarshal failed, err:%v \n", err)
					continue
				}
				fmt.Printf("get new conf %v\n", newConf)
				newConfch <- newConf
			}
		}
	}
}
