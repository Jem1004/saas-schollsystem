export interface School {
  id: number
  name: string
  address?: string
  phone?: string
  email?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateSchoolRequest {
  name: string
  address?: string
  phone?: string
  email?: string
}

export interface UpdateSchoolRequest {
  name?: string
  address?: string
  phone?: string
  email?: string
}

export interface SchoolListResponse {
  data: School[]
  total: number
  page: number
  pageSize: number
}

export interface TenantStats {
  totalSchools: number
  activeSchools: number
  inactiveSchools: number
}
