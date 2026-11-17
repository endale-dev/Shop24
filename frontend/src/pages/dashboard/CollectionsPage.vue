<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Icon from '@/components/ui/Icon.vue'
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
      <Button to="/collections/new"><Icon name="plus" /> New collection</Button>
    </PageHeader>
    <p v-if="error" class="mb-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <p v-if="loading" class="text-sm text-muted">Loading collections…</p>
    <EmptyState v-else-if="items.length === 0" icon="layers" title="No collections yet" body="Group products such as coats or linen. Drafts stay off the public shop.">
      <Button to="/collections/new"><Icon name="plus" /> Add a collection</Button>
    </EmptyState>
    <div v-else class="overflow-x-auto rounded border border-border bg-surface">
      <table class="w-full text-left text-sm">
        <thead class="border-b border-border text-xs font-medium text-muted">
          <tr><th class="px-4 py-3 font-medium">Collection</th><th class="px-4 py-3 font-medium">Products</th><th class="px-4 py-3 font-medium">Status</th><th class="px-4 py-3 font-medium"></th></tr>
        </thead>
        <tbody>
          <tr v-for="item in items" :key="item.id" class="border-b border-border last:border-0 hover:bg-bg">
            <td class="px-4 py-3">
              <p class="font-medium">{{ item.name }}</p>
              <p class="text-xs text-muted">/c/{{ item.slug }}</p>
            </td>
            <td class="px-4 py-3">{{ item.productIds.length }}</td>
            <td class="px-4 py-3">{{ item.published ? 'Published' : 'Draft' }}</td>
            <td class="px-4 py-3 text-right">
              <RouterLink class="mr-3 text-text underline" :to="`/collections/${item.id}`">Edit</RouterLink>
              <button class="text-muted hover:text-danger" type="button" @click="pendingDelete = item">Delete</button>
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
