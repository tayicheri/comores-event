<script setup lang="ts">
defineProps<{
  src: string
  alt: string
  name: string
  place: string
}>()
</script>

<template>
  <figure class="group relative m-0 overflow-hidden rounded-media-lg bg-forest-deep">
    <img
      :src="src"
      :alt="alt"
      loading="lazy"
      class="size-full object-cover transition-transform duration-700 group-hover:scale-105 motion-reduce:transition-none"
    />
    <div class="pointer-events-none absolute inset-0 bg-hero-fade" />
    <figcaption class="absolute right-4 bottom-4 left-4 text-white md:right-5 md:bottom-5 md:left-5">
      <span class="block text-2xs font-bold tracking-eyebrow text-gold">{{ place }}</span>
      <span class="mt-1 block font-serif text-lede leading-snug">{{ name }}</span>
    </figcaption>
  </figure>
</template>
