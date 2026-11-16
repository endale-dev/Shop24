import { request } from '@/services/api'
import type { Collection, Product, PublicStore } from '@/types'

export type CollectionInput = {
  name: string
  slug: string
  description: string
  imageUrl: string
  productIds: string[]
  published: boolean
}

export function listCollections() {
  return request<{ items: Collection[] }>('/api/collections')
}

export function getCollection(id: string) {
  return request<Collection>(`/api/collections/${id}`)
}

export function createCollection(body: CollectionInput) {
  return request<Collection>('/api/collections', { method: 'POST', body: JSON.stringify(body) })
}

export function updateCollection(id: string, body: CollectionInput) {
  return request<Collection>(`/api/collections/${id}`, { method: 'PATCH', body: JSON.stringify(body) })
}

export function deleteCollection(id: string) {
  return request<void>(`/api/collections/${id}`, { method: 'DELETE' })
}

export function publicCollections(slug: string) {
  return request<{ items: Collection[] }>(`/api/public/stores/${slug}/collections`)
}

export function publicCollection(slug: string, collectionSlug: string) {
  return request<{ collection: Collection; products: Product[]; store: PublicStore }>(`/api/public/stores/${slug}/collections/${collectionSlug}`)
}
