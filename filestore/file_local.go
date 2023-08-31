package filestore

import (
	"os"
	"path"
	"path/filepath"
	"sync"
)

type fileLocalStore struct {
	dir  string            // 存放文件的目录
	name map[string]string // 如果存在，Id到名称映射
	mu   sync.RWMutex
}

func (s *fileLocalStore) List() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fi, err := os.ReadDir(s.dir)
	if err != nil {
		return nil
	}

	names := make(map[string]string, len(fi))
	for _, f := range fi {
		names[f.Name()] = s.name[f.Name()]
	}
	return names
}

func (s *fileLocalStore) Remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.name, id)
	p := path.Join(s.dir, id)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	os.Remove(p)
	return true
}

// 创建新的本地文件存储
func NewFileLocalStore(dir string) FileStore {
	return &fileLocalStore{
		dir:  filepath.Clean(dir),
		name: make(map[string]string),
	}
}
