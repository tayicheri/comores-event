<script setup lang="ts">
import { computed, ref } from 'vue'
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import ProgramDayTimeline from './ProgramDayTimeline.vue'
import { EVENT_CONFIG, formatEuro } from './eventConfig'
import { PROGRAM_DAYS, type ProgramDay } from './programme'

const activeKey = ref<ProgramDay['key']>(PROGRAM_DAYS[0]!.key)
const activeDay = computed(() => PROGRAM_DAYS.find((day) => day.key === activeKey.value) ?? PROGRAM_DAYS[0]!)

const dayTagline: Record<ProgramDay['key'], string> = {
  samedi: 'Ouverture officielle, patrimoine, traditions et concert',
  dimanche: 'Ateliers, table ronde, danses et clôture',
}
</script>

<template>
  <section id="programme" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-start gap-9 md:gap-split-md lg:grid-cols-program lg:gap-program">
      <div class="lg:sticky lg:top-8">
        <LandingEyebrow>LE PROGRAMME</LandingEyebrow>
        <LandingSectionTitle size="program">
          Venez célébrer<br />notre patrimoine !
        </LandingSectionTitle>
        <p class="mt-5 max-w-program text-sm leading-relaxed text-ink-muted">
          Conférences, rencontres, animations culturelles, ateliers et concerts : le programme complet des deux
          journées.
        </p>
        <p class="mt-4 text-eyebrow font-bold tracking-wide-label text-terracotta">
          {{ EVENT_CONFIG.dateShort.toUpperCase() }} · NANTERRE · DÈS {{ EVENT_CONFIG.openingTime.toUpperCase() }}
        </p>

        <div class="mt-8 grid gap-4 sm:grid-cols-2" role="tablist" aria-label="Choisir la journée du programme">
          <button
            v-for="day in PROGRAM_DAYS"
            :id="`programme-tab-${day.key}`"
            :key="day.key"
            type="button"
            role="tab"
            :aria-selected="activeKey === day.key"
            :aria-controls="`programme-panel-${day.key}`"
            class="flex items-center gap-4 rounded-program border-2 px-3.5 py-5 text-left transition-colors md:px-6 md:py-6"
            :class="activeKey === day.key ? 'border-forest bg-cream' : 'border-transparent bg-cream/60 hover:bg-cream'"
            @click="activeKey = day.key"
          >
            <span class="flex min-w-20 flex-col items-center border-r border-line pr-3.5">
              <span class="text-2xs font-extrabold tracking-kicker text-terracotta">{{ day.label }}</span>
              <strong class="font-serif my-1 text-date-sm font-normal">{{ day.num }}</strong>
              <small class="whitespace-nowrap text-2xs text-ink-subtle">{{ day.month }}</small>
            </span>
            <span class="text-nav text-ink-muted">{{ dayTagline[day.key] }}</span>
          </button>
        </div>

        <LandingButton href="#ticket-form" class="mt-8">
          Réserver · dès {{ formatEuro(EVENT_CONFIG.prices.day) }}
          <LandingIcon name="arrow" class="size-icon-sm" />
        </LandingButton>
      </div>

      <div
        :id="`programme-panel-${activeDay.key}`"
        role="tabpanel"
        :aria-labelledby="`programme-tab-${activeDay.key}`"
      >
        <ProgramDayTimeline :day="activeDay" />
      </div>
    </LandingContainer>
  </section>
</template>
