import api from './api'
import type {
  Grade,
  HomeroomNote,
  StudentAttendance,
  HomeroomStats,
  HomeroomClass,
  CreateGradeRequest,
  UpdateGradeRequest,
  CreateHomeroomNoteRequest,
  UpdateHomeroomNoteRequest,
  ManualAttendanceRequest,
  BatchGradeRequest,
  GradeListResponse,
  HomeroomNoteListResponse,
  StudentAttendanceListResponse,
  ClassStudentListResponse,
} from '@/types/homeroom'

export interface HomeroomPaginationParams {
  page?: number
  pageSize?: number
  studentId?: number
  search?: string
  date?: string
  startDate?: string
  endDate?: string
}

export const homeroomService = {
  // Stats & Class Info
  async getStats(): Promise<HomeroomStats> {
    const response = await api.get('/homeroom/stats')
    return response.data
  },

  async getMyClass(): Promise<HomeroomClass> {
    const response = await api.get('/homeroom/my-class')
    return response.data
  },

  async getClassStudents(params?: HomeroomPaginationParams): Promise<ClassStudentListResponse> {
    const response = await api.get('/homeroom/students', { params })
    return response.data
  },

  // Grades
  async getGrades(params?: HomeroomPaginationParams): Promise<GradeListResponse> {
    const response = await api.get('/homeroom/grades', { params })
    return response.data
  },

  async getGrade(id: number): Promise<Grade> {
    const response = await api.get(`/homeroom/grades/${id}`)
    return response.data
  },

  async getStudentGrades(studentId: number, params?: HomeroomPaginationParams): Promise<GradeListResponse> {
    const response = await api.get(`/homeroom/students/${studentId}/grades`, { params })
    return response.data
  },

  async createGrade(data: CreateGradeRequest): Promise<Grade> {
    const response = await api.post('/homeroom/grades', data)
    return response.data
  },

  async createBatchGrades(data: BatchGradeRequest): Promise<Grade[]> {
    const response = await api.post('/homeroom/grades/batch', data)
    return response.data
  },

  async updateGrade(id: number, data: UpdateGradeRequest): Promise<Grade> {
    const response = await api.put(`/homeroom/grades/${id}`, data)
    return response.data
  },

  async deleteGrade(id: number): Promise<void> {
    await api.delete(`/homeroom/grades/${id}`)
  },

  // Homeroom Notes
  async getNotes(params?: HomeroomPaginationParams): Promise<HomeroomNoteListResponse> {
    const response = await api.get('/homeroom/notes', { params })
    return response.data
  },

  async getNote(id: number): Promise<HomeroomNote> {
    const response = await api.get(`/homeroom/notes/${id}`)
    return response.data
  },

  async getStudentNotes(studentId: number, params?: HomeroomPaginationParams): Promise<HomeroomNoteListResponse> {
    const response = await api.get(`/homeroom/students/${studentId}/notes`, { params })
    return response.data
  },

  async createNote(data: CreateHomeroomNoteRequest): Promise<HomeroomNote> {
    const response = await api.post('/homeroom/notes', data)
    return response.data
  },

  async updateNote(id: number, data: UpdateHomeroomNoteRequest): Promise<HomeroomNote> {
    const response = await api.put(`/homeroom/notes/${id}`, data)
    return response.data
  },

  async deleteNote(id: number): Promise<void> {
    await api.delete(`/homeroom/notes/${id}`)
  },

  // Attendance
  async getClassAttendance(date: string): Promise<StudentAttendanceListResponse> {
    const response = await api.get('/homeroom/attendance', { params: { date } })
    return response.data
  },

  async getStudentAttendance(studentId: number, params?: HomeroomPaginationParams): Promise<StudentAttendanceListResponse> {
    const response = await api.get(`/homeroom/students/${studentId}/attendance`, { params })
    return response.data
  },

  async recordManualAttendance(data: ManualAttendanceRequest): Promise<StudentAttendance> {
    const response = await api.post('/homeroom/attendance/manual', data)
    return response.data
  },

  async updateAttendance(id: number, data: Partial<ManualAttendanceRequest>): Promise<StudentAttendance> {
    const response = await api.put(`/homeroom/attendance/${id}`, data)
    return response.data
  },
}
