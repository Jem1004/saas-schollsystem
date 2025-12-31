export interface School {
  id: number
  name: string
  address?: string
  phone?: string
  email?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
  stats?: SchoolStats
}

export interface SchoolStats {
  totalClasses: number
  totalStudents: number
  totalUsers: number
  totalDevices: number
}

export interface AdminInfo {
  id: number
  username: string
  name: string
  email: string
  isActive: boolean
  createdAt: string
}

export interface SchoolDetail extends School {
  admins?: AdminInfo[]
}

export interface AdminCredentials {
  username: string
  password: string
  name: string
  email: string
  message: string
}

export interface SchoolWithAdmin extends School {
  admin?: AdminCredentials
}

export interface CreateSchoolRequest {
  name: string
  address?: string
  phone?: string
  email?: string
  adminUsername?: string
  adminPassword?: string
  adminName?: string
  adminEmail?: string
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

export interface DeleteSchoolResponse {
  id: number
  name: string
  message: string
  deleted: {
    users: number
    students: number
    classes: number
    devices: number
    attendances: number
  }
}
