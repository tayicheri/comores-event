<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import PrimaryButton from '@/components/PrimaryButton.vue'

const props = defineProps<{
  open: boolean
  initialComments: string | null
  busy?: boolean
}>()

const emit = defineEmits<{
  close: []
  save: [comments: string]
}>()

const draft = ref('')

watch(
  () => props.open,
  (isOpen) => {
    document.body.style.overflow = isOpen ? 'hidden' : ''
    if (isOpen) draft.value = props.initialComments ?? ''
  },
)

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open && !props.busy) emit('close')
}

onMounted(() => document.addEventListener('keydown', onKey))
onUnmounted(() => {
  document.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})

function submit() {
  emit('save', draft.value.trim())
}
</script>

<template>
  <div
    v-show="open"
    class="fixed inset-0 z-50 grid place-items-center bg-ink/40 p-5"
    role="presentation"
    @click.self="!busy && emit('close')"
  >
    <div
      class="w-full max-w-md rounded-card border border-line bg-surface p-5 shadow-card"
      role="dialog"
      aria-modal="true"
      aria-labelledby="comment-edit-title"
    >
      <h2 id="comment-edit-title" class="font-serif text-xl text-ink">Modifier le commentaire</h2>
      <p class="mt-1 text-sm text-ink-muted">Visible uniquement dans l’espace admin.</p>

      <label class="mt-4 flex flex-col gap-1.5">
        <span class="text-sm font-medium text-ink">Commentaire</span>
        <textarea
          v-model="draft"
          rows="5"
          maxlength="2000"
          :disabled="busy"
          class="w-full resize-y rounded-control border border-line-strong bg-surface px-3 py-2 text-sm text-ink placeholder:text-ink-subtle focus:border-accent-600 focus:outline-none focus:ring-2 focus:ring-accent-500/30 disabled:cursor-not-allowed disabled:bg-canvas"
          placeholder="Ajouter un commentaire…"
        />
      </label>

      <div class="mt-4 flex flex-wrap justify-end gap-2">
        <PrimaryButton type="button" variant="outline" :disabled="busy" @click="emit('close')">
          Annuler
        </PrimaryButton>
        <PrimaryButton type="button" :loading="busy" @click="submit">
          Enregistrer
        </PrimaryButton>
      </div>
    </div>
  </div>
</template>
