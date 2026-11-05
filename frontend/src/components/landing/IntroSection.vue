<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
</script>

<template>
  <section id="evenement" class="pt-section-md pb-14 md:pt-section-lg md:pb-14 lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 lg:grid-cols-split lg:gap-split">
      <div>
        <LandingEyebrow>UN ARCHIPEL, MILLE FAÇONS DE PARTAGER</LandingEyebrow>
        <LandingSectionTitle>La culture nous unit.</LandingSectionTitle>
      </div>
      <div class="pt-2">
        <p class="mb-6 max-w-quote text-lede leading-lede">
          Une recette que l’on transmet. Un geste que l’on découvre. Une conversation qui rapproche.
        </p>
        <p class="max-w-intro text-sm leading-loose text-ink-muted">
          Pour sa 2e édition, Comores Events vous invite à explorer le patrimoine comorien à travers ses couleurs, ses
          savoir-faire et ses rencontres. Un rendez-vous ouvert aux curieux comme à celles et ceux qui souhaitent renouer
          avec leurs racines.
        </p>
      </div>
    </LandingContainer>
  </section>
</template>
