<script setup lang="ts">
import { ref } from 'vue'
import DemoCheckoutModal from '@/components/landing/DemoCheckoutModal.vue'
import FaqSection from '@/components/landing/FaqSection.vue'
import HeritageSection from '@/components/landing/HeritageSection.vue'
import HeroSection from '@/components/landing/HeroSection.vue'
import InfoStrip from '@/components/landing/InfoStrip.vue'
import IntroSection from '@/components/landing/IntroSection.vue'
import LandingIconSprite from '@/components/landing/LandingIconSprite.vue'
import PracticalSection from '@/components/landing/PracticalSection.vue'
import ProgramSection from '@/components/landing/ProgramSection.vue'
import SiteFooter from '@/components/landing/SiteFooter.vue'
import SiteHeader from '@/components/landing/SiteHeader.vue'
import TicketSection from '@/components/landing/TicketSection.vue'
import UniversSection from '@/components/landing/UniversSection.vue'
import { EVENT_CONFIG, formatEuro, type TicketDay, type TicketFormula } from '@/components/landing/eventConfig'

const modalOpen = ref(false)
const modalLines = ref<string[]>([])

function onCheckout(payload: {
  formula: TicketFormula
  day: TicketDay
  quantity: number
  unit: number
  total: number
}) {
  const formulaLabel =
    payload.formula === 'day'
      ? `Billet 1 journée — ${EVENT_CONFIG.dates[payload.day]}`
      : 'Pass 2 jours — vendredi et samedi'
  modalLines.value = [
    formulaLabel,
    `${payload.quantity} billet${payload.quantity > 1 ? 's' : ''} × ${formatEuro(payload.unit)}`,
    `Total : ${formatEuro(payload.total)}`,
  ]
  modalOpen.value = true
}
</script>

<template>
  <div class="bg-canvas text-ink">
    <LandingIconSprite />
    <SiteHeader />
    <main>
      <HeroSection />
      <InfoStrip />
      <IntroSection />
      <UniversSection />
      <HeritageSection />
      <ProgramSection />
      <TicketSection @checkout="onCheckout" />
      <PracticalSection />
      <FaqSection />
    </main>
    <SiteFooter />
    <DemoCheckoutModal :open="modalOpen" :lines="modalLines" @close="modalOpen = false" />
  </div>
</template>
