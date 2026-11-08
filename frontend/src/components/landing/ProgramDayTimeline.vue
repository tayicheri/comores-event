<script setup lang="ts">
import ProgramSessionItem from './ProgramSessionItem.vue'
import type { ProgramDay } from './programme'

defineProps<{
  day: ProgramDay
}>()
</script>

<template>
  <div>
    <h3 class="mb-2 text-eyebrow font-bold tracking-eyebrow text-terracotta">
      {{ day.label }} {{ day.num }} {{ day.month }}
    </h3>
    <ol class="m-0 list-none border-t border-line p-0">
      <ProgramSessionItem v-for="session in day.sessions" :key="`${session.time}-${session.title}`" :session="session" />
    </ol>
  </div>
</template>
