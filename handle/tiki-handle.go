package handle

import (
	"encoding/json"
	"selfstudy/crawl/product/configuration"
	httprequest "selfstudy/crawl/product/http-request"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/parser/tiki"
	"selfstudy/crawl/product/util"
)

func crawlHandle() {
	document, _ := httprequest.GetTikiHtmlPage("")
	if document == nil {
		panic("Error when get Tiki html page")
	}
	categoryParser := tiki.CategoryParser{}
	var categories = categoryParser.Parse(document)

	var pageNum = 1
	var totalData int = 0
	for _, category := range categories {
		productResp, err := httprequest.GetTikiProductList(pageNum,
			configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, category.Code)
		if err != nil {
			panic("Error while call product API")
		}
		util.LogDebug(category.Title, ":", productResp.Paging.Total)
		totalData += productResp.Paging.Total
		go getProductDataByCategory(category.Code, productResp.Paging.LastPage)
	}

	util.LogInfo("Total data: ", totalData)
}

func getProductDataByCategory(categoryCode string, lastPage int) {
	for pageNum := 1; pageNum <= lastPage; pageNum++ {
		productResp, err := httprequest.GetTikiProductList(pageNum,
			configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, categoryCode)
		if err != nil {
			util.LogError("Error while call product API", err)
			continue
		}

		if len(productResp.Data) == 0 {
			continue
		}

		for _, product := range productResp.Data {
			byteData, err := json.Marshal(product)
			if err != nil {
				util.LogError("Error while call product API", err)
				continue
			}
			// TODO write to file
			jsonProductData := string(byteData)
			if product.UrlPath != "" && len(product.UrlPath) > 0 {
				// TODO write to file
				jsonProducDetailData := getProductDetailJson(product.UrlPath)
			}
		}
	}

}

func getProductDetailJson(page string) string {
	document, _ := httprequest.GetTikiHtmlPage(page)
	if document == nil {
		panic("Error when get tiki html page")
	}

	productDetailParser := tiki.ProductDetailParser{}
	var productDetail metadata.ProductDetail = productDetailParser.Parse(document)
	byteData, err := json.Marshal(productDetail)
	if err != nil {
		util.LogError("Error while call product API", err)
	}
	return string(byteData)
}

var CrawlHandle = crawlHandle
