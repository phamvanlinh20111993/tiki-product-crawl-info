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
		Get(configuration.GetPageConfig().TikiBaseURL + path)

	if err != nil {
		util.LogError("request failed", slog.Any("error", err))
		return &goquery.Document{}, err
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
