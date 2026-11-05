<script setup lang="ts">
import LandingIcon from './LandingIcon.vue'

withDefaults(
  defineProps<{
    href?: string
    tone?: 'default' | 'footer'
  }>(),
  { href: '#haut', tone: 'default' },
)
</script>

<template>
  <a :href="href" class="inline-flex min-w-max items-center gap-3" aria-label="Comores Events, retour en haut">
    <span
      class="grid place-items-center"
      :class="tone === 'footer' ? 'text-gold' : 'text-forest'"
      aria-hidden="true"
    >
      <LandingIcon name="star" class="size-7 -rotate-6 md:size-8" />
    </span>
    <span class="text-nav leading-tight font-semibold tracking-brand md:text-base">
      COMORES EVENTS
      <span
        class="block text-2xs font-normal tracking-normal md:text-eyebrow md:leading-normal"
        :class="tone === 'footer' ? 'text-footer-link' : 'text-ink-subtle'"
      >
        La culture nous unit
      </span>
    </span>
  </a>
</template>
