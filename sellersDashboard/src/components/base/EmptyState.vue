<template>
  <div class="text-center py-12">
    <!-- Icon -->
    <component v-if="icon" :is="icon" class="mx-auto h-12 w-12 text-gray-400" />
    
    <!-- Title -->
    <h3 v-if="title" class="mt-4 text-lg font-medium text-gray-900">
      {{ title }}
    </h3>
    
    <!-- Description -->
    <p v-if="description" class="mt-2 text-sm text-gray-600 max-w-sm mx-auto">
      {{ description }}
    </p>
    
    <!-- Action Button -->
    <div v-if="$slots.action || actionText" class="mt-6">
      <slot name="action">
        <BaseButton v-if="actionText" @click="$emit('action')" variant="primary">
          {{ actionText }}
        </BaseButton>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { BaseButton } from './index'

defineProps({
  icon: {
    type: Object,
    default: null
  },
  title: {
    type: String,
    default: 'No data available'
  },
  description: {
    type: String,
    default: ''
  },
  actionText: {
    type: String,
    default: ''
  }
})

defineEmits(['action'])
</script>
