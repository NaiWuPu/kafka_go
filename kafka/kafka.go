package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 写日志模块

var (
	client sarama.SyncProducer // 声明一个全局的kafka 生产者
)

// Init 初始化 client
func Init(address []string) (err error) {
	config := sarama.NewConfig()
	// tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完整数据 leader 和 follow 都确认
	// 指定分区 轮询方式
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个 partition

	config.Producer.Return.Successes = true // 成功交付消息在 chanel 返回

	// 链接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	return
}

func Send2Kafka(topic, data string)  {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	// 发送一个消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

}