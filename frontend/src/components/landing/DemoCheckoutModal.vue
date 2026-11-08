<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import LandingButton from './LandingButton.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'
import CheckoutBuyerForm from './checkout/CheckoutBuyerForm.vue'
import CheckoutPaymentStep from './checkout/CheckoutPaymentStep.vue'
import { EVENT_CONFIG, formatEuro } from './eventConfig'
import type { BuyerForm, CheckoutSelection } from '@/types/payment'

const props = defineProps<{
  open: boolean
  selection: CheckoutSelection | null
}>()

const emit = defineEmits<{ close: [] }>()

type Step = 'recap' | 'buyer' | 'payment' | 'success'

const step = ref<Step>('recap')
const buyer = ref<BuyerForm | null>(null)

const lines = computed(() => {
  const sel = props.selection
  if (!sel) return [] as string[]
  const formulaLabel =
    sel.formula === 'day'
      ? `Billet 1 journée — ${EVENT_CONFIG.dates[sel.day]}`
      : 'Pass 2 jours — samedi et dimanche'
  return [
    formulaLabel,
    `${sel.quantity} billet${sel.quantity > 1 ? 's' : ''} × ${formatEuro(sel.unit)}`,
    `Total : ${formatEuro(sel.total)}`,
  ]
})

const title = computed(() => {
  switch (step.value) {
    case 'buyer':
      return 'Vos coordonnées'
    case 'payment':
      return 'Paiement sécurisé'
    case 'success':
      return 'Réservation confirmée'
    default:
      return 'Votre sélection est prête'
  }
})

watch(
  () => props.open,
  (open) => {
    if (typeof document !== 'undefined') {
      document.body.style.overflow = open ? 'hidden' : ''
    }
    if (open) {
      step.value = 'recap'
      buyer.value = null
    }
  },
)

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) emit('close')
}

onMounted(() => document.addEventListener('keydown', onKey))
onUnmounted(() => {
  document.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})

function onBuyerSubmit(payload: BuyerForm) {
  buyer.value = payload
  step.value = 'payment'
}
</script>

<template>
  <div
    v-show="open"
    class="fixed inset-0 z-50 grid place-items-center bg-forest-deep/70 p-5"
    role="presentation"
    @click.self="emit('close')"
  >
    <div
      class="relative max-h-[min(92vh,44rem)] w-full max-w-modal overflow-y-auto rounded-card bg-surface p-panel text-center shadow-modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <button
        type="button"
        class="absolute top-3.5 right-3.5 grid size-9 place-items-center rounded-full border-0 bg-transparent text-ink hover:bg-pale"
        aria-label="Fermer"
        @click="emit('close')"
      >
        <LandingIcon name="x" class="size-icon-sm" />
      </button>

      <div
        class="mx-auto mt-1 grid size-16 place-items-center rounded-full bg-gold text-forest"
        aria-hidden="true"
      >
        <LandingIcon :name="step === 'success' ? 'star' : 'ticket'" class="size-7" />
      </div>

      <LandingSectionTitle id="modal-title" as="h2" size="modal" class="my-3">
        {{ title }}
      </LandingSectionTitle>

      <template v-if="step === 'recap'">
        <p class="text-sm leading-body text-ink-muted">Voici le récapitulatif de votre réservation :</p>
        <div class="my-4 rounded-media bg-pale p-4 text-sm leading-relaxed font-semibold">
          <div v-for="line in lines" :key="line">{{ line }}</div>
        </div>
        <LandingButton class="mt-3.5" block @click="step = 'buyer'">Continuer</LandingButton>
      </template>

      <template v-else-if="step === 'buyer'">
        <p class="mb-1 text-sm leading-body text-ink-muted">
          Indiquez vos informations pour recevoir la confirmation.
        </p>
        <CheckoutBuyerForm @submit="onBuyerSubmit" @back="step = 'recap'" />
      </template>

      <template v-else-if="step === 'payment' && buyer && selection">
        <CheckoutPaymentStep
          :buyer="buyer"
          :selection="selection"
          :amount-label="formatEuro(selection.total)"
          @back="step = 'buyer'"
          @paid="step = 'success'"
        />
      </template>

      <template v-else-if="step === 'success'">
        <p class="text-sm leading-body text-ink-muted">
          Merci ! Votre paiement a bien été enregistré. Un billet valide figure désormais dans
          notre billetterie.
        </p>
        <LandingButton class="mt-4" block @click="emit('close')">Fermer</LandingButton>
      </template>
    </div>
  </div>
</template>
