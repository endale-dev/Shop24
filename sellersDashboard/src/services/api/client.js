/**
 * Enhanced API Client
 * Centralized axios instance with interceptors, error handling, and logging
 */

import axios from 'axios'
import { transformError } from './errors'

// Determine API base dynamically
function resolveApiBase() {
  const envBase = import.meta?.env?.VITE_API_BASE
  if (envBase) return envBase
  
  const host = window.location.hostname
  if (host === 'localhost' || host === '127.0.0.1') {
    return `http://${host}:8080`
  }
  return 'http://api.shop24.sbs'
}

// Create axios instance
const client = axios.create({
  baseURL: resolveApiBase(),
  withCredentials: true,
  timeout: 30000, // 30 seconds
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request logging (development only)
const isDev = import.meta.env.MODE === 'development'

// Request interceptor
client.interceptors.request.use(
  (config) => {
    if (isDev) {
      console.log('[API Request]', {
        method: config.method.toUpperCase(),
        url: config.url,
        data: config.data,
        params: config.params
      })
    }
    return config
  },
  (error) => {
    if (isDev) {
      console.error('[API Request Error]', error)
    }
    return Promise.reject(transformError(error))
  }
)

// Token refresh queue
let isRefreshing = false
let failedQueue = []
let refreshAttempts = 0
const MAX_REFRESH_ATTEMPTS = 3

function processQueue(err, token = null) {
  failedQueue.forEach(p => err ? p.reject(err) : p.resolve(token))
  failedQueue = []
}

const refreshablePaths = ['/auth/seller/me', '/seller/', '/shop/']
const nonRefreshablePaths = ['/auth/seller/login', '/auth/seller/register', '/auth/seller/refresh']

// Response interceptor
client.interceptors.response.use(
  (response) => {
    if (isDev) {
      console.log('[API Response]', {
        status: response.status,
        url: response.config.url,
        data: response.data
      })
    }
    
    // Reset refresh attempts on successful request
    if (refreshAttempts > 0) {
      refreshAttempts = 0
    }
    
    return response
  },
  async (error) => {
    const { config, response } = error
    const originalRequest = config

    if (isDev) {
      console.error('[API Response Error]', {
        status: response?.status,
        url: config?.url,
        message: response?.data?.error || error.message
      })
    }

    const is401 = response?.status === 401
    const isRefreshableRequest = refreshablePaths.some(path => originalRequest?.url?.includes(path))
    const isNonRefreshable = nonRefreshablePaths.some(path => originalRequest?.url?.includes(path))

    // Handle 401 errors with token refresh
    if (is401 && isRefreshableRequest && !isNonRefreshable && !originalRequest._retry) {
      originalRequest._retry = true

      // Queue request if already refreshing
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then(() => client(originalRequest))
          .catch(err => Promise.reject(transformError(err)))
      }

      isRefreshing = true

      try {
        // Check max refresh attempts
        if (refreshAttempts >= MAX_REFRESH_ATTEMPTS) {
          console.error('[API] Max refresh attempts exceeded')
          
          // Clear auth state
          const { useAuthStore } = await import('@/store/auth')
          const authStore = useAuthStore()
          authStore.user = null
          authStore.lastAuthCheck = null
          
          processQueue(new Error('Max refresh attempts exceeded'))
          return Promise.reject(transformError(error))
        }

        refreshAttempts++
        console.log(`[API] Refreshing token (attempt ${refreshAttempts}/${MAX_REFRESH_ATTEMPTS})`)

        // Refresh token
        await client.post('/auth/seller/refresh')

        // Reset attempts on success
        refreshAttempts = 0
        console.log('[API] Token refresh successful')

        processQueue(null)
        return client(originalRequest)

      } catch (refreshError) {
        console.error('[API] Token refresh failed:', refreshError)

        // Clear auth state
        try {
          const { useAuthStore } = await import('@/store/auth')
          const authStore = useAuthStore()
          authStore.user = null
          authStore.lastAuthCheck = null
        } catch (importError) {
          console.error('[API] Failed to import auth store:', importError)
        }

        processQueue(refreshError)
        return Promise.reject(transformError(refreshError))

      } finally {
        isRefreshing = false
      }
    }

    // Transform and reject error
    return Promise.reject(transformError(error))
  }
)

/**
 * Helper function to handle API calls with loading state
 * @param {Function} apiCall - The API call function
 * @param {Object} store - Optional Pinia store to update loading state
 * @returns {Promise}
 */
export async function withLoading(apiCall, store = null) {
  if (store && typeof store.setLoading === 'function') {
    store.setLoading(true)
  }

  try {
    const result = await apiCall()
    return result
  } finally {
    if (store && typeof store.setLoading === 'function') {
      store.setLoading(false)
    }
  }
}

/**
 * Cancel token for aborting requests
 */
export function createCancelToken() {
  return axios.CancelToken.source()
}

export default client
