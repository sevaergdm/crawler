package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err	
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return "", fmt.Errorf("An error occurred processing the request. Recieved code %d", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("Content is not html, content is: %s", resp.Header.Get("content-type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
