import api from './api'
import type { ExportFilter, MonthlyRecapFilter, MonthlyRecapResponse } from '@/types/export'

/**
 * Export service for attendance data export and monthly recap
 * Requirements: 1.1, 2.5 - Export attendance and monthly recap to Excel
 */
export const exportService = {
  /**
   * Export attendance data to Excel file
   * Requirements: 1.1 - Generate Excel file (.xlsx) containing attendance records
   * Requirements: 1.2, 1.3 - Allow filtering by date range and class
   */
  async exportAttendance(filter: ExportFilter): Promise<void> {
    const params: Record<string, string | number> = {
      start_date: filter.startDate,
      end_date: filter.endDate,
    }
    if (filter.classId) {
      params.class_id = filter.classId
    }

    const response = await api.get('/attendance/export', {
      params,
      responseType: 'blob',
    })

    // Extract filename from Content-Disposition header or generate default
    const contentDisposition = response.headers['content-disposition']
    let filename = `attendance_export_${filter.startDate}_to_${filter.endDate}.xlsx`
    
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1].replace(/['"]/g, '')
      }
    }

    // Trigger file download
    downloadBlob(response.data, filename)
  },

  /**
   * Get monthly recap data
   * Requirements: 2.1 - Display summary per student including total days present, late, very late, and absent
   */
  async getMonthlyRecap(filter: MonthlyRecapFilter): Promise<MonthlyRecapResponse> {
    const params: Record<string, number> = {
      month: filter.month,
      year: filter.year,
    }
    if (filter.classId) {
      params.class_id = filter.classId
    }

    const response = await api.get<MonthlyRecapResponse>('/attendance/monthly-recap', { params })
    return response.data
  },

  /**
   * Export monthly recap to Excel file
   * Requirements: 2.5 - Export monthly recap to Excel with summary statistics
   */
  async exportMonthlyRecap(filter: MonthlyRecapFilter): Promise<void> {
    const params: Record<string, number> = {
      month: filter.month,
      year: filter.year,
    }
    if (filter.classId) {
      params.class_id = filter.classId
    }

    const response = await api.get('/attendance/monthly-recap/export', {
      params,
      responseType: 'blob',
    })

    // Extract filename from Content-Disposition header or generate default
    const contentDisposition = response.headers['content-disposition']
    const monthNames = ['', 'januari', 'februari', 'maret', 'april', 'mei', 'juni', 
                        'juli', 'agustus', 'september', 'oktober', 'november', 'desember']
    let filename = `rekap_bulanan_${monthNames[filter.month]}_${filter.year}.xlsx`
    
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1].replace(/['"]/g, '')
      }
    }

    // Trigger file download
    downloadBlob(response.data, filename)
  },
}

/**
 * Helper function to trigger file download from blob
 */
function downloadBlob(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  
  // Cleanup
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}

export default exportService
