package conf

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	TimeOut int    `ini:"timeout"`
}

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	EtcdConf    `ini:"etcd"`
}
