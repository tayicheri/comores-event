<script setup lang="ts">
import { ref } from 'vue'
import CultureHighlightsSection from '@/components/landing/CultureHighlightsSection.vue'
import DemoCheckoutModal from '@/components/landing/DemoCheckoutModal.vue'
import FaqSection from '@/components/landing/FaqSection.vue'
import HeroSection from '@/components/landing/HeroSection.vue'
import HistorySection from '@/components/landing/HistorySection.vue'
import InfoStrip from '@/components/landing/InfoStrip.vue'
import IntroSection from '@/components/landing/IntroSection.vue'
import LandingIconSprite from '@/components/landing/LandingIconSprite.vue'
import PracticalSection from '@/components/landing/PracticalSection.vue'
import ProgramSection from '@/components/landing/ProgramSection.vue'
import SiteFooter from '@/components/landing/SiteFooter.vue'
import SiteHeader from '@/components/landing/SiteHeader.vue'
import SponsorsSection from '@/components/landing/SponsorsSection.vue'
import StandSection from '@/components/landing/StandSection.vue'
import TicketSection from '@/components/landing/TicketSection.vue'
import UnescoHeritageSection from '@/components/landing/UnescoHeritageSection.vue'
import type { CheckoutSelection } from '@/types/payment'

const modalOpen = ref(false)
const selection = ref<CheckoutSelection | null>(null)

function onCheckout(payload: CheckoutSelection) {
  selection.value = payload
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
      <SponsorsSection />
      <UnescoHeritageSection />
      <CultureHighlightsSection />
      <ProgramSection />
      <TicketSection @checkout="onCheckout" />
      <StandSection />
      <HistorySection />
      <PracticalSection />
      <FaqSection />
    </main>
    <SiteFooter />
    <DemoCheckoutModal :open="modalOpen" :selection="selection" @close="modalOpen = false" />
  </div>
</template>
