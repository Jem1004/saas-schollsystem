// Settings data types

export interface SchoolSettings {
  id: number
  schoolId: number
  
  // Attendance Settings
  attendanceStartTime: string // HH:MM format
  attendanceEndTime: string // HH:MM format
  attendanceLateThreshold: number // Minutes after start to be considered late
  attendanceVeryLateThreshold: number // Minutes after start to be considered very late
  
  // Notification Settings
  enableAttendanceNotification: boolean
  enableGradeNotification: boolean
  enableBKNotification: boolean
  enableHomeroomNotification: boolean
  
  // General Settings
  academicYear: string // e.g., "2024/2025"
  semester: number // 1 or 2
  
  createdAt: string
  updatedAt: string
}

export interface AttendanceTimeWindow {
  startTime: string
  endTime: string
  lateTime: string
  veryLateTime: string
}

// Request types
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

// Default settings
export const DEFAULT_SCHOOL_SETTINGS: Omit<SchoolSettings, 'id' | 'schoolId' | 'createdAt' | 'updatedAt'> = {
  attendanceStartTime: '07:00',
  attendanceEndTime: '07:30',
  attendanceLateThreshold: 30,
  attendanceVeryLateThreshold: 60,
  enableAttendanceNotification: true,
  enableGradeNotification: true,
  enableBKNotification: true,
  enableHomeroomNotification: true,
  academicYear: '',
  semester: 1,
}

// Semester options
export const SEMESTER_OPTIONS = [
  { value: 1, label: 'Semester 1 (Ganjil)' },
  { value: 2, label: 'Semester 2 (Genap)' },
]
