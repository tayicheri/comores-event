<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import AuthFormCard from '@/components/AuthFormCard.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import TextField from '@/components/TextField.vue'
import { getErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref<string | null>(null)

async function submit() {
  error.value = null
  try {
    await auth.login({ email: email.value.trim(), password: password.value })
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/')
  } catch (e) {
    error.value = getErrorMessage(e, 'Connexion impossible.')
  }
}
</script>

<template>
  <AuthFormCard title="Connexion" subtitle="Accédez à votre compte." :error="error">
    <form class="flex flex-col gap-4" @submit.prevent="submit">
      <TextField v-model="email" label="Adresse e-mail" type="email" autocomplete="email" required />
      <TextField v-model="password" label="Mot de passe" type="password" autocomplete="current-password" required />
      <PrimaryButton type="submit" block :loading="auth.loading">Se connecter</PrimaryButton>
    </form>

    <template #footer>
      Pas encore de compte ?
      <RouterLink :to="{ name: 'register' }" class="font-medium text-accent-700 hover:underline">
        Créer un compte
      </RouterLink>
    </template>
  </AuthFormCard>
</template>
