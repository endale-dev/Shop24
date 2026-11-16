<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { publicUrl } from '@/features/stores/host'
import { publishStore, unpublishStore, updateStore } from '@/features/stores/api'
import { useShopStore } from '@/stores/shop'

const shop = useShopStore()
const name = ref('')
const slug = ref('')
const description = ref('')
const logoUrl = ref('')
const currency = ref('USD')
const error = ref('')
const message = ref('')
const loading = ref(false)

onMounted(async () => {
  const current = await shop.load()
  if (!current) return
  name.value = current.name
  slug.value = current.slug
  description.value = current.description
  logoUrl.value = current.logoUrl
  currency.value = current.currency
})

async function save() {
  loading.value = true
  error.value = ''
  message.value = ''
  try {
    shop.set(await updateStore({ name: name.value, slug: slug.value, description: description.value, logoUrl: logoUrl.value, currency: currency.value }))
    message.value = 'Store details saved.'
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not save the store'
  } finally {
    loading.value = false
  }
}

async function toggle() {
  if (!shop.shop) return
  loading.value = true
  error.value = ''
  try {
    shop.set(shop.shop.published ? await unpublishStore() : await publishStore())
    message.value = shop.shop.published ? 'The shop is public.' : 'The shop is hidden.'
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not change publishing'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <PageHeader title="Store settings" subtitle="Customers use the shop subdomain. This dashboard stays on the main domain." />
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <p v-if="message" class="mb-4 rounded-xl bg-lime-50 px-3 py-2 text-sm text-lime-900">{{ message }}</p>
    <form v-if="shop.shop" class="space-y-4 rounded-2xl border border-line bg-card p-5" @submit.prevent="save">
      <Field label="Store name"><input v-model="name" :class="controlClass" required /></Field>
      <Field label="Address"><input v-model="slug" :class="controlClass" required /></Field>
      <Field label="Description"><textarea v-model="description" :class="controlClass" rows="4" /></Field>
      <Field label="Logo URL"><input v-model="logoUrl" :class="controlClass" placeholder="https://" /></Field>
      <Field label="Currency">
        <select v-model="currency" :class="controlClass">
          <option>USD</option>
          <option>EUR</option>
          <option>GBP</option>
        </select>
      </Field>
      <div class="rounded-xl bg-paper p-3 text-sm text-stone-600">
        <p>Shop: <a class="font-medium text-ink underline" :href="publicUrl(shop.shop.slug).subdomain" target="_blank" rel="noreferrer">{{ publicUrl(shop.shop.slug).subdomain }}</a></p>
        <p class="mt-1">Same-host preview: <a class="font-medium text-ink underline" :href="publicUrl(shop.shop.slug).path" target="_blank" rel="noreferrer">{{ publicUrl(shop.shop.slug).path }}</a></p>
      </div>
      <div class="flex flex-wrap gap-2">
        <Button type="submit" :loading="loading">Save settings</Button>
        <Button type="button" variant="secondary" @click="toggle">{{ shop.shop.published ? 'Unpublish shop' : 'Publish shop' }}</Button>
      </div>
    </form>
  </div>
</template>
