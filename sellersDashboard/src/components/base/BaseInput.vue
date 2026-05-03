<template>
  <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1">
    {{ label }}
    <span v-if="required" class="text-red-500 ml-0.5">*</span>
  </label>
  
  <div class="relative">
    <input
      :id="id"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :required="required"
      :class="inputClasses"
      @input="handleInput"
      @blur="handleBlur"
      @focus="handleFocus"
    />
    
    <!-- Icon (left) -->
    <div v-if="iconLeft" class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <component :is="iconLeft" class="h-5 w-5 text-gray-400" />
    </div>
    
    <!-- Icon (right) -->
    <div v-if="iconRight" class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
      <component :is="iconRight" class="h-5 w-5 text-gray-400" />
    </div>
  </div>
  
  <!-- Error message -->
  <p v-if="error && showError" class="mt-1 text-sm text-red-600">
    {{ error }}
  </p>
  
  <!-- Help text -->
  <p v-else-if="helpText" class="mt-1 text-sm text-gray-500">
    {{ helpText }}
  </p>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  id: {
    type: String,
    default: () => `input-${Math.random().toString(36).substr(2, 9)}`
  },
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  helpText: {
    type: String,
    default: ''
  },
  iconLeft: {
    type: Object,
    default: null
  },
  iconRight: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'blur', 'focus'])

const showError = ref(false)

const inputClasses = computed(() => {
  const baseClasses = [
    'block',
    'w-full',
    'rounded-md',
    'border',
    'shadow-sm',
    'sm:text-sm',
    'transition-colors',
    'focus:outline-none',
    'focus:ring-2',
    'focus:ring-offset-0'
  ]
  
  const paddingClass = props.iconLeft ? 'pl-10 pr-3 py-2' : props.iconRight ? 'pl-3 pr-10 py-2' : 'px-3 py-2'
  
  const stateClasses = props.error && showError.value
    ? 'border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500'
    : props.disabled
    ? 'bg-gray-100 border-gray-300 text-gray-500 cursor-not-allowed'
    : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'
  
  return [
    ...baseClasses,
    paddingClass,
    stateClasses
  ].join(' ')
})

function handleInput(event) {
  emit('update:modelValue', event.target.value)
}

function handleBlur(event) {
  showError.value = true
  emit('blur', event)
}

function handleFocus(event) {
  emit('focus', event)
}
</script>
