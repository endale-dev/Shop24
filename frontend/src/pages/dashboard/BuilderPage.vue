<script setup lang="ts">
import { computed, onMounted, ref, toRaw } from 'vue'
import Button from '@/components/ui/Button.vue'
import StorefrontView from '@/components/storefront/StorefrontView.vue'
import { controlClass } from '@/components/ui/control'
import { listProducts } from '@/features/products/api'
import { generateDesign, publishDesign, refineDesign, saveDesign } from '@/features/storefront-builder/api'
import { storePaths } from '@/features/stores/host'
import { useShopStore } from '@/stores/shop'
import type { Product, Section, StorefrontConfig } from '@/types'

const shop = useShopStore()
const products = ref<Product[]>([])
const draft = ref<StorefrontConfig | null>(null)
const prompt = ref('Create a modern minimalist clothing store with a black and white theme, large product images, a hero section, featured products, and a clean navigation bar.')
const note = ref('')
const error = ref('')
const busy = ref('')
const selected = ref(0)
const examples = [
  'Create a modern minimalist clothing store with a black and white theme, large product images, a hero section, featured products, and a clean navigation bar.',
  'Make the background darker.',
  'Move the featured products above the about section.',
  'Change the theme to a luxury style and make the hero section larger.',
]

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
    products.value = (await listProducts({ status: 'published' })).items
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

function onTheme(event: Event) {
  const name = (event.target as HTMLSelectElement).value
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
  <div v-if="draft && previewStore && shop.shop" class="grid gap-6 xl:grid-cols-[420px_1fr]">
    <div class="space-y-4">
      <div>
        <h1 class="font-serif text-4xl">Storefront builder</h1>
        <p class="mt-1 text-sm text-stone-600">The designer returns a layout, not code. Product grids use your published catalog.</p>
      </div>
      <p v-if="error" class="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-800">{{ error }}</p>
      <p v-if="note" class="rounded-xl bg-card px-3 py-2 text-sm text-stone-700">{{ note }}</p>
      <textarea v-model="prompt" :class="controlClass" rows="5" />
      <div class="flex flex-wrap gap-2">
        <Button :loading="busy === 'generate'" @click="run('generate')">Generate</Button>
        <Button variant="secondary" :loading="busy === 'refine'" @click="run('refine')">Refine</Button>
        <Button variant="secondary" :loading="busy === 'save'" @click="save">Save draft</Button>
        <Button variant="secondary" :loading="busy === 'publish'" @click="publish">Publish design</Button>
      </div>
      <div class="flex flex-wrap gap-2">
        <button v-for="example in examples" :key="example" type="button" class="rounded-full border border-line bg-white px-3 py-1 text-left text-xs text-stone-600" @click="prompt = example">{{ example.slice(0, 42) }}…</button>
      </div>
      <div class="grid grid-cols-2 gap-2">
        <label class="text-xs text-stone-500">Theme
          <select :value="draft.theme.name" :class="controlClass" class="mt-1" @change="onTheme">
            <option v-for="name in ['editorial', 'minimal', 'mono', 'luxury', 'bold', 'playful', 'dark']" :key="name" :value="name">{{ name }}</option>
          </select>
        </label>
        <label class="text-xs text-stone-500">Images
          <select v-model="draft.theme.imageRatio" :class="controlClass" class="mt-1">
            <option value="square">square</option>
            <option value="portrait">portrait</option>
            <option value="landscape">landscape</option>
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
        <label v-for="key in ['background', 'surface', 'text', 'muted', 'primary', 'onPrimary', 'accent', 'border']" :key="key" class="text-[10px] uppercase text-stone-500">
          {{ key }}
          <input :value="colorValue(key)" type="color" class="mt-1 h-8 w-full" @input="onColor(key, $event)" />
        </label>
      </div>
      <ul class="space-y-2">
        <li v-for="(item, index) in draft.sections" :key="item.id" class="rounded-xl border bg-card p-3" :class="selected === index ? 'border-rust' : 'border-line'">
          <button class="flex w-full items-center justify-between text-left text-sm" type="button" @click="selected = index">
            <span class="capitalize">{{ item.type.replaceAll('_', ' ') }}</span>
            <span class="text-xs text-stone-500">{{ item.visible ? 'Shown' : 'Hidden' }}</span>
          </button>
          <div class="mt-2 flex gap-2 text-xs">
            <button type="button" @click="item.visible = !item.visible">{{ item.visible ? 'Hide' : 'Show' }}</button>
            <button type="button" @click="move(index, -1)">Up</button>
            <button type="button" @click="move(index, 1)">Down</button>
          </div>
        </li>
      </ul>
      <div v-if="section?.hero" class="space-y-2">
        <input v-model="section.hero.headline" :class="controlClass" />
        <textarea v-model="section.hero.subheadline" :class="controlClass" rows="2" />
        <input v-model="section.hero.ctaLabel" :class="controlClass" />
        <select v-model="section.hero.size" :class="controlClass"><option value="sm">small</option><option value="md">medium</option><option value="lg">large</option></select>
      </div>
      <div v-else-if="section?.productGrid" class="space-y-2">
        <input v-model="section.productGrid.title" :class="controlClass" />
        <select v-model.number="section.productGrid.columns" :class="controlClass"><option :value="2">2 columns</option><option :value="3">3 columns</option><option :value="4">4 columns</option></select>
      </div>
      <div v-else-if="section?.featured" class="space-y-2">
        <input v-model="section.featured.title" :class="controlClass" />
      </div>
      <div v-else-if="section?.about" class="space-y-2">
        <input v-model="section.about.title" :class="controlClass" />
        <textarea v-model="section.about.body" :class="controlClass" rows="4" />
      </div>
      <div v-else-if="section?.announcement"><input v-model="section.announcement.text" :class="controlClass" /></div>
      <div v-else-if="section?.categories"><input v-model="section.categories.title" :class="controlClass" /></div>
      <div v-else-if="section?.newsletter" class="space-y-2">
        <input v-model="section.newsletter.title" :class="controlClass" />
        <textarea v-model="section.newsletter.text" :class="controlClass" rows="2" />
      </div>
    </div>
    <div class="xl:sticky xl:top-4 xl:h-[calc(100vh-2rem)] xl:overflow-auto xl:rounded-2xl xl:border xl:border-line">
      <StorefrontView :store="previewStore" :products="products" :paths="storePaths(shop.shop.slug)" :interactive="false" />
    </div>
  </div>
</template>
