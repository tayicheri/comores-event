<script setup lang="ts">
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { EVENT_CONFIG, formatEuro } from './eventConfig'

const offers = [
  { key: 'pass', title: 'Pass 2 jours', caption: 'Samedi 16 et dimanche 17 janvier', price: EVENT_CONFIG.standPrices.pass },
  { key: 'day', title: '1 journée', caption: 'Samedi ou dimanche, au choix', price: EVENT_CONFIG.standPrices.day },
] as const

const telHref = (phone: string) => `tel:${phone.replace(/\s/g, '')}`
</script>

<template>
  <section id="stand" class="bg-forest py-section-md text-white md:py-section-lg lg:py-section" aria-labelledby="stand-title">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 md:gap-split-md lg:grid-cols-split lg:gap-split">
      <div>
        <LandingEyebrow tone="gold">EXPOSANTS · RÉSERVATION DE STAND</LandingEyebrow>
        <LandingSectionTitle id="stand-title" class="mb-5">Tenez votre stand au Salon.</LandingSectionTitle>
        <p class="max-w-intro text-sm leading-loose text-heritage-body">
          Artisans, restaurateurs, créateurs, associations, entrepreneurs : présentez vos savoir-faire et vos produits
          aux visiteurs du Salon de la culture et de la gastronomie comorienne en France. Réservez votre emplacement
          en contactant directement l’organisation.
        </p>
      </div>

      <div class="rounded-card bg-cream-soft p-5 text-ink md:p-panel">
        <ul class="m-0 grid list-none grid-cols-2 gap-2.5 p-0 md:gap-3.5" aria-label="Tarifs des stands">
          <li
            v-for="offer in offers"
            :key="offer.key"
            class="flex min-h-32 flex-col rounded-media border-2 p-3.5 md:min-h-36 md:p-5"
            :class="offer.key === 'pass' ? 'border-forest bg-option' : 'border-line-option bg-white'"
          >
            <span class="text-caption font-bold md:text-body">{{ offer.title }}</span>
            <span class="font-serif my-2.5 text-panel font-medium md:my-4 md:text-price">
              {{ formatEuro(offer.price) }}
              <small class="font-sans text-caption font-normal text-ink-subtle">/ stand</small>
            </span>
            <span class="text-xs text-ink-subtle md:text-eyebrow">{{ offer.caption }}</span>
          </li>
        </ul>

        <div class="mt-6 border-t border-line pt-6">
          <strong class="block text-nav">Réserver un stand</strong>
          <p class="mt-1 text-eyebrow text-ink-subtle">Contactez-nous par téléphone :</p>
          <div class="mt-4 grid gap-2.5 sm:grid-cols-2">
            <LandingButton v-for="phone in EVENT_CONFIG.venue.phones" :key="phone" :href="telHref(phone)">
              <LandingIcon name="phone" class="size-icon-sm text-gold" />
              {{ phone }}
            </LandingButton>
          </div>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
