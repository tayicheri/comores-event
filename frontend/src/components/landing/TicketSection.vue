<script setup lang="ts">
import { computed, ref } from 'vue'
import LandingButton from './LandingButton.vue'
import LandingContainer from './LandingContainer.vue'
import LandingEyebrow from './LandingEyebrow.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import { EVENT_CONFIG, formatEuro, type TicketDay, type TicketFormula } from './eventConfig'

const formula = ref<TicketFormula>('day')
const day = ref<TicketDay>('samedi')
const quantity = ref(1)

const unit = computed(() => EVENT_CONFIG.prices[formula.value])
const total = computed(() => unit.value * quantity.value)
const isDayTicket = computed(() => formula.value === 'day')
const summary = computed(() => {
  const dayLabel = isDayTicket.value ? EVENT_CONFIG.dates[day.value] : 'Pass samedi + dimanche'
  return `${quantity.value} billet${quantity.value > 1 ? 's' : ''} · ${formatEuro(unit.value)} / personne · ${dayLabel}`
})

const emit = defineEmits<{
  checkout: [{ formula: TicketFormula; day: TicketDay; quantity: number; unit: number; total: number }]
}>()

const formulas = [
  { value: 'day' as const, title: '1 journée', price: 20, caption: 'Samedi ou dimanche, au choix' },
  { value: 'pass' as const, title: 'Pass 2 jours', price: 30, caption: 'Samedi et dimanche' },
]

const days = [
  { value: 'samedi' as const, label: 'Sam. 16 janvier' },
  { value: 'dimanche' as const, label: 'Dim. 17 janvier' },
]

const steps = [
  { t: 'Choisissez votre formule', d: 'Une journée ou le week-end complet.' },
  { t: 'Composez votre visite', d: 'Votre journée et le nombre de personnes.' },
  { t: 'Poursuivez votre réservation', d: 'La billetterie FOODYFLY prend le relais.' },
]

function dec() {
  quantity.value = Math.max(1, quantity.value - 1)
}
function inc() {
  quantity.value = Math.min(EVENT_CONFIG.maxTickets, quantity.value + 1)
}
function checkout() {
  emit('checkout', {
    formula: formula.value,
    day: day.value,
    quantity: quantity.value,
    unit: unit.value,
    total: total.value,
  })
}
</script>

