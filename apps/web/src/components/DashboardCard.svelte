<script lang="ts">
  import * as Card from '$lib/components/ui/card/index.js'
  import type { LucideIcon } from '@lucide/svelte'
  import type { Snippet } from 'svelte'

  let {
    headingId,
    title,
    icon: Icon,
    description,
    children,
    footer,
  }: {
    headingId: string
    title: string
    icon: LucideIcon
    description: string
    children: Snippet
    footer?: Snippet
  } = $props()
</script>

<Card.Root size="sm" class="min-w-0 w-full max-w-full gap-3 shadow-none">
  <Card.Header class="min-w-0 grid-cols-[minmax(0,1fr)]">
    <Card.Title>
      <h2 id={headingId} class="flex items-center gap-2 text-sm font-semibold tracking-tight">
        <Icon size={16} class="shrink-0 text-primary" stroke-width={1.6} aria-hidden="true" />
        {title}
      </h2>
    </Card.Title>
    <Card.Description class="mt-1 text-[10px]">{description}</Card.Description>
  </Card.Header>

  <Card.Content class="min-w-0 w-full">
    {@render children()}
  </Card.Content>

  {#if footer}
    <Card.Footer
      class="flex-wrap justify-between gap-x-3 gap-y-1 text-[10px] text-muted-foreground"
    >
      {@render footer()}
    </Card.Footer>
  {/if}
</Card.Root>
