package http_request

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"selfstudy/crawl/product/configuration"
	"selfstudy/crawl/product/metadata"
	"selfstudy/crawl/product/util"
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
		Get(configuration.GetPageConfig().TikiProductAPIURL)

	if err != nil {
		util.LogError("request failed", slog.Any("error", err))
		return metadata.Response{}, err
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
		util.LogDebug("", slog.Any("bad status", resp.Status()))
		return metadata.Response{}, err
	}

	//	bodyString := string(resp.Body())
	//	fmt.Println(bodyString)
	var tikiProductListResponse metadata.Response
	err = json.Unmarshal(resp.Body(), &tikiProductListResponse)
	if err != nil {
		util.LogError("Error", slog.Any("Error ", err))
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
		Get(configuration.GetPageConfig().TikiBaseURL + path)

	if err != nil {
		util.LogError("request failed", slog.Any("error", err),
			slog.Any("url", configuration.GetPageConfig().TikiBaseURL+path))
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
	match, err := regexp.MatchString("^text/html", responseContentType)
	if err != nil || !match {
		panic("The webpage should return an html page")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))

	if err != nil {
		util.LogError("Error while parsing html response: ", err)
	}

	return doc, nil
}

// GetTikiProductList GetTikiProducts upper the first letter make it public outside of package
var GetTikiProductList = getTikiProductList
var GetTikiHtmlPage = getTikiHTMLPage
