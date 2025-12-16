package http_request

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/util"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
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
		client := resty.NewWithClient(httpClient())
		//		client.SetDebug(true)
		client.SetCloseConnection(false)
		client.SetTimeout(30 * time.Second)
		client.SetRetryWaitTime(5 * time.Second)
		client.SetRetryMaxWaitTime(20 * time.Second)
		client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= http.StatusInternalServerError // Retry on server errors
			},
		)

		restyClient = client
	})
	return restyClient
}

func GetAPIData[T any](url string, requestParams map[string]map[string]string) (T, error) {
	client := getRestyClientInstance()

	var request *resty.Request = client.R()

	if len(requestParams["queryParams"]) > 0 {
		request = request.SetQueryParams(requestParams["queryParams"])
	}
	if len(requestParams["headers"]) > 0 {
		request = request.SetHeaders(requestParams["headers"])
	}
	if len(requestParams["token"]) > 0 && len(requestParams["token"]) < 2 {
		request = request.SetAuthToken(requestParams["token"]["0"])
	}

	resp, err := request.
		EnableTrace(). // => tracing request/response information
		Get(url)

	var dataResponse T
	if err != nil {
		util.LogError("request failed", slog.Any("error", err))
		return dataResponse, err
	}

	if configuration.GetLoggerConfig().IsTraceRequest {
		ti := resp.Request.TraceInfo()
		// Explore trace info
		util.LogDebug("Request Trace Info:")
		util.LogDebug("  DNSLookup     :", ti.DNSLookup)
		util.LogDebug("  ConnTime      :", ti.ConnTime)
		util.LogDebug("  TCPConnTime   :", ti.TCPConnTime)
		util.LogDebug("  TLSHandshake  :", ti.TLSHandshake)
		util.LogDebug("  ServerTime    :", ti.ServerTime)
		util.LogDebug("  ResponseTime  :", ti.ResponseTime)
		util.LogDebug("  TotalTime     :", ti.TotalTime)
		util.LogDebug("  IsConnReused  :", ti.IsConnReused)
		util.LogDebug("  IsConnWasIdle :", ti.IsConnWasIdle)
		util.LogDebug("  ConnIdleTime  :", ti.ConnIdleTime)
		util.LogDebug("  RequestAttempt:", ti.RequestAttempt)
		util.LogDebug("  RemoteAddr    :", ti.RemoteAddr.String())
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		util.LogDebug("", slog.Any("bad status: ", resp.Status()))
		return dataResponse, err
	}

	//	bodyString := string(resp.Body())
	//	fmt.Println(bodyString)
	err = json.Unmarshal(resp.Body(), &dataResponse)
	if err != nil {
		util.LogError("Error", slog.Any("Error ", err))
	}

	return dataResponse, nil
}

func getHTMLPage(url string, requestParams map[string]map[string]string) (*goquery.Document, error) {
	client := getRestyClientInstance()

	var request *resty.Request = client.R()

	if len(requestParams["queryParams"]) > 0 {
		request = request.SetQueryParams(requestParams["queryParams"])
	}
	if len(requestParams["headers"]) > 0 {
		request = request.SetHeaders(requestParams["headers"])
	}
	if len(requestParams["token"]) > 0 && len(requestParams["token"]) < 2 {
		request = request.SetAuthToken(requestParams["token"]["0"])
	}

	resp, err := request.
		EnableTrace(). // => tracing request/response information
		SetHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36").
		SetHeader("sec-ch-ua", "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"").
		SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Get(url)

	if err != nil {
		util.LogError("request failed", slog.Any("error", err))
		return &goquery.Document{}, err
	}

	if configuration.GetLoggerConfig().IsTraceRequest {
		ti := resp.Request.TraceInfo()
		// Explore trace info
		util.LogDebug("Request Trace Info:")
		util.LogDebug("  DNSLookup     :", ti.DNSLookup)
		util.LogDebug("  ConnTime      :", ti.ConnTime)
		util.LogDebug("  TCPConnTime   :", ti.TCPConnTime)
		util.LogDebug("  TLSHandshake  :", ti.TLSHandshake)
		util.LogDebug("  ServerTime    :", ti.ServerTime)
		util.LogDebug("  ResponseTime  :", ti.ResponseTime)
		util.LogDebug("  TotalTime     :", ti.TotalTime)
		util.LogDebug("  IsConnReused  :", ti.IsConnReused)
		util.LogDebug("  IsConnWasIdle :", ti.IsConnWasIdle)
		util.LogDebug("  ConnIdleTime  :", ti.ConnIdleTime)
		util.LogDebug("  RequestAttempt:", ti.RequestAttempt)
		util.LogDebug("  RemoteAddr    :", ti.RemoteAddr.String())
	}

	if resp.StatusCode() > http.StatusBadRequest {
		util.LogDebug("Get page false", slog.Any("bad status", resp.Status()))
		return &goquery.Document{}, err
	}

	var responseContentType string = resp.Header().Get("content-type")
	match, err := regexp.MatchString("text/html", responseContentType)
	if err != nil || !match {
		panic("The webpage should return an html page")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))

	if err != nil {
		util.LogError("Error while parsing html response: ", err)
	}

	return doc, nil
}

var GetHTMLPage = getHTMLPage
