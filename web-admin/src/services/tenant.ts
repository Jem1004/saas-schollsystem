import api from './api'
import type {
  School,
  CreateSchoolRequest,
  UpdateSchoolRequest,
  SchoolListResponse,
  TenantStats,
} from '@/types/tenant'

export interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
}

export const tenantService = {
  async getSchools(params?: PaginationParams): Promise<SchoolListResponse> {
    const response = await api.get('/tenants', { params })
    return response.data
  },

  async getSchool(id: number): Promise<School> {
    const response = await api.get(`/tenants/${id}`)
    return response.data
  },

  async createSchool(data: CreateSchoolRequest): Promise<School> {
    const response = await api.post('/tenants', data)
    return response.data
  },

  async updateSchool(id: number, data: UpdateSchoolRequest): Promise<School> {
    const response = await api.put(`/tenants/${id}`, data)
    return response.data
  },

  async activateSchool(id: number): Promise<School> {
    const response = await api.post(`/tenants/${id}/activate`)
    return response.data
  },

  async deactivateSchool(id: number): Promise<School> {
    const response = await api.post(`/tenants/${id}/deactivate`)
    return response.data
  },

  async getStats(): Promise<TenantStats> {
    const response = await api.get('/tenants/stats')
    return response.data
  },
}
