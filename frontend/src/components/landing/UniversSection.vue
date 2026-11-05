<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const cards = [
  {
    img: '/landing/gastronomie.jpg',
    alt: 'Plat de gastronomie comorienne servi à table',
    category: '1 / SAVEURS',
    title: 'Gastronomie',
    body: 'Recettes, produits et traditions culinaires : une invitation à découvrir les saveurs qui se transmettent.',
  },
  {
    img: '/landing/artisanat.jpg',
    alt: 'Tressage artisanal d’un panier',
    category: '2 / TRADITIONS',
    title: 'Culture & patrimoine',
    body: 'Métiers, traditions et savoir-faire comoriens. Prendre le temps de regarder, de comprendre et d’échanger.',
  },
  {
    img: '/landing/rencontres.jpg',
    alt: 'Personnes réunies autour de traditions musicales',
    category: '3 / ARTS & PARTAGE',
    title: 'Animations & rencontres',
    body: 'La musique comme lien, les échanges entre générations comme horizon. Une culture qui se vit ensemble.',
  },
]
</script>

<template>
  <section id="univers" class="pt-4 pb-univers-b-sm lg:pb-univers-b">
    <LandingContainer>
      <div class="mb-8 block items-end justify-between gap-9 md:flex">
        <div>
          <LandingEyebrow>LES UNIVERS DU SALON</LandingEyebrow>
          <LandingSectionTitle size="univers">Trois univers, une même générosité.</LandingSectionTitle>
        </div>
        <p class="mt-4 max-w-none text-nav leading-body text-ink-muted md:mt-0 md:max-w-univers-aside">
          Laissez-vous guider par vos envies, au fil des découvertes.
        </p>
      </div>

      <div class="grid grid-cols-1 gap-6 sm:grid-cols-cards-md lg:grid-cols-cards">
        <article
          v-for="(card, index) in cards"
          :key="card.title"
          class="group"
          :class="index === 2 ? 'sm:col-span-2 lg:col-span-1' : ''"
        >
          <img
            :src="card.img"
            :alt="card.alt"
            class="h-card-img-xs w-full rounded-media object-cover sm:h-card-img-sm lg:h-card-img"
            :class="index === 2 ? 'sm:h-card-img-wide lg:h-card-img' : ''"
          />
          <div class="mt-5 text-xs tracking-tight-label text-terracotta">{{ card.category }}</div>
          <div class="flex items-center justify-between gap-2">
            <LandingSectionTitle as="h3" size="card" class="my-1.5">{{ card.title }}</LandingSectionTitle>
            <span
              class="mt-1.5 mb-2.5 grid size-8 shrink-0 place-items-center rounded-full border border-line text-forest transition-colors group-hover:border-forest group-hover:bg-forest group-hover:text-gold"
              aria-hidden="true"
            >
              <LandingIcon name="arrow" class="size-icon-sm" />
            </span>
          </div>
          <p class="m-0 text-nav leading-card-body text-ink-muted">{{ card.body }}</p>
        </article>
      </div>
      <p class="mt-9 text-caption text-ink-subtle">
        Le programme et les intervenants seront communiqués par l’organisation.
      </p>
    </LandingContainer>
  </section>
</template>
