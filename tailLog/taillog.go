package tailLog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log_agent/kafka"
)

// TailTask 一个日志收集的任务
type TailTask struct {
	path     string
	Topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:     path,
		Topic:    topic,
		instance: nil,
	}
	tailObj.init() // 根据路径去打开对应的日志
	return
}

// 收集日志模块
/************************ 初始化 ******************/
func (t *TailTask) init() () {
	// 打开日志文件收集日志
	config := tail.Config{
		Location:    &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件哪个位置开始读
		ReOpen:      true,                                 // 日志大小满足，重新新建一个新文件
		MustExist:   false,                                // 文件不存在不报错
		Poll:        true,                                 // 轮询文件更改而不是使用inotify
		Pipe:        false,                                // 是命名管道（mkfifo）
		RateLimiter: nil,                                  //
		Follow:      true,                                 // 跟随新文件
		MaxLineSize: 0,                                    // 如果非零，则将较长的行拆分为多行
		Logger:      nil,                                  // 当设置为nil时，Logger为tail.DefaultLogger //禁用日志记录：将字段设置为尾部丢弃记录器
	}
	var err error
	// 监听本地文件
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file ffailed, err :", err)
	}
	go t.run()	// 去采集日志发送到kafka
}

/************************ 读取日志 ******************/
func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}

func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines:
			kafka.Send2Chan(t.Topic, line.Text) // 函数调函数
		}
	}
}
