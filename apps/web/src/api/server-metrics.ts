import type { ServerMetricsResponse } from './generated/schema'

export async function fetchServerMetrics(signal?: AbortSignal): Promise<ServerMetricsResponse> {
  const response = await fetch('/api/v1/server/metrics', { signal })

  if (!response.ok) {
    throw new Error(`サーバー情報の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
