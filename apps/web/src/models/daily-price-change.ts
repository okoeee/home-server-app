import type { DailyPriceChangeResponse } from '../api/generated/schema'

const dailyPriceChangeDisplayLimit = 10

const changeRateFormatter = new Intl.NumberFormat('ja-JP', {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
  useGrouping: false,
})

export function formatChangeRate(rate: string): string {
  return changeRateFormatter.format(Number(rate))
}

export function selectTopDailyPriceChanges(
  changes: DailyPriceChangeResponse[],
): DailyPriceChangeResponse[] {
  return [...changes]
    .sort((left, right) => Number(right.change_rate) - Number(left.change_rate))
    .slice(0, dailyPriceChangeDisplayLimit)
}
