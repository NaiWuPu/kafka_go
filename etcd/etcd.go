package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

func Init(addr string, timeOut time.Duration) (err error){
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:            []string{addr},
		DialTimeout:          timeOut,
	})

	// watch 操作 用来获取未来更改的通知
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
	}
	fmt.Println("connect to etcd success")
	return err
}

func GetConf(key string) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "key")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd faild, err:%v\n", err)
		return
	}

	for _, v := range resp.Kvs {
		fmt.Printf("%s:%s\n", v.Key, v.Value)
	}
}