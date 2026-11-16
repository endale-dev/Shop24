import { request } from '@/services/api'
import type { Product, VariantOption } from '@/types'

export type ProductInput = {
  name: string
  description: string
  images: string[]
  category: string
  optionNames: string[]
  variants: {
    id: string
    options: VariantOption[]
    sku: string
    priceCents: number
    stock: number
    imageUrl: string
  }[]
  published: boolean
}

export function priceVaries(product: Product) {
  const prices = new Set((product.variants ?? []).map((variant) => variant.priceCents))
  return prices.size > 1
}

export function listProducts(params: { q?: string; status?: string } = {}) {
  const query = new URLSearchParams()
  if (params.q) query.set('q', params.q)
  if (params.status) query.set('status', params.status)
  const suffix = query.toString() ? `?${query}` : ''
  return request<{ items: Product[] }>(`/api/products${suffix}`)
}

export function getProduct(id: string) {
  return request<Product>(`/api/products/${id}`)
}

export function createProduct(body: ProductInput) {
  return request<Product>('/api/products', { method: 'POST', body: JSON.stringify(body) })
}

export function updateProduct(id: string, body: ProductInput) {
  return request<Product>(`/api/products/${id}`, { method: 'PATCH', body: JSON.stringify(body) })
}

export function deleteProduct(id: string) {
  return request<void>(`/api/products/${id}`, { method: 'DELETE' })
}

export function setStock(id: string, variantId: string, stock: number) {
  return request<Product>(`/api/products/${id}/stock`, { method: 'POST', body: JSON.stringify({ stock, variantId }) })
}
