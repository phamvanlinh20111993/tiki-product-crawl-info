package opensearch

import (
	"crypto/tls"
	"github.com/opensearch-project/opensearch-go"
	"log/slog"
	"net"
	"net/http"
	"selfstudy/crawl/product/configuration"
	"strconv"
	"sync"
	"time"
)

type OpenSearchDataSource struct {
	opensearchClient *opensearch.Client
}

var (
	once             sync.Once
	opensearchClient *opensearch.Client
)

func getOpensearchClient() OpenSearchDataSource {
	once.Do(func() {
		var openSearchConfig = configuration.GetOpenSearchConfig()
		client, err := opensearch.NewClient(opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:          100,              // Max total idle connections
				MaxIdleConnsPerHost:   10,               // Max idle connections per host
				MaxConnsPerHost:       0,                // Unlimited active connections per host (0 means no limit)
				IdleConnTimeout:       60 * time.Second, // How long an idle connection is kept alive
				ResponseHeaderTimeout: time.Second,      // Timeout for reading response headers
				DialContext: (&net.Dialer{
					Timeout: time.Second, // Connection timeout
				}).DialContext,
			},
			EnableRetryOnTimeout: true,
			EnableDebugLogger:    true, // get from config
			MaxRetries:           5,
			RetryOnStatus:        []int{http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}, // 502,503,504
			Addresses:            []string{openSearchConfig.URL + ":" + strconv.Itoa(openSearchConfig.Port)},
			//Username:  openSearchConfig.Username, // For testing only. Don't store credentials in code.
			//Password:  openSearchConfig.Password,
		})

		if err != nil {
			slog.Error("operation failed", slog.Any("error", err))
			panic(err)
		}

		opensearchClient = client
	})

	return OpenSearchDataSource{opensearchClient}
}

func getOpenSearchDataSourceInstance() OpenSearchDataSource {
	return getOpensearchClient()
}

var GetOpenSearchDataSourceInstance = getOpenSearchDataSourceInstance

func (o OpenSearchDataSource) insert() {

}

func (o OpenSearchDataSource) insertBatch() {

}

func (o OpenSearchDataSource) update() {

}

func (o OpenSearchDataSource) delete() {

}

func (o OpenSearchDataSource) close() {
	o.close()
}
