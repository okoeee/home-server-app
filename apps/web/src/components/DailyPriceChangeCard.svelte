<script lang="ts">
  import TrendingUp from '@lucide/svelte/icons/trending-up'
  import type { DailyPriceChangeResponse } from '../api/generated/schema'
  import DailyPriceChangeTable from './DailyPriceChangeTable.svelte'
  import DashboardCard from './DashboardCard.svelte'

  let {
    date,
    changes,
    loading,
    error,
  }: {
    date: string | null
    changes: DailyPriceChangeResponse[]
    loading: boolean
    error: string
  } = $props()
</script>

<section class="min-w-0 w-full max-w-3xl" aria-labelledby="stocks-heading" aria-busy={loading}>
  <DashboardCard
    headingId="stocks-heading"
    title="株価の騰落率"
    icon={TrendingUp}
    description={date
      ? `${date}時点・保存済みの最新取引日`
      : '保存済みの日足から、騰落率の高い順に表示'}
  >
    {#if loading}
      <p class="py-12 text-center text-sm text-muted-foreground" role="status">読み込み中です。</p>
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
  </DashboardCard>
</section>
