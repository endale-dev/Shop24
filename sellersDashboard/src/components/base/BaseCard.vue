<template>
  <div :class="cardClasses">
    <!-- Header -->
    <div v-if="$slots.header || title" class="px-6 py-4 border-b border-gray-200">
      <slot name="header">
        <h3 v-if="title" class="text-lg font-semibold text-gray-900">
          {{ title }}
        </h3>
        <p v-if="subtitle" class="mt-1 text-sm text-gray-600">
          {{ subtitle }}
        </p>
      </slot>
    </div>

    <!-- Body -->
    <div :class="bodyClasses">
      <slot />
    </div>

    <!-- Footer -->
    <div v-if="$slots.footer" class="px-6 py-4 border-t border-gray-200 bg-gray-50">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  /**
   * Card title
   */
  title: {
    type: String,
    default: ''
  },
  /**
   * Card subtitle
   */
  subtitle: {
    type: String,
    default: ''
  },
  /**
   * Remove default padding from body
   */
  noPadding: {
    type: Boolean,
    default: false
  },
  /**
   * Clickable card (adds hover effect)
   */
  clickable: {
    type: Boolean,
    default: false
  },
  /**
   * Remove shadow
   */
  flat: {
    type: Boolean,
    default: false
  }
})

const cardClasses = computed(() => {
  const baseClasses = [
    'bg-white',
    'rounded-lg',
    'overflow-hidden'
  ]

  const shadowClass = props.flat ? 'border border-gray-200' : 'shadow-sm'
  const hoverClass = props.clickable ? 'hover:shadow-md transition-shadow cursor-pointer' : ''

  return [
    ...baseClasses,
    shadowClass,
    hoverClass
  ].filter(Boolean).join(' ')
})

const bodyClasses = computed(() => {
  return props.noPadding ? '' : 'px-6 py-4'
})
</script>
