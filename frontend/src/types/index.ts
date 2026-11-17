export type User = {
  id: string
  name: string
  email: string
  hasPassword: boolean
  googleLinked: boolean
  createdAt: string
  updatedAt: string
}

export type NavItem = { label: string; target: 'home' | 'products' | 'about' | 'cart' | string }

export type Colors = {
  background: string
  surface: string
  text: string
  muted: string
  primary: string
  onPrimary: string
  accent: string
  border: string
}

export type Theme = {
  name: string
  colors: Colors
  fontHeading: string
  fontBody: string
  radius: string
  imageRatio: string
  density?: string
  buttons?: string
  cards?: string
  collections?: string
  heroLayout?: string
  productDisplay?: string
}

export type HeroProps = {
  eyebrow: string
  headline: string
  subheadline: string
  ctaLabel: string
  ctaTarget: string
  imageUrl: string
  align: string
  size: string
}

export type Section = {
  id: string
  type: string
  visible: boolean
  hero?: HeroProps
  productGrid?: { title: string; columns: number; limit: number }
  featured?: { title: string; limit: number }
  about?: { title: string; body: string; imageUrl: string }
  announcement?: { text: string }
  categories?: { title: string }
  newsletter?: { title: string; text: string }
}

export type StorefrontConfig = {
  version: number
  theme: Theme
  header: { announcement: string; showAnnouncement: boolean; alignment: string; nav: NavItem[] }
  sections: Section[]
  footer: { text: string; links: NavItem[] }
}

export type DesignTurn = { role: 'user' | 'designer' | string; text: string }

export type DesignNote = { provider: string; mode: string; prompt: string; summary: string; turns?: DesignTurn[] }

export type Store = {
  id: string
  ownerId: string
  name: string
  slug: string
  description: string
  logoUrl: string
  currency: string
  published: boolean
  storefront: StorefrontConfig
  draft: StorefrontConfig
  designNote?: DesignNote
  createdAt: string
  updatedAt: string
}

export type PublicStore = {
  name: string
  slug: string
  description: string
  logoUrl: string
  currency: string
  storefront: StorefrontConfig
}

export type VariantOption = { name: string; value: string }

export type Variant = {
  id: string
  options: VariantOption[]
  sku: string
  priceCents: number
  stock: number
  imageUrl: string
}

export type Product = {
  id: string
  storeId: string
  name: string
  description: string
  priceCents: number
  images: string[]
  category: string
  optionNames: string[]
  variants: Variant[]
  sku: string
  stock: number
  lowStock: boolean
  published: boolean
  createdAt: string
  updatedAt: string
}

export type Collection = {
  id: string
  storeId: string
  name: string
  slug: string
  description: string
  imageUrl: string
  productIds: string[]
  published: boolean
  createdAt: string
  updatedAt: string
}

export type OrderStatus = 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled'

export type OrderItem = {
  productId: string
  variantId: string
  name: string
  optionLabel: string
  sku: string
  unitPriceCents: number
  quantity: number
  lineTotalCents: number
}

export type Customer = {
  name: string
  email: string
  phone: string
  address: string
  city: string
  notes: string
}

export type Order = {
  id: string
  storeId: string
  status: OrderStatus
  customer: Customer
  items: OrderItem[]
  subtotalCents: number
  totalCents: number
  paymentMethod: string
  createdAt: string
  updatedAt: string
}

export type RevenuePoint = { date: string; revenueCents: number; orders: number }
export type BestSeller = { productId: string; name: string; quantity: number; revenueCents: number }

export type Overview = {
  rangeDays: number
  totalRevenueCents: number
  totalOrders: number
  averageOrderValueCents: number
  totalProducts: number
  unitsInStock: number
  lowStockThreshold: number
  lowStock: Product[]
  recentOrders: Order[]
  bestSellers: BestSeller[]
  revenueOverTime: RevenuePoint[]
}

export type DesignResult = { store: Store; provider: string; summary: string }
