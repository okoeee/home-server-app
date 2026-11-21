const refreshInterval = 15_000

export function startServerMetricsPolling(
  load: (signal: AbortSignal) => Promise<void>,
): () => void {
  const page = document
  const controller = new AbortController()
  let timer: ReturnType<typeof setInterval> | undefined
  let inFlight = false
  let stopped = false

  async function refresh() {
    if (stopped || page.hidden || inFlight) return

    inFlight = true
    try {
      await load(controller.signal)
    } finally {
      inFlight = false
    }
  }

  function updateVisibility() {
    clearInterval(timer)
    timer = undefined

    if (!page.hidden) {
      void refresh()
      timer = setInterval(() => void refresh(), refreshInterval)
    }
  }

  page.addEventListener('visibilitychange', updateVisibility)
  updateVisibility()

  return () => {
    stopped = true
    clearInterval(timer)
    page.removeEventListener('visibilitychange', updateVisibility)
    controller.abort()
  }
}
