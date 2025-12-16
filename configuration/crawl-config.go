package configuration

import (
	"reflect"
)

type PageConfig struct {
	TikiName                 string              `conf:"crawl.page.tiki-name"`
	TikiBaseURL              string              `conf:"crawl.page.tiki-base-url"`
	TikiProductAPIURL        string              `conf:"crawl.page.tiki-product-api-url"`
	TikiProductAPIQueryParam TikiProductAPIQuery `conf:"crawl.page.tiki-product-api-query-param"`
}

type TikiProductAPIQuery struct {
	Limit int `conf:"crawl.page.tiki-product-api-query-param.limit"`
}

type OpenSearchConfig struct {
	Port     int    `conf:"crawl.opensearch.port"`
	URL      string `conf:"crawl.opensearch.url"`
	Username string `conf:"crawl.opensearch.username"`
	Password string `conf:"crawl.opensearch.password"`
}

type PostgresConfig struct {
	DatabaseURL  string `conf:"crawl.postgres.database-url"`
	Username     string `conf:"crawl.postgres.username"`
	Password     string `conf:"crawl.postgres.password"`
	Host         string `conf:"crawl.postgres.host"`
	Port         int    `conf:"crawl.postgres.port"`
	DatabaseName string `conf:"crawl.postgres.databaseName"`
}

type FileConfig struct {
	Path string `conf:"crawl.file.path"`
	Name string `conf:"crawl.file.name"`
}

type LoggerConfig struct {
	Level          string `conf:"crawl.logger.level"`
	IsAddSource    bool   `conf:"crawl.logger.isAddSource"`
	IsTraceRequest bool   `conf:"crawl.logger.isTraceRequest"`
}

const DescriptionConfigStruct = "conf"

/*
*
TODO need more field type to be handle
refer: chatgpt support: https://chatgpt.com/
*/
func setFieldValue(field reflect.Value, value interface{}) {
	switch field.Kind() {

	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		}

	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case float64:
			field.SetInt(int64(v))
		}

	case reflect.Float32, reflect.Float64:
		if v, ok := value.(float64); ok {
			field.SetFloat(v)
		}

	default:
		getLogger().Info("Unsupported type:", field.Kind())
	}
}

func getStructTypeConfig(target interface{}) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	if t.Kind() != reflect.Ptr {
		getLogger().Error("Must be a pointer to a struct type:", t)
	}

	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		getLogger().Error("Provided type is not a struct")
	}
	configuration := loadConfiguration()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			getLogger().Error("Can't set field ", t.Field(i).Name)
			continue
		}
		confKey := fieldType.Tag.Get(DescriptionConfigStruct)

		// Handle nested struct recursively
		if fieldType.Type.Kind() == reflect.Struct {
			nested := fieldValue.Addr().Interface()
			getStructTypeConfig(nested)
			continue
		}

		if confKey == "" {
			getLogger().Debug("confKey is empty: ", confKey)
			continue
		}
		rawValue, ok := configuration[confKey]
		if !ok {
			getLogger().Debug("value for confKey is empty: ", confKey)
			continue
		}
		setFieldValue(fieldValue, rawValue)
	}
}

func getPageConfig() PageConfig {
	pageConfig := &PageConfig{}
	getStructTypeConfig(pageConfig)
	return *pageConfig
}

func getOpenSearchConfig() OpenSearchConfig {
	openSearchConfig := &OpenSearchConfig{}
	getStructTypeConfig(openSearchConfig)
	return *openSearchConfig
}

func getLoggerConfig() LoggerConfig {
	loggerConfig := &LoggerConfig{}
	getStructTypeConfig(loggerConfig)
	return *loggerConfig
}

func getPostgresConfig() PostgresConfig {
	postgresConfig := &PostgresConfig{}
	getStructTypeConfig(postgresConfig)
	return *postgresConfig
}

func getFileConfig() FileConfig {
	fileConfig := &FileConfig{}
	getStructTypeConfig(fileConfig)
	return *fileConfig
}

var GetPageConfig = getPageConfig
var GetOpenSearchConfig = getOpenSearchConfig
var GetLoggerConfig = getLoggerConfig
var GetPostgresConfig = getPostgresConfig
var GetFileConfig = getFileConfig
