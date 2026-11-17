<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Icon from '@/components/ui/Icon.vue'
import { useAuthStore } from '@/stores/auth'
import { useShopStore } from '@/stores/shop'

const auth = useAuthStore()
const shop = useShopStore()
const router = useRouter()
const route = useRoute()
const open = ref(false)
const links = [
  { to: '/', label: 'Overview', store: false, icon: 'grid' },
  { to: '/products', label: 'Products', store: true, icon: 'box' },
  { to: '/collections', label: 'Collections', store: true, icon: 'layers' },
  { to: '/orders', label: 'Orders', store: true, icon: 'clipboard' },
  { to: '/builder', label: 'Builder', store: true, icon: 'panel' },
  { to: '/settings', label: 'Settings', store: true, icon: 'gear' },
] as const
const visibleLinks = computed(() => links.filter((link) => !link.store || shop.shop))

function active(to: string) {
  return route.path === to || (to !== '/' && route.path.startsWith(to))
}

async function signOut() {
  await auth.logout()
  shop.reset()
  await router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-bg text-text md:grid md:grid-cols-[200px_minmax(0,1fr)]">
    <aside class="bg-surface md:sticky md:top-0 md:h-screen md:border-r md:border-border">
      <div class="flex h-12 items-center justify-between border-b border-border px-3">
        <RouterLink to="/" class="text-sm font-semibold text-text" @click="open = false">Sable</RouterLink>
        <button
          class="inline-flex h-8 w-8 items-center justify-center rounded border border-border text-text hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent md:hidden"
          type="button"
          :aria-label="open ? 'Close menu' : 'Open menu'"
          @click="open = !open"
        >
          <Icon :name="open ? 'x' : 'menu'" />
        </button>
      </div>
      <nav class="space-y-1 border-b border-border px-2 py-2 md:border-b-0" :class="open ? 'block' : 'hidden md:block'">
        <RouterLink
          v-for="link in visibleLinks"
          :key="link.to"
          :to="link.to"
          class="flex h-8 items-center gap-2 rounded px-2 text-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          :class="active(link.to) ? 'bg-accent text-white' : 'text-text hover:bg-bg'"
          @click="open = false"
        >
          <Icon :name="link.icon" />
          {{ link.label }}
        </RouterLink>
      </nav>
    </aside>
    <div class="min-w-0">
      <header class="flex min-h-12 flex-wrap items-center justify-between gap-x-3 gap-y-2 border-b border-border px-4 py-2 md:h-12 md:flex-nowrap md:px-6 md:py-0">
        <div class="flex min-w-0 items-center gap-2">
          <p class="truncate text-sm font-medium text-text">{{ shop.shop?.name || 'No store yet' }}</p>
          <Badge v-if="shop.shop" :tone="shop.shop.published ? 'good' : 'neutral'">{{ shop.shop.published ? 'Published' : 'Unpublished' }}</Badge>
          <p v-if="shop.shop" class="truncate text-xs text-muted">/{{ shop.shop.slug }}</p>
        </div>
        <Button class="shrink-0" variant="ghost" @click="signOut"><Icon name="logout" /> Sign out</Button>
      </header>
      <main class="px-4 py-4 md:px-6 md:py-5"><RouterView /></main>
    </div>
  </div>
</template>
