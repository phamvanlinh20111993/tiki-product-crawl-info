package main

import (
	"fmt"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/http-request"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/parser/tiki"
	"selfstudy/crawl/product/util"
	"strconv"
	"testing"
)

func Test_Common(t *testing.T) {
	util.LogInfo("Hello world Michael Pham")

	util.LogError("Hello world Michael Pham")

	util.LogWarn("Hello world Michael Pham")

	util.LogDebug("Hello world Michael Pham")

	fmt.Println("###################################################################")

	ma := configuration.LoadConfiguration()
	for k, v := range ma {
		fmt.Println(k, v)
	}

	fmt.Println("###################################################################")

	fmt.Println("GetOpenSearchConfig ", configuration.GetOpenSearchConfig())

	fmt.Println("GetLoggerConfig ", configuration.GetLoggerConfig())

	fmt.Println("GetFileConfig ", configuration.GetFileConfig())

	fmt.Println("GetPostgresConfig ", configuration.GetPostgresConfig())

	fmt.Println("GetPageConfig ", configuration.GetPageConfig())

	fmt.Println("###################################################################")

	fmt.Println(http_request.GetTikiHtmlPage(""))
}

func Test_TikiProduct(t *testing.T) {
	for pageNum := 1; pageNum < 2; pageNum++ {
		products, err := http_request.GetTikiProductList(pageNum, configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, "15078")
		if err != nil {
			return
		}
		fmt.Println("#########################" + strconv.Itoa(pageNum) + " ##########################################")
		fmt.Println("######################## ########################################################## #############")

		fmt.Println("product size ", len(products.Data))
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
	}
}

func Test1(t *testing.T) {
	document, _ := http_request.GetTikiHtmlPage("")
	if document == nil {
		panic("Error when get tiki html page")
	}
	fmt.Println(document)

	categoryParser := tiki.CategoryParser{}
	var categories = categoryParser.Parse(document)

	for _, category := range categories {
		fmt.Println(category.Code + ", " + category.Title + ", " + category.Path + ", " + category.CategoryImagePresentation)
	}
}

func Test2_ParserProductDetail(t *testing.T) {

	var pageDetails []string = []string{
		"tuyen-tap-kiet-tac-fujiko-f-fujio-f-the-best-tang-kem-set-sticker-kho-lon-p278948895.html?spid=278948896",
		"den-hau-xe-dap-kiotool-2-mau-xanh-do-k02-p275389988.html?spid=275389989",
		"vali-keo-du-lich-cao-cap-bao-hanh-chinh-hang-size-24inch-ks-219-hong-nhat-p171574040.html?spid=171574041",
		"dong-ho-casio-nam-w-218h-1bv-chinh-hang-p271299443.html?itm_campaign=CTP_YPD_TKA_PLA_UNK_ALL_UNK_UNK_UNK_UNK_X.295663_Y.1877983_Z.3979643_CN.Đong-ho-nam&itm_medium=CPC&itm_source=tiki-ads&spid=271299444",
	}

	for _, page := range pageDetails {
		document, _ := http_request.GetTikiHtmlPage(page)
		if document == nil {
			panic("Error when get tiki html page")
		}

		productDetailParser := tiki.ProductDetailParser{}
		var productDetail metadata.ProductDetail = productDetailParser.Parse(document)
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$ %%%%%%%%%%%%%%%%%%%%%%%%%%%%%% ############################################################")
		fmt.Println(productDetail)
	}

}
