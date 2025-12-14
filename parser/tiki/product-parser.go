package tiki

import (
	"github.com/PuerkitoBio/goquery"
	"selfstudy/crawl/product/metadata"
)

type ProductParser struct{}

func (p ProductParser) Parse(document *goquery.Document) metadata.Product {
	product := metadata.Product{}

	return product
}
