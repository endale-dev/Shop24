/**
 * Common utility functions for formatting and validation
 */

/**
 * Format currency values
 * @param {number} amount - Amount to format
 * @param {string} currency - Currency symbol (default: $)
 * @returns {string}
 */
export function formatCurrency(amount, currency = '$') {
  const formatted = (amount || 0).toFixed(2)
  return `${currency}${formatted}`
}

/**
 * Format date to readable string
 * @param {string|Date} date - Date to format
 * @param {string} format - Format string (default: 'MMM d, yyyy')
 * @returns {string}
 */
export function formatDate(date, formatStr = 'MMM d, yyyy') {
  if (!date) return 'N/A'
  try {
    const { format } = require('date-fns')
    return format(new Date(date), formatStr)
  } catch (e) {
    return 'Invalid Date'
  }
}

/**
 * Debounce function calls
 * @param {Function} func - Function to debounce
 * @param {number} wait - Wait time in ms
 * @returns {Function}
 */
export function debounce(func, wait = 300) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

/**
 * Truncate text to specified length
 * @param {string} text - Text to truncate
 * @param {number} length - Max length
 * @returns {string}
 */
export function truncate(text, length = 50) {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}

/**
 * Capitalize first letter
 * @param {string} text - Text to capitalize
 * @returns {string}
 */
export function capitalize(text) {
  if (!text) return ''
  return text.charAt(0).toUpperCase() + text.slice(1)
}

/**
 * Get initials from name
 * @param {string} name - Full name
 * @returns {string}
 */
export function getInitials(name) {
  if (!name) return ''
  const parts = name.split(' ')
  if (parts.length === 1) return parts[0].charAt(0).toUpperCase()
  return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase()
}
