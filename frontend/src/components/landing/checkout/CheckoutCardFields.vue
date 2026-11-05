<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { loadStripe } from '@stripe/stripe-js'
import type {
  Stripe,
  StripeCardCvcElement,
  StripeCardExpiryElement,
  StripeCardNumberElement,
  StripeElements,
} from '@stripe/stripe-js'
import CheckoutCardFieldNav, {
  type CheckoutCardFieldKey,
} from '@/components/landing/checkout/CheckoutCardFieldNav.vue'
import CheckoutCardVisual from '@/components/landing/checkout/CheckoutCardVisual.vue'
import {
  networkFromStripeBrand,
  stripeCardElementStyleOnCard,
  stripeCardElementStyleOnCardCvc,
  stripeCardFieldsComplete,
  type CardNetwork,
  type CheckoutStripeFieldsApi,
} from '@/lib/stripe-elements'

const holder = defineModel<string>('holder', { default: '' })

const holderInput = ref<HTMLInputElement | null>(null)
const numberMount = ref<HTMLDivElement | null>(null)
const expiryMount = ref<HTMLDivElement | null>(null)
const cvcMount = ref<HTMLDivElement | null>(null)
const numberInvalid = ref(false)
const expiryInvalid = ref(false)
const cvcInvalid = ref(false)
const brand = ref<CardNetwork>('CARD')
const flipped = ref(false)
const activeField = ref<CheckoutCardFieldKey | null>(null)

const complete = reactive({ number: false, expiry: false, cvc: false })

let cancelled = false
let stripeInst: Stripe | null = null
let elements: StripeElements | null = null
let cardNumberEl: StripeCardNumberElement | null = null
let cardExpiryEl: StripeCardExpiryElement | null = null
let cardCvcEl: StripeCardCvcElement | null = null

function destroyElements() {
  cardNumberEl?.destroy()
  cardExpiryEl?.destroy()
  cardCvcEl?.destroy()
  cardNumberEl = null
  cardExpiryEl = null
  cardCvcEl = null
  elements = null
}

async function mountElements() {
  const publishable = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY ?? ''
  if (!publishable) return
  const stripe = await loadStripe(publishable)
  if (cancelled || !stripe) return
  stripeInst = stripe
  await nextTick()
  if (cancelled || !numberMount.value || !expiryMount.value || !cvcMount.value || cardNumberEl) {
    return
  }
  elements = stripe.elements()
  cardNumberEl = elements.create('cardNumber', {
    style: stripeCardElementStyleOnCard,
    placeholder: '•••• •••• •••• ••••',
    showIcon: false,
  })
  cardExpiryEl = elements.create('cardExpiry', {
    style: stripeCardElementStyleOnCard,
    placeholder: 'MM/AA',
  })
  cardCvcEl = elements.create('cardCvc', {
    style: stripeCardElementStyleOnCardCvc,
    placeholder: '•••',
  })
  cardNumberEl.on('change', (event) => {
    complete.number = event.complete
    numberInvalid.value = Boolean(event.error)
    brand.value = networkFromStripeBrand(event.brand)
  })
  cardExpiryEl.on('change', (event) => {
    complete.expiry = event.complete
    expiryInvalid.value = Boolean(event.error)
  })
  cardCvcEl.on('change', (event) => {
    complete.cvc = event.complete
    cvcInvalid.value = Boolean(event.error)
  })
  cardNumberEl.on('focus', () => {
    activeField.value = 'number'
    flipped.value = false
  })
  cardNumberEl.on('blur', () => {
    if (activeField.value === 'number') activeField.value = null
  })
  cardExpiryEl.on('focus', () => {
    activeField.value = 'expiry'
    flipped.value = false
  })
  cardExpiryEl.on('blur', () => {
    if (activeField.value === 'expiry') activeField.value = null
  })
  cardCvcEl.on('focus', () => {
    activeField.value = 'cvc'
    flipped.value = true
  })
  cardCvcEl.on('blur', () => {
    flipped.value = false
    if (activeField.value === 'cvc') activeField.value = null
  })
  cardNumberEl.mount(numberMount.value)
  cardExpiryEl.mount(expiryMount.value)
  cardCvcEl.mount(cvcMount.value)
}

function focusField(field: CheckoutCardFieldKey) {
  if (field === 'holder') {
    flipped.value = false
    holderInput.value?.focus()
    return
  }
  if (field === 'number') {
    flipped.value = false
    cardNumberEl?.focus()
    return
  }
  if (field === 'expiry') {
    flipped.value = false
    cardExpiryEl?.focus()
    return
  }
  flipped.value = true
  cardCvcEl?.focus()
}

const api: CheckoutStripeFieldsApi = {
  getStripe: () => stripeInst,
  getCardNumber: () => cardNumberEl,
  fieldsComplete: () => stripeCardFieldsComplete(complete),
}

defineExpose(api)

onMounted(() => {
  void mountElements()
})

onBeforeUnmount(() => {
  cancelled = true
  destroyElements()
})
</script>

<template>
  <div class="grid w-full min-w-0 gap-4">
    <CheckoutCardVisual :brand="brand" :flipped="flipped">
      <template #number>
        <div
          ref="numberMount"
          class="relative box-border block min-h-7 w-full min-w-0 bg-transparent px-0 py-0.5 outline-0 [&>div]:!w-full [&_iframe]:!border-0 [&_iframe]:!bg-transparent"
          :class="{ 'shadow-[inset_0_-1px_0_var(--color-gold-soft)]': numberInvalid }"
        />
      </template>
      <template #holder>
        <input
          ref="holderInput"
          v-model="holder"
          class="mt-0.5 box-border block w-full max-w-[11rem] appearance-none border-0 bg-transparent py-0.5 font-serif text-[clamp(12px,3.4vw,15px)] font-bold tracking-[0.04em] text-surface uppercase caret-transparent outline-0 placeholder:text-white/40 placeholder:uppercase"
          type="text"
          autocomplete="cc-name"
          placeholder="JEAN DUPONT"
          @focus="activeField = 'holder'; flipped = false"
          @blur="activeField = activeField === 'holder' ? null : activeField"
        />
      </template>
      <template #expiry>
        <div
          ref="expiryMount"
          class="relative box-border block min-h-7 w-full min-w-0 bg-transparent px-0 py-0.5 outline-0 [&>div]:!w-full [&_iframe]:!border-0 [&_iframe]:!bg-transparent"
          :class="{ 'shadow-[inset_0_-1px_0_var(--color-gold-soft)]': expiryInvalid }"
        />
      </template>
      <template #cvc>
        <div
          ref="cvcMount"
          class="relative box-border block min-h-7 w-full min-w-0 bg-transparent px-0 py-0.5 outline-0 [&>div]:!w-full [&_iframe]:!border-0 [&_iframe]:!bg-transparent"
          :class="{ 'shadow-[inset_0_-1px_0_var(--color-danger-600)]': cvcInvalid }"
        />
      </template>
    </CheckoutCardVisual>

    <CheckoutCardFieldNav :active="activeField" @select="focusField" />
  </div>
</template>
