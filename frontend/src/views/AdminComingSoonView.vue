<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '@/components/AppShell.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loggingOut = ref(false)

async function logout() {
  loggingOut.value = true
  try {
    await auth.logout()
    await router.replace({ name: 'lekeza-login' })
  } finally {
    loggingOut.value = false
  }
}
</script>

<template>
  <AppShell :user-label="auth.displayName" :logging-out="loggingOut" @logout="logout">
    <section class="mx-auto flex max-w-lg flex-col items-start gap-4 rounded-card border border-line bg-surface p-8 shadow-card">
      <p class="text-xs font-semibold tracking-eyebrow text-accent-700 uppercase">Lekeza</p>
      <h1 class="font-serif text-3xl font-normal tracking-tight text-ink">Coming soon</h1>
      <p class="text-sm leading-relaxed text-ink-muted">
        L’espace d’administration arrive bientôt. La connexion admin est déjà opérationnelle.
      </p>
      <PrimaryButton type="button" :loading="loggingOut" @click="logout">Se déconnecter</PrimaryButton>
    </section>
  </AppShell>
</template>
