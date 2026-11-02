<script lang="ts">
  import { onMount } from 'svelte'
  import * as Card from '$lib/components/ui/card/index.js'
  import House from '@lucide/svelte/icons/house'
  import TrendingUp from '@lucide/svelte/icons/trending-up'
  import { fetchLatestDailyPriceChanges } from './api/daily-price-changes'
  import { fetchWeatherForecast } from './api/weather-forecast'
  import DailyPriceChangeTable from './components/DailyPriceChangeTable.svelte'
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

<main class="mx-auto grid max-w-7xl gap-6 px-5 py-7 sm:px-8 sm:py-9">
  <WeatherForecastCard {forecast} loading={weatherLoading} error={weatherError} />

  <section aria-labelledby="stocks-heading" aria-busy={loading}>
    <Card.Root class="gap-5 shadow-none">
      <Card.Header>
        <Card.Title>
          <h2
            id="stocks-heading"
            class="flex items-center gap-3 text-lg font-semibold tracking-tight"
          >
            <span
              class="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary"
              ><TrendingUp size={20} stroke-width={1.6} aria-hidden="true" /></span
            >
            株価の騰落率
          </h2>
        </Card.Title>
        <Card.Description class="mt-1 pl-12">
          {#if date}{date}時点・保存済みの最新取引日{:else}保存済みの日足から、騰落率の高い順に表示{/if}
        </Card.Description>
      </Card.Header>
      <Card.Content>
        {#if loading}
          <p class="py-12 text-center text-sm text-muted-foreground" role="status">
            読み込み中です。
          </p>
        {:else if error}
          <p
            class="rounded-xl border border-destructive/30 bg-destructive/10 p-5 text-sm text-destructive"
            role="alert"
          >
            {error}
          </p>
        {:else if changes.length === 0}
          <p class="py-12 text-center text-sm text-muted-foreground" role="status">
            データがありません。
          </p>
        {:else}
          <DailyPriceChangeTable {changes} />
        {/if}
      </Card.Content>
    </Card.Root>
  </section>
</main>
