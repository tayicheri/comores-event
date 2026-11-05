<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppShell from '@/components/AppShell.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loggingOut = ref(false)

async function logout() {
  loggingOut.value = true
  try {
    await auth.logout()
    await router.replace({ name: 'login' })
  } finally {
    loggingOut.value = false
  }
}
</script>

<template>
  <AppShell :user-label="auth.displayName" :logging-out="loggingOut" @logout="logout">
    <section class="rounded-card border border-line bg-surface p-6 shadow-card">
      <h1 class="text-xl font-semibold text-ink">Bonjour{{ auth.displayName ? `, ${auth.displayName}` : '' }}</h1>
      <p class="mt-1 text-sm text-ink-muted">Vous êtes connecté. Cet espace est prêt à accueillir vos fonctionnalités.</p>
    </section>
  </AppShell>
</template>
