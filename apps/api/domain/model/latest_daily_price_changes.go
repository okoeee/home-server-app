package model

import "time"

type LatestDailyPriceChanges struct {
	TradeDate *time.Time
	Changes   []DailyPriceChange
}
