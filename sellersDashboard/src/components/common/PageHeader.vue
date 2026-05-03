<!-- Enhanced Page Header Component -->
<template>
  <div class="mb-6 sm:mb-8 animate-slide-down">
    <div class="sm:flex sm:items-center sm:justify-between">
      <!-- Title Section -->
      <div class="flex-1">
        <div class="flex items-center space-x-3 mb-2">
          <!-- Icon -->
          <div 
            v-if="icon" 
            class="flex-shrink-0 p-2.5 rounded-lg transition-transform hover:scale-105"
            :class="iconBgClass"
          >
            <component :is="icon" class="h-6 w-6" :class="iconColorClass" />
          </div>
          
          <!-- Title -->
          <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 tracking-tight">
            {{ title }}
          </h1>
        </div>
        
        <!-- Subtitle/Description -->
        <p v-if="subtitle" class="mt-1 text-sm sm:text-base text-gray-600 max-w-3xl">
          {{ subtitle }}
        </p>
        
        <!-- Breadcrumbs Slot -->
        <div v-if="$slots.breadcrumbs" class="mt-2">
          <slot name="breadcrumbs" />
        </div>
      </div>

      <!-- Actions -->
      <div v-if="$slots.actions" class="mt-4 sm:mt-0 sm:ml-6 flex-shrink-0">
        <slot name="actions" />
      </div>
    </div>

    <!-- Stats or Meta Info -->
    <div v-if="$slots.meta" class="mt-4 pt-4 border-t border-gray-200">
      <slot name="meta" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  subtitle: {
    type: String,
    default: ''
  },
  icon: {
    type: Object,
    default: null
  },
  iconColor: {
    type: String,
    default: 'blue',
    validator: (value) => ['blue', 'green', 'amber', 'red', 'purple', 'gray'].includes(value)
  }
})

const colorMap = {
  blue: { bg: 'bg-blue-100', text: 'text-blue-600' },
  green: { bg: 'bg-green-100', text: 'text-green-600' },
  amber: { bg: 'bg-amber-100', text: 'text-amber-600' },
  red: { bg: 'bg-red-100', text: 'text-red-600' },
  purple: { bg: 'bg-purple-100', text: 'text-purple-600' },
  gray: { bg: 'bg-gray-100', text: 'text-gray-600' }
}

const iconBgClass = computed(() => colorMap[props.iconColor].bg)
const iconColorClass = computed(() => colorMap[props.iconColor].text)
</script>
