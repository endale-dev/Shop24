<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Icon from '@/components/ui/Icon.vue'
import { formatMoney } from '@/features/money'
import { priceVaries } from '@/features/products/api'
import type { Collection, Product, PublicStore, Section } from '@/types'

const props = withDefaults(
  defineProps<{
    store: PublicStore
    products: Product[]
    collections?: Collection[]
    paths: { home: string; cart: string; checkout: string; product: (id: string) => string; collection: (slug: string) => string }
    cartCount?: number
    interactive?: boolean
    showSections?: boolean
    compact?: boolean
  }>(),
  { collections: () => [], cartCount: 0, interactive: true, showSections: true, compact: false },
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

function linkIs(target: string) {
  if (!props.interactive) return 'span'
  if (target === 'cart') return RouterLink
  return 'a'
}

function linkAttrs(target: string) {
  if (!props.interactive) return {}
  const href = hrefFor(target)
  return target === 'cart' ? { to: href } : { href }
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
  if (props.compact) return columns <= 2 ? 'grid-cols-1' : 'grid-cols-2'
  if (columns <= 2) return 'grid-cols-1 sm:grid-cols-2'
  if (columns >= 4) return 'grid-cols-2 lg:grid-cols-4'
  return 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
}

function heroSize(size: string) {
  if (props.compact) return 'py-10'
  if (size === 'sm') return 'py-12 md:py-16'
  if (size === 'lg') return 'py-16 md:py-28'
  return 'py-14 md:py-20'
}

function heroType(size: string) {
  if (props.compact) return 'text-4xl'
  if (size === 'sm') return 'text-4xl sm:text-5xl'
  if (size === 'lg') return 'text-4xl sm:text-6xl md:text-7xl'
  return 'text-4xl sm:text-5xl md:text-6xl'
}

function headerClass() {
  if (props.compact) return 'flex-col items-center text-center'
  if (config.value.header.alignment === 'center') return 'flex-col items-center text-center sm:flex-row sm:justify-center'
  return 'flex-col items-start sm:flex-row sm:items-center sm:justify-between'
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
    <header class="sticky top-0 z-20 border-b backdrop-blur" :style="{ borderColor: 'var(--sf-border)', background: 'color-mix(in srgb, var(--sf-bg) 92%, transparent)' }">
      <div class="mx-auto flex max-w-6xl gap-3 px-4 py-3 sm:gap-4 sm:py-4" :class="headerClass()">
        <component :is="linkIs('home')" v-bind="linkAttrs('home')" class="max-w-full truncate text-2xl font-semibold leading-none tracking-tight focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)] sm:text-[1.65rem]" :style="{ fontFamily: 'var(--sf-heading)' }">
          {{ store.name }}
        </component>
        <nav class="flex flex-wrap items-center gap-x-4 gap-y-1 text-sm" :style="{ color: 'var(--sf-muted)' }">
          <component :is="linkIs(item.target)" v-for="item in config.header.nav" :key="item.target" v-bind="linkAttrs(item.target)" class="inline-flex items-center gap-1 hover:opacity-70 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]">
            {{ item.label }}
            <template v-if="item.target === 'cart'">
              <Icon name="cart" />
              {{ cartCount }}
            </template>
          </component>
        </nav>
      </div>
    </header>

    <main v-if="showSections">
      <section v-if="collections.length" class="mx-auto max-w-6xl px-4 pt-8 sm:pt-10">
        <h2 class="mb-4 text-2xl sm:text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">Collections</h2>
        <div class="grid gap-4" :class="compact ? 'grid-cols-1' : 'sm:grid-cols-2 lg:grid-cols-3'">
          <component :is="interactive ? RouterLink : 'div'" v-for="item in collections" :key="item.id" :to="paths.collection(item.slug)" class="group block focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]">
            <div class="overflow-hidden" :style="{ background: 'var(--sf-surface)', borderRadius: 'var(--sf-radius)' }">
              <img v-if="item.imageUrl" :src="item.imageUrl" :alt="item.name" class="aspect-[4/3] w-full object-cover" />
              <div v-else class="grid aspect-[4/3] place-items-center text-3xl" :style="{ fontFamily: 'var(--sf-heading)', color: 'var(--sf-muted)' }">{{ item.name.slice(0, 1) }}</div>
            </div>
            <h3 class="mt-3 text-base font-medium">{{ item.name }}</h3>
          </component>
        </div>
      </section>
      <template v-for="section in config.sections" :key="section.id">
        <section v-if="section.visible && section.type === 'hero' && section.hero" :class="heroSize(section.hero.size)">
          <div class="mx-auto grid max-w-6xl items-center gap-8 px-4" :class="section.hero.imageUrl && !compact ? 'md:grid-cols-2' : ''">
            <div :class="section.hero.align === 'center' && !section.hero.imageUrl ? 'mx-auto max-w-3xl text-center' : 'max-w-2xl'">
              <p v-if="section.hero.eyebrow" class="mb-3 text-xs uppercase tracking-[0.22em]" :style="{ color: 'var(--sf-accent)' }">{{ section.hero.eyebrow }}</p>
              <h1 class="leading-[0.95]" :class="heroType(section.hero.size)" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.hero.headline }}</h1>
              <p v-if="section.hero.subheadline" class="mt-4 max-w-xl text-base leading-relaxed sm:text-lg" :style="{ color: 'var(--sf-muted)' }">{{ section.hero.subheadline }}</p>
              <component
                :is="linkIs(section.hero.ctaTarget)"
                v-bind="linkAttrs(section.hero.ctaTarget)"
                class="mt-6 inline-flex px-5 py-2.5 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]"
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
            <h2 class="text-2xl sm:text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">
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
              class="group block focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:[outline-color:var(--sf-primary)]"
            >
              <div class="overflow-hidden" :style="{ background: 'var(--sf-surface)', borderRadius: 'var(--sf-radius)' }">
                <img v-if="product.images[0]" :src="product.images[0]" :alt="product.name" class="w-full object-cover transition duration-300 group-hover:scale-[1.02]" :class="ratioClass()" />
                <div v-else class="grid place-items-center text-3xl" :class="ratioClass()" :style="{ color: 'var(--sf-muted)', fontFamily: 'var(--sf-heading)' }">{{ product.name.slice(0, 1) }}</div>
              </div>
              <div class="mt-3 flex flex-col gap-1" :class="compact ? '' : 'sm:flex-row sm:items-start sm:justify-between sm:gap-3'">
                <div>
                  <h3 class="text-base font-medium">{{ product.name }}</h3>
                  <p v-if="product.category" class="text-xs" :style="{ color: 'var(--sf-muted)' }">{{ product.category }}</p>
                </div>
                <p class="text-sm">{{ priceVaries(product) ? 'From ' : '' }}{{ formatMoney(product.priceCents, store.currency) }}</p>
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

        <section v-else-if="section.visible && section.type === 'about' && section.about" id="about" class="mx-auto grid max-w-6xl gap-8 px-4 py-12" :class="compact ? '' : 'md:grid-cols-2'">
          <div>
            <h2 class="text-2xl sm:text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.about.title }}</h2>
            <p class="mt-4 whitespace-pre-line leading-relaxed" :style="{ color: 'var(--sf-muted)' }">{{ section.about.body }}</p>
          </div>
          <img v-if="section.about.imageUrl" :src="section.about.imageUrl" alt="" class="aspect-[4/3] w-full object-cover" :style="{ borderRadius: 'var(--sf-radius)' }" />
        </section>

        <section v-else-if="section.visible && section.type === 'newsletter' && section.newsletter" class="mx-auto max-w-3xl px-4 py-14 text-center">
          <h2 class="text-2xl sm:text-3xl" :style="{ fontFamily: 'var(--sf-heading)' }">{{ section.newsletter.title }}</h2>
          <p class="mt-3" :style="{ color: 'var(--sf-muted)' }">{{ section.newsletter.text }}</p>
        </section>
      </template>
    </main>
    <main v-else class="mx-auto max-w-6xl px-4 py-10"><slot /></main>

    <footer class="mt-8 border-t" :style="{ borderColor: 'var(--sf-border)' }">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 px-4 py-8 text-sm" :style="{ color: 'var(--sf-muted)' }">
        <p>{{ config.footer.text }}</p>
        <div class="flex flex-wrap gap-x-4 gap-y-2">
          <component :is="linkIs(item.target)" v-for="item in config.footer.links" :key="item.target" v-bind="linkAttrs(item.target)">{{ item.label }}</component>
        </div>
      </div>
    </footer>
  </div>
</template>
