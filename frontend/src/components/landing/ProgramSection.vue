<script setup lang="ts">
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { EVENT_CONFIG } from './eventConfig'

const items = [
  {
    title: 'Patrimoine mondial',
    body: 'Expositions vidéo et photographiques des biens comoriens inscrits au patrimoine mondial de l’humanité.',
  },
  {
    title: 'Conférences & ateliers',
    body: 'Conférences-débats et ateliers animés par des professionnels et des diplômés de l’Unesco.',
  },
  {
    title: 'Danses traditionnelles',
    body: 'Danses traditionnelles comoriennes.',
  },
  {
    title: 'Gastronomie des îles',
    body: 'Gastronomie comorienne : saveurs et traditions de nos îles.',
  },
  {
    title: 'Concerts & spectacles',
    body: 'Concert et spectacles musicaux.',
  },
]

const days = [
  { label: 'SAMEDI', num: '16', year: 'JANVIER 2027' },
  { label: 'DIMANCHE', num: '17', year: 'JANVIER 2027' },
]
</script>

<template>
  <section id="programme" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-start gap-9 md:gap-split-md lg:grid-cols-program lg:gap-program">
      <div>
        <LandingEyebrow>AU PROGRAMME</LandingEyebrow>
        <LandingSectionTitle size="program">
          Venez célébrer<br />notre patrimoine !
        </LandingSectionTitle>
        <p class="mt-5 max-w-program text-sm leading-relaxed text-ink-muted">
          Un patrimoine à préserver, une histoire à transmettre, une culture à partager.
        </p>
        <p class="mt-4 text-eyebrow font-bold tracking-wide-label text-terracotta">
          {{ EVENT_CONFIG.dateShort.toUpperCase() }} · NANTERRE
        </p>
        <p class="mt-6 max-w-program text-nav leading-body text-ink-muted">
          Rendez-vous les 16 et 17 janvier 2027 à Nanterre pour deux jours de célébration, de découverte et de partage
          autour de la richesse du patrimoine comorien.
        </p>
        <LandingButton href="#ticket-form" class="mt-8">
          Réserver · dès 20 €
          <LandingIcon name="arrow" class="size-icon-sm" />
        </LandingButton>

        <div class="mt-10 grid gap-4 sm:grid-cols-2">
          <article
            v-for="day in days"
            :key="day.num"
            class="flex items-center gap-4 rounded-program bg-cream px-3.5 py-5 md:px-6 md:py-6"
          >
            <div class="flex min-w-20 flex-col items-center border-r border-line pr-3.5">
              <span class="text-2xs font-extrabold tracking-kicker text-terracotta">{{ day.label }}</span>
              <strong class="font-serif my-1 text-date-sm font-normal">{{ day.num }}</strong>
              <small class="whitespace-nowrap text-2xs text-ink-subtle">{{ day.year }}</small>
            </div>
            <p class="m-0 text-nav text-ink-muted">Journée découverte & partage</p>
          </article>
        </div>
      </div>

      <ol class="m-0 grid list-none gap-0 border-t border-line p-0">
        <li
          v-for="(item, index) in items"
          :key="item.title"
          class="grid grid-cols-step items-start gap-x-3.5 border-b border-line py-6"
        >
          <span
            class="row-span-2 grid size-7 place-items-center rounded-full border text-xs"
            :class="index === 0 ? 'border-gold bg-gold text-forest' : 'border-line-strong text-ink-muted'"
          >
            {{ String(index + 1).padStart(2, '0') }}
          </span>
          <strong class="inline-flex items-center gap-2 text-nav">
            {{ item.title }}
            <span
              v-if="index === 0"
              class="rounded-badge bg-forest px-2 py-0.5 text-2xs font-bold tracking-label text-gold"
            >
              TEMPS FORT
            </span>
          </strong>
          <p class="mt-1 text-nav leading-body text-ink-muted">{{ item.body }}</p>
        </li>
      </ol>
    </LandingContainer>
  </section>
</template>
