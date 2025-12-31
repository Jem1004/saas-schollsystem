import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, LoginResponse, TokenPair } from '@/types/user'
import { authService } from '@/services'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role || null)
  const mustChangePassword = computed(() => user.value?.mustResetPwd || false)

  function setUser(userData: User) {
    user.value = userData
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('accessToken', access)
    localStorage.setItem('refreshToken', refresh)
  }

  function clearAuth() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  function initAuth() {
    const storedAccess = localStorage.getItem('accessToken')
    const storedRefresh = localStorage.getItem('refreshToken')
    const storedUser = localStorage.getItem('user')
    
    if (storedAccess && storedRefresh) {
      accessToken.value = storedAccess
      refreshToken.value = storedRefresh
    }
    
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch {
        localStorage.removeItem('user')
      }
    }
  }

  async function login(credentials: LoginRequest): Promise<LoginResponse> {
    isLoading.value = true
    try {
      const response = await authService.login(credentials)
      setTokens(response.accessToken, response.refreshToken)
      setUser(response.user)
      return response
    } finally {
      isLoading.value = false
    }
  }

  async function logout(): Promise<void> {
    isLoading.value = true
    try {
      await authService.logout()
    } catch {
      // Ignore logout errors - clear local state anyway
    } finally {
      clearAuth()
      isLoading.value = false
    }
  }

  async function refreshAccessToken(): Promise<TokenPair | null> {
    if (!refreshToken.value) {
      clearAuth()
      return null
    }

    try {
      const tokens = await authService.refreshToken(refreshToken.value)
      setTokens(tokens.accessToken, tokens.refreshToken)
      return tokens
    } catch {
      clearAuth()
      return null
    }
  }

  async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
    isLoading.value = true
    try {
      await authService.changePassword(oldPassword, newPassword)
      if (user.value) {
        setUser({ ...user.value, mustResetPwd: false })
      }
    } finally {
      isLoading.value = false
    }
  }

  // Check if user has specific role
  function hasRole(role: string | string[]): boolean {
    if (!user.value) return false
    const roles = Array.isArray(role) ? role : [role]
    return roles.includes(user.value.role)
  }

  // Check if user can access a specific school
  function canAccessSchool(schoolId: number): boolean {
    if (!user.value) return false
    // Super admin can access all schools
    if (user.value.role === 'super_admin') return true
    // Other users can only access their own school
    return user.value.schoolId === schoolId
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    isLoading,
    // Computed
    isAuthenticated,
    userRole,
    mustChangePassword,
    // Actions
    setUser,
    setTokens,
    clearAuth,
    initAuth,
    login,
    logout,
    refreshAccessToken,
    changePassword,
    hasRole,
    canAccessSchool,
  }
})
