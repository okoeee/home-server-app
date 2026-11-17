<script lang="ts">
  import Server from '@lucide/svelte/icons/server'
  import Cpu from '@lucide/svelte/icons/cpu'
  import MemoryStick from '@lucide/svelte/icons/memory-stick'
  import Clock from '@lucide/svelte/icons/clock'
  import type { ServerMetricsResponse } from '../api/generated/schema'
  import DashboardCard from './DashboardCard.svelte'

  let {
    metrics,
    loading,
    error,
  }: {
    metrics: ServerMetricsResponse | null
    loading: boolean
    error: string
  } = $props()
</script>

{#snippet footer()}
  {#if metrics}
    <p class="flex items-center gap-1.5">
      <Clock size={10} aria-hidden="true" />
      取得：<time datetime={metrics.fetched_at}
        >{new Date(metrics.fetched_at).toLocaleString('ja-JP', {
          timeZone: 'Asia/Tokyo',
        })}（日本時間）</time
      >
    </p>
  {/if}
{/snippet}

<section class="min-w-0 w-full" aria-labelledby="server-heading" aria-busy={loading}>
  <DashboardCard
    headingId="server-heading"
    title="サーバー"
    icon={Server}
    description="CPUは直近1分の全コア平均"
    footer={metrics && !loading && !error ? footer : undefined}
  >
    {#if loading}
      <p
        class="flex min-h-24 items-center justify-center text-sm text-muted-foreground"
        role="status"
      >
        サーバー情報を読み込み中です。
      </p>
    {:else if error}
      <p
        class="rounded-xl border border-destructive/30 bg-destructive/10 p-3 text-xs text-destructive"
        role="alert"
      >
        {error}
      </p>
    {:else if !metrics}
      <p class="py-6 text-center text-xs text-muted-foreground" role="status">
        サーバー情報のデータがありません。
      </p>
    {:else}
      <dl class="grid grid-cols-2 gap-3">
        {#each [{ label: 'CPU', value: metrics.cpu_usage_percent, icon: Cpu }, { label: 'メモリ', value: metrics.memory_usage_percent, icon: MemoryStick }] as metric (metric.label)}
          <div class="min-w-0 rounded-xl border bg-muted/40 p-3">
            <dt class="flex items-center gap-1.5 text-xs text-muted-foreground">
              <metric.icon size={14} class="shrink-0" stroke-width={1.6} aria-hidden="true" />
              {metric.label}<span class="sr-only">使用率</span>
            </dt>
            <dd class="mt-2">
              <div class="flex items-baseline gap-1 tabular-nums">
                <span class="text-2xl font-semibold tracking-tight">{metric.value.toFixed(1)}</span>
                <span class="text-xs text-muted-foreground">%</span>
              </div>
              <div class="mt-3 h-1.5 overflow-hidden rounded-full bg-border" aria-hidden="true">
                <div
                  class="h-full rounded-full bg-primary/80"
                  style:width={`${metric.value}%`}
                ></div>
              </div>
            </dd>
          </div>
        {/each}
      </dl>
    {/if}
  </DashboardCard>
</section>
