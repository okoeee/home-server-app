package model

import "time"

type Stock struct {
	Symbol              string
	EffectiveDate       time.Time
	CompanyName         string
	Sector33Code        string
	Sector33Name        string
	MarketCode          string
	MarketName          string
	ProductCategoryCode string
}
