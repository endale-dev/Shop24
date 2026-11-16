<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Modal from '@/components/ui/Modal.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { formatMoney } from '@/features/money'
import { deleteProduct, listProducts, priceVaries } from '@/features/products/api'
import { useShopStore } from '@/stores/shop'
import type { Product } from '@/types'

const shop = useShopStore()
const items = ref<Product[]>([])
const query = ref('')
const status = ref('all')
const error = ref('')
const loading = ref(true)
const pendingDelete = ref<Product | null>(null)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await listProducts({ q: query.value, status: status.value })
    items.value = result.items
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load products'
  } finally {
    loading.value = false
  }
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await deleteProduct(pendingDelete.value.id)
    pendingDelete.value = null
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not delete the product'
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Products" subtitle="Prices, photos, and the stock customers can buy.">
      <Button to="/products/new">New product</Button>
    </PageHeader>
    <form class="mb-4 flex flex-wrap gap-2" @submit.prevent="load">
      <input v-model="query" :class="controlClass" class="max-w-xs" placeholder="Search name or SKU" />
      <select v-model="status" :class="controlClass" class="max-w-[10rem]">
        <option value="all">All</option>
        <option value="published">Published</option>
        <option value="draft">Drafts</option>
      </select>
      <Button type="submit" variant="secondary">Filter</Button>
    </form>
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <p v-if="loading" class="text-sm text-stone-500">Loading products…</p>
    <EmptyState v-else-if="items.length === 0" title="No products yet" body="Add something you actually sell. Drafts stay off the public shop.">
      <Button to="/products/new">Add a product</Button>
    </EmptyState>
    <div v-else class="overflow-hidden rounded-2xl border border-line bg-card">
      <table class="w-full text-left text-sm">
        <thead class="border-b border-line text-xs uppercase tracking-wide text-stone-500">
          <tr><th class="px-4 py-3">Product</th><th class="px-4 py-3">Price</th><th class="px-4 py-3">Stock</th><th class="px-4 py-3">Status</th><th class="px-4 py-3"></th></tr>
        </thead>
        <tbody>
          <tr v-for="product in items" :key="product.id" class="border-b border-line last:border-0">
            <td class="px-4 py-3">
              <p class="font-medium">{{ product.name }}</p>
              <p class="text-xs text-stone-500">{{ product.variants?.length > 1 ? `${product.variants.length} versions` : product.sku }}<span v-if="product.category"> · {{ product.category }}</span></p>
            </td>
            <td class="px-4 py-3">{{ priceVaries(product) ? 'From ' : '' }}{{ formatMoney(product.priceCents, shop.shop?.currency) }}</td>
            <td class="px-4 py-3"><Badge :tone="product.stock === 0 ? 'bad' : product.lowStock ? 'warn' : 'good'">{{ product.stock }}</Badge></td>
            <td class="px-4 py-3">{{ product.published ? 'Published' : 'Draft' }}</td>
            <td class="px-4 py-3 text-right">
              <RouterLink class="mr-3 font-medium text-rust" :to="`/products/${product.id}`">Edit</RouterLink>
              <button class="text-stone-500" type="button" @click="pendingDelete = product">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <Modal title="Delete this product?" :open="!!pendingDelete" @close="pendingDelete = null">
      <p>{{ pendingDelete?.name }} will disappear from the shop. Past orders keep the name and price from checkout.</p>
      <template #actions>
        <Button variant="ghost" @click="pendingDelete = null">Keep</Button>
        <Button variant="danger" @click="confirmDelete">Delete</Button>
      </template>
    </Modal>
  </div>
</template>
