package util

import (
	"bytes"
	"encoding/json"
	"time"
)

/*
*
format to yyyy-mm-dd hh:mm:ss
*/
func timeToString(time time.Time, pattern string) string {
	if pattern == "" {
		pattern = "2006-01-02 15:04:05"
	}
	return time.Format(pattern)
}

func ConvertJsonData[T any](jsonData []byte, instanceType T) {
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.DisallowUnknownFields()
	// decoder.UseStrictFields()
	err := decoder.Decode(&instanceType)

	if err != nil {
		logError("Error while convert JsonData to instanceType. Error: " + err.Error())
	}
}

var TimeToString = timeToString
