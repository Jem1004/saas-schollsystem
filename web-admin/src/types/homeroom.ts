// Homeroom (Wali Kelas) data types

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

export interface HomeroomNote {
  id: number
  studentId: number
  student_id?: number // Backend uses snake_case
  studentName?: string
  student_name?: string
  studentNis?: string
  student_nis?: string
  studentClass?: string
  class_name?: string
  teacherId: number
  teacher_id?: number
  teacherName?: string
  teacher_name?: string
  content: string
  createdAt: string
  created_at?: string
  updatedAt: string
  updated_at?: string
}

export interface StudentAttendance {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: 'present' | 'absent' | 'late' | 'excused' | 'on_time' | 'very_late' | 'sick'
  method: 'rfid' | 'manual' | ''
  createdAt: string
  updatedAt: string
}

export interface ClassStudent {
  id: number
  nis: string
  nisn: string
  name: string
  rfidCode?: string
  isActive: boolean
}

export interface HomeroomStats {
  classId: number
  className: string
  totalStudents: number
  todayAttendance: {
    present: number
    absent: number
    late: number
    excused: number
  }
  recentGrades: Grade[]
  recentNotes: HomeroomNote[]
}

export interface HomeroomClass {
  id: number
  name: string
  grade: number
  year: string
  totalStudents?: number
  studentCount?: number
  homeroomTeacherId?: number
  homeroomTeacherName?: string
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

export interface CreateHomeroomNoteRequest {
  studentId: number
  content: string
}

export interface UpdateHomeroomNoteRequest {
  content: string
}

export interface ManualAttendanceRequest {
  studentId: number
  scheduleId: number
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: 'present' | 'absent' | 'late' | 'excused' | 'on_time' | 'very_late' | 'sick'
}

export interface BatchGradeRequest {
  title: string
  description?: string
  grades: {
    studentId: number
    score: number
  }[]
}

// Response types
export interface GradeListResponse {
  data: Grade[]
  pagination?: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
  total?: number
  page?: number
  pageSize?: number
}

export interface HomeroomNoteListResponse {
  notes?: HomeroomNote[]
  data?: HomeroomNote[]
  pagination?: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
  total?: number
  page?: number
  pageSize?: number
}

export interface StudentAttendanceListResponse {
  data: StudentAttendance[]
  pagination?: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
  total?: number
  page?: number
  pageSize?: number
}

export interface ClassStudentListResponse {
  data: ClassStudent[]
  pagination?: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
  total?: number
}
