package http_request

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/logger"
	"strconv"
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
			MaxIdleConns:        16, // Max idle connections in pool
			MaxIdleConnsPerHost: 20, // Max idle connections per host
			IdleConnTimeout:     60 * time.Second,
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
		client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			// You can add logic here to modify the request just before it's executed
			logger.LogDebug("Preparing to send request to ", request.URL, ", Method ", request.Method)
			if request.QueryParam != nil {
				marshalIndent, err := json.MarshalIndent(request.QueryParam, "", "  ")
				if err == nil {
					logger.LogDebug("QueryParam ", string(marshalIndent))
				}
			}
			return nil
		})

		restyClient = client
	})
	return restyClient
}

const QueryParams = "queryParams"
const Headers = "headers"
const Token = "token"
const PathParams = "pathParams"

func GetAPIData[T any](url string, requestParams map[string]map[string]string) (T, error) {
	var request *resty.Request = getRestyClientInstance().R()

	if len(requestParams[QueryParams]) > 0 {
		request = request.SetQueryParams(requestParams[QueryParams])
	}
	if len(requestParams[Headers]) > 0 {
		request = request.SetHeaders(requestParams[Headers])
	}
	if len(requestParams[PathParams]) > 0 {
		request = request.SetPathParams(requestParams[PathParams])
	}
	if len(requestParams[Token]) > 0 && len(requestParams[Token]) < 2 {
		request = request.SetAuthToken(requestParams[Token]["0"])
	}

	resp, err := request.
		EnableTrace(). // => tracing request/response information
		Get(url)

	var dataResponse T
	if err != nil {
		logger.LogError("request failed", slog.Any("error", err))
		return dataResponse, err
	}

	if configuration.GetLoggerConfig().IsTraceRequest {
		ti := resp.Request.TraceInfo()
		// Explore trace info
		logger.LogDebug("Request Trace Info:")
		logger.LogDebug("  DNSLookup     :", ti.DNSLookup)
		logger.LogDebug("  ConnTime      :", ti.ConnTime)
		logger.LogDebug("  TCPConnTime   :", ti.TCPConnTime)
		logger.LogDebug("  TLSHandshake  :", ti.TLSHandshake)
		logger.LogDebug("  ServerTime    :", ti.ServerTime)
		logger.LogDebug("  ResponseTime  :", ti.ResponseTime)
		logger.LogDebug("  TotalTime     :", ti.TotalTime)
		logger.LogDebug("  IsConnReused  :", ti.IsConnReused)
		logger.LogDebug("  IsConnWasIdle :", ti.IsConnWasIdle)
		logger.LogDebug("  ConnIdleTime  :", ti.ConnIdleTime)
		logger.LogDebug("  RequestAttempt:", ti.RequestAttempt)
		logger.LogDebug("  RemoteAddr    :", ti.RemoteAddr.String())
	}

	if resp.StatusCode() >= http.StatusBadRequest {
		logger.LogDebug("Could not handle ", slog.Any("status: ", resp.Status()))
		return dataResponse, errors.New("could not handle >= bad status " + strconv.Itoa(resp.StatusCode()))
	}

	//	bodyString := string(resp.Body())
	//	fmt.Println(bodyString)
	err = json.Unmarshal(resp.Body(), &dataResponse)
	if err != nil {
		logger.LogError("Error", slog.Any("Error ", err))
	}

	return dataResponse, err
}

func getHTMLPage(url string, requestParams map[string]map[string]string) (*goquery.Document, error) {
	var request *resty.Request = getRestyClientInstance().R()

	if len(requestParams[QueryParams]) > 0 {
		request = request.SetQueryParams(requestParams[QueryParams])
	}
	if len(requestParams[PathParams]) > 0 {
		request = request.SetPathParams(requestParams[PathParams])
	}
	if len(requestParams[Headers]) > 0 {
		request = request.SetHeaders(requestParams[Headers])
	}
	if len(requestParams[Token]) > 0 && len(requestParams[Token]) < 2 {
		request = request.SetAuthToken(requestParams[Token]["0"])
	}

	resp, err := request.
		EnableTrace(). // => tracing request/response information
		Get(url)

	if err != nil {
		logger.LogError("request failed", slog.Any("error", err))
		return &goquery.Document{}, err
	}

	if configuration.GetLoggerConfig().IsTraceRequest {
		ti := resp.Request.TraceInfo()
		// Explore trace info
		logger.LogDebug("Request Trace Info:")
		logger.LogDebug("  DNSLookup     :", ti.DNSLookup)
		logger.LogDebug("  ConnTime      :", ti.ConnTime)
		logger.LogDebug("  TCPConnTime   :", ti.TCPConnTime)
		logger.LogDebug("  TLSHandshake  :", ti.TLSHandshake)
		logger.LogDebug("  ServerTime    :", ti.ServerTime)
		logger.LogDebug("  ResponseTime  :", ti.ResponseTime)
		logger.LogDebug("  TotalTime     :", ti.TotalTime)
		logger.LogDebug("  IsConnReused  :", ti.IsConnReused)
		logger.LogDebug("  IsConnWasIdle :", ti.IsConnWasIdle)
		logger.LogDebug("  ConnIdleTime  :", ti.ConnIdleTime)
		logger.LogDebug("  RequestAttempt:", ti.RequestAttempt)
		logger.LogDebug("  RemoteAddr    :", ti.RemoteAddr.String())
	}

	if resp.StatusCode() > http.StatusBadRequest {
		logger.LogDebug("Get page false", slog.Any("bad status", resp.Status()))
		return &goquery.Document{}, errors.New("could not handle >= bad status " + strconv.Itoa(resp.StatusCode()))
	}

	var responseContentType string = resp.Header().Get("content-type")
	match, err := regexp.MatchString("^text/html", responseContentType)
	if err != nil || !match {
		logger.LogError("The webpage should return an html page")
		return &goquery.Document{}, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))

	if err != nil {
		logger.LogError("Error while parsing html response: ", err)
	}

	return doc, err
}

var GetHTMLPage = getHTMLPage
