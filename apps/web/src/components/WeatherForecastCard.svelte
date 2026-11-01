<script lang="ts">
  import type { WeatherForecastResponse } from '../api/generated/schema'
  import { formatForecastDate, weatherLabels } from '../models/weather-forecast'

  let { forecast, loading, error }: {
    forecast: WeatherForecastResponse | null
    loading: boolean
    error: string
  } = $props()
</script>

<section class="weather-card" aria-labelledby="weather-heading" aria-busy={loading}>
  <h2 id="weather-heading">{forecast?.location_name ?? '東京都'}の天気</h2>
  <p class="description">今日を含む5日間・東京都庁付近の予報（日本時間）</p>

  {#if loading}
    <p role="status">天気予報を読み込み中です。</p>
  {:else if error}
    <p role="alert">{error}</p>
  {:else if !forecast || forecast.forecasts.length === 0}
    <p role="status">天気予報のデータがありません。</p>
  {:else}
    <ul class="forecast-days">
      {#each forecast.forecasts as day (day.date)}
        <li class="forecast-day">
          <h3><time datetime={day.date}>{formatForecastDate(day.date)}</time></h3>
          <p class="condition">{weatherLabels[day.weather] ?? weatherLabels.unknown}</p>
          <dl>
            <div><dt>最高気温</dt><dd>{day.temperature_max === null ? '—' : `${day.temperature_max} ℃`}</dd></div>
            <div><dt>最低気温</dt><dd>{day.temperature_min === null ? '—' : `${day.temperature_min} ℃`}</dd></div>
            <div><dt>平均湿度</dt><dd>{day.relative_humidity_mean === null ? '—' : `${day.relative_humidity_mean} %`}</dd></div>
          </dl>
        </li>
      {/each}
    </ul>
    <p class="updated-at">
      取得日時：<time datetime={forecast.fetched_at}>{new Date(forecast.fetched_at).toLocaleString('ja-JP', { timeZone: 'Asia/Tokyo' })}（日本時間）</time>
    </p>
  {/if}

  <p class="attribution">天気データ：<a href="https://open-meteo.com/">Open-Meteo</a></p>
</section>

<style>
  .weather-card { padding: 1.5rem; border: 1px solid #d8dee6; border-radius: 1rem; margin-bottom: 2rem; }
  h2 { margin: 0; }
  .description, .updated-at, .attribution { color: #526070; font-size: 0.875rem; }
  .forecast-days { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); gap: 0.75rem; list-style: none; padding: 0; margin: 1.5rem 0; }
  .forecast-day { background: #f3f5f7; border-radius: 0.75rem; padding: 1rem; }
  h3 { margin: 0; font-size: 0.9rem; }
  .condition { min-height: 3em; font-weight: bold; line-height: 1.5; }
  dl { margin: 0; font-size: 0.875rem; }
  dl div { display: flex; flex-wrap: wrap; justify-content: space-between; gap: 0.25rem; margin-top: 0.75rem; }
  dd { margin: 0; font-variant-numeric: tabular-nums; }
  [role='alert'] { color: #a32121; }
  .attribution { margin-bottom: 0; }
  @media (max-width: 800px) { .forecast-days { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
  @media (max-width: 400px) { .forecast-days { grid-template-columns: 1fr; } }
</style>
