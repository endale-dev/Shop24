// sellersDashboard/src/services/api/endpoints/themes.js
import apiClient from '../client'

/**
 * Theme API endpoints for customizing storefront appearance
 */
export const themesAPI = {
  /**
   * Get active theme for shop
   * @param {string} shopId - Shop ID
   * @returns {Promise} Active theme data
   */
  async getActiveTheme(shopId) {
    const response = await apiClient.get(`/seller/shops/${shopId}/themes/active`)
    return response.data
  },

  /**
   * List all themes for shop
   * @param {string} shopId - Shop ID
   * @returns {Promise} Array of themes
   */
  async listThemes(shopId) {
    const response = await apiClient.get(`/seller/shops/${shopId}/themes`)
    return response.data
  },

  /**
   * Create a new theme
   * @param {string} shopId - Shop ID
   * @param {Object} themeData - Theme configuration
   * @returns {Promise} Created theme
   */
  async createTheme(shopId, themeData) {
    const response = await apiClient.post(`/seller/shops/${shopId}/themes`, themeData)
    return response.data
  },

  /**
   * Activate a theme
   * @param {string} shopId - Shop ID
   * @param {string} themeId - Theme ID to activate
   * @returns {Promise} Activation result
   */
  async activateTheme(shopId, themeId) {
    const response = await apiClient.patch(`/seller/shops/${shopId}/themes/${themeId}/activate`)
    return response.data
  },

  /**
   * Update theme colors
   * @param {string} shopId - Shop ID
   * @param {string} themeId - Theme ID
   * @param {Object} colors - Colors object with 12 color fields
   * @returns {Promise} Update result
   */
  async updateColors(shopId, themeId, colors) {
    const response = await apiClient.patch(`/seller/shops/${shopId}/themes/${themeId}/colors`, colors)
    return response.data
  },

  /**
   * Update theme typography
   * @param {string} shopId - Shop ID
   * @param {string} themeId - Theme ID
   * @param {Object} typography - Typography settings
   * @returns {Promise} Update result
   */
  async updateTypography(shopId, themeId, typography) {
    const response = await apiClient.patch(`/seller/shops/${shopId}/themes/${themeId}/typography`, typography)
    return response.data
  },

  /**
   * Update theme layout
   * @param {string} shopId - Shop ID
   * @param {string} themeId - Theme ID
   * @param {Object} layout - Layout configuration
   * @returns {Promise} Update result
   */
  async updateLayout(shopId, themeId, layout) {
    const response = await apiClient.patch(`/seller/shops/${shopId}/themes/${themeId}/layout`, layout)
    return response.data
  },

  /**
   * Update theme assets (logo, favicon, etc.)
   * @param {string} shopId - Shop ID
   * @param {string} themeId - Theme ID
   * @param {Object} assets - Assets object
   * @returns {Promise} Update result
   */
  async updateAssets(shopId, themeId, assets) {
    const response = await apiClient.patch(`/seller/shops/${shopId}/themes/${themeId}/assets`, assets)
    return response.data
  }
}

/**
 * Default color palette for new themes
 */
export const defaultColors = {
  primary: '#10B981',
  secondary: '#6366F1',
  accent: '#F59E0B',
  background: '#FFFFFF',
  surface: '#F9FAFB',
  textPrimary: '#111827',
  textSecondary: '#6B7280',
  border: '#E5E7EB',
  success: '#10B981',
  warning: '#F59E0B',
  error: '#EF4444',
  info: '#3B82F6'
}

/**
 * Default typography settings
 */
export const defaultTypography = {
  headingFont: 'Inter',
  bodyFont: 'Inter',
  headingWeight: '600',
  bodyWeight: '400',
  headingSizeBase: '2rem',
  bodySizeBase: '1rem',
  lineHeight: '1.5',
  letterSpacing: 'normal'
}

/**
 * Default layout configuration
 */
export const defaultLayout = {
  headerStyle: 'classic',
  containerWidth: 'boxed',
  productGridCols: 3,
  borderRadius: '0.5rem',
  spacing: 'normal',
  cardStyle: 'elevated',
  buttonStyle: 'rounded',
  showBreadcrumbs: true,
  showProductGrid: true
}
