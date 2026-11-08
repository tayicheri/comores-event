<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import PhotoSlideshow from './PhotoSlideshow.vue'
import type { SlideshowSlide } from './slideshow'

const slides: SlideshowSlide[] = [
  {
    src: '/landing/salon-2024/01.jpg',
    alt: 'Stand Baoperry au salon 2024 : mugs, stickers et créations illustrées inspirées des Comores',
    caption: 'Baoperry — idées cadeaux et déco comorienne au style minimaliste.',
  },
  {
    src: '/landing/salon-2024/02.jpg',
    alt: 'Table de livres et visiteurs au salon de la culture comorienne 2024',
    caption: 'La librairie des auteurs comoriens.',
  },
  {
    src: '/landing/salon-2024/03.jpg',
    alt: 'Rencontre devant le stand d’une agence immobilière comorienne au salon 2024',
    caption: 'Entrepreneurs et porteurs de projets.',
  },
  {
    src: '/landing/salon-2024/04.jpg',
    alt: 'Deux exposantes derrière leur stand orné du drapeau des Comores',
    caption: 'Des stands aux couleurs de l’archipel.',
  },
  {
    src: '/landing/salon-2024/05.jpg',
    alt: 'Stand Ulanga Arts Création : artisanat en bois et créations comoriennes',
    caption: 'L’artisanat d’art d’Ulanga Arts.',
  },
  {
    src: '/landing/salon-2024/06.jpg',
    alt: 'Échanges autour d’une table de livres lors du salon 2024',
    caption: 'Des échanges entre générations.',
  },
  {
    src: '/landing/salon-2024/07.jpg',
    alt: 'Présentation de jeux traditionnels et vannerie au salon 2024',
    caption: 'Mraha et vannerie traditionnelle.',
  },
]
</script>

<template>
  <section id="historique" class="overflow-hidden py-section-md md:py-section-lg lg:py-section">
    <LandingContainer>
      <PhotoSlideshow :slides="slides">
        <template #intro>
          <LandingEyebrow tone="corail">ÉDITION 2024</LandingEyebrow>
          <LandingSectionTitle class="mb-5">Retour sur la première édition.</LandingSectionTitle>
          <p class="m-0 max-w-intro text-sm leading-loose text-ink-muted">
            Exposants, auteurs, artisans et visiteurs : en 2024, le salon a réuni toute une communauté autour de la
            culture comorienne. Quelques images pour se souvenir, avant d’écrire la suite ensemble.
          </p>
        </template>
      </PhotoSlideshow>
    </LandingContainer>
  </section>
</template>
