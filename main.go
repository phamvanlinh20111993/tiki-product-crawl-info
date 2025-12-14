package main

import (
	"fmt"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/handle"
	"selfstudy/crawl/product/parser/tiki"
	"selfstudy/crawl/product/util"
	"strconv"
)

func test() {
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

	for pageNum := 1; pageNum < 2; pageNum++ {
		products, err := handle.GetTikiProducts(pageNum, configuration.GetPageConfig().TikiProductAPIQueryParam.Limit, "15078")
		if err != nil {
			return
		}
		fmt.Println("#########################" + strconv.Itoa(pageNum) + " ##########################################")
		fmt.Println("######################## ########################################################## #############")

		fmt.Println("product size ", len(products))
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
	}

	fmt.Println("###################################################################")

	fmt.Println(handle.GetTikiHtmlPage(""))
}

func test1() {
	document, _ := handle.GetTikiHtmlPage("")
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

func main() {
	// test()

	test1()
}
