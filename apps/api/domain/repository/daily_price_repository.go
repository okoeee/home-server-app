package repository

import (
	"api/domain/model"
	"context"
	"time"
)

type DailyPriceRepository interface {
	Upsert(context.Context, model.DailyPrice) error
	ListVolumesByTradeDate(ctx context.Context, tradeDate time.Time) ([]model.DailyVolume, error)
	// ListVolumeHistory ListVolumeHistoryはbeforeDateを含めず、銘柄ごとの直近最大20件を返す
	ListVolumeHistory(ctx context.Context, beforeDate time.Time) ([]model.DailyVolume, error)
}
