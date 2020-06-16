package tailLog

import "log_agent/etcd"

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry []*etcd.LogEntry
	takMap   map[string]*TailTask
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntry: logEntryConf,
		takMap:   nil,
	}

	for _, logEnrty := range logEntryConf {
		NewTailTask(logEnrty.Path, logEnrty.Topic)
	}
}