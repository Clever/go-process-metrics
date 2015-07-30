package metrics

import (
	"log"
	"runtime"
	"time"

	"gopkg.in/Clever/kayvee-go.v2"
)

func logMetric(source, metricName, metricType string, value int) {
	// TODO: Add ENV???
	payload := map[string]interface{}{"type": metricType, "value": value}
	log.Printf(kayvee.FormatLog(source, kayvee.Info, metricName, payload))
}

// TODO: Add a nice comment!!!
// TODO: Note that this never returns :)
func Log(source string, frequency time.Duration) {
	for {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		// TODO: Should I add any others???
		logMetric(source, "HeapAlloc", "gauge", memStats.HeapAlloc)
		logMetric(source, "HeapInuse", "gauge", memStats.HeapInuse)
		logMetric(source, "NumGC", "gauge", memStats.NumGC)
		logMetric(source, "PauseTotalMs", "gauge", memStats.PauseTotalNs/1000000)

		time.Sleep(frequency)
	}
}
