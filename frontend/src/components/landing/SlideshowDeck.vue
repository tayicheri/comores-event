<script setup lang="ts">
import { computed } from 'vue'
import type { SlideshowSlide } from './slideshow'

const props = defineProps<{
  slides: SlideshowSlide[]
  index: number
}>()

const emit = defineEmits<{ next: [] }>()

const offsetClass: Record<number, string> = {
  0: 'z-30 translate-x-0 rotate-0 scale-100 cursor-pointer',
  1: 'z-20 translate-x-6 rotate-3 scale-95 cursor-pointer md:translate-x-10',
  2: 'z-10 translate-x-12 rotate-6 scale-90 opacity-70 md:translate-x-20',
}

const deck = computed(() => {
  const depth = Math.min(3, props.slides.length)
  return Array.from({ length: depth }, (_, offset) => ({
    offset,
    slide: props.slides[(props.index + offset) % props.slides.length]!,
  }))
})
</script>

<template>
  <div class="relative pr-12 md:pr-20">
    <TransitionGroup
      tag="div"
      class="relative aspect-4/5 h-hero-img-sm sm:h-hero-img-md lg:h-hero-img-lg"
      enter-active-class="transition-all duration-700 ease-out"
      enter-from-class="opacity-0! translate-x-24! rotate-12!"
      leave-active-class="transition-all duration-500 ease-in"
      leave-to-class="opacity-0! -translate-x-1/3! -rotate-6!"
    >
      <figure
        v-for="card in deck"
        :key="card.slide.src"
        class="absolute inset-0 m-0 overflow-hidden rounded-hero-sm bg-media shadow-modal transition-all duration-700 ease-out md:rounded-hero-md lg:rounded-hero"
        :class="offsetClass[card.offset]"
        :aria-hidden="card.offset !== 0 || undefined"
        @click="card.offset <= 1 && emit('next')"
      >
        <img
          :src="card.slide.src"
          :alt="card.offset === 0 ? card.slide.alt : ''"
          class="size-full object-cover transition-[filter] duration-700"
          :class="card.offset === 0 ? 'animate-ken-burns motion-reduce:animate-none' : 'brightness-75 saturate-50'"
        />
        <template v-if="card.offset === 0">
          <div class="pointer-events-none absolute inset-x-0 top-1/2 bottom-0 bg-hero-fade" />
          <figcaption class="absolute right-6 bottom-6 left-7 text-xs font-bold tracking-eyebrow text-white">
            SALON 2024 · 1<sup>re</sup> ÉDITION
          </figcaption>
        </template>
      </figure>
    </TransitionGroup>

    <span
      class="absolute -top-4 -left-4 z-40 flex size-seal-sm flex-col items-center justify-center rounded-full bg-gold font-serif text-2xl leading-none text-forest shadow-modal md:-top-6 md:-left-6 md:size-seal md:text-seal"
      aria-hidden="true"
    >
      2024
      <small class="mt-1.5 font-sans text-micro font-bold tracking-brand md:text-2xs">SOUVENIRS</small>
    </span>
  </div>
</template>
