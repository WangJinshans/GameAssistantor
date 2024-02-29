package model

type TradeSummary struct {
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	Rate    float64 `json:"rate"`
	Title   string  `json:"title"`
	Content string  `json:"cotent"`
}
