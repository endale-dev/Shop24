const slugPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/

export function slugFromHost(hostname: string, baseDomain = 'localhost') {
  const host = hostname.toLowerCase().replace(/\.$/, '')
  const base = baseDomain.toLowerCase().replace(/\.$/, '')
  if (!host || !base || host === base || host === `www.${base}`) return null
  const suffix = `.${base}`
  if (!host.endsWith(suffix)) return null
  const slug = host.slice(0, -suffix.length)
  if (!slugPattern.test(slug) || slug.length < 3 || slug.length > 40) return null
  return slug
}

export function baseDomain() {
  return import.meta.env.VITE_BASE_DOMAIN || 'localhost'
}

export const hostSlug = slugFromHost(window.location.hostname, baseDomain())

export function storePaths(slug: string) {
  const onHost = hostSlug === slug
  const root = onHost ? '' : `/s/${slug}`
  return {
    home: root || '/',
    cart: `${root}/cart`,
    checkout: `${root}/checkout`,
    product: (id: string) => `${root}/p/${id}`,
    collection: (collectionSlug: string) => `${root}/c/${collectionSlug}`,
  }
}

function hostWithPort(host: string) {
  const port = window.location.port
  return port ? `${host}:${port}` : host
}

export function publicUrl(slug: string) {
  const protocol = window.location.protocol
  return {
    path: `${window.location.origin}/s/${slug}`,
    subdomain: `${protocol}//${hostWithPort(`${slug}.${baseDomain()}`)}`,
  }
}
