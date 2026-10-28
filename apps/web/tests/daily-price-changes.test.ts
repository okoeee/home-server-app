import assert from 'node:assert/strict'
import { test } from 'node:test'
import { fetchDailyPriceChanges } from '../src/api/daily-price-changes.ts'

test('指定日の全件を文字列の数値のまま取得する', async (t) => {
  const changes = [
    { symbol: '12340', trade_date: '2026-10-19', close_price: '110.5', previous_trade_date: '2026-10-16', previous_close_price: '100.0', change_rate: '10.5' },
    { symbol: '56780', trade_date: '2026-10-19', close_price: '90.0', previous_trade_date: '2026-10-15', previous_close_price: '100.0', change_rate: '-10' },
  ]
  t.mock.method(globalThis, 'fetch', async (url) => {
    assert.equal(url, '/api/daily-price-changes?date=2026-10-19')
    return Response.json(changes)
  })
  assert.deepEqual(await fetchDailyPriceChanges('2026-10-19'), changes)
})

test('対象なしは空配列を返す', async (t) => {
  t.mock.method(globalThis, 'fetch', async () => Response.json([]))
  assert.deepEqual(await fetchDailyPriceChanges('2026-10-19'), [])
})

test('HTTPエラーでは成功結果を返さない', async (t) => {
  t.mock.method(globalThis, 'fetch', async () => new Response('internal error', { status: 500 }))
  await assert.rejects(fetchDailyPriceChanges('2026-10-19'), /HTTP 500/)
})

test('通信失敗を呼び出し元へ返す', async (t) => {
  const failure = new TypeError('Failed to fetch')
  t.mock.method(globalThis, 'fetch', async () => { throw failure })
  await assert.rejects(fetchDailyPriceChanges('2026-10-19'), (error) => error === failure)
})
