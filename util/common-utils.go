package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss_dot_zzz /*
const Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss_dot_zzz = "2006-01-02 15:04:05.000"

const Format_yyyy_mm_dd = "2006-01-02"

// Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss /*
const Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss = "2006-01-02 15:04:05"

func timeToString(time time.Time, pattern string) string {
	if pattern == "" {
		pattern = Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss
	}
	return time.Format(pattern)
}

func ConvertJsonData[T any](jsonData []byte, instanceType T) {
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.DisallowUnknownFields()
	// decoder.UseStrictFields()
	err := decoder.Decode(&instanceType)

	if err != nil {

		slog.Error("Error while convert JsonData to instanceType. Error: " + err.Error())
	}
}

// TODO we can use json.MarshalIndend(), suggest handle manually
func printStructuralData(data any) {
	if reflect.TypeOf(data).Kind() != reflect.Struct {
		slog.Error("Structural data is not a struct")
		return
	}

	value := reflect.ValueOf(data)
	numberOfFields := value.NumField()
	for i := 0; i < numberOfFields; i++ {
		fmt.Println("Field name: ", value.Type().Field(i).Name, ",Field value: ", value.Field(i).String())
	}
}

func getCurrentFolder() string {
	dir, err := os.Executable()
	if err != nil {
		slog.Error(err.Error())
	}
	// old way: https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
	//dir, err := os.Getwd()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	exePath, _ := filepath.EvalSymlinks(dir)
	exeDir := filepath.Dir(exePath)
	return exeDir
}

func CurrentTimeToString(pattern string) string {
	if pattern == "" {
		pattern = Format_yyyy_mm_dd_space_hh_dot_mm_dot_ss
	}
	return timeToString(time.Now(), pattern)
}

var TimeToString = timeToString
var PrintStructuralData = printStructuralData
var GetCurrentFolder = getCurrentFolder