<template>
  <section id="billetterie" class="py-section-md md:py-section-lg lg:py-section">
    <LandingContainer class="grid grid-cols-1 items-center gap-9 md:gap-split-md lg:grid-cols-ticket lg:gap-ticket">
      <div>
        <LandingEyebrow>BILLETTERIE · ÉDITION 2027</LandingEyebrow>
        <LandingSectionTitle size="ticket" class="mb-4 max-w-ticket">
          Votre week-end commence ici.
        </LandingSectionTitle>
        <p class="max-w-ticket-desc leading-relaxed text-ink-muted">
          Une formule, une journée, le nombre de personnes. L’essentiel, simplement.
        </p>
        <ol class="mt-8 grid list-none gap-6 p-0">
          <li
            v-for="(step, i) in steps"
            :key="step.t"
            class="grid grid-cols-step items-start gap-x-3.5"
          >
            <span
              class="row-span-2 grid size-7 place-items-center rounded-full border border-line-strong text-xs text-ink-muted"
            >
              {{ String(i + 1).padStart(2, '0') }}
            </span>
            <strong class="text-nav">{{ step.t }}</strong>
            <span class="mt-0.5 text-caption text-ink-subtle">{{ step.d }}</span>
          </li>
        </ol>
      </div>

      <div class="rounded-card border border-line bg-white p-4 shadow-card sm:p-6 md:p-panel" aria-labelledby="ticket-title">
        <div class="mb-6 flex items-center justify-between gap-2.5">
          <LandingSectionTitle id="ticket-title" as="h3" size="panel">Choisir mes billets</LandingSectionTitle>
          <span class="whitespace-nowrap rounded-badge bg-gold-soft px-3.5 py-2 text-micro font-bold md:text-xs">
            16 &amp; 17 JAN. 2027
          </span>
        </div>

        <fieldset class="m-0 grid grid-cols-2 gap-2.5 border-0 p-0 md:gap-3.5">
          <legend class="col-span-2 mb-2 text-sm font-bold">Votre formule</legend>
          <label
            v-for="opt in formulas"
            :key="opt.value"
            class="relative flex min-h-32 cursor-pointer flex-col rounded-media border-2 border-line-option p-3.5 md:min-h-36 md:p-5"
            :class="formula === opt.value ? 'border-forest bg-option' : ''"
          >
            <input
              v-model="formula"
              class="pointer-events-none absolute opacity-0"
              type="radio"
              name="ticket-type"
              :value="opt.value"
            />
            <span class="flex justify-between text-caption font-bold md:text-body">
              {{ opt.title }}
              <span
                class="inline-grid size-radio place-items-center rounded-full border-2"
                :class="formula === opt.value ? 'border-forest bg-forest text-white' : 'border-radio'"
                aria-hidden="true"
              >
                <LandingIcon v-if="formula === opt.value" name="check" class="size-check" />
              </span>
            </span>
            <span class="font-serif my-2.5 text-panel font-medium md:my-4 md:text-price">
              {{ opt.price }} € <small class="font-sans text-caption font-normal text-ink-subtle">/ personne</small>
            </span>
            <span class="text-xs text-ink-subtle md:text-eyebrow">{{ opt.caption }}</span>
          </label>
        </fieldset>

        <fieldset
          class="mt-6 border-0 p-0"
          :class="isDayTicket ? '' : 'pointer-events-none opacity-40'"
          :disabled="!isDayTicket"
        >
          <legend class="mb-2 text-sm font-bold">Votre journée</legend>
          <div class="grid grid-cols-2 gap-2.5">
            <label
              v-for="opt in days"
              :key="opt.value"
              class="relative flex cursor-pointer items-center justify-center rounded-day border px-2 py-3 text-center text-caption font-bold"
              :class="day === opt.value ? 'border-forest bg-forest text-white' : 'border-line-option bg-white'"
            >
              <input
                v-model="day"
                class="pointer-events-none absolute opacity-0"
                type="radio"
                name="ticket-day"
                :value="opt.value"
              />
              <LandingIcon
                name="check"
                class="absolute top-1/2 left-3 size-icon-sm -translate-y-1/2"
                :class="day === opt.value ? 'opacity-100' : 'opacity-0'"
              />
              {{ opt.label }}
            </label>
          </div>
        </fieldset>

        <div class="my-6 flex items-center justify-between">
          <div>
            <strong class="block text-nav">Nombre de personnes</strong>
            <small class="mt-1 block text-eyebrow text-ink-subtle">Billets par personne</small>
          </div>
          <div class="flex items-center gap-5">
            <button
              type="button"
              class="grid size-qty place-items-center rounded-full border border-line-strong bg-white text-forest disabled:cursor-default disabled:opacity-35"
              aria-label="Retirer un billet"
              :disabled="quantity <= 1"
              @click="dec"
            >
              <LandingIcon name="minus" class="size-icon-sm" />
            </button>
            <output class="min-w-3.5 text-center text-sm" aria-live="polite">{{ quantity }}</output>
            <button
              type="button"
              class="grid size-qty place-items-center rounded-full border border-line-strong bg-white text-forest disabled:cursor-default disabled:opacity-35"
              aria-label="Ajouter un billet"
              :disabled="quantity >= EVENT_CONFIG.maxTickets"
              @click="inc"
            >
              <LandingIcon name="plus" class="size-icon-sm" />
            </button>
          </div>
        </div>

        <div class="flex items-center justify-between gap-2 border-t border-line py-6">
          <div>
            <strong class="block text-nav">Total</strong>
            <small class="mt-1 block text-eyebrow text-ink-subtle">{{ summary }}</small>
          </div>
          <output class="font-serif text-price font-medium" aria-live="polite">{{ formatEuro(total) }}</output>
        </div>

        <LandingButton block @click="checkout">
          Poursuivre avec FOODYFLY
          <LandingIcon name="arrow" class="size-icon-sm" />
        </LandingButton>
        <p class="mt-3.5 text-center text-xs text-ink-subtle">
          Maquette · Intégration de l’API de paiement non active.
        </p>
      </div>
    </LandingContainer>
  </section>
</template>
