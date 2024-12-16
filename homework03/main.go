package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type message struct {
	length int
	data   []byte
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getRandomProcesses(pid, growthRate, numProcesses int) []int {
	selected := make(map[int]struct{})
	var result []int
	for len(result) < growthRate {
		r := rand.Intn(numProcesses)
		if r != pid {
			_, exists := selected[r]
			if !exists {
				selected[r] = struct{}{}
				result = append(result, r)
			}
		}
	}

	return result
}

func listen(addr *net.UDPAddr, timelimitMs int, pid int, growthRate int, numProcesses int, rootPort int) {
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()

	timeout := time.Now().Add(time.Millisecond * time.Duration(timelimitMs))
	err = conn.SetReadDeadline(timeout)
	checkError(err)

	buffer := make([]byte, 1024)
	var receivedFirstMsg bool = false
	for {
		mLen, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				return
			}

			checkError(err)
		}

		if !receivedFirstMsg {
			receivedFirstMsg = true
			msg := buffer[:mLen]
			fmt.Printf("%s", string(msg))
			response := message{data: msg, length: mLen}
			randomProcesses := getRandomProcesses(pid, growthRate, numProcesses)
			for _, processID := range randomProcesses {
				remotePort := rootPort + 1 + processID
				remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", remotePort))
				checkError(err)
				send(remoteAddr, response)
			}
		}
	}
}

func send(addr *net.UDPAddr, msg message) {
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()

	sMsg := string(msg.data[:msg.length])
	_, err = conn.Write([]byte(sMsg))
	checkError(err)
	fmt.Printf("Sent msg to address: %s\n", addr.IP.String())
	fmt.Printf("Sent msg to port: %d\n", addr.Port)
}

func main() {
	pidPtr := flag.Int("pid", 0, "# process id")
	numProcessesPtr := flag.Int("n", 2, "total number of processes")
	numMsgsPtr := flag.Int("m", 3, "number of messages sent from the first process")
	growthRatePtr := flag.Int("k", 3, "number of messages sent from every other process upon receiving it")
	flag.Parse()

	pid := *pidPtr
	numProcesses := *numProcessesPtr
	numMsg := *numMsgsPtr
	growthRate := min(*growthRatePtr, numProcesses-1)

	fmt.Printf("Process with id %d started\n", pid)

	const rootPort = 9000
	basePort := rootPort + 1 + pid
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", basePort))
	checkError(err)

	fmt.Printf("My address: %s\n", localAddr.IP.String())
	fmt.Printf("My port: %d\n", localAddr.Port)

	if pid == 0 {
		time.Sleep(time.Second)
		randomProcesses := getRandomProcesses(pid, numMsg, numProcesses)
		for i := 1; i <= numMsg; i++ {
			msg := message{
				data:   []byte(fmt.Sprintf("%d", i)),
				length: len(fmt.Sprintf("%d", i)),
			}

			processID := randomProcesses[i-1]
			remotePort := rootPort + 1 + processID
			remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", remotePort))
			checkError(err)
			send(remoteAddr, msg)
			time.Sleep(100 * time.Millisecond)
		}
	}

	listen(localAddr, 300, pid, growthRate, numProcesses, rootPort)
}
