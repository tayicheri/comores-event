<script setup lang="ts">
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const items = [
  {
    q: 'Quelle formule choisir ?',
    a: 'Le billet 1 journée est à 20 € par personne, avec le choix du samedi 16 ou du dimanche 17 janvier. Le pass 2 jours est à 30 € par personne pour profiter des deux journées.',
    open: true,
  },
  {
    q: 'Où se déroule le salon ?',
    a: 'À la Salle des congrès, 118 rue du 8 Mai 1945, 92000 Nanterre, France.',
  },
  {
    q: 'Quand connaître les horaires et les intervenants ?',
    a: 'Les horaires détaillés, les artistes et les exposants seront annoncés par l’organisation avant le salon.',
  },
  {
    q: 'La billetterie est-elle déjà active ?',
    a: 'Oui. Vous pouvez choisir votre formule, renseigner vos coordonnées et régler par carte bancaire via Stripe. Un billet valide est créé dès confirmation du paiement.',
  },
]
</script>

<template>
  <section id="questions" class="py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-start gap-6 lg:grid-cols-faq lg:gap-faq">
      <div>
        <LandingEyebrow>AVANT DE VENIR</LandingEyebrow>
        <LandingSectionTitle size="faq">
          Quelques réponses,<br />en toute simplicité.
        </LandingSectionTitle>
      </div>
      <div class="border-t border-line">
        <details
          v-for="item in items"
          :key="item.q"
          class="group border-b border-line"
          :open="item.open || undefined"
        >
          <summary
            class="flex cursor-pointer list-none items-center justify-between gap-4 py-6 pr-1 text-sm font-bold marker:content-none [&::-webkit-details-marker]:hidden"
          >
            {{ item.q }}
            <span
              class="relative size-faq-toggle shrink-0 before:absolute before:inset-x-0.5 before:top-1/2 before:h-0.5 before:-translate-y-1/2 before:bg-forest after:absolute after:inset-y-0.5 after:left-1/2 after:w-0.5 after:-translate-x-1/2 after:bg-forest group-open:after:hidden"
              aria-hidden="true"
            />
          </summary>
          <p class="mb-6 max-w-faq text-nav leading-relaxed text-ink-muted">{{ item.a }}</p>
        </details>
      </div>
    </LandingContainer>
  </section>
</template>
