<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Icon from '@/components/ui/Icon.vue'
import Logo from '@/components/ui/Logo.vue'
import { useAuthStore } from '@/stores/auth'
import { useShopStore } from '@/stores/shop'

const auth = useAuthStore()
const shop = useShopStore()
const router = useRouter()
const route = useRoute()
const open = ref(false)
const accountOpen = ref(false)
const collapsed = ref(readCollapsed())
const links = [
  { to: '/', label: 'Overview', store: false, icon: 'grid' },
  { to: '/products', label: 'Products', store: true, icon: 'box' },
  { to: '/collections', label: 'Collections', store: true, icon: 'layers' },
  { to: '/orders', label: 'Orders', store: true, icon: 'clipboard' },
  { to: '/builder', label: 'Builder', store: true, icon: 'panel' },
  { to: '/settings', label: 'Settings', store: true, icon: 'gear' },
] as const
const visibleLinks = computed(() => links.filter((link) => !link.store || shop.shop))
const isBuilder = computed(() => route.path === '/builder')
const initials = computed(() => {
  const name = auth.user?.name?.trim() ?? ''
  if (!name) return '?'
  const parts = name.split(/\s+/).filter((part) => part.length > 0)
  const first = parts[0]
  const last = parts[parts.length - 1]
  if (parts.length >= 2 && first && last) return `${first.charAt(0)}${last.charAt(0)}`.toUpperCase()
  return name.slice(0, 2).toUpperCase()
})

function readCollapsed() {
  try {
    const saved = localStorage.getItem('sable-sidebar')
    if (saved === 'collapsed') return true
    if (saved === 'open') return false
  } catch {
    return false
  }
  return window.matchMedia('(max-width: 1279px)').matches
}

function active(to: string) {
  return route.path === to || (to !== '/' && route.path.startsWith(to))
}

function toggleNav() {
  open.value = !open.value
  if (open.value) accountOpen.value = false
}

function toggleAccount() {
  accountOpen.value = !accountOpen.value
  if (accountOpen.value) open.value = false
}

