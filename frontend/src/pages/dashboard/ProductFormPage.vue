<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass, textareaClass } from '@/components/ui/control'
import { centsToDollars, dollarsToCents } from '@/features/money'
import { createProduct, getProduct, setStock, updateProduct } from '@/features/products/api'

type Draft = {
  id: string
  sku: string
  price: string
  stock: number
  originalStock: number | null
  imageUrl: string
  values: string[]
}

const route = useRoute()
const router = useRouter()
const id = typeof route.params.id === 'string' ? route.params.id : ''
const name = ref('')
const description = ref('')
const images = ref('')
const category = ref('')
const published = ref(false)
const optionNames = ref<string[]>([])
const variants = ref<Draft[]>([blank()])
const error = ref('')
const loading = ref(false)
const versions = computed(() => optionNames.value.length > 0)

function blank(values: string[] = []): Draft {
  return { id: '', sku: '', price: '0.00', stock: 0, originalStock: null, imageUrl: '', values }
}

function enableVersions() {
  optionNames.value = ['Color']
  variants.value.forEach((variant) => {
    variant.values = [variant.values[0] ?? '']
  })
  if (variants.value.length < 2) variants.value.push(blank(['']))
}

function useSingleVersion() {
  const first = variants.value[0] ?? blank()
  first.values = []
  variants.value = [first]
  optionNames.value = []
}

function addOption() {
  if (optionNames.value.length >= 3) return
  optionNames.value.push('')
  variants.value.forEach((variant) => variant.values.push(''))
}

function removeOption(index: number) {
  optionNames.value.splice(index, 1)
  variants.value.forEach((variant) => variant.values.splice(index, 1))
  if (optionNames.value.length === 0) useSingleVersion()
}

function addVariant() {
  if (variants.value.length >= 30) return
  variants.value.push(blank(optionNames.value.map(() => '')))
}

function removeVariant(index: number) {
  if (variants.value.length === 1) return
  variants.value.splice(index, 1)
}

onMounted(async () => {
  if (!id) return
  loading.value = true
  try {
    const product = await getProduct(id)
    name.value = product.name
    description.value = product.description
    images.value = product.images.join('\n')
    category.value = product.category
    published.value = product.published
    const names = product.optionNames ?? []
    optionNames.value = [...names]
    const rows = product.variants?.length
      ? product.variants
      : [{ id: product.id, sku: product.sku, priceCents: product.priceCents, stock: product.stock, imageUrl: '', options: [] }]
    variants.value = rows.map((variant) => ({
      id: variant.id,
      sku: variant.sku,
      price: centsToDollars(variant.priceCents),
      stock: variant.stock,
      originalStock: variant.stock,
      imageUrl: variant.imageUrl || '',
      values: names.map((option) => variant.options?.find((item) => item.name === option)?.value ?? ''),
    }))
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load the product'
  } finally {
    loading.value = false
  }
})

