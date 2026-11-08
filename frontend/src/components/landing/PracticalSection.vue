<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { EVENT_CONFIG } from './eventConfig'

const mapsUrl = `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(EVENT_CONFIG.venue.address)}`
</script>

<template>
  <section id="infos" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 md:gap-split-md lg:grid-cols-split lg:gap-split">
      <div>
        <LandingEyebrow>INFORMATIONS PRATIQUES</LandingEyebrow>
        <LandingSectionTitle size="practical" class="mb-4">
          Les îles en esprit.<br />Nanterre en destination.
        </LandingSectionTitle>
        <p class="max-w-practical text-sm leading-loose text-ink-muted">
          L’immersion commence ici, en Île-de-France. Retrouvez-nous à la Salle des congrès pour deux journées consacrées
          à la culture et à la gastronomie comoriennes.
        </p>
        <p class="mt-6 flex items-center gap-1.5 text-eyebrow font-bold tracking-wide-label text-terracotta">
          <LandingIcon name="pin" class="size-icon-sm" />
          NANTERRE · FRANCE
        </p>
      </div>

      <div class="rounded-card bg-cream-soft p-5 md:p-panel">
        <div class="flex gap-4 border-b border-line pb-7">
          <span class="grid size-pin shrink-0 place-items-center rounded-full bg-forest text-gold" aria-hidden="true">
            <LandingIcon name="pin" class="size-icon-md" />
          </span>
          <div>
            <LandingSectionTitle as="h3" size="panel" class="mt-0.5 mb-2.5">
              {{ EVENT_CONFIG.venue.name }}
            </LandingSectionTitle>
            <p class="text-nav leading-body text-ink-muted">{{ EVENT_CONFIG.venue.address }}</p>
            <a
              :href="mapsUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="mt-2 inline-flex items-center gap-1 text-caption font-bold underline underline-offset-4"
            >
              Voir l’adresse sur la carte
              <LandingIcon name="arrow" class="size-icon-sm" />
            </a>
          </div>
        </div>
        <div class="flex gap-4 border-b border-line py-6">
          <span class="grid size-pin shrink-0 place-items-center rounded-full bg-hours text-forest" aria-hidden="true">
            <LandingIcon name="calendar" class="size-icon-md" />
          </span>
          <div>
            <strong class="text-nav">Samedi 16 et dimanche 17 janvier 2027</strong>
            <p class="mt-1.5 text-caption text-ink-subtle">
              Ouverture des stands à {{ EVENT_CONFIG.openingTime }} les deux jours ·
              <a href="#programme" class="font-bold underline underline-offset-4">Voir le programme complet</a>
            </p>
          </div>
        </div>
        <div class="flex gap-4 pt-6">
          <span class="grid size-pin shrink-0 place-items-center rounded-full bg-hours text-forest" aria-hidden="true">
            <LandingIcon name="info" class="size-icon-md" />
          </span>
          <div>
            <strong class="text-nav">Contact</strong>
            <p class="mt-1.5 text-nav leading-body text-ink-muted">
              <a
                v-for="(phone, i) in EVENT_CONFIG.venue.phones"
                :key="phone"
                :href="`tel:${phone.replace(/\s/g, '')}`"
                class="underline underline-offset-4"
              >
                {{ phone }}<template v-if="i < EVENT_CONFIG.venue.phones.length - 1"> / </template>
              </a>
            </p>
          </div>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
