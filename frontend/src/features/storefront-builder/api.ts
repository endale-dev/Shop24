import { request } from '@/services/api'
import type { DesignResult, PublicStore, Store, StorefrontConfig } from '@/types'

export function generateDesign(prompt: string) {
  return request<DesignResult>('/api/storefront/generate', { method: 'POST', body: JSON.stringify({ prompt }) })
}

export function refineDesign(prompt: string, config: StorefrontConfig) {
  return request<DesignResult>('/api/storefront/refine', { method: 'POST', body: JSON.stringify({ prompt, config }) })
}

export function saveDesign(config: StorefrontConfig) {
  return request<Store>('/api/storefront', { method: 'PUT', body: JSON.stringify({ config }) })
}

export function publishDesign() {
  return request<Store>('/api/storefront/publish', { method: 'POST' })
}

export function publicStore(slug: string) {
  return request<PublicStore>(`/api/public/stores/${slug}`)
}

export function publicProducts(slug: string) {
  return request<{ items: import('@/types').Product[] }>(`/api/public/stores/${slug}/products`)
}

export function publicProduct(slug: string, id: string) {
  return request<{ product: import('@/types').Product; store: PublicStore }>(`/api/public/stores/${slug}/products/${id}`)
}
