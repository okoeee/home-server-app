<script lang="ts">
  import * as Card from '$lib/components/ui/card/index.js'
  import CloudSun from '@lucide/svelte/icons/cloud-sun'
  import Droplets from '@lucide/svelte/icons/droplets'
  import Clock from '@lucide/svelte/icons/clock'
  import type { WeatherForecastResponse } from '../api/generated/schema'
  import { formatForecastDate, weatherLabels } from '../models/weather-forecast'
  import WeatherIcon from './WeatherIcon.svelte'

  let {
    forecast,
    loading,
    error,
  }: {
    forecast: WeatherForecastResponse | null
    loading: boolean
    error: string
  } = $props()
</script>

<section class="w-full max-w-md" aria-labelledby="weather-heading" aria-busy={loading}>
  <Card.Root size="sm" class="gap-3 shadow-none">
    <Card.Header>
      <Card.Title>
        <h2
          id="weather-heading"
          class="flex items-center gap-2 text-sm font-semibold tracking-tight"
        >
          <CloudSun size={16} class="text-primary" stroke-width={1.6} aria-hidden="true" />
          {forecast?.location_name ?? '東京都'}の天気
        </h2>
      </Card.Title>
      <Card.Description class="mt-1 text-[10px]">最高 / 最低（℃）・日平均湿度</Card.Description>
      <Card.Action
        ><span
          class="rounded-full border border-border bg-muted/50 px-2 py-0.5 text-[10px] text-muted-foreground"
          >5日間</span
        ></Card.Action
      >
    </Card.Header>

    <Card.Content>
      {#if loading}
        <p
          class="flex min-h-24 items-center justify-center text-sm text-muted-foreground"
          role="status"
        >
          天気予報を読み込み中です。
        </p>
      {:else if error}
        <p
          class="rounded-xl border border-destructive/30 bg-destructive/10 p-3 text-xs text-destructive"
          role="alert"
        >
          {error}
        </p>
      {:else if !forecast || forecast.forecasts.length === 0}
        <p class="py-6 text-center text-xs text-muted-foreground" role="status">
          天気予報のデータがありません。
        </p>
      {:else}
        <ul class="grid grid-cols-5 gap-1">
          {#each forecast.forecasts as day, index (day.date)}
            <li class="flex min-w-0 flex-col items-center gap-1.5 py-1">
              <h3
                class="whitespace-nowrap text-[10px] font-medium text-muted-foreground"
                class:text-primary={index === 0}
              >
                {#if index === 0}<span class="sr-only">今日：</span>{/if}
                <time datetime={day.date}>{formatForecastDate(day.date)}</time>
              </h3>
              <WeatherIcon weather={day.weather} size={24} />
              <p
                class="w-full truncate text-center text-[10px] text-muted-foreground"
                title={weatherLabels[day.weather] ?? weatherLabels.unknown}
              >
                {weatherLabels[day.weather] ?? weatherLabels.unknown}
              </p>
              <dl class="w-full text-[11px] tabular-nums">
                <div class="flex flex-wrap justify-center gap-x-1">
                  <dt class="sr-only">最高気温（℃）</dt>
                  <dd class="font-medium">
                    {day.temperature_max === null ? '—' : `${day.temperature_max}°`}
                  </dd>
                  <dt class="sr-only">最低気温（℃）</dt>
                  <dd class="text-muted-foreground">
                    {day.temperature_min === null ? '—' : `${day.temperature_min}°`}
                  </dd>
                </div>
                <div
                  class="mt-1 flex items-center justify-center gap-0.5 text-[10px] text-muted-foreground"
                >
                  <dt>
                    <span class="sr-only">平均湿度</span><Droplets size={10} aria-hidden="true" />
                  </dt>
                  <dd>
                    {day.relative_humidity_mean === null ? '—' : `${day.relative_humidity_mean}%`}
                  </dd>
                </div>
              </dl>
            </li>
          {/each}
        </ul>
      {/if}
    </Card.Content>

    <Card.Footer
      class="flex-wrap justify-between gap-x-3 gap-y-1 text-[10px] text-muted-foreground"
    >
      {#if forecast && !loading && !error}
        <p class="flex items-center gap-1.5">
          <Clock size={10} aria-hidden="true" />
          取得：<time datetime={forecast.fetched_at}
            >{new Date(forecast.fetched_at).toLocaleString('ja-JP', {
              timeZone: 'Asia/Tokyo',
            })}（日本時間）</time
          >
        </p>
      {/if}
      <a
        class="ml-auto underline underline-offset-4 hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-ring"
        href="https://open-meteo.com/">Weather data by Open-Meteo</a
      >
    </Card.Footer>
  </Card.Root>
</section>
