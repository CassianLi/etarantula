package models

// CategoryInfoRequest The request of category info
type CategoryInfoRequest struct {
	ProductNo    string `json:"asin"`
	SalesChannel string `json:"channel"`
	PriceNo      string `json:"priceNo"`
	Country      string `json:"country"`
	Price        string `json:"price"`
}

// CategoryInfo The information about category
type CategoryInfo struct {
	ProductNo    string   `json:"asin"`
	SalesChannel string   `json:"channel"`
	PriceNo      string   `json:"priceNo"`
	Country      string   `json:"country"`
	Price        string   `json:"price"`
	NewPrice     string   `json:"newPrice"`
	Screenshot   string   `json:"screenshot"`
	Status       string   `json:"status"`
	Errors       []string `json:"errors"`
}
