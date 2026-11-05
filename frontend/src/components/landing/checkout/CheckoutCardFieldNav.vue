<script setup lang="ts">
export type CheckoutCardFieldKey = 'holder' | 'number' | 'expiry' | 'cvc'

defineProps<{
  active: CheckoutCardFieldKey | null
}>()

defineEmits<{
  select: [field: CheckoutCardFieldKey]
}>()

const fields: { key: CheckoutCardFieldKey; label: string }[] = [
  { key: 'holder', label: 'Nom' },
  { key: 'number', label: 'Numéro' },
  { key: 'expiry', label: 'Expiration' },
  { key: 'cvc', label: 'CVC' },
]
</script>

<template>
  <div class="flex flex-wrap justify-center gap-2" role="group" aria-label="Champs de la carte">
    <button
      v-for="field in fields"
      :key="field.key"
      type="button"
      class="m-0 inline-flex min-h-10 cursor-pointer items-center gap-2 rounded-full border px-3.5 py-2 transition-colors"
      :class="
        active === field.key
          ? 'border-transparent bg-forest text-white'
          : 'border-line bg-white text-ink hover:border-gold'
      "
      :aria-pressed="active === field.key"
      @click="$emit('select', field.key)"
    >
      <span class="text-caption leading-none font-bold">{{ field.label }}</span>
    </button>
  </div>
</template>
