import { request } from '@/services/api'
import type { Customer, Order } from '@/types'

export function listOrders(status = '') {
  const suffix = status ? `?status=${encodeURIComponent(status)}` : ''
  return request<{ items: Order[] }>(`/api/orders${suffix}`)
}

export function getOrder(id: string) {
  return request<Order>(`/api/orders/${id}`)
}

export function updateOrderStatus(id: string, status: string) {
  return request<Order>(`/api/orders/${id}`, { method: 'PATCH', body: JSON.stringify({ status }) })
}

export function placeOrder(slug: string, body: { customer: Customer; items: { productId: string; variantId: string; quantity: number; unitPriceCents: number }[] }) {
  return request<Order>(`/api/public/stores/${slug}/orders`, { method: 'POST', body: JSON.stringify(body) })
}
