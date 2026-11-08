<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import LandingIcon from './LandingIcon.vue'
import UnescoSiteCard from './UnescoSiteCard.vue'
import type { HeritageSite } from './heritageSites'

const props = defineProps<{
  sites: readonly HeritageSite[]
}>()

const root = ref<HTMLElement | null>(null)
const index = ref(0)
const hovered = ref(false)
const focused = ref(false)
const visible = ref(false)
let observer: IntersectionObserver | null = null

const current = computed(() => props.sites[index.value])
const playing = computed(() => props.sites.length > 1 && visible.value && !hovered.value && !focused.value)

function goTo(i: number) {
  const count = props.sites.length
  if (count === 0) return
  index.value = ((i % count) + count) % count
}

function pad(n: number) {
  return String(n).padStart(2, '0')
}

function fillClass(i: number) {
  if (i < index.value) return 'scale-x-100'
  if (i > index.value) return 'scale-x-0'
  return [
    'scale-x-0 animate-slide-progress motion-reduce:animate-none motion-reduce:scale-x-100',
    playing.value ? '' : '[animation-play-state:paused]',
  ]
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

const arrowClass =
  'grid size-btn-min cursor-pointer place-items-center rounded-full border border-gold/40 bg-transparent text-gold transition-colors duration-200 hover:border-gold hover:bg-gold hover:text-forest focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-gold'
</script>

<template>
  <div
    v-if="sites.length > 0 && current"
    ref="root"
    class="flex flex-col gap-5"
    role="region"
    aria-roledescription="carousel"
    aria-label="Sites du patrimoine mondial"
    @focusin="focused = true"
    @focusout="onFocusOut"
    @mouseenter="hovered = true"
    @mouseleave="hovered = false"
  >
    <div
      class="h-hero-img-sm overflow-hidden rounded-media-lg"
      :class="sites.length > 1 ? 'cursor-pointer' : undefined"
      :role="sites.length > 1 ? 'button' : undefined"
      :tabindex="sites.length > 1 ? 0 : undefined"
      :aria-label="sites.length > 1 ? 'Site suivant' : undefined"
      @click="sites.length > 1 && goTo(index + 1)"
      @keydown.enter.prevent="sites.length > 1 && goTo(index + 1)"
      @keydown.space.prevent="sites.length > 1 && goTo(index + 1)"
    >
      <Transition
        mode="out-in"
        enter-active-class="transition duration-500 ease-out motion-reduce:transition-none"
        enter-from-class="translate-x-4 opacity-0"
        leave-active-class="transition duration-300 ease-in motion-reduce:transition-none"
        leave-to-class="-translate-x-4 opacity-0"
      >
        <UnescoSiteCard
          :key="current.name"
          class="pointer-events-none h-full w-full"
          :src="current.src"
          :alt="current.alt"
          :name="current.name"
          :place="`${current.city.toUpperCase()} · ${current.island.toUpperCase()}`"
        />
      </Transition>
    </div>

    <div class="flex items-end gap-3 font-serif text-white" aria-hidden="true">
      <span class="text-date-sm leading-none tabular-nums">{{ pad(index + 1) }}</span>
      <span class="mb-1.5 text-lede text-heritage-muted tabular-nums">/ {{ pad(sites.length) }}</span>
    </div>

    <p class="sr-only" :aria-live="playing ? 'off' : 'polite'">
      Site {{ index + 1 }} sur {{ sites.length }} : {{ current.name }}, {{ current.city }},
      {{ current.island }}
    </p>

    <div class="flex items-center gap-1.5" role="group" aria-label="Choisir un site">
      <button
        v-for="(site, i) in sites"
        :key="site.name"
        type="button"
        class="group flex-1 cursor-pointer border-0 bg-transparent py-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-gold"
        :aria-label="`Site ${i + 1} : ${site.name}`"
        :aria-current="i === index || undefined"
        @click="goTo(i)"
      >
        <span
          class="block h-0.5 overflow-hidden rounded-full bg-white/20 transition-colors group-hover:bg-white/40"
        >
          <span
            :key="i === index ? `active-${index}` : `idle-${i}`"
            class="block h-full w-full origin-left bg-gold"
            :class="fillClass(i)"
            @animationend="i === index && goTo(index + 1)"
          />
        </span>
      </button>
    </div>

    <div class="flex gap-3">
      <button type="button" :class="arrowClass" aria-label="Site précédent" @click="goTo(index - 1)">
        <LandingIcon name="chevron" class="size-icon-md rotate-180" />
      </button>
      <button type="button" :class="arrowClass" aria-label="Site suivant" @click="goTo(index + 1)">
        <LandingIcon name="chevron" class="size-icon-md" />
      </button>
    </div>
  </div>
</template>
