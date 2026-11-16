import type { ComputedRef, InjectionKey, Ref } from 'vue'
import type { useCart } from '@/features/storefront/cart'
import type { storePaths } from '@/features/stores/host'
import type { Product, PublicStore } from '@/types'

export type StorefrontContext = {
  slug: ComputedRef<string>
  store: Ref<PublicStore | null>
  products: Ref<Product[]>
  paths: ComputedRef<ReturnType<typeof storePaths>>
  cart: ComputedRef<ReturnType<typeof useCart>>
}

export const storefrontKey: InjectionKey<StorefrontContext> = Symbol('storefront')
