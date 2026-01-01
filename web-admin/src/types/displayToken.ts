// Display Token types
// Requirements: 6.1 - Display token management types

export interface DisplayToken {
  id: number
  schoolId: number
  name: string
  isActive: boolean
  lastAccessedAt?: string
  expiresAt?: string
  createdAt: string
  updatedAt: string
  displayUrl?: string
}

export interface DisplayTokenWithSecret {
  id: number
  schoolId: number
  token: string // The actual token value - shown only once
  name: string
  isActive: boolean
  expiresAt?: string
  createdAt: string
  displayUrl: string
}

export interface CreateDisplayTokenRequest {
  name: string
  expires_at?: string // ISO date string, optional expiration
}

export interface UpdateDisplayTokenRequest {
  name?: string
  expires_at?: string
  is_active?: boolean
}

export interface DisplayTokenListResponse {
  tokens: DisplayToken[]
  total: number
}

export interface DisplayTokenValidation {
  valid: boolean
  schoolId?: number
  tokenId?: number
  name?: string
  error?: string
}
