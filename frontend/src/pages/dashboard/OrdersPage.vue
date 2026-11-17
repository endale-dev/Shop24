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
      <button v-for="option in statuses" :key="option || 'all'" type="button" class="capitalize focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent" :class="status === option ? 'h-8 rounded border border-accent bg-accent px-2.5 text-sm text-white' : 'h-8 rounded border border-border bg-surface px-2.5 text-sm text-text hover:bg-bg'" @click="status = option; load()">{{ option || 'All' }}</button>
    </div>
    <p v-if="error" class="mb-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
    <p v-if="loading" class="text-sm text-muted">Loading orders…</p>
    <EmptyState v-else-if="items.length === 0" icon="clipboard" title="No orders" body="Orders appear after someone checks out on the public shop." />
    <div v-else class="overflow-x-auto rounded border border-border bg-surface">
      <table class="w-full border-separate border-spacing-0 text-left text-sm">
        <thead>
          <tr>
            <th class="border-b border-border bg-bg px-3 py-2 text-left text-xs font-medium text-muted">Customer</th>
            <th class="border-b border-border bg-bg px-3 py-2 text-left text-xs font-medium text-muted">Date</th>
            <th class="border-b border-border bg-bg px-3 py-2 text-right text-xs font-medium text-muted">Total</th>
            <th class="border-b border-border bg-bg px-3 py-2 text-left text-xs font-medium text-muted">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in items" :key="order.id" class="hover:bg-bg">
            <td class="border-b border-border px-3 py-2 align-middle"><RouterLink class="font-medium hover:underline" :to="`/orders/${order.id}`">{{ order.customer.name }}</RouterLink></td>
            <td class="border-b border-border px-3 py-2 align-middle">{{ formatDate(order.createdAt) }}</td>
            <td class="border-b border-border px-3 py-2 text-right align-middle tabular-nums">{{ formatMoney(order.totalCents, shop.shop?.currency) }}</td>
            <td class="border-b border-border px-3 py-2 align-middle"><Badge :tone="tone(order.status)">{{ order.status }}</Badge></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
