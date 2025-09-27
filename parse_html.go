package main

import (
	"fmt"
	"net/url"
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
		firstPar = main.Find("p").First().Text()
	} else {
		firstPar = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(firstPar), nil
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("Unable parse html: %v", err)
	}

	hrefs := []string{}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || strings.TrimSpace(href) == "" {
			return
		}

		parsed, err := url.Parse(href)
		if err != nil {
			return
		}
		hrefURL := baseURL.ResolveReference(parsed)
		hrefs = append(hrefs, hrefURL.String())
	})
	return hrefs, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("Unable parse html: %v", err)
	}

	imgs := []string{}
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok || strings.TrimSpace(src) == "" {
			return
		}

		parsed, err := url.Parse(src)
		if err != nil {
			return
		}
		imgURL := baseURL.ResolveReference(parsed)
		imgs = append(imgs, imgURL.String())
	})
	return imgs, nil

}

func extractPageData(html, pageURL string) PageData {
	header, _ := getH1FromHTML(html)
	firstParagraph, _ := getFirstParagraphFromHTML(html)
	parsedURL, _ := url.Parse(pageURL)
	outgoingLinks, _ := getURLsFromHTML(html, parsedURL)
	imageURLs, _ := getImagesFromHTML(html, parsedURL)

	return PageData{
		URL: pageURL,
		H1: header,
		FirstParagraph: firstParagraph,
		OutgoingLinks: outgoingLinks,
		ImageURLs: imageURLs,
	}
}
