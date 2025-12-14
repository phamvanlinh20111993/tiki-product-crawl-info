package parser

import "github.com/PuerkitoBio/goquery"

type HTMLParser interface {
	Parse(document *goquery.Document)
}
