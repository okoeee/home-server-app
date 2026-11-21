import assert from 'node:assert/strict'
import { test } from 'node:test'
import { fetchServerMetrics } from '../src/api/server-metrics.ts'

test('サーバー情報APIへ1回通信し、0%を含む使用率を受け取る', async (t) => {
  // Arrange
  const metrics = {
    fetched_at: '2026-11-09T09:00:00+09:00',
    cpu_usage_percent: 0,
    memory_usage_percent: 43.2,
  }
  const fetchMock = t.mock.method(globalThis, 'fetch', async (url) => {
    assert.equal(url, '/api/v1/server/metrics')
    return Response.json(metrics)
  })

  // Act
  const result = await fetchServerMetrics()

  // Assert
  assert.deepEqual(result, metrics)
  assert.equal(fetchMock.mock.callCount(), 1)
})

for (const status of [500, 502, 503]) {
  test(`サーバー情報APIのHTTP ${status}をエラーとして返す`, async (t) => {
    // Arrange
    t.mock.method(globalThis, 'fetch', async () => new Response('private details', { status }))

    // Act
    const result = fetchServerMetrics()

    // Assert
    await assert.rejects(result, {
      message: `サーバー情報の取得に失敗しました（HTTP ${status}）。`,
    })
  })
}

test('サーバー情報APIへの通信失敗を呼び出し元へ返す', async (t) => {
  // Arrange
  const failure = new TypeError('Failed to fetch')
  t.mock.method(globalThis, 'fetch', async () => {
    throw failure
  })

  // Act
  const result = fetchServerMetrics()

  // Assert
  await assert.rejects(result, (error) => error === failure)
})

test('サーバー情報の通信にキャンセル用のsignalを渡す', async (t) => {
  // Arrange
  const controller = new AbortController()
  const fetchMock = t.mock.method(globalThis, 'fetch', async (_url, _options) => Response.json({}))

  // Act
  await fetchServerMetrics(controller.signal)

  // Assert
  assert.equal(fetchMock.mock.calls[0].arguments[1]?.signal, controller.signal)
})
