// Attendance data types

export interface Attendance {
  id: number
  student_id: number
  student_name?: string
  student_nis?: string
  student_nisn?: string
  class_name?: string
  schedule_id?: number
  schedule_name?: string // Requirements: 3.10 - Show which schedule the attendance belongs to
  date: string
  check_in_time?: string
  check_out_time?: string
  status: AttendanceStatus
  method: AttendanceMethod
  created_at: string
  updated_at: string
}

// Alias for camelCase usage in components
export interface AttendanceCamelCase {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentNisn?: string
  className?: string
  scheduleId?: number
  scheduleName?: string
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: AttendanceStatus
  method: AttendanceMethod
  createdAt: string
  updatedAt: string
}

export type AttendanceStatus = 'present' | 'absent' | 'late' | 'very_late' | 'sick' | 'excused'
export type AttendanceMethod = 'rfid' | 'manual'

export interface AttendanceRecord {
  id: number
  studentId: number
  date: string
  checkInTime?: string
  checkOutTime?: string
  method: AttendanceMethod
  createdAt: string
  updatedAt: string
}

export interface DailyAttendanceSummary {
  date: string
  totalStudents: number
  present: number
  absent: number
  late: number
  veryLate: number
  excused: number
  attendanceRate: number
}

export interface ClassAttendanceSummary {
  classId: number
  className: string
  date: string
  totalStudents: number
  present: number
  absent: number
  late: number
  veryLate: number
  excused: number
}

export interface StudentAttendanceHistory {
  studentId: number
  studentName: string
  studentNis: string
  className: string
  totalDays: number
  presentDays: number
  absentDays: number
  lateDays: number
  excusedDays: number
  attendanceRate: number
  records: Attendance[]
}

// Request types
export interface RecordRFIDAttendanceRequest {
  deviceCode: string
  rfidCode: string
  timestamp: string
}

export interface RecordManualAttendanceRequest {
  studentId: number
  date: string
  checkInTime?: string
  checkOutTime?: string
  status?: AttendanceStatus
}

export interface UpdateAttendanceRequest {
  checkInTime?: string
  checkOutTime?: string
  status?: AttendanceStatus
}

export interface GetAttendanceParams {
  page?: number
  pageSize?: number
  studentId?: number
  classId?: number
  date?: string
  startDate?: string
  endDate?: string
  status?: AttendanceStatus
}

// Response types
export interface AttendanceListResponse {
  attendances: Attendance[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface AttendanceSummaryResponse {
  summary: DailyAttendanceSummary
  byClass: ClassAttendanceSummary[]
}

export interface StudentAttendanceResponse {
  student: {
    id: number
    name: string
    nis: string
    className: string
  }
  summary: {
    totalDays: number
    presentDays: number
    absentDays: number
    lateDays: number
    excusedDays: number
    attendanceRate: number
  }
  records: Attendance[]
}

// Attendance status labels and colors
export const ATTENDANCE_STATUS_CONFIG = {
  present: { label: 'Hadir', color: 'success' },
  absent: { label: 'Tidak Hadir', color: 'error' },
  late: { label: 'Terlambat', color: 'warning' },
  very_late: { label: 'Sangat Terlambat', color: 'orange' },
  sick: { label: 'Sakit', color: 'blue' },
  excused: { label: 'Izin', color: 'cyan' },
} as const
