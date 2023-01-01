# go-process-metrics
A library for tracking Golang process metrics

## Metrics Package
```
// Log records Golang process metrics such as HeapAlloc, NumGC, etc... every
// frequency time period. This function never returns so it should be called from
// a Goroutine.
Log(source string, frequency time.Duration)
```

## Testing
The tests can be run with:
```bash
$ make test
```
test
