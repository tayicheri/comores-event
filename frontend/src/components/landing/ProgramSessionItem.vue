<script setup lang="ts">
import type { ProgramSession } from './programme'

defineProps<{
  session: ProgramSession
}>()
</script>

<template>
  <li class="grid grid-cols-1 gap-x-5 gap-y-1.5 border-b border-line py-6 sm:grid-cols-program-time">
    <span class="font-serif text-lede font-medium text-terracotta">{{ session.time }}</span>
    <div>
      <strong class="inline-flex flex-wrap items-center gap-2 text-sm">
        {{ session.title }}
        <span
          v-if="session.highlight"
          class="rounded-badge bg-forest px-2 py-0.5 text-2xs font-bold tracking-label text-gold"
        >
          TEMPS FORT
        </span>
      </strong>
      <p v-if="session.description" class="mt-1.5 text-nav leading-body text-ink-muted">
        {{ session.description }}
      </p>
      <ul v-if="session.details?.length" class="mt-2.5 grid list-none gap-2 p-0">
        <li
          v-for="detail in session.details"
          :key="detail"
          class="relative pl-4 text-nav leading-body text-ink-muted before:absolute before:top-2.5 before:left-0 before:size-1.5 before:rounded-full before:bg-gold"
        >
          {{ detail }}
        </li>
      </ul>
    </div>
  </li>
</template>
