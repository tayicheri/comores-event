<script setup lang="ts">
import { ref } from 'vue'
import { loadStripe } from '@stripe/stripe-js'
import LandingButton from '@/components/landing/LandingButton.vue'
import CheckoutCardFields from '@/components/landing/checkout/CheckoutCardFields.vue'
import { confirmPayment, createPaymentIntent, ensureCsrfCookie } from '@/api/payments'
import { getErrorMessage } from '@/api/client'
import { formatEuro } from '@/components/landing/eventConfig'
import type { CheckoutStripeFieldsApi } from '@/lib/stripe-elements'
import type { BuyerForm, CheckoutSelection } from '@/types/payment'

const props = defineProps<{
  buyer: BuyerForm
  selection: CheckoutSelection
  amountLabel: string
}>()

const emit = defineEmits<{
  back: []
  paid: []
}>()

const holder = ref(`${props.buyer.first_name} ${props.buyer.last_name}`.trim())
const cardFields = ref<CheckoutStripeFieldsApi | null>(null)
const submitting = ref(false)
const error = ref('')

async function pay() {
  error.value = ''
  submitting.value = true
  try {
    if (!holder.value.trim()) {
      throw new Error('Veuillez saisir le nom du titulaire')
    }
    if (!cardFields.value?.fieldsComplete()) {
      throw new Error('Veuillez renseigner les informations de la carte')
    }
    const publishable = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY ?? ''
    if (!publishable) {
      throw new Error('Clé Stripe publique manquante')
    }

    await ensureCsrfCookie()
    const intent = await createPaymentIntent(props.buyer, props.selection)
    const stripe = cardFields.value?.getStripe() ?? (await loadStripe(publishable))
    const cardElement = cardFields.value?.getCardNumber()
    if (!stripe || !cardElement) {
      throw new Error('Stripe indisponible')
    }

    const result = await stripe.confirmCardPayment(intent.client_secret, {
      payment_method: {
        card: cardElement,
        billing_details: {
          name: holder.value.trim(),
          email: props.buyer.email,
          address: { city: props.buyer.city },
        },
      },
    })
    if (result.error) {
      throw new Error(result.error.message || 'Paiement carte refusé')
    }

    await confirmPayment(intent.payment_intent_id)
    emit('paid')
  } catch (e) {
    error.value = getErrorMessage(e, 'Le paiement a échoué')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="mt-2 grid gap-4">
    <p class="m-0 text-sm text-ink-muted">
      Total à régler :
      <strong class="text-ink">{{ amountLabel || formatEuro(selection.total) }}</strong>
    </p>

    <CheckoutCardFields ref="cardFields" v-model:holder="holder" />

    <p v-if="error" class="m-0 text-caption text-danger-700" role="alert">{{ error }}</p>

    <div class="grid gap-2 sm:grid-cols-2">
      <LandingButton type="button" variant="yellow" block :disabled="submitting" @click="emit('back')">
        Retour
      </LandingButton>
      <LandingButton type="button" block :disabled="submitting" @click="pay">
        {{ submitting ? 'Paiement…' : 'Payer' }}
      </LandingButton>
    </div>
  </div>
</template>
