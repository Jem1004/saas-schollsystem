export * from './common'
export * from './user'
export * from './tenant'
export * from './device'
export * from './school'
// Note: student.ts types are more detailed versions, use school.ts Student for basic usage
// Export specific types from student.ts that don't conflict
export type { StudentProfile, StudentSearchResult, ImportStudentsRequest, ImportStudentsResponse, StudentSearchResponse } from './student'
// Note: attendance.ts has more detailed types than school.ts AttendanceSummary
export type { 
  Attendance, 
  AttendanceStatus, 
  AttendanceMethod,
  AttendanceRecord,
  DailyAttendanceSummary,
  ClassAttendanceSummary,
  StudentAttendanceHistory,
  RecordRFIDAttendanceRequest,
  RecordManualAttendanceRequest,
  UpdateAttendanceRequest,
  GetAttendanceParams,
  AttendanceSummaryResponse,
  StudentAttendanceResponse,
  ATTENDANCE_STATUS_CONFIG
} from './attendance'
// Note: grade.ts has more detailed types than homeroom.ts Grade
export type { GradeSummary, GetGradesParams, StudentGradesResponse } from './grade'
export * from './bk'
export * from './homeroom'
export * from './notification'
// Note: settings.ts has more detailed types than school.ts SchoolSettings
export type { AttendanceTimeWindow, DEFAULT_SCHOOL_SETTINGS, SEMESTER_OPTIONS } from './settings'
