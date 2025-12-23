package handle

import (
	"encoding/json"
	"math"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/datasource"
	"selfstudy/crawl/product/datasource/file"
	httprequest "selfstudy/crawl/product/http-request"
	"selfstudy/crawl/product/logger"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/parser/tiki"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TikiCrawlHandler struct {
	input  any // TODO
	output []datasource.DatasourceI
}

type TransferProductDetail struct {
	productFile       *file.FileDataSource
	productFileDetail *file.FileDataSource
	product           *metadata.Product
	isEnd             bool
}

var (
	productDataQueue    chan TransferProductDetail
	currentCountProduct int32
)

func notify(doneChan chan bool, categoryL int) {
	doneCount := 0
	for {
		select {
		case <-doneChan:
			doneCount++
		default:
			if doneCount == categoryL {
				break
			}
			logger.LogInfo("^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^^ Current amount of products data crawling: ", currentCountProduct, " ^^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^")
			time.Sleep(20 * time.Second)
		}
	}
}

func (crawl TikiCrawlHandler) CrawlHandle() {
	productDataQueue = make(chan TransferProductDetail, 1000)
	currentCountProduct = 0

	document, _ := httprequest.GetTikiHtmlPage("")
	if document == nil {
		panic("Error when get Tiki html page")
	}

	categoryParser := tiki.CategoryParser{}
	var categories = categoryParser.Parse(document)
	var totalData int = 0
	doneChan := make(chan bool, len(categories))
	// get category file
	categoryFile := file.NewFileDataSource(configuration.GetFileConfig().PrefixName + "categories")
	byteData, err := json.Marshal(categories)
	if err != nil {
		logger.LogError("Error while write json to file", err)
	} else {
		categoryFile.Insert(string(byteData))
		categoryFile.Close()
	}
	// get category file
	categoryFilePath := file.NewFileDataSource(configuration.GetFileConfig().PrefixName + "categories-path")
	defer categoryFilePath.Close()
	for _, category := range categories {
		categoryPaths, err := httprequest.GetTikiProductCategoryPathList(category.Code)
		if err != nil {
			logger.LogError("Error while get category path api", err)
			continue
		}
		byteData, err := json.Marshal(categoryPaths)
		if err != nil {
			logger.LogError("Error while write json to file", err)
			continue
		}
		categoryFilePath.Insert(string(byteData))
	}

	// TODO can not > 4 because http request error Tiki: stopped after 10 redirects
	crawlRoutinePool := NewWorkerPool(2)

	go getProductDetailOnRestrictPage()

	// TODO handle category manually => bad
	for i := 0; i < len(categories); i++ {
		category := categories[i]
		productResp, err := httprequest.GetTikiProductList(1,
			configuration.GetTikiPageConfig().ProductAPIQueryParam.Limit, category.Code)
		if err != nil {
			logger.LogError("Error while call product API")
			continue
		}
		logger.LogDebug(category.Title, ": total product", productResp.Paging.Total)
		totalData += productResp.Paging.Total
		// concurrency
		crawlRoutinePool.Execute(func() {
			crawl.getProductDataByCategory(category, productResp.Paging.LastPage, doneChan)
		})
	}

	logger.LogInfo("Total product data expected to crawl: ", totalData)
}

/*
*
 */
