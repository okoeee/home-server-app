<script lang="ts">
  import type { DailyPriceChangeResponse } from '../api/generated/schema'

  let { changes }: { changes: DailyPriceChangeResponse[] } = $props()
</script>

<div class="table-container">
  <table>
    <caption>騰落率の降順・上位10件（表示：{changes.length}件）</caption>
    <thead>
      <tr>
        <th scope="col">銘柄コード</th>
        <th scope="col">会社名</th>
        <th scope="col">対象日</th>
        <th scope="col">終値（円）</th>
        <th scope="col">比較元の日付</th>
        <th scope="col">比較元の終値（円）</th>
        <th scope="col">騰落率（%）</th>
      </tr>
    </thead>
    <tbody>
      {#each changes as change (change.symbol)}
        <tr>
          <td>{change.symbol}</td>
          <td>{change.company_name ?? '名称未登録'}</td>
          <td>{change.trade_date}</td>
          <td class="number">{change.close_price}</td>
          <td>{change.previous_trade_date}</td>
          <td class="number">{change.previous_close_price}</td>
          <td class="number change-rate" class:positive={Number(change.change_rate) > 0} class:negative={Number(change.change_rate) < 0}>{change.change_rate}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .table-container { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; white-space: nowrap; }
  caption { text-align: left; padding: 0 0 1rem; font-size: 0.75rem; color: var(--muted-foreground); }
  th, td { padding: 1rem 0.85rem; text-align: left; border-bottom: 1px solid var(--border); font-size: 0.8125rem; }
  th { background: var(--muted); color: var(--muted-foreground); font-size: 0.6875rem; font-weight: 500; }
  tbody tr:last-child td { border-bottom: 0; }
  tbody tr:hover { background: color-mix(in srgb, var(--muted) 45%, transparent); }
  .number { text-align: right; font-variant-numeric: tabular-nums; }
  .change-rate { font-weight: 600; }
  .positive { color: var(--primary); }
  .negative { color: var(--destructive); }
</style>
