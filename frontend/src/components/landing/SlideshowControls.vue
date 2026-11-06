<script setup lang="ts">
import LandingIcon from './LandingIcon.vue'
import type { SlideshowSlide } from './slideshow'

const props = defineProps<{
  slides: SlideshowSlide[]
  index: number
  playing: boolean
}>()

const emit = defineEmits<{
  prev: []
  next: []
  select: [index: number]
}>()

function pad(n: number) {
  return String(n).padStart(2, '0')
}

function fillClass(i: number) {
  if (i < props.index) return 'scale-x-100'
  if (i > props.index) return 'scale-x-0'
  return [
    'scale-x-0 animate-slide-progress motion-reduce:animate-none motion-reduce:scale-x-100',
    props.playing ? '' : '[animation-play-state:paused]',
  ]
}

const arrowClass =
  'grid size-btn-min cursor-pointer place-items-center rounded-full border border-line-strong bg-transparent text-forest transition-colors duration-200 hover:border-forest hover:bg-forest hover:text-white focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-gold'
</script>

<template>
  <div>
    <div class="flex items-end gap-3 font-serif text-forest" aria-hidden="true">
      <span class="text-date-sm leading-none tabular-nums md:text-date">{{ pad(index + 1) }}</span>
      <span class="mb-1.5 text-lede text-ink-subtle tabular-nums">/ {{ pad(slides.length) }}</span>
    </div>

    <div class="mt-3 min-h-16" :aria-live="playing ? 'off' : 'polite'">
      <Transition
        mode="out-in"
        enter-active-class="transition duration-500 ease-out"
        enter-from-class="translate-y-2 opacity-0"
        leave-active-class="transition duration-200 ease-in"
        leave-to-class="-translate-y-1 opacity-0"
      >
        <p :key="index" class="m-0 max-w-hero-desc font-serif text-subtitle italic leading-snug text-ink">
          <span class="sr-only">Photo {{ index + 1 }} sur {{ slides.length }} : </span>
          {{ slides[index]?.caption }}
        </p>
      </Transition>
    </div>

    <div class="mt-6 flex max-w-hero-desc items-center gap-1.5" role="group" aria-label="Choisir une photo">
      <button
        v-for="(slide, i) in slides"
        :key="slide.src"
        type="button"
        class="group flex-1 cursor-pointer border-0 bg-transparent py-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-gold"
        :aria-label="`Photo ${i + 1} : ${slide.caption}`"
        :aria-current="i === index || undefined"
        @click="emit('select', i)"
      >
        <span class="block h-0.5 overflow-hidden rounded-full bg-line-strong transition-colors group-hover:bg-ink-subtle">
          <span
            :key="i === index ? `active-${index}` : `idle-${i}`"
            class="block h-full w-full origin-left bg-forest"
            :class="fillClass(i)"
            @animationend="i === index && emit('next')"
          />
        </span>
      </button>
    </div>

    <div class="mt-7 flex gap-3">
      <button type="button" :class="arrowClass" aria-label="Photo précédente" @click="emit('prev')">
        <LandingIcon name="chevron" class="size-icon-md rotate-180" />
      </button>
      <button type="button" :class="arrowClass" aria-label="Photo suivante" @click="emit('next')">
        <LandingIcon name="chevron" class="size-icon-md" />
      </button>
    </div>
  </div>
</template>
