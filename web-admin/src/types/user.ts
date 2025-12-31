export type UserRole = 'super_admin' | 'admin_sekolah' | 'guru_bk' | 'wali_kelas'

export interface User {
  id: number
  schoolId?: number
  role: UserRole
  username: string
  email?: string
  isActive: boolean
  mustResetPwd: boolean
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export interface TokenPair {
  accessToken: string
  refreshToken: string
}
