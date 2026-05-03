<template>
  <section class="py-16 px-6 container mx-auto">
    <div class="text-center mb-12">
      <h2 class="text-3xl font-bold mb-4" style="color: var(--color-text-primary)">{{ title }}</h2>
      <p v-if="subtitle" class="text-lg opacity-75" style="color: var(--color-text-secondary)">{{ subtitle }}</p>
    </div>

    <div v-if="loading" class="flex justify-center py-10">
      <div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
    </div>
    
    <div v-else-if="products.length === 0" class="text-center py-10 opacity-50">
      No products found.
    </div>

    <div v-else class="grid gap-8" :class="gridClass">
      <div 
        v-for="product in products" 
        :key="product.id"
        class="group bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 border border-gray-100"
      >
        <div class="relative aspect-[4/5] bg-gray-50 overflow-hidden">
          <img 
            :src="product.images?.[0]?.url || 'https://via.placeholder.com/400x500?text=No+Image'" 
            :alt="product.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          />
          <!-- Hover Overlay -->
          <div class="absolute inset-x-0 bottom-0 p-4 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex justify-center">
             <button class="bg-white text-gray-900 font-medium py-2 px-6 rounded-full shadow-lg transform translate-y-4 group-hover:translate-y-0 transition-all duration-300 w-full">
               Quick View
             </button>
          </div>
        </div>
        
        <div class="p-5 flex flex-col gap-2">
          <h3 class="font-semibold text-lg line-clamp-1" style="color: var(--color-text-primary)">
            {{ product.title }}
          </h3>
          
          <div class="flex items-center justify-between mt-auto">
            <span v-if="showPrice" class="font-bold text-lg" style="color: var(--color-primary, #10B981)">
              ${{ product.price?.amount || '0.00' }}
            </span>
            <div v-if="showRating && product.rating" class="flex items-center gap-1 text-sm text-yellow-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 fill-current" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" /></svg>
              <span>{{ product.rating.toFixed(1) }}</span>
            </div>
          </div>
        </div>
      </div>
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
  subtitle: [String, Object],
  collectionId: [String, Object],
  columns: [String, Object, Number],
  limit: [String, Object, Number],
  showPrice: [Boolean, Object],
  showRating: [Boolean, Object]
})

const { getProp } = useResponsive()
const shopStore = useShopStore()

const title = computed(() => getProp(props.title, 'Featured Products'))
const subtitle = computed(() => getProp(props.subtitle, ''))
const columns = computed(() => String(getProp(props.columns, '4')))
const limit = computed(() => Number(getProp(props.limit, 8)))
const showPrice = computed(() => getProp(props.showPrice, true))
const showRating = computed(() => getProp(props.showRating, true))
const collectionId = computed(() => getProp(props.collectionId, ''))

const gridClass = computed(() => {
  if (columns.value === '2') return 'grid-cols-1 md:grid-cols-2'
  if (columns.value === '3') return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
  return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4'
})

const products = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  if (!shopStore.shopSlug) return
  try {
    let url = `/shops/${shopStore.shopSlug}/products`
    const params: any = { limit: limit.value }
    if (collectionId.value) {
      params.collection = collectionId.value
    }
    
    const res = await axios.get(url, { params })
    products.value = res.data.products || res.data || []
  } catch (error) {
    console.error('Failed to load products for grid', error)
  } finally {
    loading.value = false
  }
})
</script>
