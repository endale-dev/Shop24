<script setup lang="ts">
withDefaults(defineProps<{ tone?: 'neutral' | 'good' | 'warn' | 'bad' }>(), { tone: 'neutral' })
const tones = {
  neutral: 'bg-stone-200 text-stone-700',
  good: 'bg-lime-100 text-lime-900',
  warn: 'bg-amber-100 text-amber-900',
  bad: 'bg-red-100 text-red-800',
}
</script>

<template>
  <span class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="tones[tone]"><slot /></span>
</template>
