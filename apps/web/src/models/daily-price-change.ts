export type DailyPriceChange = {
  symbol: string
  trade_date: string
  close_price: string
  previous_trade_date: string
  previous_close_price: string
  change_rate: string
}

export type LatestDailyPriceChanges = {
  trade_date: string | null
  changes: DailyPriceChange[]
}

const dailyPriceChangeDisplayLimit = 10

export function selectTopDailyPriceChanges(changes: DailyPriceChange[]): DailyPriceChange[] {
  return [...changes]
    .sort((left, right) => Number(right.change_rate) - Number(left.change_rate))
    .slice(0, dailyPriceChangeDisplayLimit)
}
