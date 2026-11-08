<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
</script>

<template>
  <section id="evenement" class="pt-section-md pb-14 md:pt-section-lg md:pb-14 lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 lg:grid-cols-split lg:gap-split">
      <div>
        <LandingEyebrow>2E ÉDITION · LE PATRIMOINE MONDIAL À L’HONNEUR</LandingEyebrow>
        <LandingSectionTitle>Préserver, transmettre, partager.</LandingSectionTitle>
      </div>
      <div class="pt-2">
        <p class="mb-6 max-w-quote text-lede leading-lede">
          Un patrimoine à préserver. Une histoire à transmettre. Une culture à partager.
        </p>
        <p class="max-w-intro text-sm leading-loose text-ink-muted">
          Pour sa 2e édition, Comores Impact place au cœur du salon la célébration des biens comoriens inscrits au
          patrimoine mondial de l’humanité. Expositions, conférences, danses, saveurs et musique : deux jours pour
          découvrir ces lieux d’exception et l’histoire qu’ils portent, que les Comores soient vos racines ou votre
          prochaine découverte.
        </p>
      </div>
    </LandingContainer>
  </section>
</template>
