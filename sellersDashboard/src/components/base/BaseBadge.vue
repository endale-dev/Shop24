<template>
  <span :class="badgeClasses">
    <slot />
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  /**
   * Badge variant/color
   */
  variant: {
    type: String,
    default: 'neutral',
    validator: (value) => ['primary', 'success', 'warning', 'danger', 'neutral', 'info'].includes(value)
  },
  /**
   * Badge size
   */
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  /**
   * Rounded (pill shape)
   */
  rounded: {
    type: Boolean,
    default: false
  },
  /**
   * Dot indicator
   */
  dot: {
    type: Boolean,
    default: false
  }
})

const badgeClasses = computed(() => {
  const baseClasses = [
    'inline-flex',
    'items-center',
    'font-medium',
    'gap-1.5'
  ]

  // Size classes
  const sizeClasses = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-2.5 py-0.5 text-sm',
    lg: 'px-3 py-1 text-base'
  }

  // Variant classes
  const variantClasses = {
    primary: 'bg-blue-100 text-blue-800',
    success: 'bg-green-100 text-green-800',
    warning: 'bg-amber-100 text-amber-800',
    danger: 'bg-red-100 text-red-800',
    neutral: 'bg-gray-100 text-gray-800',
    info: 'bg-cyan-100 text-cyan-800'
  }

  const roundedClass = props.rounded ? 'rounded-full' : 'rounded'

  return [
    ...baseClasses,
    sizeClasses[props.size],
    variantClasses[props.variant],
    roundedClass
  ].filter(Boolean).join(' ')
})
</script>
