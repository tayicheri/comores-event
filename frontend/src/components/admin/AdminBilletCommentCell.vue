<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref } from 'vue'
import PrimaryButton from '@/components/PrimaryButton.vue'

const props = defineProps<{
  comments: string | null
}>()

const emit = defineEmits<{
  edit: []
}>()

const open = ref(false)
const anchorEl = ref<HTMLElement | null>(null)
const panelEl = ref<HTMLElement | null>(null)
const panelStyle = ref<Record<string, string>>({})
let closeTimer: ReturnType<typeof setTimeout> | null = null

const display = computed(() => {
  const t = props.comments?.trim()
  return t && t.length > 0 ? t : null
})

function clearCloseTimer() {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
}

function positionPanel() {
  const el = anchorEl.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const width = 288
  const margin = 8
  let left = rect.left
  if (left + width > window.innerWidth - margin) {
    left = Math.max(margin, window.innerWidth - width - margin)
  }
  const spaceBelow = window.innerHeight - rect.bottom
  const preferBelow = spaceBelow >= 180
  panelStyle.value = preferBelow
    ? {
        top: `${rect.bottom + 6}px`,
        left: `${left}px`,
        width: `${width}px`,
      }
    : {
        bottom: `${window.innerHeight - rect.top + 6}px`,
        left: `${left}px`,
        width: `${width}px`,
      }
}

async function show() {
  clearCloseTimer()
  open.value = true
  await nextTick()
  positionPanel()
}

function scheduleHide() {
  clearCloseTimer()
  closeTimer = setTimeout(() => {
    open.value = false
  }, 140)
}

function onFocusOut(e: FocusEvent) {
  const next = e.relatedTarget as Node | null
  if (next && (anchorEl.value?.contains(next) || panelEl.value?.contains(next))) return
  scheduleHide()
}

function onEdit(e: Event) {
  e.preventDefault()
  clearCloseTimer()
  open.value = false
  window.removeEventListener('keydown', onKeydown)
  emit('edit')
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && open.value) open.value = false
}

function onPanelEnter() {
  clearCloseTimer()
  window.addEventListener('keydown', onKeydown)
}

function onPanelLeave() {
  scheduleHide()
  window.removeEventListener('keydown', onKeydown)
}

onUnmounted(() => {
  clearCloseTimer()
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div
    ref="anchorEl"
    class="relative max-w-44"
    @mouseenter="show"
    @mouseleave="scheduleHide"
    @focusin="show"
    @focusout="onFocusOut"
  >
    <button
      type="button"
      class="block w-full max-w-44 truncate rounded-control text-left text-sm text-ink focus:outline-none focus-visible:ring-2 focus-visible:ring-accent-500/40"
      :class="display ? '' : 'text-ink-muted'"
      :aria-expanded="open"
      aria-haspopup="dialog"
    >
      {{ display ?? '—' }}
    </button>

    <Teleport to="body">
      <div
        v-if="open"
        ref="panelEl"
        class="fixed z-50 rounded-card border border-line bg-surface p-3 shadow-card"
        :style="panelStyle"
        role="dialog"
        aria-label="Commentaire"
        @mouseenter="onPanelEnter"
        @mouseleave="onPanelLeave"
      >
        <p class="max-h-40 overflow-y-auto text-sm leading-relaxed break-words whitespace-pre-wrap text-ink">
          {{ display ?? 'Aucun commentaire' }}
        </p>
        <div class="mt-3 border-t border-line pt-3">
          <PrimaryButton type="button" variant="outline" block @mousedown="onEdit">
            Modifier
          </PrimaryButton>
        </div>
      </div>
    </Teleport>
  </div>
</template>
