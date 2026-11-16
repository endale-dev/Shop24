import { request } from '@/services/api'
import type { Store } from '@/types'

export function createStore(body: { name: string; slug: string; description: string }) {
  return request<Store>('/api/stores', { method: 'POST', body: JSON.stringify(body) })
}

export function myStore() {
  return request<Store>('/api/stores/mine')
}

export function updateStore(body: { name: string; slug: string; description: string; logoUrl: string; currency: string }) {
  return request<Store>('/api/stores/mine', { method: 'PATCH', body: JSON.stringify(body) })
}

export function publishStore() {
  return request<Store>('/api/stores/mine/publish', { method: 'POST' })
}

export function unpublishStore() {
  return request<Store>('/api/stores/mine/unpublish', { method: 'POST' })
}
