<script setup lang="ts">
import { inject, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { formatMoney } from '@/features/money'
import { storefrontKey } from '@/features/storefront/context'
import { publicProduct } from '@/features/storefront-builder/api'
import type { Product } from '@/types'

const route = useRoute()
const ctx = inject(storefrontKey)
const product = ref<Product | null>(null)
const quantity = ref(1)
const error = ref('')
const added = ref(false)

watch(
  () => route.params.productId,
  async (id) => {
    if (!ctx || typeof id !== 'string') return
    error.value = ''
    added.value = false
    try {
      product.value = (await publicProduct(ctx.slug.value, id)).product
      quantity.value = 1
    } catch (err) {
      product.value = null
      error.value = err instanceof Error ? err.message : 'Product not found'
    }
  },
  { immediate: true },
)

function add() {
  if (!ctx || !product.value || product.value.stock < 1) return
  ctx.cart.value.add(
    {
      productId: product.value.id,
      name: product.value.name,
      priceCents: product.value.priceCents,
      image: product.value.images[0] || '',
      stock: product.value.stock,
    },
    quantity.value,
  )
  added.value = true
}
</script>

<template>
  <p v-if="error" class="text-sm">{{ error }}</p>
  <article v-else-if="product && ctx?.store.value" class="grid gap-8 md:grid-cols-2">
    <div class="overflow-hidden" :style="{ borderRadius: 'var(--sf-radius)', background: 'var(--sf-surface)' }">
      <img v-if="product.images[0]" :src="product.images[0]" :alt="product.name" class="aspect-[3/4] w-full object-cover" />
      <div v-else class="grid aspect-[3/4] place-items-center text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ product.name.slice(0, 1) }}</div>
    </div>
    <div>
      <p v-if="product.category" class="text-xs uppercase tracking-[0.18em]" :style="{ color: 'var(--sf-muted)' }">{{ product.category }}</p>
      <h1 class="mt-2 text-4xl md:text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ product.name }}</h1>
      <p class="mt-3 text-lg">{{ formatMoney(product.priceCents, ctx.store.value.currency) }}</p>
      <p class="mt-4 whitespace-pre-line" :style="{ color: 'var(--sf-muted)' }">{{ product.description }}</p>
      <p class="mt-4 text-sm">{{ product.stock > 0 ? `${product.stock} available` : 'Sold out' }}</p>
      <div class="mt-5 flex items-center gap-3">
        <input v-model.number="quantity" type="number" min="1" :max="product.stock" class="w-20 border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" :disabled="product.stock < 1" />
        <button type="button" class="px-5 py-2 text-sm font-semibold disabled:opacity-40" :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)', borderRadius: 'var(--sf-radius)' }" :disabled="product.stock < 1" @click="add">Add to cart</button>
      </div>
      <p v-if="added" class="mt-3 text-sm">Added. <RouterLink :to="ctx.paths.value.cart">View cart</RouterLink></p>
    </div>
  </article>
</template>
