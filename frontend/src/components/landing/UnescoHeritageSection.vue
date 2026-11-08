<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import UnescoSiteCard from './UnescoSiteCard.vue'

type HeritageSite = {
  src: string
  alt: string
  name: string
  place: string
  layout: string
}

const sites: HeritageSite[] = [
  {
    src: '/landing/patrimoine/mosquee-vendredi-moroni.jpg',
    alt: 'Ancienne mosquée du Vendredi de Moroni, blanche à arcades, au bord du port',
    name: 'Ancienne mosquée du Vendredi',
    place: 'MORONI · GRANDE COMORE',
    layout: 'sm:col-span-2 lg:row-span-2',
  },
  {
    src: '/landing/patrimoine/porte-medina.jpg',
    alt: 'Porte blanche à arcades ouvrant sur une ruelle de médina',
    name: 'Portes et ruelles des médinas',
    place: 'VILLES ANCIENNES',
    layout: 'row-span-2',
  },
  {
    src: '/landing/patrimoine/plafond-sculpte.jpg',
    alt: 'Plafond en bois sculpté et peint aux motifs géométriques',
    name: 'Plafonds de bois sculpté',
    place: 'SAVOIR-FAIRE',
    layout: '',
  },
  {
    src: '/landing/patrimoine/ruines-arcades.jpg',
    alt: 'Ruines à arcades en pierre volcanique face à l’océan',
    name: 'Arcades de pierre volcanique',
    place: 'BORD DE MER',
    layout: '',
  },
  {
    src: '/landing/patrimoine/ruines-palais.jpg',
    alt: 'Vestiges d’un palais en pierre volcanique avec son escalier',
    name: 'Palais et demeures anciennes',
    place: 'ARCHITECTURE ANCIENNE',
    layout: 'lg:col-span-2',
  },
  {
    src: '/landing/patrimoine/citadelle.jpg',
    alt: 'Vue aérienne des remparts d’une citadelle surplombant la ville et l’océan',
    name: 'Citadelles et remparts',
    place: 'FORTIFICATIONS',
    layout: 'lg:col-span-2',
  },
]
</script>

<template>
  <section
    id="patrimoine"
    class="bg-forest py-heritage-sm text-white lg:py-heritage"
    aria-labelledby="patrimoine-title"
  >
    <LandingContainer>
      <div class="mb-10 grid grid-cols-1 items-end gap-6 lg:mb-12 lg:grid-cols-split lg:gap-split">
        <div>
          <LandingEyebrow tone="gold" class="flex items-center gap-3">
            <img src="/landing/logo-unesco.svg" alt="" class="h-6 w-auto" />
            PATRIMOINE MONDIAL DE L’HUMANITÉ
          </LandingEyebrow>
          <LandingSectionTitle id="patrimoine-title">
            Nos trésors, patrimoine de l’humanité.
          </LandingSectionTitle>
        </div>
        <div>
          <p class="max-w-intro text-sm leading-heritage text-heritage-body">
            Cœur de cette 2e édition, les biens comoriens inscrits au patrimoine mondial de l’humanité seront à
            l’honneur à travers des expositions vidéo et photographiques, des conférences-débats et des ateliers animés
            par des professionnels et des diplômés de l’Unesco.
          </p>
          <p class="mt-5 inline-flex items-center gap-2 text-nav">
            <LandingIcon name="star" class="size-icon-sm text-gold" />
            Venez célébrer notre patrimoine !
          </p>
        </div>
      </div>

      <div
        class="grid grid-cols-1 auto-rows-(--spacing-unesco-row) gap-4 sm:grid-cols-2 lg:grid-cols-4 lg:auto-rows-(--spacing-unesco-row-lg) lg:gap-5"
      >
        <UnescoSiteCard
          v-for="site in sites"
          :key="site.src"
          :src="site.src"
          :alt="site.alt"
          :name="site.name"
          :place="site.place"
          :class="site.layout"
        />
      </div>
      <p class="mt-5 text-xs text-heritage-muted">
        Images des sites et savoir-faire présentés lors des expositions du salon.
      </p>
    </LandingContainer>
  </section>
</template>
