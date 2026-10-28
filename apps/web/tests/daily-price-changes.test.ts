import assert from 'node:assert/strict'
import { test } from 'node:test'
import { fetchLatestDailyPriceChanges } from '../src/api/daily-price-changes.ts'

test('1回の通信で最新取引日と全件を取得する', async (t) => {
  const changes = [
    { symbol: '12340', trade_date: '2026-10-19', close_price: '110.5', previous_trade_date: '2026-10-16', previous_close_price: '100.0', change_rate: '10.5' },
    { symbol: '56780', trade_date: '2026-10-19', close_price: '90.0', previous_trade_date: '2026-10-15', previous_close_price: '100.0', change_rate: '-10' },
  ]
  const result = { trade_date: '2026-10-19', changes }
  const fetchMock = t.mock.method(globalThis, 'fetch', async (url) => {
    assert.equal(url, '/api/daily-price-changes/latest')
    return Response.json(result)
  })
  assert.deepEqual(await fetchLatestDailyPriceChanges(), result)
  assert.equal(fetchMock.mock.callCount(), 1)
})

test('DBが空なら対象日なしと空配列を返す', async (t) => {
  const result = { trade_date: null, changes: [] }
  t.mock.method(globalThis, 'fetch', async () => Response.json(result))
  assert.deepEqual(await fetchLatestDailyPriceChanges(), result)
})

test('比較可能な銘柄がなくても最新取引日を返す', async (t) => {
  const result = { trade_date: '2026-10-19', changes: [] }
  t.mock.method(globalThis, 'fetch', async () => Response.json(result))
  assert.deepEqual(await fetchLatestDailyPriceChanges(), result)
})

test('HTTPエラーでは成功結果を返さない', async (t) => {
  t.mock.method(globalThis, 'fetch', async () => new Response('internal error', { status: 500 }))
  await assert.rejects(fetchLatestDailyPriceChanges(), /HTTP 500/)
})

test('通信失敗を呼び出し元へ返す', async (t) => {
  const failure = new TypeError('Failed to fetch')
  t.mock.method(globalThis, 'fetch', async () => { throw failure })
  await assert.rejects(fetchLatestDailyPriceChanges(), (error) => error === failure)
})
