package main

import (
	"fmt"
	"log"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/handle"
	"selfstudy/crawl/product/util"
)

func main() {
	fmt.Println("Hello World")
	log.Println("standard logger")
	util.LogInfo("Hello world Michael Pham")

	configuration.GetConfiguration()

	products, err := handle.GetTikiProducts(1, 300, "15078")
	if err != nil {
		return
	}
	fmt.Println("######################## ##################################")
	fmt.Println(products)
}
