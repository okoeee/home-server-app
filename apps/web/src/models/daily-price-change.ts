import type { DailyPriceChangeResponse } from '../api/generated/schema'

const dailyPriceChangeDisplayLimit = 10

export function selectTopDailyPriceChanges(changes: DailyPriceChangeResponse[]): DailyPriceChangeResponse[] {
  return [...changes]
    .sort((left, right) => Number(right.change_rate) - Number(left.change_rate))
    .slice(0, dailyPriceChangeDisplayLimit)
}
