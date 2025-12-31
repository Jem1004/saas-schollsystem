import api from './api'
import type {
  Student,
  StudentProfile,
  CreateStudentRequest,
  UpdateStudentRequest,
  ImportStudentsRequest,
  GetStudentsParams,
  StudentListResponse,
  StudentSearchResponse,
  ImportStudentsResponse,
} from '@/types/student'

export const studentService = {
  // Get students list
  async getStudents(params?: GetStudentsParams): Promise<StudentListResponse> {
    const response = await api.get<StudentListResponse>('/students', { params })
    return response.data
  },

  // Get single student
  async getStudent(id: number): Promise<Student> {
    const response = await api.get<Student>(`/students/${id}`)
    return response.data
  },

  // Get student profile with summary
  async getStudentProfile(id: number): Promise<StudentProfile> {
    const response = await api.get<StudentProfile>(`/students/${id}/profile`)
    return response.data
  },

  // Get student by NISN
  async getStudentByNISN(nisn: string): Promise<Student> {
    const response = await api.get<Student>(`/students/nisn/${nisn}`)
    return response.data
  },

  // Search students
  async searchStudents(query: string, classId?: number): Promise<StudentSearchResponse> {
    const response = await api.get<StudentSearchResponse>('/students/search', {
      params: { q: query, classId },
    })
    return response.data
  },

  // Create student
  async createStudent(data: CreateStudentRequest): Promise<Student> {
    const response = await api.post<Student>('/students', data)
    return response.data
  },

  // Update student
  async updateStudent(id: number, data: UpdateStudentRequest): Promise<Student> {
    const response = await api.put<Student>(`/students/${id}`, data)
    return response.data
  },

  // Delete student
  async deleteStudent(id: number): Promise<void> {
    await api.delete(`/students/${id}`)
  },

  // Import students (bulk)
  async importStudents(data: ImportStudentsRequest): Promise<ImportStudentsResponse> {
    const response = await api.post<ImportStudentsResponse>('/students/import', data)
    return response.data
  },

  // Get students by class
  async getStudentsByClass(classId: number, params?: GetStudentsParams): Promise<StudentListResponse> {
    const response = await api.get<StudentListResponse>(`/classes/${classId}/students`, { params })
    return response.data
  },

  // Assign RFID to student
  async assignRFID(studentId: number, rfidCode: string): Promise<Student> {
    const response = await api.put<Student>(`/students/${studentId}/rfid`, { rfidCode })
    return response.data
  },

  // Remove RFID from student
  async removeRFID(studentId: number): Promise<Student> {
    const response = await api.delete<Student>(`/students/${studentId}/rfid`)
    return response.data
  },
}

export default studentService
