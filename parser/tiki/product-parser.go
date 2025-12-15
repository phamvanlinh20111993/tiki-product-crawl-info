package tiki

import (
	"selfstudy/crawl/product/metadata"

	"github.com/PuerkitoBio/goquery"
)

type ProductParser struct{}

func (p ProductParser) Parse(document *goquery.Document) metadata.Product {
	product := metadata.Product{}
	// TODO no need for now
	return product
}
