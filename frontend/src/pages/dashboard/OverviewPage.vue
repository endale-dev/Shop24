<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Field from '@/components/ui/Field.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import { controlClass } from '@/components/ui/control'
import { overview } from '@/features/analytics/api'
import { formatDate, formatMoney } from '@/features/money'
import { createStore } from '@/features/stores/api'
import { useShopStore } from '@/stores/shop'
import type { Overview } from '@/types'

const shop = useShopStore()
const days = ref(30)
const data = ref<Overview | null>(null)
const error = ref('')
const loading = ref(false)
const creating = ref(false)
const form = ref({ name: '', slug: '', description: '' })

const maxRevenue = computed(() => Math.max(...(data.value?.revenueOverTime.map((point) => point.revenueCents) ?? [0]), 1))

onMounted(async () => {
  await shop.load()
  if (shop.shop) await load()
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    data.value = await overview(days.value)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not load analytics'
  } finally {
    loading.value = false
  }
}

function suggestSlug() {
  if (form.value.slug.trim()) return
  form.value.slug = form.value.name.toLowerCase().trim().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '').slice(0, 40)
}

async function submitStore() {
  creating.value = true
  error.value = ''
  try {
    shop.set(await createStore(form.value))
    await load()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not create the store'
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <div>
    <PageHeader title="Overview" subtitle="Sales, stock, and the orders that need you." />
    <p v-if="error" class="mb-4 rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
    <form v-if="!shop.shop" class="max-w-xl space-y-4 rounded-2xl border border-line bg-card p-5" @submit.prevent="submitStore">
      <h2 class="font-serif text-2xl">Open your store</h2>
      <p class="text-sm text-stone-600">This creates the shop customers will visit. You can design the page after the catalog exists.</p>
      <Field label="Store name"><input v-model="form.name" :class="controlClass" required minlength="2" @blur="suggestSlug" /></Field>
      <Field label="Store address" hint="Becomes the shop subdomain, such as this-name.localhost while you test."><input v-model="form.slug" :class="controlClass" required /></Field>
      <Field label="Description"><textarea v-model="form.description" :class="controlClass" rows="3" /></Field>
      <Button type="submit" :loading="creating">Create store</Button>
    </form>
    <div v-else>
      <div class="mb-4 flex gap-2">
        <button v-for="option in [7, 30, 90]" :key="option" type="button" class="rounded-full px-3 py-1 text-sm" :class="days === option ? 'bg-ink text-white' : 'bg-white text-stone-600'" @click="days = option; load()">{{ option }} days</button>
      </div>
      <div v-if="loading || !data" class="text-sm text-stone-500">Loading the numbers…</div>
      <template v-else>
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <article v-for="card in [
            ['Revenue', formatMoney(data.totalRevenueCents, shop.shop.currency)],
            ['Orders', String(data.totalOrders)],
            ['Average order', formatMoney(data.averageOrderValueCents, shop.shop.currency)],
            ['In stock', `${data.unitsInStock} ${data.unitsInStock === 1 ? 'unit' : 'units'} · ${data.totalProducts} ${data.totalProducts === 1 ? 'product' : 'products'}`],
          ]" :key="card[0]" class="rounded-2xl border border-line bg-card p-4">
            <p class="text-xs uppercase tracking-wide text-stone-500">{{ card[0] }}</p>
            <p class="mt-2 font-serif text-2xl">{{ card[1] }}</p>
          </article>
        </div>
        <div class="mt-6 grid gap-4 lg:grid-cols-[1.4fr_1fr]">
          <section class="rounded-2xl border border-line bg-card p-4">
            <h2 class="font-serif text-2xl">Revenue</h2>
            <div class="mt-4 flex h-36 items-end gap-1">
              <div v-for="point in data.revenueOverTime" :key="point.date" class="group relative flex-1 rounded-t bg-rust/80" :style="{ height: `${Math.max(point.revenueCents ? 8 : 2, (point.revenueCents / maxRevenue) * 100)}%` }">
                <span class="pointer-events-none absolute bottom-full left-1/2 mb-1 hidden -translate-x-1/2 whitespace-nowrap rounded bg-ink px-2 py-1 text-[10px] text-white group-hover:block">{{ point.date }} · {{ formatMoney(point.revenueCents, shop.shop.currency) }}</span>
              </div>
            </div>
          </section>
          <section class="rounded-2xl border border-line bg-card p-4">
            <h2 class="font-serif text-2xl">Best sellers</h2>
            <EmptyState v-if="data.bestSellers.length === 0" title="No sales in this range" />
            <ul v-else class="mt-3 space-y-3">
              <li v-for="item in data.bestSellers" :key="item.productId" class="flex items-center justify-between gap-3 text-sm">
                <span>{{ item.name }}</span>
                <span class="text-stone-500">{{ item.quantity }} · {{ formatMoney(item.revenueCents, shop.shop.currency) }}</span>
              </li>
            </ul>
          </section>
        </div>
        <div class="mt-4 grid gap-4 lg:grid-cols-2">
          <section class="rounded-2xl border border-line bg-card p-4">
            <h2 class="font-serif text-2xl">Recent orders</h2>
            <EmptyState v-if="data.recentOrders.length === 0" title="No orders yet" body="They will show up here after a customer checks out." />
            <ul v-else class="mt-3 divide-y divide-line">
              <li v-for="order in data.recentOrders" :key="order.id">
                <RouterLink class="flex items-center justify-between py-3 text-sm" :to="`/orders/${order.id}`">
                  <span>{{ order.customer.name }} · {{ formatDate(order.createdAt) }}</span>
                  <span>{{ formatMoney(order.totalCents, shop.shop.currency) }}</span>
                </RouterLink>
              </li>
            </ul>
          </section>
          <section class="rounded-2xl border border-line bg-card p-4">
            <h2 class="font-serif text-2xl">Low stock</h2>
            <p class="text-xs text-stone-500">At or below {{ data.lowStockThreshold }} units.</p>
            <EmptyState v-if="data.lowStock.length === 0" title="Stock looks fine" />
            <ul v-else class="mt-3 space-y-2">
              <li v-for="product in data.lowStock" :key="product.id" class="flex items-center justify-between text-sm">
                <RouterLink :to="`/products/${product.id}`">{{ product.name }}</RouterLink>
                <Badge :tone="product.stock === 0 ? 'bad' : 'warn'">{{ product.stock }} left</Badge>
              </li>
            </ul>
          </section>
        </div>
      </template>
    </div>
  </div>
</template>
