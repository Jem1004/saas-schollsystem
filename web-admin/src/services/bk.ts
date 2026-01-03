import api from './api'
import type {
  Violation,
  ViolationCategory,
  Achievement,
  Permit,
  PermitDocumentData,
  CounselingNote,
  StudentBKProfile,
  BKStats,
  CreateViolationRequest,
  CreateViolationCategoryRequest,
  UpdateViolationCategoryRequest,
  CreateAchievementRequest,
  CreatePermitRequest,
  UpdatePermitReturnRequest,
  CreateCounselingNoteRequest,
  ViolationListResponse,
  ViolationCategoryListResponse,
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

// Transform API params to backend format
function transformParams(params?: BKPaginationParams): Record<string, unknown> {
  if (!params) return {}
  return {
    page: params.page,
    page_size: params.pageSize,
    student_id: params.studentId,
    search: params.search,
    start_date: params.startDate,
    end_date: params.endDate,
    level: params.level,
    category: params.category,
  }
}

// Transform backend violation to frontend format
function transformViolation(v: Record<string, unknown>): Violation {
  return {
    id: v.id as number,
    studentId: v.student_id as number,
    studentName: v.student_name as string,
    studentNis: v.student_nis as string,
    studentNisn: v.student_nisn as string,
    studentClass: v.class_name as string,
    className: v.class_name as string,
    categoryId: v.category_id as number | undefined,
    category: v.category as string,
    level: v.level as 'ringan' | 'sedang' | 'berat',
    point: v.point as number || 0,
    description: v.description as string,
    createdBy: v.created_by as number,
    createdByName: v.creator_name as string,
    creatorName: v.creator_name as string,
    createdAt: v.created_at as string,
  }
}

// Transform backend violation category to frontend format
function transformViolationCategory(c: Record<string, unknown>): ViolationCategory {
  return {
    id: c.id as number,
    schoolId: c.school_id as number,
    name: c.name as string,
    defaultPoint: c.default_point as number,
    defaultLevel: c.default_level as 'ringan' | 'sedang' | 'berat',
    description: c.description as string,
    isActive: c.is_active as boolean,
    createdAt: c.created_at as string,
    updatedAt: c.updated_at as string,
  }
}

// Transform backend achievement to frontend format
function transformAchievement(a: Record<string, unknown>): Achievement {
  return {
    id: a.id as number,
    studentId: a.student_id as number,
    studentName: a.student_name as string,
    studentNis: a.student_nis as string,
    studentNisn: a.student_nisn as string,
    studentClass: a.class_name as string,
    className: a.class_name as string,
    title: a.title as string,
    point: a.point as number,
    description: a.description as string,
    createdBy: a.created_by as number,
    createdByName: a.creator_name as string,
    creatorName: a.creator_name as string,
    createdAt: a.created_at as string,
  }
}

// Transform backend permit to frontend format
function transformPermit(p: Record<string, unknown>): Permit {
  return {
    id: p.id as number,
    studentId: p.student_id as number,
    studentName: p.student_name as string,
    studentNis: p.student_nis as string,
    studentNisn: p.student_nisn as string,
    studentClass: p.class_name as string,
    className: p.class_name as string,
    reason: p.reason as string,
    exitTime: p.exit_time as string,
    returnTime: p.return_time as string | undefined,
    responsibleTeacher: p.responsible_teacher as number,
    responsibleTeacherId: p.responsible_teacher as number,
    responsibleTeacherName: p.teacher_name as string,
    teacherName: p.teacher_name as string,
    documentUrl: p.document_url as string,
    hasReturned: p.has_returned as boolean,
    createdBy: p.created_by as number,
    createdByName: p.creator_name as string,
    creatorName: p.creator_name as string,
    createdAt: p.created_at as string,
  }
}

// Transform backend counseling note to frontend format
function transformCounselingNote(n: Record<string, unknown>): CounselingNote {
  return {
    id: n.id as number,
    studentId: n.student_id as number,
    studentName: n.student_name as string,
    studentNis: n.student_nis as string,
    studentNisn: n.student_nisn as string,
    studentClass: n.class_name as string,
    className: n.class_name as string,
    internalNote: n.internal_note as string,
    parentSummary: n.parent_summary as string,
    createdBy: n.created_by as number,
    createdByName: n.creator_name as string,
    creatorName: n.creator_name as string,
    createdAt: n.created_at as string,
  }
}

// Transform backend student BK profile to frontend format
function transformStudentBKProfile(p: Record<string, unknown>): StudentBKProfile {
  return {
    studentId: p.student_id as number,
    studentName: p.student_name as string,
    studentNis: p.student_nis as string,
    studentNisn: p.student_nisn as string,
    className: p.class_name as string,
    totalPoints: p.total_points as number,
    totalAchievementPoints: p.total_points as number,
    violationCount: p.violation_count as number,
    achievementCount: p.achievement_count as number,
    permitCount: p.permit_count as number,
    counselingCount: p.counseling_count as number,
  }
}

// Transform backend BK stats to frontend format
function transformBKStats(s: Record<string, unknown>): BKStats {
  const recentViolations = (s.recent_violations as Record<string, unknown>[] || []).map(transformViolation)
  const recentAchievements = (s.recent_achievements as Record<string, unknown>[] || []).map(transformAchievement)
  const studentsNeedingAttention = s.students_needing_attention as { student_id: number; student_name: string; class_name: string; violation_count: number; reason: string }[] || []
  
  return {
    totalViolations: s.total_violations as number,
    totalAchievements: s.total_achievements as number,
    totalPermits: s.total_permits as number,
    activePermits: s.active_permits as number,
    totalCounseling: s.total_counseling as number,
    totalCounselingNotes: s.total_counseling as number,
    recentViolations,
    recentAchievements,
    studentsNeedingAttention: studentsNeedingAttention.map(item => ({
      studentId: item.student_id,
      studentName: item.student_name,
      className: item.class_name,
      violationCount: item.violation_count,
      reason: item.reason,
    })),
    studentsRequiringAttention: [],
  }
}

export const bkService = {
  // Stats/Dashboard
  async getStats(): Promise<BKStats> {
    const response = await api.get('/bk/dashboard')
    return transformBKStats(response.data.data || response.data)
  },

  // Student BK Profile
  async getStudentBKProfile(studentId: number): Promise<StudentBKProfile> {
    const response = await api.get(`/bk/students/${studentId}/profile`)
    return transformStudentBKProfile(response.data.data || response.data)
  },

  async getStudentAchievementPoints(studentId: number): Promise<number> {
    const response = await api.get(`/bk/achievements/student/${studentId}/points`)
    const data = response.data.data || response.data
    return data?.total_points || data?.totalPoints || 0
  },

  // Violations
  async getViolations(params?: BKPaginationParams): Promise<ViolationListResponse> {
    const response = await api.get('/bk/violations', { params: transformParams(params) })
    const data = response.data.data || response.data
    const violations = (data.violations || data.data || []).map(transformViolation)
    const pagination = data.pagination || {}
    return {
      data: violations,
      total: pagination.total || violations.length,
      page: pagination.page || 1,
      pageSize: pagination.page_size || violations.length,
    }
  },

  async getViolation(id: number): Promise<Violation> {
    const response = await api.get(`/bk/violations/${id}`)
    return transformViolation(response.data.data || response.data)
  },

  async getStudentViolations(studentId: number, params?: BKPaginationParams): Promise<ViolationListResponse> {
    const response = await api.get(`/bk/students/${studentId}/violations`, { params: transformParams(params) })
    const data = response.data.data || response.data
    // Backend returns array directly for student-specific endpoints
    const violations = Array.isArray(data) ? data.map(transformViolation) : (data.violations || data.data || []).map(transformViolation)
    return { data: violations, total: violations.length, page: 1, pageSize: violations.length }
  },

  async createViolation(data: CreateViolationRequest): Promise<Violation> {
    const response = await api.post('/bk/violations', {
      student_id: data.studentId,
      category_id: data.categoryId,
      category: data.category,
      level: data.level,
      point: data.point,
      description: data.description,
    })
    return transformViolation(response.data.data || response.data)
  },

  async deleteViolation(id: number): Promise<void> {
    await api.delete(`/bk/violations/${id}`)
  },

  async getStudentViolationPoints(studentId: number): Promise<number> {
    const response = await api.get(`/bk/students/${studentId}/violations/points`)
    const data = response.data.data || response.data
    return data?.total_points || 0
  },

  // Violation Categories
  async getViolationCategories(activeOnly = true): Promise<ViolationCategoryListResponse> {
    const response = await api.get('/bk/violation-categories', { params: { active_only: activeOnly } })
    const data = response.data.data || response.data
    const categories = (data.categories || []).map(transformViolationCategory)
    return { categories }
  },

  async createViolationCategory(data: CreateViolationCategoryRequest): Promise<ViolationCategory> {
    const response = await api.post('/bk/violation-categories', {
      name: data.name,
      default_point: data.defaultPoint,
      default_level: data.defaultLevel,
      description: data.description,
    })
    return transformViolationCategory(response.data.data || response.data)
  },

  async updateViolationCategory(id: number, data: UpdateViolationCategoryRequest): Promise<ViolationCategory> {
    const response = await api.put(`/bk/violation-categories/${id}`, {
      name: data.name,
      default_point: data.defaultPoint,
      default_level: data.defaultLevel,
      description: data.description,
      is_active: data.isActive,
    })
    return transformViolationCategory(response.data.data || response.data)
  },

  async deleteViolationCategory(id: number): Promise<void> {
    await api.delete(`/bk/violation-categories/${id}`)
  },

  async initializeDefaultCategories(): Promise<void> {
    await api.post('/bk/violation-categories/initialize')
  },

  // Achievements
  async getAchievements(params?: BKPaginationParams): Promise<AchievementListResponse> {
    const response = await api.get('/bk/achievements', { params: transformParams(params) })
    const data = response.data.data || response.data
    const achievements = (data.achievements || data.data || []).map(transformAchievement)
    const pagination = data.pagination || {}
    return {
      data: achievements,
      total: pagination.total || achievements.length,
      page: pagination.page || 1,
      pageSize: pagination.page_size || achievements.length,
    }
  },

  async getAchievement(id: number): Promise<Achievement> {
    const response = await api.get(`/bk/achievements/${id}`)
    return transformAchievement(response.data.data || response.data)
  },

  async getStudentAchievements(studentId: number, params?: BKPaginationParams): Promise<AchievementListResponse> {
    const response = await api.get(`/bk/students/${studentId}/achievements`, { params: transformParams(params) })
    const data = response.data.data || response.data
    // Backend returns array directly for student-specific endpoints
    const achievements = Array.isArray(data) ? data.map(transformAchievement) : (data.achievements || data.data || []).map(transformAchievement)
    return { data: achievements, total: achievements.length, page: 1, pageSize: achievements.length }
  },

  async createAchievement(data: CreateAchievementRequest): Promise<Achievement> {
    const response = await api.post('/bk/achievements', {
      student_id: data.studentId,
      title: data.title,
      point: data.point,
      description: data.description,
    })
    return transformAchievement(response.data.data || response.data)
  },

  async deleteAchievement(id: number): Promise<void> {
    await api.delete(`/bk/achievements/${id}`)
  },

  // Permits
  async getPermits(params?: BKPaginationParams): Promise<PermitListResponse> {
    const response = await api.get('/bk/permits', { params: transformParams(params) })
    const data = response.data.data || response.data
    const permits = (data.permits || data.data || []).map(transformPermit)
    const pagination = data.pagination || {}
    return {
      data: permits,
      total: pagination.total || permits.length,
      page: pagination.page || 1,
      pageSize: pagination.page_size || permits.length,
    }
  },

  async getPermit(id: number): Promise<Permit> {
    const response = await api.get(`/bk/permits/${id}`)
    return transformPermit(response.data.data || response.data)
  },

  async getStudentPermits(studentId: number, params?: BKPaginationParams): Promise<PermitListResponse> {
    const response = await api.get(`/bk/students/${studentId}/permits`, { params: transformParams(params) })
    const data = response.data.data || response.data
    // Backend returns array directly for student-specific endpoints
    const permits = Array.isArray(data) ? data.map(transformPermit) : (data.permits || data.data || []).map(transformPermit)
    return { data: permits, total: permits.length, page: 1, pageSize: permits.length }
  },

  async createPermit(data: CreatePermitRequest): Promise<Permit> {
    const response = await api.post('/bk/permits', {
      student_id: data.studentId,
      reason: data.reason,
      exit_time: data.exitTime,
      responsible_teacher: data.responsibleTeacherId,
    })
    return transformPermit(response.data.data || response.data)
  },

  async recordReturn(id: number, data: UpdatePermitReturnRequest): Promise<Permit> {
    const response = await api.post(`/bk/permits/${id}/return`, {
      return_time: data.returnTime,
    })
    return transformPermit(response.data.data || response.data)
  },

  async getPermitDocument(id: number): Promise<PermitDocumentData> {
    const response = await api.get(`/bk/permits/${id}/document`)
    return response.data.data || response.data
  },

  async deletePermit(id: number): Promise<void> {
    await api.delete(`/bk/permits/${id}`)
  },

  // Counseling Notes
  async getCounselingNotes(params?: BKPaginationParams): Promise<CounselingNoteListResponse> {
    const response = await api.get('/bk/counseling', { params: transformParams(params) })
    const data = response.data.data || response.data
    const notes = (data.notes || data.data || []).map(transformCounselingNote)
    const pagination = data.pagination || {}
    return {
      data: notes,
      total: pagination.total || notes.length,
      page: pagination.page || 1,
      pageSize: pagination.page_size || notes.length,
    }
  },

  async getCounselingNote(id: number): Promise<CounselingNote> {
    const response = await api.get(`/bk/counseling/${id}`)
    return transformCounselingNote(response.data.data || response.data)
  },

  async getStudentCounselingNotes(studentId: number, params?: BKPaginationParams): Promise<CounselingNoteListResponse> {
    const response = await api.get(`/bk/students/${studentId}/counseling`, { params: transformParams(params) })
    const data = response.data.data || response.data
    // Backend returns array directly for student-specific endpoints
    const notes = Array.isArray(data) ? data.map(transformCounselingNote) : (data.notes || data.data || []).map(transformCounselingNote)
    return { data: notes, total: notes.length, page: 1, pageSize: notes.length }
  },

  async createCounselingNote(data: CreateCounselingNoteRequest): Promise<CounselingNote> {
    const response = await api.post('/bk/counseling', {
      student_id: data.studentId,
      internal_note: data.internalNote,
      parent_summary: data.parentSummary,
    })
    return transformCounselingNote(response.data.data || response.data)
  },

  async deleteCounselingNote(id: number): Promise<void> {
    await api.delete(`/bk/counseling/${id}`)
  },
}
