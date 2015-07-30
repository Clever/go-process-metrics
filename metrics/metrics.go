package metrics

import (
	"log"
	"runtime"
	"time"

	"gopkg.in/Clever/kayvee-go.v2"
)

func logMetric(source, metricName, metricType string, value uint64) {
	payload := map[string]interface{}{"type": metricType, "value": value}
	log.Printf(kayvee.FormatLog(source, kayvee.Info, metricName, payload))
}

// Log records Golang process metrics such as HeapAlloc, NumGC, etc... every
// frequency time period. This function never returns so it should be called from
// a Goroutine.
func Log(source string, frequency time.Duration) {
	for _ = range time.Tick(frequency) {
		logMetric(source, "NumGoroutine", "gauge", uint64(runtime.NumGoroutine()))

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		logMetric(source, "Alloc", "gauge", memStats.Alloc)
		logMetric(source, "HeapAlloc", "gauge", memStats.HeapAlloc)
		logMetric(source, "NumGC", "counter", uint64(memStats.NumGC))
		logMetric(source, "PauseTotalMs", "counter", memStats.PauseTotalNs/1000000)
	}
}
