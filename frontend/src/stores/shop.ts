import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ApiError } from '@/services/api'
import { myStore } from '@/features/stores/api'
import type { Store } from '@/types'

export const useShopStore = defineStore('shop', () => {
  const shop = ref<Store | null>(null)
  const loaded = ref(false)

  async function load(force = false) {
    if (loaded.value && !force) return shop.value
    try {
      shop.value = await myStore()
    } catch (error) {
      if (error instanceof ApiError && error.status === 404) shop.value = null
      else throw error
    } finally {
      loaded.value = true
    }
    return shop.value
  }

  function set(next: Store | null) {
    shop.value = next
    loaded.value = true
  }

  function reset() {
    shop.value = null
    loaded.value = false
  }

  return { shop, loaded, load, set, reset }
})
