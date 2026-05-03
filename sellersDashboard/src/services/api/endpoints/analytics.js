/**
 * Analytics API Endpoint
 * All analytics-related API calls
 */

import client from '../client'

const analyticsAPI = {
  /**
   * Get analytics summary
   * @param {string} shopId - Shop ID
   * @returns {Promise<Object>}
   */
  async getSummary(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/summary`)
    return data
  },

  /**
   * Get dashboard analytics (comprehensive)
   * @param {string} shopId - Shop ID
   * @returns {Promise<Object>}
   */
  async getDashboard(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/dashboard`)
    return data
  },

  /**
   * Get revenue over time
   * @param {string} shopId - Shop ID
   * @param {number} days - Number of days (default: 30)
   * @returns {Promise<{labels: Array, values: Array}>}
   */
  async getRevenueOverTime(shopId, days = 30) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/revenue-over-time`, {
      params: { days }
    })
    return data
  },

  /**
   * Get orders over time
   * @param {string} shopId - Shop ID
   * @param {number} days - Number of days (default: 30)
   * @returns {Promise<{labels: Array, values: Array}>}
   */
  async getOrdersOverTime(shopId, days = 30) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/orders-over-time`, {
      params: { days }
    })
    return data
  },

  /**
   * Get top products
   * @param {string} shopId - Shop ID
   * @returns {Promise<Array>}
   */
  async getTopProducts(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/top-products`)
    return data
  },

  /**
   * Get top customers
   * @param {string} shopId - Shop ID
   * @returns {Promise<Array>}
   */
  async getTopCustomers(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/top-customers`)
    return data
  },

  /**
   * Get inventory status
   * @param {string} shopId - Shop ID
   * @returns {Promise<Object>}
   */
  async getInventoryStatus(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/inventory-status`)
    return data
  },

  /**
   * Get customers over time
   * @param {string} shopId - Shop ID
   * @param {number} days - Number of days (default: 30)
   * @returns {Promise<{labels: Array, values: Array}>}
   */
  async getCustomersOverTime(shopId, days = 30) {
    const { data } = await client.get(`/seller/shops/${shopId}/analytics/customers-over-time`, {
      params: { days }
    })
    return data
  }
}

export default analyticsAPI
