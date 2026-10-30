import type { LatestDailyPriceChangesResponse } from './generated/schema'

export async function fetchLatestDailyPriceChanges(): Promise<LatestDailyPriceChangesResponse> {
  const response = await fetch('/api/v1/daily-price-changes/latest')

  if (!response.ok) {
    throw new Error(`騰落率の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
