package tiki

import (
	"regexp"
	"selfstudy/crawl/product/metadata"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CategoryParser struct{}

func (t CategoryParser) Parse(document *goquery.Document) []metadata.Category {
	var categories []metadata.Category

	document.Find("main a").Each(func(_ int, tagA *goquery.Selection) {

		link, _ := tagA.Attr("href")

		match, _ := regexp.MatchString("^\\s*(/[\\w\\-]+){2}\\s*$", link)
		if !match {
			return
		}

		linkSplit := strings.Split(link, "/")
		var code = ""
		splLen := len(linkSplit)
		if splLen > 2 {
			code = linkSplit[splLen-1][1:]
		}

		imagePresent, isExist := tagA.Find("div picture source").Attr("srcset")
		if isExist {
			imagePresentSplit := strings.Split(imagePresent, " ")
			imgLen := len(imagePresentSplit)
			categoryImagePresentation := ""
			if imgLen > 0 {
				categoryImagePresentation = imagePresentSplit[0]
			}

			title, _ := tagA.Attr("title")

			category := metadata.Category{Code: code, Title: title,
				Path: link, CategoryImagePresentation: categoryImagePresentation}
			categories = append(categories, category)
		}

	})

	return categories
}
