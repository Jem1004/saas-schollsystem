import api from './api'
import type {
  School,
  SchoolWithAdmin,
  SchoolDetail,
  CreateSchoolRequest,
  UpdateSchoolRequest,
  SchoolListResponse,
  TenantStats,
  AdminCredentials,
  AdminInfo,
  DeleteSchoolResponse,
} from '@/types/tenant'

export interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
}

// Helper to transform snake_case to camelCase for school response
function transformSchool(data: Record<string, unknown>): School {
  return {
    id: data.id as number,
    name: data.name as string,
    address: data.address as string | undefined,
    phone: data.phone as string | undefined,
    email: data.email as string | undefined,
    isActive: data.is_active as boolean,
    createdAt: data.created_at as string,
    updatedAt: data.updated_at as string,
    stats: data.stats ? {
      totalClasses: (data.stats as Record<string, number>).total_classes || 0,
      totalStudents: (data.stats as Record<string, number>).total_students || 0,
      totalUsers: (data.stats as Record<string, number>).total_users || 0,
      totalDevices: (data.stats as Record<string, number>).total_devices || 0,
    } : undefined,
  }
}

// Helper to transform admin credentials
function transformAdmin(data: Record<string, unknown> | undefined): AdminCredentials | undefined {
  if (!data) return undefined
  return {
    username: data.username as string,
    password: data.password as string,
    name: data.name as string,
    email: data.email as string,
    message: data.message as string,
  }
}

// Helper to transform admin info
function transformAdminInfo(data: Record<string, unknown>): AdminInfo {
  return {
    id: data.id as number,
    username: data.username as string,
    name: data.name as string,
    email: data.email as string,
    isActive: data.is_active as boolean,
    createdAt: data.created_at as string,
  }
}

export const tenantService = {
  async getSchools(params?: PaginationParams): Promise<SchoolListResponse> {
    const queryParams: Record<string, unknown> = {}
    if (params?.page) queryParams.page = params.page
    if (params?.pageSize) queryParams.page_size = params.pageSize
    if (params?.search) queryParams.name = params.search
    
    const response = await api.get('/schools', { params: queryParams })
    const result = response.data
    
    // Transform backend response to frontend format
    if (result.success && result.data) {
      const backendData = result.data
      return {
        data: (backendData.schools || []).map(transformSchool),
        total: backendData.pagination?.total || 0,
        page: backendData.pagination?.page || 1,
        pageSize: backendData.pagination?.page_size || 20,
      }
    }
    
    return { data: [], total: 0, page: 1, pageSize: 20 }
  },

  async getSchool(id: number): Promise<School> {
    const response = await api.get(`/schools/${id}`)
    const result = response.data
    if (result.success && result.data) {
      return transformSchool(result.data)
    }
    throw new Error('School not found')
  },

  async getSchoolDetail(id: number): Promise<SchoolDetail> {
    const response = await api.get(`/schools/${id}/detail`)
    const result = response.data
    if (result.success && result.data) {
      const school = transformSchool(result.data)
      const admins = (result.data.admins as Record<string, unknown>[] || []).map(transformAdminInfo)
      return {
        ...school,
        admins,
      }
    }
    throw new Error('School not found')
  },

  async createSchool(data: CreateSchoolRequest): Promise<SchoolWithAdmin> {
    // Transform camelCase to snake_case for backend
    const backendData: Record<string, unknown> = {
      name: data.name,
      address: data.address,
      phone: data.phone,
      email: data.email,
    }
    if (data.adminUsername) backendData.admin_username = data.adminUsername
    if (data.adminPassword) backendData.admin_password = data.adminPassword
    if (data.adminName) backendData.admin_name = data.adminName
    if (data.adminEmail) backendData.admin_email = data.adminEmail
    
    const response = await api.post('/schools', backendData)
    const result = response.data
    if (result.success && result.data) {
      return {
        ...transformSchool(result.data),
        admin: transformAdmin(result.data.admin as Record<string, unknown> | undefined),
      }
    }
    throw new Error(result.error?.message || 'Failed to create school')
  },

  async updateSchool(id: number, data: UpdateSchoolRequest): Promise<School> {
    const response = await api.put(`/schools/${id}`, data)
    const result = response.data
    if (result.success && result.data) {
      return transformSchool(result.data)
    }
    throw new Error(result.error?.message || 'Failed to update school')
  },

  async activateSchool(id: number): Promise<School> {
    const response = await api.post(`/schools/${id}/activate`)
    const result = response.data
    if (result.success && result.data) {
      return transformSchool(result.data)
    }
    throw new Error(result.error?.message || 'Failed to activate school')
  },

  async deactivateSchool(id: number): Promise<School> {
    const response = await api.post(`/schools/${id}/deactivate`)
    const result = response.data
    if (result.success && result.data) {
      return transformSchool(result.data)
    }
    throw new Error(result.error?.message || 'Failed to deactivate school')
  },

  async deleteSchool(id: number): Promise<DeleteSchoolResponse> {
    const response = await api.delete(`/schools/${id}`)
    const result = response.data
    if (result.success && result.data) {
      const data = result.data
      return {
        id: data.id as number,
        name: data.name as string,
        message: data.message as string,
        deleted: {
          users: (data.deleted as Record<string, number>)?.users || 0,
          students: (data.deleted as Record<string, number>)?.students || 0,
          classes: (data.deleted as Record<string, number>)?.classes || 0,
          devices: (data.deleted as Record<string, number>)?.devices || 0,
          attendances: (data.deleted as Record<string, number>)?.attendances || 0,
        },
      }
    }
    throw new Error(result.error?.message || 'Failed to delete school')
  },

  async getStats(): Promise<TenantStats> {
    // Stats endpoint doesn't exist in backend yet, calculate from schools list
    try {
      const response = await api.get('/schools', { params: { page_size: 1000 } })
      const result = response.data
      if (result.success && result.data) {
        const schools = result.data.schools || []
        const activeSchools = schools.filter((s: Record<string, unknown>) => s.is_active).length
        return {
          totalSchools: schools.length,
          activeSchools,
          inactiveSchools: schools.length - activeSchools,
        }
      }
    } catch {
      // Return default stats on error
    }
    return { totalSchools: 0, activeSchools: 0, inactiveSchools: 0 }
  },
}
