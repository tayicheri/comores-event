<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import UnescoSiteCard from './UnescoSiteCard.vue'

type Highlight = {
  key: string
  /** `null` tant que la photo n’a pas été fournie. */
  src: string | null
  alt: string
  label: string
  title: string
  body: string
}

const highlights: Highlight[] = [
  {
    key: 'culture',
    src: '/landing/culture/culture-danse.jpg',
    alt: 'Jeunes danseurs comoriens en polos bleus, pagnes rayés et kofia, lors d’une prestation en place publique',
    label: 'CULTURE',
    title: 'Traditions vivantes',
    body: 'Zifafa, Sambe, Chigoma, Tari, Wadaha : les danses et cérémonies qui rythment la vie des îles.',
  },
  {
    key: 'gastronomie',
    src: '/landing/culture/gastronomie-samboussas.jpg',
    alt: 'Samboussas comoriens dorés dans un bol en coco, entourés de galettes et d’épices',
    label: 'GASTRONOMIE',
    title: 'Saveurs des îles',
    body: 'Ateliers culinaires et dégustation d’anciens plats traditionnels, recettes transmises de génération en génération.',
  },
  {
    key: 'artisanat',
    src: '/landing/culture/artisanat-ilatu.jpg',
    alt: 'Ilatu, sandales traditionnelles comoriennes colorées exposées sur un étal',
    label: 'ARTISANAT',
    title: 'Savoir-faire comoriens',
    body: 'Des Ilatu, chaussures traditionnelles, à la vannerie : l’artisanat qui fait rayonner la culture comorienne.',
  },
]
</script>

<template>
  <section id="culture" class="py-section-md md:py-section-lg lg:py-section" aria-labelledby="culture-title">
    <LandingContainer>
      <div class="mb-10 max-w-intro lg:mb-12">
        <LandingEyebrow>CULTURE · GASTRONOMIE · ARTISANAT</LandingEyebrow>
        <LandingSectionTitle id="culture-title">L’art de vivre comorien.</LandingSectionTitle>
      </div>
      <ul class="m-0 grid list-none grid-cols-1 gap-6 p-0 md:grid-cols-3 lg:gap-8">
        <li v-for="item in highlights" :key="item.key">
          <UnescoSiteCard
            :src="item.src"
            :alt="item.alt"
            :name="item.title"
            :place="item.label"
            class="h-unesco-row lg:h-unesco-row-lg"
          />
          <p class="mt-4 text-nav leading-body text-ink-muted">{{ item.body }}</p>
        </li>
      </ul>
    </LandingContainer>
  </section>
</template>
