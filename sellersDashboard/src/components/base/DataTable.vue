<template>
  <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 rounded-lg">
    <table class="min-w-full divide-y divide-gray-300">
      <!-- Header -->
      <thead class="bg-gray-50">
        <tr>
          <th
            v-for="column in columns"
            :key="column.key"
            scope="col"
            class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
            :class="column.headerClass"
          >
            {{ column.label }}
          </th>
          <th v-if="$slots.actions" scope="col" class="px-3 py-3.5 text-right text-sm font-semibold text-gray-900">
            Actions
          </th>
        </tr>
      </thead>
      
      <!-- Body -->
      <tbody class="divide-y divide-gray-200 bg-white">
        <!-- Empty state -->
        <tr v-if="!data || data.length === 0">
          <td :colspan="columns.length + ($slots.actions ? 1 : 0)" class="px-3 py-12">
            <EmptyState
              :icon="emptyIcon"
              :title="emptyTitle"
              :description="emptyDescription"
            />
          </td>
        </tr>
        
        <!-- Data rows -->
        <tr 
          v-else
          v-for="(item, index) in data" 
          :key="item[rowKey] || index"
          class="hover:bg-gray-50 transition-colors"
        >
          <td
            v-for="column in columns"
            :key="column.key"
            class="whitespace-nowrap px-3 py-4 text-sm text-gray-900"
            :class="column.class"
          >
            <slot :name="`cell-${column.key}`" :item="item" :value="item[column.key]">
              {{ item[column.key] }}
            </slot>
          </td>
          
          <!-- Actions column -->
          <td v-if="$slots.actions" class="whitespace-nowrap px-3 py-4 text-right text-sm font-medium">
            <slot name="actions" :item="item" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { EmptyState } from './index'

defineProps({
  columns: {
    type: Array,
    required: true,
    validator: (columns) => {
      return columns.every(col => col.key && col.label)
    }
  },
  data: {
    type: Array,
    default: () => []
  },
  rowKey: {
    type: String,
    default: 'id'
  },
  emptyIcon: {
    type: Object,
    default: null
  },
  emptyTitle: {
    type: String,
    default: 'No data'
  },
  emptyDescription: {
    type: String,
    default: 'No records found'
  }
})
</script>
