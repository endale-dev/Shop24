<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { formatMoney } from '@/features/money'
import { getOrder, updateOrderStatus } from '@/features/orders/api'
import { useShopStore } from '@/stores/shop'
import type { Order, OrderStatus } from '@/types'

const route = useRoute()
const shop = useShopStore()
const order = ref<Order | null>(null)
const status = ref<OrderStatus>('pending')
const error = ref('')
const loading = ref(false)
const statuses: OrderStatus[] = ['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled']

onMounted(async () => {
  try {
    order.value = await getOrder(String(route.params.id))
    status.value = order.value.status
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load the order'
  }
})

async function save() {
  if (!order.value) return
  loading.value = true
  error.value = ''
  try {
    order.value = await updateOrderStatus(order.value.id, status.value)
    status.value = order.value.status
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not update the order'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div v-if="order">
    <PageHeader :title="order.customer.name" :subtitle="`Cash on delivery · ${formatMoney(order.totalCents, shop.shop?.currency)}`" />
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <div class="grid gap-4 lg:grid-cols-[1.2fr_0.8fr]">
      <section class="rounded-2xl border border-line bg-card p-4">
        <table class="w-full text-left text-sm">
          <thead class="text-xs uppercase tracking-wide text-stone-500"><tr><th class="py-2">Item</th><th>Qty</th><th>Price</th><th>Total</th></tr></thead>
          <tbody>
            <tr v-for="item in order.items" :key="item.productId + item.sku" class="border-t border-line">
              <td class="py-3">{{ item.name }}<span class="block text-xs text-stone-500">{{ item.sku }}</span></td>
              <td>{{ item.quantity }}</td>
              <td>{{ formatMoney(item.unitPriceCents, shop.shop?.currency) }}</td>
              <td>{{ formatMoney(item.lineTotalCents, shop.shop?.currency) }}</td>
            </tr>
          </tbody>
        </table>
      </section>
      <section class="space-y-3 rounded-2xl border border-line bg-card p-4 text-sm">
        <p>{{ order.customer.email }}</p>
        <p v-if="order.customer.phone">{{ order.customer.phone }}</p>
        <p>{{ order.customer.address }}<br />{{ order.customer.city }}</p>
        <p v-if="order.customer.notes" class="text-stone-600">{{ order.customer.notes }}</p>
        <label class="block text-xs uppercase tracking-wide text-stone-500">Status
          <select v-model="status" :class="controlClass" class="mt-1" :disabled="order.status === 'cancelled'">
            <option v-for="option in statuses" :key="option" :value="option">{{ option }}</option>
          </select>
        </label>
        <p v-if="order.status === 'cancelled'" class="text-xs text-stone-500">This order is closed. Stock was returned for products that still exist.</p>
        <Button :disabled="order.status === 'cancelled'" :loading="loading" @click="save">Update status</Button>
      </section>
    </div>
  </div>
  <p v-else class="text-sm text-stone-500">{{ error || 'Loading order…' }}</p>
</template>
