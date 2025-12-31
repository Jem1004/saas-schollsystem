import api from './api'
import type {
  Device,
  RegisterDeviceRequest,
  UpdateDeviceRequest,
  DeviceListResponse,
  DeviceStats,
} from '@/types/device'

export interface DevicePaginationParams {
  page?: number
  pageSize?: number
  search?: string
  schoolId?: number
  isActive?: boolean
}

export const deviceService = {
  async getDevices(params?: DevicePaginationParams): Promise<DeviceListResponse> {
    const response = await api.get('/devices', { params })
    return response.data
  },

  async getDevice(id: number): Promise<Device> {
    const response = await api.get(`/devices/${id}`)
    return response.data
  },

  async registerDevice(data: RegisterDeviceRequest): Promise<Device> {
    const response = await api.post('/devices', data)
    return response.data
  },

  async updateDevice(id: number, data: UpdateDeviceRequest): Promise<Device> {
    const response = await api.put(`/devices/${id}`, data)
    return response.data
  },

  async revokeApiKey(id: number): Promise<void> {
    await api.post(`/devices/${id}/revoke-key`)
  },

  async regenerateApiKey(id: number): Promise<Device> {
    const response = await api.post(`/devices/${id}/regenerate-key`)
    return response.data
  },

  async getStats(): Promise<DeviceStats> {
    const response = await api.get('/devices/stats')
    return response.data
  },
}
