package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 写日志模块
type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer // 声明一个全局的kafka 生产者
	logDataChan chan *logData
)

// Init 初始化 client
func Init(address []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	// tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完整数据 leader 和 follow 都确认
	// 指定分区 轮询方式
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个 partition
	config.Producer.Return.Successes = true                   // 成功交付消息在 chanel 返回
	// 链接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	logDataChan = make(chan *logData, maxSize)
	// 开启goroutine 从通道中取数据，发往kafka
	go send2Kafka()
	return
}

// 真正往kafka发送日志的函数
func send2Kafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送一个消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 给外部暴露一个函数，改函数只把日志数据发送到一个内部的channel 中
func Send2Chan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}
