package conf

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type TaillogConf struct {
	FileName string `ini:"path"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	TimeOut int    `ini:"timeout"`
}

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
	EtcdConf    `ini:"etcd"`
}
