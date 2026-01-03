import api from './api'
import type {
  Attendance,
  DailyAttendanceSummary,
  ClassAttendanceSummary,
  StudentAttendanceHistory,
  RecordRFIDAttendanceRequest,
  RecordManualAttendanceRequest,
  UpdateAttendanceRequest,
  GetAttendanceParams,
  AttendanceListResponse,
  AttendanceSummaryResponse,
  StudentAttendanceResponse,
} from '@/types/attendance'

export const attendanceService = {
  // Get attendance list
  async getAttendance(params?: GetAttendanceParams): Promise<AttendanceListResponse> {
    const response = await api.get<{ success: boolean; data: AttendanceListResponse }>('/attendance', { params })
    return response.data.data
  },

  // Get single attendance record
  async getAttendanceById(id: number): Promise<Attendance> {
    const response = await api.get<Attendance>(`/attendance/${id}`)
    return response.data
  },

  // Get daily attendance summary
  async getDailySummary(date: string): Promise<AttendanceSummaryResponse> {
    const response = await api.get<AttendanceSummaryResponse>('/attendance/summary', {
      params: { date },
    })
    return response.data
  },

  // Get attendance summary by date range
  async getSummaryByRange(startDate: string, endDate: string): Promise<DailyAttendanceSummary[]> {
    const response = await api.get<DailyAttendanceSummary[]>('/attendance/summary/range', {
      params: { startDate, endDate },
    })
    return response.data
  },

  // Get class attendance for a specific date
  async getClassAttendance(classId: number, date: string): Promise<AttendanceListResponse> {
    const response = await api.get<AttendanceListResponse>(`/attendance/class/${classId}`, {
      params: { date },
    })
    return response.data
  },

  // Get class attendance summary
  async getClassAttendanceSummary(classId: number, date: string): Promise<ClassAttendanceSummary> {
    const response = await api.get<ClassAttendanceSummary>(`/attendance/class/${classId}/summary`, {
      params: { date },
    })
    return response.data
  },

  // Get student attendance history
  async getStudentAttendance(
    studentId: number,
    params?: GetAttendanceParams
  ): Promise<StudentAttendanceResponse> {
    const response = await api.get<StudentAttendanceResponse>(
      `/attendance/student/${studentId}`,
      { params }
    )
    return response.data
  },

  // Get student attendance history with summary
  async getStudentAttendanceHistory(
    studentId: number,
    startDate: string,
    endDate: string
  ): Promise<StudentAttendanceHistory> {
    const response = await api.get<StudentAttendanceHistory>(
      `/attendance/student/${studentId}/history`,
      { params: { startDate, endDate } }
    )
    return response.data
  },

  // Record RFID attendance (from device)
  async recordRFIDAttendance(data: RecordRFIDAttendanceRequest): Promise<Attendance> {
    const response = await api.post<Attendance>('/attendance/rfid', data)
    return response.data
  },

  // Record manual attendance
  async recordManualAttendance(data: RecordManualAttendanceRequest): Promise<Attendance> {
    const response = await api.post<Attendance>('/attendance/manual', data)
    return response.data
  },

  // Update attendance record
  async updateAttendance(id: number, data: UpdateAttendanceRequest): Promise<Attendance> {
    const response = await api.put<Attendance>(`/attendance/${id}`, data)
    return response.data
  },

  // Delete attendance record
  async deleteAttendance(id: number): Promise<void> {
    await api.delete(`/attendance/${id}`)
  },

  // Export attendance report
  async exportAttendance(params: {
    startDate: string
    endDate: string
    classId?: number
    format?: 'xlsx' | 'csv' | 'pdf'
  }): Promise<Blob> {
    const response = await api.get('/attendance/export', {
      params,
      responseType: 'blob',
    })
    return response.data
  },

  // Get today's attendance status for a student
  async getTodayStatus(studentId: number): Promise<Attendance | null> {
    const response = await api.get<Attendance | null>(`/attendance/student/${studentId}/today`)
    return response.data
  },
}

export default attendanceService
