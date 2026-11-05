import { apiClient } from '@/api/client'
import type { Billet, BilletListResponse, BilletStatus } from '@/types/payment'

export async function fetchAdminBillets(): Promise<BilletListResponse> {
  const { data } = await apiClient.get<BilletListResponse>('/api/admin/billets')
  return data
}

export async function updateAdminBillet(
  id: string,
  payload: { status?: BilletStatus; comments?: string },
): Promise<Billet> {
  const { data } = await apiClient.patch<Billet>(`/api/admin/billets/${id}`, payload)
  return data
}
