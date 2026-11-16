<script setup lang="ts">
import { computed, onMounted, ref, toRaw } from 'vue'
import Button from '@/components/ui/Button.vue'
import StorefrontView from '@/components/storefront/StorefrontView.vue'
import { controlClass } from '@/components/ui/control'
import { listCollections } from '@/features/collections/api'
import { listProducts } from '@/features/products/api'
import { generateDesign, publishDesign, refineDesign, saveDesign } from '@/features/storefront-builder/api'
import { publicUrl, storePaths } from '@/features/stores/host'
import { useShopStore } from '@/stores/shop'
import type { Collection, Product, Section, StorefrontConfig } from '@/types'

const shop = useShopStore()
const products = ref<Product[]>([])
const collections = ref<Collection[]>([])
const draft = ref<StorefrontConfig | null>(null)
const prompt = ref('A quiet clothing shop. Black and white, large product photos, a short hero, and clear navigation.')
const note = ref('')
const error = ref('')
const busy = ref('')
const selected = ref(0)
const panel = ref<'look' | 'page' | 'sections'>('look')
const frame = ref<'full' | 'phone'>('full')
const examples = [
  { label: 'Minimal', prompt: 'A quiet clothing shop. Black and white, large product photos, a short hero, and clear navigation.' },
  { label: 'Darker', prompt: 'Make the background darker and keep the rest of the layout.' },
  { label: 'Reorder', prompt: 'Move the featured products above the about section.' },
  { label: 'Luxury', prompt: 'Use a luxury theme and a larger hero.' },
]
const colorKeys = ['background', 'surface', 'text', 'muted', 'primary', 'onPrimary', 'accent', 'border'] as const

const shopUrl = computed(() => (shop.shop ? publicUrl(shop.shop.slug) : null))
const previewStore = computed(() => {
  if (!shop.shop || !draft.value) return null
  return {
    name: shop.shop.name,
    slug: shop.shop.slug,
    description: shop.shop.description,
    logoUrl: shop.shop.logoUrl,
    currency: shop.shop.currency,
    storefront: draft.value,
  }
})
const section = computed(() => draft.value?.sections[selected.value])
const fonts = ['Instrument Serif', 'Fraunces', 'Newsreader', 'Outfit', 'Figtree', 'Literata']
const themePresets: Record<string, Omit<StorefrontConfig['theme'], 'name'>> = {
  editorial: { fontHeading: 'Fraunces', fontBody: 'Figtree', radius: 'md', imageRatio: 'portrait', colors: { background: '#f6f1e8', surface: '#fffdf8', text: '#1f1a17', muted: '#6f655d', primary: '#9a3412', onPrimary: '#fff7ed', accent: '#b45309', border: '#e6ddd1' } },
  minimal: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'sm', imageRatio: 'square', colors: { background: '#ffffff', surface: '#fafafa', text: '#171717', muted: '#737373', primary: '#171717', onPrimary: '#ffffff', accent: '#404040', border: '#e5e5e5' } },
  mono: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'sm', imageRatio: 'square', colors: { background: '#ffffff', surface: '#f4f4f5', text: '#111111', muted: '#52525b', primary: '#111111', onPrimary: '#ffffff', accent: '#18181b', border: '#e4e4e7' } },
  luxury: { fontHeading: 'Instrument Serif', fontBody: 'Figtree', radius: 'none', imageRatio: 'portrait', colors: { background: '#100e0c', surface: '#1a1714', text: '#f6f1e7', muted: '#b7ab9a', primary: '#c6a15b', onPrimary: '#1a1408', accent: '#e7d3a1', border: '#3a342c' } },
  bold: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'lg', imageRatio: 'square', colors: { background: '#fff7ed', surface: '#ffffff', text: '#1c1917', muted: '#78716c', primary: '#ea580c', onPrimary: '#ffffff', accent: '#0f172a', border: '#fed7aa' } },
  playful: { fontHeading: 'Fraunces', fontBody: 'Outfit', radius: 'lg', imageRatio: 'square', colors: { background: '#fdf2f8', surface: '#ffffff', text: '#3b0764', muted: '#7e22ce', primary: '#db2777', onPrimary: '#ffffff', accent: '#0891b2', border: '#fbcfe8' } },
  dark: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'md', imageRatio: 'portrait', colors: { background: '#121212', surface: '#1c1c1c', text: '#f5f5f5', muted: '#a3a3a3', primary: '#f5f5f5', onPrimary: '#111111', accent: '#d4d4d4', border: '#2e2e2e' } },
}

