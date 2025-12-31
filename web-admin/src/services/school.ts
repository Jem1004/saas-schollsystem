import api from './api'
import type {
  Class,
  Student,
  Parent,
  SchoolUser,
  SchoolSettings,
  SchoolStats,
  AttendanceSummary,
  CreateClassRequest,
  UpdateClassRequest,
  CreateStudentRequest,
  UpdateStudentRequest,
  CreateParentRequest,
  UpdateParentRequest,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateSchoolSettingsRequest,
  ClassListResponse,
  StudentListResponse,
  ParentListResponse,
  UserListResponse,
  AttendanceListResponse,
} from '@/types/school'

export interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
  classId?: number
  date?: string
  startDate?: string
  endDate?: string
}

export const schoolService = {
  // Stats
  async getStats(): Promise<SchoolStats> {
    const response = await api.get('/school/stats')
    return response.data
  },

  // Classes
  async getClasses(params?: PaginationParams): Promise<ClassListResponse> {
    const response = await api.get('/school/classes', { params })
    return response.data
  },

  async getClass(id: number): Promise<Class> {
    const response = await api.get(`/school/classes/${id}`)
    return response.data
  },

  async createClass(data: CreateClassRequest): Promise<Class> {
    const response = await api.post('/school/classes', data)
    return response.data
  },

  async updateClass(id: number, data: UpdateClassRequest): Promise<Class> {
    const response = await api.put(`/school/classes/${id}`, data)
    return response.data
  },

  async deleteClass(id: number): Promise<void> {
    await api.delete(`/school/classes/${id}`)
  },

  // Students
  async getStudents(params?: PaginationParams): Promise<StudentListResponse> {
    const response = await api.get('/school/students', { params })
    return response.data
  },

  async getStudent(id: number): Promise<Student> {
    const response = await api.get(`/school/students/${id}`)
    return response.data
  },

  async createStudent(data: CreateStudentRequest): Promise<Student> {
    const response = await api.post('/school/students', data)
    return response.data
  },

  async updateStudent(id: number, data: UpdateStudentRequest): Promise<Student> {
    const response = await api.put(`/school/students/${id}`, data)
    return response.data
  },

  async deleteStudent(id: number): Promise<void> {
    await api.delete(`/school/students/${id}`)
  },

  // Parents
  async getParents(params?: PaginationParams): Promise<ParentListResponse> {
    const response = await api.get('/school/parents', { params })
    return response.data
  },

  async getParent(id: number): Promise<Parent> {
    const response = await api.get(`/school/parents/${id}`)
    return response.data
  },

  async createParent(data: CreateParentRequest): Promise<Parent> {
    const response = await api.post('/school/parents', data)
    return response.data
  },

  async updateParent(id: number, data: UpdateParentRequest): Promise<Parent> {
    const response = await api.put(`/school/parents/${id}`, data)
    return response.data
  },

  async deleteParent(id: number): Promise<void> {
    await api.delete(`/school/parents/${id}`)
  },

  async linkParentToStudent(parentId: number, studentId: number): Promise<void> {
    await api.post(`/school/parents/${parentId}/students/${studentId}`)
  },

  async unlinkParentFromStudent(parentId: number, studentId: number): Promise<void> {
    await api.delete(`/school/parents/${parentId}/students/${studentId}`)
  },

  // Users
  async getUsers(params?: PaginationParams): Promise<UserListResponse> {
    const response = await api.get('/school/users', { params })
    return response.data
  },

  async getUser(id: number): Promise<SchoolUser> {
    const response = await api.get(`/school/users/${id}`)
    return response.data
  },

  async createUser(data: CreateUserRequest): Promise<SchoolUser> {
    const response = await api.post('/school/users', data)
    return response.data
  },

  async updateUser(id: number, data: UpdateUserRequest): Promise<SchoolUser> {
    const response = await api.put(`/school/users/${id}`, data)
    return response.data
  },

  async deleteUser(id: number): Promise<void> {
    await api.delete(`/school/users/${id}`)
  },

  async resetUserPassword(id: number): Promise<{ temporaryPassword: string }> {
    const response = await api.post(`/school/users/${id}/reset-password`)
    return response.data
  },

  // Attendance
  async getAttendanceSummary(params?: PaginationParams): Promise<AttendanceListResponse> {
    const response = await api.get('/school/attendance/summary', { params })
    return response.data
  },

  async getAttendanceByClass(classId: number, date: string): Promise<AttendanceSummary> {
    const response = await api.get(`/school/attendance/class/${classId}`, { params: { date } })
    return response.data
  },

  async exportAttendance(params: { startDate: string; endDate: string; classId?: number }): Promise<Blob> {
    const response = await api.get('/school/attendance/export', {
      params,
      responseType: 'blob',
    })
    return response.data
  },

  // Settings
  async getSettings(): Promise<SchoolSettings> {
    const response = await api.get('/school/settings')
    return response.data
  },

  async updateSettings(data: UpdateSchoolSettingsRequest): Promise<SchoolSettings> {
    const response = await api.put('/school/settings', data)
    return response.data
  },

  async resetSettings(): Promise<SchoolSettings> {
    const response = await api.post('/school/settings/reset')
    return response.data
  },
}
