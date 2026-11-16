<script setup lang="ts">
import { inject } from 'vue'
import { formatMoney } from '@/features/money'
import { storefrontKey } from '@/features/storefront/context'

const ctx = inject(storefrontKey)
</script>

<template>
  <div v-if="ctx?.store.value">
    <h1 class="text-4xl" :style="{ fontFamily: 'var(--sf-heading)' }">Cart</h1>
    <p v-if="ctx.cart.value.lines.value.length === 0" class="mt-4 text-sm" :style="{ color: 'var(--sf-muted)' }">Your cart is empty.</p>
    <ul v-else class="mt-6 divide-y" :style="{ borderColor: 'var(--sf-border)' }">
      <li v-for="line in ctx.cart.value.lines.value" :key="line.productId" class="flex items-center justify-between gap-4 py-4">
        <div>
          <p class="font-medium">{{ line.name }}</p>
          <p class="text-sm" :style="{ color: 'var(--sf-muted)' }">{{ formatMoney(line.priceCents, ctx.store.value.currency) }}</p>
        </div>
        <div class="flex items-center gap-3">
          <input :value="line.quantity" type="number" min="1" :max="line.stock" class="w-16 border bg-transparent px-2 py-1" :style="{ borderColor: 'var(--sf-border)' }" @change="ctx.cart.value.setQuantity(line.productId, Number(($event.target as HTMLInputElement).value))" />
          <button type="button" class="text-sm" @click="ctx.cart.value.remove(line.productId)">Remove</button>
        </div>
      </li>
    </ul>
    <div v-if="ctx.cart.value.lines.value.length" class="mt-6 flex items-center justify-between">
      <p class="text-lg">{{ formatMoney(ctx.cart.value.subtotal.value, ctx.store.value.currency) }}</p>
      <RouterLink class="px-5 py-2 text-sm font-semibold" :to="ctx.paths.value.checkout" :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)', borderRadius: 'var(--sf-radius)' }">Checkout</RouterLink>
    </div>
  </div>
</template>
