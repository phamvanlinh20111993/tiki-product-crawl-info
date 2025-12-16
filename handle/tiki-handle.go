package handle

import (
	"encoding/json"
	"math"
	"math/rand"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/datasource/file"
	httprequest "selfstudy/crawl/product/http-request"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/parser/tiki"
	"selfstudy/crawl/product/util"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func crawlHandle() {
	document, _ := httprequest.GetTikiHtmlPage("")
	if document == nil {
		panic("Error when get Tiki html page")
	}

	categoryParser := tiki.CategoryParser{}
	var categories = categoryParser.Parse(document)

	var totalData int = 0
	var wg sync.WaitGroup
	var currentCountData int32 = 0
	doneChan := make(chan bool, len(categories))

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	categoryFile := file.NewFileDataSource(configuration.GetFileConfig().Name + "categories")
	byteData, err := json.Marshal(categories)
	if err != nil {
		util.LogError("Error while write json to file", err)
	} else {
		categoryFile.Insert(string(byteData))
		categoryFile.Close()
	}

	for i := 1; i < len(categories); i++ {
		category := categories[i]
		productResp, err := httprequest.GetTikiProductList(1,
			configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, category.Code)
		if err != nil {
			util.LogError("Error while call product API")
			continue
		}
		util.LogDebug(category.Title, ": total product", productResp.Paging.Total)
		totalData += productResp.Paging.Total
		//	wg.Go(func() {
		getProductDataByCategory(category, productResp.Paging.LastPage, &currentCountData, &wg, doneChan, random)
		//	})
	}
	util.LogInfo("Total product data expected to crawl: ", totalData)

	/*	wg.Add(1)
		go func() {
			defer wg.Done()
			doneCount := 0
			for {
				select {
				case <-doneChan:
					doneCount++
				default:
					if doneCount == len(categories) {
						break
					}
					util.LogInfo("^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^^ Current amount of products data crawling: ", currentCountData, " ^^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^")
					time.Sleep(20 * time.Second)
				}
			}
		}()

		wg.Wait() */
	//	close(doneChan)
}

func getProductDataByCategory(category metadata.Category, lastPage int, currentCountProduct *int32, wg *sync.WaitGroup, doneChan chan bool, random *rand.Rand) {
	productFile := file.NewFileDataSource(configuration.GetFileConfig().Name + category.Title + "-" + category.Code)
	productFileDetail := file.NewFileDataSource(configuration.GetFileConfig().Name + category.Title + "-" + category.Code + "-Detail")

	defer productFile.Close()
	defer productFileDetail.Close()
	//defer wg.Done()

	for pageNum := 1; pageNum <= lastPage; pageNum++ {
		util.LogInfo("@@@@@@@@@@@@@@@@@@@@@@@@@", category.Title, ": page Number ", pageNum, "@@@@@@@@@@@@@@@@@@@@@@@@@")
		productResp, err := httprequest.GetTikiProductList(pageNum,
			configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, category.Code)
		if err != nil {
			util.LogError("Error while call product API", err)
			continue
		}

		if len(productResp.Data) == 0 {
			util.LogInfo("Product Data is empty", category.Title, ", At page", pageNum)
			continue
		}

		var i int = 0
		var errorCount int = 2
		var exponential float64 = 2.0 //2.7 + random.Float64()*(6.5-2.7)
		for i = 0; i < len(productResp.Data); {
			product := productResp.Data[i]
			byteData, err := json.Marshal(product)
			if err != nil {
				util.LogError("Error while call product API", err)
				continue
			}
			if product.UrlPath != "" && len(product.UrlPath) > 0 {
				jsonProductDetailData := getProductDetailJson(product.UrlPath)

				if jsonProductDetailData != "" {
					atomic.AddInt32((*int32)(currentCountProduct), 1)
					jsonProductData := string(byteData)
					productFile.Insert(jsonProductData)
					productFileDetail.Insert(jsonProductDetailData)
					i++
					errorCount = 1
				} else {
					// more than 30s
					if errorCount > 30 {
						util.LogInfo("We cant do request forever, errorCount = ", errorCount)
						continue
					}
					util.LogInfo("Start retry with duration ", time.Duration(errorCount)*time.Second)
					time.Sleep(time.Duration(errorCount) * time.Second)
					errorCount = int(math.Round(exponential * float64(errorCount)))
				}
			}
			// time.Sleep(250 * time.Microsecond)
		}

	}
	//	doneChan <- true
}

func getProductDetailJson(page string) string {
	document, err := httprequest.GetTikiHtmlPage(page)
	if document == nil || document == (&goquery.Document{}) || err != nil {
		util.LogError("Error when get Tiki html page")
		return ""
	}

	productDetailParser := tiki.ProductDetailParser{}
	var productDetail metadata.ProductDetail = productDetailParser.Parse(document)
	if productDetail.ProductId <= 0 {
		util.LogDebug("Product Id is empty at page ", page)
		return ""
	}
	byteData, err := json.Marshal(productDetail)
	if err != nil {

		util.LogError("Error while call product API", err)
	}
	return string(byteData)
}

var CrawlHandle = crawlHandle
