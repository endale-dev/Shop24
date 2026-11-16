<script setup lang="ts">
defineProps<{ title: string; open: boolean }>()
const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-40 flex items-end justify-center bg-ink/40 p-4 sm:items-center" @click.self="emit('close')">
    <div class="w-full max-w-md rounded-2xl bg-card p-5 shadow-xl">
      <h2 class="font-serif text-2xl">{{ title }}</h2>
      <div class="mt-3 text-sm text-stone-600"><slot /></div>
      <div class="mt-5 flex justify-end gap-2"><slot name="actions" /></div>
    </div>
  </div>
</template>
