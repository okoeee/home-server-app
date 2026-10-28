<script lang="ts">
  import { onMount } from 'svelte'
  import { fetchLatestDailyPriceChanges } from './api/daily-price-changes'
  import DailyPriceChangeTable from './components/DailyPriceChangeTable.svelte'
  import type { DailyPriceChange } from './models/daily-price-change'
  import { selectTopDailyPriceChanges } from './models/daily-price-change'

  let date = $state<string | null>(null)
  let changes = $state<DailyPriceChange[]>([])
  let loading = $state(true)
  let error = $state('')

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
  })
</script>

<main>
  <h1>株価の騰落率</h1>

  {#if date}
    <p>{date}時点（DBの最新取引日）</p>
  {/if}

  <section aria-label="取得結果" aria-busy={loading}>
    {#if loading}
      <p role="status">読み込み中です。</p>
    {:else if error}
      <p role="alert">{error}</p>
    {:else if changes.length === 0}
      <p role="status">データがありません。</p>
    {:else}
      <DailyPriceChangeTable {changes} />
    {/if}
  </section>
</main>

<style>
  main { max-width: 1100px; margin: 0 auto; padding: 2rem 1rem; }
  h1 { margin-top: 0; }
  [role='alert'] { color: #a32121; }
</style>
