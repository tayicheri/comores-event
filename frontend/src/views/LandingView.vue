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
      <UniversSection />
      <HeritageSection />
      <ProgramSection />
      <TicketSection @checkout="onCheckout" />
      <PracticalSection />
      <FaqSection />
    </main>
    <SiteFooter />
    <DemoCheckoutModal :open="modalOpen" :selection="selection" @close="modalOpen = false" />
  </div>
</template>
