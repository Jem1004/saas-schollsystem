import api from './api'
import type { LoginRequest, LoginResponse, TokenPair, User } from '@/types/user'

// API response wrapper
interface ApiResponse<T> {
  success: boolean
  data: T
  error?: {
    code: number
    message: string
  }
}

// Backend login response (snake_case)
interface BackendLoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
  user: {
    id: number
    school_id: number | null
    role: string
    username: string
    email: string
    name: string
    is_active: boolean
    must_reset_pwd: boolean
    last_login_at: string | null
  }
}

// Transform backend user to frontend user
function transformUser(backendUser: BackendLoginResponse['user']): User {
  return {
    id: backendUser.id,
    schoolId: backendUser.school_id ?? undefined,
    role: backendUser.role as User['role'],
    username: backendUser.username,
    email: backendUser.email,
    isActive: backendUser.is_active,
    mustResetPwd: backendUser.must_reset_pwd,
    lastLoginAt: backendUser.last_login_at ?? undefined,
    createdAt: '',
    updatedAt: '',
  }
}

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<ApiResponse<BackendLoginResponse>>('/auth/login', credentials)
    const data = response.data.data
    
    return {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      user: transformUser(data.user),
    }
  },

  async refreshToken(refreshToken: string): Promise<TokenPair> {
    const response = await api.post<ApiResponse<{ access_token: string; refresh_token: string }>>('/auth/refresh', { refresh_token: refreshToken })
    return {
      accessToken: response.data.data.access_token,
      refreshToken: response.data.data.refresh_token,
    }
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout')
  },

  async changePassword(oldPassword: string, newPassword: string): Promise<void> {
    await api.post('/auth/change-password', { old_password: oldPassword, new_password: newPassword })
  },
}
