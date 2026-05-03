/**
 * Form Validation Utilities
 * Reusable validation rules and helpers
 */

/**
 * Common validation rules
 */
export const validationRules = {
  required: (value) => {
    if (value === null || value === undefined || value === '') {
      return 'This field is required'
    }
    return null
  },

  email: (value) => {
    if (!value) return null
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(value) ? null : 'Please enter a valid email address'
  },

  minLength: (min) => (value) => {
    if (!value) return null
    return value.length >= min ? null : `Must be at least ${min} characters`
  },

  maxLength: (max) => (value) => {
    if (!value) return null
    return value.length <= max ? null : `Must be at most ${max} characters`
  },

  numeric: (value) => {
    if (!value) return null
    return /^\d+$/.test(value) ? null : 'Must be a number'
  },

  min: (min) => (value) => {
    if (value === null || value === undefined || value === '') return null
    return Number(value) >= min ? null : `Must be at least ${min}`
  },

  max: (max) => (value) => {
    if (value === null || value === undefined || value === '') return null
    return Number(value) <= max ? null : `Must be at most ${max}`
  },

  url: (value) => {
    if (!value) return null
    try {
      new URL(value)
      return null
    } catch {
      return 'Please enter a valid URL'
    }
  },

  price: (value) => {
    if (!value) return null
    const priceRegex = /^\d+(\.\d{1,2})?$/
    return priceRegex.test(value) ? null : 'Please enter a valid price (e.g., 10.99)'
  }
}

/**
 * Validate a value against multiple rules
 * @param {any} value - Value to validate
 * @param {Array} rules - Array of validation functions
 * @returns {string|null} First error message or null
 */
export function validate(value, rules = []) {
  for (const rule of rules) {
    const error = rule(value)
    if (error) return error
  }
  return null
}

/**
 * Validate an entire form object
 * @param {Object} formData - Form data object
 * @param {Object} rules - Validation rules object
 * @returns {Object} Errors object
 */
export function validateForm(formData, rules) {
  const errors = {}
  
  for (const field in rules) {
    const error = validate(formData[field], rules[field])
    if (error) {
      errors[field] = error
    }
  }
  
  return errors
}

/**
 * Check if form has errors
 * @param {Object} errors - Errors object
 * @returns {boolean}
 */
export function hasErrors(errors) {
  return Object.keys(errors).length > 0
}
