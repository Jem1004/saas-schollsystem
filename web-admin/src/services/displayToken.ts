import api from './api'
import type {
  DisplayToken,
  DisplayTokenWithSecret,
  CreateDisplayTokenRequest,
  UpdateDisplayTokenRequest,
  DisplayTokenListResponse,
} from '@/types/displayToken'

// Transform backend response to frontend format (snake_case to camelCase)
function transformToken(data: Record<string, unknown>): DisplayToken {
  return {
    id: data.id as number,
    schoolId: data.school_id as number,
    name: data.name as string,
    isActive: data.is_active as boolean,
    lastAccessedAt: data.last_accessed_at as string | undefined,
    expiresAt: data.expires_at as string | undefined,
    createdAt: data.created_at as string,
    updatedAt: data.updated_at as string,
    displayUrl: data.display_url as string | undefined,
  }
}

function transformTokenWithSecret(data: Record<string, unknown>): DisplayTokenWithSecret {
  return {
    id: data.id as number,
    schoolId: data.school_id as number,
    token: data.token as string,
    name: data.name as string,
    isActive: data.is_active as boolean,
    expiresAt: data.expires_at as string | undefined,
    createdAt: data.created_at as string,
    displayUrl: data.display_url as string,
  }
}

export const displayTokenService = {
  // Get all display tokens for the school
  async getTokens(): Promise<DisplayTokenListResponse> {
    const response = await api.get<{ data: { tokens: Record<string, unknown>[]; total: number } }>('/display-tokens')
    return {
      tokens: (response.data.data.tokens || []).map(transformToken),
      total: response.data.data.total || 0,
    }
  },

  // Get a single display token by ID
  async getTokenById(id: number): Promise<DisplayToken> {
    const response = await api.get<{ data: Record<string, unknown> }>(`/display-tokens/${id}`)
    return transformToken(response.data.data)
  },

  // Create a new display token
  async createToken(data: CreateDisplayTokenRequest): Promise<DisplayTokenWithSecret> {
    const response = await api.post<{ data: Record<string, unknown> }>('/display-tokens', data)
    return transformTokenWithSecret(response.data.data)
  },

  // Update an existing display token
  async updateToken(id: number, data: UpdateDisplayTokenRequest): Promise<DisplayToken> {
    const response = await api.put<{ data: Record<string, unknown> }>(`/display-tokens/${id}`, data)
    return transformToken(response.data.data)
  },

  // Delete a display token
  async deleteToken(id: number): Promise<void> {
    await api.delete(`/display-tokens/${id}`)
  },

  // Revoke a display token (mark as inactive)
  async revokeToken(id: number): Promise<void> {
    await api.post(`/display-tokens/${id}/revoke`)
  },

  // Regenerate a display token (creates new token value)
  async regenerateToken(id: number): Promise<DisplayTokenWithSecret> {
    const response = await api.post<{ data: Record<string, unknown> }>(`/display-tokens/${id}/regenerate`)
    return transformTokenWithSecret(response.data.data)
  },
}

export default displayTokenService
