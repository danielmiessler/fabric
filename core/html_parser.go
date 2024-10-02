package core

import (
	"bytes"
	"fmt"
	"github.com/go-shiori/go-readability"
)

// HtmlReadability Turn any web page into a clean view
// args：
//
//	html (string): full data of web page
//
// return：
//
//	viewContent (string): html main content
//	err (error): parser error
func HtmlReadability(html string) (viewContent string, err error) {
	buf := bytes.NewBufferString(html)
	article, err := readability.FromReader(buf, nil)
	if err != nil {
		return "", err
	}
	fmt.Println("MAIN-CONTENT:", article.TextContent)

	return article.TextContent, nil
}
