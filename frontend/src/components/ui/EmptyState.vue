<script setup lang="ts">
import Icon from '@/components/ui/Icon.vue'

withDefaults(
  defineProps<{
    title: string
    body?: string
    icon?: 'grid' | 'box' | 'layers' | 'clipboard' | 'panel' | 'gear' | 'plus' | 'menu' | 'x' | 'logout' | 'inbox' | 'cart'
  }>(),
  { icon: 'inbox' },
)
</script>

<template>
  <div class="rounded border border-dashed border-border bg-surface px-4 py-8 text-center">
    <Icon :name="icon" class="mx-auto text-muted" />
    <h2 class="mt-2 text-sm font-semibold text-text">{{ title }}</h2>
    <p v-if="body" class="mx-auto mt-1 max-w-md text-sm text-muted">{{ body }}</p>
    <div v-if="$slots.default" class="mt-4"><slot /></div>
  </div>
</template>
