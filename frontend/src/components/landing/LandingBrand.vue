<script setup lang="ts">
withDefaults(
  defineProps<{
    href?: string
    tone?: 'default' | 'footer'
  }>(),
  { href: '#haut', tone: 'default' },
)
</script>

<template>
  <a
    :href="href"
    class="inline-flex min-w-max items-center gap-3"
    aria-label="Salon de la culture et de la gastronomie comorienne, retour en haut"
  >
    <span
      class="grid size-8 place-items-center rounded-full bg-forest font-serif text-base font-semibold text-gold md:size-9 md:text-lg"
      :class="tone === 'footer' ? 'ring-1 ring-gold/40' : ''"
      aria-hidden="true"
    >
      S
    </span>
    <span class="text-nav leading-tight font-semibold tracking-brand md:text-base">
      SALON DE LA CULTURE
      <span
        class="block text-2xs font-normal tracking-normal md:text-eyebrow md:leading-normal"
        :class="tone === 'footer' ? 'text-footer-link' : 'text-ink-subtle'"
      >
        et de la gastronomie comorienne
      </span>
    </span>
  </a>
</template>
