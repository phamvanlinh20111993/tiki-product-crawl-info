package main

import (
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/handle"
	"selfstudy/crawl/product/util"
)

func main() {
	util.RemoveDir(configuration.GetFileConfig().Path)
	util.CreateDir(configuration.GetFileConfig().Path)
	handle.CrawlHandle()
}
