package main

import (
	"fmt"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/util"
)

func test() {
	util.LogInfo("Hello world Michael Pham")

	util.LogError("Hello world Michael Pham")

	util.LogWarn("Hello world Michael Pham")

	util.LogDebug("Hello world Michael Pham")

	ma := configuration.LoadConfiguration()
	for k, v := range ma {
		fmt.Println(k, v)
	}

	fmt.Println("GetOpenSearchConfig ", configuration.GetOpenSearchConfig())

	fmt.Println("GetLoggerConfig ", configuration.GetLoggerConfig())

	fmt.Println("GetFileConfig ", configuration.GetFileConfig())

	fmt.Println("GetPostgresConfig ", configuration.GetPostgresConfig())

	fmt.Println("GetPageConfig ", configuration.GetPageConfig())

	//products, err := handle.GetTikiProducts(1, 1, "15078")
	//if err != nil {
	//	return
	//}
	//fmt.Println("######################## ##################################")
	//fmt.Println("%+v\n", products)
}

func main() {
	test()
}
