package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func main() {
	inputArgs := os.Args[1:]

	if len(inputArgs) < 1 {
		log.Fatal("no arguments provided")
	}

	if len(inputArgs) > 3 {
		log.Fatal("too many arguments provided")
	}

	baseURL := inputArgs[0]
	maxConcurrency, err := strconv.Atoi(inputArgs[1])
	if err != nil {
		log.Fatalf("maxConcurrency must be an integer: %v", err)
	}
	maxPages, err := strconv.Atoi(inputArgs[2])
	if err != nil {
		log.Fatalf("maxPages must be an integer: %v", err)
	}

	cfg, err := configure(baseURL, maxConcurrency, maxPages)
	if err != nil {
		log.Fatalf("Unable to create configuration: %v", err)
	}

	fmt.Printf("starting crawl of: %s\n", baseURL)
	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	cfg.wg.Wait()

	writeCSVReport(cfg.pages, "report.csv")
}
