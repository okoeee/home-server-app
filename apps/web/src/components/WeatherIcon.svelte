<script lang="ts">
  import Sun from '@lucide/svelte/icons/sun'
  import CloudSun from '@lucide/svelte/icons/cloud-sun'
  import Cloud from '@lucide/svelte/icons/cloud'
  import CloudFog from '@lucide/svelte/icons/cloud-fog'
  import CloudDrizzle from '@lucide/svelte/icons/cloud-drizzle'
  import CloudRain from '@lucide/svelte/icons/cloud-rain'
  import CloudRainWind from '@lucide/svelte/icons/cloud-rain-wind'
  import CloudSnow from '@lucide/svelte/icons/cloud-snow'
  import Snowflake from '@lucide/svelte/icons/snowflake'
  import CloudLightning from '@lucide/svelte/icons/cloud-lightning'
  import CloudHail from '@lucide/svelte/icons/cloud-hail'
  import CircleQuestionMark from '@lucide/svelte/icons/circle-question-mark'
  import type { WeatherCondition } from '../api/generated/schema'

  let { weather, size = 40 }: { weather: WeatherCondition; size?: number } = $props()

  const icons = {
    clear: { icon: Sun, color: 'text-amber-300' },
    mainly_clear: { icon: Sun, color: 'text-amber-300' },
    partly_cloudy: { icon: CloudSun, color: 'text-amber-200' },
    overcast: { icon: Cloud, color: 'text-stone-300' },
    fog: { icon: CloudFog, color: 'text-stone-400' },
    drizzle: { icon: CloudDrizzle, color: 'text-sky-300' },
    freezing_drizzle: { icon: CloudDrizzle, color: 'text-cyan-200' },
    rain: { icon: CloudRain, color: 'text-sky-300' },
    freezing_rain: { icon: CloudRain, color: 'text-cyan-200' },
    snow: { icon: CloudSnow, color: 'text-cyan-100' },
    snow_grains: { icon: Snowflake, color: 'text-cyan-100' },
    rain_showers: { icon: CloudRainWind, color: 'text-sky-300' },
    snow_showers: { icon: CloudSnow, color: 'text-cyan-100' },
    thunderstorm: { icon: CloudLightning, color: 'text-violet-300' },
    thunderstorm_with_hail: { icon: CloudHail, color: 'text-violet-300' },
    unknown: { icon: CircleQuestionMark, color: 'text-muted-foreground' },
  } satisfies Record<WeatherCondition, { icon: typeof Sun; color: string }>

  const presentation = $derived(icons[weather] ?? icons.unknown)
  const Icon = $derived(presentation.icon)
</script>

<Icon {size} stroke-width={1.5} class={presentation.color} aria-hidden="true" />
