<template>
  <div>
    <PageRenderer v-if="snapshot && snapshot.blocks" :blocks="snapshot.blocks" />
    <div v-else class="py-24 text-center">Loading page...</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import PageRenderer from '@/components/PageRenderer.vue'
import { getCurrentShopSlug, getCurrentStoreId } from '@/services/shop'

const snapshot = ref<any | null>(null)

async function fetchSnapshot(path: string) {
  const shopSlug = getCurrentShopSlug()
  if (!shopSlug) return null
  const url = new URL(`/shops/${shopSlug}/page`, window.location.origin)
  url.searchParams.set('path', path)
  const storeId = await getCurrentStoreId()
  if (storeId) url.searchParams.set('storeId', storeId)

  const res = await fetch(url.toString(), { credentials: 'include' })
  if (!res.ok) { return null }
  const data = await res.json()
  return data.snapshot
}

onMounted(async () => {
  const p = window.location.pathname || '/'
  const snap = await fetchSnapshot(p)
  if (snap) {
    if (snap.tokens) {
      for (const [k, v] of Object.entries(snap.tokens)) {
        document.documentElement.style.setProperty(`--token-${k}`, String(v))
      }
    }
    snapshot.value = snap
  }
})
</script>
