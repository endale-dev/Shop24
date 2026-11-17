<script setup lang="ts">
withDefaults(defineProps<{ tone?: 'neutral' | 'good' | 'warn' | 'bad' }>(), { tone: 'neutral' })
const tones = {
  neutral: 'bg-bg text-text',
  good: 'bg-good-bg text-good-fg',
  warn: 'bg-warn-bg text-warn-fg',
  bad: 'bg-bad-bg text-bad-fg',
}
</script>

<template>
  <span class="inline-flex rounded px-1.5 py-0.5 text-xs font-medium" :class="tones[tone]"><slot /></span>
</template>