func (crawl TikiCrawlHandler) getProductDataByCategory(category metadata.CategoryRoot, lastPage int, doneChan chan bool) {
	productFile := file.NewFileDataSource(configuration.GetFileConfig().PrefixName + category.Title + "-" + category.Code)
	productFileDetail := file.NewFileDataSource(configuration.GetFileConfig().PrefixName + category.Title + "-" + category.Code + "-Detail")

	for pageNum := 1; pageNum <= lastPage; pageNum++ {
		logger.LogInfo("@@@@@@@@@@@@@@@@@@@@@@@@@", category.Title, ": page Number ", pageNum, "@@@@@@@@@@@@@@@@@@@@@@@@@")
		productResp, err := httprequest.GetTikiProductList(pageNum,
			configuration.GetTikiPageConfig().ProductAPIQueryParam.Limit, category.Code)

		if err != nil {
			logger.LogError("Error while call product API", err)
			if pageNum == lastPage {
				productFile.Close()
				productFileDetail.Close()
			}
			continue
		}

		if len(productResp.Data) == 0 {
			logger.LogInfo("Product Data is empty", category.Title, ", At page", pageNum)
			if pageNum == lastPage {
				productFile.Close()
				productFileDetail.Close()
			}
			continue
		}
		// TODO here pass chan to another function
		productResLen := len(productResp.Data)
		for ind, product := range productResp.Data {
			var isEnd bool = false
			if pageNum == lastPage && ind == productResLen-1 {
				isEnd = true
			}
			productDataQueue <- TransferProductDetail{
				productFile,
				productFileDetail,
				&product,
				isEnd,
			}
		}
	}
	doneChan <- true
}

func getProductDetailOnRestrictPage() {
	emptyCount := 0
	var isLoop bool = true

	for isLoop {
		select {
		case productDetail, ok := <-productDataQueue:
			logger.LogInfo("Read product data from queue ", productDetail.product.Name)
			if ok {
				emptyCount = 0
				byteData, err := json.Marshal(productDetail.product)
				var i int = 0
				var errorCount int = 2
				var exponential float64 = 2.0 //2.7 + random.Float64()*(6.5-2.7)
				var product *metadata.Product = productDetail.product

				if err != nil {
					logger.LogError("Error while call product API", err)
					if productDetail.isEnd {
						productDetail.productFile.Close()
						productDetail.productFileDetail.Close()
					}
					continue
				}

				if product.UrlPath == "" || len(product.UrlPath) <= 0 {
					if productDetail.isEnd {
						productDetail.productFile.Close()
						productDetail.productFileDetail.Close()
					}
					continue
				}

				for {
					jsonProductDetailData, err := getProductDetailJson(product.UrlPath)
					if err != nil {
						jsonProductData := string(byteData)
						productDetail.productFile.Insert(jsonProductData)
						i++
						errorCount = 1
						break
					}

					if jsonProductDetailData != "" {
						atomic.AddInt32(&currentCountProduct, 1)
						jsonProductData := string(byteData)
						productDetail.productFile.Insert(jsonProductData)
						productDetail.productFileDetail.Insert(jsonProductDetailData)
						i++
						errorCount = 1
						break
					}
					// more than 129s 2^7, 7 time
					if errorCount > 129 {
						logger.LogInfo("We can't do request forever, errorCount = ", errorCount)
						jsonProductData := string(byteData)
						productDetail.productFile.Insert(jsonProductData)
						i++
						errorCount = 1
						break
					}

					logger.LogInfo("Start retry with duration ", time.Duration(errorCount)*time.Second)
					time.Sleep(time.Duration(errorCount) * time.Second)
					errorCount = int(math.Round(exponential * float64(errorCount)))
				}

				if productDetail.isEnd {
					productDetail.productFile.Close()
					productDetail.productFileDetail.Close()
				}
			}
		default:
			if emptyCount > 5 {
				isLoop = false
				logger.LogDebug("Channel is empty (or no value is ready to be read).")
			}
			time.Sleep(3 * time.Second)
			emptyCount++
		}
	}

}

func getProductDetailJson(page string) (string, error) {
	document, err := httprequest.GetTikiHtmlPage(page)
	if document == nil || document == (&goquery.Document{}) || err != nil {
		logger.LogError("Error when get Tiki product detail html page")
		return "", err
	}

	productDetailParser := tiki.ProductDetailParser{}
	var productDetail metadata.ProductDetail = productDetailParser.Parse(document)
	if productDetail.ProductId <= 0 {
		logger.LogDebug("Product Id is empty at page ", page)
		return "", nil
	}
	byteData, err := json.Marshal(productDetail)
	if err != nil {

		logger.LogError("Error while call product API", err)
	}
	return string(byteData), nil
}
