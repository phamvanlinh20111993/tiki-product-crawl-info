package main

import (
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/datasource/file"
	"selfstudy/crawl/product/util"
)

func main() {
	util.RemoveDir(configuration.GetFileConfig().Path)
	util.CreateDir(configuration.GetFileConfig().Path)
}

// TODO handle later
func createFileDataSource() file.FileDataSource {
	return file.FileDataSource{}
}
