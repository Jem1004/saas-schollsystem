// School data types for Admin Sekolah

export interface Class {
  id: number
  schoolId: number
  name: string
  grade: number
  year: string
  homeroomTeacherId?: number
  homeroomTeacherName?: string
  studentCount?: number
  createdAt: string
  updatedAt: string
}

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
  hasAccount?: boolean
  username?: string
  temporaryPassword?: string
  createdAt: string
  updatedAt: string
}

export interface Parent {
  id: number
  schoolId: number
  userId: number
  name: string
  phone?: string
  email?: string
  username?: string
  temporaryPassword?: string
  studentIds?: number[]
  studentNames?: string[]
  createdAt: string
  updatedAt: string
}

export interface ResetPasswordResponse {
  username: string
  temporaryPassword: string
  message: string
}

export interface AssignedClassInfo {
  id: number
  name: string
}

export interface SchoolUser {
  id: number
  schoolId: number
  role: 'guru' | 'wali_kelas' | 'guru_bk'
  username: string
  email?: string
  name?: string
  isActive: boolean
  mustResetPwd: boolean
  assignedClassId?: number      // For wali_kelas
  assignedClassName?: string    // For wali_kelas
  assignedClasses?: AssignedClassInfo[]  // For guru_bk
  lastLoginAt?: string
  createdAt: string
  updatedAt: string
}

export interface AttendanceSummary {
  date: string
  classId?: number
  className?: string
  totalStudents: number
  present: number
  absent: number
  late: number
  sick: number
  excused: number
}

// Backend response for school attendance summary
export interface SchoolAttendanceSummaryResponse {
  school_id: number
  date: string
  summary: {
    date: string
    total_count: number
    present: number
    late: number
    very_late: number
    absent: number
    sick: number
    excused: number
  }
  by_class: {
    class_id: number
    class_name: string
    total_students: number
    present: number
    late: number
    absent: number
    sick: number
    excused: number
  }[]
}

export interface SchoolSettings {
  id: number
  schoolId: number
  attendanceStartTime: string
  attendanceEndTime: string
  attendanceLateThreshold: number
  attendanceVeryLateThreshold: number
  enableAttendanceNotification: boolean
  enableGradeNotification: boolean
  enableBKNotification: boolean
  enableHomeroomNotification: boolean
  academicYear: string
  semester: number
  createdAt: string
  updatedAt: string
}

export interface SchoolStats {
  totalStudents: number
  totalClasses: number
  totalTeachers: number
  totalParents: number
  todayAttendance: {
    present: number
    absent: number
    late: number
    total: number
  }
}

// Request types - using snake_case to match backend API
export interface CreateClassRequest {
  name: string
  grade: number
  year: string
  homeroom_teacher_id?: number
}

export interface UpdateClassRequest {
  name?: string
  grade?: number
  year?: string
  homeroom_teacher_id?: number
}

export interface CreateStudentRequest {
  class_id: number
  nis: string
  nisn: string
  name: string
  rfid_code?: string
  create_account?: boolean // If true, create user account for mobile login
}

export interface UpdateStudentRequest {
  class_id?: number
  nis?: string
  name?: string
  rfid_code?: string
  is_active?: boolean
}

export interface CreateParentRequest {
  name: string
  phone: string // Required, used as username
  email?: string
  password?: string // Optional, will auto-generate if empty
  student_ids: number[]
}

export interface UpdateParentRequest {
  name?: string
  phone?: string
  email?: string
  student_ids?: number[]
}

export interface CreateUserRequest {
  role: 'guru' | 'wali_kelas' | 'guru_bk'
  username: string
  email?: string
  name?: string
  password: string
  assigned_class_id?: number     // For wali_kelas
  assigned_class_ids?: number[]  // For guru_bk
}

export interface UpdateUserRequest {
  email?: string
  name?: string
  is_active?: boolean
  assigned_class_id?: number     // For wali_kelas
  assigned_class_ids?: number[]  // For guru_bk
}

export interface UpdateSchoolSettingsRequest {
  attendance_start_time?: string
  attendance_end_time?: string
  attendance_late_threshold?: number
  attendance_very_late_threshold?: number
  enable_attendance_notification?: boolean
  enable_grade_notification?: boolean
  enable_bk_notification?: boolean
  enable_homeroom_notification?: boolean
  academic_year?: string
  semester?: number
}

// Response types
export interface ClassListResponse {
  classes: Class[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface StudentListResponse {
  students: Student[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface ParentListResponse {
  parents: Parent[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface UserListResponse {
  users: SchoolUser[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface AttendanceListResponse {
  data: AttendanceSummary[]
  total: number
  page: number
  pageSize: number
}

// Device types
export interface Device {
  id: number
  schoolId: number
  deviceCode: string
  description: string
  isActive: boolean
  lastSeenAt?: string
  createdAt: string
  updatedAt: string
}

// Pairing types
export interface PairingSessionResponse {
  active: boolean
  studentId?: number
  studentName?: string
  deviceId?: number
  expiresAt?: string
  message: string
}
