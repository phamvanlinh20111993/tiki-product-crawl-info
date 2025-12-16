package util

import (
	"errors"
	"os"
)

// refer: https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func isExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		// path/to/whatever exists
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

func createDir(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false
	}
	return true
}

func removeDir(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		return false
	}
	return true
}

func getLineSeperator() string {
	newline := "\n"
	if os.PathSeparator == '\\' { // window os
		newline = "\r\n"
	}
	return newline
}

var IsExist = isExist
var CreateDir = createDir
var GetLineSeperator = getLineSeperator
var RemoveDir = removeDir
