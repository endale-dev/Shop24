<template>
  <section class="py-16 px-6 container mx-auto">
    <div class="text-center mb-12">
      <h2 class="text-3xl font-bold mb-4" style="color: var(--color-text-primary)">{{ title }}</h2>
    </div>

    <div v-if="loading" class="flex justify-center py-10">
      <div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
    </div>
    
    <div v-else-if="collections.length === 0" class="text-center py-10 opacity-50">
      No collections found.
    </div>

    <div v-else class="grid gap-6" :class="gridClass">
      <router-link 
        v-for="collection in collections" 
        :key="collection.id"
        :to="`/collections/${collection.handle || collection.id}`"
        class="group relative overflow-hidden rounded-2xl aspect-[3/2] shadow-sm hover:shadow-xl transition-all duration-300 bg-gray-100"
      >
        <img 
          :src="collection.image?.url || 'https://via.placeholder.com/600x400?text=No+Image'" 
          :alt="collection.title"
          class="absolute inset-0 w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
        />
        <!-- Gradient Overlay -->
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>
        
        <div class="absolute inset-0 p-6 flex flex-col justify-end">
          <h3 class="text-2xl font-bold text-white mb-2 transform translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
            {{ collection.title }}
          </h3>
          <p 
            v-if="showDescription && collection.description" 
            class="text-gray-200 line-clamp-2 text-sm opacity-0 group-hover:opacity-100 transform translate-y-4 group-hover:translate-y-0 transition-all duration-300 delay-75"
          >
            {{ collection.description }}
          </p>
        </div>
      </router-link>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useResponsive } from '../../composables/useResponsive'
import { useShopStore } from '../../stores/shop'
import axios from 'axios'

const props = defineProps({
  title: [String, Object],
  columns: [String, Object, Number],
  showDescription: [Boolean, Object]
})

const { getProp } = useResponsive()
const shopStore = useShopStore()

const title = computed(() => getProp(props.title, 'Shop by Category'))
const columns = computed(() => String(getProp(props.columns, '3')))
const showDescription = computed(() => getProp(props.showDescription, true))

const gridClass = computed(() => {
  if (columns.value === '2') return 'grid-cols-1 md:grid-cols-2'
  if (columns.value === '3') return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
  return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4'
})

const collections = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  if (!shopStore.shopSlug) return
  try {
    const res = await axios.get(`/shops/${shopStore.shopSlug}/collections`)
    collections.value = res.data.collections || res.data || []
  } catch (error) {
    console.error('Failed to load collections', error)
  } finally {
    loading.value = false
  }
})
</script>
