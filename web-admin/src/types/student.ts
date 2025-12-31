// Student data types (for student-specific operations)

export interface Student {
  id: number
  schoolId: number
  classId: number
  className?: string
  nis: string
  nisn: string
  name: string
  rfidCode?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface StudentProfile {
  id: number
  schoolId: number
  classId: number
  className: string
  nis: string
  nisn: string
  name: string
  rfidCode?: string
  isActive: boolean
  
  // Summary statistics
  attendanceSummary: {
    totalDays: number
    presentDays: number
    absentDays: number
    lateDays: number
    attendanceRate: number
  }
  
  academicSummary: {
    totalGrades: number
    averageScore: number
  }
  
  bkSummary: {
    totalViolations: number
    totalAchievements: number
    achievementPoints: number
    totalPermits: number
    totalCounseling: number
  }
  
  createdAt: string
  updatedAt: string
}

export interface StudentSearchResult {
  id: number
  nis: string
  nisn: string
  name: string
  className: string
  classId: number
  isActive: boolean
}

// Request types
export interface CreateStudentRequest {
  classId: number
  nis: string
  nisn: string
  name: string
  rfidCode?: string
}

export interface UpdateStudentRequest {
  classId?: number
  nis?: string
  name?: string
  rfidCode?: string
  isActive?: boolean
}

export interface ImportStudentsRequest {
  classId: number
  students: {
    nis: string
    nisn: string
    name: string
    rfidCode?: string
  }[]
}

export interface GetStudentsParams {
  page?: number
  pageSize?: number
  classId?: number
  search?: string
  isActive?: boolean
}

// Response types
export interface StudentListResponse {
  data: Student[]
  total: number
  page: number
  pageSize: number
}

export interface StudentSearchResponse {
  data: StudentSearchResult[]
  total: number
}

export interface ImportStudentsResponse {
  imported: number
  failed: number
  errors: {
    row: number
    nisn: string
    error: string
  }[]
}
