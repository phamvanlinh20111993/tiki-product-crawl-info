package handle

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
	"time"
)

var (
	restyClient *resty.Client
	once        sync.Once
)

func httpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100, // Max idle connections in pool
			MaxIdleConnsPerHost: 20,  // Max idle connections per host
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
		},
		Timeout: 30 * time.Second,
	}

	return client
}

func getRestyClientInstance() *resty.Client {
	once.Do(func() {
		restyClient := resty.NewWithClient(httpClient())
		restyClient.SetDebug(true)
		restyClient.SetCloseConnection(false)
		restyClient.SetTimeout(30 * time.Second)
		restyClient.SetRetryWaitTime(5 * time.Second)
		restyClient.SetRetryMaxWaitTime(20 * time.Second)
		restyClient.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= 500 // Retry on server errors
			},
		)
	})
	return restyClient
}

func getApi() {
	client := getRestyClientInstance()
	resp, err := client.R().
		EnableTrace().
		Get("https://httpbin.org/get")

	if err != nil {
		fmt.Println("request failed: %w", err)
	}
	if resp.StatusCode() >= 400 {
		fmt.Println("bad status: %s", resp.Status())
	}
}
