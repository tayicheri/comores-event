<script setup lang="ts">
import { ref } from 'vue'
import LandingBrand from './LandingBrand.vue'
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingIcon from './LandingIcon.vue'

const menuOpen = ref(false)

const links = [
  { href: '#evenement', label: 'L’événement' },
  { href: '#univers', label: 'Les univers' },
  { href: '#programme', label: 'Le programme' },
  { href: '#infos', label: 'Infos pratiques' },
] as const

function closeMenu() {
  menuOpen.value = false
}

function toggleMenu() {
  menuOpen.value = !menuOpen.value
}
</script>

<template>
  <header id="haut" class="relative z-30 h-header-sm bg-canvas md:h-header">
    <LandingContainer class="flex h-full items-center justify-between gap-4 md:gap-5">
      <LandingBrand />

      <nav
        id="navigation"
        class="max-md:absolute max-md:top-full max-md:right-0 max-md:left-0 max-md:flex-col max-md:items-stretch max-md:gap-0 max-md:bg-canvas max-md:px-6 max-md:py-3 max-md:shadow-nav md:flex md:items-center md:gap-nav-md lg:gap-nav"
        :class="menuOpen ? 'flex' : 'hidden md:flex'"
        aria-label="Navigation principale"
      >
        <a
          v-for="item in links"
          :key="item.href"
          :href="item.href"
          class="text-nav font-medium hover:underline hover:underline-offset-4 max-md:border-t max-md:border-line max-md:py-3.5"
          @click="closeMenu"
        >
          {{ item.label }}
        </a>
      </nav>

      <LandingButton href="#billetterie" compact class="ml-auto hidden sm:inline-flex md:ml-0">
        Réserver ma visite
        <LandingIcon name="arrow" class="size-icon-sm" />
      </LandingButton>

      <button
        type="button"
        class="flex flex-col gap-1.5 border-0 bg-transparent p-2.5 md:hidden"
        :aria-expanded="menuOpen"
        aria-controls="navigation"
        :aria-label="menuOpen ? 'Fermer le menu' : 'Ouvrir le menu'"
        @click="toggleMenu"
      >
        <span class="block h-0.5 w-kicker-bar bg-forest" />
        <span class="block h-0.5 w-kicker-bar bg-forest" />
        <span class="block h-0.5 w-kicker-bar bg-forest" />
      </button>
    </LandingContainer>
  </header>
</template>
