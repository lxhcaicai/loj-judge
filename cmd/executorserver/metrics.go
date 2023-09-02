package main

import (
	"github.com/lxhcaicai/loj-judge/env/pool"
	"github.com/lxhcaicai/loj-judge/filestore"
	"sync"
)

type metricsFileStore struct {
	mu sync.Mutex
	filestore.FileStore
	fileSize map[string]int64
}

func newMetricsFileStore(fs filestore.FileStore) filestore.FileStore {
	return &metricsFileStore{
		FileStore: fs,
		fileSize:  make(map[string]int64),
	}
}

type metriceEnvBuilder struct {
	pool.EnvBuilder
}
