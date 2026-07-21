package parser

import (
	"net/url"
)

func ResolveLink(pageURL string, href string) string {
	base, err := url.Parse(pageURL)
	if err != nil {
		return ""
	}

	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	abs := base.ResolveReference(ref)
	if abs.Scheme != "http" && abs.Scheme != "https" {
		return ""
	}
	abs.Fragment = ""
	return abs.String()
}
