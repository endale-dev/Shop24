import { computed, ref, watch, type Ref } from 'vue'

export type CartLine = {
  productId: string
  name: string
  priceCents: number
  image: string
  quantity: number
  stock: number
}

const carts = new Map<string, Ref<CartLine[]>>()

function storageKey(slug: string) {
  return `sable.cart.${slug}`
}

function read(slug: string): CartLine[] {
  try {
    const parsed = JSON.parse(localStorage.getItem(storageKey(slug)) || '[]') as CartLine[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export function useCart(slug: string) {
  let lines = carts.get(slug)
  if (!lines) {
    lines = ref(read(slug))
    watch(lines, (value) => localStorage.setItem(storageKey(slug), JSON.stringify(value)), { deep: true })
    carts.set(slug, lines)
  }
  const count = computed(() => lines!.value.reduce((sum, line) => sum + line.quantity, 0))
  const subtotal = computed(() => lines!.value.reduce((sum, line) => sum + line.priceCents * line.quantity, 0))

  function add(line: Omit<CartLine, 'quantity'>, quantity = 1) {
    const existing = lines!.value.find((item) => item.productId === line.productId)
    if (existing) {
      existing.quantity = Math.min(existing.stock, existing.quantity + quantity)
      existing.priceCents = line.priceCents
      existing.stock = line.stock
      return
    }
    lines!.value.push({ ...line, quantity: Math.min(line.stock, quantity) })
  }

  function setQuantity(productId: string, quantity: number) {
    const line = lines!.value.find((item) => item.productId === productId)
    if (!line) return
    line.quantity = Math.max(1, Math.min(line.stock, quantity))
  }

  function remove(productId: string) {
    lines!.value = lines!.value.filter((item) => item.productId !== productId)
  }

  function clear() {
    lines!.value = []
  }

  return { lines, count, subtotal, add, setQuantity, remove, clear }
}
