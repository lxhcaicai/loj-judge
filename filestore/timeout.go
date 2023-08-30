package filestore

import (
	"container/heap"
	"sync"
	"time"
)

var (
	_ FileStore      = &Timeout{}
	_ heap.Interface = &Timeout{}
)

// 文件系统是否具有最大TTL
type Timeout struct {
	mu sync.Mutex
	FileStore
	timeout   time.Duration
	files     []timeoutFile
	idToIndex map[string]int
}

type timeoutFile struct {
	id   string
	time time.Time
}

// 为文件创建一个具有最大TTL的超时文件系统
func NewTimeout(fs FileStore, timeout time.Duration, checkInterval time.Duration) FileStore {
	t := &Timeout{
		FileStore: fs,
		timeout:   timeout,
		files:     make([]timeoutFile, 0),
		idToIndex: make(map[string]int),
	}
	go t.checkTimeoutLoop(checkInterval)
	return t
}

func (t *Timeout) checkTimeoutLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		t.checkTimeoutAndRemove()
		<-ticker.C
	}
}

func (t *Timeout) checkTimeoutAndRemove() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	for len(t.files) > 0 && t.files[0].time.Add(t.timeout).Before(now) {
		f := t.files[0]
		t.FileStore.Remove(f.id)
		heap.Pop(t)
	}
}

func (t *Timeout) Len() int {
	//TODO implement me
	panic("implement me")
}

func (t *Timeout) Less(i, j int) bool {
	//TODO implement me
	panic("implement me")
}

func (t *Timeout) Swap(i, j int) {
	//TODO implement me
	panic("implement me")
}

func (t *Timeout) Push(x any) {
	//TODO implement me
	panic("implement me")
}

func (t *Timeout) Pop() any {
	//TODO implement me
	panic("implement me")
}
