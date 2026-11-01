import type { WeatherCondition } from '../api/generated/schema'

const weekdayFormatter = new Intl.DateTimeFormat('ja-JP', {
  weekday: 'short',
  timeZone: 'UTC',
})

export function formatForecastDate(date: string): string {
  const [year, month, day] = date.split('-').map(Number)
  // APIの暦日を保ったまま曜日を計算し、ブラウザのタイムゾーンへ変換しない。
  const weekday = weekdayFormatter.format(new Date(Date.UTC(year, month - 1, day)))
  return `${month}/${day}（${weekday}）`
}

export const weatherLabels: Record<WeatherCondition, string> = {
  clear: '快晴',
  mainly_clear: '晴れ',
  partly_cloudy: '晴れ（一部に雲）',
  overcast: '曇り',
  fog: '霧',
  drizzle: '霧雨',
  freezing_drizzle: '着氷性の霧雨',
  rain: '雨',
  freezing_rain: '着氷性の雨',
  snow: '雪',
  snow_grains: '霧雪',
  rain_showers: 'にわか雨',
  snow_showers: 'にわか雪',
  thunderstorm: '雷雨',
  thunderstorm_with_hail: 'ひょうを伴う雷雨',
  unknown: '不明',
}
