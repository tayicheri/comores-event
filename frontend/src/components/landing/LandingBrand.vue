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
  <a :href="href" class="inline-flex min-w-max items-center gap-3" aria-label="Comores Impact, retour en haut">
    <span class="grid place-items-center" aria-hidden="true">
      <img
        src="/landing/logo-comores-impact.png"
        alt=""
        width="32"
        height="32"
        class="size-7 rounded-sm object-contain md:size-8"
      />
    </span>
    <span class="text-nav leading-tight font-semibold tracking-brand md:text-base">
      COMORES IMPACT
      <span
        class="block text-2xs font-normal tracking-normal md:text-eyebrow md:leading-normal"
        :class="tone === 'footer' ? 'text-footer-link' : 'text-ink-subtle'"
      >
        La culture nous unit
      </span>
    </span>
  </a>
</template>
