<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const sponsors = [
  {
    name: 'UNESCO',
    src: '/landing/logo-unesco.svg',
    alt: 'Logo de l’UNESCO',
  },
  {
    name: 'Ville de Nanterre',
    src: '/landing/logo-ville-nanterre.svg',
    alt: 'Logo de la Ville de Nanterre',
  },
  {
    name: 'Comores Impact',
    src: '/landing/logo-comores-impact.png',
    alt: 'Logo de Comores Impact',
  },
] as const
</script>

<template>
  <section id="partenaires" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer>
      <div class="mb-10 text-center md:mb-12">
        <LandingEyebrow>ILS NOUS SOUTIENNENT</LandingEyebrow>
        <LandingSectionTitle class="mt-2">Nos partenaires.</LandingSectionTitle>
      </div>

      <ul
        class="m-0 flex list-none flex-col items-center justify-center gap-10 p-0 sm:flex-row sm:flex-wrap sm:gap-12 md:gap-16"
        aria-label="Partenaires de l’événement"
      >
        <li v-for="sponsor in sponsors" :key="sponsor.name" class="flex items-center justify-center">
          <img
            :src="sponsor.src"
            :alt="sponsor.alt"
            class="h-12 w-auto max-w-[10rem] object-contain md:h-14 md:max-w-[12rem]"
            loading="lazy"
            decoding="async"
          />
        </li>
      </ul>
    </LandingContainer>
  </section>
</template>
