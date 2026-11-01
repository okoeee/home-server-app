import type { WeatherForecastResponse } from './generated/schema'

export async function fetchWeatherForecast(): Promise<WeatherForecastResponse> {
  const response = await fetch('/api/v1/weather/forecast')

  if (!response.ok) {
    throw new Error(`天気予報の取得に失敗しました（HTTP ${response.status}）。`)
  }

  return response.json()
}
