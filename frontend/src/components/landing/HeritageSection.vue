<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
</script>

<template>
  <section class="bg-forest py-heritage-sm text-white md:py-heritage-sm lg:py-heritage">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 md:gap-split-md lg:grid-cols-heritage lg:gap-program">
      <div>
        <img
          src="/landing/heritage.jpg"
          alt="Scène artisanale évoquant les savoir-faire et les traditions comoriennes"
          class="h-heritage-img-sm w-full rounded-media-lg object-cover md:h-heritage-img-md lg:h-heritage-img"
        />
        <span class="mt-3 block text-xs text-heritage-muted">
          Aux Comores, les savoir-faire racontent notre histoire · Image d’illustration
        </span>
      </div>
      <div>
        <LandingEyebrow tone="gold">PLUS QU’UN RENDEZ-VOUS</LandingEyebrow>
        <LandingSectionTitle size="heritage" class="mb-4 max-w-intro">
          Un héritage vivant, une histoire à partager.
        </LandingSectionTitle>
        <p class="max-w-intro text-sm leading-heritage text-heritage-body">
          Le patrimoine se raconte dans les gestes, les matières et les souvenirs. Il se transmet à table, en musique,
          dans les échanges entre générations.
        </p>
        <p class="mt-4 max-w-intro text-sm leading-heritage text-heritage-body">
          À Nanterre, faisons vivre ces liens. Que les Comores soient votre histoire ou votre prochaine découverte.
        </p>
        <div class="mt-8 flex flex-wrap gap-7 text-nav">
          <span class="inline-flex items-center gap-2">
            <LandingIcon name="star" class="size-icon-sm text-gold" />
            Transmettre les savoir-faire
          </span>
          <span class="inline-flex items-center gap-2">
            <LandingIcon name="star" class="size-icon-sm text-gold" />
            Faire place à la rencontre
          </span>
        </div>
      </div>
    </LandingContainer>
  </section>
</template>
