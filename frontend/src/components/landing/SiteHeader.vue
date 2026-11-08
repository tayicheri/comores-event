<script setup lang="ts">
import { ref } from 'vue'
import LandingBrand from './LandingBrand.vue'
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingIcon from './LandingIcon.vue'

const menuOpen = ref(false)

const links = [
  { href: '#evenement', label: 'L’événement' },
  { href: '#patrimoine', label: 'Le patrimoine' },
  { href: '#programme', label: 'Le programme' },
  { href: '#stand', label: 'Stand' },
  { href: '#historique', label: 'Édition 2024' },
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
        class="max-lg:absolute max-lg:top-full max-lg:right-0 max-lg:left-0 max-lg:flex-col max-lg:items-stretch max-lg:gap-0 max-lg:bg-canvas max-lg:px-6 max-lg:py-3 max-lg:shadow-nav lg:flex lg:items-center lg:gap-nav-md xl:gap-nav"
        :class="menuOpen ? 'flex' : 'hidden lg:flex'"
        aria-label="Navigation principale"
      >
        <a
          v-for="item in links"
          :key="item.href"
          :href="item.href"
          class="whitespace-nowrap text-nav font-medium hover:underline hover:underline-offset-4 max-lg:border-t max-lg:border-line max-lg:py-3.5"
          @click="closeMenu"
        >
          {{ item.label }}
        </a>
      </nav>

      <LandingButton href="#ticket-form" compact class="ml-auto hidden whitespace-nowrap sm:inline-flex lg:ml-0">
        <span class="md:hidden">Réserver</span>
        <span class="hidden md:inline">Réserver ma visite</span>
        <LandingIcon name="arrow" class="hidden size-icon-sm md:block" />
      </LandingButton>

      <button
        type="button"
        class="flex flex-col gap-1.5 border-0 bg-transparent p-2.5 lg:hidden"
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
