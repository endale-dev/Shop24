<script setup lang="ts">
defineProps<{ title: string; body?: string }>()
</script>

<template>
  <div class="rounded-2xl border border-dashed border-line bg-card px-6 py-12 text-center">
    <h2 class="font-serif text-2xl text-ink">{{ title }}</h2>
    <p v-if="body" class="mx-auto mt-2 max-w-md text-sm text-stone-600">{{ body }}</p>
    <div v-if="$slots.default" class="mt-5"><slot /></div>
  </div>
</template>
