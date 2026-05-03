<template>
  <component 
    :is="resolvedComponent" 
    v-bind="block.props" 
    :class="block.props?.className"
    :style="block.props?.style"
  >
    <template v-if="block.children && block.children.length > 0">
      <DynamicBlock 
        v-for="child in block.children" 
        :key="child.id" 
        :block="child" 
      />
    </template>
  </component>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'

const props = defineProps({
  block: {
    type: Object,
    required: true
  }
})

// Registry of available block components
const componentRegistry: Record<string, any> = {
  // Built-in HTML elements
  'div': 'div',
  'span': 'span',
  'p': 'p',
  'h1': 'h1',
  'h2': 'h2',
  'h3': 'h3',
  'img': 'img',
  'button': 'button',
  'a': 'a',
  
  // Custom Components (Async loaded)
  'HeroBanner': defineAsyncComponent(() => import('./blocks/HeroBanner.vue')),
  'ProductGrid': defineAsyncComponent(() => import('./blocks/ProductGrid.vue')),
  'CollectionGrid': defineAsyncComponent(() => import('./blocks/CollectionGrid.vue')),
  'Features': defineAsyncComponent(() => import('./blocks/Features.vue')),
  'Newsletter': defineAsyncComponent(() => import('./blocks/Newsletter.vue')),
  'RichText': defineAsyncComponent(() => import('./blocks/RichText.vue')),
}

const resolvedComponent = computed(() => {
  // SellersDashboard builder uses `type` or componentKey. Let's handle both.
  const key = props.block.componentKey || props.block.type
  return componentRegistry[key] || 'div' // Fallback to div if component not found
})
</script>
