<script setup lang="ts">
import { computed, inject, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { formatMoney } from '@/features/money'
import { storefrontKey } from '@/features/storefront/context'
import { publicProduct } from '@/features/storefront-builder/api'
import type { Product, Variant } from '@/types'

const route = useRoute()
const ctx = inject(storefrontKey)
const product = ref<Product | null>(null)
const selected = ref<Record<string, string>>({})
const quantity = ref(1)
const error = ref('')
const added = ref(false)

const variant = computed<Variant | null>(() => {
  const item = product.value
  if (!item?.variants?.length) return null
  if (item.optionNames.length === 0) return item.variants[0]
  return item.variants.find((candidate) => item.optionNames.every((name) => candidate.options.find((option) => option.name === name)?.value === selected.value[name])) ?? null
})

const image = computed(() => variant.value?.imageUrl || product.value?.images[0] || '')

function valuesFor(name: string) {
  const values = new Set<string>()
  for (const candidate of product.value?.variants ?? []) {
    const value = candidate.options.find((option) => option.name === name)?.value
    if (value) values.add(value)
  }
  return [...values]
}

function choose(name: string, value: string) {
  selected.value = { ...selected.value, [name]: value }
  quantity.value = 1
  added.value = false
}

watch(
  () => route.params.productId,
  async (id) => {
    if (!ctx || typeof id !== 'string') return
    error.value = ''
    added.value = false
    try {
      product.value = (await publicProduct(ctx.slug.value, id)).product
      const first = product.value.variants?.[0]
      selected.value = {}
      for (const option of first?.options ?? []) selected.value[option.name] = option.value
      quantity.value = 1
    } catch (err) {
      product.value = null
      error.value = err instanceof Error ? err.message : 'Product not found'
    }
  },
  { immediate: true },
)

function add() {
  if (!ctx || !product.value || !variant.value || variant.value.stock < 1) return
  ctx.cart.value.add(
    {
      productId: product.value.id,
      variantId: variant.value.id,
      name: product.value.name,
      optionLabel: variant.value.options.map((option) => option.value).filter(Boolean).join(' / '),
      priceCents: variant.value.priceCents,
      image: image.value,
      stock: variant.value.stock,
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
      <img v-if="image" :src="image" :alt="product.name" class="aspect-[3/4] w-full object-cover" />
      <div v-else class="grid aspect-[3/4] place-items-center text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ product.name.slice(0, 1) }}</div>
    </div>
    <div>
      <p v-if="product.category" class="text-xs uppercase tracking-[0.18em]" :style="{ color: 'var(--sf-muted)' }">{{ product.category }}</p>
      <h1 class="mt-2 text-4xl md:text-5xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ product.name }}</h1>
      <p class="mt-3 text-lg">{{ formatMoney(variant?.priceCents ?? product.priceCents, ctx.store.value.currency) }}</p>
      <p class="mt-4 whitespace-pre-line" :style="{ color: 'var(--sf-muted)' }">{{ product.description }}</p>
      <div v-for="name in product.optionNames" :key="name" class="mt-5">
        <p class="text-xs uppercase tracking-[0.16em]" :style="{ color: 'var(--sf-muted)' }">{{ name }}</p>
        <div class="mt-2 flex flex-wrap gap-2">
          <button
            v-for="value in valuesFor(name)"
            :key="value"
            type="button"
            class="border px-3 py-1.5 text-sm"
            :style="{ borderColor: 'var(--sf-border)', borderRadius: '999px', background: selected[name] === value ? 'var(--sf-primary)' : 'transparent', color: selected[name] === value ? 'var(--sf-on)' : 'var(--sf-text)' }"
            @click="choose(name, value)"
          >
            {{ value }}
          </button>
        </div>
      </div>
      <p v-if="!variant" class="mt-4 text-sm">That combination is not available.</p>
      <p v-else class="mt-4 text-sm">{{ variant.stock > 0 ? `${variant.stock} available` : 'Sold out' }}</p>
      <div class="mt-5 flex items-center gap-3">
        <input v-model.number="quantity" type="number" min="1" :max="variant?.stock || 1" class="w-20 border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" :disabled="!variant || variant.stock < 1" />
        <button type="button" class="px-5 py-2 text-sm font-semibold disabled:opacity-40" :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)', borderRadius: 'var(--sf-radius)' }" :disabled="!variant || variant.stock < 1" @click="add">Add to cart</button>
      </div>
      <p v-if="added" class="mt-3 text-sm">Added. <RouterLink :to="ctx.paths.value.cart">View cart</RouterLink></p>
    </div>
  </article>
</template>
