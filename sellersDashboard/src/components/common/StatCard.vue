/**
 * Statistics Card Component
 * Displays a single metric with icon, title, value, and change indicator
 */
<template>
  <BaseCard no-padding>
    <div class="px-6 py-5">
      <div class="flex items-center">
        <div v-if="icon" class="flex-shrink-0">
          <div class="p-3 rounded-lg" :class="iconBgClass">
            <component :is="icon" class="h-6 w-6" :class="iconColorClass" />
          </div>
        </div>
        <div :class="icon ? 'ml-5' : ''" class="flex-1">
          <p class="text-sm font-medium text-gray-600 truncate">
            {{ title }}
          </p>
          <div class="mt-1 flex items-baseline">
            <p class="text-2xl font-semibold text-gray-900">
              {{ value }}
            </p>
            <p v-if="change !== null && change !== undefined" class="ml-2 flex items-baseline text-sm font-medium" :class="changeColorClass">
              <component :is="changeIcon" class="self-center flex-shrink-0 h-4 w-4" :class="changeColorClass" />
              <span class="ml-1">{{ Math.abs(change) }}%</span>
            </p>
          </div>
          <p v-if="subtitle" class="mt-1 text-sm text-gray-500">
            {{ subtitle }}
          </p>
        </div>
      </div>
    </div>
  </BaseCard>
</template>

<script setup>
import { computed } from 'vue'
import { ArrowUpIcon, ArrowDownIcon } from '@heroicons/vue/24/solid'
import { BaseCard } from './index'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  subtitle: {
    type: String,
    default: ''
  },
  change: {
    type: Number,
    default: null
  },
  icon: {
    type: Object,
    default: null
  },
  iconColor: {
    type: String,
    default: 'blue',
    validator: (value) => ['blue', 'green', 'amber', 'red', 'purple'].includes(value)
  }
})

const iconColorMap = {
  blue: 'text-blue-600',
  green: 'text-green-600',
  amber: 'text-amber-600',
  red: 'text-red-600',
  purple: 'text-purple-600'
}

const iconBgMap = {
  blue: 'bg-blue-100',
  green: 'bg-green-100',
  amber: 'bg-amber-100',
  red: 'bg-red-100',
  purple: 'bg-purple-100'
}

const iconColorClass = computed(() => iconColorMap[props.iconColor])
const iconBgClass = computed(() => iconBgMap[props.iconColor])

const changeIcon = computed(() => {
  return props.change >= 0 ? ArrowUpIcon : ArrowDownIcon
})

const changeColorClass = computed(() => {
  return props.change >= 0 ? 'text-green-600' : 'text-red-600'
})
</script>
