import api from './api'
import type {
  Grade,
  GradeSummary,
  CreateGradeRequest,
  UpdateGradeRequest,
  BatchGradeRequest,
  GetGradesParams,
  GradeListResponse,
  StudentGradesResponse,
} from '@/types/grade'

export const gradeService = {
  // Get grades list
  async getGrades(params?: GetGradesParams): Promise<GradeListResponse> {
    const response = await api.get<GradeListResponse>('/grades', { params })
    return response.data
  },

  // Get single grade
  async getGrade(id: number): Promise<Grade> {
    const response = await api.get<Grade>(`/grades/${id}`)
    return response.data
  },

  // Get student grades
  async getStudentGrades(studentId: number, params?: GetGradesParams): Promise<StudentGradesResponse> {
    const response = await api.get<StudentGradesResponse>(`/grades/student/${studentId}`, { params })
    return response.data
  },

  // Get student grade summary
  async getStudentGradeSummary(studentId: number): Promise<GradeSummary> {
    const response = await api.get<GradeSummary>(`/grades/student/${studentId}/summary`)
    return response.data
  },

  // Get grades by class
  async getGradesByClass(classId: number, params?: GetGradesParams): Promise<GradeListResponse> {
    const response = await api.get<GradeListResponse>(`/grades/class/${classId}`, { params })
    return response.data
  },

  // Create grade
  async createGrade(data: CreateGradeRequest): Promise<Grade> {
    const response = await api.post<Grade>('/grades', data)
    return response.data
  },

  // Create batch grades
  async createBatchGrades(data: BatchGradeRequest): Promise<Grade[]> {
    const response = await api.post<Grade[]>('/grades/batch', data)
    return response.data
  },

  // Update grade
  async updateGrade(id: number, data: UpdateGradeRequest): Promise<Grade> {
    const response = await api.put<Grade>(`/grades/${id}`, data)
    return response.data
  },

  // Delete grade
  async deleteGrade(id: number): Promise<void> {
    await api.delete(`/grades/${id}`)
  },

  // Export grades
  async exportGrades(params: {
    classId?: number
    studentId?: number
    startDate?: string
    endDate?: string
    format?: 'xlsx' | 'csv' | 'pdf'
  }): Promise<Blob> {
    const response = await api.get('/grades/export', {
      params,
      responseType: 'blob',
    })
    return response.data
  },
}

export default gradeService
