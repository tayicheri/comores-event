<script setup lang="ts">
import { ref } from 'vue'
import AdminBilletCommentCell from '@/components/admin/AdminBilletCommentCell.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import type { Billet } from '@/types/payment'

defineProps<{
  billets: Billet[]
  busyId?: string | null
}>()

const emit = defineEmits<{
  markUsed: [id: string]
  editComment: [billet: Billet]
}>()

function formatDate(iso: string): string {
  return new Intl.DateTimeFormat('fr-FR', {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(iso))
}

function daysLabel(b: Billet): string {
  const parts: string[] = []
  if (b.day_vendredi) parts.push('Ven')
  if (b.day_samedi) parts.push('Sam')
  return parts.join(' + ') || '—'
}

function euro(cents: number): string {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
  }).format(cents / 100)
}

const confirming = ref<string | null>(null)

function requestMark(id: string) {
  confirming.value = id
}

function confirmMark(id: string) {
  confirming.value = null
  emit('markUsed', id)
}
</script>

<template>
  <div class="overflow-x-auto rounded-card border border-line bg-surface shadow-card">
    <table class="min-w-full border-collapse text-left text-sm">
      <thead class="border-b border-line bg-pale text-caption tracking-tight text-ink-muted uppercase">
        <tr>
          <th class="px-4 py-3 font-semibold">Achat</th>
          <th class="px-4 py-3 font-semibold">Acheteur</th>
          <th class="px-4 py-3 font-semibold">Qté</th>
          <th class="px-4 py-3 font-semibold">Jours</th>
          <th class="px-4 py-3 font-semibold">Commentaire</th>
          <th class="px-4 py-3 font-semibold">Montant</th>
          <th class="px-4 py-3 font-semibold">Statut</th>
          <th class="px-4 py-3 font-semibold">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="billets.length === 0">
          <td colspan="8" class="px-4 py-8 text-center text-ink-muted">Aucun billet pour le moment.</td>
        </tr>
        <tr
          v-for="b in billets"
          :key="b.id"
          class="border-b border-line last:border-0"
        >
          <td class="px-4 py-3 whitespace-nowrap text-ink-muted">{{ formatDate(b.purchased_at) }}</td>
          <td class="px-4 py-3">
            <div class="font-medium text-ink">{{ b.first_name }} {{ b.last_name }}</div>
            <div class="text-caption text-ink-muted">{{ b.email }} · {{ b.city }}</div>
          </td>
          <td class="px-4 py-3">{{ b.quantity }}</td>
          <td class="px-4 py-3">{{ daysLabel(b) }}</td>
          <td class="px-4 py-3">
            <AdminBilletCommentCell :comments="b.comments" @edit="emit('editComment', b)" />
          </td>
          <td class="px-4 py-3 whitespace-nowrap">{{ euro(b.amount_cents) }}</td>
          <td class="px-4 py-3">
            <span
              class="inline-flex rounded-full px-2.5 py-1 text-2xs font-bold uppercase"
              :class="
                b.status === 'used'
                  ? 'bg-pale text-ink-muted'
                  : 'bg-gold-soft text-forest'
              "
            >
              {{ b.status === 'used' ? 'Utilisé' : 'Valide' }}
            </span>
          </td>
          <td class="px-4 py-3">
            <template v-if="b.status === 'valid'">
              <div v-if="confirming === b.id" class="flex flex-wrap gap-2">
                <PrimaryButton type="button" :loading="busyId === b.id" @click="confirmMark(b.id)">
                  Confirmer
                </PrimaryButton>
                <PrimaryButton type="button" variant="outline" @click="confirming = null">
                  Annuler
                </PrimaryButton>
              </div>
              <PrimaryButton
                v-else
                type="button"
                variant="outline"
                :disabled="busyId === b.id"
                @click="requestMark(b.id)"
              >
                Marquer utilisé
              </PrimaryButton>
            </template>
            <span v-else class="text-caption text-ink-subtle">—</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
