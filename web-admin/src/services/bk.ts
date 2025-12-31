import api from './api'
import type {
  Violation,
  Achievement,
  Permit,
  CounselingNote,
  StudentBKProfile,
  BKStats,
  CreateViolationRequest,
  CreateAchievementRequest,
  CreatePermitRequest,
  UpdatePermitReturnRequest,
  CreateCounselingNoteRequest,
  ViolationListResponse,
  AchievementListResponse,
  PermitListResponse,
  CounselingNoteListResponse,
} from '@/types/bk'

export interface BKPaginationParams {
  page?: number
  pageSize?: number
  studentId?: number
  search?: string
  startDate?: string
  endDate?: string
  level?: string
  category?: string
}

export const bkService = {
  // Stats
  async getStats(): Promise<BKStats> {
    const response = await api.get('/bk/stats')
    return response.data
  },

  // Student BK Profile
  async getStudentBKProfile(studentId: number): Promise<StudentBKProfile> {
    const response = await api.get(`/bk/students/${studentId}/profile`)
    return response.data
  },

  async getStudentAchievementPoints(studentId: number): Promise<number> {
    const response = await api.get(`/bk/students/${studentId}/achievement-points`)
    return response.data.totalPoints
  },

  // Violations
  async getViolations(params?: BKPaginationParams): Promise<ViolationListResponse> {
    const response = await api.get('/bk/violations', { params })
    return response.data
  },

  async getViolation(id: number): Promise<Violation> {
    const response = await api.get(`/bk/violations/${id}`)
    return response.data
  },

  async getStudentViolations(studentId: number, params?: BKPaginationParams): Promise<ViolationListResponse> {
    const response = await api.get(`/bk/students/${studentId}/violations`, { params })
    return response.data
  },

  async createViolation(data: CreateViolationRequest): Promise<Violation> {
    const response = await api.post('/bk/violations', data)
    return response.data
  },

  async deleteViolation(id: number): Promise<void> {
    await api.delete(`/bk/violations/${id}`)
  },

  // Achievements
  async getAchievements(params?: BKPaginationParams): Promise<AchievementListResponse> {
    const response = await api.get('/bk/achievements', { params })
    return response.data
  },

  async getAchievement(id: number): Promise<Achievement> {
    const response = await api.get(`/bk/achievements/${id}`)
    return response.data
  },

  async getStudentAchievements(studentId: number, params?: BKPaginationParams): Promise<AchievementListResponse> {
    const response = await api.get(`/bk/students/${studentId}/achievements`, { params })
    return response.data
  },

  async createAchievement(data: CreateAchievementRequest): Promise<Achievement> {
    const response = await api.post('/bk/achievements', data)
    return response.data
  },

  async deleteAchievement(id: number): Promise<void> {
    await api.delete(`/bk/achievements/${id}`)
  },

  // Permits
  async getPermits(params?: BKPaginationParams): Promise<PermitListResponse> {
    const response = await api.get('/bk/permits', { params })
    return response.data
  },

  async getPermit(id: number): Promise<Permit> {
    const response = await api.get(`/bk/permits/${id}`)
    return response.data
  },

  async getStudentPermits(studentId: number, params?: BKPaginationParams): Promise<PermitListResponse> {
    const response = await api.get(`/bk/students/${studentId}/permits`, { params })
    return response.data
  },

  async createPermit(data: CreatePermitRequest): Promise<Permit> {
    const response = await api.post('/bk/permits', data)
    return response.data
  },

  async recordReturn(id: number, data: UpdatePermitReturnRequest): Promise<Permit> {
    const response = await api.put(`/bk/permits/${id}/return`, data)
    return response.data
  },

  async getPermitDocument(id: number): Promise<Blob> {
    const response = await api.get(`/bk/permits/${id}/document`, {
      responseType: 'blob',
    })
    return response.data
  },

  async deletePermit(id: number): Promise<void> {
    await api.delete(`/bk/permits/${id}`)
  },

  // Counseling Notes
  async getCounselingNotes(params?: BKPaginationParams): Promise<CounselingNoteListResponse> {
    const response = await api.get('/bk/counseling-notes', { params })
    return response.data
  },

  async getCounselingNote(id: number): Promise<CounselingNote> {
    const response = await api.get(`/bk/counseling-notes/${id}`)
    return response.data
  },

  async getStudentCounselingNotes(studentId: number, params?: BKPaginationParams): Promise<CounselingNoteListResponse> {
    const response = await api.get(`/bk/students/${studentId}/counseling-notes`, { params })
    return response.data
  },

  async createCounselingNote(data: CreateCounselingNoteRequest): Promise<CounselingNote> {
    const response = await api.post('/bk/counseling-notes', data)
    return response.data
  },

  async deleteCounselingNote(id: number): Promise<void> {
    await api.delete(`/bk/counseling-notes/${id}`)
  },
}
