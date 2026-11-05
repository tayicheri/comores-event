<script setup lang="ts">
withDefaults(
  defineProps<{
    as?: 'h1' | 'h2' | 'h3'
    size?: 'display' | 'section' | 'univers' | 'program' | 'ticket' | 'heritage' | 'practical' | 'faq' | 'footer' | 'panel' | 'card' | 'modal' | 'subtitle'
  }>(),
  { as: 'h2', size: 'section' },
)

const sizeClass: Record<string, string> = {
  display: 'font-serif text-display-md tracking-display font-normal md:text-display',
  section: 'font-serif text-section-md tracking-section font-normal lg:text-section',
  univers: 'font-serif text-section-univers tracking-section font-normal',
  program: 'font-serif text-section-md tracking-section font-normal lg:text-section-program',
  ticket: 'font-serif text-section-md tracking-section font-normal lg:text-section-ticket',
  heritage: 'font-serif text-section-md tracking-section font-normal leading-heritage lg:text-section-heritage',
  practical: 'font-serif text-section-md tracking-section font-normal lg:text-section-practical',
  faq: 'font-serif text-section-md tracking-section font-normal lg:text-section-faq',
  footer: 'font-serif text-section-footer tracking-section font-normal',
  panel: 'font-serif text-panel font-normal md:text-price',
  card: 'font-serif text-card tracking-card font-normal',
  modal: 'font-serif text-modal font-normal',
  subtitle: 'font-serif text-subtitle font-normal md:text-card',
}
</script>

<template>
  <component :is="as" class="m-0" :class="sizeClass[size]">
    <slot />
  </component>
</template>
