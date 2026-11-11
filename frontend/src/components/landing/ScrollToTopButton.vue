<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import LandingIcon from './LandingIcon.vue'

/** Affiche le bouton après ~1 viewport (hero + début de page). */
const SHOW_AFTER_PX = () => Math.min(window.innerHeight, 900)

const visible = ref(false)

function onScroll() {
  visible.value = window.scrollY >= SHOW_AFTER_PX()
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<template>
  <button
    type="button"
    class="fixed right-5 bottom-5 z-40 flex size-btn-min items-center justify-center rounded-full bg-forest text-white shadow-nav transition duration-200 hover:bg-forest-hover hover:-translate-y-0.5 focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-gold md:right-8 md:bottom-8"
    :class="visible ? 'opacity-100 translate-y-0' : 'pointer-events-none opacity-0 translate-y-2'"
    :tabindex="visible ? 0 : -1"
    :aria-hidden="!visible"
    aria-label="Remonter en haut"
    @click="scrollToTop"
  >
    <LandingIcon name="chevron" class="size-icon-md -rotate-90" />
  </button>
</template>
