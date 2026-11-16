<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<template>
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <h1 class="font-serif text-3xl tracking-tight text-ink md:text-4xl">{{ title }}</h1>
      <p v-if="subtitle" class="mt-1 max-w-2xl text-sm text-stone-600">{{ subtitle }}</p>
    </div>
    <slot />
  </div>
</template>
