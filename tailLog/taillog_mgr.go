package tailLog

import (
	"fmt"
	"log_agent/etcd"
	"time"
)

var tskMgr *tailLogMgr

// tailTask 管理者
type tailLogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry:    logEntryConf, // 把当前的日志手机项配置保存起来
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区的通道
	}

	for _, logEnrty := range logEntryConf {
		// 初始化的时候起了多少个tailtask 都记下来
		// logEntry.Path 要手机的日志文件路径
		tailObj := NewTailTask(logEnrty.Path, logEnrty.Topic)

		mk := fmt.Sprintf("%s_%s", logEnrty.Path, logEnrty.Topic)
		tskMgr.tskMap[mk] = tailObj
	}

	go tskMgr.run()
}

// 监听自己的newConfChan 有了新的配置过来则处理
// 1. 新增配置
// 2. 删除配置
// 3. 配置变更
func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("新的配置", newConf)
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					// 原来就有不需要操作
					continue
				} else {
					// 新增
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			// 删除
			for _, c1 := range t.logEntry {	// 去新的配置中逐一进行比较
			isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

// 一个函数，向外暴露tskMgr 的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
