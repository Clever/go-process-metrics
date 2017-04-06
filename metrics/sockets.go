package metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/Clever/kayvee-go.v6/logger"
)

var (
	lg = logger.New("go-process-metrics")
)

// getSocketCount returns the number of open sockets
func getSocketCount() uint64 {
	var socketCount uint64
	pid := os.Getpid()
	base := fmt.Sprintf("/proc/%d/fd", pid)
	fds, err := ioutil.ReadDir(base)
	if err != nil {
		lg.ErrorD("failed-get-fds", logger.M{"err": err.Error()})
		return 0
	}

	hasLogged := false
	for _, fd := range fds {
		sl, err := os.Readlink(fmt.Sprintf("%s/%s", base, fd.Name()))
		// The fd opened by ReadDir is often closed at this point.
		// In general, fds may close between the check and readlink and we should ignore those
		if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
			if !hasLogged {
				lg.ErrorD("failed-readlink", logger.M{"err": err.Error()})
				hasLogged = true
			}
			continue
		}

		if strings.Contains(sl, "socket") {
			socketCount++
		}
	}
	return socketCount
}

func getFDCount() uint64 {
	pid := os.Getpid()
	base := fmt.Sprintf("/proc/%d/fd", pid)
	fds, err := ioutil.ReadDir(base)
	if err != nil {
		lg.ErrorD("failed-get-fds", logger.M{"err": err.Error()})
		return 0
	}
	return uint64(len(fds))
}
