import { defineStore } from 'pinia'
import { ref } from 'vue'
import { currentUser, login as loginRequest, logout as logoutRequest, register as registerRequest } from '@/features/auth/api'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const ready = ref(false)

  async function bootstrap() {
    try {
      user.value = await currentUser()
    } catch {
      user.value = null
    } finally {
      ready.value = true
    }
  }

  async function login(email: string, password: string) {
    user.value = await loginRequest({ email, password })
  }

  async function register(name: string, email: string, password: string) {
    user.value = await registerRequest({ name, email, password })
  }

  async function logout() {
    await logoutRequest()
    user.value = null
  }

  return { user, ready, bootstrap, login, register, logout }
})
