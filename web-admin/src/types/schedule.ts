// Attendance Schedule types
// Requirements: 3.1 - Schedule management types

export interface AttendanceSchedule {
  id: number
  schoolId: number
  name: string
  startTime: string // Format: HH:MM
  endTime: string // Format: HH:MM
  lateThreshold: number // minutes after start_time
  veryLateThreshold?: number // optional, minutes after start_time
  daysOfWeek: string // e.g., "1,2,3,4,5" (Mon-Fri)
  isActive: boolean
  isDefault: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateScheduleRequest {
  name: string
  start_time: string // Format: HH:MM
  end_time: string // Format: HH:MM
  late_threshold: number
  very_late_threshold?: number
  days_of_week?: string
  is_active?: boolean
}

export interface UpdateScheduleRequest {
  name?: string
  start_time?: string
  end_time?: string
  late_threshold?: number
  very_late_threshold?: number
  days_of_week?: string
  is_active?: boolean
}

export interface ScheduleListResponse {
  schedules: AttendanceSchedule[]
  total: number
}

export interface ActiveScheduleResponse {
  schedule?: AttendanceSchedule
  message?: string
}

// Days of week configuration
export const DAYS_OF_WEEK = [
  { value: '1', label: 'Senin' },
  { value: '2', label: 'Selasa' },
  { value: '3', label: 'Rabu' },
  { value: '4', label: 'Kamis' },
  { value: '5', label: 'Jumat' },
  { value: '6', label: 'Sabtu' },
  { value: '7', label: 'Minggu' },
] as const

// Helper function to parse days of week string to array
export function parseDaysOfWeek(daysString: string): string[] {
  if (!daysString) return []
  return daysString.split(',').filter(d => d.trim())
}

// Helper function to format days of week array to string
export function formatDaysOfWeek(days: string[]): string {
  return days.join(',')
}

// Helper function to get day labels from days string
export function getDayLabels(daysString: string): string {
  const days = parseDaysOfWeek(daysString)
  return days
    .map(d => DAYS_OF_WEEK.find(day => day.value === d)?.label || d)
    .join(', ')
}
