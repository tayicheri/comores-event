<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const days = [
  {
    label: 'SAMEDI',
    num: '16',
    year: 'JANVIER 2027',
    title: 'Le plaisir de la découverte',
    body: 'Salon, découvertes culturelles et gastronomiques : une première journée pour explorer les univers comoriens.',
  },
  {
    label: 'DIMANCHE',
    num: '17',
    year: 'JANVIER 2027',
    title: 'Le bonheur de se retrouver',
    body: 'Rencontres, animations et découvertes : prolonger les échanges et partager de nouvelles histoires.',
  },
]
</script>

<template>
  <section id="programme" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 md:gap-split-md lg:grid-cols-program lg:gap-program">
      <div>
        <LandingEyebrow>IL Y A 2 JOURNÉES</LandingEyebrow>
        <LandingSectionTitle size="program">
          Deux jours. Une<br />culture à<br />partager.
        </LandingSectionTitle>
        <p class="mt-5 max-w-program text-sm leading-relaxed text-ink-muted">
          Venez pour une journée de découvertes, ou prenez le temps de vivre tout le week-end.
        </p>
        <p class="mt-4 flex items-start gap-2 text-caption leading-relaxed text-ink-muted">
          <LandingIcon name="info" class="mt-0.5 size-icon-md text-terracotta" />
          Les horaires détaillés et les intervenants seront annoncés prochainement.
        </p>
      </div>

      <div class="grid gap-4">
        <article
          v-for="day in days"
          :key="day.num"
          class="flex items-center gap-4 rounded-program bg-cream px-3.5 py-5 md:gap-8 md:px-9 md:py-7"
        >
          <div class="flex min-w-20 flex-col items-center border-r border-line pr-3.5 md:min-w-24 md:pr-6">
            <span class="text-2xs font-extrabold tracking-kicker text-terracotta">{{ day.label }}</span>
            <strong class="font-serif my-1 text-date-sm font-normal md:text-date">{{ day.num }}</strong>
            <small class="whitespace-nowrap text-2xs text-ink-subtle">{{ day.year }}</small>
          </div>
          <div>
            <LandingSectionTitle as="h3" size="subtitle" class="mb-2">{{ day.title }}</LandingSectionTitle>
            <p class="m-0 text-eyebrow leading-body text-ink-muted md:text-nav">{{ day.body }}</p>
          </div>
        </article>
      </div>
    </LandingContainer>
  </section>
</template>
