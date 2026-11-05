<script setup lang="ts">
import type { CardNetwork } from '@/lib/stripe-elements'

defineProps<{
  brand: CardNetwork
  flipped: boolean
}>()

const faceBody =
  'relative flex h-full box-border flex-col justify-between overflow-hidden rounded-card border border-white/15 bg-gradient-to-br from-forest-hover to-forest-deep p-[clamp(14px,4.2vw,24px)] text-surface [backface-visibility:hidden] [&_*]:[backface-visibility:hidden]'
</script>

<template>
  <div
    class="mx-auto h-auto w-[min(100%,22rem)] shrink-0 justify-self-center [aspect-ratio:85.6/53.98] [perspective:1000px]"
  >
    <div
      class="relative size-full transition-transform duration-[650ms] ease-[cubic-bezier(0.34,1.56,0.64,1)] [transform-style:preserve-3d]"
      :class="{ '[transform:rotateY(180deg)]': flipped }"
    >
      <div class="absolute inset-0 size-full [backface-visibility:hidden] [transform:rotateY(0deg)_translateZ(1px)] [transform-style:flat]">
        <div :class="faceBody">
          <div
            class="pointer-events-none absolute -top-10 -right-10 size-[140px] rounded-full bg-[radial-gradient(circle,color-mix(in_srgb,var(--color-gold)_28%,transparent)_0%,transparent_70%)]"
          />
          <div
            class="pointer-events-none absolute -bottom-7.5 -left-7.5 size-[120px] rounded-full bg-[radial-gradient(circle,color-mix(in_srgb,var(--color-gold-soft)_22%,transparent)_0%,transparent_70%)]"
          />
          <div class="relative z-[1] flex min-h-0 flex-1 flex-col justify-between">
            <div class="flex items-start justify-between">
              <div
                class="h-9 w-12 rounded-md bg-gradient-to-br from-gold-soft via-gold to-ochre"
                aria-hidden="true"
              />
              <div class="text-right">
                <span class="font-serif text-sm font-bold italic">{{ brand }}</span>
                <small class="block text-2xs tracking-[0.16em] text-gold-soft uppercase">
                  Comores Event
                </small>
              </div>
            </div>
            <div
              class="m-0 min-w-0 overflow-hidden font-mono text-[clamp(13px,4.2vw,18px)] tracking-[0.12em] text-ellipsis whitespace-nowrap"
            >
              <slot name="number" />
            </div>
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0 flex-1">
                <span class="block text-2xs tracking-[0.16em] text-gold-soft uppercase">Titulaire</span>
                <slot name="holder" />
              </div>
              <div class="min-w-[4.5rem] text-right">
                <span class="block text-2xs tracking-[0.16em] text-gold-soft uppercase">Expire</span>
                <div class="mt-0.5 min-w-[4.25rem]">
                  <slot name="expiry" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div
        class="absolute inset-0 size-full [backface-visibility:hidden] [transform:rotateY(180deg)_translateZ(1px)] [transform-style:flat]"
      >
        <div :class="faceBody">
          <div class="absolute inset-x-0 top-4 m-0 h-11 bg-ink" />
          <div class="relative pt-[52px] text-left">
            <span class="text-caption text-gold-soft">Code de sécurité</span>
            <div class="mt-1.5 rounded-lg bg-surface px-3 py-2 text-right font-mono font-bold text-ink">
              <slot name="cvc" />
            </div>
            <p class="mt-2 text-2xs leading-relaxed text-gold-muted">
              Paiement sécurisé. Aucune donnée sensible n’est enregistrée en clair.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
