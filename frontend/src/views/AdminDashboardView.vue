<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AdminBilletCommentEditModal from '@/components/admin/AdminBilletCommentEditModal.vue'
import AdminBilletsTable from '@/components/admin/AdminBilletsTable.vue'
import AdminStatsStrip from '@/components/admin/AdminStatsStrip.vue'
import AppShell from '@/components/AppShell.vue'
import PrimaryButton from '@/components/PrimaryButton.vue'
import { fetchAdminBillets, updateAdminBillet } from '@/api/billets'
import { getErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { Billet, BilletStats } from '@/types/payment'

const auth = useAuthStore()
const router = useRouter()

const loggingOut = ref(false)
const loading = ref(true)
const error = ref('')
const billets = ref<Billet[]>([])
const stats = ref<BilletStats>({
  total_count: 0,
  valid_count: 0,
  used_count: 0,
  revenue_cents: 0,
  vendredi_places: 0,
  samedi_places: 0,
})
const busyId = ref<string | null>(null)
const editing = ref<Billet | null>(null)
const savingComment = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await fetchAdminBillets()
    billets.value = data.billets
    stats.value = data.stats
  } catch (e) {
    error.value = getErrorMessage(e, 'Impossible de charger les billets')
  } finally {
    loading.value = false
  }
}

async function markUsed(id: string) {
  busyId.value = id
  error.value = ''
  try {
    const updated = await updateAdminBillet(id, { status: 'used' })
    billets.value = billets.value.map((b) => (b.id === id ? updated : b))
    const data = await fetchAdminBillets()
    stats.value = data.stats
  } catch (e) {
    error.value = getErrorMessage(e, 'Mise à jour impossible')
  } finally {
    busyId.value = null
  }
}

function openCommentEdit(billet: Billet) {
  editing.value = billet
}

async function saveComment(comments: string) {
  const target = editing.value
  if (!target) return
  savingComment.value = true
  error.value = ''
  try {
    const updated = await updateAdminBillet(target.id, { comments })
    billets.value = billets.value.map((b) => (b.id === target.id ? updated : b))
    editing.value = null
  } catch (e) {
    error.value = getErrorMessage(e, 'Enregistrement du commentaire impossible')
  } finally {
    savingComment.value = false
  }
}

async function logout() {
  loggingOut.value = true
  try {
    await auth.logout()
    await router.replace({ name: 'lekeza-login' })
  } finally {
    loggingOut.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <AppShell :user-label="auth.displayName" :logging-out="loggingOut" @logout="logout">
    <div class="flex flex-col gap-6">
      <header class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <p class="text-xs font-semibold tracking-eyebrow text-accent-700 uppercase">Lekeza</p>
          <h1 class="font-serif text-3xl font-normal tracking-tight text-ink">Billetterie</h1>
          <p class="mt-1 text-sm text-ink-muted">Liste des billets encaissés et indicateurs clés.</p>
        </div>
        <PrimaryButton type="button" variant="outline" :loading="loading" @click="load">
          Actualiser
        </PrimaryButton>
      </header>

      <p v-if="error" class="rounded-control border border-danger-600/30 bg-danger-50 px-4 py-3 text-sm text-danger-700" role="alert">
        {{ error }}
      </p>

      <AdminStatsStrip :stats="stats" />
      <AdminBilletsTable
        :billets="billets"
        :busy-id="busyId"
        @mark-used="markUsed"
        @edit-comment="openCommentEdit"
      />
    </div>

    <AdminBilletCommentEditModal
      :open="editing !== null"
      :initial-comments="editing?.comments ?? null"
      :busy="savingComment"
      @close="editing = null"
      @save="saveComment"
    />
  </AppShell>
</template>
