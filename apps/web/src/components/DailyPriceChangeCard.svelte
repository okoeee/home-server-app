<script lang="ts">
  import * as Card from '$lib/components/ui/card/index.js'
  import TrendingUp from '@lucide/svelte/icons/trending-up'
  import type { DailyPriceChangeResponse } from '../api/generated/schema'
  import DailyPriceChangeTable from './DailyPriceChangeTable.svelte'

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
  <Card.Root class="min-w-0 w-full max-w-full gap-5 shadow-none">
    <Card.Header class="min-w-0 grid-cols-[minmax(0,1fr)]">
      <Card.Title>
        <h2
          id="stocks-heading"
          class="flex items-center gap-3 text-sm font-semibold tracking-tight"
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
    <Card.Content class="min-w-0 w-full">
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
