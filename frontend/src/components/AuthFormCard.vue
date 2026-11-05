<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  error?: string | null
}>()

const appName = import.meta.env.VITE_APP_NAME || 'Comores Event'
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-canvas px-4 py-10">
    <div class="w-full max-w-sm">
      <p class="mb-6 text-center text-sm font-semibold tracking-wide text-accent-700">{{ appName }}</p>

      <section class="rounded-card border border-line bg-surface p-6 shadow-card">
        <header class="mb-5">
          <h1 class="text-xl font-semibold text-ink">{{ title }}</h1>
          <p v-if="subtitle" class="mt-1 text-sm text-ink-muted">{{ subtitle }}</p>
        </header>

        <p
          v-if="error"
          role="alert"
          class="mb-4 rounded-control border border-danger-600/20 bg-danger-50 px-3 py-2 text-sm text-danger-700"
        >
          {{ error }}
        </p>

        <slot />
      </section>

      <p class="mt-4 text-center text-sm text-ink-muted">
        <slot name="footer" />
      </p>
    </div>
  </main>
</template>
