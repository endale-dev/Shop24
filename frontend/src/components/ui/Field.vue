<script setup lang="ts">
defineProps<{ label: string; hint?: string; error?: string }>()
</script>

<template>
  <label class="block">
    <span class="mb-1.5 block text-sm font-medium text-stone-700">{{ label }}</span>
    <slot />
    <p v-if="hint" class="mt-1 text-xs text-stone-500">{{ hint }}</p>
    <p v-if="error" class="mt-1 text-sm text-red-700">{{ error }}</p>
  </label>
</template>
