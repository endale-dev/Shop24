<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Badge from '@/components/ui/Badge.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { formatDate, formatMoney } from '@/features/money'
import { listOrders } from '@/features/orders/api'
import { useShopStore } from '@/stores/shop'
import type { Order, OrderStatus } from '@/types'

const shop = useShopStore()
const status = ref('')
const items = ref<Order[]>([])
const error = ref('')
const loading = ref(true)
const statuses: Array<'' | OrderStatus> = ['', 'pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled']

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = (await listOrders(status.value)).items
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load orders'
  } finally {
    loading.value = false
  }
}

function tone(value: OrderStatus) {
  if (value === 'delivered') return 'good'
  if (value === 'cancelled') return 'bad'
  if (value === 'pending') return 'warn'
  return 'neutral'
}

onMounted(load)
</script>

<template>
  <div>
    <PageHeader title="Orders" subtitle="Cash on delivery. Changing a status does not charge a card." />
    <div class="mb-4 flex flex-wrap gap-2">
      <button v-for="option in statuses" :key="option || 'all'" type="button" class="rounded border px-3 py-1 text-sm capitalize focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent" :class="status === option ? 'border-text bg-text text-white' : 'border-border bg-surface text-text hover:bg-bg'" @click="status = option; load()">{{ option || 'All' }}</button>
    </div>
    <p v-if="error" class="mb-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <p v-if="loading" class="text-sm text-muted">Loading orders…</p>
    <EmptyState v-else-if="items.length === 0" icon="clipboard" title="No orders" body="Orders appear after someone checks out on the public shop." />
    <div v-else class="overflow-x-auto rounded border border-border bg-surface">
      <table class="w-full text-left text-sm">
        <thead class="border-b border-border text-xs font-medium text-muted">
          <tr><th class="px-4 py-3 font-medium">Customer</th><th class="px-4 py-3 font-medium">Date</th><th class="px-4 py-3 font-medium">Total</th><th class="px-4 py-3 font-medium">Status</th></tr>
        </thead>
        <tbody>
          <tr v-for="order in items" :key="order.id" class="border-b border-border last:border-0 hover:bg-bg">
            <td class="px-4 py-3"><RouterLink class="font-medium text-text underline" :to="`/orders/${order.id}`">{{ order.customer.name }}</RouterLink></td>
            <td class="px-4 py-3">{{ formatDate(order.createdAt) }}</td>
            <td class="px-4 py-3">{{ formatMoney(order.totalCents, shop.shop?.currency) }}</td>
            <td class="px-4 py-3"><Badge :tone="tone(order.status)">{{ order.status }}</Badge></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
