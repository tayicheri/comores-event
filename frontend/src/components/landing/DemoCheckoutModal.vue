<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
import LandingButton from './LandingButton.vue'
import LandingIcon from './LandingIcon.vue'
import LandingSectionTitle from './LandingSectionTitle.vue'

const props = defineProps<{
  open: boolean
  lines: string[]
}>()

const emit = defineEmits<{ close: [] }>()

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) emit('close')
}

watch(
  () => props.open,
  (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
  },
)

onMounted(() => document.addEventListener('keydown', onKey))
onUnmounted(() => {
  document.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})
</script>

<template>
  <div
    v-show="open"
    class="fixed inset-0 z-50 grid place-items-center bg-forest-deep/70 p-5"
    role="presentation"
    @click.self="emit('close')"
  >
    <div
      class="relative w-full max-w-modal rounded-card bg-surface p-panel text-center shadow-modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <button
        type="button"
        class="absolute top-3.5 right-3.5 grid size-9 place-items-center rounded-full border-0 bg-transparent text-ink hover:bg-pale"
        aria-label="Fermer"
        @click="emit('close')"
      >
        <LandingIcon name="x" class="size-icon-sm" />
      </button>
      <div class="mx-auto mt-1 grid size-16 place-items-center rounded-full bg-gold text-forest" aria-hidden="true">
        <LandingIcon name="star" class="size-7" />
      </div>
      <LandingSectionTitle id="modal-title" as="h2" size="modal" class="my-3">
        Votre sélection est prête
      </LandingSectionTitle>
      <p class="text-sm leading-body text-ink-muted">Voici le récapitulatif de votre réservation :</p>
      <div class="my-4 rounded-media bg-pale p-4 text-sm leading-relaxed font-semibold">
        <div v-for="line in lines" :key="line">{{ line }}</div>
      </div>
      <p class="text-caption leading-body text-ink-muted">
        Cette maquette ne déclenche aucun paiement. Le bouton final sera relié à la billetterie FOODYFLY via son API.
      </p>
      <LandingButton class="mt-3.5" @click="emit('close')">Revenir au salon</LandingButton>
    </div>
  </div>
</template>
