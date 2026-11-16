import { request } from '@/services/api'
import type { Product } from '@/types'

export type ProductInput = {
  name: string
  description: string
  priceCents: number
  images: string[]
  category: string
  sku: string
  stock: number
  published: boolean
}

export function listProducts(params: { q?: string; status?: string; category?: string; lowStock?: boolean } = {}) {
  const query = new URLSearchParams()
  if (params.q) query.set('q', params.q)
  if (params.status) query.set('status', params.status)
  if (params.category) query.set('category', params.category)
  if (params.lowStock) query.set('lowStock', '1')
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

export function setStock(id: string, stock: number) {
  return request<Product>(`/api/products/${id}/stock`, { method: 'POST', body: JSON.stringify({ stock }) })
}
