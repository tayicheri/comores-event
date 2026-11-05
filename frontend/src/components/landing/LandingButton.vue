<script setup lang="ts">
withDefaults(
  defineProps<{
    href?: string
    type?: 'button' | 'submit'
    variant?: 'dark' | 'yellow'
    block?: boolean
    compact?: boolean
    disabled?: boolean
  }>(),
  { type: 'button', variant: 'dark', block: false, compact: false, disabled: false },
)

defineEmits<{ click: [MouseEvent] }>()
</script>

<template>
  <component
    :is="href ? 'a' : 'button'"
    :href="disabled ? undefined : href"
    :type="href ? undefined : type"
    :disabled="href ? undefined : disabled"
    :aria-disabled="disabled || undefined"
    class="inline-flex items-center justify-center gap-2.5 rounded-pill text-nav font-bold transition-transform duration-200 hover:-translate-y-0.5 focus-visible:outline focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-gold disabled:pointer-events-none disabled:opacity-60"
    :class="[
      compact
        ? 'min-h-btn-min-sm px-4 py-2.5 text-eyebrow md:min-h-btn-min md:px-6 md:py-btn-y md:text-nav'
        : 'min-h-btn-min px-btn-x py-btn-y',
      block ? 'w-full rounded-checkout' : '',
      variant === 'dark' ? 'bg-forest text-white hover:bg-forest-hover' : 'bg-gold text-forest-deep hover:bg-gold-hover',
    ]"
    @click="$emit('click', $event)"
  >
    <slot />
  </component>
</template>
