// Attendance data types

export interface Attendance {
  id: number
  studentId: number
  studentName?: string
  studentNis?: string
  studentClass?: string
  scheduleId?: number
  scheduleName?: string // Requirements: 3.10 - Show which schedule the attendance belongs to
  date: string
  checkInTime?: string
  checkOutTime?: string
  status: AttendanceStatus
  method: AttendanceMethod
  createdAt: string
  updatedAt: string
}

export type AttendanceStatus = 'present' | 'absent' | 'late' | 'very_late' | 'excused'
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
  data: Attendance[]
  total: number
  page: number
  pageSize: number
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
  excused: { label: 'Izin', color: 'processing' },
} as const
