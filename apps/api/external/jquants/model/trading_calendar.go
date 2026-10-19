package model

type TradingCalendar struct {
	Date   string `json:"Date"`
	HolDiv string `json:"HolDiv"`
}

type ListedTradingCalendar struct {
	Data []TradingCalendar `json:"data"`
}
