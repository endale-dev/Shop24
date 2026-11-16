<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Modal from '@/components/ui/Modal.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { deleteCollection, listCollections } from '@/features/collections/api'
import type { Collection } from '@/types'

const items = ref<Collection[]>([])
const error = ref('')
const loading = ref(true)
const pendingDelete = ref<Collection | null>(null)

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = (await listCollections()).items
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load collections'
  } finally {
    loading.value = false
  }
}

async function confirmDelete() {
  if (!pendingDelete.value) return
  try {
    await deleteCollection(pendingDelete.value.id)
    pendingDelete.value = null
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not delete the collection'
    pendingDelete.value = null
  }
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Collections" subtitle="A named group of products with its own image. You choose what goes in.">
      <Button to="/collections/new">New collection</Button>
    </PageHeader>
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <p v-if="loading" class="text-sm text-stone-500">Loading collections…</p>
    <EmptyState v-else-if="items.length === 0" title="No collections yet" body="Group products such as coats or linen. Drafts stay off the public shop.">
      <Button to="/collections/new">Add a collection</Button>
    </EmptyState>
    <div v-else class="overflow-hidden rounded-2xl border border-line bg-card">
      <table class="w-full text-left text-sm">
        <thead class="border-b border-line text-xs uppercase tracking-wide text-stone-500">
          <tr><th class="px-4 py-3">Collection</th><th class="px-4 py-3">Products</th><th class="px-4 py-3">Status</th><th class="px-4 py-3"></th></tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id" class="border-b border-line last:border-0">
            <td class="px-4 py-3">
              <p class="font-medium">{{ item.name }}</p>
              <p class="text-xs text-stone-500">/c/{{ item.slug }}</p>
            </td>
            <td class="px-4 py-3">{{ item.productIds.length }}</td>
            <td class="px-4 py-3">{{ item.published ? 'Published' : 'Draft' }}</td>
            <td class="px-4 py-3 text-right">
              <RouterLink class="mr-3 font-medium text-rust" :to="`/collections/${item.id}`">Edit</RouterLink>
              <button class="text-stone-500" type="button" @click="pendingDelete = item">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <Modal title="Delete this collection?" :open="!!pendingDelete" @close="pendingDelete = null">
      <p>{{ pendingDelete?.name }} will leave the shop. The products stay.</p>
      <template #actions>
        <Button variant="ghost" @click="pendingDelete = null">Keep</Button>
        <Button variant="danger" @click="confirmDelete">Delete</Button>
      </template>
    </Modal>
  </div>
</template>
