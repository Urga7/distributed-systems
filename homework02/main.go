package main

import (
	"fmt"
	"homework02/socialNetwork"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Task = socialNetwork.Task
type Queue = socialNetwork.Q

const queueCapacity int = 10000
const hardMaxWorkers int = 128
const numWorkersConst int = hardMaxWorkers / queueCapacity

var wg sync.WaitGroup
var lock sync.Mutex

func cleanseString(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	cleaned := re.ReplaceAllString(s, " ")
	return strings.ToLower(cleaned)
}

func writeResults(mainIndex map[string][]uint64, localIndex map[string][]uint64) {
	for word, ids := range localIndex {
		mainIndex[word] = append(mainIndex[word], ids...)
	}
}

func indexer(input <-chan Task, index map[string][]uint64, stopChannel <-chan struct{}) {
	localIndex := make(map[string][]uint64)
	defer wg.Done()
	defer func() {
		lock.Lock()
		writeResults(index, localIndex)
		lock.Unlock()
	}()

	for {
		select {
		case post, ok := <-input:
			if !ok {
				return
			}

			data := cleanseString(post.Data)
			words := strings.Fields(data)
			for _, word := range words {
				if len(word) >= 4 {
					localIndex[word] = append(localIndex[word], post.Id)
				}
			}

		case <-stopChannel:
			return
		}
	}
}

func controller(producer Queue, index map[string][]uint64, maxWorkers int) {
	fmt.Println("Controller launched")
	stopChannel := make(chan struct{})
	taskChannel := producer.TaskChan
	var activeWorkers int = 1
	wg.Add(1)
	go indexer(taskChannel, index, stopChannel)
	fmt.Println("First indexer started, starting producer...")
	producer.Run()
	for !producer.QueueEmpty() && len(taskChannel) > 0 {
		queueLength := len(producer.TaskChan)
		fmt.Printf("Currently %d items in queue\n", queueLength)
		requiredWorkers := queueLength * numWorkersConst
		if requiredWorkers < 1 {
			requiredWorkers = 1
		} else if requiredWorkers > maxWorkers {
			requiredWorkers = maxWorkers
		}

		if requiredWorkers > activeWorkers {
			newWorkers := requiredWorkers - activeWorkers
			fmt.Printf("There are %d active workers. Adding %d new workers.", activeWorkers, newWorkers)
			wg.Add(newWorkers)
			for i := 0; i < newWorkers; i++ {
				go indexer(taskChannel, index, stopChannel)
			}
		} else if requiredWorkers < activeWorkers {
			obsoleteWorkers := activeWorkers - requiredWorkers
			fmt.Printf("There are %d active workers. Removing %d workers.", activeWorkers, obsoleteWorkers)
			for i := 0; i < obsoleteWorkers; i++ {
				stopChannel <- struct{}{}
			}
		}

		activeWorkers = requiredWorkers
		time.Sleep(time.Millisecond)
	}

	close(taskChannel)
	for i := 0; i < activeWorkers; i++ {
		stopChannel <- struct{}{}
	}

	close(stopChannel)
	wg.Wait()
}

func main() {
	start := time.Now()
	var producer Queue
	producer.New(5000)
	index := make(map[string][]uint64)
	controller(producer, index, 64)
	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %f", elapsed.Seconds())
}
