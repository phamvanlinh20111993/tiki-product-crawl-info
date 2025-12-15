package parser

import "github.com/PuerkitoBio/goquery"

type HTMLParser interface {
	Parse(document *goquery.Document) interface{}
}

func ParseData(htmlParser HTMLParser, document *goquery.Document) interface{} {
	return htmlParser.Parse(document)
}
