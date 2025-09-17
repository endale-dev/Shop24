import api, { getShopUrl, getShopSlugFromSubdomain } from './api'

export interface Shop {
  id: string
  name: string
  slug: string
  image: string
  banner?: string
  description: string
}

export function setCurrentShopSlug(slug: string) {
  localStorage.setItem('currentShopSlug', slug);
}

export function getPersistedShopSlug(): string | null {
  return localStorage.getItem('currentShopSlug');
}

export function setCurrentShopData(shop: Shop) {
  localStorage.setItem('currentShopData', JSON.stringify(shop));
}

export function getPersistedShopData(): Shop | null {
  const data = localStorage.getItem('currentShopData');
  if (!data) return null;
  try {
    return JSON.parse(data);
  } catch {
    return null;
  }
}

export function getCurrentShopSlug(): string | null {
  // 1) Prefer subdomain when available
  const fromSubdomain = getShopSlugFromSubdomain();
  if (fromSubdomain) return fromSubdomain;

  // 2) Check query params for explicit shop slug
  try {
    const qp = new URLSearchParams(window.location.search)
    const shopParam = qp.get('shop') || qp.get('shopSlug')
    if (shopParam) {
      setCurrentShopSlug(shopParam)
      return shopParam
    }
  } catch {}

  // 3) Fallback to persisted value
  const persisted = getPersistedShopSlug();
  if (persisted) return persisted;

  // 4) Dev convenience: allow an env-provided default when running without subdomain
  const envDefault = (import.meta as any).env?.VITE_DEFAULT_SHOP_SLUG as string | undefined
  if (envDefault) return envDefault;

  return null;
}

// Resolve a shop slug -> shop object (including id) using public storefront API
export async function resolveShopBySlug(slug: string): Promise<Shop | null> {
  if (!slug) return null
  try {
    const response = await api.get(getShopUrl(slug, ''))
    return response.data as Shop
  } catch (err) {
    console.error('Failed to resolve shop by slug', err)
    return null
  }
}

// Convenience: get a storeId (ObjectID hex) for the current context. Will try query param, window.__STORE_ID__, or resolve via shop slug.
export async function getCurrentStoreId(): Promise<string | null> {
  // 1) explicit query param
  try {
    const qp = new URLSearchParams(window.location.search)
    const storeId = qp.get('storeId')
    if (storeId) return storeId
  } catch {}

  // 2) injected global
  const globalId = (window as any).__STORE_ID__
  if (globalId) return String(globalId)

  // 3) resolve from slug
  const slug = getCurrentShopSlug()
  if (!slug) return null
  const shop = await resolveShopBySlug(slug)
  return shop?.id ?? null
}

export async function fetchShop(shopSlug?: string): Promise<Shop | null> {
  const slug = shopSlug || getShopSlugFromSubdomain() || getCurrentShopSlug();
  if (!slug) return null;
  try {
    const response = await api.get(getShopUrl(slug, ''))
    return response.data
  } catch (error) {
    console.error('Failed to fetch shop:', error)
    return null
  }
}

export function clearPersistedShop() {
  localStorage.removeItem('currentShopSlug');
  localStorage.removeItem('currentShopData');
}
