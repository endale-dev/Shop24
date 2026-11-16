<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useShopStore } from '@/stores/shop'

const auth = useAuthStore()
const shop = useShopStore()
const router = useRouter()
const route = useRoute()
const open = ref(false)
const links = [
  { to: '/', label: 'Overview', store: false },
  { to: '/products', label: 'Products', store: true },
  { to: '/collections', label: 'Collections', store: true },
  { to: '/orders', label: 'Orders', store: true },
  { to: '/builder', label: 'Builder', store: true },
  { to: '/settings', label: 'Settings', store: true },
]
const visibleLinks = computed(() => links.filter((link) => !link.store || shop.shop))

async function signOut() {
  await auth.logout()
  shop.reset()
  await router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-paper text-ink md:grid md:grid-cols-[240px_1fr]">
    <aside class="bg-ink text-stone-100 md:min-h-screen">
      <div class="flex items-center justify-between px-5 py-5">
        <RouterLink to="/" class="font-serif text-3xl">Sable</RouterLink>
        <button class="rounded-full border border-white/20 px-3 py-1 text-xs md:hidden" type="button" @click="open = !open">Menu</button>
      </div>
      <nav class="space-y-1 px-3 pb-5" :class="open ? 'block' : 'hidden md:block'">
        <RouterLink
          v-for="link in visibleLinks"
          :key="link.to"
          :to="link.to"
          class="block rounded-xl px-3 py-2 text-sm text-stone-300 transition hover:bg-white/10 hover:text-white"
          :class="route.path === link.to || (link.to !== '/' && route.path.startsWith(link.to)) ? 'bg-white/10 text-white' : ''"
          @click="open = false"
        >
          {{ link.label }}
        </RouterLink>
      </nav>
    </aside>
    <div class="min-w-0">
      <header class="flex items-center justify-between border-b border-line px-5 py-4 md:px-8">
        <div>
          <p class="text-sm font-medium">{{ shop.shop?.name || 'No store yet' }}</p>
          <p v-if="shop.shop" class="text-xs text-stone-500">{{ shop.shop.published ? 'Published' : 'Unpublished' }} · /{{ shop.shop.slug }}</p>
        </div>
        <button class="text-sm text-stone-600 hover:text-ink" type="button" @click="signOut">Sign out</button>
      </header>
      <main class="px-5 py-6 md:px-8 md:py-8"><RouterView /></main>
    </div>
  </div>
</template>
