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
  studentName?: string
  studentNis?: string
  studentClass?: string
  teacherId: number
  teacherName?: string
  content: string
  createdAt: string
  updatedAt: string
}

export interface StudentAttendance {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: 'present' | 'absent' | 'late' | 'excused'
  method: 'rfid' | 'manual'
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
  studentCount: number
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
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: 'present' | 'absent' | 'late' | 'excused'
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
  total: number
  page: number
  pageSize: number
}

export interface HomeroomNoteListResponse {
  data: HomeroomNote[]
  total: number
  page: number
  pageSize: number
}

export interface StudentAttendanceListResponse {
  data: StudentAttendance[]
  total: number
  page: number
  pageSize: number
}

export interface ClassStudentListResponse {
  data: ClassStudent[]
  total: number
}
