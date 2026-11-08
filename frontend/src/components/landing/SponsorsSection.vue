<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { PARTNER_GROUPS } from './partners'
</script>

<template>
  <section id="partenaires" class="bg-pale py-section-md md:py-section-lg lg:py-section">
    <LandingContainer>
      <div class="mb-10 text-center md:mb-12">
        <LandingEyebrow>ILS NOUS SOUTIENNENT</LandingEyebrow>
        <LandingSectionTitle class="mt-2">Organisateurs et partenaires.</LandingSectionTitle>
      </div>

      <div class="grid gap-10 md:gap-12">
        <div v-for="group in PARTNER_GROUPS" :key="group.title">
          <h3 class="mb-5 text-center text-eyebrow font-bold tracking-eyebrow text-terracotta">
            {{ group.title.toUpperCase() }}
          </h3>
          <ul
            class="m-0 flex list-none flex-wrap items-stretch justify-center gap-4 p-0 md:gap-5"
            :aria-label="group.title"
          >
            <li
              v-for="partner in group.partners"
              :key="partner.src"
              class="flex h-24 w-36 items-center justify-center rounded-media bg-white p-3 shadow-card md:h-28 md:w-44"
            >
              <img
                :src="partner.src"
                :alt="`Logo ${partner.name}`"
                :title="partner.name"
                class="max-h-full max-w-full object-contain"
                loading="lazy"
                decoding="async"
              />
            </li>
          </ul>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