function cloneConfig(config: StorefrontConfig): StorefrontConfig {
  return structuredClone(toRaw(config))
}

onMounted(async () => {
  try {
    const current = await shop.load()
    if (!current) return
    draft.value = cloneConfig(current.draft)
    note.value = current.designNote?.summary || ''
    const [catalog, groups] = await Promise.all([listProducts({ status: 'published' }), listCollections()])
    products.value = catalog.items
    collections.value = groups.items.filter((item) => item.published)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not open the builder'
  }
})

async function run(kind: 'generate' | 'refine') {
  busy.value = kind
  error.value = ''
  try {
    const result = kind === 'generate' ? await generateDesign(prompt.value) : await refineDesign(prompt.value)
    shop.set(result.store)
    draft.value = cloneConfig(result.store.draft)
    note.value = result.summary
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'The designer failed'
  } finally {
    busy.value = ''
  }
}

async function save() {
  if (!draft.value) return
  busy.value = 'save'
  error.value = ''
  try {
    shop.set(await saveDesign(draft.value))
    draft.value = cloneConfig(shop.shop!.draft)
    note.value = 'Draft saved. Publish it when the preview looks right.'
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not save the draft'
  } finally {
    busy.value = ''
  }
}

async function publish() {
  busy.value = 'publish'
  error.value = ''
  try {
    if (draft.value) shop.set(await saveDesign(draft.value))
    shop.set(await publishDesign())
    draft.value = cloneConfig(shop.shop!.draft)
    note.value = shop.shop!.published ? 'This design is live.' : 'Design published. The shop itself is still unpublished in settings.'
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not publish the design'
  } finally {
    busy.value = ''
  }
}

function pickTheme(name: string) {
  const preset = themePresets[name]
  if (!draft.value || !preset) return
  draft.value.theme = { name, ...structuredClone(preset) }
}

function colorValue(key: string) {
  if (!draft.value || !(key in draft.value.theme.colors)) return '#000000'
  return draft.value.theme.colors[key as keyof typeof draft.value.theme.colors]
}

function onColor(key: string, event: Event) {
  setColor(key, (event.target as HTMLInputElement).value)
}

function setColor(key: string, value: string) {
  if (!draft.value || !(key in draft.value.theme.colors)) return
  draft.value.theme.colors[key as keyof typeof draft.value.theme.colors] = value
}

function move(index: number, direction: number) {
  if (!draft.value) return
  const next = index + direction
  if (next < 0 || next >= draft.value.sections.length) return
  const [item] = draft.value.sections.splice(index, 1)
  draft.value.sections.splice(next, 0, item as Section)
  selected.value = next
}
</script>

