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
  studentIds?: number[]
  studentNames?: string[]
  createdAt: string
  updatedAt: string
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
  assignedClassId?: number
  assignedClassName?: string
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
  excused: number
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

// Request types
export interface CreateClassRequest {
  name: string
  grade: number
  year: string
  homeroomTeacherId?: number
}

export interface UpdateClassRequest {
  name?: string
  grade?: number
  year?: string
  homeroomTeacherId?: number
}

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

export interface CreateParentRequest {
  name: string
  phone?: string
  email?: string
  password: string
  studentIds: number[]
}

export interface UpdateParentRequest {
  name?: string
  phone?: string
  email?: string
  studentIds?: number[]
}

export interface CreateUserRequest {
  role: 'guru' | 'wali_kelas' | 'guru_bk'
  username: string
  email?: string
  name?: string
  password: string
  assignedClassId?: number
}

export interface UpdateUserRequest {
  email?: string
  name?: string
  isActive?: boolean
  assignedClassId?: number
}

export interface UpdateSchoolSettingsRequest {
  attendanceStartTime?: string
  attendanceEndTime?: string
  attendanceLateThreshold?: number
  attendanceVeryLateThreshold?: number
  enableAttendanceNotification?: boolean
  enableGradeNotification?: boolean
  enableBKNotification?: boolean
  enableHomeroomNotification?: boolean
  academicYear?: string
  semester?: number
}

// Response types
export interface ClassListResponse {
  data: Class[]
  total: number
  page: number
  pageSize: number
}

export interface StudentListResponse {
  data: Student[]
  total: number
  page: number
  pageSize: number
}

export interface ParentListResponse {
  data: Parent[]
  total: number
  page: number
  pageSize: number
}

export interface UserListResponse {
  data: SchoolUser[]
  total: number
  page: number
  pageSize: number
}

export interface AttendanceListResponse {
  data: AttendanceSummary[]
  total: number
  page: number
  pageSize: number
}
