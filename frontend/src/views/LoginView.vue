<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AuthFormCard from '@/components/AuthFormCard.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import TextField from '@/components/TextField.vue'
import { getErrorMessage } from '@/api/client'
import { usePageSeo } from '@/composables/usePageSeo'
import { useAuthStore } from '@/stores/auth'

usePageSeo()

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref<string | null>(null)

function safeRedirect(raw: unknown): string {
  if (typeof raw !== 'string') return '/lekeza'
  if (!raw.startsWith('/') || raw.startsWith('//')) return '/lekeza'
  if (!raw.startsWith('/lekeza')) return '/lekeza'
  return raw
}

async function submit() {
  error.value = null
  try {
    await auth.login({ email: email.value.trim(), password: password.value })
    if (!auth.isAdmin) {
      await auth.logout()
      error.value = 'Accès réservé aux administrateurs.'
      return
    }
    await router.replace(safeRedirect(route.query.redirect))
  } catch (e) {
    error.value = getErrorMessage(e, 'Connexion impossible.')
  }
}
</script>

<template>
  <AuthFormCard
    title="Espace admin"
    subtitle="Connectez-vous pour accéder à Lekeza."
    :error="error"
  >
    <form class="flex flex-col gap-4" @submit.prevent="submit">
      <TextField v-model="email" label="Adresse e-mail" type="email" autocomplete="username" required />
      <TextField
        v-model="password"
        label="Mot de passe"
        type="password"
        autocomplete="current-password"
        required
      />
      <PrimaryButton type="submit" block :loading="auth.loading">Se connecter</PrimaryButton>
    </form>
  </AuthFormCard>
</template>
