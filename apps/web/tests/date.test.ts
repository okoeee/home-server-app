import assert from 'node:assert/strict'
import { test } from 'node:test'
import { getTodayInJapan } from '../src/date.ts'

test('日本時間の午前0時を境に対象日を切り替える', () => {
  assert.equal(getTodayInJapan(new Date('2026-10-27T14:59:59Z')), '2026-10-27')
  assert.equal(getTodayInJapan(new Date('2026-10-27T15:00:00Z')), '2026-10-28')
})

test('日本時間で年をまたぐ日付を返す', () => {
  assert.equal(getTodayInJapan(new Date('2026-12-31T15:00:00Z')), '2027-01-01')
})
