<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import SlideshowControls from './SlideshowControls.vue'
import SlideshowDeck from './SlideshowDeck.vue'
import type { SlideshowSlide } from './slideshow'

const props = defineProps<{
  slides: SlideshowSlide[]
}>()

const root = ref<HTMLElement | null>(null)
const index = ref(0)
const hovered = ref(false)
const focused = ref(false)
const visible = ref(false)
let observer: IntersectionObserver | null = null

const playing = computed(() => props.slides.length > 1 && visible.value && !hovered.value && !focused.value)

function goTo(i: number) {
  const count = props.slides.length
  if (count === 0) return
  index.value = ((i % count) + count) % count
}

function onFocusOut(event: FocusEvent) {
  const next = event.relatedTarget
  if (!(next instanceof Node) || !root.value?.contains(next)) focused.value = false
}

onMounted(() => {
  if (!root.value) return
  observer = new IntersectionObserver(
    ([entry]) => {
      visible.value = entry?.isIntersecting ?? false
    },
    { threshold: 0.35 },
  )
  observer.observe(root.value)
})

onUnmounted(() => observer?.disconnect())
</script>

<template>
  <div
    v-if="slides.length > 0"
    ref="root"
    class="grid grid-cols-1 gap-12 lg:grid-cols-history lg:gap-x-split lg:gap-y-10"
    role="region"
    aria-roledescription="carousel"
    aria-label="Photos du salon 2024"
    @focusin="focused = true"
    @focusout="onFocusOut"
  >
    <div class="lg:col-start-1 lg:row-start-1">
      <slot name="intro" />
    </div>

    <div
      class="flex justify-center lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:items-center lg:justify-end"
      @mouseenter="hovered = true"
      @mouseleave="hovered = false"
    >
      <SlideshowDeck :slides="slides" :index="index" @next="goTo(index + 1)" />
    </div>

    <SlideshowControls
      class="lg:col-start-1 lg:row-start-2 lg:self-end"
      :slides="slides"
      :index="index"
      :playing="playing"
      @prev="goTo(index - 1)"
      @next="goTo(index + 1)"
      @select="goTo"
    />
  </div>
</template>
