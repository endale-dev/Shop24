/**
 * API Error Classes
 * Standardized error types for better error handling across the application
 */

export class APIError extends Error {
  constructor(message, statusCode = 500, originalError = null) {
    super(message)
    this.name = 'APIError'
    this.statusCode = statusCode
    this.originalError = originalError
    this.timestamp = new Date().toISOString()
  }

  toJSON() {
    return {
      name: this.name,
      message: this.message,
      statusCode: this.statusCode,
      timestamp: this.timestamp
    }
  }
}

export class NetworkError extends APIError {
  constructor(message = 'Network connection failed', originalError = null) {
    super(message, 0, originalError)
    this.name = 'NetworkError'
  }
}

export class AuthenticationError extends APIError {
  constructor(message = 'Authentication failed', originalError = null) {
    super(message, 401, originalError)
    this.name = 'AuthenticationError'
  }
}

export class AuthorizationError extends APIError {
  constructor(message = 'Access denied', originalError = null) {
    super(message, 403, originalError)
    this.name = 'AuthorizationError'
  }
}

export class ValidationError extends APIError {
  constructor(message = 'Validation failed', errors = {}, originalError = null) {
    super(message, 400, originalError)
    this.name = 'ValidationError'
    this.errors = errors
  }

  toJSON() {
    return {
      ...super.toJSON(),
      errors: this.errors
    }
  }
}

export class NotFoundError extends APIError {
  constructor(message = 'Resource not found', originalError = null) {
    super(message, 404, originalError)
    this.name = 'NotFoundError'
  }
}

export class ServerError extends APIError {
  constructor(message = 'Server error occurred', originalError = null) {
    super(message, 500, originalError)
    this.name = 'ServerError'
  }
}

/**
 * Transform axios error into our custom error types
 */
export function transformError(error) {
  // Network errors (no response)
  if (!error.response) {
    if (error.request) {
      return new NetworkError('Unable to reach server. Please check your connection.', error)
    }
    return new APIError(error.message, 0, error)
  }

  const { status, data } = error.response
  const message = data?.error || data?.message || error.message

  // Transform based on status code
  switch (status) {
    case 400:
      return new ValidationError(message, data?.errors || {}, error)
    case 401:
      return new AuthenticationError(message, error)
    case 403:
      return new AuthorizationError(message, error)
    case 404:
      return new NotFoundError(message, error)
    case 500:
    case 502:
    case 503:
      return new ServerError(message, error)
    default:
      return new APIError(message, status, error)
  }
}

/**
 * User-friendly error messages
 */
export function getUserMessage(error) {
  if (error instanceof NetworkError) {
    return 'Unable to connect. Please check your internet connection.'
  }
  if (error instanceof AuthenticationError) {
    return 'Please log in to continue.'
  }
  if (error instanceof AuthorizationError) {
    return 'You don\'t have permission to perform this action.'
  }
  if (error instanceof ValidationError) {
    return error.message || 'Please check your input and try again.'
  }
  if (error instanceof NotFoundError) {
    return 'The requested resource was not found.'
  }
  if (error instanceof ServerError) {
    return 'A server error occurred. Please try again later.'
  }
  return error.message || 'An unexpected error occurred.'
}
