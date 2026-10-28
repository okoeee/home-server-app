<script lang="ts">
  import type { DailyPriceChange } from '../models/daily-price-change'

  let { changes }: { changes: DailyPriceChange[] } = $props()
</script>

<div class="table-container">
  <table>
    <caption>騰落率の降順・上位10件（表示：{changes.length}件）</caption>
    <thead>
      <tr>
        <th scope="col">銘柄コード</th>
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
          <td>{change.trade_date}</td>
          <td class="number">{change.close_price}</td>
          <td>{change.previous_trade_date}</td>
          <td class="number">{change.previous_close_price}</td>
          <td class="number">{change.change_rate}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .table-container { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; white-space: nowrap; }
  caption { text-align: left; padding: 1rem 0; font-weight: bold; }
  th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #d8dee6; }
  th { background: #f3f5f7; }
  .number { text-align: right; font-variant-numeric: tabular-nums; }
</style>
