<script setup lang="ts">
import { computed, nextTick, onMounted, ref, toRaw } from 'vue'
import Button from '@/components/ui/Button.vue'
import Icon from '@/components/ui/Icon.vue'
import StorefrontView from '@/components/storefront/StorefrontView.vue'
import { controlClass, textareaClass } from '@/components/ui/control'
import { listCollections } from '@/features/collections/api'
import { listProducts } from '@/features/products/api'
import { generateDesign, publishDesign, refineDesign, saveDesign } from '@/features/storefront-builder/api'
import { publicUrl, storePaths } from '@/features/stores/host'
import { useShopStore } from '@/stores/shop'
import type { Collection, DesignNote, Product, Section, StorefrontConfig } from '@/types'

const shop = useShopStore()
const products = ref<Product[]>([])
const collections = ref<Collection[]>([])
const draft = ref<StorefrontConfig | null>(null)
const prompt = ref('')
const note = ref('')
const messages = ref<{ role: 'user' | 'designer'; text: string }[]>([])
const thread = ref<HTMLElement | null>(null)
const split = ref<HTMLElement | null>(null)
const chatWidth = ref(readChatWidth())
const dragging = ref(false)
const error = ref('')
const busy = ref('')
const selected = ref(0)
const panel = ref<'look' | 'page' | 'sections' | null>(null)
const frame = ref<'full' | 'phone'>('full')
const tools = [
  { id: 'look', label: 'Look', icon: 'palette' },
  { id: 'page', label: 'Page', icon: 'page' },
  { id: 'sections', label: 'Sections', icon: 'layers' },
] as const
const iconButton =
  'inline-flex h-8 w-8 shrink-0 items-center justify-center rounded border border-border bg-surface text-text hover:bg-bg focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50'
const iconButtonOn =
  'inline-flex h-8 w-8 shrink-0 items-center justify-center rounded border border-accent bg-accent text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent'
const iconGhost =
  'inline-flex h-7 w-7 shrink-0 items-center justify-center rounded text-muted hover:bg-bg hover:text-text focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-40'
const examples = [
  { label: 'Black and white', prompt: 'A quiet clothing shop. Black and white, large product photos, a short hero, and clear navigation.' },
  { label: 'Full-bleed hero', prompt: 'Use a luxury theme, a full-bleed hero, and more space.' },
  { label: 'Product list', prompt: 'Show the products as a list, use outline buttons, and hide the collections.' },
  { label: 'Darker', prompt: 'Make the background darker and keep the rest of the layout.' },
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
  editorial: { fontHeading: 'Fraunces', fontBody: 'Figtree', radius: 'md', imageRatio: 'portrait', density: 'roomy', buttons: 'fill', cards: 'border', collections: 'strip', heroLayout: 'split', productDisplay: 'feature', colors: { background: '#f6f1e8', surface: '#fffdf8', text: '#1f1a17', muted: '#6f655d', primary: '#9a3412', onPrimary: '#fff7ed', accent: '#b45309', border: '#e6ddd1' } },
  minimal: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'sm', imageRatio: 'square', density: 'tight', buttons: 'link', cards: 'flat', collections: 'strip', heroLayout: 'stack', productDisplay: 'grid', colors: { background: '#ffffff', surface: '#fafafa', text: '#171717', muted: '#737373', primary: '#171717', onPrimary: '#ffffff', accent: '#404040', border: '#e5e5e5' } },
  mono: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'sm', imageRatio: 'square', density: 'tight', buttons: 'outline', cards: 'flat', collections: 'hidden', heroLayout: 'stack', productDisplay: 'list', colors: { background: '#ffffff', surface: '#f4f4f5', text: '#111111', muted: '#52525b', primary: '#111111', onPrimary: '#ffffff', accent: '#18181b', border: '#e4e4e7' } },
  luxury: { fontHeading: 'Instrument Serif', fontBody: 'Figtree', radius: 'none', imageRatio: 'portrait', density: 'roomy', buttons: 'outline', cards: 'flat', collections: 'feature', heroLayout: 'bleed', productDisplay: 'feature', colors: { background: '#100e0c', surface: '#1a1714', text: '#f6f1e7', muted: '#b7ab9a', primary: '#c6a15b', onPrimary: '#1a1408', accent: '#e7d3a1', border: '#3a342c' } },
  bold: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'lg', imageRatio: 'square', density: 'regular', buttons: 'fill', cards: 'shadow', collections: 'strip', heroLayout: 'bleed', productDisplay: 'grid', colors: { background: '#fff7ed', surface: '#ffffff', text: '#1c1917', muted: '#78716c', primary: '#ea580c', onPrimary: '#ffffff', accent: '#0f172a', border: '#fed7aa' } },
  playful: { fontHeading: 'Fraunces', fontBody: 'Outfit', radius: 'lg', imageRatio: 'square', density: 'regular', buttons: 'fill', cards: 'shadow', collections: 'strip', heroLayout: 'stack', productDisplay: 'grid', colors: { background: '#fdf2f8', surface: '#ffffff', text: '#3b0764', muted: '#7e22ce', primary: '#db2777', onPrimary: '#ffffff', accent: '#0891b2', border: '#fbcfe8' } },
  dark: { fontHeading: 'Outfit', fontBody: 'Figtree', radius: 'md', imageRatio: 'portrait', density: 'roomy', buttons: 'fill', cards: 'border', collections: 'feature', heroLayout: 'bleed', productDisplay: 'feature', colors: { background: '#121212', surface: '#1c1c1c', text: '#f5f5f5', muted: '#a3a3a3', primary: '#f5f5f5', onPrimary: '#111111', accent: '#d4d4d4', border: '#2e2e2e' } },
}

