// BK (Bimbingan Konseling) data types for Guru BK

export interface Violation {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  studentClass?: string
  className?: string
  categoryId?: number
  category: string
  level: 'ringan' | 'sedang' | 'berat'
  point: number
  description: string
  createdBy: number
  createdByName?: string
  creatorName?: string
  createdAt: string
}

export interface ViolationCategory {
  id: number
  schoolId: number
  name: string
  defaultPoint: number
  defaultLevel: 'ringan' | 'sedang' | 'berat'
  description?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Achievement {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  studentClass?: string
  className?: string
  title: string
  point: number
  description?: string
  createdBy: number
  createdByName?: string
  creatorName?: string
  createdAt: string
}

export interface Permit {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  studentClass?: string
  className?: string
  reason: string
  exitTime: string
  returnTime?: string
  responsibleTeacher?: number
  responsibleTeacherId?: number
  responsibleTeacherName?: string
  teacherName?: string
  documentUrl?: string
  hasReturned?: boolean
  createdBy: number
  createdByName?: string
  creatorName?: string
  createdAt: string
}

export interface PermitDocumentData {
  permit_id: number
  student_name: string
  student_nis: string
  student_nisn: string
  class_name: string
  school_name: string
  reason: string
  exit_time: string
  return_time?: string
  responsible_teacher: string
  generated_at: string
}

export interface CounselingNote {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  studentClass?: string
  className?: string
  internalNote: string
  parentSummary?: string
  createdBy: number
  createdByName?: string
  creatorName?: string
  createdAt: string
}

export interface StudentBKProfile {
  studentId: number
  studentName: string
  studentNis: string
  studentNisn: string
  className: string
  totalPoints: number
  totalAchievementPoints?: number
  violationCount: number
  achievementCount: number
  permitCount: number
  counselingCount: number
  // Legacy format support
  student?: {
    id: number
    name: string
    nis: string
    nisn: string
    className: string
    classId: number
  }
}

export interface StudentAttentionItem {
  studentId: number
  studentName: string
  className: string
  violationCount: number
  reason: string
}

export interface BKStats {
  totalViolations: number
  totalAchievements: number
  totalPermits: number
  activePermits?: number
  totalCounseling?: number
  totalCounselingNotes?: number
  recentViolations: Violation[]
  recentAchievements: Achievement[]
  studentsNeedingAttention?: StudentAttentionItem[]
  studentsRequiringAttention?: StudentBKProfile[]
}

// Request types
export interface CreateViolationRequest {
  studentId: number
  categoryId?: number
  category: string
  level: 'ringan' | 'sedang' | 'berat'
  point?: number
  description: string
}

export interface CreateViolationCategoryRequest {
  name: string
  defaultPoint: number
  defaultLevel: 'ringan' | 'sedang' | 'berat'
  description?: string
}

export interface UpdateViolationCategoryRequest {
  name?: string
  defaultPoint?: number
  defaultLevel?: 'ringan' | 'sedang' | 'berat'
  description?: string
  isActive?: boolean
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

export interface ViolationCategoryListResponse {
  categories: ViolationCategory[]
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

// Violation categories - now loaded from API, these are fallback defaults
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
  { value: 'ringan', label: 'Ringan', color: 'warning', defaultPoint: -5 },
  { value: 'sedang', label: 'Sedang', color: 'orange', defaultPoint: -15 },
  { value: 'berat', label: 'Berat', color: 'error', defaultPoint: -30 },
]
