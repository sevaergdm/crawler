package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	inputArgs := os.Args[1:]

	if len(inputArgs) < 1 {
		log.Fatal("no website provided")
	}

	if len(inputArgs) > 1 {
		log.Fatal("too many arguments provided")
	}

	baseURL := inputArgs[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)
	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)

	for k, v := range pages {
		fmt.Printf("Key: %s, Value: %d\n", k, v)
	}
}
