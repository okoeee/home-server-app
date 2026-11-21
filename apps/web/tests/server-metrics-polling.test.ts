import assert from 'node:assert/strict'
import { test, type TestContext } from 'node:test'
import { startServerMetricsPolling } from '../src/lib/server-metrics-polling.ts'

function preparePage(t: TestContext, hidden = false) {
  const page = Object.assign(new EventTarget(), { hidden })
  const original = Object.getOwnPropertyDescriptor(globalThis, 'document')
  Object.defineProperty(globalThis, 'document', { configurable: true, value: page })
  t.mock.timers.enable({ apis: ['setInterval'] })
  t.after(() => {
    if (original) Object.defineProperty(globalThis, 'document', original)
    else Reflect.deleteProperty(globalThis, 'document')
  })
  return page
}

test('初回に取得し、その後15秒ごとに更新する', async (t) => {
  // Arrange
  preparePage(t)
  const load = t.mock.fn(async () => {})

  // Act
  const stop = startServerMetricsPolling(load)
  t.after(stop)

  // Assert
  assert.equal(load.mock.callCount(), 1)
  await Promise.resolve()
  t.mock.timers.tick(14_999)
  assert.equal(load.mock.callCount(), 1)
  t.mock.timers.tick(1)
  assert.equal(load.mock.callCount(), 2)
  await Promise.resolve()
  t.mock.timers.tick(15_000)
  assert.equal(load.mock.callCount(), 3)
})

test('通信中はタイマーやタブの再表示で取得を重複させない', async (t) => {
  // Arrange
  const page = preparePage(t)
  const pending = Promise.withResolvers<void>()
  const load = t.mock.fn(() => pending.promise)
  const stop = startServerMetricsPolling(load)
  t.after(stop)

  // Act
  t.mock.timers.tick(30_000)
  page.hidden = true
  page.dispatchEvent(new Event('visibilitychange'))
  page.hidden = false
  page.dispatchEvent(new Event('visibilitychange'))

  // Assert
  assert.equal(load.mock.callCount(), 1)
  pending.resolve()
  await Promise.resolve()
  t.mock.timers.tick(15_000)
  assert.equal(load.mock.callCount(), 2)
})

test('非表示中は取得せず、表示に戻ったら即取得して定期更新を再開する', async (t) => {
  // Arrange
  const page = preparePage(t, true)
  const load = t.mock.fn(async () => {})
  const stop = startServerMetricsPolling(load)
  t.after(stop)

  // Act
  t.mock.timers.tick(30_000)

  // Assert
  assert.equal(load.mock.callCount(), 0)
  page.hidden = false
  page.dispatchEvent(new Event('visibilitychange'))
  assert.equal(load.mock.callCount(), 1)
  await Promise.resolve()
  page.hidden = true
  page.dispatchEvent(new Event('visibilitychange'))
  t.mock.timers.tick(60_000)
  assert.equal(load.mock.callCount(), 1)
  page.hidden = false
  page.dispatchEvent(new Event('visibilitychange'))
  assert.equal(load.mock.callCount(), 2)
  await Promise.resolve()
  t.mock.timers.tick(15_000)
  assert.equal(load.mock.callCount(), 3)
})

test('停止時に通信をキャンセルし、タイマーとイベントによる再取得を止める', async (t) => {
  // Arrange
  const page = preparePage(t)
  const pending = Promise.withResolvers<void>()
  const load = t.mock.fn((_signal: AbortSignal) => pending.promise)
  const stop = startServerMetricsPolling(load)
  t.after(stop)
  const signal = load.mock.calls[0].arguments[0]

  // Act
  stop()
  pending.resolve()
  await Promise.resolve()
  t.mock.timers.tick(30_000)
  page.dispatchEvent(new Event('visibilitychange'))

  // Assert
  assert.equal(signal.aborted, true)
  assert.equal(load.mock.callCount(), 1)
})
