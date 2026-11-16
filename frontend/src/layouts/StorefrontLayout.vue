<script setup lang="ts">
import { computed, provide, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { publicCollections } from '@/features/collections/api'
import { storefrontKey } from '@/features/storefront/context'
import { publicProducts, publicStore } from '@/features/storefront-builder/api'
import { storePaths } from '@/features/stores/host'
import { useCart } from '@/features/storefront/cart'
import StorefrontView from '@/components/storefront/StorefrontView.vue'
import type { Collection, Product, PublicStore } from '@/types'

const props = defineProps<{ slug?: string }>()
const route = useRoute()
const slug = computed(() => props.slug || String(route.params.slug || ''))
const paths = computed(() => storePaths(slug.value))
const store = ref<PublicStore | null>(null)
const products = ref<Product[]>([])
const collections = ref<Collection[]>([])
const loading = ref(true)
const error = ref('')
const cart = computed(() => useCart(slug.value))

async function load(next: string) {
  loading.value = true
  error.value = ''
  try {
    const [shop, catalog, groups] = await Promise.all([
      publicStore(next),
      publicProducts(next),
      publicCollections(next).catch(() => ({ items: [] as Collection[] })),
    ])
    store.value = shop
    products.value = catalog.items
    collections.value = groups.items
    document.title = shop.name
  } catch (err) {
    store.value = null
    products.value = []
    collections.value = []
    error.value = err instanceof Error ? err.message : 'This shop is not available'
  } finally {
    loading.value = false
  }
}

watch(slug, (value) => void load(value), { immediate: true })
provide(storefrontKey, { slug, store, products, paths, cart })
</script>

<template>
  <div v-if="loading" class="grid min-h-screen place-items-center bg-paper text-sm text-stone-500">Opening the shop…</div>
  <div v-else-if="error || !store" class="grid min-h-screen place-items-center bg-paper px-6 text-center">
    <div>
      <p class="font-serif text-4xl">Shop unavailable</p>
      <p class="mt-2 text-sm text-stone-600">{{ error || 'This store is not published.' }}</p>
    </div>
  </div>
  <StorefrontView v-else :store="store" :products="products" :collections="collections" :paths="paths" :cart-count="cart.count.value" :show-sections="route.name === 'store-home'">
    <RouterView />
  </StorefrontView>
</template>
