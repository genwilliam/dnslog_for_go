package dnslog

import (
	"sync"
	"time"
)

// Record 表示一条 DNS 请求日志
type Record struct {
	Domain    string `json:"domain"`
	ClientIP  string `json:"client_ip"`
	Protocol  string `json:"protocol"`
	QType     string `json:"qtype"`
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Token     string `json:"token"`
}

type logStore struct {
	mu      sync.RWMutex
	maxSize int
	data    []Record
}

var store *logStore

// InitStore 初始化存储（比如最多保存 1000 条记录）
func InitStore(maxSize int) {
	store = &logStore{
		maxSize: maxSize,
		data:    make([]Record, 0, maxSize),
	}
}

// AddRecord 添加一条记录，如果超过容量则丢弃最旧的一条
func AddRecord(rec Record) {
	if store == nil {

		InitStore(1000)
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if len(store.data) >= store.maxSize {
		// 丢掉最前面的一条
		store.data = store.data[1:]
	}
	store.data = append(store.data, rec)
}

// GetRecords 返回当前所有记录
func GetRecords() []Record {
	if store == nil {
		return nil
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	out := make([]Record, len(store.data))
	copy(out, store.data)
	return out
}

// nowMillis 返回当前时间的毫秒时间戳
func nowMillis() int64 {
	return time.Now().UnixMilli()
}
