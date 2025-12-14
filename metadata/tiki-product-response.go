package metadata

type Response struct {
	Block  Block     `json:"block"`
	Data   []Product `json:"data"`
	Paging Paging    `json:"paging"`
}

type Paging struct {
	CurrentPage int `json:"current_page"`
	From        int `json:"from"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	To          int `json:"to"`
	Total       int `json:"total"`
}

type Block struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
}
