<script lang="ts">
  import { cn, type WithElementRef } from '$lib/utils.js'
  import type { HTMLTableAttributes } from 'svelte/elements'

  let {
    ref = $bindable(null),
    class: className,
    children,
    ...restProps
  }: WithElementRef<HTMLTableAttributes> = $props()
</script>

<div data-slot="table-container" class="relative min-w-0 w-full max-w-full overflow-x-auto">
  <table
    bind:this={ref}
    data-slot="table"
    class={cn('w-full caption-bottom text-sm', className)}
    {...restProps}
  >
    {@render children?.()}
  </table>
</div>
