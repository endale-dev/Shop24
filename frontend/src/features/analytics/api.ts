import { request } from '@/services/api'
import type { Overview } from '@/types'

export function overview(days: number) {
  return request<Overview>(`/api/analytics?days=${days}`)
}
