<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { formatMoney } from '@/features/money'
import type { Product, PublicStore, Section, StorefrontConfig } from '@/types'

const props = withDefaults(
  defineProps<{
    store: PublicStore
    products: Product[]
    paths: { home: string; cart: string; checkout: string; product: (id: string) => string }
    cartCount?: number
    interactive?: boolean
    showSections?: boolean
  }>(),
  { cartCount: 0, interactive: true, showSections: true },
)

const config = computed(() => props.store.storefront)
const catalog = computed(() => props.products.filter((product) => product.published))
const category = ref('all')
const categories = computed(() => [...new Set(catalog.value.map((product) => product.category).filter(Boolean))])

const frame = computed(() => {
  const theme = config.value.theme
  const radius = { none: '0px', sm: '10px', md: '18px', lg: '28px' }[theme.radius] ?? '18px'
  return {
    background: theme.colors.background,
    color: theme.colors.text,
    fontFamily: `"${theme.fontBody}", Figtree, sans-serif`,
    '--sf-bg': theme.colors.background,
    '--sf-surface': theme.colors.surface,
    '--sf-text': theme.colors.text,
    '--sf-muted': theme.colors.muted,
    '--sf-primary': theme.colors.primary,
    '--sf-on': theme.colors.onPrimary,
    '--sf-accent': theme.colors.accent,
    '--sf-border': theme.colors.border,
    '--sf-radius': radius,
    '--sf-heading': `"${theme.fontHeading}", Georgia, serif`,
  }
})

function hrefFor(target: string) {
  if (target === 'cart') return props.paths.cart
  if (target === 'home') return props.paths.home
  if (target === 'about') return `${props.paths.home}#about`
  return `${props.paths.home}#products`
}

function productsFor(section: Section) {
  const limit = section.type === 'featured_products' ? section.featured?.limit ?? 4 : section.productGrid?.limit ?? 12
  const source = section.type === 'product_grid' && category.value !== 'all' ? catalog.value.filter((item) => item.category === category.value) : catalog.value
  return source.slice(0, limit)
}

function ratioClass() {
  if (config.value.theme.imageRatio === 'portrait') return 'aspect-[3/4]'
  if (config.value.theme.imageRatio === 'landscape') return 'aspect-[4/3]'
  return 'aspect-square'
}

function columnsClass(columns = 3) {
  if (columns <= 2) return 'grid-cols-1 sm:grid-cols-2'
  if (columns >= 4) return 'grid-cols-2 lg:grid-cols-4'
  return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
}

function heroSize(size: string) {
  if (size === 'sm') return 'py-12 md:py-16'
  if (size === 'lg') return 'py-20 md:py-32'
  return 'py-16 md:py-24'
}

function heroType(size: string) {
  if (size === 'sm') return 'text-4xl md:text-5xl'
  if (size === 'lg') return 'text-5xl md:text-7xl'
  return 'text-5xl md:text-6xl'
}
</script>

