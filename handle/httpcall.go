package handle

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
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

var restyClient *resty.Client

func getRestyClientInstance() *resty.Client {
	once.Do(func() {
		restyClient := resty.NewWithClient(httpClient())
		restyClient.SetDebug(true)
		restyClient.SetCloseConnection(false)
	})
	return restyClient
}

func getApi() {
	client := getRestyClientInstance()
	defer client.

	res, err := client.R().
		EnableTrace().
		Get("https://httpbin.org/get")
	fmt.Println(err, res)
	fmt.Println(res.Request.TraceInfo())

}
