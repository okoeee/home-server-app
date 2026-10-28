import assert from 'node:assert/strict'
import { test } from 'node:test'
import { selectTopDailyPriceChanges } from '../src/models/daily-price-change.ts'
import type { DailyPriceChange } from '../src/models/daily-price-change.ts'

function priceChange(symbol: string, rate: string): DailyPriceChange {
  return {
    symbol, change_rate: rate, trade_date: '2026-10-27',
    close_price: '100', previous_trade_date: '2026-10-26', previous_close_price: '100',
  }
}

test('文字列ではなく数値の降順で並べ、元の配列は変更しない', () => {
  const changes = ['2', '-10', '10.5', '0', '-1.5'].map((rate, index) => priceChange(String(index), rate))
  const original = [...changes]
  assert.deepEqual(selectTopDailyPriceChanges(changes).map(row => row.change_rate), ['10.5', '2', '0', '-1.5', '-10'])
  assert.deepEqual(changes, original)
})

test('11件以上なら並べ替えた後の先頭10件だけを返す', () => {
  const changes = Array.from({ length: 12 }, (_, index) => priceChange(String(index), String(index)))
  assert.deepEqual(selectTopDailyPriceChanges(changes).map(row => row.change_rate), ['11', '10', '9', '8', '7', '6', '5', '4', '3', '2'])
})

test('同じ騰落率の銘柄は元の順番を保つ', () => {
  const changes = [priceChange('12340', '2.0'), priceChange('56780', '2')]
  assert.deepEqual(selectTopDailyPriceChanges(changes), changes)
})

test('空配列は空配列を返す', () => {
  assert.deepEqual(selectTopDailyPriceChanges([]), [])
})
