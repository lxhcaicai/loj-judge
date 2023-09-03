package main

import (
	"github.com/lxhcaicai/loj-judge/env/pool"
	"github.com/lxhcaicai/loj-judge/filestore"
	"github.com/lxhcaicai/loj-judge/worker"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const (
	metricsNamespace   = "executorserver"
	execSubsystem      = "exec"
	filestoreSubsystem = "file"
)

var (

	// 1ms -> 10s
	timeBuckets = []float64{
		0.001, 0.002, 0.005, 0.008, 0.010, 0.025, 0.050, 0.075, 0.1, 0.2,
		0.4, 0.6, 0.8, 1.0, 1.5, 2, 5, 10,
	}

	fileSizeBucket = prometheus.ExponentialBuckets(1<<8, 2, 20)

	execErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Subsystem: execSubsystem,
		Name:      "error_count",
		Help:      "Number of exec query returns error",
	})

	execTimeHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: execSubsystem,
		Name:      "time_seconds",
		Help:      "Histogram for the command execution time",
		Buckets:   timeBuckets,
	}, []string{"status"})

	execMemHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: filestoreSubsystem,
		Name:      "size_bytes",
		Help:      "Histgram for the file size created in the file store",
		Buckets:   fileSizeBucket,
	}, []string{"status"})
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

type metricsEnvPool struct {
	worker.EnvironmentPool
}

func execObserve(res worker.Response) {
	if res.Error != nil {
		execErrorCount.Inc()
	}
	for _, r := range res.Results {
		status := r.Status.String()
		time := r.Time.Seconds()
		memory := float64(r.Memory)

		execTimeHist.WithLabelValues(status).Observe(time)
		execMemHist.WithLabelValues(status).Observe(memory)
	}
}
