/**
 * Customers API Endpoint
 * All customer-related API calls
 */

import client from '../client'

const customersAPI = {
  /**
   * Fetch all customers for a shop
   * @param {string} shopId - Shop ID
   * @param {Object} params - Query parameters
   * @returns {Promise<Array>}
   */
  async getAll(shopId, params = {}) {
    const { data } = await client.get(`/seller/shops/${shopId}/customers`, { params })
    return data
  },

  /**
   * Get customer statistics
   * @param {string} shopId - Shop ID
   * @returns {Promise<Object>}
   */
  async getStats(shopId) {
    const { data } = await client.get(`/seller/shops/${shopId}/customers/stats`)
    return data
  },

  /**
   * Get a single customer by ID
   * @param {string} shopId - Shop ID
   * @param {string} customerId - Customer ID
   * @returns {Promise<Object>}
   */
  async getOne(shopId, customerId) {
    const { data } = await client.get(`/seller/shops/${shopId}/customers/${customerId}`)
    return data
  },

  /**
   * Get customer order history
   * @param {string} shopId - Shop ID
   * @param {string} customerId - Customer ID
   * @returns {Promise<Array>}
   */
  async getOrderHistory(shopId, customerId) {
    const { data } = await client.get(`/seller/shops/${shopId}/customers/${customerId}/history`)
    return data
  }
}

export default customersAPI
