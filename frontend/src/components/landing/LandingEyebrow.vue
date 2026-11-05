<script setup lang="ts">
withDefaults(
  defineProps<{
    tone?: 'default' | 'gold' | 'corail'
  }>(),
  { tone: 'default' },
)
</script>

<template>
  <p
    class="mb-5 text-eyebrow font-bold tracking-eyebrow"
    :class="{
      'text-ochre': tone === 'default',
      'text-gold': tone === 'gold',
      'text-corail': tone === 'corail',
    }"
  >
    <slot />
  </p>
</template>
