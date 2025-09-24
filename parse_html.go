package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) (string, error) {
	htmlReader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", fmt.Errorf("Unable to parse html: %v", err)
	}

	h1 := doc.Find("h1").Text()
	return h1, nil
}

func getFirstParagraphFromHTML(html string) (string, error) {
	htmlReader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", fmt.Errorf("Unable to parse html: %v", err)
	}

	main := doc.Find("main")
	var firstPar string
	if main.Length() > 0 {
		firstPar = doc.Find("p").First().Text()
	} else {
		firstPar = main.Find("p").First().Text()
	}

	return strings.TrimSpace(firstPar), nil
}
