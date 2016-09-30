package metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// getSocketCount returns the number of open sockets
func getSocketCount() uint64 {
	var socketCount uint64
	pid := os.Getpid()
	base := fmt.Sprintf("/proc/%d/fd", pid)
	fds, err := ioutil.ReadDir(base)
	if err != nil {
		// Should we log something here?
		return 0
	}

	for _, fd := range fds {
		sl, err := os.Readlink(fmt.Sprintf("%s/%s", base, fd.Name()))
		if err != nil {
			// Should we log something here?
			continue
		}

		if strings.Contains(sl, "socket") {
			socketCount++
		}
	}
	return socketCount
}
