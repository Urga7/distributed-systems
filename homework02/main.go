package main

import (
	"fmt"
	"homework02/socialNetwork"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Task = socialNetwork.Task
type Queue = socialNetwork.Q
type Duration = time.Duration
type indexCountEntry struct {
	key   string
	count int
}

const producerDelay int = 5000
const producerActiveTime Duration = time.Millisecond * 150
const controllerDelay Duration = time.Microsecond * 300
const maxWorkers = 64
const minWordLength = 4
const queueCapacity int = 10000
const maxIntensityWorkers int = 128
const numWorkersConst float64 = float64(maxIntensityWorkers) / float64(queueCapacity)

var live bool = false
var stopWorker chan bool
var queueStatistics []float64
var workersStatistics []float64
var timestamps []time.Time
var wg sync.WaitGroup
var lock sync.Mutex

func createPlotData(values []float64, timestamps []time.Time) plotter.XYs {
	pts := make(plotter.XYs, len(values))
	start := timestamps[0]
	for i, v := range values {
		pts[i].X = timestamps[i].Sub(start).Seconds()
		pts[i].Y = v
	}

	return pts
}

func plotGraph(queueLengths []float64, workerCounts []float64, timestamps []time.Time) {
	p := plot.New()
	p.Title.Text = "Queue Length and Active Workers Over Time"
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "Count"

	queueData := createPlotData(queueLengths, timestamps)
	queuePlot, err := plotter.NewLine(queueData)
	if err != nil {
		panic(err)
	}

	queuePlot.Color = plotutil.Color(0)
	p.Add(queuePlot)
	p.Legend.Add("Queue Length in hundreds", queuePlot)

	workerData := createPlotData(workerCounts, timestamps)
	workerPlot, err := plotter.NewLine(workerData)
	if err != nil {
		panic(err)
	}

	workerPlot.Color = plotutil.Color(1)
	p.Add(workerPlot)
	p.Legend.Add("Active Workers", workerPlot)
	p.Legend.Top = true
	p.Legend.Left = false
	if err := p.Save(10*vg.Inch, 5*vg.Inch, "queue_workers_graph.png"); err != nil {
		panic(err)
	}

	fmt.Printf("Graph saved as queue_workers_graph.png\n\n")
}

func cleanseString(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	cleaned := re.ReplaceAllString(s, " ")
	return strings.ToLower(cleaned)
}

func printIndex(index map[string][]uint64) {
	entries := make([]indexCountEntry, 0, len(index))

	for key, ids := range index {
		entries = append(entries, indexCountEntry{key, len(ids)})
	}

	sort.Slice(entries, func(i int, j int) bool {
		return entries[i].count > entries[j].count
	})

	fmt.Printf("\nTop 15 largest entries:\n")
	for i := 0; i < len(entries) && i < 15; i++ {
		fmt.Printf("%d. %s (Appeared in %d posts)\n", i+1, entries[i].key, entries[i].count)
	}

	fmt.Println()
}

func writeResults(mainIndex map[string][]uint64, localIndex map[string][]uint64) {
	for word, ids := range localIndex {
		mainIndex[word] = append(mainIndex[word], ids...)
	}
}

func saveStats(queueLength int, activeWorkers int) {
	queueStatistics = append(queueStatistics, float64(queueLength)/100)
	workersStatistics = append(workersStatistics, float64(activeWorkers))
	timestamps = append(timestamps, time.Now())
}

func calculateRequiredWorkers(queueLength int) int {
	requiredWorkers := int(float64(queueLength) * numWorkersConst)
	if requiredWorkers < 1 {
		return 1
	} else if requiredWorkers > maxWorkers {
		return maxWorkers
	}

	return requiredWorkers
}

func indexer(input <-chan Task, index map[string][]uint64) {
	localIndex := make(map[string][]uint64)
	defer func() {
		lock.Lock()
		writeResults(index, localIndex)
		lock.Unlock()
		wg.Done()
	}()

	for {
		select {
		case <-stopWorker:
			return

		case post, ok := <-input:
			if !ok {
				return
			}

			data := cleanseString(post.Data)
			words := strings.Fields(data)
			for _, word := range words {
				if len(word) >= minWordLength {
					localIndex[word] = append(localIndex[word], post.Id)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "live" {
		live = true
	}

	start := time.Now()
	var producer Queue
	producer.New(producerDelay)
	stopWorker = make(chan bool, maxWorkers)
	taskChannel := producer.TaskChan
	index := make(map[string][]uint64)

	var activeWorkers int
	var producerStopped chan bool = make(chan bool)
	go producer.Run()
	go func() {
		time.Sleep(producerActiveTime)
		producer.Stop()
		producerStopped <- true
		if live {
			fmt.Println("Producer stopped")
		}
	}()

	var producerActive bool = true
	for !producer.QueueEmpty() || producerActive {
		select {
		case <-producerStopped:
			producerActive = false
		default:
		}

		queueLength := len(producer.TaskChan)
		if live && queueLength > 0 {
			fmt.Printf("Queue: %d\n", queueLength)
		}

		saveStats(queueLength, activeWorkers)
		requiredWorkers := calculateRequiredWorkers(queueLength)
		if requiredWorkers > activeWorkers {
			newWorkers := requiredWorkers - activeWorkers
			if live {
				fmt.Printf("Workers: %d (+ %d)\n", activeWorkers, newWorkers)
			}

			wg.Add(newWorkers)
			for i := 0; i < newWorkers; i++ {
				go indexer(taskChannel, index)
			}
		} else if requiredWorkers < activeWorkers {
			obsoleteWorkers := activeWorkers - requiredWorkers
			if live {
				fmt.Printf("Workers: %d (- %d)\n", activeWorkers, obsoleteWorkers)
			}

			for i := 0; i < obsoleteWorkers; i++ {
				stopWorker <- true
			}
		}

		activeWorkers = requiredWorkers
		time.Sleep(controllerDelay)
	}

	for i := 0; i < activeWorkers; i++ {
		stopWorker <- true
	}

	wg.Wait()
	close(stopWorker)
	elapsed := time.Since(start)

	printIndex(index)
	fmt.Printf("Elapsed time: %f\n", elapsed.Seconds())
	fmt.Printf("Processing rate: %f MReqs/s\n", float64(producer.N)/float64(elapsed.Seconds())/1000000.0)
	fmt.Printf("Average queue length: %.2f %%\n", producer.GetAverageQueueLength())
	fmt.Printf("Max queue length: %.2f %%\n", producer.GetMaxQueueLength())
	fmt.Printf("Number of produced posts: %d\n\n", producer.N)
	plotGraph(queueStatistics, workersStatistics, timestamps)
}
