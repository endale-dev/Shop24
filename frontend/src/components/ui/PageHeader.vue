<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<template>
  <div class="mb-4 flex flex-wrap items-end justify-between gap-3 border-b border-border pb-4">
    <div>
      <h1 class="text-lg font-semibold leading-6 text-text">{{ title }}</h1>
      <p v-if="subtitle" class="mt-1 max-w-2xl text-sm text-muted">{{ subtitle }}</p>
    </div>
    <slot />
  </div>
</template>
