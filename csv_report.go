package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	err = writer.Write([]string{
		"page_url", 
		"h1", 
		"first_paragraph", 
		"outgoing_link_urls", 
		"image_urls",
	})
	if err != nil {
		fmt.Errorf("Unable to write to csv file: %v", err)
	}

	for _, page := range pages {
		err = writer.Write([]string{
			page.URL,
			page.H1,
			page.FirstParagraph,
			strings.Join(page.OutgoingLinks, ";"),
			strings.Join(page.ImageURLs, ";"),
		})
		if err != nil {
			fmt.Printf("Unable to write page data for page '%s'\n", page.URL)
			continue
		}
	}
	return nil
}
