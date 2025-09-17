<template>
  <div :style="gridStyle">
    <div v-for="p in products" :key="p._id" style="padding:8px">
      <div style="border:1px solid #eee;padding:8px">{{ p.title }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const props = defineProps({ columns: { type: Number, default: 3 }, limit: { type: Number, default: 12 }, query: { type: Object, default: () => ({}) } })
const products = ref([])

const gridStyle = {
  display: 'grid',
  gridTemplateColumns: `repeat(${props.columns}, 1fr)`,
  gap: '16px'
}

onMounted(async () => {
  try {
    // Simple API call; storefront proxies /api
    const res = await fetch(`/shops?api=products&limit=${props.limit}`)
    const data = await res.json()
    products.value = data.items || data
  } catch (e) {
    console.warn('ProductGrid fetch failed', e)
  }
})
</script>
