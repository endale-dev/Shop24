<!-- Skeleton Loader Component -->
<template>
  <div class="animate-pulse">
    <!-- Card Skeleton -->
    <div v-if="type === 'card'" class="bg-white rounded-lg shadow-sm p-6 space-y-3">
      <div class="h-4 bg-gray-200 rounded w-3/4"></div>
      <div class="h-4 bg-gray-200 rounded w-1/2"></div>
      <div class="h-8 bg-gray-200 rounded w-full"></div>
    </div>

    <!-- Table Row Skeleton -->
    <div v-else-if="type === 'table-row'" class="flex items-center space-x-4 p-4">
      <div class="h-10 w-10 bg-gray-200 rounded"></div>
      <div class="flex-1 space-y-2">
        <div class="h-4 bg-gray-200 rounded w-3/4"></div>
        <div class="h-3 bg-gray-200 rounded w-1/2"></div>
      </div>
      <div class="h-4 w-16 bg-gray-200 rounded"></div>
    </div>

    <!-- Text Skeleton -->
    <div v-else-if="type === 'text'" class="space-y-2">
      <div v-for="i in lines" :key="i" class="h-4 bg-gray-200 rounded" :class="i === lines ? 'w-2/3' : 'w-full'"></div>
    </div>

    <!-- Avatar Skeleton -->
    <div v-else-if="type === 'avatar'" :class="avatarSizeClass" class="bg-gray-200 rounded-full"></div>

    <!-- Custom Skeleton -->
    <div v-else>
      <slot />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'custom',
    validator: (value) => ['card', 'table-row', 'text', 'avatar', 'custom'].includes(value)
  },
  lines: {
    type: Number,
    default: 3
  },
  avatarSize: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  }
})

const avatarSizeClass = computed(() => {
  const sizes = {
    sm: 'h-8 w-8',
    md: 'h-12 w-12',
    lg: 'h-16 w-16'
  }
  return sizes[props.avatarSize]
})
</script>
