package file

import (
	"os"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/logger"
	"selfstudy/crawl/product/util"
)

type FileDataSource struct {
	connection   *os.File
	fullFilePath string
}

const DOT string = "."

func NewFileDataSource(fileName string) *FileDataSource {
	var fileOpen *os.File
	var fullFilePath string = configuration.GetFileConfig().Path + string(os.PathSeparator) + fileName + DOT + configuration.GetFileConfig().Extension

	if util.IsExist(fullFilePath) {
		fOpen, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			logger.LogError("Error while open file: ", fileName, err)
			panic(err)
		}
		fileOpen = fOpen
	} else {
		// open output file
		fo, err := os.Create(fullFilePath) // get from config
		if err != nil {
			panic(err)
		}
		fileOpen = fo
	}

	return &FileDataSource{connection: fileOpen, fullFilePath: fullFilePath}
}

func (fd *FileDataSource) Insert(data string) {
	// fd.connection.WriteString(data + util.GetLineSeparator()) ()
	if _, err := fd.connection.Write([]byte(data + util.GetLineSeparator())); err != nil {
		logger.LogError("Error while writing to file: ", fd.fullFilePath)
	}
}

func (fd *FileDataSource) InsertBatch(listData []string) {
	for _, data := range listData {
		fd.Insert(data)
	}
}

func (fd *FileDataSource) Close() {
	err := fd.connection.Close()
	if err != nil {
		logger.LogError("Error while closing file: ", fd.fullFilePath, err)
	}
}
