<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import UnescoSiteCard from './UnescoSiteCard.vue'
import UnescoSitesSlideshow from './UnescoSitesSlideshow.vue'
import { HERITAGE_SITES } from './heritageSites'
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
          <p class="max-w-intro text-base leading-heritage text-heritage-body">
            Cœur de cette 2e édition, les biens comoriens inscrits au patrimoine mondial de l’humanité seront à
            l’honneur à travers des expositions vidéo et photographiques, des conférences-débats et des ateliers animés
            par des spécialistes du patrimoine et des représentants de l’UNESCO.
          </p>
          <p class="mt-5 inline-flex items-center gap-2 text-sm">
            <LandingIcon name="star" class="size-icon-sm text-gold" />
            Venez célébrer notre patrimoine !
          </p>
        </div>
      </div>

      <UnescoSitesSlideshow class="sm:hidden" :sites="HERITAGE_SITES" />

      <div
        class="hidden grid-flow-row-dense auto-rows-(--spacing-unesco-row) gap-4 sm:grid sm:grid-cols-2 lg:grid-cols-4 lg:auto-rows-(--spacing-unesco-row-lg) lg:gap-5"
      >
        <UnescoSiteCard
          v-for="site in HERITAGE_SITES"
          :key="site.name"
          :src="site.src"
          :alt="site.alt"
          :name="site.name"
          :place="`${site.city.toUpperCase()} · ${site.island.toUpperCase()}`"
          :class="site.layout"
        />
      </div>
      <p class="mt-5 text-sm text-heritage-muted">
        Sites patrimoniaux de Moroni, Iconi, Itsandra, Domoni, Mutsamudu et Ntsoudjini présentés lors des expositions
        du Salon.
      </p>
    </LandingContainer>
  </section>
</template>
