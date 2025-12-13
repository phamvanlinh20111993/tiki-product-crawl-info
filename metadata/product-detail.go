package metadata

type ProductDetail struct {
	ProductId         int64 `json:"productId"`
	describeImage     []string
	description       string
	detailInformation map[string]string
}
