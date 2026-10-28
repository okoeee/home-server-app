import type { DailyPriceChange } from '../models/daily-price-change'

export async function fetchDailyPriceChanges(date: string): Promise<DailyPriceChange[]> {
  const query = new URLSearchParams({ date })
  const response = await fetch(`/api/daily-price-changes?${query}`)

  if (!response.ok) {
    throw new Error(`騰落率の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
