package main

import (
	"fmt"
	"log"
)

func main() {
	htmlText := `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
		<p>This is an outside paragraph</p>
		<p>This is another paragraph</p>
  </body>
</html>
`

	bodyMain, err := getFirstParagraphFromHTML(htmlText)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bodyMain)
}
