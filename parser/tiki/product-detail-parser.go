package tiki

import (
	"encoding/json"
	"regexp"
	"selfstudy/crawl/product/logger"
	"selfstudy/crawl/product/metadata"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ProductDetailParser struct{}

const mapKeySplitMark string = ":"

func (pd ProductDetailParser) Parse(document *goquery.Document) metadata.ProductDetail {
	productDetail := metadata.ProductDetail{}
	// fmt.Println(document.Html())

	// get image list
	var imageList []string
	document.Find("a[data-view-id='pdp_main_view_photo'] picture img").Each(func(_ int, a *goquery.Selection) {
		imgLink, _ := a.Attr("srcset")
		linkSplit := strings.Split(imgLink, " ")
		splLen := len(linkSplit)
		if splLen > 1 {
			imageList = append(imageList, linkSplit[0])
		}
	})
	productDetail.DescribeImage = imageList

	// get product id
	productLinkUrl, isOk := document.Find("link[rel='canonical']").First().Attr("href")
	if isOk {
		r := regexp.MustCompile("p(\\d+).html\\s*$")
		matches := r.FindStringSubmatch(productLinkUrl)
		// Check if a match was found
		if len(matches) > 1 {
			value, err := strconv.Atoi(matches[1])
			if err == nil {
				productDetail.ProductId = int64(value)
			}
		}
	}

	/*
		const productDescription = "Mô tả sản phẩm"
		const detailInformation = "Thông tin chi tiết"
		var productDescriptionSelection *goquery.Selection
		var detailInformationSelection *goquery.Selection
		document.Find("div").Each(func(_ int, div *goquery.Selection) {
			if div.Text() == productDescription {
				productDescriptionSelection = div.Siblings().First().Find("span")
				return
			}
			if div.Text() == detailInformation {
				detailInformationSelection = div.Siblings().First()
				return
			}
			//	if productDescriptionSelection != nil && detailInformationSelection != nil {
			//		return
			//	}
		})

		// product description
		if productDescriptionSelection != nil {
			var index int64 = 0
			var key string = ""
			var productDescriptionMap = make(map[string]string)
			productDescriptionSelection.Each(func(_ int, span *goquery.Selection) {
				if index%2 == 0 {
					key = span.Text()
				} else {
					productDescriptionMap[key] = span.Text()
				}
				index++
			})
			productDetail.Description = productDescriptionMap[key]
		}

		//product detail information
		if detailInformationSelection != nil {
			productDetail.DescriptionHTML, _ = detailInformationSelection.Html()
			productDetail.Description = detailInformationSelection.Text()
		}
	*/

	var jsonMap map[string]any
	document.Find("script[id='__NEXT_DATA__'][type='application/json']").Each(func(index int, element *goquery.Selection) {
		// Get the content of the <script> tag
		jsonData := element.Text()
		err := json.Unmarshal([]byte(jsonData), &jsonMap)
		if err != nil {
			logger.LogError("Error unmarshalling JSON:", err)
			return
		}
	})

	if len(jsonMap) > 0 {
		var productDataMap map[string]any = foundKeyInMap(jsonMap, "productData"+mapKeySplitMark+"response"+mapKeySplitMark+"data")
		if len(productDataMap) > 0 {
			// product description
			productDetail.Description = productDataMap["description"].(string)

			//product detail information
			var specificationArr []interface{} = productDataMap["specifications"].([]interface{})
			if len(specificationArr) > 0 {
				for _, specification := range specificationArr {
					attributes, isOk := specification.(map[string]interface{})
					if isOk && attributes["attributes"] != nil {
						attributeArr, isCorrectVal := attributes["attributes"].([]interface{})
						if isCorrectVal {
							var detailInformation map[string]string = map[string]string{}
							for _, attribute := range attributeArr {
								mapKeyValue, ok := attribute.(map[string]interface{})
								if ok {
									detailInformation[mapKeyValue["name"].(string)] = mapKeyValue["value"].(string)
								}
							}
							productDetail.DetailInformation = detailInformation
						}

						break
					}
				}
			}
		}

	}

	return productDetail
}

func foundKeyInMap(inp map[string]any, keyMaps string) map[string]any {

	var isFind bool = false
	var recursive func(input map[string]any, key string)
	var result map[string]any

	recursive = func(input map[string]any, key string) {
		if !isFind {
			for k, v := range input {
				mapStringAny, isOke := v.(map[string]any)
				if isOke {
					if k == key {
						isFind = true
						result = mapStringAny
						break
					}
					recursive(mapStringAny, key)
				}
			}
		}
	}

	splitMapKey := strings.Split(keyMaps, mapKeySplitMark)

	var temp map[string]any = inp
	for _, keyMap := range splitMapKey {
		isFind = false
		result = make(map[string]any)
		recursive(temp, keyMap)
		if len(result) > 0 {
			temp = result
		} else {
			temp = nil
			logger.LogDebug("Can not found data for key ", keyMap)
			break
		}
	}

	return temp
}
