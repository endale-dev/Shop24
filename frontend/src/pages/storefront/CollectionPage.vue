<script setup lang="ts">
import { inject, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { publicCollection } from '@/features/collections/api'
import { formatMoney } from '@/features/money'
import { priceVaries } from '@/features/products/api'
import { storefrontKey } from '@/features/storefront/context'
import type { Collection, Product } from '@/types'

const route = useRoute()
const ctx = inject(storefrontKey)
const collection = ref<Collection | null>(null)
const products = ref<Product[]>([])
const error = ref('')

watch(
  () => [ctx?.slug.value, route.params.collectionSlug],
  async () => {
    if (!ctx || typeof route.params.collectionSlug !== 'string') return
    error.value = ''
    try {
      const result = await publicCollection(ctx.slug.value, route.params.collectionSlug)
      collection.value = result.collection
      products.value = result.products
    } catch (err) {
      collection.value = null
      products.value = []
      error.value = err instanceof Error ? err.message : 'Collection not found'
    }
  },
  { immediate: true },
)
</script>

<template>
  <p v-if="error" class="text-sm">{{ error }}</p>
  <article v-else-if="collection && ctx?.store.value">
    <div class="grid gap-8 md:grid-cols-[0.8fr_1.2fr]">
      <div class="overflow-hidden" :style="{ borderRadius: 'var(--sf-radius)', background: 'var(--sf-surface)' }">
        <img v-if="collection.imageUrl" :src="collection.imageUrl" :alt="collection.name" class="aspect-[4/3] w-full object-cover" />
        <div v-else class="grid aspect-[4/3] place-items-center text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ collection.name.slice(0, 1) }}</div>
      </div>
      <div>
        <h1 class="text-4xl md:text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ collection.name }}</h1>
        <p v-if="collection.description" class="mt-4 whitespace-pre-line" :style="{ color: 'var(--sf-muted)' }">{{ collection.description }}</p>
      </div>
    </div>
    <div v-if="products.length === 0" class="mt-10 text-sm" :style="{ color: 'var(--sf-muted)' }">Nothing in this collection is published yet.</div>
    <div v-else class="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
      <RouterLink v-for="product in products" :key="product.id" :to="ctx.paths.value.product(product.id)" class="group block focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]">
        <div class="overflow-hidden" :style="{ background: 'var(--sf-surface)', borderRadius: 'var(--sf-radius)' }">
          <img v-if="product.images[0]" :src="product.images[0]" :alt="product.name" class="aspect-square w-full object-cover" />
          <div v-else class="grid aspect-square place-items-center text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ product.name.slice(0, 1) }}</div>
        </div>
        <div class="mt-3 flex items-start justify-between gap-3">
          <h2 class="text-base font-medium">{{ product.name }}</h2>
          <p class="text-sm">{{ priceVaries(product) ? 'From ' : '' }}{{ formatMoney(product.priceCents, ctx.store.value.currency) }}</p>
        </div>
      </RouterLink>
    </div>
  </article>
</template>
