<script setup lang="ts">
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { EVENT_CONFIG } from './eventConfig'
</script>

<template>
  <section class="py-5 md:py-10 md:pb-11" aria-labelledby="hero-title">
    <LandingContainer
      class="grid grid-cols-1 items-center gap-8 md:min-h-hero-min md:grid-cols-hero md:gap-9 lg:gap-hero-gap"
    >
      <div class="pt-5 md:py-2.5">
        <div class="flex items-center gap-3 text-xs font-bold tracking-kicker text-terracotta md:text-caption">
          <span class="inline-block h-0.5 w-kicker-bar bg-gold" aria-hidden="true" />
          {{ EVENT_CONFIG.dateKicker }}
        </div>
        <LandingSectionTitle id="hero-title" as="h1" size="display" class="my-7 max-w-hero-copy md:my-10">
          Les Comores<br />à l’honneur.
        </LandingSectionTitle>
        <p class="mb-3.5 text-lede font-semibold leading-lede">
          Salon de la culture et de la<br class="hidden sm:inline" />
          gastronomie comorienne
        </p>
        <p class="mb-6 max-w-hero-desc text-sm leading-body text-ink-muted md:text-body">
          Deux journées pour découvrir, partager et célébrer. Toute la richesse d’une culture, le temps d’un week-end à
          Nanterre.
        </p>
        <div class="flex items-center gap-4 md:gap-5">
          <LandingButton href="#billetterie">
            Choisir mes billets
            <LandingIcon name="arrow" class="size-icon-sm" />
          </LandingButton>
          <a href="#evenement" class="text-nav font-bold underline underline-offset-4">Explorer le salon</a>
        </div>
        <p class="mt-5 text-caption text-ink-subtle">Dès 20 € / personne · Pass 2 jours à 30 €</p>
      </div>

      <div
        class="relative h-hero-img-sm max-w-hero-copy overflow-hidden rounded-hero-sm bg-media sm:h-hero-img-md sm:w-full sm:rounded-hero-md lg:h-hero-img lg:rounded-hero"
      >
        <img
          src="/landing/hero_moroni.jpg"
          alt="Vue panoramique de Moroni, ses maisons et le littoral comorien"
          class="size-full object-cover object-center saturate-75"
        />
        <div class="pointer-events-none absolute inset-x-0 top-1/3 bottom-0 bg-hero-fade" />
        <span
          class="absolute top-5 right-6 z-10 flex size-seal-sm flex-col items-center justify-center rounded-full bg-gold font-serif text-2xl leading-none text-forest md:size-seal md:text-seal"
        >
          2e
          <small class="mt-1.5 font-sans text-micro font-bold tracking-brand md:text-xs">ÉDITION</small>
        </span>
        <div class="absolute right-5 bottom-7 left-7 z-10 text-white">
          <p class="font-serif mb-1.5 text-subtitle font-normal italic leading-snug md:text-card">
            L’esprit des îles, le plaisir d’être ensemble.
          </p>
          <span class="text-xs tracking-eyebrow">{{ EVENT_CONFIG.dateOverlay }}</span>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
