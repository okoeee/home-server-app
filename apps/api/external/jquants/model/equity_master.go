package model

// EquityMaster J-Quantsから取得する銘柄情報
type EquityMaster struct {
	// 銘柄情報の適用日
	Date                string `json:"Date"`
	Code                string `json:"Code"`
	CompanyName         string `json:"CoName"`
	Sector33Code        string `json:"S33"`
	Sector33Name        string `json:"S33Nm"`
	MarketCode          string `json:"Mkt"`
	MarketName          string `json:"MktNm"`
	ProductCategoryCode string `json:"ProdCat"`
}

// ListedEquityMaster J-Quantsのレスポンス全体
type ListedEquityMaster struct {
	Data          []EquityMaster `json:"data"`
	PaginationKey string         `json:"pagination_key"`
}
