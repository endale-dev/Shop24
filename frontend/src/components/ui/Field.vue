<script setup lang="ts">
defineProps<{ label: string; hint?: string; error?: string }>()
</script>

<template>
  <label class="block">
    <span class="mb-1.5 block text-sm font-medium text-text">{{ label }}</span>
    <slot />
    <p v-if="hint" class="mt-1 text-xs text-muted">{{ hint }}</p>
    <p v-if="error" class="mt-1 text-sm text-bad-fg">{{ error }}</p>
  </label>
</template>
