package metrics

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"gopkg.in/Clever/kayvee-go.v6/logger"
)

const (
	numPauseQuantiles = 5
)

// Log records Golang process metrics such as HeapAlloc, NumGC, etc... every
// frequency time period. This function never returns so it should be called from
// a Goroutine.
func Log(source string, frequency time.Duration) {
	lg := logger.New("go-process-metrics")

	logMetric := func(metricName, metricType string, value uint64) {
		lg.TraceD(metricName, logger.M{
			"type":  metricType,
			"value": value,
			"via":   "process-metrics",
		})
	}

	for range time.Tick(frequency) {
		// get a variety of runtime stats
		start := time.Now()
		numGoRoutines := runtime.NumGoroutine()
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		gcStats := &debug.GCStats{
			// allocate 5 slots for 0,25,50,75,100th percentiles
			PauseQuantiles: make([]time.Duration, numPauseQuantiles),
		}
		debug.ReadGCStats(gcStats)
		duration := time.Now().Sub(start)

		// Track how long gathering all these metrics actually take so we don't
		// start to lag our services too much.
		// Since this is in it's own go-routine it isn't an exact measure, but
		// something is better than nothing. It will give us the upper bound cost.
		// This may turn important since some calls like ReadMemStats must "stop the world"
		// to gather their information.
		logMetric("MetricsCostNs", "gauge", uint64(duration.Nanoseconds()))

		// log the # of go routines we have running
		logMetric("NumGoroutine", "gauge", uint64(numGoRoutines))

		// log various memory allocation stats
		logMetric("Alloc", "gauge", memStats.Alloc)
		logMetric("HeapAlloc", "gauge", memStats.HeapAlloc)
		logMetric("NumConns", "gauge", getSocketCount())
		logMetric("NumFDs", "gauge", getFDCount())

		// log various GC stats
		logMetric("NumGC", "counter", uint64(gcStats.NumGC))
		// log the min, 25th, 50th, 75th and max GC pause percentiles
		for idx := 0; idx < numPauseQuantiles; idx++ {
			percent := idx * 25
			title := fmt.Sprintf("GCPauseNs-%d", percent)
			logMetric(title, "guage", uint64(gcStats.PauseQuantiles[idx].Nanoseconds()))
		}
	}
}
