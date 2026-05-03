<template>
  <div v-if="shopStore.isLoading && !shopStore.shop" class="min-h-screen flex items-center justify-center">
    <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-gray-900"></div>
  </div>
  
  <div v-else-if="shopStore.error" class="min-h-screen flex flex-col items-center justify-center p-4">
    <h1 class="text-2xl font-bold text-red-600 mb-2">Store Unavailable</h1>
    <p class="text-gray-600">{{ shopStore.error }}</p>
  </div>
  
  <div v-else class="dynamic-page" :style="cssVariables">
    <!-- Header Blocks -->
    <DynamicBlock v-for="block in headerBlocks" :key="block.id" :block="block" />
    
    <!-- Main Page Content Blocks -->
    <main class="min-h-[50vh]">
      <DynamicBlock v-for="block in pageContentBlocks" :key="block.id" :block="block" />
      <div v-if="pageContentBlocks.length === 0" class="text-center py-20 text-gray-500">
        No content available for this page
      </div>
    </main>

    <!-- Footer Blocks -->
    <DynamicBlock v-for="block in footerBlocks" :key="block.id" :block="block" />
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useShopStore } from '../stores/shop'
import DynamicBlock from '../components/DynamicBlock.vue'

const route = useRoute()
const shopStore = useShopStore()

onMounted(async () => {
  await shopStore.initialize()
  await shopStore.loadPage(route.path)
})

watch(() => route.path, async (newPath) => {
  if (shopStore.shop) {
    await shopStore.loadPage(newPath)
  }
})

const headerBlocks = computed(() => {
  return shopStore.theme?.headerBlocks || []
})

const footerBlocks = computed(() => {
  return shopStore.theme?.footerBlocks || []
})

const pageContentBlocks = computed(() => {
  if (shopStore.currentPage?.blocks) {
    return shopStore.currentPage.blocks
  }
  return [] // fallback
})

const cssVariables = computed(() => {
  if (!shopStore.theme?.colors) return {}
  const colors = shopStore.theme.colors
  return {
    '--color-primary': colors.primary,
    '--color-secondary': colors.secondary,
    '--color-accent': colors.accent,
    '--color-background': colors.background,
    '--color-surface': colors.surface,
    '--color-text-primary': colors.textPrimary,
    '--color-text-secondary': colors.textSecondary,
    '--color-border': colors.border,
  }
})
</script>

<style>
.dynamic-page {
  background-color: var(--color-background, #ffffff);
  color: var(--color-text-primary, #111827);
  font-family: var(--font-body, 'Inter', sans-serif);
}
</style>
