package configuration

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"selfstudy/crawl/product/util"
	"time"
)

const CONFIG_FILE_NAME = "crawl-config.yml"

var currentTime time.Time

func reload() {

}

func getConfiguration() {
	dir, err := os.Getwd()
	if err != nil {
		util.LogError("Error while get current working directory" + err.Error())
	}
	configPath := filepath.Join(dir, CONFIG_FILE_NAME)

	f, err := os.ReadFile(configPath)
	if err != nil {
		util.LogError("Error while read " + err.Error())
	}

	// get all the input data into an interface
	var config map[string]interface{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		util.LogError("error ", err)
	}

	fmt.Println(config)
}

var GetConfiguration = getConfiguration
