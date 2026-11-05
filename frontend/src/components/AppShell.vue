<script setup lang="ts">
import PrimaryButton from '@/components/PrimaryButton.vue'

defineProps<{
  userLabel?: string
  loggingOut?: boolean
}>()

defineEmits<{ logout: [] }>()

const appName = import.meta.env.VITE_APP_NAME || 'Comores Event'
</script>

<template>
  <div class="min-h-screen bg-canvas text-ink">
    <header class="border-b border-line bg-surface">
      <div class="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
        <span class="text-sm font-semibold text-accent-700">{{ appName }}</span>
        <div class="flex items-center gap-3">
          <span v-if="userLabel" class="hidden text-sm text-ink-muted sm:inline">{{ userLabel }}</span>
          <PrimaryButton variant="outline" :loading="loggingOut" @click="$emit('logout')">Se déconnecter</PrimaryButton>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-8">
      <slot />
    </main>
  </div>
</template>
