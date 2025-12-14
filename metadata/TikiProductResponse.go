package metadata

type Response struct {
	Block Block     `json:"block"`
	Data  []Product `json:"data"`
}

type Block struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
}
