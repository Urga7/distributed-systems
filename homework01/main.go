package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/laspp/PS-2024/vaje/naloga-1/koda/xkcd"
)

type Comic = xkcd.Comic
type WordCount struct {
	Word  string
	Count int
}

var wg sync.WaitGroup
var lock sync.Mutex
var num_comics int
var processed_comics int

const num_displayed_words int = 15

func print_map(m map[string]int) {
	var wordCountPairs []WordCount
	for word, count := range m {
		wordCountPairs = append(wordCountPairs, WordCount{Word: word, Count: count})
	}

	sort.Slice(wordCountPairs, func(i, j int) bool {
		if wordCountPairs[i].Count == wordCountPairs[j].Count {
			return wordCountPairs[i].Word < wordCountPairs[j].Word
		}
		return wordCountPairs[i].Count > wordCountPairs[j].Count
	})

	for i, pair := range wordCountPairs {
		if i >= num_displayed_words {
			break
		}
		fmt.Printf("%s, %d\n", pair.Word, pair.Count)
	}
}

func cleanse_string(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	cleaned := re.ReplaceAllString(s, " ")
	return strings.ToLower(cleaned)
}

func extract_text(c Comic) string {
	if c.Transcript == "" {
		return strings.Join([]string{c.Title, c.Tooltip}, " ")
	}

	return strings.Join([]string{c.Title, c.Transcript}, " ")
}

func get_num_comics() int {
	comic, err := xkcd.FetchComic(0)
	if err != nil {
		fmt.Println(err)
	}

	return comic.Id
}

func get_unprocessed_comic_id() int {
	var comic_id int
	lock.Lock()
	if processed_comics < num_comics {
		processed_comics++
		comic_id = processed_comics
	} else {
		comic_id = -1
	}

	lock.Unlock()
	return comic_id
}

func write_results(main_map *map[string]int, temp_map map[string]int) {
	for word, count := range temp_map {
		(*main_map)[word] += count
	}
}

func count_words(word_count_map *map[string]int) {
	defer wg.Done()
	for {
		comic_id := get_unprocessed_comic_id()
		if comic_id == -1 {
			break
		}

		comic, err := xkcd.FetchComic(comic_id)
		if err != nil {
			fmt.Printf("Error fetching comic %d: %s\n", comic_id, err)
			continue
		}

		comic_text := extract_text(comic)
		cleaned_comic := cleanse_string(comic_text)
		word_occurrences := make(map[string]int)
		words := strings.Fields(cleaned_comic)
		for _, word := range words {
			if len(word) >= 4 {
				word_occurrences[word]++
			}
		}

		lock.Lock()
		write_results(word_count_map, word_occurrences)
		lock.Unlock()
		//fmt.Printf("Comic %d done\n", comic_id)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <comic_id>")
		return
	}

	num_workers, err := strconv.Atoi(os.Args[1])
	if err != nil || num_workers <= 0 {
		fmt.Println("Invalid input. Please provide a valid positive integer for workers.")
		return
	}

	num_comics = get_num_comics()
	processed_comics = 0
	word_count_map := make(map[string]int)
	for i := 0; i < num_workers; i++ {
		wg.Add(1)
		go count_words(&word_count_map)
	}

	wg.Wait()
	print_map(word_count_map)
}
