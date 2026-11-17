<script setup lang="ts">
import { inject, ref } from 'vue'
import { formatMoney } from '@/features/money'
import { placeOrder } from '@/features/orders/api'
import { storefrontKey } from '@/features/storefront/context'
import type { Order } from '@/types'

const ctx = inject(storefrontKey)
const form = ref({ name: '', email: '', phone: '', address: '', city: '', notes: '' })
const error = ref('')
const loading = ref(false)
const order = ref<Order | null>(null)

async function submit() {
  if (!ctx) return
  loading.value = true
  error.value = ''
  try {
    order.value = await placeOrder(ctx.slug.value, {
      customer: form.value,
      items: ctx.cart.value.lines.value.map((line) => ({ productId: line.productId, variantId: line.variantId || '', quantity: line.quantity, unitPriceCents: line.priceCents })),
    })
    ctx.cart.value.clear()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not place the order'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div v-if="ctx?.store.value">
    <div v-if="order">
      <h1 class="text-4xl" :style="{ fontFamily: 'var(--sf-heading)' }">Order received</h1>
      <p class="mt-3 max-w-lg" :style="{ color: 'var(--sf-muted)' }">Pay {{ formatMoney(order.totalCents, ctx.store.value.currency) }} when it arrives. The shop will move this order from pending when they confirm it.</p>
      <RouterLink class="mt-6 inline-block text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]" :to="ctx.paths.value.home">Back to the shop</RouterLink>
    </div>
    <form v-else class="grid gap-8 md:grid-cols-[1.1fr_0.9fr]" @submit.prevent="submit">
      <div class="space-y-3">
        <h1 class="text-4xl" :style="{ fontFamily: 'var(--sf-heading)' }">Checkout</h1>
        <p class="text-sm" :style="{ color: 'var(--sf-muted)' }">Cash on delivery. No card is charged. If a price changed since you added it, checkout will ask you to review the cart.</p>
        <p v-if="error" class="text-sm">{{ error }}</p>
        <input v-model="form.name" required placeholder="Name" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <input v-model="form.email" required type="email" placeholder="Email" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <input v-model="form.phone" placeholder="Phone" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <input v-model="form.address" required minlength="4" placeholder="Address" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <input v-model="form.city" required placeholder="City" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <textarea v-model="form.notes" rows="3" placeholder="Notes for the seller" class="w-full border bg-transparent px-3 py-2" :style="{ borderColor: 'var(--sf-border)', borderRadius: 'var(--sf-radius)' }" />
        <button class="px-5 py-2 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)] disabled:opacity-40" :disabled="loading || ctx.cart.value.lines.value.length === 0" :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)', borderRadius: 'var(--sf-radius)' }">Place order</button>
      </div>
      <aside>
        <h2 class="text-2xl" :style="{ fontFamily: 'var(--sf-heading)' }">Your bag</h2>
        <p v-if="ctx.cart.value.lines.value.length === 0" class="mt-3 text-sm">Add a product before checking out.</p>
        <ul v-else class="mt-4 space-y-2 text-sm">
          <li v-for="line in ctx.cart.value.lines.value" :key="`${line.productId}:${line.variantId || ''}`" class="flex justify-between gap-3">
            <span>{{ line.quantity }} × {{ line.name }}<span v-if="line.optionLabel"> ({{ line.optionLabel }})</span></span>
            <span>{{ formatMoney(line.priceCents * line.quantity, ctx.store.value.currency) }}</span>
          </li>
        </ul>
        <p class="mt-4 text-lg">{{ formatMoney(ctx.cart.value.subtotal.value, ctx.store.value.currency) }}</p>
      </aside>
    </form>
  </div>
</template>
