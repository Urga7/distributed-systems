package main

import (
	"fmt"

	"github.com/laspp/PS-2024/vaje/naloga-1/koda/xkcd"
)

func count_words()

func main() {

	comic, err := xkcd.FetchComic(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	var num_comics int = comic.Id
}
