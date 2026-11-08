<script setup lang="ts">
import LandingIcon from './LandingIcon.vue'
import { EVENT_CONFIG } from './eventConfig'
</script>

<template>
  <div>
    <p class="text-eyebrow font-bold tracking-eyebrow text-ochre">SIX VILLES DES COMORES À L’HONNEUR</p>
    <ul class="mt-4 flex list-none flex-wrap gap-2.5 p-0" aria-label="Les six villes célébrées par le Salon">
      <li
        v-for="city in EVENT_CONFIG.cities"
        :key="city"
        class="inline-flex items-center gap-1.5 rounded-badge border border-forest/20 bg-cream px-3.5 py-2 text-nav font-semibold text-forest"
      >
        <LandingIcon name="star" class="size-icon-sm text-gold" />
        {{ city }}
      </li>
    </ul>
  </div>
</template>
