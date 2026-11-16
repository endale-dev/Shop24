export class ApiError extends Error {
  status: number
  code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json')
  const response = await fetch(path, { ...init, headers, credentials: 'include' })
  if (response.status === 204) return undefined as T
  const text = await response.text()
  const body = text ? (JSON.parse(text) as { data?: T; error?: { code?: string; message?: string } }) : null
  if (!response.ok) {
    throw new ApiError(response.status, body?.error?.code ?? 'error', body?.error?.message ?? 'Something went wrong')
  }
  return body?.data as T
}
