<script lang="ts">
  import * as Table from '$lib/components/ui/table/index.js'
  import type { DailyPriceChangeResponse } from '../api/generated/schema'
  import { formatChangeRate } from '../models/daily-price-change'

  let { changes }: { changes: DailyPriceChangeResponse[] } = $props()
</script>

<Table.Root class="min-w-[520px] table-fixed">
  <Table.Caption class="sr-only">株価の騰落率</Table.Caption>
  <colgroup>
    <col class="w-[14%]" />
    <col class="w-[32%]" />
    <col class="w-[18%]" />
    <col class="w-[18%]" />
    <col class="w-[18%]" />
  </colgroup>
  <Table.Header class="bg-muted">
    <Table.Row>
      <Table.Head scope="col" class="whitespace-normal px-3 py-2 text-xs text-muted-foreground"
        >銘柄コード</Table.Head
      >
      <Table.Head scope="col" class="whitespace-normal px-3 py-2 text-xs text-muted-foreground"
        >会社名</Table.Head
      >
      <Table.Head
        scope="col"
        class="whitespace-normal px-3 py-2 text-right text-xs text-muted-foreground"
        >終値（円）</Table.Head
      >
      <Table.Head
        scope="col"
        class="whitespace-normal px-3 py-2 text-right text-xs text-muted-foreground"
        >比較元の終値（円）</Table.Head
      >
      <Table.Head
        scope="col"
        class="whitespace-normal px-3 py-2 text-right text-xs text-muted-foreground"
        >騰落率（%）</Table.Head
      >
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#each changes as change (change.symbol)}
      <Table.Row>
        <Table.Cell class="px-3 py-4 tabular-nums">{change.symbol}</Table.Cell>
        <Table.Cell class="whitespace-normal px-3 py-4 leading-relaxed [overflow-wrap:anywhere]"
          >{change.company_name ?? '名称未登録'}</Table.Cell
        >
        <Table.Cell class="px-3 py-4 text-right tabular-nums">{change.close_price}</Table.Cell>
        <Table.Cell class="px-3 py-4 text-right tabular-nums"
          >{change.previous_close_price}</Table.Cell
        >
        <Table.Cell
          class={`px-3 py-4 text-right font-semibold tabular-nums ${Number(change.change_rate) > 0 ? 'text-primary' : Number(change.change_rate) < 0 ? 'text-destructive' : 'text-muted-foreground'}`}
        >
          {formatChangeRate(change.change_rate)}
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
