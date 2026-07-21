package parser

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Page struct {
	URL   string
	Title string
	Text  string
	Links []string
}

func Parse(pageURL string, body []byte) (*Page, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parser: parse html: %w", err)
	}
	var text strings.Builder

	page := &Page{URL: pageURL}
	walk(doc, page, &text, pageURL)
	page.Text = strings.TrimSpace(text.String())
	return page, nil
}

func walk(node *html.Node, page *Page, text *strings.Builder, url string) {
	if node == nil {
		return
	}

	if node.Type == html.ElementNode {
		switch node.Data {
		case "script", "style", "noscript":
			return
		}
	}
	if node.Type == html.TextNode {
		text.WriteString(node.Data)
		text.WriteByte(' ')
	}

	if node.Type == html.ElementNode && node.Data == "title" && page.Title == "" {
		var titleSB strings.Builder
		collectText(node, &titleSB)
		page.Title = strings.TrimSpace(titleSB.String())
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				link := ResolveLink(url, attr.Val)
				if link != "" {
					page.Links = append(page.Links, link)
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		walk(c, page, text, url)
	}
}

func collectText(node *html.Node, sb *strings.Builder) {
	if node.Type == html.TextNode {
		sb.WriteString(node.Data)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, sb)
	}
}
