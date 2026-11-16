<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { centsToDollars, dollarsToCents } from '@/features/money'
import { createProduct, getProduct, setStock, updateProduct } from '@/features/products/api'

const route = useRoute()
const router = useRouter()
const id = typeof route.params.id === 'string' ? route.params.id : ''
const name = ref('')
const description = ref('')
const price = ref('0.00')
const images = ref('')
const category = ref('')
const sku = ref('')
const stock = ref(0)
const published = ref(false)
const originalStock = ref<number | null>(null)
const error = ref('')
const loading = ref(false)

onMounted(async () => {
  if (!id) return
  loading.value = true
  try {
    const product = await getProduct(id)
    name.value = product.name
    description.value = product.description
    price.value = centsToDollars(product.priceCents)
    images.value = product.images.join('\n')
    category.value = product.category
    sku.value = product.sku
    stock.value = product.stock
    originalStock.value = product.stock
    published.value = product.published
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load the product'
  } finally {
    loading.value = false
  }
})

async function submit() {
  const priceCents = dollarsToCents(price.value)
  if (priceCents == null || priceCents < 1) {
    error.value = 'Enter a price like 24.00'
    return
  }
  error.value = ''
  loading.value = true
  const body = {
    name: name.value,
    description: description.value,
    priceCents,
    images: images.value.split('\n').map((line) => line.trim()).filter(Boolean),
    category: category.value,
    sku: sku.value,
    stock: Number(stock.value),
    published: published.value,
  }
  try {
    if (id) {
      await updateProduct(id, body)
      if (originalStock.value !== null && Number(stock.value) !== originalStock.value) {
        await setStock(id, Number(stock.value))
      }
    } else {
      await createProduct(body)
    }
    await router.push('/products')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not save the product'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="max-w-2xl" @submit.prevent="submit">
    <PageHeader :title="id ? 'Edit product' : 'New product'" subtitle="Stock is the number you can still sell. Orders reduce it automatically." />
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <div class="space-y-4 rounded-2xl border border-line bg-card p-5">
      <Field label="Name"><input v-model="name" :class="controlClass" required /></Field>
      <Field label="Description"><textarea v-model="description" :class="controlClass" rows="4" /></Field>
      <div class="grid gap-4 sm:grid-cols-2">
        <Field label="Price"><input v-model="price" :class="controlClass" inputmode="decimal" required /></Field>
        <Field label="Stock"><input v-model.number="stock" :class="controlClass" type="number" min="0" required /></Field>
        <Field label="SKU"><input v-model="sku" :class="controlClass" required /></Field>
        <Field label="Category"><input v-model="category" :class="controlClass" /></Field>
      </div>
      <Field label="Image URLs" hint="One http or https URL per line."><textarea v-model="images" :class="controlClass" rows="3" /></Field>
      <label class="flex items-center gap-2 text-sm"><input v-model="published" type="checkbox" /> Published on the storefront</label>
      <div class="flex gap-2">
        <Button type="submit" :loading="loading">Save product</Button>
        <RouterLink to="/products"><Button variant="ghost" type="button">Cancel</Button></RouterLink>
      </div>
    </div>
  </form>
</template>
