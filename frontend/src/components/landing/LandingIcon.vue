<script setup lang="ts">
defineProps<{
  name:
    | 'star'
    | 'arrow'
    | 'calendar'
    | 'pin'
    | 'ticket'
    | 'info'
    | 'check'
    | 'plus'
    | 'minus'
    | 'x'
    | 'chevron'
}>()
</script>

<template>
  <svg class="size-icon shrink-0 overflow-visible" aria-hidden="true">
    <use :href="`#i-${name}`" />
  </svg>
</template>
