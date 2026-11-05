<script setup lang="ts">
import { reactive, ref } from 'vue'
import LandingButton from '@/components/landing/LandingButton.vue'
import type { BuyerForm } from '@/types/payment'

const emit = defineEmits<{
  submit: [buyer: BuyerForm]
  back: []
}>()

const form = reactive<BuyerForm>({
  first_name: '',
  last_name: '',
  email: '',
  city: '',
})
const error = ref('')

function onSubmit() {
  error.value = ''
  const first = form.first_name.trim()
  const last = form.last_name.trim()
  const email = form.email.trim()
  const city = form.city.trim()
  if (!first || !last || !email || !city) {
    error.value = 'Tous les champs sont obligatoires.'
    return
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    error.value = 'Adresse email invalide.'
    return
  }
  emit('submit', { first_name: first, last_name: last, email, city })
}

const fieldClass =
  'w-full rounded-media border border-line bg-cream px-3.5 py-3 text-sm text-ink outline-none transition-colors placeholder:text-ink-subtle focus:border-forest'
</script>

<template>
  <form class="mt-2 grid gap-3 text-left" @submit.prevent="onSubmit">
    <div class="grid gap-3 sm:grid-cols-2">
      <label class="grid gap-1.5">
        <span class="text-caption font-semibold text-ink-muted">Prénom</span>
        <input v-model="form.first_name" :class="fieldClass" type="text" autocomplete="given-name" required />
      </label>
      <label class="grid gap-1.5">
        <span class="text-caption font-semibold text-ink-muted">Nom</span>
        <input v-model="form.last_name" :class="fieldClass" type="text" autocomplete="family-name" required />
      </label>
    </div>
    <label class="grid gap-1.5">
      <span class="text-caption font-semibold text-ink-muted">Email</span>
      <input v-model="form.email" :class="fieldClass" type="email" autocomplete="email" required />
    </label>
    <label class="grid gap-1.5">
      <span class="text-caption font-semibold text-ink-muted">Ville</span>
      <input v-model="form.city" :class="fieldClass" type="text" autocomplete="address-level2" required />
    </label>

    <p v-if="error" class="m-0 text-caption text-danger-700" role="alert">{{ error }}</p>

    <div class="mt-2 grid gap-2 sm:grid-cols-2">
      <LandingButton type="button" variant="yellow" block @click="emit('back')">Retour</LandingButton>
      <LandingButton type="submit" block>Valider</LandingButton>
    </div>
  </form>
</template>
