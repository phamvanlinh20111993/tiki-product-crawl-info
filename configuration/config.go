package configuration

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = "crawl-config.yml"

var (
	lastConfigFileModifiedTime time.Time
	configMapKeyValue          map[string]any
	once                       sync.Once
	logger                     *slog.Logger
	onceLogger                 sync.Once
)

func getLogger() *slog.Logger {
	onceLogger.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: false}))
	})
	return logger
}

func reload() {
	for {
		time.Sleep(300 * time.Second)
		getLogger().Debug("^^ Start reload config file ^^")
		dir, err := os.Getwd()
		if err != nil {
			getLogger().Error("Error while get current working directory" + err.Error())
		}
		configPath := filepath.Join(dir, ConfigFileName)
		info, err := os.Stat(configPath)

		if err != nil {
			getLogger().Error("Error while get config file information" + err.Error())
			panic(errors.New("Error while get config file information" + err.Error()))
		}
		// if time.Now().UTC().Sub(currentTime) > 15*time.Minute {
		if info.ModTime().After(lastConfigFileModifiedTime) {
			lastConfigFileModifiedTime = time.Now().UTC()
			getConfigFromFile()
		}
	}
}

var lock sync.Mutex

func getConfigFromFile() {
	lock.Lock()
	dir, err := os.Getwd()
	if err != nil {
		getLogger().Error("Error while get current working directory" + err.Error())
	}
	configPath := filepath.Join(dir, ConfigFileName)

	f, err := os.ReadFile(configPath)
	if err != nil {
		getLogger().Error("Error while read " + err.Error())
	}

	// get all the input data into an interface
	var config map[string]interface{}
	if err := yaml.Unmarshal(f, &config); err != nil {
		getLogger().Error("error ", err)
	}

	configMapKeyValue = make(map[string]any)
	var flatMapConfig func(config map[string]interface{}, key string)

	flatMapConfig = func(config map[string]interface{}, key string) {
		for k, v := range config {
			var configMapKey string = key
			if key != "" {
				configMapKey += "."
			}

			if reflect.ValueOf(v).Kind() == reflect.Map {
				value, ok := v.(map[string]interface{})
				if ok {
					flatMapConfig(value, configMapKey+k)
				}
			} else {
				configMapKeyValue[configMapKey+k] = v
			}
		}
	}

	flatMapConfig(config, "")
	lock.Unlock()
}

func loadConfiguration() map[string]any {
	once.Do(func() {
		go reload()
	})

	if len(configMapKeyValue) < 1 {
		getConfigFromFile()
	}
	return configMapKeyValue
}

var LoadConfiguration = loadConfiguration
