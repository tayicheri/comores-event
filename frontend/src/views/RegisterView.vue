<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import AuthFormCard from '@/components/AuthFormCard.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import TextField from '@/components/TextField.vue'
import { getErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const displayName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref<string | null>(null)

async function submit() {
  error.value = null
  if (password.value !== confirmPassword.value) {
    error.value = 'Les mots de passe ne correspondent pas.'
    return
  }
  try {
    const name = displayName.value.trim()
    await auth.register({
      ...(name ? { display_name: name } : {}),
      email: email.value.trim(),
      password: password.value,
    })
    await router.replace(auth.isAuthenticated ? { name: 'home' } : { name: 'login' })
  } catch (e) {
    error.value = getErrorMessage(e, 'Inscription impossible.')
  }
}
</script>

<template>
  <AuthFormCard title="Créer un compte" subtitle="Quelques informations suffisent." :error="error">
    <form class="flex flex-col gap-4" @submit.prevent="submit">
      <TextField v-model="displayName" label="Nom affiché" autocomplete="name" />
      <TextField v-model="email" label="Adresse e-mail" type="email" autocomplete="email" required />
      <TextField
        v-model="password"
        label="Mot de passe"
        type="password"
        autocomplete="new-password"
        :minlength="8"
        required
      />
      <TextField
        v-model="confirmPassword"
        label="Confirmer le mot de passe"
        type="password"
        autocomplete="new-password"
        :minlength="8"
        required
      />
      <PrimaryButton type="submit" block :loading="auth.loading">Créer mon compte</PrimaryButton>
    </form>

    <template #footer>
      Déjà un compte ?
      <RouterLink :to="{ name: 'login' }" class="font-medium text-accent-700 hover:underline">Se connecter</RouterLink>
    </template>
  </AuthFormCard>
</template>
