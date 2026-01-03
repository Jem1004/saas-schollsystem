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
import type { AttendanceSchedule } from '@/types/schedule'

export interface HomeroomPaginationParams {
  page?: number
  pageSize?: number
  studentId?: number
  search?: string
  date?: string
  startDate?: string
  endDate?: string
}

// API response wrapper
interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

export const homeroomService = {
  // Stats & Class Info
  async getStats(): Promise<HomeroomStats> {
    const response = await api.get<ApiResponse<HomeroomStats>>('/homeroom/stats')
    return response.data.data
  },

  async getMyClass(): Promise<HomeroomClass> {
    const response = await api.get<ApiResponse<HomeroomClass>>('/homeroom/my-class')
    return response.data.data
  },

  async getClassStudents(params?: HomeroomPaginationParams): Promise<ClassStudentListResponse> {
    const response = await api.get<ApiResponse<ClassStudentListResponse>>('/homeroom/students', { params })
    // Handle nested data structure from backend
    const data = response.data.data
    return {
      data: data.data || [],
      pagination: data.pagination,
    }
  },

  // Grades
  async getGrades(params?: HomeroomPaginationParams): Promise<GradeListResponse> {
    const response = await api.get<ApiResponse<GradeListResponse>>('/homeroom/grades', { params })
    // Handle nested data structure from backend
    const data = response.data.data
    return {
      data: data.data || [],
      pagination: data.pagination,
    }
  },

  async getGrade(id: number): Promise<Grade> {
    const response = await api.get<ApiResponse<Grade>>(`/homeroom/grades/${id}`)
    return response.data.data
  },

  async getStudentGrades(studentId: number, params?: HomeroomPaginationParams): Promise<GradeListResponse> {
    const response = await api.get<ApiResponse<GradeListResponse>>(`/homeroom/students/${studentId}/grades`, { params })
    // Handle nested data structure from backend
    const data = response.data.data
    return {
      data: data.data || [],
      pagination: data.pagination,
    }
  },

  async createGrade(data: CreateGradeRequest): Promise<Grade> {
    const response = await api.post<ApiResponse<Grade>>('/homeroom/grades', data)
    return response.data.data
  },

  async createBatchGrades(data: BatchGradeRequest): Promise<Grade[]> {
    const response = await api.post<ApiResponse<Grade[]>>('/homeroom/grades/batch', data)
    return response.data.data
  },

  async updateGrade(id: number, data: UpdateGradeRequest): Promise<Grade> {
    const response = await api.put<ApiResponse<Grade>>(`/homeroom/grades/${id}`, data)
    return response.data.data
  },

  async deleteGrade(id: number): Promise<void> {
    await api.delete(`/homeroom/grades/${id}`)
  },

  // Homeroom Notes
  async getNotes(params?: HomeroomPaginationParams): Promise<HomeroomNoteListResponse> {
    const response = await api.get<ApiResponse<HomeroomNoteListResponse>>('/homeroom/notes', { params })
    return response.data.data
  },

  async getNote(id: number): Promise<HomeroomNote> {
    const response = await api.get<ApiResponse<HomeroomNote>>(`/homeroom/notes/${id}`)
    return response.data.data
  },

  async getStudentNotes(studentId: number, params?: HomeroomPaginationParams): Promise<HomeroomNoteListResponse> {
    const response = await api.get<ApiResponse<HomeroomNoteListResponse>>(`/homeroom/students/${studentId}/notes`, { params })
    return response.data.data
  },

  async createNote(data: CreateHomeroomNoteRequest): Promise<HomeroomNote> {
    const response = await api.post<ApiResponse<HomeroomNote>>('/homeroom/notes', data)
    return response.data.data
  },

  async updateNote(id: number, data: UpdateHomeroomNoteRequest): Promise<HomeroomNote> {
    const response = await api.put<ApiResponse<HomeroomNote>>(`/homeroom/notes/${id}`, data)
    return response.data.data
  },

  async deleteNote(id: number): Promise<void> {
    await api.delete(`/homeroom/notes/${id}`)
  },

  // Attendance
  async getClassAttendance(date: string): Promise<StudentAttendanceListResponse> {
    const response = await api.get<ApiResponse<StudentAttendanceListResponse>>('/homeroom/attendance', { params: { date } })
    // Handle nested data structure from backend
    const data = response.data.data
    return {
      data: data.data || [],
      pagination: data.pagination,
    }
  },

  async getStudentAttendance(studentId: number, params?: HomeroomPaginationParams): Promise<StudentAttendanceListResponse> {
    const response = await api.get<ApiResponse<StudentAttendanceListResponse>>(`/homeroom/students/${studentId}/attendance`, { params })
    return response.data.data
  },

  async recordManualAttendance(data: ManualAttendanceRequest): Promise<StudentAttendance> {
    const response = await api.post<ApiResponse<StudentAttendance>>('/homeroom/attendance/manual', data)
    return response.data.data
  },

  async updateAttendance(id: number, data: Partial<ManualAttendanceRequest>): Promise<StudentAttendance> {
    const response = await api.put<ApiResponse<StudentAttendance>>(`/homeroom/attendance/${id}`, data)
    return response.data.data
  },

  // Schedules
  async getActiveSchedules(date: string): Promise<AttendanceSchedule[]> {
    const response = await api.get<ApiResponse<AttendanceSchedule[]>>('/homeroom/schedules', { params: { date } })
    return response.data.data || []
  },
}
