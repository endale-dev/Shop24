/**
 * Products API Endpoint
 * All product-related API calls
 */

import client from './client'

const productsAPI = {
  /**
   * Fetch all products for a shop with pagination
   * @param {string} shopId - Shop ID
   * @param {Object} params - Query parameters (page, limit, search, collection_id, stock_status)
   * @returns {Promise<{products: Array, pagination: Object, stats: Object}>}
   */
  async getAll(shopId, params = {}) {
    const { data } = await client.get(`/seller/shops/${shopId}/products`, { params })
    return data
  },

  /**
   * Get a single product by ID
   * @param {string} shopId - Shop ID
   * @param {string} productId - Product ID
   * @returns {Promise<Object>}
   */
  async getOne(shopId, productId) {
    const { data } = await client.get(`/seller/shops/${shopId}/products/${productId}`)
    return data
  },

  /**
   * Create a new product
   * @param {string} shopId - Shop ID
   * @param {Object} productData - Product data
   * @returns {Promise<Object>}
   */
  async create(shopId, productData) {
    const { data } = await client.post(`/seller/shops/${shopId}/products`, productData)
    return data
  },

  /**
   * Update an existing product
   * @param {string} shopId - Shop ID
   * @param {string} productId - Product ID
   * @param {Object} updates - Product updates
   * @returns {Promise<Object>}
   */
  async update(shopId, productId, updates) {
    const { data } = await client.patch(`/seller/shops/${shopId}/products/${productId}`, updates)
    return data
  },

  /**
   * Delete a product
   * @param {string} shopId - Shop ID
   * @param {string} productId - Product ID
   * @returns {Promise<Object>}
   */
  async delete(shopId, productId) {
    const { data} = await client.delete(`/seller/shops/${shopId}/products/${productId}`)
    return data
  }
}

export default productsAPI
