package configuration

type TikiPageConfig struct {
	Name                 string              `conf:"crawl.tiki-page.name"`
	BaseURL              string              `conf:"crawl.tiki-page.base-url"`
	ProductAPIURL        string              `conf:"crawl.tiki-page.product-api-url"`
	CategoryPathAPIRL    string              `conf:"crawl.tiki-page.category-path-api-url"`
	ProductAPIQueryParam TikiProductAPIQuery `conf:"crawl.tiki-page.product-api-query-param"`
}

type LazadaPageConfig struct {
	Name                 string              `conf:"crawl.lazada-page.name"`
	BaseURL              string              `conf:"crawl.lazada-page.base-url"`
	ProductAPIURL        string              `conf:"crawl.lazada-page.product-api-url"`
	ProductAPISearchURL  string              `conf:"crawl.lazada-page.product-api-search-url"`
	ProductAPIQueryParam TikiProductAPIQuery `conf:"crawl.lazada-page.product-api-query-param"`
}

type TikiProductAPIQuery struct {
	Limit int `conf:"crawl.tiki-page.product-api-query-param.limit"`
}

type OpenSearchConfig struct {
	Port     int    `conf:"crawl.datasource.opensearch.port"`
	URL      string `conf:"crawl.datasource.opensearch.url"`
	Username string `conf:"crawl.datasource.opensearch.username"`
	Password string `conf:"crawl.datasource.opensearch.password"`
}

type PostgresConfig struct {
	DatabaseURL  string `conf:"crawl.datasource.postgres.database-url"`
	Username     string `conf:"crawl.datasource.postgres.username"`
	Password     string `conf:"crawl.datasource.postgres.password"`
	Host         string `conf:"crawl.datasource.postgres.host"`
	Port         int    `conf:"crawl.datasource.postgres.port"`
	DatabaseName string `conf:"crawl.datasource.postgres.databaseName"`
}

type FileConfig struct {
	Path       string `conf:"crawl.datasource.file-local.path"`
	PrefixName string `conf:"crawl.datasource.file-local.prefix-name"`
	Extension  string `conf:"crawl.datasource.file-local.extension"`
}

type LoggerConfig struct {
	Level          string   `conf:"crawl.logger.level"`
	IsAddSource    bool     `conf:"crawl.logger.add-source"`
	IsTraceRequest bool     `conf:"crawl.logger.trace-request"`
	Target         []string `conf:"crawl.logger.target"`
	//	Nums           []float64 `conf:"crawl.logger.nums"` // not use, just for testing
	FilePath    string `conf:"crawl.logger.file-path"`
	FilePattern string `conf:"crawl.logger.file-pattern"`
	KeepLogDays int    `conf:"crawl.logger.keep-log-days"`
	//	Booleans       []bool    `conf:"crawl.logger.booleans"` // not use, just for testing
}

const DescriptionConfigStruct = "conf"
