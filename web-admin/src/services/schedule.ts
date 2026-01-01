import api from './api'
import type {
  AttendanceSchedule,
  CreateScheduleRequest,
  UpdateScheduleRequest,
  ScheduleListResponse,
  ActiveScheduleResponse,
} from '@/types/schedule'

// Transform backend response to frontend format (snake_case to camelCase)
function transformSchedule(data: Record<string, unknown>): AttendanceSchedule {
  return {
    id: data.id as number,
    schoolId: data.school_id as number,
    name: data.name as string,
    startTime: data.start_time as string,
    endTime: data.end_time as string,
    lateThreshold: data.late_threshold as number,
    veryLateThreshold: data.very_late_threshold as number | undefined,
    daysOfWeek: data.days_of_week as string,
    isActive: data.is_active as boolean,
    isDefault: data.is_default as boolean,
    createdAt: data.created_at as string,
    updatedAt: data.updated_at as string,
  }
}

export const scheduleService = {
  // Get all schedules for the school
  async getSchedules(): Promise<ScheduleListResponse> {
    const response = await api.get<{ schedules: Record<string, unknown>[]; total: number }>('/schedules')
    return {
      schedules: (response.data.schedules || []).map(transformSchedule),
      total: response.data.total || 0,
    }
  },

  // Get a single schedule by ID
  async getScheduleById(id: number): Promise<AttendanceSchedule> {
    const response = await api.get<Record<string, unknown>>(`/schedules/${id}`)
    return transformSchedule(response.data)
  },

  // Create a new schedule
  async createSchedule(data: CreateScheduleRequest): Promise<AttendanceSchedule> {
    const response = await api.post<Record<string, unknown>>('/schedules', data)
    return transformSchedule(response.data)
  },

  // Update an existing schedule
  async updateSchedule(id: number, data: UpdateScheduleRequest): Promise<AttendanceSchedule> {
    const response = await api.put<Record<string, unknown>>(`/schedules/${id}`, data)
    return transformSchedule(response.data)
  },

  // Delete a schedule
  async deleteSchedule(id: number): Promise<void> {
    await api.delete(`/schedules/${id}`)
  },

  // Set a schedule as default
  async setDefaultSchedule(id: number): Promise<AttendanceSchedule> {
    const response = await api.post<Record<string, unknown>>(`/schedules/${id}/default`)
    return transformSchedule(response.data)
  },

  // Get the currently active schedule
  async getActiveSchedule(): Promise<ActiveScheduleResponse> {
    const response = await api.get<{ schedule?: Record<string, unknown>; message?: string }>('/schedules/active')
    return {
      schedule: response.data.schedule ? transformSchedule(response.data.schedule) : undefined,
      message: response.data.message,
    }
  },
}

export default scheduleService
