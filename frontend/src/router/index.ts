import { createRouter, createWebHistory } from 'vue-router'
import { hostSlug } from '@/features/stores/host'
import AuthLayout from '@/layouts/AuthLayout.vue'
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import StorefrontLayout from '@/layouts/StorefrontLayout.vue'
import LoginPage from '@/pages/auth/LoginPage.vue'
import RegisterPage from '@/pages/auth/RegisterPage.vue'
import BuilderPage from '@/pages/dashboard/BuilderPage.vue'
import OrderDetailPage from '@/pages/dashboard/OrderDetailPage.vue'
import OrdersPage from '@/pages/dashboard/OrdersPage.vue'
import OverviewPage from '@/pages/dashboard/OverviewPage.vue'
import CollectionFormPage from '@/pages/dashboard/CollectionFormPage.vue'
import CollectionsPage from '@/pages/dashboard/CollectionsPage.vue'
import ProductFormPage from '@/pages/dashboard/ProductFormPage.vue'
import ProductsPage from '@/pages/dashboard/ProductsPage.vue'
import SettingsPage from '@/pages/dashboard/SettingsPage.vue'
import NotFoundPage from '@/pages/NotFoundPage.vue'
import CartPage from '@/pages/storefront/CartPage.vue'
import CheckoutPage from '@/pages/storefront/CheckoutPage.vue'
import CollectionPage from '@/pages/storefront/CollectionPage.vue'
import ProductPage from '@/pages/storefront/ProductPage.vue'
import StoreHomePage from '@/pages/storefront/StoreHomePage.vue'
import { useAuthStore } from '@/stores/auth'
import { useShopStore } from '@/stores/shop'

const storefrontChildren = [
  { path: '', name: 'store-home', component: StoreHomePage },
  { path: 'c/:collectionSlug', name: 'store-collection', component: CollectionPage },
  { path: 'p/:productId', name: 'store-product', component: ProductPage },
  { path: 'cart', name: 'store-cart', component: CartPage },
  { path: 'checkout', name: 'store-checkout', component: CheckoutPage },
]

const dashboardRoutes = [
  { path: '/login', component: AuthLayout, meta: { guest: true }, children: [{ path: '', component: LoginPage }] },
  { path: '/register', component: AuthLayout, meta: { guest: true }, children: [{ path: '', component: RegisterPage }] },
  {
    path: '/',
    component: DashboardLayout,
    meta: { auth: true },
    children: [
      { path: '', component: OverviewPage },
      { path: 'products', component: ProductsPage, meta: { store: true } },
      { path: 'products/new', component: ProductFormPage, meta: { store: true } },
      { path: 'products/:id', component: ProductFormPage, meta: { store: true } },
      { path: 'collections', component: CollectionsPage, meta: { store: true } },
      { path: 'collections/new', component: CollectionFormPage, meta: { store: true } },
      { path: 'collections/:id', component: CollectionFormPage, meta: { store: true } },
      { path: 'orders', component: OrdersPage, meta: { store: true } },
      { path: 'orders/:id', component: OrderDetailPage, meta: { store: true } },
      { path: 'builder', component: BuilderPage, meta: { store: true } },
      { path: 'settings', component: SettingsPage, meta: { store: true } },
    ],
  },
  { path: '/s/:slug', component: StorefrontLayout, children: storefrontChildren },
  { path: '/:pathMatch(.*)*', component: NotFoundPage },
]

const hostRoutes = [
  { path: '/', component: StorefrontLayout, props: { slug: hostSlug }, children: storefrontChildren },
  { path: '/:pathMatch(.*)*', component: NotFoundPage },
]

export const router = createRouter({
  history: createWebHistory(),
  routes: hostSlug ? hostRoutes : dashboardRoutes,
})

router.beforeEach(async (to) => {
  if (hostSlug) return true
  const auth = useAuthStore()
  if (!auth.ready) await auth.bootstrap()
  if (to.meta.auth && !auth.user) return { path: '/login', query: { redirect: to.fullPath } }
  if (to.meta.guest && auth.user) return '/'
  if (to.matched.some((record) => record.meta.store)) {
    const shop = useShopStore()
    await shop.load()
    if (!shop.shop) return '/'
  }
  return true
})