<template>
  <div class="min-h-screen" :style="frame">
    <div
      v-if="config.header.showAnnouncement && config.header.announcement"
      class="px-4 py-2 text-center text-xs tracking-wide"
      :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)' }"
    >
      {{ config.header.announcement }}
    </div>
    <header class="sticky top-0 z-20 border-b backdrop-blur" :style="{ borderColor: 'var(--sf-border)', background: 'color-mix(in srgb, var(--sf-bg) 88%, transparent)' }">
      <div class="mx-auto flex max-w-6xl items-center gap-4 px-4 py-4" :class="config.header.alignment === 'center' ? 'justify-center text-center' : 'justify-between'">
        <component :is="interactive ? RouterLink : 'span'" :to="paths.home" class="font-semibold tracking-tight" :style="{ fontFamily: 'var(--sf-heading)', fontSize: '1.45rem' }">
          {{ store.name }}
        </component>
        <nav class="flex flex-wrap items-center gap-4 text-sm" :style="{ color: 'var(--sf-muted)' }">
          <template v-for="item in config.header.nav" :key="item.target">
            <component :is="interactive ? RouterLink : 'span'" :to="hrefFor(item.target)" class="hover:opacity-70">
              {{ item.label }}
              <span v-if="item.target === 'cart'"> ({{ cartCount }})</span>
            </component>
          </template>
        </nav>
      </div>
    </header>

    <main v-if="showSections">
      <template v-for="section in config.sections" :key="section.id">
        <section v-if="section.visible && section.type === 'hero' && section.hero" :class="heroSize(section.hero.size)">
          <div class="mx-auto grid max-w-6xl items-center gap-8 px-4" :class="section.hero.imageUrl ? 'md:grid-cols-2' : ''">
            <div :class="section.hero.align === 'center' && !section.hero.imageUrl ? 'mx-auto max-w-3xl text-center' : 'max-w-2xl'">
              <p v-if="section.hero.eyebrow" class="mb-3 text-xs uppercase tracking-[0.22em]" :style="{ color: 'var(--sf-accent)' }">{{ section.hero.eyebrow }}</p>
              <h1 class="leading-[0.95]" :class="heroType(section.hero.size)" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.hero.headline }}</h1>
              <p v-if="section.hero.subheadline" class="mt-4 text-base md:text-lg" :style="{ color: 'var(--sf-muted)' }">{{ section.hero.subheadline }}</p>
              <component
                :is="interactive ? RouterLink : 'span'"
                :to="hrefFor(section.hero.ctaTarget)"
                class="mt-6 inline-flex px-5 py-2.5 text-sm font-semibold"
                :style="{ background: 'var(--sf-primary)', color: 'var(--sf-on)', borderRadius: 'var(--sf-radius)' }"
              >
                {{ section.hero.ctaLabel }}
              </component>
            </div>
            <img v-if="section.hero.imageUrl" :src="section.hero.imageUrl" alt="" class="w-full object-cover" :class="ratioClass()" :style="{ borderRadius: 'var(--sf-radius)' }" />
          </div>
        </section>

        <section v-else-if="section.visible && section.type === 'announcement' && section.announcement" class="px-4 py-8 text-center">
          <p class="text-sm tracking-wide" :style="{ color: 'var(--sf-muted)' }">{{ section.announcement.text }}</p>
        </section>

        <section v-else-if="section.visible && (section.type === 'featured_products' || section.type === 'product_grid')" :id="section.type === 'product_grid' ? 'products' : undefined" class="mx-auto max-w-6xl px-4 py-12">
          <div class="mb-6 flex items-end justify-between gap-4">
            <h2 class="text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">
              {{ section.type === 'featured_products' ? section.featured?.title : section.productGrid?.title }}
            </h2>
          </div>
          <div v-if="productsFor(section).length === 0" class="rounded-2xl border border-dashed px-4 py-10 text-sm" :style="{ borderColor: 'var(--sf-border)', color: 'var(--sf-muted)' }">
            Published products from this shop will show up here.
          </div>
          <div v-else class="grid gap-5" :class="columnsClass(section.type === 'product_grid' ? section.productGrid?.columns : 3)">
            <component
              :is="interactive ? RouterLink : 'div'"
              v-for="product in productsFor(section)"
              :key="product.id"
              :to="paths.product(product.id)"
              class="group block"
            >
              <div class="overflow-hidden" :style="{ background: 'var(--sf-surface)', borderRadius: 'var(--sf-radius)' }">
                <img v-if="product.images[0]" :src="product.images[0]" :alt="product.name" class="w-full object-cover transition duration-300 group-hover:scale-[1.02]" :class="ratioClass()" />
                <div v-else class="grid place-items-center text-3xl" :class="ratioClass()" :style="{ color: 'var(--sf-muted)', fontFamily: 'var(--sf-heading)' }">{{ product.name.slice(0, 1) }}</div>
              </div>
              <div class="mt-3 flex items-start justify-between gap-3">
                <div>
                  <h3 class="text-base font-medium">{{ product.name }}</h3>
                  <p v-if="product.category" class="text-xs" :style="{ color: 'var(--sf-muted)' }">{{ product.category }}</p>
                </div>
                <p class="text-sm">{{ formatMoney(product.priceCents, store.currency) }}</p>
              </div>
            </component>
          </div>
        </section>

        <section v-else-if="section.visible && section.type === 'categories' && section.categories" class="mx-auto max-w-6xl px-4 py-4">
          <h2 class="mb-4 text-2xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.categories.title }}</h2>
          <div v-if="categories.length === 0" class="text-sm" :style="{ color: 'var(--sf-muted)' }">Categories appear when products have them.</div>
          <div v-else class="flex flex-wrap gap-2">
            <button type="button" class="border px-3 py-1.5 text-sm" :style="{ borderColor: 'var(--sf-border)', borderRadius: '999px', background: category === 'all' ? 'var(--sf-primary)' : 'transparent', color: category === 'all' ? 'var(--sf-on)' : 'var(--sf-text)' }" @click="category = 'all'">All</button>
            <button v-for="item in categories" :key="item" type="button" class="border px-3 py-1.5 text-sm" :style="{ borderColor: 'var(--sf-border)', borderRadius: '999px', background: category === item ? 'var(--sf-primary)' : 'transparent', color: category === item ? 'var(--sf-on)' : 'var(--sf-text)' }" @click="category = item">{{ item }}</button>
          </div>
        </section>

        <section v-else-if="section.visible && section.type === 'about' && section.about" id="about" class="mx-auto grid max-w-6xl gap-8 px-4 py-14 md:grid-cols-2">
          <div>
            <h2 class="text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.about.title }}</h2>
            <p class="mt-4 whitespace-pre-line leading-relaxed" :style="{ color: 'var(--sf-muted)' }">{{ section.about.body }}</p>
          </div>
          <img v-if="section.about.imageUrl" :src="section.about.imageUrl" alt="" class="aspect-[4/3] w-full object-cover" :style="{ borderRadius: 'var(--sf-radius)' }" />
        </section>

        <section v-else-if="section.visible && section.type === 'newsletter' && section.newsletter" class="mx-auto max-w-3xl px-4 py-14 text-center">
          <h2 class="text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.newsletter.title }}</h2>
          <p class="mt-3" :style="{ color: 'var(--sf-muted)' }">{{ section.newsletter.text }}</p>
        </section>
      </template>
    </main>
    <main v-else class="mx-auto max-w-6xl px-4 py-10"><slot /></main>

    <footer class="mt-8 border-t" :style="{ borderColor: 'var(--sf-border)' }">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-4 py-8 text-sm" :style="{ color: 'var(--sf-muted)' }">
        <p>{{ config.footer.text }}</p>
        <div class="flex gap-4">
          <component :is="interactive ? RouterLink : 'span'" v-for="item in config.footer.links" :key="item.target" :to="hrefFor(item.target)">{{ item.label }}</component>
        </div>
      </div>
    </footer>
  </div>
</template>
