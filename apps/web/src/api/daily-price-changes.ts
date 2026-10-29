import type { LatestDailyPriceChanges } from '../models/daily-price-change'

export async function fetchLatestDailyPriceChanges(): Promise<LatestDailyPriceChanges> {
  const response = await fetch('/api/v1/daily-price-changes/latest')

  if (!response.ok) {
    throw new Error(`騰落率の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
