<script setup lang="ts">
import type { BilletStats } from '@/types/payment'

defineProps<{
  stats: BilletStats
}>()

function euro(cents: number): string {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
    maximumFractionDigits: 0,
  }).format(cents / 100)
}
</script>

<template>
  <div class="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-6">
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">Billets</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ stats.total_count }}</p>
    </article>
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">CA</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ euro(stats.revenue_cents) }}</p>
    </article>
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">Valides</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ stats.valid_count }}</p>
    </article>
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">Utilisés</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ stats.used_count }}</p>
    </article>
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">Samedi</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ stats.vendredi_places }}</p>
    </article>
    <article class="rounded-card border border-line bg-surface p-4 shadow-card">
      <p class="text-2xs font-semibold tracking-eyebrow text-accent-700 uppercase">Dimanche</p>
      <p class="mt-1 font-serif text-2xl text-ink">{{ stats.samedi_places }}</p>
    </article>
  </div>
</template>
