package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.hitPageLimit() {
		return
	}

	parseCurrent, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to parse currentURL '%s': %v", rawCurrentURL, err)
		return
	}

	if cfg.baseURL.Hostname() != parseCurrent.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to normalize URL: %v", err)
		return
	}

	if isFirst := cfg.addPageVisit(normalizedCurrentURL); !isFirst {
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)
	htmlContent, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Unable to get HTML content: %v\n", err)
		return
	}

	pageData := extractPageData(htmlContent, rawCurrentURL)
	cfg.setPageData(normalizedCurrentURL, pageData)

	for _, rawURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(rawURL)
	}
}
