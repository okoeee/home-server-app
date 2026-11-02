import assert from 'node:assert/strict'
import { test } from 'node:test'
import { fetchWeatherForecast } from '../src/api/weather-forecast.ts'

test('天気APIへ1回通信し、数値とnullをそのまま受け取る', async (t) => {
  // Arrange
  const result = {
    location_name: '東京都',
    fetched_at: '2026-11-01T00:30:00+09:00',
    forecasts: [
      {
        date: '2026-11-01',
        weather: 'clear',
        temperature_max: 21.5,
        temperature_min: 0,
        relative_humidity_mean: null,
      },
    ],
  }
  const fetchMock = t.mock.method(globalThis, 'fetch', async (url) => {
    assert.equal(url, '/api/v1/weather/forecast')
    return Response.json(result)
  })

  // Act
  const got = await fetchWeatherForecast()

  // Assert
  assert.deepEqual(got, result)
  assert.equal(fetchMock.mock.callCount(), 1)
})

test('天気APIのHTTPエラーを呼び出し元へ返す', async (t) => {
  // Arrange
  t.mock.method(globalThis, 'fetch', async () => new Response('upstream error', { status: 502 }))

  // Act
  const result = fetchWeatherForecast()

  // Assert
  await assert.rejects(result, /天気予報の取得に失敗しました（HTTP 502）/)
})

test('天気APIへの通信失敗を呼び出し元へ返す', async (t) => {
  // Arrange
  const failure = new TypeError('Failed to fetch')
  t.mock.method(globalThis, 'fetch', async () => {
    throw failure
  })

  // Act
  const result = fetchWeatherForecast()

  // Assert
  await assert.rejects(result, (error) => error === failure)
})