async function submit() {
  const names = optionNames.value.map((item) => item.trim())
  if (names.some((item) => item === '')) {
    error.value = 'Name each option, or remove the empty one'
    return
  }
  const built = []
  for (const variant of variants.value) {
    const priceCents = dollarsToCents(variant.price)
    if (priceCents == null || priceCents < 1) {
      error.value = 'Enter a price like 24.00'
      return
    }
    built.push({
      id: variant.id,
      sku: variant.sku.trim(),
      priceCents,
      stock: Number(variant.stock),
      imageUrl: variant.imageUrl.trim(),
      options: names.map((option, index) => ({ name: option, value: variant.values[index]?.trim() ?? '' })),
    })
  }
  error.value = ''
  loading.value = true
  const body = {
    name: name.value,
    description: description.value,
    images: images.value.split('\n').map((line) => line.trim()).filter(Boolean),
    category: category.value,
    optionNames: names,
    variants: built,
    published: published.value,
  }
  try {
    if (id) {
      await updateProduct(id, body)
      for (const variant of variants.value) {
        if (variant.id && variant.originalStock !== null && Number(variant.stock) !== variant.originalStock) {
          await setStock(id, variant.id, Number(variant.stock))
        }
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
  <form class="max-w-3xl" @submit.prevent="submit">
    <PageHeader :title="id ? 'Edit product' : 'New product'" subtitle="Each version keeps its own price, SKU, stock, and image. Nothing else is generated for you." />
    <p v-if="error" class="mb-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <div class="space-y-4 rounded border border-border bg-surface p-4">
      <Field label="Name"><input v-model="name" :class="controlClass" required /></Field>
      <Field label="Description"><textarea v-model="description" :class="textareaClass" rows="4" /></Field>
      <Field label="Category"><input v-model="category" :class="controlClass" /></Field>
      <Field label="Image URLs" hint="One http or https URL per line. A version can override this with its own image."><textarea v-model="images" :class="textareaClass" rows="3" /></Field>

      <div v-if="!versions" class="grid gap-4 sm:grid-cols-3">
        <Field label="Price"><input v-model="variants[0].price" :class="controlClass" inputmode="decimal" required /></Field>
        <Field label="Stock"><input v-model.number="variants[0].stock" :class="controlClass" type="number" min="0" required /></Field>
        <Field label="SKU"><input v-model="variants[0].sku" :class="controlClass" required /></Field>
      </div>
      <button v-if="!versions" class="text-sm font-medium text-text underline" type="button" @click="enableVersions">Sell more than one version</button>

      <div v-else class="space-y-3">
        <div class="flex items-center justify-between gap-3">
          <p class="text-sm font-semibold text-text">Options</p>
          <button class="text-sm text-muted hover:text-danger" type="button" @click="useSingleVersion">Use one version</button>
        </div>
        <div v-for="(_, index) in optionNames" :key="index" class="flex gap-2">
          <input v-model="optionNames[index]" :class="controlClass" placeholder="Color, size, capacity…" required />
          <button class="shrink-0 text-sm text-muted hover:text-danger" type="button" @click="removeOption(index)">Remove</button>
        </div>
        <button v-if="optionNames.length < 3" class="text-sm font-medium text-text underline" type="button" @click="addOption">Add option</button>
        <div class="space-y-3">
          <div v-for="(variant, index) in variants" :key="index" class="rounded border border-border bg-bg p-3">
            <div class="mb-3 flex items-center justify-between">
              <p class="text-sm font-semibold text-text">Version {{ index + 1 }}</p>
              <button v-if="variants.length > 1" class="text-sm text-muted hover:text-danger" type="button" @click="removeVariant(index)">Remove</button>
            </div>
            <div class="grid gap-3 sm:grid-cols-2">
              <Field v-for="(option, optionIndex) in optionNames" :key="optionIndex" :label="option || 'Value'">
                <input v-model="variant.values[optionIndex]" :class="controlClass" required />
              </Field>
              <Field label="Price"><input v-model="variant.price" :class="controlClass" inputmode="decimal" required /></Field>
              <Field label="Stock"><input v-model.number="variant.stock" :class="controlClass" type="number" min="0" required /></Field>
              <Field label="SKU"><input v-model="variant.sku" :class="controlClass" required /></Field>
              <Field label="Image URL"><input v-model="variant.imageUrl" :class="controlClass" placeholder="https://" /></Field>
            </div>
          </div>
        </div>
        <button v-if="variants.length < 30" class="text-sm font-medium text-text underline" type="button" @click="addVariant">Add version</button>
      </div>

      <label class="flex items-center gap-2 text-sm"><input v-model="published" type="checkbox" /> Published on the storefront</label>
      <div class="flex gap-2">
        <Button type="submit" :loading="loading">Save product</Button>
        <Button to="/products" variant="ghost">Cancel</Button>
      </div>
    </div>
  </form>
</template>
