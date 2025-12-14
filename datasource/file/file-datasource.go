package file

import (
	"bufio"
	"io"
	"os"
	"selfstudy/crawl/product/configuration"
	"sync"
)

type FileDataSource struct {
	connection *os.File
}

var (
	once sync.Once
)

func NewFileDataSource(connection string) *FileDataSource {
	var fileOpen *os.File
	once.Do(func() {
		// open output file
		fo, err := os.Create(configuration.GetFileConfig().Path) // get from config
		if err != nil {
			panic(err)
		}
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		fileOpen = fo
	})

	return &FileDataSource{connection: fileOpen}
}

func (fd *FileDataSource) insert() {
	// make a write buffer
	w := bufio.NewWriter(fd.connection)
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)

	//TODO recheck
	r := bufio.NewReader(fd.connection)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}
