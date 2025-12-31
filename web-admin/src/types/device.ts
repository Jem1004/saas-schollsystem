export interface Device {
  id: number
  schoolId: number
  schoolName?: string
  deviceCode: string
  apiKey: string
  description?: string
  isActive: boolean
  lastSeenAt?: string
  createdAt: string
  updatedAt: string
}

export interface RegisterDeviceRequest {
  schoolId: number
  deviceCode: string
  description?: string
}

export interface UpdateDeviceRequest {
  description?: string
}

export interface DeviceListResponse {
  data: Device[]
  total: number
  page: number
  pageSize: number
}

export interface SchoolDevices {
  schoolId: number
  schoolName: string
  isActive: boolean
  deviceCount: number
  devices: Device[]
}

export interface GroupedDevicesResponse {
  schools: SchoolDevices[]
  total: number
}

export interface DeviceStats {
  totalDevices: number
  activeDevices: number
  inactiveDevices: number
  onlineDevices: number
}
