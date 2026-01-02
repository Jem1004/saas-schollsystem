// Export and Monthly Recap types
// Requirements: 1.2, 2.1 - Export filter and monthly recap types

/**
 * Filter options for exporting attendance data
 * Requirements: 1.2, 1.3 - Allow filtering by date range and class
 */
export interface ExportFilter {
  startDate: string // Format: YYYY-MM-DD
  endDate: string   // Format: YYYY-MM-DD
  classId?: number
}

/**
 * Filter options for monthly recap
 * Requirements: 2.3, 2.4 - Allow filtering by month, year, and class
 */
export interface MonthlyRecapFilter {
  month: number  // 1-12
  year: number   // e.g., 2024
  classId?: number
}

/**
 * Student attendance summary for monthly recap
 * Requirements: 2.1, 2.2 - Summary per student with attendance percentage
 */
export interface StudentRecapSummary {
  student_id: number
  student_nis: string
  student_nisn: string
  student_name: string
  class_name: string
  total_present: number
  total_late: number
  total_very_late: number
  total_absent: number
  attendance_percent: number // (present / total_days) * 100
}

/**
 * Monthly recap response data
 * Requirements: 2.1 - Display summary per student including total days present, late, very late, and absent
 */
export interface MonthlyRecapResponse {
  month: number
  year: number
  total_days: number       // Total school days in the month
  class_id?: number
  class_name?: string
  student_recaps: StudentRecapSummary[]
}

/**
 * Export attendance record for display
 * Requirements: 1.4, 1.5 - Include student info and attendance details
 */
export interface ExportAttendanceRecord {
  studentNis: string
  studentNisn: string
  studentName: string
  className: string
  date: string
  checkInTime: string
  checkOutTime: string
  status: string
  scheduleName?: string
}

// Month options for selector
export const MONTH_OPTIONS = [
  { value: 1, label: 'Januari' },
  { value: 2, label: 'Februari' },
  { value: 3, label: 'Maret' },
  { value: 4, label: 'April' },
  { value: 5, label: 'Mei' },
  { value: 6, label: 'Juni' },
  { value: 7, label: 'Juli' },
  { value: 8, label: 'Agustus' },
  { value: 9, label: 'September' },
  { value: 10, label: 'Oktober' },
  { value: 11, label: 'November' },
  { value: 12, label: 'Desember' },
] as const

// Helper function to get month name
export function getMonthName(month: number): string {
  const monthOption = MONTH_OPTIONS.find(m => m.value === month)
  return monthOption?.label ?? ''
}

// Generate year options (current year and 5 years back)
export function getYearOptions(): { value: number; label: string }[] {
  const currentYear = new Date().getFullYear()
  const years: { value: number; label: string }[] = []
  for (let i = 0; i <= 5; i++) {
    const year = currentYear - i
    years.push({ value: year, label: year.toString() })
  }
  return years
}