<template>
  <p v-if="error && !draft" class="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
  <div v-if="draft && previewStore && shop.shop" class="grid items-start gap-6 xl:grid-cols-[340px_1fr]">
    <div class="space-y-4">
      <div class="flex items-end justify-between gap-3">
        <div>
          <h1 class="font-serif text-4xl">Builder</h1>
          <p class="mt-1 text-sm text-stone-600">The preview is the draft. Publishing updates the shop host.</p>
        </div>
      </div>
      <div v-if="shopUrl" class="flex flex-wrap gap-x-4 gap-y-1 text-sm">
        <a class="font-medium text-rust" :href="shopUrl.subdomain" target="_blank" rel="noreferrer">Open shop</a>
        <a class="text-stone-500" :href="shopUrl.path" target="_blank" rel="noreferrer">Same-host preview</a>
      </div>
      <p v-if="error" class="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
      <p v-if="note" class="rounded-xl bg-card px-3 py-2 text-sm text-stone-700">{{ note }}</p>
      <textarea v-model="prompt" :class="controlClass" rows="3" />
      <div class="flex flex-wrap gap-2">
        <button v-for="example in examples" :key="example.label" type="button" class="rounded-full border border-line bg-white px-3 py-1 text-xs text-stone-600" @click="prompt = example.prompt">{{ example.label }}</button>
      </div>
      <div class="flex flex-wrap gap-2">
        <Button :loading="busy === 'generate'" @click="run('generate')">Generate</Button>
        <Button variant="secondary" :loading="busy === 'refine'" @click="run('refine')">Refine</Button>
        <Button variant="secondary" :loading="busy === 'save'" @click="save">Save</Button>
        <Button variant="secondary" :loading="busy === 'publish'" @click="publish">Publish</Button>
      </div>

      <div class="grid grid-cols-3 rounded-xl border border-line bg-white p-1 text-sm">
        <button v-for="item in (['look', 'page', 'sections'] as const)" :key="item" type="button" class="rounded-lg px-2 py-1.5 capitalize" :class="panel === item ? 'bg-ink text-white' : 'text-stone-600'" @click="panel = item">{{ item }}</button>
      </div>

      <div v-if="panel === 'look'" class="space-y-3">
        <div class="flex flex-wrap gap-2">
          <button v-for="name in Object.keys(themePresets)" :key="name" type="button" class="rounded-full border px-3 py-1 text-xs capitalize" :class="draft.theme.name === name ? 'border-ink bg-ink text-white' : 'border-line bg-white text-stone-600'" @click="pickTheme(name)">{{ name }}</button>
        </div>
        <div class="grid grid-cols-2 gap-2">
          <label class="text-xs text-stone-500">Images
            <select v-model="draft.theme.imageRatio" :class="controlClass" class="mt-1">
              <option value="square">Square</option>
              <option value="portrait">Portrait</option>
              <option value="landscape">Landscape</option>
            </select>
          </label>
          <label class="text-xs text-stone-500">Corners
            <select v-model="draft.theme.radius" :class="controlClass" class="mt-1">
              <option value="none">Sharp</option>
              <option value="sm">Small</option>
              <option value="md">Medium</option>
              <option value="lg">Large</option>
            </select>
          </label>
          <label class="text-xs text-stone-500">Heading
            <select v-model="draft.theme.fontHeading" :class="controlClass" class="mt-1"><option v-for="font in fonts" :key="font">{{ font }}</option></select>
          </label>
          <label class="text-xs text-stone-500">Body
            <select v-model="draft.theme.fontBody" :class="controlClass" class="mt-1"><option v-for="font in fonts" :key="font">{{ font }}</option></select>
          </label>
        </div>
        <div class="grid grid-cols-4 gap-2">
          <label v-for="key in colorKeys" :key="key" class="text-[10px] uppercase tracking-wide text-stone-500">
            {{ key }}
            <input :value="colorValue(key)" type="color" class="mt-1 h-8 w-full cursor-pointer rounded border border-line bg-white" @input="onColor(key, $event)" />
          </label>
        </div>
      </div>

      <div v-else-if="panel === 'page'" class="space-y-2">
        <label class="flex items-center gap-2 text-sm"><input v-model="draft.header.showAnnouncement" type="checkbox" /> Show the bar above the header</label>
        <input v-model="draft.header.announcement" :class="controlClass" placeholder="Announcement" />
        <select v-model="draft.header.alignment" :class="controlClass">
          <option value="left">Logo on the left</option>
          <option value="center">Centered header</option>
        </select>
        <label v-for="item in draft.header.nav" :key="item.target" class="block text-xs capitalize text-stone-500">{{ item.target }}
          <input v-model="item.label" :class="controlClass" class="mt-1" />
        </label>
        <label class="block text-xs text-stone-500">Footer
          <input v-model="draft.footer.text" :class="controlClass" class="mt-1" />
        </label>
      </div>

      <div v-else class="space-y-3">
        <ul class="space-y-2">
          <li v-for="(item, index) in draft.sections" :key="item.id" class="rounded-xl border bg-card p-3" :class="selected === index ? 'border-rust' : 'border-line'">
            <button class="flex w-full items-center justify-between text-left text-sm" type="button" @click="selected = index">
              <span class="capitalize">{{ item.type.replaceAll('_', ' ') }}</span>
              <span class="text-xs text-stone-500">{{ item.visible ? 'Shown' : 'Hidden' }}</span>
            </button>
            <div class="mt-2 flex gap-3 text-xs">
              <button type="button" @click="item.visible = !item.visible">{{ item.visible ? 'Hide' : 'Show' }}</button>
              <button type="button" :disabled="index === 0" @click="move(index, -1)">Up</button>
              <button type="button" :disabled="index === draft.sections.length - 1" @click="move(index, 1)">Down</button>
            </div>
          </li>
        </ul>
        <div v-if="section?.hero" class="space-y-2">
          <input v-model="section.hero.eyebrow" :class="controlClass" placeholder="Eyebrow" />
          <input v-model="section.hero.headline" :class="controlClass" placeholder="Headline" />
          <textarea v-model="section.hero.subheadline" :class="controlClass" rows="2" placeholder="Subheadline" />
          <input v-model="section.hero.ctaLabel" :class="controlClass" placeholder="Button label" />
          <select v-model="section.hero.ctaTarget" :class="controlClass"><option value="products">Shop</option><option value="about">About</option><option value="cart">Cart</option><option value="home">Home</option></select>
          <div class="grid grid-cols-2 gap-2">
            <select v-model="section.hero.size" :class="controlClass"><option value="sm">Small</option><option value="md">Medium</option><option value="lg">Large</option></select>
            <select v-model="section.hero.align" :class="controlClass"><option value="left">Left</option><option value="center">Center</option></select>
          </div>
          <input v-model="section.hero.imageUrl" :class="controlClass" placeholder="Hero image URL" />
        </div>
        <div v-else-if="section?.productGrid" class="space-y-2">
          <input v-model="section.productGrid.title" :class="controlClass" placeholder="Title" />
          <select v-model.number="section.productGrid.columns" :class="controlClass"><option :value="2">2 columns</option><option :value="3">3 columns</option><option :value="4">4 columns</option></select>
          <label class="block text-xs text-stone-500">Products shown
            <input v-model.number="section.productGrid.limit" :class="controlClass" class="mt-1" type="number" min="1" max="24" />
          </label>
        </div>
        <div v-else-if="section?.featured" class="space-y-2">
          <input v-model="section.featured.title" :class="controlClass" placeholder="Title" />
          <label class="block text-xs text-stone-500">Products shown
            <input v-model.number="section.featured.limit" :class="controlClass" class="mt-1" type="number" min="1" max="8" />
          </label>
        </div>
        <div v-else-if="section?.about" class="space-y-2">
          <input v-model="section.about.title" :class="controlClass" placeholder="Title" />
          <textarea v-model="section.about.body" :class="controlClass" rows="4" placeholder="About the shop" />
          <input v-model="section.about.imageUrl" :class="controlClass" placeholder="Image URL" />
        </div>
        <div v-else-if="section?.announcement"><input v-model="section.announcement.text" :class="controlClass" placeholder="Announcement" /></div>
        <div v-else-if="section?.categories"><input v-model="section.categories.title" :class="controlClass" placeholder="Title" /></div>
        <div v-else-if="section?.newsletter" class="space-y-2">
          <input v-model="section.newsletter.title" :class="controlClass" placeholder="Title" />
          <textarea v-model="section.newsletter.text" :class="controlClass" rows="2" />
          <p class="text-xs text-stone-500">Copy only. This does not collect email addresses.</p>
        </div>
      </div>
    </div>
    <div class="xl:sticky xl:top-4">
      <div class="mb-2 flex justify-end gap-2 text-xs">
        <button type="button" class="rounded-full border px-3 py-1" :class="frame === 'full' ? 'border-ink bg-ink text-white' : 'border-line bg-white'" @click="frame = 'full'">Desktop</button>
        <button type="button" class="rounded-full border px-3 py-1" :class="frame === 'phone' ? 'border-ink bg-ink text-white' : 'border-line bg-white'" @click="frame = 'phone'">Phone</button>
      </div>
      <div class="overflow-hidden rounded-2xl border border-line bg-white" :class="frame === 'phone' ? 'mx-auto max-w-[390px]' : ''">
        <div class="max-h-[calc(100vh-5rem)] overflow-auto">
          <StorefrontView :store="previewStore" :products="products" :collections="collections" :paths="storePaths(shop.shop.slug)" :interactive="false" :compact="frame === 'phone'" />
        </div>
      </div>
    </div>
  </div>
</template>
