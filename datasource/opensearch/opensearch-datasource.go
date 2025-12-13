package opensearch

import (
	"crypto/tls"
	"github.com/opensearch-project/opensearch-go"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type OpenSearchDataSource struct {
	opensearchClient *opensearch.Client
}

var (
	opensearchClient OpenSearchDataSource
	once             sync.Once
)

func getOpensearchClient() OpenSearchDataSource {
	var opensearchClient = OpenSearchDataSource{}
	once.Do(func() {
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
			RetryOnStatus:        []int{502, 503, 504},
			Addresses:            []string{"https://localhost:9200"},
			//Username:  "admin", // For testing only. Don't store credentials in code.
			//Password:  "admin",
		})

		if err != nil {
			slog.Error("operation failed",
				slog.Any("error", err))
			panic(err)
		}

		opensearchClient.opensearchClient = client
	})

	return opensearchClient
}

func (o OpenSearchDataSource) insert() {

}

func (o OpenSearchDataSource) insertBatch() {

}

func (o OpenSearchDataSource) update() {

}

func (o OpenSearchDataSource) delete() {

}
