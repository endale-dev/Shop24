/**
 * Orders API Endpoint
 * All order-related API calls
 */

import client from '../client'

const ordersAPI = {
  /**
   * Fetch all orders for a shop with pagination
   * @param {string} shopId - Shop ID
   * @param {Object} params - Query parameters (page, limit, status)
   * @returns {Promise<{orders: Array, pagination: Object}>}
   */
  async getAll(shopId, params = {}) {
    const { data } = await client.get(`/seller/shops/${shopId}/orders`, { params })
    return data
  },

  /**
   * Get orders for dashboard view
   * @param {string} shopId - Shop ID
   * @returns {Promise<Array>}
   */
  async getDashboard(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/orders/dashboard`)
    return data
  },

  /**
   * Get order statistics
   * @param {string} shopId - Shop ID
   * @returns {Promise<Object>}
   */
  async getStats(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/orders/stats`)
    return data
  },

  /**
   * Get a single order by ID
   * @param {string} shopId - Shop ID
   * @param {string} orderId - Order ID
   * @returns {Promise<Object>}
   */
  async getOne(shopId, orderId) {
    const { data } = await client.get(`/seller/shops/${shopId}/orders/${orderId}`)
    return data
  },

  /**
   * Get order with customer details
   * @param {string} shopId - Shop ID
   * @param {string} orderId - Order ID
   * @returns {Promise<Object>}
   */
  async getWithDetails(shopId, orderId) {
    const { data } = await client.get(`/seller/shops/${shopId}/orders/${orderId}/details`)
    return data
  },

  /**
   * Update an order (status, etc.)
   * @param {string} shopId - Shop ID
   * @param {string} orderId - Order ID
   * @param {Object} updates - Order updates
   * @returns {Promise<Object>}
   */
  async update(shopId, orderId, updates) {
    const { data } = await client.patch(`/seller/shops/${shopId}/orders/${orderId}`, updates)
    return data
  },

  /**
   * Delete an order
   * @param {string} shopId - Shop ID
   * @param {string} orderId - Order ID
   * @returns {Promise<Object>}
   */
  async delete(shopId, orderId) {
    const { data } = await client.delete(`/seller/shops/${shopId}/orders/${orderId}`)
    return data
  }
}

export default ordersAPI
