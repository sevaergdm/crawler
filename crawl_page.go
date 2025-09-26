package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	parseBase, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Unable to parse baseURL '%s': %v", rawBaseURL, err)
		return
	}

	parseCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to parse currentURL '%s': %v", rawCurrentURL, err)
		return
	}

	if parseBase.Hostname() != parseCurrent.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to normalize URL: %v", err)
		return
	}

	if _, ok := pages[normalizedCurrentURL]; ok {
		pages[normalizedCurrentURL]++
		return
	}

	pages[normalizedCurrentURL] = 1

	fmt.Printf("crawling: %s\n", rawCurrentURL)
	htmlContent, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to get HTML content: %v\n", err)
		return
	}

	contentURLs, err := getURLsFromHTML(htmlContent, parseBase)
	if err != nil {
		fmt.Printf("Unable to get URLs from HTML: %v", err)
		return
	}

	for _, rawURL := range contentURLs {
		crawlPage(rawBaseURL, rawURL, pages)
	}
}
