package models

// CategoryInfoRequest The request of category info
type CategoryInfoRequest struct {
	ProductNo    string `json:"asin"`
	SalesChannel string `json:"salesChannel"`
	PriceNo      string `json:"priceNo"`
	Country      string `json:"country"`
	Price        string `json:"price"`
}

// CategoryInfo The information about category
type CategoryInfo struct {
	ProductNo    string   `json:"asin"`
	SalesChannel string   `json:"salesChannel"`
	PriceNo      string   `json:"priceNo"`
	Country      string   `json:"country"`
	Price        string   `json:"price"`
	NewPrice     string   `json:"newPrice"`
	Screenshot   string   `json:"screenshot"`
	Errors       []string `json:"errors"`
}

type PriceSelector struct {
	Whole    string `mapstracture:"whole"`
	Fraction string `mapstracture:"fraction"`
	Tag      string `mapstracture:"tag"`
}

type DetailSelector struct {
	Section string `mapstracture:"section"`
	Tr      string `mapstracture:"tr"`
	Tag     string `mapstracture:"tag"`
}

type AmazonConfig struct {
	Urls            map[string]string `mapstracture:"urls"`
	PriceSelectors  []PriceSelector   `mapstracture:"price-selectors"`
	DetailSelectors []DetailSelector  `mapstracture:"detail-selectors"`
}
