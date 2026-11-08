<script setup lang="ts">
import LandingBrand from './LandingBrand.vue'
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const links = [
  { href: '#evenement', label: 'L’événement' },
  { href: '#patrimoine', label: 'Le patrimoine' },
  { href: '#programme', label: 'Le programme' },
  { href: '#billetterie', label: 'Billetterie' },
  { href: '#stand', label: 'Stand' },
  { href: '#infos', label: 'Infos pratiques' },
] as const
</script>

<template>
  <footer class="relative bg-forest-deep bg-footer-pattern bg-size-footer bg-top pt-footer-y pb-6 text-white">
    <LandingContainer class="relative z-10">
      <div
        class="flex flex-col items-start justify-between gap-7 border-b border-strip-line pb-16 md:flex-row md:items-center"
      >
        <div>
          <LandingEyebrow tone="corail">SALON DE LA CULTURE ET DE LA GASTRONOMIE COMORIENNE · 2027</LandingEyebrow>
          <LandingSectionTitle size="footer" class="my-3.5">On se retrouve à Nanterre ?</LandingSectionTitle>
          <p class="text-nav text-heritage-muted">16 &amp; 17 janvier 2027 · Un patrimoine mondial à célébrer.</p>
        </div>
        <LandingButton href="#ticket-form" variant="yellow">
          Choisir mes billets
          <LandingIcon name="arrow" class="size-icon-sm" />
        </LandingButton>
      </div>

      <div class="flex flex-col items-start justify-between gap-5 py-10 md:flex-row md:items-center">
        <LandingBrand tone="footer" />
        <div class="flex flex-wrap gap-3.5 md:gap-7">
          <a
            v-for="link in links"
            :key="link.href"
            :href="link.href"
            class="text-nav font-medium hover:underline hover:underline-offset-4"
          >
            {{ link.label }}
          </a>
        </div>
      </div>

      <div class="flex flex-col justify-between gap-2 pt-4 text-eyebrow text-footer-muted md:flex-row md:gap-5">
        <span>© 2027 Salon de la culture et de la gastronomie comorienne en France</span>
        <span>Site réalisé par Foodyfly</span>
      </div>
    </LandingContainer>
  </footer>
</template>
