// Grade data types

export interface Grade {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentClass?: string
  title: string
  score: number
  description?: string
  createdBy: number
  createdByName?: string
  createdAt: string
  updatedAt: string
}

export interface GradeSummary {
  studentId: number
  studentName: string
  studentNis: string
  className: string
  totalGrades: number
  averageScore: number
  highestScore: number
  lowestScore: number
}

// Request types
export interface CreateGradeRequest {
  studentId: number
  title: string
  score: number
  description?: string
}

export interface UpdateGradeRequest {
  title?: string
  score?: number
  description?: string
}

export interface BatchGradeRequest {
  title: string
  description?: string
  grades: {
    studentId: number
    score: number
  }[]
}

export interface GetGradesParams {
  page?: number
  pageSize?: number
  studentId?: number
  classId?: number
  search?: string
  startDate?: string
  endDate?: string
}

// Response types
export interface GradeListResponse {
  data: Grade[]
  total: number
  page: number
  pageSize: number
}

export interface StudentGradesResponse {
  student: {
    id: number
    name: string
    nis: string
    className: string
  }
  summary: GradeSummary
  grades: Grade[]
}
