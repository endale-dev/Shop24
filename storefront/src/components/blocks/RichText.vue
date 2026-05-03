<template>
  <section class="py-16 px-6" :style="{ backgroundColor: bgColor }">
    <div class="container mx-auto max-w-3xl prose prose-lg" :class="alignmentClass">
      <h2 v-if="heading" class="mb-6 font-bold text-3xl md:text-4xl" style="color: var(--color-text-primary)">
        {{ heading }}
      </h2>
      <div 
        class="text-gray-700 leading-relaxed"
        style="color: var(--color-text-secondary)"
        v-html="formattedContent"
      ></div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useResponsive } from '../../composables/useResponsive'

const props = defineProps({
  heading: [String, Object],
  content: [String, Object],
  alignment: [String, Object],
  backgroundColor: [String, Object],
})

const { getProp } = useResponsive()

const heading = computed(() => getProp(props.heading, ''))
const content = computed(() => getProp(props.content, 'Enter your content here...'))
const bgColor = computed(() => getProp(props.backgroundColor, 'transparent'))
const alignment = computed(() => getProp(props.alignment, 'left'))

const alignmentClass = computed(() => {
  if (alignment.value === 'center') return 'text-center mx-auto'
  if (alignment.value === 'right') return 'text-right ms-auto me-0'
  return 'text-left' // default
})

// Simple newline to <br> formatter if plain text
const formattedContent = computed(() => {
  if (!content.value) return ''
  // If it looks like HTML, return as is
  if (content.value.includes('<p>') || content.value.includes('<br>')) {
    return content.value
  }
  // Otherwise replace newlines with breaks
  return content.value.replace(/\n/g, '<br />')
})
</script>
