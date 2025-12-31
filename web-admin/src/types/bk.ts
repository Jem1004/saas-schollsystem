// BK (Bimbingan Konseling) data types for Guru BK

export interface Violation {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentClass?: string
  category: string
  level: 'ringan' | 'sedang' | 'berat'
  description: string
  createdBy: number
  createdByName?: string
  createdAt: string
}

export interface Achievement {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentClass?: string
  title: string
  point: number
  description?: string
  createdBy: number
  createdByName?: string
  createdAt: string
}

export interface Permit {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  studentClass?: string
  reason: string
  exitTime: string
  returnTime?: string
  responsibleTeacherId: number
  responsibleTeacherName?: string
  documentUrl?: string
  createdBy: number
  createdByName?: string
  createdAt: string
}

export interface CounselingNote {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentClass?: string
  internalNote: string
  parentSummary?: string
  createdBy: number
  createdByName?: string
  createdAt: string
}

export interface StudentBKProfile {
  student: {
    id: number
    name: string
    nis: string
    nisn: string
    className: string
    classId: number
  }
  totalAchievementPoints: number
  violationCount: number
  achievementCount: number
  permitCount: number
  counselingCount: number
}

export interface BKStats {
  totalViolations: number
  totalAchievements: number
  totalPermits: number
  totalCounselingNotes: number
  recentViolations: Violation[]
  recentAchievements: Achievement[]
  studentsRequiringAttention: StudentBKProfile[]
}

// Request types
export interface CreateViolationRequest {
  studentId: number
  category: string
  level: 'ringan' | 'sedang' | 'berat'
  description: string
}

export interface CreateAchievementRequest {
  studentId: number
  title: string
  point: number
  description?: string
}

export interface CreatePermitRequest {
  studentId: number
  reason: string
  exitTime: string
  responsibleTeacherId: number
}

export interface UpdatePermitReturnRequest {
  returnTime: string
}

export interface CreateCounselingNoteRequest {
  studentId: number
  internalNote: string
  parentSummary?: string
}

// Response types
export interface ViolationListResponse {
  data: Violation[]
  total: number
  page: number
  pageSize: number
}

export interface AchievementListResponse {
  data: Achievement[]
  total: number
  page: number
  pageSize: number
}

export interface PermitListResponse {
  data: Permit[]
  total: number
  page: number
  pageSize: number
}

export interface CounselingNoteListResponse {
  data: CounselingNote[]
  total: number
  page: number
  pageSize: number
}

// Violation categories
export const VIOLATION_CATEGORIES = [
  'Keterlambatan',
  'Bolos',
  'Seragam',
  'Perilaku',
  'Kekerasan',
  'Bullying',
  'Merokok',
  'Narkoba',
  'Pencurian',
  'Vandalisme',
  'Lainnya',
]

export const VIOLATION_LEVELS = [
  { value: 'ringan', label: 'Ringan', color: 'warning' },
  { value: 'sedang', label: 'Sedang', color: 'orange' },
  { value: 'berat', label: 'Berat', color: 'error' },
]
