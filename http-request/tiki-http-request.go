package http_request

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/logger"
	"selfstudy/crawl/product/metadata"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func getTikiProductList(pageNum int, limit int, category string) (metadata.Response, error) {
	client := getRestyClientInstance()
	resp, err := client.R().
		EnableTrace(). // => tracing request/response information
		SetQueryParam("page", strconv.Itoa(pageNum)).
		SetQueryParam("limit", strconv.Itoa(limit)).
		SetQueryParam("category", category).
		Get(configuration.GetPageConfig().ProductAPIURL)

	if err != nil {
		logger.LogError("request failed", slog.Any("error", err))
		return metadata.Response{}, err
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
		logger.LogDebug("", slog.Any("bad status", resp.Status()))
		return metadata.Response{}, err
	}

	//	bodyString := string(resp.Body())
	//	fmt.Println(bodyString)
	var tikiProductListResponse metadata.Response
	err = json.Unmarshal(resp.Body(), &tikiProductListResponse)
	if err != nil {
		logger.LogError("Error", slog.Any("Error ", err))
	}

	return tikiProductListResponse, nil
}

func getTikiHTMLPage(path string) (*goquery.Document, error) {
	client := getRestyClientInstance()
	resp, err := client.R().
		EnableTrace(). // => tracing request/response information
		//	SetHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36").
		//		SetHeader("sec-ch-ua", "\"Google Chrome\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"").
		//	SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Get(configuration.GetPageConfig().BaseURL + path)

	if err != nil {
		logger.LogError("request failed", slog.Any("error", err),
			slog.Any("url", configuration.GetPageConfig().BaseURL+path))
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
		return &goquery.Document{}, err
	}

	var responseContentType string = resp.Header().Get("content-type")
	match, err := regexp.MatchString("^text/html", responseContentType)
	if err != nil || !match {
		panic("The webpage should return an html page")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))

	if err != nil {
		logger.LogError("Error while parsing html response: ", err)
	}

	return doc, nil
}

// GetTikiProductList GetTikiProducts upper the first letter make it public outside of package
var GetTikiProductList = getTikiProductList
var GetTikiHtmlPage = getTikiHTMLPage
