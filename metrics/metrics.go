package metrics

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"gopkg.in/Clever/kayvee-go.v6"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

// Log records Golang process metrics such as HeapAlloc, NumGC, etc... every
// frequency time period. This function never returns so it should be called from
// a Goroutine.
func Log(source string, frequency time.Duration) {
	logMetric := func(metricName, metricType string, value uint64) {
		log.Printf(kayvee.FormatLog(source, kayvee.Info, metricName, logger.M{
			"type":  metricType,
			"value": value,
			"via":   "process-metrics",
		}))
	}

	for range time.Tick(frequency) {
		// get a variety of runtime stats
		start := time.Now()
		numGoRoutines := runtime.NumGoroutine()
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		gcStats := &debug.GCStats{
			// allocate 5 slots for 0,25,50,75,100th percentiles
			PauseQuantiles: make([]time.Duration, 5),
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
		for percent := 0.0; percent <= 1.0; percent += 0.25 {
			logMetric(
				fmt.Sprintf("GCPauseNs-%d", int(percent*100)),
				"guage",
				uint64(gcStats.PauseQuantiles[int(percent*4.0)].Nanoseconds()))
		}
	}
}
