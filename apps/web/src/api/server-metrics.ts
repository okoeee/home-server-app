import type { ServerMetricsResponse } from './generated/schema'

export async function fetchServerMetrics(): Promise<ServerMetricsResponse> {
  const response = await fetch('/api/v1/server/metrics')

  if (!response.ok) {
    throw new Error(`サーバー情報の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
