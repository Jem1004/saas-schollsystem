import api from './api'
import type {
  Device,
  RegisterDeviceRequest,
  UpdateDeviceRequest,
  DeviceListResponse,
  DeviceStats,
  SchoolDevices,
  GroupedDevicesResponse,
} from '@/types/device'

export interface DevicePaginationParams {
  page?: number
  pageSize?: number
  search?: string
  schoolId?: number
  isActive?: boolean
}

// Helper to transform snake_case to camelCase for device response
function transformDevice(data: Record<string, unknown>): Device {
  return {
    id: data.id as number,
    schoolId: data.school_id as number,
    schoolName: data.school_name as string | undefined,
    deviceCode: data.device_code as string,
    apiKey: data.api_key as string || '',
    description: data.description as string | undefined,
    isActive: data.is_active as boolean,
    lastSeenAt: data.last_seen_at as string | undefined,
    createdAt: data.created_at as string,
    updatedAt: data.updated_at as string,
  }
}

// Helper to transform school devices
function transformSchoolDevices(data: Record<string, unknown>): SchoolDevices {
  const devices = (data.devices as Record<string, unknown>[] || []).map(transformDevice)
  return {
    schoolId: data.school_id as number,
    schoolName: data.school_name as string,
    isActive: data.is_active as boolean,
    deviceCount: data.device_count as number,
    devices,
  }
}

export const deviceService = {
  async getDevices(params?: DevicePaginationParams): Promise<DeviceListResponse> {
    const queryParams: Record<string, unknown> = {}
    if (params?.page) queryParams.page = params.page
    if (params?.pageSize) queryParams.page_size = params.pageSize
    if (params?.schoolId) queryParams.school_id = params.schoolId
    if (params?.search) queryParams.device_code = params.search
    if (params?.isActive !== undefined) queryParams.is_active = params.isActive
    
    const response = await api.get('/devices', { params: queryParams })
    const result = response.data
    
    // Transform backend response to frontend format
    if (result.success && result.data) {
      const backendData = result.data
      return {
        data: (backendData.devices || []).map(transformDevice),
        total: backendData.pagination?.total || 0,
        page: backendData.pagination?.page || 1,
        pageSize: backendData.pagination?.page_size || 20,
      }
    }
    
    return { data: [], total: 0, page: 1, pageSize: 20 }
  },

  async getDevicesGrouped(): Promise<GroupedDevicesResponse> {
    const response = await api.get('/devices/grouped')
    const result = response.data
    
    if (result.success && result.data) {
      const backendData = result.data
      return {
        schools: (backendData.schools || []).map(transformSchoolDevices),
        total: backendData.total || 0,
      }
    }
    
    return { schools: [], total: 0 }
  },

  async getDevice(id: number): Promise<Device> {
    const response = await api.get(`/devices/${id}`)
    const result = response.data
    if (result.success && result.data) {
      return transformDevice(result.data)
    }
    throw new Error('Device not found')
  },

  async registerDevice(data: RegisterDeviceRequest): Promise<Device> {
    // Transform camelCase to snake_case for backend
    const backendData = {
      school_id: data.schoolId,
      device_code: data.deviceCode,
      description: data.description,
    }
    const response = await api.post('/devices', backendData)
    const result = response.data
    if (result.success && result.data) {
      return transformDevice(result.data)
    }
    throw new Error(result.error?.message || 'Failed to register device')
  },

  async updateDevice(id: number, data: UpdateDeviceRequest): Promise<Device> {
    const response = await api.put(`/devices/${id}`, data)
    const result = response.data
    if (result.success && result.data) {
      return transformDevice(result.data)
    }
    throw new Error(result.error?.message || 'Failed to update device')
  },

  async deleteDevice(id: number): Promise<void> {
    const response = await api.delete(`/devices/${id}`)
    const result = response.data
    if (!result.success) {
      throw new Error(result.error?.message || 'Failed to delete device')
    }
  },

  async revokeApiKey(id: number): Promise<void> {
    const response = await api.post(`/devices/${id}/revoke`)
    const result = response.data
    if (!result.success) {
      throw new Error(result.error?.message || 'Failed to revoke API key')
    }
  },

  async regenerateApiKey(id: number): Promise<Device> {
    const response = await api.post(`/devices/${id}/regenerate`)
    const result = response.data
    if (result.success && result.data) {
      return transformDevice(result.data)
    }
    throw new Error(result.error?.message || 'Failed to regenerate API key')
  },

  async getStats(): Promise<DeviceStats> {
    // Stats endpoint doesn't exist in backend yet, calculate from devices list
    try {
      const response = await api.get('/devices', { params: { page_size: 1000 } })
      const result = response.data
      if (result.success && result.data) {
        const devices = result.data.devices || []
        const activeDevices = devices.filter((d: Record<string, unknown>) => d.is_active).length
        const now = new Date()
        const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000)
        const onlineDevices = devices.filter((d: Record<string, unknown>) => {
          if (!d.is_active || !d.last_seen_at) return false
          return new Date(d.last_seen_at as string) > fiveMinutesAgo
        }).length
        return {
          totalDevices: devices.length,
          activeDevices,
          inactiveDevices: devices.length - activeDevices,
          onlineDevices,
        }
      }
    } catch {
      // Return default stats on error
    }
    return { totalDevices: 0, activeDevices: 0, inactiveDevices: 0, onlineDevices: 0 }
  },
}
