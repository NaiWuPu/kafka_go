package tailLog

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var (
	tailObj *tail.Tail
)

// 收集日志模块

/************************ 初始化 ******************/
func Init(fileName string) (err error) {
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

	// 监听本地文件
	tailObj, err = tail.TailFile(fileName, config)

	if err != nil {
		fmt.Println("tail file failed err:", err)
	}

	return
}

/************************ 读取日志 ******************/
func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}