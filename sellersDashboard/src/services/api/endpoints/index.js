/**
 * API Endpoints Index
 * Centralized export of all API endpoints
 */

import productsAPI from './products'
import analyticsAPI from './analytics'
import ordersAPI from './orders'
import customersAPI from './customers'
import { themesAPI } from './themes'

export {
  productsAPI,
  analyticsAPI,
  ordersAPI,
  customersAPI,
  themesAPI
}

// Default export for convenience
export default {
  products: productsAPI,
  analytics: analyticsAPI,
  orders: ordersAPI,
  customers: customersAPI,
  themes: themesAPI
}
