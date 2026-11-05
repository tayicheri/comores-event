<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingIcon from './LandingIcon.vue'
import { EVENT_CONFIG } from './eventConfig'

type IconName = 'calendar' | 'pin' | 'ticket'

const items: Array<{
  icon: IconName
  title: string
  body?: string
  accentPrices?: boolean
}> = [
  {
    icon: 'calendar',
    title: 'UN WEEK-END POUR SE RETROUVER',
    body: 'Ven. 15 & sam. 16 janvier 2027',
  },
  {
    icon: 'pin',
    title: 'RENDEZ-VOUS EN ÎLE-DE-FRANCE',
    body: 'Salle des congrès · Nanterre, France',
  },
  {
    icon: 'ticket',
    title: 'À CHACUN SA VISITE',
    accentPrices: true,
  },
]
</script>

<template>
  <section class="bg-forest text-canvas" aria-label="Informations essentielles">
    <LandingContainer
      class="grid grid-cols-1 items-center py-4 md:min-h-strip md:grid-cols-strip md:py-0"
    >
      <div
        v-for="(item, index) in items"
        :key="item.title"
        class="flex items-center gap-2.5 border-strip-line py-3 md:gap-4 md:border-0 md:py-0 md:pr-6"
        :class="index > 0 ? 'border-t md:border-t-0 md:border-l md:pl-5 lg:pl-12' : ''"
      >
        <span
          class="grid size-icon-ring shrink-0 place-items-center rounded-full bg-gold/20 text-gold shadow-gold-ring"
          aria-hidden="true"
        >
          <LandingIcon :name="item.icon" class="size-icon-md md:size-icon-lg" />
        </span>
        <div class="min-w-0">
          <strong class="block text-xs font-semibold tracking-label text-gold-muted">{{ item.title }}</strong>
          <p v-if="item.accentPrices" class="mt-2 text-nav">
            <span class="text-gold">{{ EVENT_CONFIG.prices.day }} €</span> la journée ·
            <span class="text-gold">{{ EVENT_CONFIG.prices.pass }} €</span> les 2 jours
          </p>
          <p v-else class="mt-2 text-nav">{{ item.body }}</p>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
