package g

import "sync"

// LogEntry 日志条目
type LogEntry struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

// LogRingBuffer 线程安全的环形缓冲区
type LogRingBuffer struct {
	mu   sync.RWMutex
	buf  []LogEntry
	cap  int
	head int
	size int

	subscribers map[int64]chan LogEntry
	subMu       sync.RWMutex
	subIDGen    int64
}

// NewLogRingBuffer 创建环形缓冲区
func NewLogRingBuffer(cap int) *LogRingBuffer {
	return &LogRingBuffer{
		buf:         make([]LogEntry, cap),
		cap:         cap,
		subscribers: make(map[int64]chan LogEntry),
	}
}

// Append 追加日志条目
func (rb *LogRingBuffer) Append(entry LogEntry) {
	rb.mu.Lock()
	rb.buf[rb.head] = entry
	rb.head = (rb.head + 1) % rb.cap
	if rb.size < rb.cap {
		rb.size++
	}
	rb.mu.Unlock()

	rb.subMu.RLock()
	for _, ch := range rb.subscribers {
		select {
		case ch <- entry:
		default:
		}
	}
	rb.subMu.RUnlock()
}

// GetAll 返回所有日志条目（按时间顺序）
func (rb *LogRingBuffer) GetAll() []LogEntry {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	result := make([]LogEntry, rb.size)
	for i := 0; i < rb.size; i++ {
		idx := (rb.head - rb.size + i + rb.cap) % rb.cap
		result[i] = rb.buf[idx]
	}
	return result
}

// Subscribe 创建订阅通道
func (rb *LogRingBuffer) Subscribe() (int64, chan LogEntry) {
	rb.subMu.Lock()
	defer rb.subMu.Unlock()
	rb.subIDGen++
	id := rb.subIDGen
	ch := make(chan LogEntry, 256)
	rb.subscribers[id] = ch
	return id, ch
}

// Unsubscribe 取消订阅
func (rb *LogRingBuffer) Unsubscribe(id int64) {
	rb.subMu.Lock()
	defer rb.subMu.Unlock()
	if ch, ok := rb.subscribers[id]; ok {
		close(ch)
		delete(rb.subscribers, id)
	}
}

// LogBuffer 全局日志缓冲区
var LogBuffer *LogRingBuffer
