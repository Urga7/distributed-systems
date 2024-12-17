package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

var Logger *govec.GoLog
var opts govec.GoLogOptions
var pid int
var numProcesses int

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

func listen(addr *net.UDPAddr, persistance time.Duration, pid int, growthRate int, numProcesses int, rootPort int) {
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()

	timeout := time.Now().Add(persistance)
	err = conn.SetReadDeadline(timeout)
	checkError(err)

	buffer := make([]byte, 1024)
	var msg []byte
	receivedFirstMsg := false
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				return
			}

			checkError(err)
		}

		Logger.UnpackReceive("Prejeto sporocilo ", buffer, &msg, opts)
		mLen := len(msg)

		if !receivedFirstMsg {
			receivedFirstMsg = true
			rMsg := message{
				data:   msg[:mLen],
				length: mLen,
			}

			randomProcesses := getRandomProcesses(pid, growthRate, numProcesses)
			for _, processID := range randomProcesses {
				remotePort := rootPort + 1 + processID
				remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", remotePort))
				checkError(err)
				send(remoteAddr, rMsg)
			}
		}
	}
}

func send(addr *net.UDPAddr, msg message) {
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()

	Logger.LogLocalEvent("Priprava sporocila ", opts)
	sMsg := msg.data[:msg.length]
	sMsgVC := Logger.PrepareSend("Poslano sporocilo ", []byte(sMsg), opts)
	_, err = conn.Write(sMsgVC)
	checkError(err)
}

func sendInitialMessages(numMsg int, rootPort int, sendDelay time.Duration) {
	time.Sleep(sendDelay)
	randomProcesses := getRandomProcesses(pid, numMsg, numProcesses)
	for i := 1; i <= numMsg; i++ {
		msgData := fmt.Sprintf("%d", i)
		msg := message{
			data:   []byte(msgData),
			length: len(msgData),
		}

		processID := randomProcesses[i-1]
		remotePort := rootPort + processID
		remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", remotePort))
		checkError(err)
		send(remoteAddr, msg)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	pidPtr := flag.Int("pid", 0, "# process id")
	numProcessesPtr := flag.Int("n", 2, "total number of processes")
	numMsgsPtr := flag.Int("m", 3, "number of messages sent from the first process")
	growthRatePtr := flag.Int("k", 3, "number of messages sent from every other process upon receiving it")
	flag.Parse()

	pid = *pidPtr
	numProcesses = *numProcessesPtr
	numMsg := *numMsgsPtr
	growthRate := min(*growthRatePtr, numProcesses-1)

	var rootPort = 11044
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort+pid))
	checkError(err)

	Logger = govec.InitGoVector("Process-"+strconv.Itoa(pid), "Log-Process-"+strconv.Itoa(pid), govec.GetDefaultConfig())
	opts = govec.GetDefaultLogOptions()

	const persistance = time.Second
	const sendDelay = time.Millisecond * 300

	if pid == 0 {
		sendInitialMessages(numMsg, rootPort, sendDelay)
	}

	listen(localAddr, persistance, pid, growthRate, numProcesses, rootPort)
}
