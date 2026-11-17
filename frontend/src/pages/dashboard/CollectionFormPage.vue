<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass, textareaClass } from '@/components/ui/control'
import { createCollection, getCollection, updateCollection } from '@/features/collections/api'
import { listProducts } from '@/features/products/api'
import type { Product } from '@/types'

const route = useRoute()
const router = useRouter()
const id = typeof route.params.id === 'string' ? route.params.id : ''
const name = ref('')
const slug = ref('')
const slugTouched = ref(false)
const description = ref('')
const imageUrl = ref('')
const published = ref(false)
const productIds = ref<string[]>([])
const products = ref<Product[]>([])
const error = ref('')
const loading = ref(false)

function slugify(value: string) {
  return value.toLowerCase().trim().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '').slice(0, 40)
}

function onName(value: string) {
  name.value = value
  if (!slugTouched.value) slug.value = slugify(value)
}

const ordered = computed(() => {
  const byId = new Map(products.value.map((product) => [product.id, product]))
  const selected = productIds.value.map((productId) => byId.get(productId)).filter((product): product is Product => !!product)
  const rest = products.value.filter((product) => !productIds.value.includes(product.id))
  return [...selected, ...rest]
})

function toggle(productId: string) {
  const index = productIds.value.indexOf(productId)
  if (index >= 0) productIds.value.splice(index, 1)
  else if (productIds.value.length < 100) productIds.value.push(productId)
}

function move(productId: string, direction: number) {
  const index = productIds.value.indexOf(productId)
  const next = index + direction
  if (index < 0 || next < 0 || next >= productIds.value.length) return
  const copy = [...productIds.value]
  const [item] = copy.splice(index, 1)
  copy.splice(next, 0, item)
  productIds.value = copy
}

onMounted(async () => {
  loading.value = true
  try {
    products.value = (await listProducts()).items
    if (!id) return
    const collection = await getCollection(id)
    name.value = collection.name
    slug.value = collection.slug
    slugTouched.value = true
    description.value = collection.description
    imageUrl.value = collection.imageUrl
    published.value = collection.published
    productIds.value = [...collection.productIds]
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load the collection'
  } finally {
    loading.value = false
  }
})

async function submit() {
  error.value = ''
  loading.value = true
  const body = {
    name: name.value,
    slug: slug.value,
    description: description.value,
    imageUrl: imageUrl.value.trim(),
    productIds: productIds.value,
    published: published.value,
  }
  try {
    if (id) await updateCollection(id, body)
    else await createCollection(body)
    await router.push('/collections')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not save the collection'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="max-w-2xl" @submit.prevent="submit">
    <PageHeader :title="id ? 'Edit collection' : 'New collection'" subtitle="One image, a short address, and the products you put in it." />
    <p v-if="error" class="mb-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <div class="space-y-4 rounded border border-border bg-surface p-4">
      <Field label="Name"><input :value="name" :class="controlClass" required @input="onName(($event.target as HTMLInputElement).value)" /></Field>
      <Field label="Address" hint="Shown as /c/this-address"><input v-model="slug" :class="controlClass" required @input="slugTouched = true" /></Field>
      <Field label="Description"><textarea v-model="description" :class="textareaClass" rows="3" /></Field>
      <Field label="Image URL"><input v-model="imageUrl" :class="controlClass" placeholder="https://" /></Field>
      <div>
        <p class="mb-2 text-sm font-semibold text-text">Products</p>
        <p v-if="products.length === 0" class="text-sm text-muted">Add a product before filling a collection.</p>
        <ul v-else class="divide-y divide-border rounded border border-border">
          <li v-for="product in ordered" :key="product.id" class="flex items-center gap-3 px-3 py-2 text-sm">
            <input :checked="productIds.includes(product.id)" type="checkbox" @change="toggle(product.id)" />
            <span class="min-w-0 flex-1 truncate">{{ product.name }}</span>
            <button v-if="productIds.includes(product.id)" class="text-xs text-muted hover:text-text" type="button" @click="move(product.id, -1)">Up</button>
            <button v-if="productIds.includes(product.id)" class="text-xs text-muted hover:text-text" type="button" @click="move(product.id, 1)">Down</button>
          </li>
        </ul>
      </div>
      <label class="flex items-center gap-2 text-sm"><input v-model="published" type="checkbox" /> Published on the storefront</label>
      <div class="flex gap-2">
        <Button type="submit" :loading="loading">Save collection</Button>
        <Button to="/collections" variant="ghost">Cancel</Button>
      </div>
    </div>
  </form>
</template>
