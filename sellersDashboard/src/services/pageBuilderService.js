import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * Page Builder Service - API client for page CRUD operations
 */

// Get auth headers helper
function getAuthHeaders() {
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

/**
 * Fetch all pages for a shop
 * @param {string} shopId - The shop ObjectID
 * @returns {Promise<Array>} List of pages
 */
export async function fetchPages(shopId) {
  const response = await axios.get(`${API_URL}/seller/shops/${shopId}/pages`, {
    headers: getAuthHeaders()
  });
  return response.data.pages || [];
}

/**
 * Fetch a single page by ID
 * @param {string} shopId - The shop ObjectID
 * @param {string} pageId - The page ObjectID
 * @returns {Promise<Object>} Page object
 */
export async function fetchPage(shopId, pageId) {
  const response = await axios.get(`${API_URL}/seller/shops/${shopId}/pages/${pageId}`, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Create a new page
 * @param {string} shopId - The shop ObjectID
 * @param {Object} pageData - Page data { path, title, blocks, seo, visibility, status }
 * @returns {Promise<Object>} Created page
 */
export async function createPage(shopId, pageData) {
  const response = await axios.post(`${API_URL}/seller/shops/${shopId}/pages`, pageData, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Update an existing page
 * @param {string} shopId - The shop ObjectID
 * @param {string} pageId - The page ObjectID
 * @param {Object} pageData - Updated page data
 * @returns {Promise<Object>} Result
 */
export async function updatePage(shopId, pageId, pageData) {
  const response = await axios.patch(`${API_URL}/seller/shops/${shopId}/pages/${pageId}`, pageData, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Delete a page
 * @param {string} shopId - The shop ObjectID
 * @param {string} pageId - The page ObjectID
 * @returns {Promise<Object>} Result
 */
export async function deletePage(shopId, pageId) {
  const response = await axios.delete(`${API_URL}/seller/shops/${shopId}/pages/${pageId}`, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Publish a page (create snapshot for storefront)
 * @param {string} shopId - The shop ObjectID
 * @param {string} pageId - The page ObjectID
 * @param {Object} pageData - Full page data for snapshot
 * @returns {Promise<Object>} Result
 */
export async function publishPage(shopId, pageId, pageData) {
  const response = await axios.post(`${API_URL}/seller/shops/${shopId}/pages/${pageId}/publish`, pageData, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Fetch active theme for a shop
 * @param {string} shopId
 * @returns {Promise<Object>} Theme object
 */
export async function fetchActiveTheme(shopId) {
  const response = await axios.get(`${API_URL}/seller/shops/${shopId}/themes/active`, {
    headers: getAuthHeaders()
  });
  return response.data.theme;
}

/**
 * Update theme colors
 * @param {string} shopId
 * @param {string} themeId
 * @param {Object} colors
 */
export async function updateThemeColors(shopId, themeId, colors) {
  const response = await axios.patch(`${API_URL}/seller/shops/${shopId}/themes/${themeId}/colors`, colors, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Update theme typography
 * @param {string} shopId
 * @param {string} themeId
 * @param {Object} typography
 */
export async function updateThemeTypography(shopId, themeId, typography) {
  const response = await axios.patch(`${API_URL}/seller/shops/${shopId}/themes/${themeId}/typography`, typography, {
    headers: getAuthHeaders()
  });
  return response.data;
}

/**
 * Update theme layout
 * @param {string} shopId
 * @param {string} themeId
 * @param {Object} layout
 */
export async function updateThemeLayout(shopId, themeId, layout) {
  const response = await axios.patch(`${API_URL}/seller/shops/${shopId}/themes/${themeId}/layout`, layout, {
    headers: getAuthHeaders()
  });
  return response.data;
}

export default {
  fetchPages,
  fetchPage,
  createPage,
  updatePage,
  deletePage,
  publishPage,
  fetchActiveTheme,
  updateThemeColors,
  updateThemeTypography,
  updateThemeLayout
};