function toggleSidebar() {
  collapsed.value = !collapsed.value
  accountOpen.value = false
  try {
    localStorage.setItem('sable-sidebar', collapsed.value ? 'collapsed' : 'open')
  } catch {
    /* preference is optional */
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') accountOpen.value = false
}

async function signOut() {
  await auth.logout()
  shop.reset()
  await router.push('/login')
}

watch(
  () => route.path,
  () => {
    open.value = false
    accountOpen.value = false
  },
)

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="min-h-screen bg-bg text-text md:flex">
    <div
      class="sticky top-0 flex h-12 items-center justify-between border-b border-border bg-surface px-3 md:hidden"
      :class="accountOpen ? 'z-40' : 'z-30'"
    >
      <RouterLink to="/" class="inline-flex items-center gap-2 text-text" @click="open = false">
        <Logo />
        <span class="text-sm font-medium leading-none text-text">Sable</span>
      </RouterLink>
      <div class="flex items-center gap-1">
        <div class="relative">
          <button
            class="inline-flex h-8 w-8 items-center justify-center rounded bg-text text-[11px] font-medium text-surface focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            type="button"
            aria-label="Account"
            aria-haspopup="menu"
            :aria-expanded="accountOpen"
            @click="toggleAccount"
          >
            {{ initials }}
          </button>
          <div
            v-if="accountOpen"
            class="absolute right-0 top-9 z-40 w-[196px] rounded border border-border bg-surface p-1 shadow-[0_8px_24px_rgba(27,23,20,0.12)]"
            role="menu"
          >
            <div class="px-2 py-1">
              <p class="truncate text-sm font-medium text-text">{{ auth.user?.name }}</p>
              <p class="truncate text-xs text-muted">{{ auth.user?.email }}</p>
            </div>
            <div class="my-1 border-t border-border" role="separator"></div>
            <RouterLink
              class="flex h-8 w-full items-center gap-2 rounded px-2 text-sm hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              to="/account"
              role="menuitem"
            >
              <Icon name="user" /> Account
            </RouterLink>
            <button
              class="flex h-8 w-full items-center gap-2 rounded px-2 text-left text-sm hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              type="button"
              @click="signOut"
            >
              <Icon name="logout" /> Sign out
            </button>
          </div>
        </div>
        <button
          class="inline-flex h-8 w-8 items-center justify-center rounded border border-border text-text hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          type="button"
          :aria-label="open ? 'Close menu' : 'Open menu'"
          @click="toggleNav"
        >
          <Icon :name="open ? 'x' : 'menu'" />
        </button>
      </div>
    </div>

    <aside
      class="bg-surface transition-[width] duration-200 ease-out md:sticky md:top-0 md:flex md:h-dvh md:shrink-0 md:flex-col md:border-r md:border-border"
      :class="[collapsed ? 'md:w-14' : 'md:w-[220px]', accountOpen ? 'z-40' : '']"
    >
      <div class="hidden min-h-12 items-center gap-2 px-2 py-2 md:flex" :class="collapsed ? 'flex-col' : 'flex-row-reverse justify-between'">
        <button
          class="inline-flex h-8 w-8 items-center justify-center rounded text-muted hover:bg-bg hover:text-text focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          type="button"
          :aria-label="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          :title="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          @click="toggleSidebar"
        >
          <Icon name="sidebar" />
        </button>
        <RouterLink to="/" class="inline-flex items-center gap-2 text-text" :title="collapsed ? 'Sable' : undefined">
          <Logo />
          <span class="text-sm font-medium leading-none text-text" :class="collapsed ? 'md:sr-only' : ''">Sable</span>
        </RouterLink>
      </div>

      <div class="hidden border-b border-border px-3 py-2 md:block" :class="collapsed ? 'md:px-2 md:py-2' : ''">
        <p class="truncate text-sm font-medium text-text" :class="collapsed ? 'md:sr-only' : ''">{{ shop.shop?.name || 'No store yet' }}</p>
        <div class="mt-1 flex min-w-0 items-center gap-2" :class="collapsed ? 'md:mt-0 md:justify-center' : ''">
          <span
            class="h-1.5 w-1.5 shrink-0 rounded-full"
            :class="shop.shop?.published ? 'bg-good-fg' : 'bg-border'"
            :title="shop.shop ? (shop.shop.published ? 'Published' : 'Unpublished') : 'No store yet'"
          ></span>
          <p v-if="shop.shop" class="min-w-0 truncate text-xs text-muted" :class="collapsed ? 'md:sr-only' : ''">
            {{ shop.shop.published ? 'Published' : 'Unpublished' }} · /{{ shop.shop.slug }}
          </p>
        </div>
      </div>

      <nav class="min-h-0 flex-1 space-y-0.5 overflow-y-auto px-2 py-2" :class="open ? 'block' : 'hidden md:block'">
        <div class="-mx-2 -mt-2 border-b border-border px-3 py-2 md:hidden">
          <p class="truncate text-sm font-medium text-text">{{ shop.shop?.name || 'No store yet' }}</p>
          <p v-if="shop.shop" class="mt-1 truncate text-xs text-muted">
            {{ shop.shop.published ? 'Published' : 'Unpublished' }} · /{{ shop.shop.slug }}
          </p>
        </div>
        <RouterLink
          v-for="link in visibleLinks"
          :key="link.to"
          :to="link.to"
          class="flex h-9 items-center gap-2.5 rounded px-2 text-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          :class="[
            active(link.to) ? 'bg-accent text-white' : 'text-text hover:bg-bg',
            collapsed ? 'md:justify-center md:px-0' : '',
          ]"
          :aria-label="link.label"
          :title="collapsed ? link.label : undefined"
          @click="open = false; accountOpen = false"
        >
          <Icon :name="link.icon" />
          <span :class="collapsed ? 'md:sr-only' : ''">{{ link.label }}</span>
        </RouterLink>
      </nav>

      <div class="relative mt-auto hidden border-t border-border px-2 py-2 md:block" :class="accountOpen ? 'z-40' : ''">
        <button
          class="flex h-9 w-full items-center gap-2.5 rounded px-1 text-left text-sm hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          :class="collapsed ? 'md:justify-center md:px-0' : ''"
          type="button"
          :aria-label="auth.user?.name || 'Account'"
          aria-haspopup="menu"
          :aria-expanded="accountOpen"
          :title="collapsed ? auth.user?.name : undefined"
          @click="toggleAccount"
        >
          <span class="inline-flex h-6 w-6 shrink-0 items-center justify-center rounded bg-text text-[11px] font-medium text-surface">{{ initials }}</span>
          <span class="min-w-0 truncate" :class="collapsed ? 'md:sr-only' : ''">{{ auth.user?.name }}</span>
        </button>
        <div
          v-if="accountOpen"
          class="absolute z-40 w-[196px] rounded border border-border bg-surface p-1 shadow-[0_8px_24px_rgba(27,23,20,0.12)]"
          :class="collapsed ? 'bottom-0 left-full ml-2' : 'bottom-11 left-2'"
          role="menu"
        >
          <div class="px-2 py-1">
            <p class="truncate text-sm font-medium text-text">{{ auth.user?.name }}</p>
            <p class="truncate text-xs text-muted">{{ auth.user?.email }}</p>
          </div>
          <div class="my-1 border-t border-border" role="separator"></div>
          <RouterLink
            class="flex h-8 w-full items-center gap-2 rounded px-2 text-sm hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            to="/account"
            role="menuitem"
          >
            <Icon name="user" /> Account
          </RouterLink>
          <button
            class="flex h-8 w-full items-center gap-2 rounded px-2 text-left text-sm hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            type="button"
            @click="signOut"
          >
            <Icon name="logout" /> Sign out
          </button>
        </div>
      </div>
    </aside>

    <div class="min-w-0 flex-1" :class="isBuilder ? 'lg:h-dvh lg:overflow-hidden' : ''">
      <main :class="isBuilder ? 'p-0 lg:h-full' : 'px-4 py-4 md:px-6 md:py-5'">
        <RouterView />
      </main>
    </div>

    <button
      v-if="accountOpen"
      class="fixed inset-0 z-30"
      type="button"
      aria-label="Close account"
      @click="accountOpen = false"
    ></button>
  </div>
</template>
