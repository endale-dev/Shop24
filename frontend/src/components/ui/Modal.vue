<script setup lang="ts">
defineProps<{ title: string; open: boolean }>()
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-40 flex items-end justify-center bg-black/40 p-4 sm:items-center" @click.self="emit('close')">
    <div class="w-full max-w-md rounded border border-border bg-surface p-4 shadow-[0_8px_24px_rgba(17,20,24,0.12)]">
      <h2 class="text-sm font-semibold text-text">{{ title }}</h2>
      <div class="mt-3 text-sm text-muted"><slot /></div>
      <div class="mt-4 flex justify-end gap-2"><slot name="actions" /></div>
    </div>
  </div>
</template>
