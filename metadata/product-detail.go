package metadata

type ProductDetail struct {
	ProductId         int64 `json:"productId"`
	DescribeImage     []string
	Description       string
	DetailInformation map[string]string
}
