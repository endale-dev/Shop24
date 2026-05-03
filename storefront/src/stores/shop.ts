import { defineStore } from 'pinia'
import axios from 'axios'

interface ShopState {
    shopSlug: string
    shop: any | null
    theme: any | null
    currentPage: any | null
    isLoading: boolean
    error: string | null
}

export const useShopStore = defineStore('shop', {
    state: (): ShopState => ({
        shopSlug: '',
        shop: null,
        theme: null,
        currentPage: null,
        isLoading: true,
        error: null
    }),

    actions: {
        extractSubdomain() {
            const hostname = window.location.hostname
            // Extract subdomain (e.g., shupslug.localhost or shupslug.shop24.sbs)
            const parts = hostname.split('.')

            if (parts.length >= 2 && parts[0] !== 'www' && parts[0] !== 'localhost') {
                this.shopSlug = parts[0]
            } else {
                // Fallback for direct IP or missing subdomain
                this.shopSlug = 'demo' // Or redirect
            }
        },

        async initialize() {
            this.extractSubdomain()
            this.isLoading = true
            this.error = null

            try {
                const [shopRes, themeRes] = await Promise.all([
                    axios.get(`/shops/${this.shopSlug}`),
                    axios.get(`/shops/${this.shopSlug}/theme`)
                ])

                this.shop = shopRes.data
                this.theme = themeRes.data.theme || themeRes.data
            } catch (err: any) {
                this.error = err.response?.data?.error || 'Failed to load shop details'
                console.error('Failed to init shop:', err)
            } finally {
                this.isLoading = false
            }
        },

        async loadPage(path: string) {
            this.isLoading = true
            try {
                const res = await axios.get(`/shops/${this.shopSlug}/page`, {
                    params: { path }
                })
                this.currentPage = res.data.page || res.data
            } catch (err) {
                console.error('Failed to load page:', path, err)
                // We do not overwrite error to allow global shop load failure to persist
            } finally {
                this.isLoading = false
            }
        }
    }
})
