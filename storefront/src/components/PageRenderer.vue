<template>
  <div>
    <component
      v-for="block in blocks"
      :key="block.id"
      :is="resolveComponent(block.componentKey)"
      v-bind="block.props"
    >
      <PageRenderer v-if="block.children && block.children.length" :blocks="block.children" />
    </component>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import HeaderBlock from '@/components/sections/Header.vue'
import ContainerBlock from '@/components/sections/Container.vue'
import ProductGrid from '@/components/sections/ProductGrid.vue'

const props = defineProps({ blocks: { type: Array, default: () => [] } })

const registry = {
  'Header': HeaderBlock,
  'Container': ContainerBlock,
  'ProductGrid': ProductGrid,
}

function resolveComponent(key) {
  return registry[key] || ContainerBlock
}
</script>
