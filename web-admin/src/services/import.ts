import api from './api'

// Import result types based on backend DTOs
export interface ImportError {
  row: number
  field: string
  message: string
}

export interface ImportWarning {
  row: number
  field: string
  message: string
}

export interface ImportResult {
  total_rows: number
  success_count: number
  failed_count: number
  warning_count: number
  errors?: ImportError[]
  warnings?: ImportWarning[]
  students_without_class?: number
}

export interface BulkAssignClassRequest {
  student_ids: number[]
  class_id: number
}

export interface BulkAssignClassResponse {
  updated_count: number
  students: Array<{
    id: number
    nis: string
    nisn: string
    name: string
    class_id: number
    class_name: string
    is_active: boolean
  }>
}

export const importService = {
  /**
   * Download student import template
   * Requirements: 1.1 - Template with columns: NIS, NISN, Nama, Kelas
   */
  async downloadStudentTemplate(): Promise<Blob> {
    const response = await api.get('/school/import/template/students', {
      responseType: 'blob',
    })
    return response.data
  },

  /**
   * Download parent import template
   * Requirements: 1.2 - Template with columns: Nama, No_HP, Email
   */
  async downloadParentTemplate(): Promise<Blob> {
    const response = await api.get('/school/import/template/parents', {
      responseType: 'blob',
    })
    return response.data
  },

  /**
   * Import students from Excel file
   * Requirements: 3.1-3.11 - Student import with validation and class matching
   */
  async importStudents(file: File): Promise<ImportResult> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post<{ success: boolean; data: ImportResult; message: string }>(
      '/school/import/students',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
    return response.data.data
  },

  /**
   * Import parents from Excel file
   * Requirements: 4.1-4.7 - Parent import with validation
   */
  async importParents(file: File): Promise<ImportResult> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await api.post<{ success: boolean; data: ImportResult; message: string }>(
      '/school/import/parents',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
    return response.data.data
  },

  /**
   * Get students without class assignment
   * Requirements: 6.1 - Filter for students without class
   */
  async getStudentsWithoutClass(): Promise<Array<{
    id: number
    nis: string
    nisn: string
    name: string
    is_active: boolean
  }>> {
    const response = await api.get<{ success: boolean; data: Array<{
      id: number
      nis: string
      nisn: string
      name: string
      is_active: boolean
    }> }>('/school/students/without-class')
    return response.data.data
  },

  /**
   * Bulk assign class to multiple students
   * Requirements: 6.2, 6.3, 6.4 - Bulk class assignment
   */
  async bulkAssignClass(studentIds: number[], classId: number): Promise<BulkAssignClassResponse> {
    const response = await api.post<{ success: boolean; data: BulkAssignClassResponse; message: string }>(
      '/school/students/bulk-assign-class',
      {
        student_ids: studentIds,
        class_id: classId,
      }
    )
    return response.data.data
  },
}

export default importService