function readChatWidth() {
  try {
    const width = Number(localStorage.getItem('sable-chat-width'))
    if (width >= 280 && width <= 560) return width
  } catch {
    return 380
  }
  return 380
}

function ensureLayout(config: StorefrontConfig) {
  const theme = config.theme
  theme.density = theme.density || 'regular'
  theme.buttons = theme.buttons || 'fill'
  theme.cards = theme.cards || 'flat'
  theme.collections = theme.collections || 'strip'
  theme.heroLayout = theme.heroLayout || 'split'
  theme.productDisplay = theme.productDisplay || 'grid'
  return config
}

function cloneConfig(config: StorefrontConfig): StorefrontConfig {
  return ensureLayout(structuredClone(toRaw(config)))
}

onMounted(async () => {
  try {
    const current = await shop.load()
    if (!current) return
    draft.value = cloneConfig(current.draft)
    note.value = current.designNote?.summary || ''
    messages.value = remembered(current.designNote)
    const [catalog, groups] = await Promise.all([listProducts({ status: 'published' }), listCollections()])
    products.value = catalog.items
    collections.value = groups.items.filter((item) => item.published)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Could not open the builder'
  }
})

async function scrollThread() {
  await nextTick()
  if (thread.value) thread.value.scrollTop = thread.value.scrollHeight
}

function remembered(note?: DesignNote) {
  const stored = note?.turns?.filter((turn) => turn.text && (turn.role === 'user' || turn.role === 'designer'))
  if (stored?.length) return stored.map((turn) => ({ role: turn.role === 'user' ? 'user' as const : 'designer' as const, text: turn.text }))
  const turns: { role: 'user' | 'designer'; text: string }[] = []
  if (note?.prompt) turns.push({ role: 'user', text: note.prompt })
  if (note?.summary) turns.push({ role: 'designer', text: note.summary })
  if (!turns.length) {
    turns.push({ role: 'designer', text: 'Describe the shop, or a change. I can restyle the hero, the product layout, the buttons, the collections, and the colors.' })
  }
  return turns
}

async function send(kind: 'generate' | 'refine') {
  const text = prompt.value.trim()
  if (!text || busy.value || !draft.value) return
  messages.value.push({ role: 'user', text })
  await scrollThread()
  const ok = await run(kind)
  if (ok) {
    messages.value = remembered(shop.shop?.designNote)
    prompt.value = ''
  } else if (error.value) {
    messages.value.push({ role: 'designer', text: error.value })
  }
  await scrollThread()
}

function onPromptKey(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    void send('refine')
  }
}

function startDrag(event: PointerEvent) {
  dragging.value = true
  const handle = event.currentTarget as HTMLElement
  try {
    handle.setPointerCapture(event.pointerId)
  } catch {
    /* the pointer can end before capture is available */
  }
}

function onDrag(event: PointerEvent) {
  if (!dragging.value || !split.value) return
  const bounds = split.value.getBoundingClientRect()
  const max = Math.min(560, bounds.width - 420)
  chatWidth.value = Math.round(Math.min(Math.max(event.clientX - bounds.left, 280), Math.max(280, max)))
}

function rememberWidth() {
  try {
    localStorage.setItem('sable-chat-width', String(chatWidth.value))
  } catch {
    /* preference is optional */
  }
}

function stopDrag() {
  if (!dragging.value) return
  dragging.value = false
  rememberWidth()
}

