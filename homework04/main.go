package main

import (
	"flag"
	"strconv"
)

func main() {
	pidPtr := flag.Int("pid", 0, "# process id")
	numProcessesPtr := flag.Int("n", 2, "total number of processes")
	flag.Parse()
	pid := *pidPtr
	numProcesses := *numProcessesPtr
	chainStorageLength := numProcesses - 1

	const basePort int = 11044

	if pid == 0 {
		putAddr := ":" + strconv.Itoa(basePort+1)
		getAddr := ":" + strconv.Itoa(basePort+chainStorageLength)
		Client(putAddr, getAddr)
	} else {
		addr := ":" + strconv.Itoa(basePort+pid)
		nextNodeAddr := ""
		prevNodeAddr := ""

		if pid < chainStorageLength {
			nextNodeAddr = ":" + strconv.Itoa(basePort+pid+1)
		}

		if pid > 1 {
			prevNodeAddr = ":" + strconv.Itoa(basePort+pid-1)
		}

		Server(addr, nextNodeAddr, prevNodeAddr)
	}
}
