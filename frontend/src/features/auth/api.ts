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

export function authProviders() {
  return request<{ google: boolean }>('/api/auth/providers')
}

export function requestReset(email: string) {
  return request<{ status: string }>('/api/auth/forgot', { method: 'POST', body: JSON.stringify({ email }) })
}

export function resetPassword(token: string, password: string) {
  return request<User>('/api/auth/reset', { method: 'POST', body: JSON.stringify({ token, password }) })
}

export function updateName(name: string) {
  return request<User>('/api/auth/profile', { method: 'PATCH', body: JSON.stringify({ name }) })
}

export function changePassword(body: { current: string; password: string }) {
  return request<User>('/api/auth/password', { method: 'PATCH', body: JSON.stringify(body) })
}
