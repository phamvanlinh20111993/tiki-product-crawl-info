package handle

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/util"
	"strconv"
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

func getTikiProducts(pageNum int, limit int, category string) ([]metadata.Product, error) {
	// client := getRestyClientInstance()
	resp, err := resty.New().R().
		EnableTrace().
		SetQueryParam("page", strconv.Itoa(pageNum)).
		SetQueryParam("limit", strconv.Itoa(limit)).
		SetQueryParam("category", category).
		Get("https://tiki.vn/api/personalish/v1/blocks/listings")

	if err != nil {
		util.LogError("request failed", slog.Any("error", err))
		return []metadata.Product{}, err
	}

	if resp.StatusCode() >= 400 {
		util.LogDebug("", slog.Any("bad status", resp.Status()))
		return []metadata.Product{}, err
	}

	bodyString := string(resp.Body())
	fmt.Println(bodyString)
	var tikiProductResponse metadata.Response
	err = json.Unmarshal(resp.Body(), &tikiProductResponse)
	if err != nil {
		util.LogError("Error", slog.Any("Error ", err))
	}

	return tikiProductResponse.Data, nil
}

var GetTikiProducts = getTikiProducts
