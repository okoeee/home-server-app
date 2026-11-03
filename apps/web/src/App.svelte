<script lang="ts">
  import { onMount } from 'svelte'
  import House from '@lucide/svelte/icons/house'
  import { fetchLatestDailyPriceChanges } from './api/daily-price-changes'
  import { fetchWeatherForecast } from './api/weather-forecast'
  import DailyPriceChangeCard from './components/DailyPriceChangeCard.svelte'
  import WeatherForecastCard from './components/WeatherForecastCard.svelte'
  import type { DailyPriceChangeResponse, WeatherForecastResponse } from './api/generated/schema'
  import { selectTopDailyPriceChanges } from './models/daily-price-change'

  let date = $state<string | null>(null)
  let changes = $state<DailyPriceChangeResponse[]>([])
  let loading = $state(true)
  let error = $state('')
  let forecast = $state<WeatherForecastResponse | null>(null)
  let weatherLoading = $state(true)
  let weatherError = $state('')

  async function loadWeather() {
    try {
      forecast = await fetchWeatherForecast()
    } catch (cause) {
      weatherError = cause instanceof Error ? cause.message : '天気予報を取得できませんでした。'
    } finally {
      weatherLoading = false
    }
  }

  async function loadChanges() {
    try {
      const result = await fetchLatestDailyPriceChanges()
      date = result.trade_date
      changes = selectTopDailyPriceChanges(result.changes)
    } catch (cause) {
      error = cause instanceof Error ? cause.message : '騰落率を取得できませんでした。'
    } finally {
      loading = false
    }
  }

  onMount(() => {
    void loadChanges()
    void loadWeather()
  })
</script>

<header class="border-b border-border/70">
  <div class="mx-auto flex max-w-7xl items-center gap-3 px-5 py-5 sm:px-8">
    <span
      class="flex size-10 items-center justify-center rounded-xl bg-primary text-primary-foreground"
      ><House size={21} stroke-width={1.6} aria-hidden="true" /></span
    >
    <div>
      <h1 class="text-base font-semibold tracking-tight">My Home</h1>
      <p class="mt-0.5 text-xs text-muted-foreground">暮らしとマーケットのダッシュボード</p>
    </div>
  </div>
</header>

<main
  class="mx-auto box-border grid w-full min-w-0 max-w-7xl grid-cols-1 gap-6 px-5 py-7 sm:px-8 sm:py-9"
>
  <WeatherForecastCard {forecast} loading={weatherLoading} error={weatherError} />

  <DailyPriceChangeCard {date} {changes} {loading} {error} />
</main>
