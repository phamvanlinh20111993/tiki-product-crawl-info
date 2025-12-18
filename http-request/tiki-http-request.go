package http_request

import (
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/logger"
	"selfstudy/crawl/product/metadata"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func getTikiProductList(pageNum int, limit int, category string) (metadata.ProductAPIResponse, error) {
	var requestParams map[string]map[string]string = map[string]map[string]string{}
	requestParams[QueryParams] = map[string]string{
		"pageNum":  strconv.Itoa(pageNum),
		"limit":    strconv.Itoa(limit),
		"category": category,
	}
	requestParams[Headers] = map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
		"sec-ch-ua":  "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"",
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}

	tikiProductListResponse, err := GetAPIData[metadata.ProductAPIResponse](configuration.GetTikiPageConfig().ProductAPIURL, requestParams)

	return tikiProductListResponse, err
}

func getTikiHTMLPage(path string) (*goquery.Document, error) {
	var requestParams map[string]map[string]string = map[string]map[string]string{}
	requestParams[Headers] = map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
		"sec-ch-ua":  "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"",
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}
	doc, err := getHTMLPage(configuration.GetTikiPageConfig().BaseURL+path, requestParams)

	if err != nil {
		logger.LogError("Error while parsing html response: ", err)
	}
	return doc, err
}

func getTikiProductCategoryPathList(category string) (metadata.CategoryResponse, error) {
	var requestParams map[string]map[string]string = map[string]map[string]string{}
	requestParams[QueryParams] = map[string]string{
		"include":   "children",
		"parent_id": category,
	}
	requestParams[Headers] = map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
		"sec-ch-ua":  "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"",
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}

	tikiCategoryPathResponse, err := GetAPIData[metadata.CategoryResponse](configuration.GetTikiPageConfig().CategoryPathAPIRL, requestParams)

	return tikiCategoryPathResponse, err
}

// GetTikiProductList GetTikiProducts upper the first letter make it public outside of package
var GetTikiProductList = getTikiProductList
var GetTikiHtmlPage = getTikiHTMLPage
var GetTikiProductCategoryPathList = getTikiProductCategoryPathList
