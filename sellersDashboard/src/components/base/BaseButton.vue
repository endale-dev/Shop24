<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <!-- Loading spinner -->
    <svg
      v-if="loading"
      class="animate-spin -ml-1 mr-2 h-4 w-4"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>

    <!-- Icon (left) -->
    <component v-if="iconLeft && !loading" :is="iconLeft" :class="iconClass" />

    <!-- Content -->
    <span v-if="$slots.default" :class="{ 'opacity-0': loading }">
      <slot />
    </span>

    <!-- Icon (right) -->
    <component v-if="iconRight && !loading" :is="iconRight" :class="iconClass" />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  /**
   * Button type attribute
   */
  type: {
    type: String,
    default: 'button',
    validator: (value) => ['button', 'submit', 'reset'].includes(value)
  },
  /**
   * Button visual variant
   */
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'danger', 'ghost', 'link'].includes(value)
  },
  /**
   * Button size
   */
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  /**
   * Full width button
   */
  fullWidth: {
    type: Boolean,
    default: false
  },
  /**
   * Disabled state
   */
  disabled: {
    type: Boolean,
    default: false
  },
  /**
   * Loading state
   */
  loading: {
    type: Boolean,
    default: false
  },
  /**
   * Icon component to display on the left
   */
  iconLeft: {
    type: Object,
    default: null
  },
  /**
   * Icon component to display on the right
   */
  iconRight: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['click'])

const handleClick = (event) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}

const buttonClasses = computed(() => {
  const baseClasses = [
    'inline-flex',
    'items-center',
    'justify-center',
    'font-medium',
    'rounded-lg',
    'transition-all',
    'duration-200',
    'focus:outline-none',
    'focus:ring-2',
    'focus:ring-offset-2'
  ]

  // Size classes
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  }

  // Variant classes
  const variantClasses = {
    primary: [
      'bg-blue-600',
      'text-white',
      'hover:bg-blue-700',
      'focus:ring-blue-500',
      'disabled:bg-blue-300',
      'disabled:cursor-not-allowed'
    ],
    secondary: [
      'bg-gray-200',
      'text-gray-900',
      'hover:bg-gray-300',
      'focus:ring-gray-500',
      'disabled:bg-gray-100',
      'disabled:text-gray-400',
      'disabled:cursor-not-allowed'
    ],
    danger: [
      'bg-red-600',
      'text-white',
      'hover:bg-red-700',
      'focus:ring-red-500',
      'disabled:bg-red-300',
      'disabled:cursor-not-allowed'
    ],
    ghost: [
      'bg-transparent',
      'text-gray-700',
      'hover:bg-gray-100',
      'focus:ring-gray-500',
      'disabled:text-gray-400',
      'disabled:hover:bg-transparent',
      'disabled:cursor-not-allowed'
    ],
    link: [
      'bg-transparent',
      'text-blue-600',
      'hover:text-blue-800',
      'hover:underline',
      'focus:ring-blue-500',
      'disabled:text-blue-300',
      'disabled:cursor-not-allowed',
      'px-0'
    ]
  }

  const widthClass = props.fullWidth ? 'w-full' : ''

  return [
    ...baseClasses,
    sizeClasses[props.size],
    ...variantClasses[props.variant],
    widthClass
  ].filter(Boolean).join(' ')
})

const iconClass = computed(() => {
  const sizes = {
    sm: 'h-4 w-4',
    md: 'h-5 w-5',
    lg: 'h-6 w-6'
  }
  return sizes[props.size]
})
</script>
