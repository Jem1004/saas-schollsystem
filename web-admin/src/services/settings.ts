import api from './api'
import type {
  SchoolSettings,
  AttendanceTimeWindow,
  UpdateSchoolSettingsRequest,
} from '@/types/settings'

export const settingsService = {
  // Get school settings
  async getSettings(): Promise<SchoolSettings> {
    const response = await api.get<SchoolSettings>('/settings')
    return response.data
  },

  // Update school settings
  async updateSettings(data: UpdateSchoolSettingsRequest): Promise<SchoolSettings> {
    const response = await api.put<SchoolSettings>('/settings', data)
    return response.data
  },

  // Reset settings to defaults
  async resetSettings(): Promise<SchoolSettings> {
    const response = await api.post<SchoolSettings>('/settings/reset')
    return response.data
  },

  // Get attendance time window for a specific date
  async getAttendanceTimeWindow(date?: string): Promise<AttendanceTimeWindow> {
    const response = await api.get<AttendanceTimeWindow>('/settings/attendance-window', {
      params: { date },
    })
    return response.data
  },

  // Update attendance settings only
  async updateAttendanceSettings(data: {
    attendanceStartTime?: string
    attendanceEndTime?: string
    attendanceLateThreshold?: number
    attendanceVeryLateThreshold?: number
  }): Promise<SchoolSettings> {
    const response = await api.put<SchoolSettings>('/settings/attendance', data)
    return response.data
  },

  // Update notification settings only
  async updateNotificationSettings(data: {
    enableAttendanceNotification?: boolean
    enableGradeNotification?: boolean
    enableBKNotification?: boolean
    enableHomeroomNotification?: boolean
  }): Promise<SchoolSettings> {
    const response = await api.put<SchoolSettings>('/settings/notifications', data)
    return response.data
  },

  // Update academic settings only
  async updateAcademicSettings(data: {
    academicYear?: string
    semester?: number
  }): Promise<SchoolSettings> {
    const response = await api.put<SchoolSettings>('/settings/academic', data)
    return response.data
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
