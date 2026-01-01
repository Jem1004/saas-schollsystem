import api from './api'
import type {
  SchoolSettings,
  AttendanceTimeWindow,
  UpdateSchoolSettingsRequest,
} from '@/types/settings'

// Backend response types for Settings (snake_case)
interface SettingsApiResponse {
  id: number
  school_id: number
  attendance_start_time: string
  attendance_end_time: string
  attendance_late_threshold: number
  attendance_very_late_threshold: number
  enable_attendance_notification: boolean
  enable_grade_notification: boolean
  enable_bk_notification: boolean
  enable_homeroom_notification: boolean
  academic_year: string
  semester: number
  created_at: string
  updated_at: string
}

// Backend response types for AttendanceTimeWindow (snake_case)
interface AttendanceTimeWindowApiResponse {
  start_time: string
  end_time: string
  late_time: string
  very_late_time: string
}

// API response wrapper
interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

// Transform settings from API response to frontend format
function transformSettings(apiSettings: SettingsApiResponse): SchoolSettings {
  return {
    id: apiSettings.id,
    schoolId: apiSettings.school_id,
    attendanceStartTime: apiSettings.attendance_start_time,
    attendanceEndTime: apiSettings.attendance_end_time,
    attendanceLateThreshold: apiSettings.attendance_late_threshold,
    attendanceVeryLateThreshold: apiSettings.attendance_very_late_threshold,
    enableAttendanceNotification: apiSettings.enable_attendance_notification,
    enableGradeNotification: apiSettings.enable_grade_notification,
    enableBKNotification: apiSettings.enable_bk_notification,
    enableHomeroomNotification: apiSettings.enable_homeroom_notification,
    academicYear: apiSettings.academic_year,
    semester: apiSettings.semester,
    createdAt: apiSettings.created_at,
    updatedAt: apiSettings.updated_at,
  }
}

// Transform attendance time window from API response to frontend format
function transformAttendanceTimeWindow(apiResponse: AttendanceTimeWindowApiResponse): AttendanceTimeWindow {
  return {
    startTime: apiResponse.start_time,
    endTime: apiResponse.end_time,
    lateTime: apiResponse.late_time,
    veryLateTime: apiResponse.very_late_time,
  }
}

// Transform request from frontend format to backend format (camelCase to snake_case)
function transformUpdateRequest(data: UpdateSchoolSettingsRequest): Record<string, unknown> {
  return {
    attendance_start_time: data.attendanceStartTime,
    attendance_end_time: data.attendanceEndTime,
    attendance_late_threshold: data.attendanceLateThreshold,
    attendance_very_late_threshold: data.attendanceVeryLateThreshold,
    enable_attendance_notification: data.enableAttendanceNotification,
    enable_grade_notification: data.enableGradeNotification,
    enable_bk_notification: data.enableBKNotification,
    enable_homeroom_notification: data.enableHomeroomNotification,
    academic_year: data.academicYear,
    semester: data.semester,
  }
}

export const settingsService = {
  // Get school settings
  async getSettings(): Promise<SchoolSettings> {
    const response = await api.get<ApiResponse<SettingsApiResponse>>('/settings')
    return transformSettings(response.data.data)
  },

  // Update school settings
  async updateSettings(data: UpdateSchoolSettingsRequest): Promise<SchoolSettings> {
    const response = await api.put<ApiResponse<SettingsApiResponse>>('/settings', transformUpdateRequest(data))
    return transformSettings(response.data.data)
  },

  // Reset settings to defaults
  async resetSettings(): Promise<SchoolSettings> {
    const response = await api.post<ApiResponse<SettingsApiResponse>>('/settings/reset')
    return transformSettings(response.data.data)
  },

  // Get attendance time window for a specific date
  async getAttendanceTimeWindow(date?: string): Promise<AttendanceTimeWindow> {
    const response = await api.get<ApiResponse<AttendanceTimeWindowApiResponse>>('/settings/attendance-window', {
      params: { date },
    })
    return transformAttendanceTimeWindow(response.data.data)
  },

  // Update attendance settings only
  async updateAttendanceSettings(data: {
    attendanceStartTime?: string
    attendanceEndTime?: string
    attendanceLateThreshold?: number
    attendanceVeryLateThreshold?: number
  }): Promise<SchoolSettings> {
    const response = await api.put<ApiResponse<SettingsApiResponse>>('/settings/attendance', {
      attendance_start_time: data.attendanceStartTime,
      attendance_end_time: data.attendanceEndTime,
      attendance_late_threshold: data.attendanceLateThreshold,
      attendance_very_late_threshold: data.attendanceVeryLateThreshold,
    })
    return transformSettings(response.data.data)
  },

  // Update notification settings only
  async updateNotificationSettings(data: {
    enableAttendanceNotification?: boolean
    enableGradeNotification?: boolean
    enableBKNotification?: boolean
    enableHomeroomNotification?: boolean
  }): Promise<SchoolSettings> {
    const response = await api.put<ApiResponse<SettingsApiResponse>>('/settings/notifications', {
      enable_attendance_notification: data.enableAttendanceNotification,
      enable_grade_notification: data.enableGradeNotification,
      enable_bk_notification: data.enableBKNotification,
      enable_homeroom_notification: data.enableHomeroomNotification,
    })
    return transformSettings(response.data.data)
  },

  // Update academic settings only
  async updateAcademicSettings(data: {
    academicYear?: string
    semester?: number
  }): Promise<SchoolSettings> {
    const response = await api.put<ApiResponse<SettingsApiResponse>>('/settings/academic', {
      academic_year: data.academicYear,
      semester: data.semester,
    })
    return transformSettings(response.data.data)
  },

  // Check if notification is enabled for a specific type
  async isNotificationEnabled(type: 'attendance' | 'grade' | 'bk' | 'homeroom'): Promise<boolean> {
    const settings = await this.getSettings()
    switch (type) {
      case 'attendance':
        return settings.enableAttendanceNotification
      case 'grade':
        return settings.enableGradeNotification
      case 'bk':
        return settings.enableBKNotification
      case 'homeroom':
        return settings.enableHomeroomNotification
      default:
        return false
    }
  },
}

export default settingsService
