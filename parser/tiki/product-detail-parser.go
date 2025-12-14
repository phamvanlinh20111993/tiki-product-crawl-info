package tiki

import (
	"github.com/PuerkitoBio/goquery"
	"selfstudy/crawl/product/metadata"
)

type ProductDetailParser struct{}

func (pd ProductDetailParser) Parse(document *goquery.Document) metadata.ProductDetail {
	productDetail := metadata.ProductDetail{}

	return productDetail
}
