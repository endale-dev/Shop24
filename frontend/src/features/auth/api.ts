import { request } from '@/services/api'
import type { User } from '@/types'

export function register(body: { name: string; email: string; password: string }) {
  return request<User>('/api/auth/register', { method: 'POST', body: JSON.stringify(body) })
}

export function login(body: { email: string; password: string }) {
  return request<User>('/api/auth/login', { method: 'POST', body: JSON.stringify(body) })
}

export function logout() {
  return request<void>('/api/auth/logout', { method: 'POST' })
}

export function currentUser() {
  return request<User>('/api/auth/me')
}