function nudgeChat(event: KeyboardEvent) {
  if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight') return
  event.preventDefault()
  chatWidth.value = Math.min(560, Math.max(280, chatWidth.value + (event.key === 'ArrowLeft' ? -16 : 16)))
  rememberWidth()
}

async function run(kind: 'generate' | 'refine') {
  if (!draft.value) return false
  busy.value = kind
  error.value = ''
  try {
    const result = kind === 'generate' ? await generateDesign(prompt.value) : await refineDesign(prompt.value, draft.value)
    shop.set(result.store)
    draft.value = cloneConfig(result.store.draft)
    note.value = result.summary
    return true
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'The designer failed'
    return false
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
  <h1 class="sr-only">Builder</h1>
  <p v-if="error && !draft" class="m-4 rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
  <div v-if="draft && previewStore && shop.shop" ref="split" class="w-full min-w-0 lg:flex lg:h-full lg:overflow-hidden" :class="dragging ? 'select-none' : ''">
    <div class="flex flex-col border-b border-border bg-surface lg:h-full lg:w-[var(--chat-width)] lg:shrink-0 lg:border-b-0 lg:border-r lg:border-border" :style="{ '--chat-width': chatWidth + 'px' }">
      <div class="sticky top-12 z-20 flex shrink-0 items-center gap-1 border-b border-border bg-surface px-2 py-1.5 md:top-0 lg:static">
        <button
          v-for="item in tools"
          :key="item.id"
          type="button"
          :class="panel === item.id ? iconButtonOn : iconButton"
          :aria-label="item.label"
          :aria-pressed="panel === item.id"
          :title="item.label"
          @click="panel = panel === item.id ? null : item.id"
        >
          <Icon :name="item.icon" />
        </button>
      </div>

      <div v-if="!panel" ref="thread" class="flex min-h-[280px] flex-1 flex-col gap-5 overflow-y-auto px-4 py-4 lg:min-h-0">
        <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>
        <div v-for="(message, index) in messages" :key="index" class="flex" :class="message.role === 'user' ? 'justify-end' : 'justify-start'">
          <div v-if="message.role === 'designer'" class="max-w-[92%]">
            <p class="mb-1 text-xs font-medium text-muted">Sable</p>
            <p class="text-sm leading-relaxed text-text">{{ message.text }}</p>
          </div>
          <p v-else class="max-w-[85%] rounded bg-bg px-3 py-2 text-sm leading-relaxed text-text">{{ message.text }}</p>
        </div>
        <p v-if="busy === 'generate' || busy === 'refine'" class="text-xs text-muted">Designing…</p>
      </div>
      <div v-else class="space-y-3 px-3 py-3 lg:min-h-0 lg:flex-1 lg:overflow-y-auto">
        <p v-if="error" class="rounded bg-bad-bg px-3 py-2 text-sm text-bad-fg">{{ error }}</p>

        <div v-if="panel === 'look'" class="space-y-3">
          <div class="flex flex-wrap gap-1.5">
            <button v-for="name in Object.keys(themePresets)" :key="name" type="button" class="h-7 rounded border px-2 text-xs capitalize focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent" :class="draft.theme.name === name ? 'border-accent bg-accent text-white' : 'border-border bg-surface text-text hover:bg-bg'" @click="pickTheme(name)">{{ name }}</button>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <label class="text-xs text-muted">Images
              <select v-model="draft.theme.imageRatio" :class="controlClass" class="mt-1">
                <option value="square">Square</option>
                <option value="portrait">Portrait</option>
                <option value="landscape">Landscape</option>
              </select>
            </label>
            <label class="text-xs text-muted">Corners
              <select v-model="draft.theme.radius" :class="controlClass" class="mt-1">
                <option value="none">Sharp</option>
                <option value="sm">Small</option>
                <option value="md">Medium</option>
                <option value="lg">Large</option>
              </select>
            </label>
            <label class="text-xs text-muted">Heading
              <select v-model="draft.theme.fontHeading" :class="controlClass" class="mt-1"><option v-for="font in fonts" :key="font">{{ font }}</option></select>
            </label>
            <label class="text-xs text-muted">Body
              <select v-model="draft.theme.fontBody" :class="controlClass" class="mt-1"><option v-for="font in fonts" :key="font">{{ font }}</option></select>
            </label>
          </div>
          <div>
            <p class="text-xs text-muted">Colors</p>
            <div class="mt-1 flex flex-wrap gap-1.5">
              <label v-for="key in colorKeys" :key="key" class="relative h-8 w-8" :title="key">
                <span class="sr-only">{{ key }}</span>
                <input :value="colorValue(key)" type="color" class="h-8 w-8 cursor-pointer rounded border border-border bg-surface" @input="onColor(key, $event)" />
              </label>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <label class="text-xs text-muted">Space
              <select v-model="draft.theme.density" :class="controlClass" class="mt-1"><option value="tight">Tight</option><option value="regular">Regular</option><option value="roomy">Roomy</option></select>
            </label>
            <label class="text-xs text-muted">Buttons
              <select v-model="draft.theme.buttons" :class="controlClass" class="mt-1"><option value="fill">Fill</option><option value="outline">Outline</option><option value="link">Text</option></select>
            </label>
            <label class="text-xs text-muted">Cards
              <select v-model="draft.theme.cards" :class="controlClass" class="mt-1"><option value="flat">Flat</option><option value="border">Border</option><option value="shadow">Shadow</option></select>
            </label>
            <label class="text-xs text-muted">Collections
              <select v-model="draft.theme.collections" :class="controlClass" class="mt-1"><option value="strip">Strip</option><option value="feature">Featured</option><option value="hidden">Hidden</option></select>
            </label>
            <label class="text-xs text-muted">Hero
              <select v-model="draft.theme.heroLayout" :class="controlClass" class="mt-1"><option value="stack">Stacked</option><option value="split">Split</option><option value="bleed">Full bleed</option></select>
            </label>
            <label class="text-xs text-muted">Products
              <select v-model="draft.theme.productDisplay" :class="controlClass" class="mt-1"><option value="grid">Grid</option><option value="feature">Feature first</option><option value="list">List</option></select>
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
          <label v-for="item in draft.header.nav" :key="item.target" class="block text-xs capitalize text-muted">{{ item.target }}
            <input v-model="item.label" :class="controlClass" class="mt-1" />
          </label>
          <label class="block text-xs text-muted">Footer
            <input v-model="draft.footer.text" :class="controlClass" class="mt-1" />
          </label>
        </div>

        <div v-else-if="panel === 'sections'" class="space-y-3">
          <ul class="space-y-1.5">
            <li v-for="(item, index) in draft.sections" :key="item.id" class="rounded border bg-surface px-2 py-1.5" :class="selected === index ? 'border-text' : 'border-border'">
              <div class="flex items-center gap-1">
                <button class="min-w-0 flex-1 truncate text-left text-sm capitalize" type="button" @click="selected = index">{{ item.type.replaceAll('_', ' ') }}</button>
                <button type="button" :class="iconGhost" :aria-label="item.visible ? 'Hide section' : 'Show section'" :title="item.visible ? 'Hide' : 'Show'" @click="item.visible = !item.visible">
                  <Icon :name="item.visible ? 'eye' : 'eye-off'" />
                </button>
                <button type="button" :class="iconGhost" aria-label="Move up" title="Move up" :disabled="index === 0" @click="move(index, -1)"><Icon name="chevron-up" /></button>
                <button type="button" :class="iconGhost" aria-label="Move down" title="Move down" :disabled="index === draft.sections.length - 1" @click="move(index, 1)"><Icon name="chevron-down" /></button>
              </div>
            </li>
          </ul>
          <div v-if="section?.hero" class="space-y-2">
            <input v-model="section.hero.eyebrow" :class="controlClass" placeholder="Eyebrow" />
            <input v-model="section.hero.headline" :class="controlClass" placeholder="Headline" />
            <textarea v-model="section.hero.subheadline" :class="textareaClass" rows="2" placeholder="Subheadline" />
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
            <label class="block text-xs text-muted">Products shown
              <input v-model.number="section.productGrid.limit" :class="controlClass" class="mt-1" type="number" min="1" max="24" />
            </label>
          </div>
          <div v-else-if="section?.featured" class="space-y-2">
            <input v-model="section.featured.title" :class="controlClass" placeholder="Title" />
            <label class="block text-xs text-muted">Products shown
              <input v-model.number="section.featured.limit" :class="controlClass" class="mt-1" type="number" min="1" max="8" />
            </label>
          </div>
          <div v-else-if="section?.about" class="space-y-2">
            <input v-model="section.about.title" :class="controlClass" placeholder="Title" />
            <textarea v-model="section.about.body" :class="textareaClass" rows="4" placeholder="About the shop" />
            <input v-model="section.about.imageUrl" :class="controlClass" placeholder="Image URL" />
          </div>
          <div v-else-if="section?.announcement"><input v-model="section.announcement.text" :class="controlClass" placeholder="Announcement" /></div>
          <div v-else-if="section?.categories"><input v-model="section.categories.title" :class="controlClass" placeholder="Title" /></div>
          <div v-else-if="section?.newsletter" class="space-y-2">
            <input v-model="section.newsletter.title" :class="controlClass" placeholder="Title" />
            <textarea v-model="section.newsletter.text" :class="textareaClass" rows="2" />
            <p class="text-xs text-muted">Copy only. This does not collect email addresses.</p>
          </div>
        </div>
      </div>

      <div class="shrink-0 border-t border-border bg-surface px-3 py-3">
        <p v-if="shop.shop.published" class="mb-2 text-xs text-muted">Sending applies this design to the live shop.</p>
        <div v-if="messages.length < 3" class="mb-2 flex flex-wrap gap-x-3 gap-y-1">
          <button v-for="example in examples" :key="example.label" type="button" class="text-xs text-muted hover:text-text focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent" @click="prompt = example.prompt">{{ example.label }}</button>
        </div>
        <div class="rounded border border-border bg-bg focus-within:border-accent">
          <label class="sr-only" for="builder-prompt">Message</label>
          <textarea id="builder-prompt" v-model="prompt" rows="2" class="w-full resize-none bg-transparent px-3 py-2 text-sm text-text outline-none placeholder:text-muted" placeholder="Describe a change" @keydown="onPromptKey" />
          <div class="flex items-center justify-between gap-2 px-2 pb-2">
            <button type="button" class="text-xs text-muted hover:text-text disabled:opacity-50" :disabled="busy !== '' || prompt.trim() === ''" @click="send('generate')">New design</button>
            <button type="button" class="inline-flex h-8 items-center gap-1.5 rounded bg-accent px-2.5 text-sm font-medium text-white hover:bg-accent-hover focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:opacity-50" :disabled="busy !== '' || prompt.trim() === ''" @click="send('refine')">
              <span v-if="busy === 'refine' || busy === 'generate'" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-r-transparent"></span>
              <Icon v-else name="send" />
              Send
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      class="hidden w-1.5 shrink-0 cursor-col-resize touch-none items-center justify-center bg-surface hover:bg-border lg:flex"
      :class="dragging ? 'bg-accent/40' : ''"
      role="separator"
      aria-orientation="vertical"
      aria-label="Resize chat and preview"
      aria-valuemin="280"
      aria-valuemax="560"
      :aria-valuenow="chatWidth"
      tabindex="0"
      @pointerdown="startDrag"
      @pointermove="onDrag"
      @pointerup="stopDrag"
      @pointercancel="stopDrag"
      @keydown="nudgeChat"
    >
      <span class="h-8 w-px bg-border"></span>
    </div>

    <div class="flex min-h-[70vh] min-w-0 flex-1 flex-col bg-bg lg:h-full lg:min-h-0">
      <div class="flex min-h-10 shrink-0 flex-wrap items-center justify-between gap-2 border-b border-border bg-surface px-3 py-1">
        <p class="text-xs text-muted">Draft preview</p>
        <div class="flex flex-wrap items-center justify-end gap-1">
          <button type="button" :class="frame === 'full' ? iconButtonOn : iconButton" aria-label="Desktop preview" title="Desktop" @click="frame = 'full'"><Icon name="monitor" /></button>
          <button type="button" :class="frame === 'phone' ? iconButtonOn : iconButton" aria-label="Phone preview" title="Phone" @click="frame = 'phone'"><Icon name="phone" /></button>
          <span class="mx-1 h-4 w-px bg-border"></span>
          <button type="button" :class="iconButton" aria-label="Save" title="Save" :disabled="busy !== ''" @click="save">
            <span v-if="busy === 'save'" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-r-transparent"></span>
            <Icon v-else name="save" />
          </button>
          <Button :loading="busy === 'publish'" :disabled="busy !== '' && busy !== 'publish'" @click="publish">Publish</Button>
          <a v-if="shopUrl" :class="iconButton" :href="shopUrl.subdomain" target="_blank" rel="noreferrer" aria-label="Open shop" title="Open shop"><Icon name="external" /></a>
          <a v-if="shopUrl" :class="iconButton" :href="shopUrl.path" target="_blank" rel="noreferrer" aria-label="Same-host preview" title="Same-host preview"><Icon name="eye" /></a>
        </div>
      </div>
      <div class="min-h-0 flex-1 overflow-auto">
        <div :class="frame === 'phone' ? 'mx-auto min-h-full max-w-[390px] border-x border-border' : ''">
          <StorefrontView :store="previewStore" :products="products" :collections="collections" :paths="storePaths(shop.shop.slug)" :interactive="false" :compact="frame === 'phone'" />
        </div>
      </div>
    </div>
  </div>
</template>
